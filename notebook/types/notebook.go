package types

// MarkdownLanguageID is an ID of the markup language.
const MarkdownLanguageID = "markdown"

// Notebook Cell Kind.
const (
	NotebookCellKindMarkup NotebookCellKind = 1
	NotebookCellKindCode   NotebookCellKind = 2
)

// Book type.
const (
	BookTypeJavaBook BookType = "javabook"
)

// NotebookData represents notebook data model.
type NotebookData struct {
	Cells    []NotebookCellData     `json:"cells,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NotebookCellData represents notebook cell data model.
type NotebookCellData struct {
	LanguageID string                 `json:"languageId"`
	Content    string                 `json:"content"`
	Kind       NotebookCellKind       `json:"kind"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// NotebookCellKind represents notebook cell kind.
//
// if Markup: cell content contains markdown source code that is used to display.
// If Code:   cell content contains source code that can be executed and that produces output.
type NotebookCellKind int

// BookType represents book type.
type BookType string

// SupportedBookTypes returns a slice of supported book type names.
func SupportedBookTypes() []string {
	return []string{
		string(BookTypeJavaBook),
	}
}
