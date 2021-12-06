package comments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/MonkeyBuisness/celli/notebook/types"
	"github.com/MonkeyBuisness/celli/notebook/utils"
)

const (
	filePrefix = "file://"
)

// CodeCommentSerializer represents <!-- code:{...} --> comment serializer.
type CodeCommentSerializer struct{}

type codeCommentPayload struct {
	LanguageID string                 `json:"lang"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
	Content    string                 `json:"content,omitempty"`
	URI        string                 `json:"uri,omitempty"`
}

// NewCodeCommentSerializer returns new CodeCommentSerializer instance.
func NewCodeCommentSerializer() CodeCommentSerializer {
	return CodeCommentSerializer{}
}

// Key returns the name of the serializable comment key.
func (s CodeCommentSerializer) Key() string {
	return "code"
}

// Render renders serializer data to the notebook.
func (s CodeCommentSerializer) Render(notebook *types.NotebookData, payload []byte) error {
	var code codeCommentPayload
	if err := json.Unmarshal(payload, &code); err != nil {
		return err
	}

	if code.URI != "" {
		content, err := readURIContent(code.URI)
		if err != nil {
			return fmt.Errorf("could not read URI content: %v", err)
		}
		code.Content = string(content)
	}

	notebook.Cells = append(notebook.Cells, types.NotebookCellData{
		LanguageID: code.LanguageID,
		Content:    code.Content,
		Kind:       types.NotebookCellKindCode,
		Metadata:   code.Meta,
	})

	return nil
}

func readURIContent(uri string) ([]byte, error) {
	isFile := strings.HasPrefix(uri, filePrefix)
	if isFile {
		filePath := strings.TrimPrefix(uri, filePrefix)

		data, err := os.ReadFile(filepath.Clean(filePath))
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	//nolint:gosec,gocritic,bodyclose // itэs entirely the responsibility of the author.
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer utils.Close(resp.Body)

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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
