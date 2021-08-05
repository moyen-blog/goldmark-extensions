package extensions

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
)

var markdownTasklist = goldmark.New(
	goldmark.WithExtensions(
		TasklistExtension,
	),
)

func TestUnorderedTasklist(t *testing.T) {
	source := `- [ ] Unchecked
- [x] Checked`

	var buf bytes.Buffer
	if err := markdownTasklist.Convert([]byte(source), &buf); err != nil {
		t.Error("Failed to convert Markdown")
	}
	expected := `<ul>
<li class="task"><input disabled="" type="checkbox"> Unchecked</li>
<li class="task"><input checked="" disabled="" type="checkbox"> Checked</li>
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
		t.Error("Failed to convert Markdown")
	}
	expected := `<ol>
<li class="task"><input disabled="" type="checkbox"> Unchecked</li>
<li class="task"><input checked="" disabled="" type="checkbox"> Checked</li>
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
		t.Error("Failed to convert Markdown")
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

func TestMissingListItem(t *testing.T) {
	source := `* One
*
* Three`

	var buf bytes.Buffer
	if err := markdownTasklist.Convert([]byte(source), &buf); err != nil {
		t.Error("Failed to convert Markdown")
	}
	expected := `<ul>
<li>One</li>
<li></li>
<li>Three</li>
</ul>
`
	if buf.String() != expected {
		t.Errorf("Should render '%s', but got '%s'", expected, buf.String())
	}
}

func TestEmbeddedTaskList(t *testing.T) {
	source := `1. One
    - [ ] Unchecked
    - [x] Checked
2. Two`

	var buf bytes.Buffer
	if err := markdownTasklist.Convert([]byte(source), &buf); err != nil {
		t.Error("Failed to convert Markdown")
	}
	expected := `<ol>
<li>One
<ul>
<li class="task"><input disabled="" type="checkbox"> Unchecked</li>
<li class="task"><input checked="" disabled="" type="checkbox"> Checked</li>
</ul>
</li>
<li>Two</li>
</ol>
`
	if buf.String() != expected {
		t.Errorf("Should render '%s', but got '%s'", expected, buf.String())
	}
}
