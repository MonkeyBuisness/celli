package comments

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/MonkeyBuisness/celli/notebook/types"
)

//go:embed authors.tpl.md
var authorsTemplate embed.FS

// AuthorCommentSerializer represents <!-- author:[...] --> comment serializer.
type AuthorCommentSerializer struct{}

type authorCommentPayload struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
	Link   string `json:"link,omitempty"`
	About  string `json:"about,omitempty"`
}

// NewAuthorCommentSerializer returns new AuthorCommentSerializer instance.
func NewAuthorCommentSerializer() *AuthorCommentSerializer {
	return &AuthorCommentSerializer{}
}

// Key returns the name of the serializable comment key.
func (s *AuthorCommentSerializer) Key() string {
	return "author"
}

// Render renders serializer data to the notebook.
func (s *AuthorCommentSerializer) Render(notebook *types.NotebookData, payload []byte) error {
	var authors []authorCommentPayload
	if err := json.Unmarshal(payload, &authors); err != nil {
		return err
	}

	t, err := template.ParseFS(authorsTemplate, "*.tpl.md")
	if err != nil {
		return fmt.Errorf("could not parse authors template: %v", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, authors); err != nil {
		return fmt.Errorf("could not execute authors template: %v", err)
	}

	notebook.Cells = append(notebook.Cells, types.NotebookCellData{
		LanguageID: types.MarkdownLanguageID,
		Content:    buf.String(),
		Kind:       types.NotebookCellKindMarkup,
	})

	return nil
}
