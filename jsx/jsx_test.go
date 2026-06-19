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
			"const a = (tag(\"b\"))",
		},
		{ // 2
			"const a = <b />",
			`import tag from '@tag'; tag('TAG_NAME')`,
			"import def from\"@tag\"\nconst a = (def(\"b\"))",
		},
		{ // 3
			"const a = <b />",
			`import * as tag from '@tag'; tag.T('TAG_NAME')`,
			"import*as ns from\"@tag\"\nconst a = (ns.T(\"b\"))",
		},
		{ // 4
			"const a = <b c='1'/>",
			`tag('TAG_NAME', PARAMS)`,
			"const a = (tag(\"b\", {\"c\":\"1\"}))",
		},
		{ // 5
			"const a = <b><c d=<e />/></b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [(tag(\"c\", {\"d\":(tag(\"e\", {}, []))}, []))]))",
		},
		{ // 6
			"const a = <b><c d=<a />/></b>",
			`{{ if .InHTML }}import {TAG_NAME} from '@html';TAG_NAME(PARAMS, CHILDREN){{else}}tag('TAG_NAME', PARAMS, CHILDREN){{end}}`,
			"import{a as a_1,b}from\"@html\"\nconst a = (b({}, [(tag(\"c\", {\"d\":(a_1({}, []))}, []))]))",
		},
		{ // 7
			"import {b as z} from '@html';const a = <b><c d=<a />/></b>",
			`{{ if .InHTML }}import {TAG_NAME} from '@html';TAG_NAME(PARAMS, CHILDREN){{else}}tag('TAG_NAME', PARAMS, CHILDREN){{end}}`,
			"import {a as a_1,b as z} from '@html';const a = (z({}, [(tag(\"c\", {\"d\":(a_1({}, []))}, []))]))",
		},
		{ // 8
			"const a = <b />",
			`import '@OTHER';tag('TAG_NAME')`,
			"import\"@OTHER\"\nconst a = (tag(\"b\"))",
		},
		{ // 9
			"import ns from '@MODULE';\n\nconst a = <b />;",
			`import ns from '@MODULE';ns.tag('TAG_NAME')`,
			"import ns from '@MODULE';\n\nconst a = (ns.tag(\"b\"));",
		},
		{ // 10
			"import * as ns from '@MODULE';\n\nconst a = <b />;",
			`import * as ns from '@MODULE';ns.tag('TAG_NAME')`,
			"import * as ns from '@MODULE';\n\nconst a = (ns.tag(\"b\"));",
		},
		{ // 11
			"const a = <b {...d}/>",
			`tag('TAG_NAME', PARAMS)`,
			"const a = (tag(\"b\", {...d}))",
		},
		{ // 12
			"const a = <b><c /></b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [(tag(\"c\", {}, []))]))",
		},
		{ // 13
			"const a = <b><c />d{...e}</b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [(tag(\"c\", {}, [])),\"d\",...e]))",
		},
		{ // 14
			"const a = <></>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = []",
		},
		{ // 15
			"const a = <><b /></>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = [(tag(\"b\", {}, []))]",
		},
		{ // 16
			"const a = <b><><c /></></b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [[(tag(\"c\", {}, []))]]))",
		},
		{ // 17
			"const a = <b c=<></> />",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {\"c\":[]}, []))",
		},
		{ // 18
			"const a = <div></div>",
			`{{ if or .InHTML .InSVG }}import {TAG_NAME} from '@{{.Namespace}}';{{ end }}TAG_NAME(PARAMS, CHILDREN)`,
			"import{div}from\"@html\"\nconst a = (div({}, []))",
		},
		{ // 19
			"const a = <div><a /></div>",
			`{{ if or .InHTML .InSVG }}import {TAG_NAME} from '@{{.Namespace}}';{{ end }}TAG_NAME(PARAMS, CHILDREN)`,
			"import{a as a_1,div}from\"@html\"\nconst a = (div({}, [(a_1({}, []))]))",
		},
		{ // 20
			"const a = <div><svg><a /></svg></div>",
			`{{ if or .InHTML .InSVG }}import {TAG_NAME} from '@{{.Namespace}}';{{ end }}TAG_NAME(PARAMS, CHILDREN)`,
			"import{a as a_1,svg}from\"@svg\";import{div}from\"@html\"\nconst a = (div({}, [(svg({}, [(a_1({}, []))]))]))",
		},
		{ // 21
			"const a = <div><svg><a /></svg></div>",
			`{{ if or .InHTML .InSVG }}import {TAG_NAME} from '@{{.Namespace}}';{{ end }}TAG_NAME({{ if .HasParams }}PARAMS{{ end }}{{ if .HasChildren}}{{if .HasParams }}, {{ end }}CHILDREN{{ end }})`,
			"import{a as a_1,svg}from\"@svg\";import{div}from\"@html\"\nconst a = (div([(svg([(a_1())]))]))",
		},
		{ // 22
			"const a = <b>{// A\nc // B\n}</b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [// A\nc // B\n]))",
		},
		{ // 23
			"const a = <b>{// A\nc // B\n\n// C\n}</b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [// A\nc // B\n\n// C\n]))",
		},
		{ // 24
			"const a = <b>{// A\n... /* B */ c\n// C\n}</b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [// A\n... /* B */ c\n// C\n]))",
		},
		{ // 25
			"const a = <b {// A\n... /* B */ d // C\n}/>",
			`tag('TAG_NAME', PARAMS)`,
			"const a = (tag(\"b\", {// A\n... /* B */ d // C\n}))",
		},
		{ // 26
			"const a = <b {\n// A\n... /* B */ d // C\n}/>",
			`tag('TAG_NAME', PARAMS)`,
			"const a = (tag(\"b\", {\n// A\n... /* B */ d // C\n}))",
		},
		{ // 27
			"const a = <div><svg:a /></div>",
			`{{ if or .InHTML .InSVG }}import {TAG_NAME} from '@{{.Namespace}}';{{ end }}TAG_NAME({{ if .HasParams }}PARAMS{{ end }}{{ if .HasChildren}}{{if .HasParams }}, {{ end }}CHILDREN{{ end }})`,
			"import{a as a_1}from\"@svg\";import{div}from\"@html\"\nconst a = (div([(a_1())]))",
		},
		{ // 28
			"const a = <dialog open />",
			`TAG_NAME(PARAMS, CHILDREN)`,
			`const a = (dialog({"open":true}, []))`,
		},
		{ // 29
			"const a = <b // C\n/>",
			`TAG_NAME(PARAMS, CHILDREN)`,
			"const a = (b // C\n({}, []))",
		},
		{ // 30
			"const a = <b.c />",
			`TAG_NAME()`,
			"const a = (b.c())",
		},
		{ // 31
			"const a = <b.c.d />",
			`TAG_NAME()`,
			"const a = (b.c.d())",
		},
		{ // 32
			"const a = </*A*/b/*B*/./*C*/c/*D*/./*E*/d/*F*/ />",
			`TAG_NAME()`,
			"const a = (/*A*/b/*B*/./*C*/c/*D*/./*E*/d/*F*/())",
		},
		{ // 33
			"const a = </*A*/b/*B*/:/*C*/c/*D*//>",
			`TAG_NAME()`,
			"const a = (/*A*//*B*//*C*/c/*D*/())",
		},
		{ // 34
			"const a = </*A*/></*B*// /*C*/>",
			`TAG_NAME()`,
			"const a = [/*A*//*B*/ /*C*/]",
		},
		{ // 35
			"const a = <b c/*A*/=/*B*/\"d\"/*C*/></b>",
			`tag('TAG_NAME', PARAMS)`,
			"const a = (tag(\"b\", {\"c\"/*A*/:/*B*/\"d\"/*C*/}))",
		},
		{ // 36
			"const a = <b c/*A*/:/*B*/d/*C*/=/*D*/{/*E*/\"d\"/*F*/}/*G*/></b>",
			`tag('TAG_NAME', PARAMS)`,
			"const a = (tag(\"b\", {/*A*//*B*/\"d\"/*C*/:/*D*//*E*/\"d\"/*F*//*G*/}))",
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)

		if m, err := javascript.ParseModule(javascript.AsJSX(&tk)); err != nil {
			t.Errorf("test %d: unexpected error parsing input: %s", n+1, err)
		} else if tmp, err := template.New("").Parse(test.Template); err != nil {
			t.Errorf("test %d: unexpected error parsing template: %s", n+1, err)
		} else if err := Process(m, tmp); err != nil {
			t.Errorf("test %d: unexpected error processing: %s", n+1, err)
		} else if output := fmt.Sprintf("%#s", m); output != test.Output {
			t.Errorf("test %d: expecting output %q, got %q", n+1, test.Output, output)
		}
	}
}
