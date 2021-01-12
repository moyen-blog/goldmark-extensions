package meta

import (
	"bytes"
	"errors"
	"regexp"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"gopkg.in/yaml.v2"
)

var contextKeyMeta = parser.NewContextKey()

// Unmarshal attempts to unmarshal YAML metadata into the provided struct
func Unmarshal(pc parser.Context, out interface{}) error {
	v := pc.Get(contextKeyMeta)
	if v == nil {
		return nil // Buffer never set, no YAML
	}
	buf, ok := v.(bytes.Buffer)
	if !ok {
		return errors.New("Failed to get YAML buffer")
	}
	return yaml.Unmarshal(buf.Bytes(), out)
}

func isSeparator(line []byte) bool {
	r := regexp.MustCompile(`\s*-{3,}\s*`)
	return r.Match(line)
}

type metaParser struct{}

var defaultMetaParser = &metaParser{}

func (b *metaParser) Trigger() []byte {
	return []byte{'-'}
}

func (b *metaParser) Open(parent gast.Node, reader text.Reader, pc parser.Context) (gast.Node, parser.State) {
	linenum, _ := reader.Position()
	if linenum != 0 {
		return nil, parser.NoChildren
	}
	line, _ := reader.PeekLine()
	if isSeparator(line) {
		return gast.NewTextBlock(), parser.NoChildren
	}
	return nil, parser.NoChildren
}

func (b *metaParser) Continue(node gast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()
	if isSeparator(line) {
		reader.Advance(segment.Len())
		return parser.Close
	}
	node.Lines().Append(segment)
	return parser.Continue | parser.NoChildren
}

func (b *metaParser) Close(node gast.Node, reader text.Reader, pc parser.Context) {
	lines := node.Lines()
	var buf bytes.Buffer
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		buf.Write(segment.Value(reader.Source()))
	}
	pc.Set(contextKeyMeta, buf)
	node.Parent().RemoveChild(node.Parent(), node)
}

func (b *metaParser) CanInterruptParagraph() bool {
	return false
}

func (b *metaParser) CanAcceptIndentedLine() bool {
	return false
}

type metadataExtension struct{}

// MetadataExtension is a Goldmark extension
var MetadataExtension = &metadataExtension{}

func (e *metadataExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(defaultMetaParser, 0),
		),
	)
}
