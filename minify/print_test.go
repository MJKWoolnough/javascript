package minify

import (
	"strings"
	"testing"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/parser"
)

func TestPrint(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output string
	}{
		{ // 1
			"var a = 1;",
			"var a=1",
		},
		{ // 2
			"var [a] = 1;",
			"var[a]=1",
		},
		{ // 3
			"async function a(){}",
			"async function a(){}",
		},
		{ // 4
			"typeof []",
			"typeof[]",
		},
		{ // 4
			"[] instanceof [].prototype",
			"[]instanceof[].prototype",
		},
		{ // 5
			"export * from 'a';",
			"export*from'a'",
		},
		{ // 6
			"export * as a from 'b';",
			"export*as a from'b'",
		},
		{ // 7
			"export {a, b as c, d} from 'f';",
			"export{a,b as c,d}from'f'",
		},
		{ // 8
			"import * as a from 'b';",
			"import*as a from'b'",
		},
		{ // 9
			"import {a, b as c, d} from 'e';",
			"import{a,b as c,d}from'e'",
		},
		{ // 10
			"import a from 'b';",
			"import a from'b'",
		},
		{ // 11
			"import a, {b, c} from 'e';",
			"import a,{b,c}from'e'",
		},
		{ // 12
			"a\nb\nc",
			"a;b;c",
		},
		{ // 13
			"a\n{}\nb",
			"a;{}b",
		},
		{ // 14
			"{a\nb\nc}",
			"{a;b;c}",
		},
		{ // 15
			"{a\n{}\nb}",
			"{a;{}b}",
		},
		{ // 16
			"if (a) b\nelse c",
			"if(a)b;else c",
		},
		{ // 17
			"if (a){\nb\n}\nelse{\nc\n}",
			"if(a){b}else{c}",
		},
		{ // 18
			"do a()\nwhile (1)",
			"do a();while(1)",
		},
		{ // 19
			"do{\na()\n}\nwhile (1)",
			"do{a()}while(1)",
		},
		{ // 20
			"switch(a){case a:\nb\ncase b:\n{}\ncase c: c}",
			"switch(a){case a:b;case b:{}case c:c}",
		},
		{ // 21
			"switch(a){case a:\nb\ndefault:\nc}",
			"switch(a){case a:b;default:c}",
		},
		{ // 22
			"switch(a){default:\na\ncase b: c}",
			"switch(a){default:a;case b:c}",
		},
		{ // 23
			"switch ( a ) { case []:\n1\ncase b: 2}",
			"switch(a){case[]:1;case b:2}",
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)
		m, err := javascript.ParseModule(&tk)
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else {
			var sb strings.Builder
			if _, err := Print(&sb, m); err != nil {
				t.Errorf("test %d: unexpected error: %s", n+1, err)
			} else if str := sb.String(); str != test.Output {
				t.Errorf("test %d: expecting output %q, got %q", n+1, test.Output, str)
			}
		}
	}
}
