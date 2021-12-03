package converter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	e "github.com/MonkeyBuisness/cellementary-cli/notebook/errors"
	"github.com/MonkeyBuisness/cellementary-cli/notebook/serializer/comments"
	"github.com/MonkeyBuisness/cellementary-cli/notebook/types"
)

// Proceed converts notebook to the template data.
//
// Not it only supports `br:`, `code:{}` and `notebook:{}` type of serializable comments.
func Proceed(source io.Reader) ([]byte, error) {
	// read notebook content.
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(source); err != nil {
		return nil, e.ErrReadNotebookSource.New(err.Error())
	}

	// parse notebook data.
	var notebook types.NotebookData
	if err := json.Unmarshal(buf.Bytes(), &notebook); err != nil {
		return nil, e.ErrParseNotebookContent.New(err.Error())
	}

	// create template based on notebook data.
	return createTemplateData(&notebook)
}

func createTemplateData(notebook *types.NotebookData) ([]byte, error) {
	buf := make([]byte, 0, len(notebook.Cells))

	// convert notebook metadata.
	if len(notebook.Metadata) != 0 {
		buf = append(buf, createMetadataComment(notebook.Metadata)...)
	}

	// convert cells meatadata.
	for i := range notebook.Cells {
		c := &notebook.Cells[i]

		if c.Kind == types.NotebookCellKindMarkup {
			buf = append(buf, createMarkupComment(c)...)
			continue
		}

		codeComment, err := createCodeComment(c)
		if err != nil {
			return nil, e.ErrCreateTemplateContent.New(err.Error())
		}
		buf = append(buf, codeComment...)
	}

	return buf, nil
}

func createMetadataComment(meta map[string]interface{}) []byte {
	return []byte(fmt.Sprintf("%s\n", comments.NewNotebook(meta)))
}

func createMarkupComment(cell *types.NotebookCellData) []byte {
	return []byte(fmt.Sprintf("%s\n\n%s\n\n", cell.Content, comments.NewBr()))
}

func createCodeComment(cell *types.NotebookCellData) ([]byte, error) {
	codeComment, err := comments.NewCode(cell)
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("\n\n%s\n\n", string(codeComment))), nil
}
