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
		{ // 29
			"do\na\nwhile\n(\n1\n)",
			"do a; while (1);",
			"do a; while (\n\t1\n);",
		},
		{ // 30
			"do\n{\n}\nwhile\n(\n1\n)",
			"do {} while (1);",
			"do {} while (\n\t1\n);",
		},
		{ // 31
			"while(a)b",
			"while (a) b;",
			"while (a) b;",
		},
		{ // 32
			"while\n(\na\n)\nb\n;",
			"while (a) b;",
			"while (\n	a\n) b;",
		},
		{ // 33
			"for\n(\n;\n;\n)\na",
			"for (;;) a;",
			"for (;;) a;",
		},
		{ // 34
			"for\n(a;;) b",
			"for (a;;) b;",
			"for (\n\ta;;\n) b;",
		},
		{ // 35
			"for(var a=b;c<d;e++){}",
			"for (var a = b; c < d; e++) {}",
			"for (var a = b; c < d; e++) {}",
		},
		{ // 36
			"for(\nvar a=b;\nc<d;\ne++){}",
			"for (var a = b; c < d; e++) {}",
			"for (\n	var a = b;\n	c < d;\n	e++\n) {}",
		},
		{ // 37
			"for(let a=b;c<d;e++){}",
			"for (let a = b; c < d; e++) {}",
			"for (let a = b; c < d; e++) {}",
		},
		{ // 38
			"for(\nlet a=b;\nc<d;\ne++){}",
			"for (let a = b; c < d; e++) {}",
			"for (\n	let a = b;\n	c < d;\n	e++\n) {}",
		},
		{ // 39
			"for(const a=b;c<d;e++){}",
			"for (const a = b; c < d; e++) {}",
			"for (const a = b; c < d; e++) {}",
		},
		{ // 40
			"for(\nconst a=b;\nc<d;\ne++){}",
			"for (const a = b; c < d; e++) {}",
			"for (\n	const a = b;\n	c < d;\n	e++\n) {}",
		},
		{ // 41
			"for(a in b){}",
			"for (a in b) {}",
			"for (a in b) {}",
		},
		{ // 42
			"for\n(a\nin\nb\n)\n{}",
			"for (a in b) {}",
			"for (\n	a in b\n) {}",
		},
		{ // 43
			"for(var a in b){}",
			"for (var a in b) {}",
			"for (var a in b) {}",
		},
		{ // 44
			"for\n(var\na\nin\nb\n)\n{}",
			"for (var a in b) {}",
			"for (var a in b) {}",
		},
		{ // 45
			"for(let a in b){}",
			"for (let a in b) {}",
			"for (let a in b) {}",
		},
		{ // 46
			"for\n(let\na\nin\nb\n)\n{}",
			"for (let a in b) {}",
			"for (let a in b) {}",
		},
		{ // 47
			"for(const a in b){}",
			"for (const a in b) {}",
			"for (const a in b) {}",
		},
		{ // 48
			"for\n(const\na\nin\nb\n)\n{}",
			"for (const a in b) {}",
			"for (const a in b) {}",
		},
		{ // 49
			"for(a of b){}",
			"for (a of b) {}",
			"for (a of b) {}",
		},
		{ // 50
			"for\n(a\nof\nb\n)\n{}",
			"for (a of b) {}",
			"for (\n	a of b\n) {}",
		},
		{ // 51
			"for(var a of b){}",
			"for (var a of b) {}",
			"for (var a of b) {}",
		},
		{ // 52
			"for\n(var\na\nof\nb\n)\n{}",
			"for (var a of b) {}",
			"for (var a of b) {}",
		},
		{ // 53
			"for(let a of b){}",
			"for (let a of b) {}",
			"for (let a of b) {}",
		},
		{ // 54
			"for\n(let\na\nof\nb\n)\n{}",
			"for (let a of b) {}",
			"for (let a of b) {}",
		},
		{ // 55
			"for(const a of b){}",
			"for (const a of b) {}",
			"for (const a of b) {}",
		},
		{ // 56
			"for\n(const\na\nof\nb\n)\n{}",
			"for (const a of b) {}",
			"for (const a of b) {}",
		},
		{ // 57
			"async () => {\nfor await(a of b) {}\n}",
			"async () => {\n	for await (a of b) {}\n};",
			"async () => {\n	for await (a of b) {}\n};",
		},
		{ // 58
			"async () => {\nfor\nawait(a\nof\nb)\n{}\n}",
			"async () => {\n	for await (a of b) {}\n};",
			"async () => {\n	for await (\n		a of b\n	) {}\n};",
		},
		{ // 59
			"switch(a) {}",
			"switch (a) {}",
			"switch (a) {}",
		},
		{ // 60
			"switch\n(\na\n)\n{\n}",
			"switch (a) {}",
			"switch (\n	a\n) {}",
		},
		{ // 61
			"switch(a){case b:case c:default:case d:case e:}",
			"switch (a) {\ncase b:\ncase c:\ndefault:\ncase d:\ncase e:\n}",
			"switch (a) {\ncase b:\ncase c:\ndefault:\ncase d:\ncase e:\n}",
		},
		{ // 62
			"switch\n\n(\n\na\n\n)\n\n{\n\ncase\n\nb\n\n:\n\ncase\n\nc\n\n:\n\ndefault\n\n:\n\ncase\n\nd\n\n:\n\ncase\n\ne\n\n:\n\n}",
			"switch (a) {\ncase b:\ncase c:\ndefault:\ncase d:\ncase e:\n}",
			"switch (\n	a\n) {\ncase b:\ncase c:\ndefault:\ncase d:\ncase e:\n}",
		},
		{ // 63
			"with(a)b",
			"with (a) b;",
			"with (a) b;",
		},
		{ // 64
			"with\n(\na\n)\nb",
			"with (a) b;",
			"with (\n	a\n) b;",
		},
		{ // 65
			"function a(){}",
			"function a() {}",
			"function a() {}",
		},
		{ // 66
			"function a(b){}",
			"function a(b) {}",
			"function a(b) {}",
		},
		{ // 67
			"function a(b,c){}",
			"function a(b, c) {}",
			"function a(b, c) {}",
		},
		{ // 68
			"function\na(\nb\n,\nc\n){}",
			"function a(b, c) {}",
			"function a(b, c) {}",
		},
		{ // 69
			"function*a(){}",
			"function* a() {}",
			"function* a() {}",
		},
		{ // 70
			"function* a(b){}",
			"function* a(b) {}",
			"function* a(b) {}",
		},
		{ // 71
			"function *a(b,c){}",
			"function* a(b, c) {}",
			"function* a(b, c) {}",
		},
		{ // 72
			"function\n*a(\nb\n,\nc\n){}",
			"function* a(b, c) {}",
			"function* a(b, c) {}",
		},
		{ // 73
			"async function a(){}",
			"async function a() {}",
			"async function a() {}",
		},
		{ // 74
			"async function a(b){}",
			"async function a(b) {}",
			"async function a(b) {}",
		},
		{ // 75
			"async function a(b,c){}",
			"async function a(b, c) {}",
			"async function a(b, c) {}",
		},
		{ // 76
			"async function\na(\nb\n,\nc\n){}",
			"async function a(b, c) {}",
			"async function a(b, c) {}",
		},
		{ // 77
			"async function*a(){}",
			"async function* a() {}",
			"async function* a() {}",
		},
		{ // 78
			"async function* a(b){}",
			"async function* a(b) {}",
			"async function* a(b) {}",
		},
		{ // 79
			"async function *a(b,c){}",
			"async function* a(b, c) {}",
			"async function* a(b, c) {}",
		},
		{ // 80
			"async function\n*a(\nb\n,\nc\n){}",
			"async function* a(b, c) {}",
			"async function* a(b, c) {}",
		},
		{ // 81
			"a = function(){}",
			"a = function () {};",
			"a = function () {};",
		},
		{ // 82
			"a=function(b){}",
			"a = function (b) {};",
			"a = function (b) {};",
		},
		{ // 83
			"a=function *(b,c){}",
			"a = function* (b, c) {};",
			"a = function* (b, c) {};",
		},
		{ // 84
			"a=function\n(\nb\n,\nc\n){}",
			"a = function (b, c) {};",
			"a = function (b, c) {};",
		},
		{ // 85
			"try{}catch{}",
			"try {} catch {}",
			"try {} catch {}",
		},
		{ // 86
			"try\n{\n}\ncatch\n{\n}",
			"try {} catch {}",
			"try {} catch {}",
		},
		{ // 87
			"try{}catch(a){}",
			"try {} catch (a) {}",
			"try {} catch (a) {}",
		},
		{ // 88
			"try\n{\n}\ncatch\n(\na\n)\n{\n}",
			"try {} catch (a) {}",
			"try {} catch (a) {}",
		},
		{ // 89
			"try{}catch({}){}",
			"try {} catch ({}) {}",
			"try {} catch ({}) {}",
		},
		{ // 90
			"try{}catch([]){}",
			"try {} catch ([]) {}",
			"try {} catch ([]) {}",
		},
		{ // 91
			"try{}finally{}",
			"try {} finally {}",
			"try {} finally {}",
		},
		{ // 92
			"try\n{\n}\nfinally\n{\n}",
			"try {} finally {}",
			"try {} finally {}",
		},
		{ // 93
			"try{}catch{}finally{}",
			"try {} catch {} finally {}",
			"try {} catch {} finally {}",
		},
		{ // 94
			"try\n{\n}\ncatch\n{\n}\nfinally\n{\n}",
			"try {} catch {} finally {}",
			"try {} catch {} finally {}",
		},
		{ // 95
			"try{}catch(a){}finally{}",
			"try {} catch (a) {} finally {}",
			"try {} catch (a) {} finally {}",
		},
		{ // 96
			"try\n{\n}\ncatch\n(\na\n)\n{\n}\nfinally\n{\n}",
			"try {} catch (a) {} finally {}",
			"try {} catch (a) {} finally {}",
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
