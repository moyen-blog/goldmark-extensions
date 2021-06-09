package extensions

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
)

var markdownTable = goldmark.New(
	goldmark.WithExtensions(
		TableExtension,
	),
)

func TestTable(t *testing.T) {
	source := `One | Two
--- | ---
Foo | Bar`

	var buf bytes.Buffer
	if err := markdownTable.Convert([]byte(source), &buf); err != nil {
		t.Error("Failed to convert markdown")
	}
	expected := `<div class="table-wrapper">
<table>
<thead>
<tr>
<th>One</th>
<th>Two</th>
</tr>
</thead>
<tbody>
<tr>
<td>Foo</td>
<td>Bar</td>
</tr>
</tbody>
</table>
</div>
`
	if buf.String() != expected {
		t.Errorf("Should render '%s', but got '%s'", expected, buf.String())
	}
}
