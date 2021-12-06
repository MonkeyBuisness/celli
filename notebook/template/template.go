package template

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/MonkeyBuisness/cellementary-cli/notebook/types"
	"github.com/tcnksm/go-gitconfig"
)

const (
	bookTemplateFile = "book.tpl.md"
)

//go:embed book.tpl.md
var bookTemplate embed.FS

//go:embed book.settings.json
var bookSettingsData []byte

type bookSettings struct {
	Notebook json.RawMessage `json:"notebook,omitempty"`
	Code     json.RawMessage `json:"code,omitempty"`
	Author   json.RawMessage `json:"author,omitempty"`
}

// NewBookTemplate creates a book template according to the provided book type.
func NewBookTemplate(bookType types.BookType) ([]byte, error) {
	t := template.Must(template.
		New(bookTemplateFile).
		Funcs(defaultTemplateFuncs()).
		ParseFS(bookTemplate, "*.tpl.md"),
	)

	var booksSettings map[types.BookType]bookSettings
	if err := getBookSettings(booksSettings); err != nil {
		return nil, fmt.Errorf("could not read book settings: %v", err)
	}

	bookSettings, ok := booksSettings[bookType]
	if !ok {
		return nil, fmt.Errorf("could not find book settings for type %s", bookType)
	}

	var buf bytes.Buffer
	if err := t.Funcs(defaultTemplateFuncs()).Execute(&buf, bookSettings); err != nil {
		return nil, fmt.Errorf("could not execute book template: %v", err)
	}

	return buf.Bytes(), nil
}

func getBookSettings(settings map[types.BookType]bookSettings) error {
	return json.Unmarshal(bookSettingsData, &settings)
}

func defaultTemplateFuncs() template.FuncMap {
	return map[string]interface{}{
		"now": func() string {
			return time.Now().UTC().Format("02 Jan 2006")
		},
		"asJSON": func(data []uint8) string {
			var buf bytes.Buffer
			if err := json.Indent(&buf, data, "", "\t"); err != nil {
				return buf.String()
			}

			return strings.ReplaceAll(
				strings.ReplaceAll(
					strings.TrimSuffix(buf.String()[1:buf.Len()-2], "\n"), "\\n", "\n"), "\\t", "\t",
			)
		},
		"authorName": func() string {
			if githubUser, err := gitconfig.GithubUser(); err == nil {
				return githubUser
			}

			if username, err := gitconfig.Username(); err == nil {
				return username
			}

			return ""
		},
		"authorLink": func() string {
			if githubUser, err := gitconfig.GithubUser(); err == nil {
				return fmt.Sprintf("https://github.com/%s", githubUser)
			}

			if url, err := gitconfig.OriginURL(); err == nil {
				return url
			}

			return ""
		},
		"authorAvatar": func() string {
			if githubUser, err := gitconfig.GithubUser(); err == nil {
				return fmt.Sprintf("https://github.com/%s.png", githubUser)
			}

			return ""
		},
		"authorAbout": func() string {
			if email, err := gitconfig.Email(); err == nil {
				return fmt.Sprintf("Email: %s", email)
			}

			return ""
		},
	}
}
