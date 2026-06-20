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
			`import TAG_NAME from '@tag'; TAG_NAME()`,
			"import def from\"@tag\"\nconst a = (def())",
		},
		{ // 4
			"const a = <b />",
			`import * as TAG_NAME from '@tag'; TAG_NAME()`,
			"import*as ns from\"@tag\"\nconst a = (ns())",
		},
		{ // 5
			"const a = <b />",
			`import * as tag from '@tag'; tag.T('TAG_NAME')`,
			"import*as ns from\"@tag\"\nconst a = (ns.T(\"b\"))",
		},
		{ // 6
			"const a = <b c='1'/>",
			`tag('TAG_NAME', PARAMS)`,
			"const a = (tag(\"b\", {\"c\":\"1\"}))",
		},
		{ // 7
			"const a = <b><c d=<e />/></b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [(tag(\"c\", {\"d\":(tag(\"e\", {}, []))}, []))]))",
		},
		{ // 8
			"const a = <b><c d=<a />/></b>",
			`{{ if .InHTML }}import {TAG_NAME} from '@html';TAG_NAME(PARAMS, CHILDREN){{else}}tag('TAG_NAME', PARAMS, CHILDREN){{end}}`,
			"import{a as a_1,b}from\"@html\"\nconst a = (b({}, [(tag(\"c\", {\"d\":(a_1({}, []))}, []))]))",
		},
		{ // 9
			"import {b as z} from '@html';const a = <b><c d=<a />/></b>",
			`{{ if .InHTML }}import {TAG_NAME} from '@html';TAG_NAME(PARAMS, CHILDREN){{else}}tag('TAG_NAME', PARAMS, CHILDREN){{end}}`,
			"import {a as a_1,b as z} from '@html';const a = (z({}, [(tag(\"c\", {\"d\":(a_1({}, []))}, []))]))",
		},
		{ // 10
			"const a = <b />",
			`import '@OTHER';tag('TAG_NAME')`,
			"import\"@OTHER\"\nconst a = (tag(\"b\"))",
		},
		{ // 11
			"import ns from '@MODULE';\n\nconst a = <b />;",
			`import ns from '@MODULE';ns.tag('TAG_NAME')`,
			"import ns from '@MODULE';\n\nconst a = (ns.tag(\"b\"));",
		},
		{ // 12
			"import * as ns from '@MODULE';\n\nconst a = <b />;",
			`import * as ns from '@MODULE';ns.tag('TAG_NAME')`,
			"import * as ns from '@MODULE';\n\nconst a = (ns.tag(\"b\"));",
		},
		{ // 13
			"const a = <b {...d}/>",
			`tag('TAG_NAME', PARAMS)`,
			"const a = (tag(\"b\", {...d}))",
		},
		{ // 14
			"const a = <b><c /></b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [(tag(\"c\", {}, []))]))",
		},
		{ // 15
			"const a = <b><c />d{...e}</b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [(tag(\"c\", {}, [])),\"d\",...e]))",
		},
		{ // 16
			"const a = <></>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = []",
		},
		{ // 17
			"const a = <><b /></>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = [(tag(\"b\", {}, []))]",
		},
		{ // 18
			"const a = <b><><c /></></b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [[(tag(\"c\", {}, []))]]))",
		},
		{ // 19
			"const a = <b c=<></> />",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {\"c\":[]}, []))",
		},
		{ // 20
			"const a = <div></div>",
			`{{ if or .InHTML .InSVG }}import {TAG_NAME} from '@{{.Namespace}}';{{ end }}TAG_NAME(PARAMS, CHILDREN)`,
			"import{div}from\"@html\"\nconst a = (div({}, []))",
		},
		{ // 21
			"const a = <div><a /></div>",
			`{{ if or .InHTML .InSVG }}import {TAG_NAME} from '@{{.Namespace}}';{{ end }}TAG_NAME(PARAMS, CHILDREN)`,
			"import{a as a_1,div}from\"@html\"\nconst a = (div({}, [(a_1({}, []))]))",
		},
		{ // 22
			"const a = <div><svg><a /></svg></div>",
			`{{ if or .InHTML .InSVG }}import {TAG_NAME} from '@{{.Namespace}}';{{ end }}TAG_NAME(PARAMS, CHILDREN)`,
			"import{a as a_1,svg}from\"@svg\";import{div}from\"@html\"\nconst a = (div({}, [(svg({}, [(a_1({}, []))]))]))",
		},
		{ // 23
			"const a = <div><svg><a /></svg></div>",
			`{{ if or .InHTML .InSVG }}import {TAG_NAME} from '@{{.Namespace}}';{{ end }}TAG_NAME({{ if .HasParams }}PARAMS{{ end }}{{ if .HasChildren}}{{if .HasParams }}, {{ end }}CHILDREN{{ end }})`,
			"import{a as a_1,svg}from\"@svg\";import{div}from\"@html\"\nconst a = (div([(svg([(a_1())]))]))",
		},
		{ // 24
			"const a = <b>{// A\nc // B\n}</b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [// A\nc // B\n]))",
		},
		{ // 25
			"const a = <b>{// A\nc // B\n\n// C\n}</b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [// A\nc // B\n\n// C\n]))",
		},
		{ // 26
			"const a = <b>{// A\n... /* B */ c\n// C\n}</b>",
			`tag('TAG_NAME', PARAMS, CHILDREN)`,
			"const a = (tag(\"b\", {}, [// A\n... /* B */ c\n// C\n]))",
		},
		{ // 27
			"const a = <b {// A\n... /* B */ d // C\n}/>",
			`tag('TAG_NAME', PARAMS)`,
			"const a = (tag(\"b\", {// A\n... /* B */ d // C\n}))",
		},
		{ // 28
			"const a = <b {\n// A\n... /* B */ d // C\n}/>",
			`tag('TAG_NAME', PARAMS)`,
			"const a = (tag(\"b\", {\n// A\n... /* B */ d // C\n}))",
		},
		{ // 29
			"const a = <div><svg:a /></div>",
			`{{ if or .InHTML .InSVG }}import {TAG_NAME} from '@{{.Namespace}}';{{ end }}TAG_NAME({{ if .HasParams }}PARAMS{{ end }}{{ if .HasChildren}}{{if .HasParams }}, {{ end }}CHILDREN{{ end }})`,
			"import{a as a_1}from\"@svg\";import{div}from\"@html\"\nconst a = (div([(a_1())]))",
		},
		{ // 30
			"const a = <dialog open />",
			`TAG_NAME(PARAMS, CHILDREN)`,
			`const a = (dialog({"open":true}, []))`,
		},
		{ // 31
			"const a = <b // C\n/>",
			`TAG_NAME(PARAMS, CHILDREN)`,
			"const a = (b // C\n({}, []))",
		},
		{ // 32
			"const a = <b.c />",
			`TAG_NAME()`,
			"const a = (b.c())",
		},
		{ // 33
			"const a = <b.c.d />",
			`TAG_NAME()`,
			"const a = (b.c.d())",
		},
		{ // 34
			"const a = </*A*/b/*B*/./*C*/c/*D*/./*E*/d/*F*/ />",
			`TAG_NAME()`,
			"const a = (/*A*/b/*B*/./*C*/c/*D*/./*E*/d/*F*/())",
		},
		{ // 35
			"const a = </*A*/b/*B*/:/*C*/c/*D*//>",
			`TAG_NAME()`,
			"const a = (/*A*//*B*//*C*/c/*D*/())",
		},
		{ // 36
			"const a = </*A*/></*B*// /*C*/>",
			`TAG_NAME()`,
			"const a = [/*A*//*B*/ /*C*/]",
		},
		{ // 37
			"const a = <b c/*A*/=/*B*/\"d\"/*C*/></b>",
			`tag('TAG_NAME', PARAMS)`,
			"const a = (tag(\"b\", {\"c\"/*A*/:/*B*/\"d\"/*C*/}))",
		},
		{ // 38
			"const a = <b c/*A*/:/*B*/d/*C*/=/*D*/{/*E*/\"d\"/*F*/}/*G*/></b>",
			`tag('TAG_NAME', PARAMS)`,
			"const a = (tag(\"b\", {/*A*//*B*/\"c:d\"/*C*/:/*D*//*E*/\"d\"/*F*//*G*/}))",
		},
		{ // 39
			"function* a() {return <b c=/*A*/{/*B*/yield/*C*/d/*E*/}/*F*/></b>}",
			`TAG_NAME(PARAMS)`,
			"function* a() {return (b({\"c\":/*A*//*B*/yield/*C*/d/*E*//*F*/}))}",
		},
		{ // 40
			"async function a() {return <b c=/*A*/{/*B*/await/*C*/d/*E*/}/*F*/></b>}",
			`TAG_NAME(PARAMS)`,
			"async function a() {return (b({\"c\":/*A*//*B*/await/*C*/d/*E*//*F*/}))}",
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
