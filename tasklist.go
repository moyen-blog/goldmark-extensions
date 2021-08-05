package extensions

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type taskCheckBoxParser struct {
	super parser.InlineParser
}

var defaultTaskCheckBoxParser = &taskCheckBoxParser{
	extension.NewTaskCheckBoxParser(),
}

func (s *taskCheckBoxParser) Trigger() []byte {
	return s.super.Trigger()
}

func (s *taskCheckBoxParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	n := s.super.Parse(parent, block, pc)
	if n != nil && parent.Parent().Kind() == ast.KindListItem {
		parent.Parent().SetAttributeString("class", []byte("task"))
	}
	return n
}

type tasklistExtension struct{}

func (e *tasklistExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(defaultTaskCheckBoxParser, 0),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(extension.NewTaskCheckBoxHTMLRenderer(), 500),
	))
}

// TasklistExtension is a GoldMark extension thinly wrapping the tasklist extension
// It adds "task" as a class to tasklist items
// See https://github.com/yuin/goldmark/blob/master/extension/tasklist.go
var TasklistExtension = &tasklistExtension{}
