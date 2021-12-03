package comments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/MonkeyBuisness/cellementary-cli/notebook/types"
)

const (
	filePrefix = "file://"
)

// CodeCommentSerializer represents <!-- code:{...} --> comment serializer.
type CodeCommentSerializer struct {
	payload *codeCommentPayload
}

type codeCommentPayload struct {
	LanguageID string                 `json:"lang"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
	Content    string                 `json:"content,omitempty"`
	URI        string                 `json:"uri,omitempty"`
}

// NewCodeCommentSerializer returns new CodeCommentSerializer instance.
func NewCodeCommentSerializer() CodeCommentSerializer {
	return CodeCommentSerializer{
		payload: &codeCommentPayload{},
	}
}

// Key returns the name of the serializable comment key.
func (s CodeCommentSerializer) Key() string {
	return "code"
}

// Render renders serializer data to the notebook.
func (s CodeCommentSerializer) Render(notebook *types.NotebookData) error {
	if s.payload.URI != "" {
		if err := readURIContent(s.payload.URI, &s.payload.Content); err != nil {
			return fmt.Errorf("could not read URI content: %v", err)
		}
	}

	notebook.Cells = append(notebook.Cells, types.NotebookCellData{
		LanguageID: s.payload.LanguageID,
		Content:    s.payload.Content,
		Kind:       types.NotebookCellKindCode,
		Metadata:   s.payload.Meta,
	})

	return nil
}

// SetPayload sets payload data to the serializer.
func (s CodeCommentSerializer) SetPayload(data []byte) error {
	return json.Unmarshal(data, s.payload)
}

func readURIContent(uri string, content *string) error {
	var dataStr string
	defer func() {
		content = &dataStr
	}()

	isFile := strings.HasPrefix(uri, filePrefix)
	if isFile {
		filePath := strings.TrimPrefix(uri, filePrefix)

		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		dataStr = string(data)

		return nil
	}

	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return err
	}
	dataStr = buf.String()

	return nil
}

// NewCode creates new <!-- code:{} --> comment string.
func NewCode(cell *types.NotebookCellData) ([]byte, error) {
	data, err := json.Marshal(codeCommentPayload{
		LanguageID: cell.LanguageID,
		Meta:       cell.Metadata,
		Content:    cell.Content,
	})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, data, "", "\t"); err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("<!-- code:%s -->", buf.String())), nil
}
