package comments

import (
	"github.com/MonkeyBuisness/cellementary-cli/notebook/types"
)

// BrCommentSerializer represents <!-- br: --> comment serializer.
type BrCommentSerializer struct{}

// NewBrCommentSerializer returns new BrCommentSerializer instance.
func NewBrCommentSerializer() BrCommentSerializer {
	return BrCommentSerializer{}
}

// Key returns the name of the serializable comment key.
func (s BrCommentSerializer) Key() string {
	return "br"
}

// Render renders serializer data to the notebook.
func (s BrCommentSerializer) Render(notebook *types.NotebookData) error {
	return nil
}

// SetPayload sets payload data to the serializer.
func (s BrCommentSerializer) SetPayload(data []byte) error {
	return nil
}

// NewBr creates new <!-- br: --> comment string.
func NewBr() string {
	return "<!-- br: -->"
}
