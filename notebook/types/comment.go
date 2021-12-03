package types

// SerializableComment represents API to work with serializable comments in the document.
type SerializableComment interface {
	Key() string
	Render(notebook *NotebookData) error
	SetPayload(data []byte) error
}
