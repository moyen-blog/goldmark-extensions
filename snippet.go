package extensions

import (
	"errors"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var contextKeySnippet = parser.NewContextKey()

// Snippet returns the snippet made up of text from the underlying Markdown document
// Only paragraph nodes are used to generate the snippet
func Snippet(pc parser.Context) (string, error) {
	v := pc.Get(contextKeySnippet)
	s, ok := v.(string)
	if !ok {
		return "", errors.New("failed to get snippet")
	}
	return s, nil
}

type snippetTransformer struct {
	max int
}

func (r snippetTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	buf := newSnippetBuffer(r.max)
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if buf.IsFull() {
			return ast.WalkStop, nil
		}
		if entering && n.Kind() == ast.KindParagraph {
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				if c.Kind() == ast.KindImage {
					continue // Skip image alt text
				}
				buf.Write(c.Text(reader.Source()))
				if t, ok := c.(*ast.Text); ok {
					if t.SoftLineBreak() {
						buf.WriteByte(' ')
					}
				}
			}
			buf.WriteByte(' ')
		}
		if n.Type() == ast.TypeBlock { // Don't go deeper than block nodes
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	pc.Set(contextKeySnippet, buf.String())
}

type snippetExtension struct {
	max int
}

func (e *snippetExtension) Extend(m goldmark.Markdown) {
	p := int(^uint(0) >> 1) // Lowest priority
	m.Parser().AddOptions(parser.WithASTTransformers(
		util.Prioritized(snippetTransformer{e.max}, p), // Generate snippet after all other transformers applied
	))
}

// SnippetExtension returns a GoldMark extension
func SnippetExtension(max int) goldmark.Extender {
	return &snippetExtension{max}
}
