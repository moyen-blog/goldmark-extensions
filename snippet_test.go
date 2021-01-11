package meta

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
		t.Error("Failed to get snippet")
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

func TestTrucatedSnippet(t *testing.T) {
	source := `# Hello
Paragraph text here.`

	markdown := goldmark.New(
		goldmark.WithExtensions(
			SnippetExtension(5),
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
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

	var buf bytes.Buffer
	context := parser.NewContext()
	context.Set(contextKeySnippet, 0) // Not the expected string
	if err := markdownSnippet.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		t.Error("Failed to convert markdown")
	}
	if _, err := Snippet(context); err == nil {
		t.Error("Should throw snippet not available error")
	}
}
