package comments

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MonkeyBuisness/cellementary-cli/notebook/types"
)

// NotebookCommentSerializer represents <!-- notebook:{...} --> comment serializer.
type NotebookCommentSerializer struct {
	payload notebookCommentPayload
}

type notebookCommentPayload map[string]interface{}

// NewNotebookCommentSerializer returns new NotebookCommentSerializer instance.
func NewNotebookCommentSerializer() NotebookCommentSerializer {
	return NotebookCommentSerializer{
		payload: make(notebookCommentPayload),
	}
}

// Key returns the name of the serializable comment key.
func (s NotebookCommentSerializer) Key() string {
	return "notebook"
}

// Render renders serializer data to the notebook.
func (s NotebookCommentSerializer) Render(notebook *types.NotebookData) error {
	for key, value := range s.payload {
		notebook.Metadata[key] = value
	}

	return nil
}

// SetPayload sets payload data to the serializer.
func (s NotebookCommentSerializer) SetPayload(data []byte) error {
	return json.Unmarshal(data, &s.payload)
}

// NewNotebook creates new <!-- notebook:{} --> comment string.
func NewNotebook(meta map[string]interface{}) string {
	var metaFields string
	for key, value := range meta {
		metaFields = fmt.Sprintf("%s\t\"%s\": \"%v\",\n", metaFields, key, value)
	}
	metaFields = strings.TrimSuffix(metaFields, "\n")
	metaFields = metaFields[:len(metaFields)-1]

	return fmt.Sprintf("<!-- notebook:{\n%s\n} -->", metaFields)
}
