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
		Input, SimpleOutput, VerboseOutput string
	}{
		{ // 1
			"1;",
			"1;",
			"1;",
		},
		{ // 2
			"1",
			"1;",
			"1;",
		},
		{ // 3
			"1;2;",
			"1;\n\n2;",
			"1;\n\n2;",
		},
		{ // 4
			"1\n2;",
			"1;\n\n2;",
			"1;\n\n2;",
		},
		{ // 5
			"continue",
			"continue;",
			"continue;",
		},
		{ // 6
			"break",
			"break;",
			"break;",
		},
		{ // 7
			"break a",
			"break a;",
			"break a;",
		},
		{ // 8
			"() => {return}",
			"() => {\n	return;\n};",
			"() => { return; };",
		},
		{ // 9
			"() => {return a}",
			"() => {\n	return a;\n};",
			"() => { return a; };",
		},
		{ // 10
			"throw a",
			"throw a;",
			"throw a;",
		},
		{ // 11
			"{\n1\n}",
			"{\n	1;\n}",
			"{\n	1;\n}",
		},
		{ // 12
			"{\n1\n2\n}",
			"{\n	1;\n	2;\n}",
			"{\n	1;\n	2;\n}",
		},
		{ // 13
			"{1;}",
			"{\n	1;\n}",
			"{ 1; }",
		},
		{ // 14
			"{1;2;}",
			"{\n	1;\n	2;\n}",
			"{ 1; 2; }",
		},
		{ // 15
			"var\na;",
			"var a;",
			"var a;",
		},
		{ // 16
			"var\na\n=\n1;",
			"var a = 1;",
			"var a = 1;",
		},
		{ // 17
			"var\na\n=\n1\n,\nb\n=\n2",
			"var a = 1, b = 2;",
			"var a = 1,\nb = 2;",
		},
		{ // 18
			"var a=1,b=2",
			"var a = 1, b = 2;",
			"var a = 1, b = 2;",
		},
		{ // 19
			"a,b,c,d",
			"a, b, c, d;",
			"a, b, c, d;",
		},
		{ // 20
			"a\n,\nb\n,\nc\n,\nd",
			"a, b, c, d;",
			"a,\nb,\nc,\nd;",
		},
		{ // 21
			"if(a){}",
			"if (a) {}",
			"if (a) {}",
		},
		{ // 22
			"if\n(\na\n)\n{\n}",
			"if (a) {}",
			"if (\n	a\n) {}",
		},
		{ // 23
			"if(a)b; else c",
			"if (a) b; else c;",
			"if (a) b; else c;",
		},
		{ // 24
			"if\n(\na\n)\nb\nelse\nc",
			"if (a) b; else c;",
			"if (\n	a\n) b; else c;",
		},
		{ // 25
			"if(a){b}else{c}",
			"if (a) {\n	b;\n} else {\n	c;\n}",
			"if (a) { b; } else { c; }",
		},
		{ // 26
			"if\n(\na\n)\n{\nb\n}\nelse\n{\nc\n}",
			"if (a) {\n	b;\n} else {\n	c;\n}",
			"if (\n	a\n) {\n	b;\n} else {\n	c;\n}",
		},
		{ // 27
			"do\n	a\nwhile(1)",
			"do a; while (1);",
			"do a; while (1);",
		},
		{ // 28
			"do{}while(1)",
			"do {} while (1);",
			"do {} while (1);",
		},
		{ // 27
			"do\na\nwhile\n(\n1\n)",
			"do a; while (1);",
			"do a; while (1);",
		},
		{ // 28
			"do\n{\n}\nwhile\n(\n1\n)",
			"do {} while (1);",
			"do {} while (1);",
		},
	} {
		s, err := ParseScript(parser.NewStringTokeniser(test.Input))
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
			continue
		}
		st.Verbose = false
		st.Buffer = st.Buffer[:0]
		s.Format(&st, 's')
		if str := string(st.Buffer); str != test.SimpleOutput {
			t.Errorf("test %d.1: expecting %q, got %q\n%s", n+1, test.SimpleOutput, str, s)
		}
		st.Verbose = true
		st.Buffer = st.Buffer[:0]
		s.Format(&st, 's')
		if str := string(st.Buffer); str != test.VerboseOutput {
			t.Errorf("test %d.2: expecting %q, got %q\n%s", n+1, test.VerboseOutput, str, s)
		}
	}
}

func TestModulePrinting(t *testing.T) {
	var st state
	for n, test := range [...]struct {
		Input, Output string
		Verbose       bool
	}{
		{ // 1
			Input:  `1;`,
			Output: `1;`,
		},
	} {
		m, err := ParseModule(parser.NewStringTokeniser(test.Input))
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
			continue
		}
		st.Buffer = st.Buffer[:0]
		st.Verbose = test.Verbose
		m.Format(&st, 's')
		if str := string(st.Buffer); str != test.Output {
			t.Errorf("test %d: expecting %q, got %q\n%s", n+1, test.Output, str, m)
		}
	}
}
