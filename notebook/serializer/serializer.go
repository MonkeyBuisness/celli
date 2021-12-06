package serializer

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	e "github.com/MonkeyBuisness/cellementary-cli/notebook/errors"
	"github.com/MonkeyBuisness/cellementary-cli/notebook/types"
	"github.com/sirupsen/logrus"
)

const (
	nodeKindText    nodeKind = iota
	nodeKindComment nodeKind = iota
)

const (
	subExpCommentKey     = "key"
	subExpCommentPayload = "payload"
)

var (
	commentRegexp     = regexp.MustCompile(`(<!--[\s\S]*?-->)`)
	commentMetaRegexp = regexp.MustCompile(
		fmt.Sprintf(`(?P<%s>\w*[\s\S]):(?P<%s>[{|\[]*[\s\S]*[}|\]])?`,
			subExpCommentKey, subExpCommentPayload),
	)
)

// Option represents serializer option model.
type Option func(*Options)

// Options represents serializer configuration model.
type Options struct {
	serializers map[string]types.SerializableComment
}

// Serializer represents notebook serializer implementation.
type Serializer struct{}

type documentNode interface {
	nodeKind() nodeKind
	render(notebook *types.NotebookData) error
}

type baseNode struct {
	start int
	end   int
	kind  nodeKind
}

type nodeKind int

type commentNode struct {
	*baseNode

	serializer types.SerializableComment
}

type textNode struct {
	*baseNode

	content string
}

// New returns new notebook Serializer instance.
func New() Serializer {
	return Serializer{}
}

// SerializeNotebook converts markup text to the notebook data implementation.
func (s *Serializer) SerializeNotebook(
	source io.Reader, opt ...Option) (*types.NotebookData, error) {
	// apply incoming options.
	opts := Options{
		serializers: make(map[string]types.SerializableComment),
	}
	for _, o := range opt {
		o(&opts)
	}

	// read markdown content.
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(source); err != nil {
		return nil, e.ErrReadMarkdownSource.New(err.Error())
	}
	content := buf.String()

	// parse markup content.
	nodes, err := s.parseMarkupContent(content, &opts)
	if err != nil {
		return nil, err
	}

	// render nodes to the notebook document data.
	return s.renderNotebook(nodes)
}

func (s *Serializer) parseMarkupContent(content string, opts *Options) ([]documentNode, error) {
	// find all HTML comment blocks inside the document.
	commentIndices := commentRegexp.FindAllStringIndex(content, -1)
	if len(commentIndices) == 0 {
		// insert "fake" comment to the start.
		commentIndices = [][]int{{}}
	}
	// insert a "fake" comment to the end.
	commentIndices = append(commentIndices, []int{
		len(content),
		len(content),
	})

	nodes := make([]documentNode, 0)
	// parse each comment to extract meta value.
	var prevCommentPos int
	expKeyIndex := commentMetaRegexp.SubexpIndex(subExpCommentKey)
	expPayloadIndex := commentMetaRegexp.SubexpIndex(subExpCommentPayload)
	for i := range commentIndices {
		iStart, iEnd := commentIndices[i][0], commentIndices[i][1]

		// check if document contains text node before comment.
		if textNodeLen := iStart - prevCommentPos; textNodeLen > 1 {
			tNode := newTextNode(prevCommentPos+1, iStart-1, content)
			nodes = append(nodes, tNode)
		}

		prevCommentPos = iEnd
		if iStart == iEnd {
			continue
		}

		metaIndices := commentMetaRegexp.FindStringSubmatch(content[iStart:iEnd])

		// check if comment is a not serializable comment.
		if len(metaIndices) == 0 {
			tNode := newTextNode(iStart, iEnd, content)
			nodes = append(nodes, tNode)
			continue
		}

		// detect comment key.
		cKey := metaIndices[expKeyIndex]
		serializer, ok := opts.serializers[cKey]
		if !ok {
			logrus.Warnf("could not serialize comment at position %d:%d: unknown key %s",
				iStart, iEnd, cKey)
			tNode := newTextNode(iStart, iEnd, content)
			nodes = append(nodes, tNode)
			continue
		}

		// detect comment payload.
		cPayload := metaIndices[expPayloadIndex]
		if err := serializer.SetPayload([]byte(cPayload)); err != nil {
			return nil, e.ErrMarshalCommentPayload.New(err.Error(),
				fmt.Sprintf("at position %d:%d", iStart, iEnd))
		}

		// create comment node.
		cNode := commentNode{
			baseNode: &baseNode{
				start: iStart,
				end:   iEnd,
				kind:  nodeKindComment,
			},
			serializer: serializer,
		}
		nodes = append(nodes, cNode)
	}

	return optimizeNodes(nodes), nil
}

func (s *Serializer) renderNotebook(nodes []documentNode) (*types.NotebookData, error) {
	notebookData := types.NotebookData{
		Cells:    make([]types.NotebookCellData, 0, len(nodes)),
		Metadata: make(map[string]interface{}),
	}

	for i := range nodes {
		n := nodes[i]

		if err := n.render(&notebookData); err != nil {
			return nil, e.ErrRenderNotebook.New(err.Error())
		}
	}

	return &notebookData, nil
}

func (n *baseNode) nodeKind() nodeKind {
	return n.kind
}

func (n textNode) render(notebook *types.NotebookData) error {
	content := strings.TrimSpace(n.content)
	if content == "" {
		return nil
	}

	notebook.Cells = append(notebook.Cells, types.NotebookCellData{
		LanguageID: types.MarkdownLanguageID,
		Content:    content,
		Kind:       types.NotebookCellKindMarkup,
	})

	return nil
}

func (n commentNode) render(notebook *types.NotebookData) error {
	return n.serializer.Render(notebook)
}

// WithCommentSerializer adds a new comment serializer.
func WithCommentSerializer(s ...types.SerializableComment) Option {
	return func(o *Options) {
		for i := range s {
			o.serializers[s[i].Key()] = s[i]
		}
	}
}

func newTextNode(start, end int, content string) textNode {
	return textNode{
		baseNode: &baseNode{
			start: start,
			end:   end,
			kind:  nodeKindText,
		},
		content: content[start:end],
	}
}

func optimizeNodes(nodes []documentNode) []documentNode {
	optimizedNodes := make([]documentNode, 0, len(nodes))

	var (
		prevNode          documentNode
		prevTextNodeIndex int
	)
	for i := range nodes {
		n := nodes[i]
		if prevNode != nil && n.nodeKind() == nodeKindText && prevNode.nodeKind() == nodeKindText {
			tPrevNode := prevNode.(textNode)
			tNode := n.(textNode)
			tPrevNode.content = fmt.Sprintf("%s\n%s", tPrevNode.content, tNode.content)
			tPrevNode.end = tNode.end
			optimizedNodes[prevTextNodeIndex] = tPrevNode
			prevNode = optimizedNodes[prevTextNodeIndex]
			continue
		}

		optimizedNodes = append(optimizedNodes, n)
		prevNode = n
		if n.nodeKind() == nodeKindText {
			prevTextNodeIndex = len(optimizedNodes) - 1
		}
	}

	return optimizedNodes
}
