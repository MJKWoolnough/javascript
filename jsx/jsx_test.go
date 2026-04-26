package jsx

import (
	"fmt"
	"testing"
	"text/template"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/parser"
)

func TestProcess(t *testing.T) {
	for n, test := range [...]struct {
		Input, Template, Output string
	}{
		{ // 1
			"const a = <b />",
			`tag('TAG_NAME')`,
			"const a = (tag(\"b\"));",
		},
		{ // 2
			"const a = <b />",
			`import tag from '@tag'; tag('TAG_NAME')`,
			"import def from \"@tag\";\n\nconst a = (def(\"b\"));",
		},
		{ // 3
			"const a = <b />",
			`import * as tag from '@tag'; tag.T('TAG_NAME')`,
			"import * as ns from \"@tag\";\n\nconst a = (ns.T(\"b\"));",
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)

		if m, err := javascript.ParseModule(javascript.AsJSX(&tk)); err != nil {
			t.Errorf("test %d: unexpected error parsing input: %s", n+1, err)
		} else if tmp, err := template.New("").Parse(test.Template); err != nil {
			t.Errorf("test %d: unexpected error parsing template: %s", n+1, err)
		} else if err := Process(m, tmp); err != nil {
			t.Errorf("test %d: unexpected error processing: %s", n+1, err)
		} else if output := fmt.Sprintf("%s", m); output != test.Output {
			t.Errorf("test %d: expecting output %q, got %q", n+1, test.Output, output)
		}
	}
}
