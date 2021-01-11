// Package meta is a extension for the goldmark(http://github.com/yuin/goldmark).
//
// This extension parses YAML metadata blocks and store metadata to a parser.Context.
package meta

import (
	"errors"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var contextKeySnippet = parser.NewContextKey()

// Snippet returns the snippet made up of text from the underlying markdown document
// Only paragraph nodes are used to generate the snippet
func Snippet(pc parser.Context) (string, error) {
	if v := pc.Get(contextKeySnippet); v != nil {
		if s, ok := v.(string); ok {
			return s, nil
		}
		return "", errors.New("No snippet found in context")
	}
	return "", nil
}

type snippetParagraphTransformer struct {
	buf *snippetBuffer
}

func (p *snippetParagraphTransformer) Transform(node *ast.Paragraph, reader text.Reader, pc parser.Context) {
	lines := node.Lines()
	if lines.Len() == 0 || p.buf.IsFull() {
		return
	}
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		b := segment.Value(reader.Source())
		p.buf.Write(b)
	}
	pc.Set(contextKeySnippet, p.buf.String())
}

type snippetExtension struct {
	max int
}

// SnippetExtension returns a Goldmark extension
func SnippetExtension(max int) goldmark.Extender {
	return &snippetExtension{max}
}

func (e *snippetExtension) Extend(m goldmark.Markdown) {
	transformer := util.Prioritized(&snippetParagraphTransformer{
		newSnippetBuffer(e.max),
	}, 999)
	paragraphTransformers := append(parser.DefaultParagraphTransformers(), transformer)
	m.Parser().AddOptions(
		parser.WithParagraphTransformers(paragraphTransformers...),
	)
}
