package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/MonkeyBuisness/celli/notebook/converter"
	"github.com/MonkeyBuisness/celli/notebook/serializer"
	"github.com/MonkeyBuisness/celli/notebook/serializer/comments"
	"github.com/MonkeyBuisness/celli/notebook/template"
	"github.com/MonkeyBuisness/celli/notebook/types"
	"github.com/MonkeyBuisness/celli/notebook/utils"
)

const (
	defaultTemplateFileName = "template.md"
)

// CreateTemplate creates a new template based on the type.
func CreateTemplate(bookType, dest string) error {
	templateData, err := template.NewBookTemplate(types.BookType(strings.ToLower(bookType)))
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(dest)
	if err != nil {
		return fmt.Errorf("could not resolve template file path: %v", err)
	}

	if fileInfo.IsDir() {
		dest = path.Join(dest, defaultTemplateFileName)
	}

	file, err := os.OpenFile(filepath.Clean(dest), os.O_CREATE|os.O_WRONLY, types.DefaultFileMode)
	if err != nil {
		return fmt.Errorf("could not open template file: %v", err)
	}
	defer utils.Close(file)

	if _, err := file.Write(templateData); err != nil {
		return fmt.Errorf("could not write data to the template file: %v", err)
	}

	return nil
}

// ConvertToTemplate converts notebook file to the template implementation.
func ConvertToTemplate(notebookPath string) error {
	file, err := os.OpenFile(filepath.Clean(notebookPath), os.O_RDONLY, types.DefaultFileMode)
	if err != nil {
		return fmt.Errorf("could not open notebook file: %v", err)
	}
	defer utils.Close(file)

	data, err := converter.Proceed(file)
	if err != nil {
		return fmt.Errorf("could not convert notebook data: %v", err)
	}

	if _, err := os.Stdout.Write(data); err != nil {
		return err
	}

	return nil
}

// ConvertToNotebook converts template file to the notebook implementation.
func ConvertToNotebook(templatePath string, pretty bool) error {
	file, err := os.OpenFile(filepath.Clean(templatePath), os.O_RDONLY, types.DefaultFileMode)
	if err != nil {
		return fmt.Errorf("could not open notebook file: %v", err)
	}
	defer utils.Close(file)

	s := serializer.New()
	notebookData, err := s.SerializeNotebook(file,
		serializer.WithCommentSerializer(
			comments.NewCodeCommentSerializer(),
			comments.NewBrCommentSerializer(),
			comments.NewNotebookCommentSerializer(),
			comments.NewAuthorCommentSerializer(),
		),
	)
	if err != nil {
		return fmt.Errorf("could not serialize notebook data: %v", err)
	}

	data, err := json.Marshal(notebookData)
	if err != nil {
		return err
	}

	if pretty {
		var buf bytes.Buffer
		if err := json.Indent(&buf, data, "", "\t"); err != nil {
			return err
		}
		data = buf.Bytes()
	}

	if _, err := os.Stdout.Write(data); err != nil {
		return err
	}

	return nil
}
