package extensions

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var markdownTasklist = goldmark.New(
	goldmark.WithExtensions(
		extension.TaskList,
		TasklistExtension,
	),
)

func TestUnorderedTasklist(t *testing.T) {
	source := `- [ ] Unchecked
- [x] Checked`

	var buf bytes.Buffer
	if err := markdownTasklist.Convert([]byte(source), &buf); err != nil {
		t.Error("Failed to convert markdown")
	}
	expected := `<ul>
<li style="list-style-type:none"><input disabled="" type="checkbox"> Unchecked</li>
<li style="list-style-type:none"><input checked="" disabled="" type="checkbox"> Checked</li>
</ul>
`
	if buf.String() != expected {
		t.Errorf("Should render '%s', but got '%s'", expected, buf.String())
	}
}

func TestOrderedTasklist(t *testing.T) {
	source := `1. [ ] Unchecked
2. [x] Checked`

	var buf bytes.Buffer
	if err := markdownTasklist.Convert([]byte(source), &buf); err != nil {
		t.Error("Failed to convert markdown")
	}
	expected := `<ol>
<li style="list-style-type:none"><input disabled="" type="checkbox"> Unchecked</li>
<li style="list-style-type:none"><input checked="" disabled="" type="checkbox"> Checked</li>
</ol>
`
	if buf.String() != expected {
		t.Errorf("Should render '%s', but got '%s'", expected, buf.String())
	}
}

func TestStandardList(t *testing.T) {
	source := `- One
- Two
- Three`

	var buf bytes.Buffer
	if err := markdownTasklist.Convert([]byte(source), &buf); err != nil {
		t.Error("Failed to convert markdown")
	}
	expected := `<ul>
<li>One</li>
<li>Two</li>
<li>Three</li>
</ul>
`
	if buf.String() != expected {
		t.Errorf("Should render '%s', but got '%s'", expected, buf.String())
	}
}
