package extensions

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type tasklistTransformer struct{}

var defaultTasklistTransformer = &tasklistTransformer{}

func (b *tasklistTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && node.Kind() == ast.KindListItem {
			if node.FirstChild().FirstChild().Kind() == east.KindTaskCheckBox {
				node.SetAttributeString("class", []byte("task"))
			}
		}
		if node.Kind() == ast.KindDocument || node.Kind() == ast.KindList {
			return ast.WalkContinue, nil
		}
		return ast.WalkSkipChildren, nil
	})
}

type tasklistExtension struct{}

func (e *tasklistExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(defaultTasklistTransformer, 0),
		),
	)
}

// TasklistExtension is a Goldmark extension
var TasklistExtension = &tasklistExtension{}
