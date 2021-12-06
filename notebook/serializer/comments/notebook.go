package comments

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MonkeyBuisness/celli/notebook/types"
)

// NotebookCommentSerializer represents <!-- notebook:{...} --> comment serializer.
type NotebookCommentSerializer struct{}

// NewNotebookCommentSerializer returns new NotebookCommentSerializer instance.
func NewNotebookCommentSerializer() NotebookCommentSerializer {
	return NotebookCommentSerializer{}
}

// Key returns the name of the serializable comment key.
func (s NotebookCommentSerializer) Key() string {
	return "notebook"
}

// Render renders serializer data to the notebook.
func (s NotebookCommentSerializer) Render(notebook *types.NotebookData, payload []byte) error {
	var meta map[string]interface{}
	if err := json.Unmarshal(payload, &meta); err != nil {
		return err
	}

	for key, value := range meta {
		notebook.Metadata[key] = value
	}

	return nil
}

// NewNotebook creates new <!-- notebook:{} --> comment string.
func NewNotebook(meta map[string]interface{}) string {
	var metaFields string
	for key, value := range meta {
		metaFields = fmt.Sprintf("%s\t%q: \"%v\",\n", metaFields, key, value)
	}
	metaFields = strings.TrimSuffix(metaFields, "\n")
	metaFields = metaFields[:len(metaFields)-1]

	return fmt.Sprintf("<!-- notebook:{\n%s\n} -->", metaFields)
}
