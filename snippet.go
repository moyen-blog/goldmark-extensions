package meta

import (
	"errors"
	"io"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

// Snippet returns the snippet made up of text from the underlying markdown document
// Only paragraph nodes are used to generate the snippet
func Snippet(g goldmark.Markdown) (string, error) {
	if r, ok := g.Renderer().(snippetRenderer); ok {
		if r.buf == nil {
			return "", errors.New("Snippet buffer is nil; renderer.Render() may not have been called")
		}
		return r.buf.String(), nil
	}
	return "", errors.New("Failed to get underlying snippet renderer")
}

type snippetRenderer struct {
	baseRenderer renderer.Renderer
	buf          *snippetBuffer
}

func (r snippetRenderer) Walker(source []byte) ast.Walker {
	return func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if r.buf.IsFull() {
			return ast.WalkStop, nil
		}
		if entering && n.Kind() == ast.KindParagraph {
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				if c.Kind() == ast.KindImage {
					continue // Skip image alt text
				}
				r.buf.Write(c.Text(source))
				if t, ok := c.(*ast.Text); ok {
					if t.SoftLineBreak() {
						r.buf.WriteByte(' ')
					}
				}
			}
			r.buf.WriteByte(' ')
		}
		if n.Type() == ast.TypeBlock { // Don't go deeper than block nodes
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	}
}

func (r snippetRenderer) Render(w io.Writer, source []byte, n ast.Node) error {
	r.buf.Reset() // Clear buffer in case of previous render
	ast.Walk(n, r.Walker(source))
	return r.baseRenderer.Render(w, source, n)
}

func (r snippetRenderer) AddOptions(opts ...renderer.Option) {
	r.baseRenderer.AddOptions(opts...)
}

type snippetExtension struct {
	max int
}

func (e *snippetExtension) Extend(m goldmark.Markdown) {
	m.SetRenderer(snippetRenderer{
		m.Renderer(),
		newSnippetBuffer(e.max),
	})
}

// SnippetExtension returns a Goldmark extension
func SnippetExtension(max int) goldmark.Extender {
	return &snippetExtension{max}
}
