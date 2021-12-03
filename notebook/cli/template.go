package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/MonkeyBuisness/cellementary-cli/notebook/converter"
	"github.com/MonkeyBuisness/cellementary-cli/notebook/serializer"
	"github.com/MonkeyBuisness/cellementary-cli/notebook/serializer/comments"
	"github.com/MonkeyBuisness/cellementary-cli/notebook/template"
	"github.com/MonkeyBuisness/cellementary-cli/notebook/types"
)

const (
	defaultTemplateFileName = "template.md"
)

// CreateTemplate creates a new template based on the type.
func CreateTemplate(bookType string, dest string) error {
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

	file, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("could not open template file: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(templateData); err != nil {
		return fmt.Errorf("could not write data to the template file: %v", err)
	}

	return nil
}

// ConvertToTemplate converts notebook file to the template implementation.
func ConvertToTemplate(notebookPath string) error {
	file, err := os.OpenFile(notebookPath, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("could not open notebook file: %v", err)
	}
	defer file.Close()

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
func ConvertToNotebook(templatePath string) error {
	file, err := os.OpenFile(templatePath, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("could not open notebook file: %v", err)
	}
	defer file.Close()

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

	if err := json.NewEncoder(os.Stdout).Encode(notebookData); err != nil {
		return err
	}

	return nil
}
