package serializer

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	e "github.com/MonkeyBuisness/celli/notebook/errors"
	"github.com/MonkeyBuisness/celli/notebook/types"
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

	payload    []byte
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
	nodes := s.parseMarkupContent(content, &opts)

	// render nodes to the notebook document data.
	return s.renderNotebook(nodes)
}

func (s *Serializer) parseMarkupContent(content string, opts *Options) []documentNode {
	// find all HTML comment blocks inside the document.
	commentIndices := commentRegexp.FindAllStringIndex(content, -1)

	// split document into nodes.
	docNodes := splicDocIntoNodes(commentIndices, len(content))
	expKeyIndex := commentMetaRegexp.SubexpIndex(subExpCommentKey)
	expPayloadIndex := commentMetaRegexp.SubexpIndex(subExpCommentPayload)

	nodes := make([]documentNode, len(docNodes))
	for i := range docNodes {
		node := &docNodes[i]
		nodeContent := content[node.start:node.end]
		var docNode documentNode = newTextNode(node.start, node.end, nodeContent)

		switch node.kind {
		case nodeKindComment:
			// parse comment to extract meta value.
			metaIndices := commentMetaRegexp.FindStringSubmatch(nodeContent)

			// check if comment is a not serializable comment.
			if len(metaIndices) == 0 {
				break
			}

			// detect comment key.
			cKey := metaIndices[expKeyIndex]
			serializer, ok := opts.serializers[cKey]
			if !ok {
				logrus.Warnf("could not serialize comment at position %d:%d: unknown key %s",
					node.start, node.end, cKey)
				break
			}

			// detect comment payload.
			cPayload := metaIndices[expPayloadIndex]

			// create comment node.
			docNode = commentNode{
				baseNode: &baseNode{
					start: node.start,
					end:   node.end,
					kind:  nodeKindComment,
				},
				serializer: serializer,
				payload:    []byte(cPayload),
			}
		}

		nodes[i] = docNode
	}

	return optimizeNodes(nodes)
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
	return n.serializer.Render(notebook, n.payload)
}

// WithCommentSerializer adds a new comment serializer.
func WithCommentSerializer(s ...types.SerializableComment) Option {
	return func(o *Options) {
		for i := range s {
			o.serializers[s[i].Key()] = s[i]
		}
	}
}

func splicDocIntoNodes(commentIndices [][]int, docLen int) []baseNode {
	nodes := make([]baseNode, 0, len(commentIndices))

	commentByIndex := func(index int) (start, end int) {
		start, end = -1, -1

		for i := range commentIndices {
			cm := commentIndices[i]

			if index >= cm[0] && index <= cm[1] {
				start, end = cm[0], cm[1]
				return
			}
		}

		return
	}

	txtNode := baseNode{
		start: 0,
		end:   0,
		kind:  nodeKindText,
	}
	for i := 0; i < docLen; i++ {
		if commStart, commEnd := commentByIndex(i); commStart != -1 {
			if txtNode.end != 0 {
				nodes = append(nodes, txtNode)
			}
			nodes = append(nodes, baseNode{
				start: commStart,
				end:   commEnd,
				kind:  nodeKindComment,
			})
			i = commEnd
			txtNode.start = i + 1
			txtNode.end = 0
			continue
		}

		txtNode.end = i
	}

	if txtNode.end != 0 {
		nodes = append(nodes, txtNode)
	}

	return nodes
}

func newTextNode(start, end int, content string) textNode {
	return textNode{
		baseNode: &baseNode{
			start: start,
			end:   end,
			kind:  nodeKindText,
		},
		content: content,
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
