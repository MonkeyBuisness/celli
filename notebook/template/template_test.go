package template

import (
	"testing"

	"github.com/MonkeyBuisness/celli/notebook/types"
	"github.com/stretchr/testify/require"
)

func Test_NewBookTemplate(t *testing.T) {
	t.Run("get book settings error", func(t *testing.T) {
		tmpBooksData := make([]byte, len(bookSettingsData))
		copy(tmpBooksData, bookSettingsData)
		bookSettingsData = []byte{}
		defer func() {
			bookSettingsData = tmpBooksData
		}()

		_, err := NewBookTemplate(types.BookTypeJavaBook)
		require.EqualError(t, err, "could not read book settings: unexpected end of JSON input")
	})
	t.Run("find book settings error", func(t *testing.T) {
		_, err := NewBookTemplate(types.BookType("invalid"))
		require.EqualError(t, err, "could not find book settings for type invalid")
	})
	t.Run("all ok", func(t *testing.T) {
		data, err := NewBookTemplate(types.BookTypeJavaBook)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		const expData = `# _Put name of the book here..._

<!-- notebook:{ 
	"version": "1.0"
} -->

## Start your book here

<!-- br: -->

# Code Example

<!-- code:{ 
	"lang": "java",
	"content": "package main;

public class Main {
	public static void main(String[] args) {
		System.out.println(\"Hello World!\");
	}
}
",
	"meta": {}
} -->

## License

MIT

<!-- author:[`
		require.Contains(t, string(data), expData)
	})
}
