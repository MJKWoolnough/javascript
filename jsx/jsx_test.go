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
		{ // 4
			"const a = <b c='1'/>",
			`tag('TAG_NAME', PARAMS)`,
			"const a = (tag(\"b\", {c: \"1\"}));",
		},
		{ // 5
			"const a = <b><c d=<e />/></b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [(tag(\"c\", {d: (tag(\"e\", {}, []))}, []))]));",
		},
		{ // 6
			"const a = <b><c d=<a />/></b>",
			`{{ if .InHTML }}import {TAG_NAME} from '@html';TAG_NAME(PARAMS, CHILDREN){{else}}tag('TAG_NAME', PARAMS, CHILDREN){{end}}`,
			"import {a as a_1, b} from \"@html\";\n\nconst a = (b({}, [(tag(\"c\", {d: (a_1({}, []))}, []))]));",
		},
		{ // 7
			"import {b as z} from '@html';const a = <b><c d=<a />/></b>",
			`{{ if .InHTML }}import {TAG_NAME} from '@html';TAG_NAME(PARAMS, CHILDREN){{else}}tag('TAG_NAME', PARAMS, CHILDREN){{end}}`,
			"import {a as a_1, b as z} from '@html';\n\nconst a = (z({}, [(tag(\"c\", {d: (a_1({}, []))}, []))]));",
		},
		{ // 8
			"const a = <b />",
			`import '@OTHER';tag('TAG_NAME')`,
			"import \"@OTHER\";\n\nconst a = (tag(\"b\"));",
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
