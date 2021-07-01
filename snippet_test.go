package extensions

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
)

var markdownSnippet = goldmark.New(
	goldmark.WithExtensions(
		SnippetExtension(100),
	),
)

func TestSnippet(t *testing.T) {
	source := `# Hello
Paragraph text here.

Another one.
And continued here.`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdownSnippet.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		t.Error("Failed to convert markdown")
	}
	s, err := Snippet(context)
	if err != nil {
		t.Error("Failed to get snippet", err.Error())
	}
	if s != "Paragraph text here. Another one. And continued here." {
		t.Errorf("Snippet must be 'Paragraph text here. Another one. And continued here.', but got '%s'", s)
	}
	expected := "<h1>Hello</h1>\n<p>Paragraph text here.</p>\n<p>Another one.\nAnd continued here.</p>\n"
	if buf.String() != expected {
		t.Errorf("Should render '%s', but got '%s'", expected, buf.String())
	}
}

func TestEmptySnippet(t *testing.T) {
	source := `# Hello`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdownSnippet.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		t.Error("Failed to convert markdown")
	}
	s, err := Snippet(context)
	if err != nil {
		t.Error("Failed to get snippet")
	}
	if s != "" {
		t.Errorf("Snippet must be empty, but got '%s'", s)
	}
	if buf.String() != "<h1>Hello</h1>\n" {
		t.Errorf("Should render '<h1>Hello</h1>\n', but got '%s'", buf.String())
	}
}

func TestComplexSnippet(t *testing.T) {
	source := `# Heading
[Link](link) *italics*
![image](image)continued.

**Bold** [link.](link)
> Block quote.

    Code.
1. Ordered.
- Unordered
continued.

## Subheading

![Image](image)` + "`Inline.`"

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdownSnippet.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		t.Error("Failed to convert markdown")
	}
	s, err := Snippet(context)
	if err != nil {
		t.Error("Failed to get snippet")
	}
	if s != "Link italics continued. Bold link. Inline." {
		t.Errorf("Snippet must be 'Link italics continued. Bold link. Inline.', but got '%s'", s)
	}
}

func TestTruncatedSnippet(t *testing.T) {
	source := `# Hello
Paragraph text here.`
	markdownSnippet := goldmark.New(
		goldmark.WithExtensions(
			SnippetExtension(5),
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdownSnippet.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		t.Error("Failed to convert markdown")
	}
	s, err := Snippet(context)
	if err != nil {
		t.Error("Failed to get snippet")
	}
	if s != "Parag" {
		t.Errorf("Snippet must be 'Parag', but got '%s'", s)
	}
	expected := "<h1>Hello</h1>\n<p>Paragraph text here.</p>\n"
	if buf.String() != expected {
		t.Errorf("Should render '%s', but got '%s'", expected, buf.String())
	}
}

func TestSnippetError(t *testing.T) {
	source := `# Hello`
	markdownSnippet := goldmark.New()

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdownSnippet.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		t.Error("Failed to convert markdown")
	}
	if _, err := Snippet(context); err == nil {
		t.Error("Should throw because of missing extension")
	}
}

func TestSnippetInvalidString(t *testing.T) {
	source := `# Hello`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdownSnippet.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		t.Error("Failed to convert markdown")
	}
	context.Set(contextKeySnippet, 0) // Not the expected string
	if _, err := Snippet(context); err == nil {
		t.Error("Should throw because of invalid snippet string")
	}
}

func TestSnippetBuffer(t *testing.T) {
	b := newSnippetBuffer(10)
	if !b.IsEmpty() {
		t.Error("New snippet buffer should be empty")
	}
	b.Write([]byte(""))
	if !b.IsEmpty() {
		t.Error("Snippet buffer should be empty after writing empty string")
	}
	b.Write([]byte("12345"))
	if b.IsEmpty() || b.IsFull() {
		t.Error("Snippet buffer should be neither empty nor full")
	}
	b.Write([]byte("67890"))
	if !b.IsFull() {
		t.Error("Snippet buffer should be full")
	}
	b.Reset()
	if !b.IsEmpty() {
		t.Error("Snippet buffer should be empty after reset")
	}
}
