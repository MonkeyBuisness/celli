package types

import "io/fs"

// DefaultFileMode contains default file mode type.
const DefaultFileMode fs.FileMode = 0o666

// SerializableComment represents API to work with serializable comments in the document.
type SerializableComment interface {
	Key() string
	Render(notebook *NotebookData) error
	SetPayload(data []byte) error
}
