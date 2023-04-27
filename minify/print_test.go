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
		{
			"var a = 1;",
			"var a=1",
		},
		{
			"var [a] = 1;",
			"var[a]=1",
		},
		{
			"async function a(){}",
			"async function a(){}",
		},
		{
			"typeof []",
			"typeof[]",
		},
		{
			"[] instanceof [].prototype",
			"[]instanceof[].prototype",
		},
		{
			"export * from 'a';",
			"export*from'a'",
		},
		{
			"export * as a from 'b';",
			"export*as a from'b'",
		},
		{
			"export {a, b as c, d} from 'f';",
			"export{a,b as c,d}from'f'",
		},
		{
			"import * as a from 'b';",
			"import*as a from'b'",
		},
		{
			"import {a, b as c, d} from 'e';",
			"import{a,b as c,d}from'e'",
		},
		{
			"import a from 'b';",
			"import a from'b'",
		},
		{
			"import a, {b, c} from 'e';",
			"import a,{b,c}from'e'",
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
