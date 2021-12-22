package comments

import (
	"fmt"

	"github.com/MonkeyBuisness/celli/notebook/types"
	"gopkg.in/yaml.v2"
)

// YCodeCommentSerializer represents <!-- ycode:{...} --> comment serializer.
type YCodeCommentSerializer struct{}

type ycodeCommentPayload struct {
	LanguageID string                 `yaml:"lang"`
	Content    string                 `yaml:"code,omitempty,flow"`
	URI        string                 `yaml:"uri,omitempty"`
	Meta       map[string]interface{} `yaml:"meta,omitempty"`
}

// NewYCodeCommentSerializer returns new YCodeCommentSerializer instance.
func NewYCodeCommentSerializer() YCodeCommentSerializer {
	return YCodeCommentSerializer{}
}

// Key returns the name of the serializable comment key.
func (s YCodeCommentSerializer) Key() string {
	return "ycode"
}

// Render renders serializer data to the notebook.
func (s YCodeCommentSerializer) Render(notebook *types.NotebookData, payload []byte) error {
	var code ycodeCommentPayload
	if err := yaml.Unmarshal(payload[1:len(payload)-1], &code); err != nil {
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
