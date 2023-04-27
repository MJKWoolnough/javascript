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
