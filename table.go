package extensions

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type tableHTMLRenderer struct {
	super renderer.NodeRenderer
}

var defaultTableHTMLRenderer = &tableHTMLRenderer{
	extension.NewTableHTMLRenderer(),
}

func (r *tableHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	r.super.RegisterFuncs(reg)
	reg.Register(east.KindTable, r.renderTable) // Override default render function
}

// renderTable overrides the table extension default table render function
// The only difference should be the table element being wrapped in a div
func (r *tableHTMLRenderer) renderTable(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<div class=\"table-wrapper\">\n<table")
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, extension.TableAttributeFilter)
		}
		_, _ = w.WriteString(">\n")
	} else {
		_, _ = w.WriteString("</table>\n</div>\n")
	}
	return ast.WalkContinue, nil
}

type tableExtension struct{}

func (e *tableExtension) Extend(m goldmark.Markdown) {
	extension.Table.Extend(m)
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(defaultTableHTMLRenderer, 450),
	))
}

// TableExtension is a Goldmark extension thinly wrapping the table extension
// It wraps each table in a containing div
// See https://github.com/yuin/goldmark/blob/master/extension/table.go
var TableExtension = &tableExtension{}
