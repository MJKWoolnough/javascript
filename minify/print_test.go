package minify

import (
	"fmt"
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
		{ // 5
			"[] instanceof [].prototype",
			"[]instanceof[].prototype",
		},
		{ // 6
			"export * from 'a';",
			"export*from'a'",
		},
		{ // 7
			"export * as a from 'b';",
			"export*as a from'b'",
		},
		{ // 8
			"export {a, b as c, d} from 'f';",
			"export{a,b as c,d}from'f'",
		},
		{ // 9
			"import * as a from 'b';",
			"import*as a from'b'",
		},
		{ // 10
			"import {a, b as c, d} from 'e';",
			"import{a,b as c,d}from'e'",
		},
		{ // 11
			"import a from 'b';",
			"import a from'b'",
		},
		{ // 12
			"import a, {b, c} from 'e';",
			"import a,{b,c}from'e'",
		},
		{ // 13
			"a\nb\nc",
			"a;b;c",
		},
		{ // 14
			"a\n{}\nb",
			"a;{}b",
		},
		{ // 15
			"{a\nb\nc}",
			"{a;b;c}",
		},
		{ // 16
			"{a\n{}\nb}",
			"{a;{}b}",
		},
		{ // 17
			"if (a) b\nelse c",
			"if(a)b;else c",
		},
		{ // 18
			"if (a){\nb\n}\nelse{\nc\n}",
			"if(a){b}else{c}",
		},
		{ // 19
			"do a()\nwhile (1)",
			"do a();while(1)",
		},
		{ // 20
			"do{\na()\n}\nwhile (1)",
			"do{a()}while(1)",
		},
		{ // 21
			"switch(a){case a:\nb\ncase b:\n{}\ncase c: c}",
			"switch(a){case a:b;case b:{}case c:c}",
		},
		{ // 22
			"switch(a){case a:\nb\ndefault:\nc}",
			"switch(a){case a:b;default:c}",
		},
		{ // 23
			"switch(a){default:\na\ncase b: c}",
			"switch(a){default:a;case b:c}",
		},
		{ // 24
			"switch ( a ) { case []:\n1\ncase b: 2}",
			"switch(a){case[]:1;case b:2}",
		},
		{ // 25
			"switch ( a ) { case a:\na\nb\nc }",
			"switch(a){case a:a;b;c}",
		},
		{ // 26
			"class A {a\nb\nc\nd(){}\ne\n}",
			"class A{a;b;c;d(){}e}",
		},
		{ // 27
			"class A {static a = 1;static b(){}}",
			"class A{static a=1;static b(){}}",
		},
		{ // 28
			"class A {static [a] = 1;static [b](){}}",
			"class A{static[a]=1;static[b](){}}",
		},
		{ // 29
			"#a in b;",
			"#a in b",
		},
		{ // 30
			"#a in[b];",
			"#a in[b]",
		},
		{ // 31
			"import {a as b} from './c';",
			"import{a as b}from'./c'",
		},
		{ // 32
			"import * as a from './b';",
			"import*as a from'./b'",
		},
		{ // 33
			"var a = 1;",
			"var a=1",
		},
		{ // 34
			"var [a] = [1];",
			"var[a]=[1]",
		},
		{ // 35
			"function a(){}",
			"function a(){}",
		},
		{ // 36
			"(function (){})",
			"(function(){})",
		},
		{ // 37
			"async function a(){}",
			"async function a(){}",
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)
		m, err := javascript.ParseModule(&tk)
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else {
			var sb strings.Builder
			if _, err := Print(&sb, m); err != nil {
				t.Errorf("test %d.1: unexpected error: %s", n+1, err)
			} else if str := sb.String(); str != test.Output {
				t.Errorf("test %d.1: expecting output %q, got %q", n+1, test.Output, str)
			} else {
				normalStr := fmt.Sprint(m)
				tk = parser.NewStringTokeniser(str)
				m, err := javascript.ParseModule(&tk)
				if err != nil {
					t.Errorf("test %d.2: unexpected error: %s", n+1, err)
				} else if otherStr := fmt.Sprint(m); normalStr != otherStr {
					t.Errorf("test %d.2: normal output not equal, expecting %s, got %s", n+1, normalStr, otherStr)
				}
			}
		}
	}
}
