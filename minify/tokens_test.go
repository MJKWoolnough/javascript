package minify

import "testing"

func TestIsIdentifier(t *testing.T) {
	for n, test := range [...]struct {
		Input        string
		IsIdentifier bool
	}{
		{
			"a",
			true,
		},
		{
			"aa",
			true,
		},
		{
			"_",
			true,
		},
		{
			"__",
			true,
		},
		{
			"0",
			false,
		},
		{
			"true",
			true,
		},
		{
			"a a",
			false,
		},
	} {
		if ii := isIdentifier(test.Input); ii != test.IsIdentifier {
			t.Errorf("test %d: for input %q, expecting IsIdentifier to return %v, got %v", n+1, test.Input, test.IsIdentifier, ii)
		}
	}
}

func TestIsSimpleNumber(t *testing.T) {
	for n, test := range [...]struct {
		Input          string
		IsSimpleNumber bool
	}{
		{
			"",
			false,
		},
		{
			"a",
			false,
		},
		{
			"0",
			true,
		},
		{
			"01",
			false,
		},
		{
			"1",
			true,
		},
		{
			"1234567890",
			true,
		},
		{
			"1234567890a",
			false,
		},
		{
			"9007199254740990",
			true,
		},
		{
			"9007199254740991",
			true,
		},
		{
			"9007199254740992",
			false,
		},
		{
			"19007199254740992",
			false,
		},
	} {
		if isn := isSimpleNumber(test.Input); isn != test.IsSimpleNumber {
			t.Errorf("test %d: for input %s, got %v, when expecting %v", n+1, test.Input, test.IsSimpleNumber, isn)
		}
	}
}
