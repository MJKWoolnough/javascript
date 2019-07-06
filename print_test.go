package javascript

import (
	"testing"

	"vimagination.zapto.org/memio"
	"vimagination.zapto.org/parser"
)

type state struct {
	memio.Buffer
	Verbose bool
}

func (state) Width() (int, bool)     { return 0, false }
func (state) Precision() (int, bool) { return 0, false }
func (s *state) Flag(c int) bool     { return c == '+' == s.Verbose }

func TestPrintingScript(t *testing.T) {
	var st state
	for n, test := range [...]struct {
		Input, Output string
		Verbose       bool
	}{
		{ // 1
			Input:  `1;`,
			Output: `1;`,
		},
		{ // 2
			Input:  `1;2;`,
			Output: "1;\n\n2;",
		},
		{ // 3
			Input: `function	a	(  ){   }`,
			Output: `function a() {}`,
		},
		{ // 3
			Input:  "const a = function(){\n};",
			Output: "const a = function () {};",
		},
		{ // 4
			Input:  "function a(){}function b(){}",
			Output: "function a() {}\n\nfunction b() {}",
		},
		{ // 5
			Input: `class	a	extends	b	{c(d){alert(1);}e(f,
g){alert(2);}}`,
			Output: `class a extends b {
	c(d) {
		alert(1);
	}
	e(f, g) {
		alert(2);
	}
}`,
		},
		{ // 6
			Input: `a.Method(b=>1);c.Map((d,e)=>{f(e, d);});`,
			Output: `a.Method(b => 1);

c.Map((d, e) => {
	f(e, d);
});`,
		},
		{ // 7
			Input:  "var a = b&c==d||e^f==g&&h?i+j:k**l, m = n();",
			Output: "var a = b & c == d || e ^ f == g && h ? i + j : k**l, m = n();",
		},
		{ // 8
			Input:  "var a = b&c==d||e^f==g&&h?i+j:k**l,\nm = n();",
			Output: "var a = b & c == d || e ^ f == g && h ? i + j : k**l, m = n();",
		},
		{ // 9
			Input:   "var a = b&c==d||e^f==g&&h?i+j:k**l,\nm = n();",
			Output:  "var a = b & c == d || e ^ f == g && h ? i + j : k**l,\nm = n();",
			Verbose: true,
		},
		{ // 10
			Input: `for(
let a = 0;

a < 10;

a++
) {
	console.log(a);
}`,
			Output: `for (let a = 0; a < 10; a++) {
	console.log(a);
}`,
		},
		{ // 11
			Input: `for(
let a = 0;

a < 10;

a++
) {
	console.log(a);
}`,
			Output: `for (
	let a = 0;
	a < 10;
	a++
) {
	console.log(a);
}`,
			Verbose: true,
		},
	} {
		s, err := ParseScript(parser.NewStringTokeniser(test.Input))
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
			continue
		}
		st.Buffer = st.Buffer[:0]
		st.Verbose = test.Verbose
		s.Format(&st, 's')
		if str := string(st.Buffer); str != test.Output {
			t.Errorf("test %d: expecting %q, got %q\n%s", n+1, test.Output, str, s)
		}
	}
}
