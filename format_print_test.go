package javascript

import (
	"bytes"
	"testing"

	"vimagination.zapto.org/parser"
)

type state struct {
	bytes.Buffer
	Verbose bool
}

func (state) Width() (int, bool) { return 0, false }

func (state) Precision() (int, bool) { return 0, false }

func (s *state) Flag(c int) bool { return c == '+' == s.Verbose }

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
			"() => {\n\treturn;\n};",
			"() => {\n\treturn;\n};",
		},
		{ // 9
			"() => {return a}",
			"() => {\n\treturn a;\n};",
			"() => {\n\treturn a;\n};",
		},
		{ // 10
			"throw a",
			"throw a;",
			"throw a;",
		},
		{ // 11
			"{\n1\n}",
			"{\n\t1;\n}",
			"{\n\t1;\n}",
		},
		{ // 12
			"{\n1\n2\n}",
			"{\n\t1;\n\t2;\n}",
			"{\n\t1;\n\t2;\n}",
		},
		{ // 13
			"{1;}",
			"{\n\t1;\n}",
			"{\n\t1;\n}",
		},
		{ // 14
			"{1;2;}",
			"{\n\t1;\n\t2;\n}",
			"{\n\t1;\n\t2;\n}",
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
			"if (a) {}",
		},
		{ // 23
			"if(a)b; else c",
			"if (a) b; else c;",
			"if (a) b; else c;",
		},
		{ // 24
			"if\n(\na\n)\nb\nelse\nc",
			"if (a) b; else c;",
			"if (a) b; else c;",
		},
		{ // 25
			"if(a){b}else{c}",
			"if (a) {\n\tb;\n} else {\n\tc;\n}",
			"if (a) {\n\tb;\n} else {\n\tc;\n}",
		},
		{ // 26
			"if\n(\na\n)\n{\nb\n}\nelse\n{\nc\n}",
			"if (a) {\n\tb;\n} else {\n\tc;\n}",
			"if (a) {\n\tb;\n} else {\n\tc;\n}",
		},
		{ // 27
			"if ( // A\na) b",
			"if (a) b;",
			"if ( // A\n\n\ta\n) b;",
		},
		{ // 28
			"do\n\ta\nwhile(1)",
			"do a; while (1);",
			"do a; while (1);",
		},
		{ // 29
			"do{}while(1)",
			"do {} while (1);",
			"do {} while (1);",
		},
		{ // 30
			"do\na\nwhile\n(\n1\n)",
			"do a; while (1);",
			"do a; while (1);",
		},
		{ // 31
			"do\n{\n}\nwhile\n(\n1\n)",
			"do {} while (1);",
			"do {} while (1);",
		},
		{ // 32
			"do\n{\n}\nwhile\n( // A\na)",
			"do {} while (a);",
			"do {} while ( // A\n\n\ta\n);",
		},
		{ // 33
			"while(a)b",
			"while (a) b;",
			"while (a) b;",
		},
		{ // 34
			"while\n(\na\n)\nb\n;",
			"while (a) b;",
			"while (a) b;",
		},
		{ // 35
			"while ( // A\na) b;",
			"while (a) b;",
			"while ( // A\n\n\ta\n) b;",
		},
		{ // 36
			"for\n(\n;\n;\n)\na",
			"for (;;) a;",
			"for (;;) a;",
		},
		{ // 37
			"for\n(a;;) b",
			"for (a;;) b;",
			"for (\n\ta;;\n) b;",
		},
		{ // 38
			"for(var a=b;c<d;e++){}",
			"for (var a = b; c < d; e++) {}",
			"for (var a = b; c < d; e++) {}",
		},
		{ // 39
			"for(\nvar a=b;\nc<d;\ne++){}",
			"for (var a = b; c < d; e++) {}",
			"for (\n\tvar a = b;\n\tc < d;\n\te++\n) {}",
		},
		{ // 40
			"for(let a=b;c<d;e++){}",
			"for (let a = b; c < d; e++) {}",
			"for (let a = b; c < d; e++) {}",
		},
		{ // 41
			"for(\nlet a=b;\nc<d;\ne++){}",
			"for (let a = b; c < d; e++) {}",
			"for (\n\tlet a = b;\n\tc < d;\n\te++\n) {}",
		},
		{ // 42
			"for(const a=b;c<d;e++){}",
			"for (const a = b; c < d; e++) {}",
			"for (const a = b; c < d; e++) {}",
		},
		{ // 43
			"for(\nconst a=b;\nc<d;\ne++){}",
			"for (const a = b; c < d; e++) {}",
			"for (\n\tconst a = b;\n\tc < d;\n\te++\n) {}",
		},
		{ // 44
			"for(a in b){}",
			"for (a in b) {}",
			"for (a in b) {}",
		},
		{ // 45
			"for\n(a\nin\nb\n)\n{}",
			"for (a in b) {}",
			"for (\n\ta in b\n) {}",
		},
		{ // 46
			"for(var a in b){}",
			"for (var a in b) {}",
			"for (var a in b) {}",
		},
		{ // 47
			"for\n(var\na\nin\nb\n)\n{}",
			"for (var a in b) {}",
			"for (var a in b) {}",
		},
		{ // 48
			"for(let a in b){}",
			"for (let a in b) {}",
			"for (let a in b) {}",
		},
		{ // 49
			"for\n(let\na\nin\nb\n)\n{}",
			"for (let a in b) {}",
			"for (let a in b) {}",
		},
		{ // 50
			"for(const a in b){}",
			"for (const a in b) {}",
			"for (const a in b) {}",
		},
		{ // 51
			"for\n(const\na\nin\nb\n)\n{}",
			"for (const a in b) {}",
			"for (const a in b) {}",
		},
		{ // 52
			"for(a of b){}",
			"for (a of b) {}",
			"for (a of b) {}",
		},
		{ // 53
			"for\n(a\nof\nb\n)\n{}",
			"for (a of b) {}",
			"for (\n\ta of b\n) {}",
		},
		{ // 54
			"for(var a of b){}",
			"for (var a of b) {}",
			"for (var a of b) {}",
		},
		{ // 55
			"for\n(var\na\nof\nb\n)\n{}",
			"for (var a of b) {}",
			"for (var a of b) {}",
		},
		{ // 56
			"for(let a of b){}",
			"for (let a of b) {}",
			"for (let a of b) {}",
		},
		{ // 57
			"for\n(let\na\nof\nb\n)\n{}",
			"for (let a of b) {}",
			"for (let a of b) {}",
		},
		{ // 58
			"for(const a of b){}",
			"for (const a of b) {}",
			"for (const a of b) {}",
		},
		{ // 59
			"for\n(const\na\nof\nb\n)\n{}",
			"for (const a of b) {}",
			"for (const a of b) {}",
		},
		{ // 60
			"async () => {\nfor await(a of b) {}\n}",
			"async () => {\n\tfor await (a of b) {}\n};",
			"async () => {\n\tfor await (a of b) {}\n};",
		},
		{ // 61
			"async () => {\nfor\nawait(a\nof\nb)\n{}\n}",
			"async () => {\n\tfor await (a of b) {}\n};",
			"async () => {\n\tfor await (\n\t	a of b\n\t) {}\n};",
		},
		{ // 62
			"switch(a) {}",
			"switch (a) {}",
			"switch (a) {}",
		},
		{ // 63
			"switch\n(\na\n)\n{\n}",
			"switch (a) {}",
			"switch (\n\ta\n) {}",
		},
		{ // 64
			"switch(a){case b:case c:default:case d:case e:}",
			"switch (a) {\ncase b:\ncase c:\ndefault:\ncase d:\ncase e:\n}",
			"switch (a) {\ncase b:\ncase c:\ndefault:\ncase d:\ncase e:\n}",
		},
		{ // 65
			"switch\n\n(\n\na\n\n)\n\n{\n\ncase\n\nb\n\n:\n\ncase\n\nc\n\n:\n\ndefault\n\n:\n\ncase\n\nd\n\n:\n\ncase\n\ne\n\n:\n\n}",
			"switch (a) {\ncase b:\ncase c:\ndefault:\ncase d:\ncase e:\n}",
			"switch (\n\ta\n) {\ncase b:\ncase c:\ndefault:\ncase d:\ncase e:\n}",
		},
		{ // 66
			"with(a)b",
			"with (a) b;",
			"with (a) b;",
		},
		{ // 67
			"with\n(\na\n)\nb",
			"with (a) b;",
			"with (\n\ta\n) b;",
		},
		{ // 68
			"function a(){}",
			"function a() {}",
			"function a() {}",
		},
		{ // 69
			"function a(b){}",
			"function a(b) {}",
			"function a(b) {}",
		},
		{ // 70
			"function a(b,c){}",
			"function a(b, c) {}",
			"function a(b, c) {}",
		},
		{ // 71
			"function\na(\nb\n,\nc\n){}",
			"function a(b, c) {}",
			"function a(b, c) {}",
		},
		{ // 72
			"function*a(){}",
			"function* a() {}",
			"function* a() {}",
		},
		{ // 73
			"function* a(b){}",
			"function* a(b) {}",
			"function* a(b) {}",
		},
		{ // 74
			"function *a(b,c){}",
			"function* a(b, c) {}",
			"function* a(b, c) {}",
		},
		{ // 75
			"function\n*a(\nb\n,\nc\n){}",
			"function* a(b, c) {}",
			"function* a(b, c) {}",
		},
		{ // 76
			"async function a(){}",
			"async function a() {}",
			"async function a() {}",
		},
		{ // 77
			"async function a(b){}",
			"async function a(b) {}",
			"async function a(b) {}",
		},
		{ // 78
			"async function a(b,c){}",
			"async function a(b, c) {}",
			"async function a(b, c) {}",
		},
		{ // 79
			"async function\na(\nb\n,\nc\n){}",
			"async function a(b, c) {}",
			"async function a(b, c) {}",
		},
		{ // 80
			"async function*a(){}",
			"async function* a() {}",
			"async function* a() {}",
		},
		{ // 81
			"async function* a(b){}",
			"async function* a(b) {}",
			"async function* a(b) {}",
		},
		{ // 82
			"async function *a(b,c){}",
			"async function* a(b, c) {}",
			"async function* a(b, c) {}",
		},
		{ // 83
			"async function\n*a(\nb\n,\nc\n){}",
			"async function* a(b, c) {}",
			"async function* a(b, c) {}",
		},
		{ // 84
			"a = function(){}",
			"a = function () {};",
			"a = function () {};",
		},
		{ // 85
			"a=function(b){}",
			"a = function (b) {};",
			"a = function (b) {};",
		},
		{ // 86
			"a=function *(b,c){}",
			"a = function* (b, c) {};",
			"a = function* (b, c) {};",
		},
		{ // 87
			"a=function\n(\nb\n,\nc\n){}",
			"a = function (b, c) {};",
			"a = function (b, c) {};",
		},
		{ // 88
			"try{}catch{}",
			"try {} catch {}",
			"try {} catch {}",
		},
		{ // 89
			"try\n{\n}\ncatch\n{\n}",
			"try {} catch {}",
			"try {} catch {}",
		},
		{ // 90
			"try{}catch(a){}",
			"try {} catch (a) {}",
			"try {} catch (a) {}",
		},
		{ // 91
			"try\n{\n}\ncatch\n(\na\n)\n{\n}",
			"try {} catch (a) {}",
			"try {} catch (a) {}",
		},
		{ // 92
			"try{}catch({}){}",
			"try {} catch ({}) {}",
			"try {} catch ({}) {}",
		},
		{ // 93
			"try{}catch([]){}",
			"try {} catch ([]) {}",
			"try {} catch ([]) {}",
		},
		{ // 94
			"try{}finally{}",
			"try {} finally {}",
			"try {} finally {}",
		},
		{ // 95
			"try\n{\n}\nfinally\n{\n}",
			"try {} finally {}",
			"try {} finally {}",
		},
		{ // 96
			"try{}catch{}finally{}",
			"try {} catch {} finally {}",
			"try {} catch {} finally {}",
		},
		{ // 97
			"try\n{\n}\ncatch\n{\n}\nfinally\n{\n}",
			"try {} catch {} finally {}",
			"try {} catch {} finally {}",
		},
		{ // 98
			"try{}catch(a){}finally{}",
			"try {} catch (a) {} finally {}",
			"try {} catch (a) {} finally {}",
		},
		{ // 99
			"try\n{\n}\ncatch\n(\na\n)\n{\n}\nfinally\n{\n}",
			"try {} catch (a) {} finally {}",
			"try {} catch (a) {} finally {}",
		},
		{ // 100
			"class a{}",
			"class a {}",
			"class a {}",
		},
		{ // 101
			"class\na\n{\n}\n",
			"class a {}",
			"class a {}",
		},
		{ // 102
			"class a extends b {}",
			"class a extends b {}",
			"class a extends b {}",
		},
		{ // 103
			"class\na\nextends\nb\n{\n}",
			"class a extends b {}",
			"class a extends b {}",
		},
		{ // 104
			"a = class{}",
			"a = class {};",
			"a = class {};",
		},
		{ // 105
			"a\n=\nclass\nb\n{\n}",
			"a = class b {};",
			"a = class b {};",
		},
		{ // 106
			"a\n=\nclass\nextends\nb\n{\n}",
			"a = class extends b {};",
			"a = class extends b {};",
		},
		{ // 107
			"let a = 1",
			"let a = 1;",
			"let a = 1;",
		},
		{ // 108
			"let\na\n=\n1\n",
			"let a = 1;",
			"let a = 1;",
		},
		{ // 109
			"let a=1,b=2,c=3",
			"let a = 1, b = 2, c = 3;",
			"let a = 1,\nb = 2,\nc = 3;",
		},
		{ // 110
			"const a = 1",
			"const a = 1;",
			"const a = 1;",
		},
		{ // 111
			"const\na\n=\n1\n",
			"const a = 1;",
			"const a = 1;",
		},
		{ // 112
			"const a=1,b=2,c=3",
			"const a = 1, b = 2, c = 3;",
			"const a = 1,\nb = 2,\nc = 3;",
		},
		{ // 113
			"let a",
			"let a;",
			"let a;",
		},
		{ // 114
			"let\na\n,\nb\n=\n1\n,\nc\n",
			"let a, b = 1, c;",
			"let a,\nb = 1,\nc;",
		},
		{ // 115
			"const a",
			"const a;",
			"const a;",
		},
		{ // 116
			"const\na\n,\nb\n=\n1\n,\nc\n",
			"const a, b = 1, c;",
			"const a,\nb = 1,\nc;",
		},
		{ // 117
			"let [a]=1",
			"let [a] = 1;",
			"let [a] = 1;",
		},
		{ // 118
			"const\n[\na\n]\n=\n1",
			"const [a] = 1;",
			"const [a] = 1;",
		},
		{ // 119
			"let {a}=1",
			"let {a} = 1;",
			"let {a: a} = 1;",
		},
		{ // 120
			"const\n{\na\n}\n=\n1",
			"const {a} = 1;",
			"const {a: a} = 1;",
		},
		{ // 121
			"function* a() {yield a}",
			"function* a() {\n\tyield a;\n}",
			"function* a() {\n\tyield a;\n}",
		},
		{ // 122
			"() => {}",
			"() => {};",
			"() => {};",
		},
		{ // 123
			"a=b",
			"a = b;",
			"a = b;",
		},
		{ // 124
			"a/=b",
			"a /= b;",
			"a /= b;",
		},
		{ // 125
			"a%=b",
			"a %= b;",
			"a %= b;",
		},
		{ // 126
			"a+=b",
			"a += b;",
			"a += b;",
		},
		{ // 127
			"a-=b",
			"a -= b;",
			"a -= b;",
		},
		{ // 128
			"a<<=b",
			"a <<= b;",
			"a <<= b;",
		},
		{ // 129
			"a>>=b",
			"a >>= b;",
			"a >>= b;",
		},
		{ // 130
			"a>>>=b",
			"a >>>= b;",
			"a >>>= b;",
		},
		{ // 131
			"a&=b",
			"a &= b;",
			"a &= b;",
		},
		{ // 132
			"a^=b",
			"a ^= b;",
			"a ^= b;",
		},
		{ // 133
			"a|=b",
			"a |= b;",
			"a |= b;",
		},
		{ // 134
			"a**=b",
			"a **= b;",
			"a **= b;",
		},
		{ // 135
			"a?b:c",
			"a ? b : c;",
			"a ? b : c;",
		},
		{ // 136
			"new a",
			"new a;",
			"new a;",
		},
		{ // 137
			"a()",
			"a();",
			"a();",
		},
		{ // 138
			"var {a} = 1",
			"var {a} = 1;",
			"var {a: a} = 1;",
		},
		{ // 139
			"var { a , b, ...c } = 1",
			"var {a, b, ...c} = 1;",
			"var {a: a, b: b, ...c} = 1;",
		},
		{ // 140
			"var [a] = 1",
			"var [a] = 1;",
			"var [a] = 1;",
		},
		{ // 141
			"var [ a , b, ...c ] = 1",
			"var [a, b, ...c] = 1;",
			"var [a, b, ...c] = 1;",
		},
		{ // 142
			"switch (a) {case 1:}",
			"switch (a) {\ncase 1:\n}",
			"switch (a) {\ncase 1:\n}",
		},
		{ // 143
			"switch (a) {case 1:b;c;d}",
			"switch (a) {\ncase 1:\n\tb;\n\tc;\n\td;\n}",
			"switch (a) {\ncase 1:\n\tb;\n\tc;\n\td;\n}",
		},
		{ // 144
			"function a(b){}",
			"function a(b) {}",
			"function a(b) {}",
		},
		{ // 145
			"function\na(\nb\n)\n{\n}\n",
			"function a(b) {}",
			"function a(b) {}",
		},
		{ // 146
			"function a(b,c,...d){}",
			"function a(b, c, ...d) {}",
			"function a(b, c, ...d) {}",
		},
		{ // 147
			"class\na{b(){}c\n(){}}",
			"class a {\n\tb() {}\n\tc() {}\n}",
			"class a {\n\tb() {}\n\tc() {}\n}",
		},
		{ // 148
			"class\na{*b(){}\n*\nc\n(){}}",
			"class a {\n\t* b() {}\n\t* c() {}\n}",
			"class a {\n\t* b() {}\n\t* c() {}\n}",
		},
		{ // 149
			"class\na{async b(){}\nasync c\n(){}}",
			"class a {\n\tasync b() {}\n\tasync c() {}\n}",
			"class a {\n\tasync b() {}\n\tasync c() {}\n}",
		},
		{ // 150
			"class\na{async *b(){}\nasync *\nc\n(){}}",
			"class a {\n\tasync * b() {}\n\tasync * c() {}\n}",
			"class a {\n\tasync * b() {}\n\tasync * c() {}\n}",
		},
		{ // 151
			"class\na{get\nb(){}\nget c\n(){}}",
			"class a {\n\tget b() {}\n\tget c() {}\n}",
			"class a {\n\tget b() {}\n\tget c() {}\n}",
		},
		{ // 152
			"class\na{set\nb(c){}\nset d\n(e){}}",
			"class a {\n\tset b(c) {}\n\tset d(e) {}\n}",
			"class a {\n\tset b(c) {}\n\tset d(e) {}\n}",
		},
		{ // 153
			"class\na{static\nb(){}\nstatic c\n(){}}",
			"class a {\n\tstatic b() {}\n\tstatic c() {}\n}",
			"class a {\n\tstatic b() {}\n\tstatic c() {}\n}",
		},
		{ // 154
			"class\na{static\n*b(){}\nstatic *\nc\n(){}}",
			"class a {\n\tstatic * b() {}\n\tstatic * c() {}\n}",
			"class a {\n\tstatic * b() {}\n\tstatic * c() {}\n}",
		},
		{ // 155
			"class\na{static\nasync b(){}\nstatic async c\n(){}}",
			"class a {\n\tstatic async b() {}\n\tstatic async c() {}\n}",
			"class a {\n\tstatic async b() {}\n\tstatic async c() {}\n}",
		},
		{ // 156
			"class\na{static\nasync *b(){}\nstatic async *\nc(){}}",
			"class a {\n\tstatic async * b() {}\n\tstatic async * c() {}\n}",
			"class a {\n\tstatic async * b() {}\n\tstatic async * c() {}\n}",
		},
		{ // 157
			"class\na{static\nget\nb(){}static get c\n(){}}",
			"class a {\n\tstatic get b() {}\n\tstatic get c() {}\n}",
			"class a {\n\tstatic get b() {}\n\tstatic get c() {}\n}",
		},
		{ // 158
			"class\na{static\nset\nb(c){}static set d\n(e){}}",
			"class a {\n\tstatic set b(c) {}\n\tstatic set d(e) {}\n}",
			"class a {\n\tstatic set b(c) {}\n\tstatic set d(e) {}\n}",
		},
		{ // 159
			"a",
			"a;",
			"a;",
		},
		{ // 160
			"a?b:c",
			"a ? b : c;",
			"a ? b : c;",
		},
		{ // 161
			"a\n?\nb\n:\nc\n",
			"a ? b : c;",
			"a ? b : c;",
		},
		{ // 162
			"a=>b",
			"a => b;",
			"a => b;",
		},
		{ // 163
			"a =>\nb",
			"a => b;",
			"a => b;",
		},
		{ // 164
			"async a => b",
			"async a => b;",
			"async a => b;",
		},
		{ // 165
			"(a,b)=>b",
			"(a, b) => b;",
			"(a, b) => b;",
		},
		{ // 166
			"async (a,b)=>c",
			"async (a, b) => c;",
			"async (a, b) => c;",
		},
		{ // 167
			"a=>{}",
			"a => {};",
			"a => {};",
		},
		{ // 168
			"async a=>{}",
			"async a => {};",
			"async a => {};",
		},
		{ // 169
			"(a,b)=>{}",
			"(a, b) => {};",
			"(a, b) => {};",
		},
		{ // 170
			"async(a,b)=>{}",
			"async (a, b) => {};",
			"async (a, b) => {};",
		},
		{ // 171
			"new a",
			"new a;",
			"new a;",
		},
		{ // 172
			"new\nnew \n new	new\na",
			"new new new new a;",
			"new new new new a;",
		},
		{ // 173
			"super()",
			"super();",
			"super();",
		},
		{ // 174
			"super\n()",
			"super();",
			"super();",
		},
		{ // 175
			"import(a)",
			"import(a);",
			"import(a);",
		},
		{ // 176
			"import\n(a)",
			"import(a);",
			"import(a);",
		},
		{ // 177
			"a()",
			"a();",
			"a();",
		},
		{ // 178
			"a\n()",
			"a();",
			"a();",
		},
		{ // 179
			"a\n()\n()",
			"a()();",
			"a()();",
		},
		{ // 180
			"a()[b]",
			"a()[b];",
			"a()[b];",
		},
		{ // 181
			"a().b",
			"a().b;",
			"a().b;",
		},
		{ // 182
			"a()`b`",
			"a()`b`;",
			"a()`b`;",
		},
		{ // 183
			"var{a}=b",
			"var {a} = b;",
			"var {a: a} = b;",
		},
		{ // 184
			"var\n{\na\n:\nb\n}\n=\nc\n",
			"var {a: b} = c;",
			"var {a: b} = c;",
		},
		{ // 185
			"var[a=b]=c",
			"var [a = b] = c;",
			"var [a = b] = c;",
		},
		{ // 186
			"var[a=[b]] = c",
			"var [a = [b]] = c;",
			"var [a = [b]] = c;",
		},
		{ // 187
			"var[a={b}] = c",
			"var [a = {b}] = c;",
			"var [a = {b: b}] = c;",
		},
		{ // 188
			"var a={[\"b\"]:c}",
			"var a = {[\"b\"]: c};",
			"var a = {[\"b\"]: c};",
		},
		{ // 189
			"a||b",
			"a || b;",
			"a || b;",
		},
		{ // 190
			"(a,b,c)",
			"(a, b, c);",
			"(a, b, c);",
		},
		{ // 191
			"(\na\n,\nb\n,\nc\n)\n",
			"(a, b, c);",
			"(a, b, c);",
		},
		{ // 192
			"var a=(b,c,...d)=>{}",
			"var a = (b, c, ...d) => {};",
			"var a = (b, c, ...d) => {};",
		},
		{ // 193
			"var a=(b,c,...[e])=>{}",
			"var a = (b, c, ...[e]) => {};",
			"var a = (b, c, ...[e]) => {};",
		},
		{ // 194
			"var a=(b,c,...{...e})=>{}",
			"var a = (b, c, ...{...e}) => {};",
			"var a = (b, c, ...{...e}) => {};",
		},
		{ // 195
			"new a()",
			"new a();",
			"new a();",
		},
		{ // 196
			"new\nnew\na\n(\n)\n(\n)\n",
			"new new a()();",
			"new new a()();",
		},
		{ // 197
			"a\n[\n1\n]\n",
			"a[1];",
			"a[1];",
		},
		{ // 198
			"a\n.\nb\n",
			"a.b;",
			"a\n.b;",
		},
		{ // 199
			"a\n`b`",
			"a`b`;",
			"a`b`;",
		},
		{ // 200
			"new\nsuper\n[\na\n]\n[\nb\n]\n.\nc`d`\n(\nnew\n.\ntarget\n)\n",
			"new super[a][b].c`d`(new.target);",
			"new super[a][b]\n.c`d`(new.target);",
		},
		{ // 201
			"a(b,c,...d)",
			"a(b, c, ...d);",
			"a(b, c, ...d);",
		},
		{ // 202
			"a\n(\n...\nb\n)\n",
			"a(...b);",
			"a(...b);",
		},
		{ // 203
			"`a`",
			"`a`;",
			"`a`;",
		},
		{ // 204
			"`a${b}c`",
			"`a${b}c`;",
			"`a${b}c`;",
		},
		{ // 205
			"`a${\nb\n}c${\nd\n}e`",
			"`a${b}c${d}e`;",
			"`a${b}c${d}e`;",
		},
		{ // 206
			"{\n`a`\n}",
			"{\n\t`a`;\n}",
			"{\n\t`a`;\n}",
		},
		{ // 207
			"{\n`a\nb`\n}",
			"{\n\t`a\nb`;\n}",
			"{\n\t`a\nb`;\n}",
		},
		{ // 208
			"{\n`a\nb${c}d\ne${f}g\nh`\n}",
			"{\n\t`a\nb${c}d\ne${f}g\nh`;\n}",
			"{\n\t`a\nb${c}d\ne${f}g\nh`;\n}",
		},
		{ // 209
			"a&&b",
			"a && b;",
			"a && b;",
		},
		{ // 210
			"this",
			"this;",
			"this;",
		},
		{ // 211
			"a",
			"a;",
			"a;",
		},
		{ // 212
			"1",
			"1;",
			"1;",
		},
		{ // 213
			"[\n]\n",
			"[];",
			"[];",
		},
		{ // 214
			"var a={}",
			"var a = {};",
			"var a = {};",
		},
		{ // 215
			"var a=function(){}",
			"var a = function () {};",
			"var a = function () {};",
		},
		{ // 216
			"var a=class{}",
			"var a = class {};",
			"var a = class {};",
		},
		{ // 217
			"`a`",
			"`a`;",
			"`a`;",
		},
		{ // 218
			"(a)",
			"(a);",
			"(a);",
		},
		{ // 219
			"a|b",
			"a | b;",
			"a | b;",
		},
		{ // 220
			"[a,b,...c]",
			"[a, b, ...c];",
			"[a, b, ...c];",
		},
		{ // 221
			"[...a]",
			"[...a];",
			"[...a];",
		},
		{ // 222
			"[a]",
			"[a];",
			"[a];",
		},
		{ // 223
			"var a={b:c}",
			"var a = {b: c};",
			"var a = {b: c};",
		},
		{ // 224
			"var a={b:c,d:e}",
			"var a = {b: c, d: e};",
			"var a = {b: c, d: e};",
		},
		{ // 225
			"var a={\nb\n:\nc\n}",
			"var a = {b: c};",
			"var a = {b: c};",
		},
		{ // 226
			"var a={\nb\n:\nc\n,\nd\n:\ne\n}",
			"var a = {b: c, d: e};",
			"var a = {b: c, d: e};",
		},
		{ // 227
			"a^b",
			"a ^ b;",
			"a ^ b;",
		},
		{ // 228
			"var a\n=\n{\nb\n:\nc\n,\nd\n,\ne\n=\nf\n,\ng\n(\n)\n{\n}\n,\n...\nh\n}\n",
			"var a = {b: c, d, e = f, g() {}, ...h};",
			"var a = {b: c, d: d, e = f, g() {}, ...h};",
		},
		{ // 229
			"a&b",
			"a & b;",
			"a & b;",
		},
		{ // 230
			"a==b",
			"a == b;",
			"a == b;",
		},
		{ // 231
			"a!=b",
			"a != b;",
			"a != b;",
		},
		{ // 232
			"a===b",
			"a === b;",
			"a === b;",
		},
		{ // 233
			"a!==b",
			"a !== b;",
			"a !== b;",
		},
		{ // 234
			"a<b",
			"a < b;",
			"a < b;",
		},
		{ // 235
			"a>b",
			"a > b;",
			"a > b;",
		},
		{ // 236
			"a<=b",
			"a <= b;",
			"a <= b;",
		},
		{ // 237
			"a>=b",
			"a >= b;",
			"a >= b;",
		},
		{ // 238
			"a instanceof b",
			"a instanceof b;",
			"a instanceof b;",
		},
		{ // 239
			"a in b",
			"a in b;",
			"a in b;",
		},
		{ // 240
			"a<<b",
			"a << b;",
			"a << b;",
		},
		{ // 241
			"a>>b",
			"a >> b;",
			"a >> b;",
		},
		{ // 242
			"a>>>b",
			"a >>> b;",
			"a >>> b;",
		},
		{ // 243
			"a+b",
			"a + b;",
			"a + b;",
		},
		{ // 244
			"a-b",
			"a - b;",
			"a - b;",
		},
		{ // 245
			"a*b",
			"a * b;",
			"a * b;",
		},
		{ // 246
			"a/b",
			"a / b;",
			"a / b;",
		},
		{ // 247
			"a%b",
			"a % b;",
			"a % b;",
		},
		{ // 248
			"a**b",
			"a ** b;",
			"a ** b;",
		},
		{ // 249
			"delete a",
			"delete a;",
			"delete a;",
		},
		{ // 250
			"void a",
			"void a;",
			"void a;",
		},
		{ // 251
			"typeof a",
			"typeof a;",
			"typeof a;",
		},
		{ // 252
			"+\na",
			"+a;",
			"+a;",
		},
		{ // 253
			"-\na",
			"-a;",
			"-a;",
		},
		{ // 254
			"~\na",
			"~a;",
			"~a;",
		},
		{ // 255
			"!\na",
			"!a;",
			"!a;",
		},
		{ // 256
			"async function a(){await b}",
			"async function a() {\n\tawait b;\n}",
			"async function a() {\n\tawait b;\n}",
		},
		{ // 257
			"a ++",
			"a++;",
			"a++;",
		},
		{ // 258
			"a --",
			"a--;",
			"a--;",
		},
		{ // 259
			"++\na",
			"++a;",
			"++a;",
		},
		{ // 260
			"--\na",
			"--a;",
			"--a;",
		},
		{ // 261
			"a: function b(){}",
			"a: function b() {}",
			"a: function b() {}",
		},
		{ // 262
			"a: b",
			"a: b;",
			"a: b;",
		},
		{ // 263
			"continue a",
			"continue a;",
			"continue a;",
		},
		{ // 264
			"debugger",
			"debugger;",
			"debugger;",
		},
		{ // 265
			"for(var a,b,\nc;;){}",
			"for (var a, b, c;;) {}",
			"for (var a, b,\n\tc;;) {}",
		},
		{ // 266
			"for(var{a}in b){}",
			"for (var {a} in b) {}",
			"for (var {a: a} in b) {}",
		},
		{ // 267
			"for(var[a]in b){}",
			"for (var [a] in b) {}",
			"for (var [a] in b) {}",
		},
		{ // 268
			"switch(a){default:b}",
			"switch (a) {\ndefault:\n\tb;\n}",
			"switch (a) {\ndefault:\n\tb;\n}",
		},
		{ // 269
			"function*a(){yield *b}",
			"function* a() {\n\tyield * b;\n}",
			"function* a() {\n\tyield * b;\n}",
		},
		{ // 270
			"a*=b",
			"a *= b;",
			"a *= b;",
		},
		{ // 271
			"var[[a]]=b",
			"var [[a]] = b;",
			"var [[a]] = b;",
		},
		{ // 272
			"var[{a}]=b",
			"var [{a}] = b;",
			"var [{a: a}] = b;",
		},
		{ // 273
			"super\n.\na\n",
			"super.a;",
			"super.a;",
		},
		{ // 274
			"a\n?.\nb",
			"a?.b;",
			"a?.b;",
		},
		{ // 275
			"a\n??\nb",
			"a ?? b;",
			"a ?? b;",
		},
		{ // 276
			"a\n??\nb\n??\nc",
			"a ?? b ?? c;",
			"a ?? b ?? c;",
		},
		{ // 277
			"a = ([b]) => b",
			"a = ([b]) => b;",
			"a = ([b]) => b;",
		},
		{ // 278
			"a?.b().c",
			"a?.b().c;",
			"a?.b().c;",
		},
		{ // 279
			"a?.b()?.c",
			"a?.b()?.c;",
			"a?.b()?.c;",
		},
		{ // 280
			"a&&=1",
			"a &&= 1;",
			"a &&= 1;",
		},
		{ // 281
			"a||=1",
			"a ||= 1;",
			"a ||= 1;",
		},
		{ // 282
			"a??=1",
			"a ??= 1;",
			"a ??= 1;",
		},
		{ // 283
			"[a, b] = [b, a]",
			"[a, b] = [b, a];",
			"[a, b] = [b, a];",
		},
		{ // 284
			"[a.b, a.c] = [a.c, a.b]",
			"[a.b, a.c] = [a.c, a.b];",
			"[a.b, a.c] = [a.c, a.b];",
		},
		{ // 285
			"{a}",
			"{\n\ta;\n}",
			"{\n\ta;\n}",
		},
		{ // 286
			"{a;b}",
			"{\n\ta;\n\tb;\n}",
			"{\n\ta;\n\tb;\n}",
		},
		{ // 287
			"{a;\nb}",
			"{\n\ta;\n\tb;\n}",
			"{\n\ta;\n\tb;\n}",
		},
		{ // 288
			"{\na;\nb\n}",
			"{\n\ta;\n\tb;\n}",
			"{\n\ta;\n\tb;\n}",
		},
		{ // 289
			"({a, b} = {a: 1, b: 2})",
			"({a, b} = {a: 1, b: 2});",
			"({a: a, b: b} = {a: 1, b: 2});",
		},
		{ // 290
			"[a,b,...c] = [b, a]",
			"[a, b, ...c] = [b, a];",
			"[a, b, ...c] = [b, a];",
		},
		{ // 291
			"({a,b,...c}=d)",
			"({a, b, ...c} = d);",
			"({a: a, b: b, ...c} = d);",
		},
		{ // 292
			"({a:{b,c: d,...e},...f}=g)",
			"({a: {b, c: d, ...e}, ...f} = g);",
			"({a: {b: b, c: d, ...e}, ...f} = g);",
		},
		{ // 293
			"[a, ,[b,{c},,...d],,...e]=f",
			"[a, , [b, {c}, , ...d], , ...e] = f;",
			"[a, , [b, {c: c}, , ...d], , ...e] = f;",
		},
		{ // 294
			"a() ?.\nb",
			"a()?.b;",
			"a()?.b;",
		},
		{ // 295
			"a ?. [1]",
			"a?.[1];",
			"a?.[1];",
		},
		{ // 296
			"a ?. `1`",
			"a?.`1`;",
			"a?.`1`;",
		},
		{ // 297
			"a()\n.then()\n.catch()",
			"a().then().catch();",
			"a()\n.then()\n.catch();",
		},
		{ // 298
			"[a=b]=c",
			"[a = b] = c;",
			"[a = b] = c;",
		},
		{ // 299
			"({a=b}=c)",
			"({a = b} = c);",
			"({a: a = b} = c);",
		},
		{ // 300
			"a.#b",
			"a.#b;",
			"a.#b;",
		},
		{ // 301
			"a\n.#b",
			"a.#b;",
			"a\n.#b;",
		},
		{ // 302
			"a.#b.c",
			"a.#b.c;",
			"a.#b.c;",
		},
		{ // 303
			"a\n.#b\n.c",
			"a.#b.c;",
			"a\n.#b\n.c;",
		},
		{ // 304
			"class\na\n{\nb\n}",
			"class a {\n\tb;\n}",
			"class a {\n\tb;\n}",
		},
		{ // 305
			"class a { b () {} }",
			"class a {\n\tb() {}\n}",
			"class a {\n\tb() {}\n}",
		},
		{ // 306
			"class\na\n{\n#b\n}",
			"class a {\n\t#b;\n}",
			"class a {\n\t#b;\n}",
		},
		{ // 307
			"class a { #b () {} }",
			"class a {\n\t#b() {}\n}",
			"class a {\n\t#b() {}\n}",
		},
		{ // 308
			"class a { #b = 1 }",
			"class a {\n\t#b = 1;\n}",
			"class a {\n\t#b = 1;\n}",
		},
		{ // 309
			"class a { #b = 1; #c = 2 }",
			"class a {\n\t#b = 1;\n\t#c = 2;\n}",
			"class a {\n\t#b = 1;\n\t#c = 2;\n}",
		},
		{ // 310
			"class a { #b = 1\n#c = 2 }",
			"class a {\n\t#b = 1;\n\t#c = 2;\n}",
			"class a {\n\t#b = 1;\n\t#c = 2;\n}",
		},
		{ // 311
			"class a { #b(){}#c = 2 }",
			"class a {\n\t#b() {}\n\t#c = 2;\n}",
			"class a {\n\t#b() {}\n\t#c = 2;\n}",
		},
		{ // 312
			"class a { #b\n#c(){}}",
			"class a {\n\t#b;\n\t#c() {}\n}",
			"class a {\n\t#b;\n\t#c() {}\n}",
		},
		{ // 313
			"class a { #b = 1\n#c(){}}",
			"class a {\n\t#b = 1;\n\t#c() {}\n}",
			"class a {\n\t#b = 1;\n\t#c() {}\n}",
		},
		{ // 314
			"class a { #b = 1;#c(){}}",
			"class a {\n\t#b = 1;\n\t#c() {}\n}",
			"class a {\n\t#b = 1;\n\t#c() {}\n}",
		},
		{ // 315
			"class a { #b;#c(){}}",
			"class a {\n\t#b;\n\t#c() {}\n}",
			"class a {\n\t#b;\n\t#c() {}\n}",
		},
		{ // 316
			"class a {static a;static b\nstatic c = 2;static d(){};static{}static{e}static{e;f}}",
			"class a {\n\tstatic a;\n\tstatic b;\n\tstatic c = 2;\n\tstatic d() {}\n\tstatic {}\n\tstatic {\n\t\te;\n\t}\n\tstatic {\n\t\te;\n\t\tf;\n\t}\n}",
			"class a {\n\tstatic a;\n\tstatic b;\n\tstatic c = 2;\n\tstatic d() {}\n\tstatic {}\n\tstatic {\n\t\te;\n\t}\n\tstatic {\n\t\te;\n\t\tf;\n\t}\n}",
		},
		{ // 317
			"#a in b",
			"#a in b;",
			"#a in b;",
		},
		{ // 318
			"#a\nin\nb",
			"#a in b;",
			"#a in b;",
		},
		{ // 319
			"a().#b",
			"a().#b;",
			"a().#b;",
		},
		{ // 320
			"a\n(\n)\n.\n#b",
			"a().#b;",
			"a()\n.#b;",
		},
		{ // 321
			"a\n?.\n#b",
			"a?.#b;",
			"a?.#b;",
		},
		{ // 322
			"a?.c.#d",
			"a?.c.#d;",
			"a?.c.#d;",
		},
		{ // 323
			"// A\n// B\n\na();\n// C\n// D\n",
			"a();",
			"// A\n// B\n\na();\n// C\n// D\n",
		},
		{ // 324
			"/* A *//* B */\n// C\n\na();\n// D\n/* E */   /* F */\n",
			"a();",
			"/* A */ /* B */\n// C\n\na();\n// D\n/* E */ /* F */",
		},
		{ // 325
			"/*\nA\n*//* B */\n// C\n\na();\n// D\n/* E */   /*\n\nF\n\n*/\n",
			"a();",
			"/*\nA\n*/ /* B */\n// C\n\na();\n// D\n/* E */ /*\n\nF\n\n*/",
		},
		{ // 326
			"/*\nA\n*//* B */\n// C\n\n// D\na(); // E\n// F\n\n// G\n/* H */   /*\n\nI\n\n*/\n",
			"a();",
			"/*\nA\n*/ /* B */\n// C\n\n// D\na(); // E\n     // F\n\n// G\n/* H */ /*\n\nI\n\n*/",
		},
		{ // 327
			"// A\n\n// B\nsuper // C\n[ // D\n1\n // E\n]// F\n",
			"super[1];",
			"// A\n\n// B\nsuper // C\n[ // D\n\n\t1\n// E\n] // F\n",
		},
		{ // 328
			"// A\n\n// B\nsuper /* C */ . /* D */ a // E\n",
			"super.a;",
			"// A\n\n// B\nsuper /* C */ . /* D */ a // E\n",
		},
		{ // 329
			"// A\n\n// B\nnew /* C */./* D */target /* E */",
			"new.target;",
			"// A\n\n// B\nnew /* C */ . /* D */ target /* E */;",
		},
		{ // 330
			"// A\n\n/* B */import/* C */./* D */meta/* E */",
			"import.meta;",
			"// A\n\n/* B */ import /* C */ . /* D */ meta /* E */;",
		},
		{ // 331
			"// A\n\n// B\nnew/* C */1/* D */() // E\n",
			"new 1();",
			"// A\n\n// B\nnew /* C */ 1 /* D */ () // E\n",
		},
		{ // 332
			"// A\n\n// B\na // C\n",
			"a;",
			"// A\n\n// B\na // C\n",
		},
		{ // 333
			"// A\n\n// B\na /* C */``/* D */",
			"a``;",
			"// A\n\n// B\na /* C */ `` /* D */;",
		},
		{ // 334
			"// A\n\n/* B */a/* C */./* D */#b/* E */",
			"a.#b;",
			"// A\n\n/* B */ a /* C */ . /* D */ #b /* E */;",
		},
		{ // 335
			"// A\n\n/* B */a/* C */./* D */#b/* E */./* F */c // G\n",
			"a.#b.c;",
			"// A\n\n/* B */ a /* C */ . /* D */ #b /* E */ . /* F */ c // G\n",
		},
		{ // 336
			"// A\n\n// B\na /* C */ . /* D */ #b /* E */ [ // F\n\"c\"\n// G\n] /* H */",
			"a.#b[\"c\"];",
			"// A\n\n// B\na /* C */ . /* D */ #b /* E */ [ // F\n\n\t\"c\"\n// G\n] /* H */;",
		},
		{ // 337
			"super[ // C\n\n// D\n1 // E\n\n// F\n]",
			"super[1];",
			"super[ // C\n\n\t// D\n\t1 // E\n\n// F\n];",
		},
		{ // 338
			"// A\n\n// B\na /* C */ . /* D */ #b /* E */ [ // F\n\n// G\n\"c\" // H\n\n// I\n] /* J */",
			"a.#b[\"c\"];",
			"// A\n\n// B\na /* C */ . /* D */ #b /* E */ [ // F\n\n\t// G\n\t\"c\" // H\n\n// I\n] /* J */;",
		},
		{ // 339
			"a( // A\n\n// B\n)",
			"a();",
			"a( // A\n\n// B\n);",
		},
		{ // 340
			"a( // A\n\n// B\nb // C\n\n// D\n, // E\nc // F\n\n//G\n)",
			"a(b, c);",
			"a( // A\n\n\t// B\n\tb // C\n\n\t// D\n\t, // E\n\tc // F\n\n//G\n);",
		},
		{ // 341
			"a( // A\n\n// B\nb // C\n\n// D\n, // E\n...// F\nc // G\n\n//H\n)",
			"a(b, ...c);",
			"a( // A\n\n\t// B\n\tb // C\n\n\t// D\n\t, // E\n\t... // F\n\tc // G\n\n//H\n);",
		},
		{ // 342
			"( // A\n\n// B\na // C\n\n// D\n, // E\nb // F\n\n//G\n)",
			"(a, b);",
			"( // A\n\n\t// B\n\ta // C\n\n\t// D\n\t, // E\n\tb // F\n\n//G\n);",
		},
		{ // 343
			"[\n// A\n...// B\na // C\n]",
			"[...a];",
			"[\n\t// A\n\t... // B\n\ta // C\n];",
		},
		{ // 344
			"[ // A\n// B\n\n// C\na // D\n\n// E\n, // F\n\n// G\nb // H\n// I\n\n// J\n]",
			"[a, b];",
			"[ // A\n  // B\n\n\t// C\n\ta // D\n\n\t// E\n\t, // F\n\n\t// G\n\tb // H\n\t  // I\n\n// J\n];",
		},
		{ // 345
			"var a = {[ // A\n\n// B\nb // C\n\n// D\n]:c}",
			"var a = {[b]: c};",
			"var a = {[ // A\n\n\t// B\n\tb // C\n\n// D\n]: c};",
		},
		{ // 346
			"var a = {[ // A\n\nb]:c}",
			"var a = {[b]: c};",
			"var a = {[ // A\n\n\tb]: c};",
		},
		{ // 347
			"() => { // A\n\n// B\na // C\n\n// D\n}",
			"() => {\n\ta;\n};",
			"() => { // A\n\n\t// B\n\ta // C\n\n// D\n};",
		},
		{ // 348
			"let [\n// A\na// B\n] = b;",
			"let [a] = b;",
			"let [\n\t// A\n\ta // B\n] = b;",
		},
		{ // 349
			"let {\n// A\na// B\n} = b;",
			"let {a} = b;",
			"let {\n\ta: // A\n\ta // B\n} = b;",
		},
		{ // 350
			"let {\n// A\na// B\n: // C\nb // D\n} = c;",
			"let {a: b} = c;",
			"let {\n\t// A\n\ta // B\n\t: // C\n\tb // D\n} = c;",
		},
		{ // 351
			"let {\n// A\na// B\n= // C\nb // D\n} = c;",
			"let {a = b} = c;",
			"let {\n\ta: // A\n\ta // B\n\t= // C\n\tb // D\n} = c;",
		},
		{ // 352
			"let {\n// A\na// B\n: // C\nb // D\n= // E\nc // F\n} = d;",
			"let {a: b = c} = d;",
			"let {\n\t// A\n\ta // B\n\t: // C\n\tb // D\n\t= // E\n\tc // F\n} = d;",
		},
		{ // 353
			"let { // A\n\n// B\na // C\n\n// D\n} = b",
			"let {a} = b;",
			"let { // A\n\n\ta: // B\n\ta // C\n\n// D\n} = b;",
		},
		{ // 354
			"let { // A\n\n// B\n...// C\na // D\n// E\n\n//F\n} = b",
			"let {...a} = b;",
			"let { // A\n\n\t// B\n\t... // C\n\ta // D\n\t  // E\n\n//F\n} = b;",
		},
		{ // 355
			"let { // A\n\n// B\na // C\n, // D\n...// E\nb // F\n\n// G\n} = c",
			"let {a, ...b} = c;",
			"let { // A\n\n\ta: // B\n\ta // C\n\t,\n\t// D\n\t... // E\n\tb // F\n\n// G\n} = c;",
		},
		{ // 356
			"let [ // A\n\n// B\n] = a",
			"let [] = a;",
			"let [ // A\n\n\n// B\n] = a;",
		},
		{ // 357
			"let [ // A\n\n// B\n...// C\na // D\n\n// E\n] = b",
			"let [...a] = b;",
			"let [ // A\n\n\t// B\n\t... // C\n\ta // D\n\n// E\n] = b;",
		},
		{ // 358
			"let [ // A\n\n// B\na // C\n, // D\n\n// E\n, // F\n... // G\nb // H\n\n// I\n] = c",
			"let [a, , ...b] = c;",
			"let [ // A\n\n\t// B\n\ta // C\n\t,\n\t// D\n\n\t// E\n\t,\n\t// F\n\t... // G\n\tb // H\n\n// I\n] = c;",
		},
		{ // 359
			"function a( // A\n\n// B\n){}",
			"function a() {}",
			"function a( // A\n\n// B\n) {}",
		},
		{ // 360
			"function a( // A\n\n// B\n... // C\nb // D\n\n// E\n){}",
			"function a(...b) {}",
			"function a( // A\n\n\t// B\n\t... // C\n\tb // D\n\n// E\n) {}",
		},
		{ // 361
			"function a( // A\n\n// B\nb // C\n\n// D\n){}",
			"function a(b) {}",
			"function a( // A\n\n\t// B\n\tb // C\n\n// D\n) {}",
		},
		{ // 362
			"function a( // A\n\n// B\nb // C\n, // D\nc // E\n\n// F\n){}",
			"function a(b, c) {}",
			"function a( // A\n\n\t// B\n\tb // C\n\t, // D\n\tc // E\n\n// F\n) {}",
		},
		{ // 363
			"function a( // A\n\n// B\nb // C\n, // D\n... // E\nc // F\n\n// G\n){}",
			"function a(b, ...c) {}",
			"function a( // A\n\n\t// B\n\tb // C\n\t, // D\n\t... // E\n\tc // F\n\n// G\n) {}",
		},
		{ // 364
			"function a( // A\n\n// B\n... // C\n[]// D\n\n// E\n){}",
			"function a(...[]) {}",
			"function a( // A\n\n\t// B\n\t... // C\n\t[] // D\n\n// E\n) {}",
		},
		{ // 365
			"function a( // A\n\n// B\n... // C\n{}// D\n\n// E\n){}",
			"function a(...{}) {}",
			"function a( // A\n\n\t// B\n\t... // C\n\t{} // D\n\n// E\n) {}",
		},
		{ // 366
			"class a {\n// A\nb /* B */(){}\n}",
			"class a {\n\tb() {}\n}",
			"class a {\n\t// A\n\tb /* B */() {}\n}",
		},
		{ // 367
			"class a {\n// A\na /* B */(){}\n/* C */ b// D\n(){}\n}",
			"class a {\n\ta() {}\n\tb() {}\n}",
			"class a {\n\t// A\n\ta /* B */() {}\n\t/* C */ b // D\n\t() {}\n}",
		},
		{ // 368
			"class a {static //A\nb /* B */() {} }",
			"class a {\n\tstatic b() {}\n}",
			"class a {\n\tstatic //A\n\tb /* B */() {}\n}",
		},
		{ // 369
			"class a {static /* A */ [\"b\"]// B\n() {} }",
			"class a {\n\tstatic [\"b\"]() {}\n}",
			"class a {\n\tstatic /* A */ [\"b\"] // B\n\t() {}\n}",
		},
		{ // 370
			"class a {static // A\n#b// B\n() {} }",
			"class a {\n\tstatic #b() {}\n}",
			"class a {\n\tstatic // A\n\t#b // B\n\t() {}\n}",
		},
		{ // 371
			"class a {static // A\nb// B\n}",
			"class a {\n\tstatic b;\n}",
			"class a {\n\tstatic // A\n\tb // B\n}",
		},
		{ // 372
			"class a {static // A\nb // B\n= 1}",
			"class a {\n\tstatic b = 1;\n}",
			"class a {\n\tstatic // A\n\tb // B\n\t= 1;\n}",
		},
		{ // 373
			"class a {static // A\n[b] // B\n}",
			"class a {\n\tstatic [b];\n}",
			"class a {\n\tstatic // A\n\t[b] // B\n}",
		},
		{ // 374
			"// A\n\n// B\nnew // C\na // D\n",
			"new a;",
			"// A\n\n// B\nnew // C\na // D\n",
		},
		{ // 375
			"// A\n\n// B\nnew // C\na // D\n() // E\n",
			"new a();",
			"// A\n\n// B\nnew // C\na // D\n() // E\n",
		},
		{ // 376
			"// A\n\n// B\nnew // C\nnew // D\na // E\n() // F\n",
			"new new a();",
			"// A\n\n// B\nnew // C\nnew // D\na // E\n() // F\n",
		},
		{ // 377
			"({\n// A\nget // B\na // C\n( // D\n\n// E\n) // F\n{} // G\n})",
			"({get a() {}});",
			"({\n\t// A\n\tget // B\n\ta // C\n\t( // D\n\n\t// E\n\t) // F\n\t{} // G\n\n});",
		},
		{ // 378
			"({\n// A\nset // B\na // C\n( // D\n\n// E\nb // F\n\n// G\n) // H\n{} // I\n})",
			"({set a(b) {}});",
			"({\n\t// A\n\tset // B\n\ta // C\n\t( // D\n\n\t\t// E\n\t\tb // F\n\n\t// G\n\t) // H\n\t{} // I\n\n});",
		},
		{ // 379
			"({\n// A\na // B\n( // C\n\n// D\nb // E\n\n// F\n) // G\n{} // H\n})",
			"({a(b) {}});",
			"({\n\t// A\n\ta // B\n\t( // C\n\n\t\t// D\n\t\tb // E\n\n\t// F\n\t) // G\n\t{} // H\n\n});",
		},
		{ // 380
			"({\n// A\nasync /* B */ a // C\n( // D\n\n// E\n) // F\n{} // G\n})",
			"({async a() {}});",
			"({\n\t// A\n\tasync /* B */ a // C\n\t( // D\n\n\t// E\n\t) // F\n\t{} // G\n\n});",
		},
		{ // 381
			"({\n// A\n* // B\na // C\n() {}})",
			"({* a() {}});",
			"({\n\t// A\n\t* // B\n\ta // C\n\t() {}\n});",
		},
		{ // 382
			"({\n// A\nasync /* B*/ * // C\na(){}})",
			"({async * a() {}});",
			"({\n\t// A\n\tasync /* B*/ * // C\n\ta() {}\n});",
		},
		{ // 383
			"({\n// A\nasync // B\n(){}})",
			"({async() {}});",
			"({\n\t// A\n\tasync // B\n\t() {}\n});",
		},
		{ // 384
			"({\n// A\nget // B\n(){}})",
			"({get() {}});",
			"({\n\t// A\n\tget // B\n\t() {}\n});",
		},
		{ // 385
			"({\n// A\n... // B\na // C\n})",
			"({...a});",
			"({\n\t// A\n\t... // B\n\ta // C\n\n});",
		},
		{ // 386
			"({\n// A\na // B\n,})",
			"({a});",
			"({\n\t// A\n\ta // B\n\t: a\n});",
		},
		{ // 387
			"({\n// A\na // B\n= // C\nb // D\n})",
			"({a = b});",
			"({\n\t// A\n\ta // B\n\t= // C\n\tb // D\n\n});",
		},
		{ // 388
			"({\n// A\na // B\n: // C\nb // D\n})",
			"({a: b});",
			"({\n\t// A\n\ta // B\n\t: // C\n\tb // D\n\n});",
		},
		{ // 389
			"({\n// A\n[ // B\na\n// C\n] // D\n: // E\nb // F\n})",
			"({[a]: b});",
			"({\n\t// A\n\t[ // B\n\n\t\ta\n\t// C\n\t] // D\n\t: // E\n\tb // F\n\n});",
		},
		{ // 390
			"({ // A\n// B\n\n// C\n})",
			"({});",
			"({ // A\n   // B\n\n// C\n});",
		},
		{ // 391
			"({ // A\n// B\n\n// C\na // D\n// E\n\n// F\n})",
			"({a});",
			"({ // A\n   // B\n\n\t// C\n\ta // D\n\t  // E\n\t: a\n\n// F\n});",
		},
		{ // 392
			"({ // A\n\n// B\na // C\n, // D\nb // E\n\n// F\n})",
			"({a, b});",
			"({ // A\n\n\t// B\n\ta // C\n\t: a,\n\t// D\n\tb // E\n\t: b\n\n// F\n});",
		},
		{ // 393
			"function // A\na // B\n()// C\n{}",
			"function a() {}",
			"function // A\na // B\n() // C\n{}",
		},
		{ // 394
			"async /* A */ function // B\na // C\n()// D\n{}",
			"async function a() {}",
			"async /* A */ function // B\na // C\n() // D\n{}",
		},
		{ // 395
			"function // A\n* // B\na // C\n()// D\n{}",
			"function* a() {}",
			"function // A\n* // B\na // C\n() // D\n{}",
		},
		{ // 396
			"(function // A\n()// B\n{})",
			"(function () {});",
			"(function // A\n() // B\n{});",
		},
		{ // 397
			"(async /* A */ function // B\n()// C\n{})",
			"(async function () {});",
			"(async /* A */ function // B\n() // C\n{});",
		},
		{ // 398
			"(function // A\n* // B\n()// C\n{})",
			"(function* () {});",
			"(function // A\n* // B\n() // C\n{});",
		},
		{ // 399
			"(async /* A */ function // B\n* // C\n()// D\n{})",
			"(async function* () {});",
			"(async /* A */ function // B\n* // C\n() // D\n{});",
		},
		{ // 400
			"`${ // A\n\n// B\na // C\n\n// D\n}`",
			"`${a}`;",
			"`${ // A\n\n// B\na // C\n\n// D\n}`;",
		},
		{ // 401
			"`${ // A\na // B\n\n// C\n}${ // D\nb // E\n}`",
			"`${a}${b}`;",
			"`${ // A\na // B\n\n// C\n}${ // D\nb // E\n}`;",
		},
		{ // 402
			"// A\n\n// B\na /* C */ ++ // D\n",
			"a++;",
			"// A\n\n// B\na /* C */++ // D\n",
		},
		{ // 403
			"// A\n\n// B\na /* C */ -- // D\n",
			"a--;",
			"// A\n\n// B\na /* C */-- // D\n",
		},
		{ // 404
			"// A\n\n// B\n++ // C\na // D\n",
			"++a;",
			"// A\n\n// B\n++ // C\na // D\n",
		},
		{ // 405
			"// A\n\n// B\n-- // C\na // D\n",
			"--a;",
			"// A\n\n// B\n-- // C\na // D\n",
		},
		{ // 406
			"// A\n\n// B\ntypeof // C\na // D\n",
			"typeof a;",
			"// A\n\n// B\ntypeof // C\na // D\n",
		},
		{ // 407
			"// A\n\n// B\nvoid // C\n+ // D\na // E\n",
			"void +a;",
			"// A\n\n// B\nvoid // C\n+ // D\na // E\n",
		},
		{ // 408
			"// A\n\n// B\nsuper // C\n()// D\n",
			"super();",
			"// A\n\n// B\nsuper // C\n() // D\n",
		},
		{ // 409
			"// A\n\n// B\nimport // C\n( // D\n\n// E\na // F\n\n// G\n) // H\n",
			"import(a);",
			"// A\n\n// B\nimport // C\n( // D\n\n\t// E\n\ta // F\n\n// G\n) // H\n",
		},
		{ // 410
			"// A\n\n// B\nsuper // C\n()// D\n`` // E",
			"super()``;",
			"// A\n\n// B\nsuper // C\n() // D\n`` // E\n",
		},
		{ // 411
			"// A\n\n// B\nsuper // C\n()// D\n()// E\n",
			"super()();",
			"// A\n\n// B\nsuper // C\n() // D\n() // E\n",
		},
		{ // 412
			"// A\n\n// B\nsuper // C\n()// D\n[ // E\n\n// F\na // G\n\n// H\n]// I",
			"super()[a];",
			"// A\n\n// B\nsuper // C\n() // D\n[ // E\n\n\t// F\n\ta // G\n\n// H\n] // I\n",
		},
		{ // 413
			"// A\n\n// B\nsuper // C\n()// D\n. // E\na // F",
			"super().a;",
			"// A\n\n// B\nsuper // C\n() // D\n. // E\na // F\n",
		},
		{ // 414
			"// A\n\n// B\nsuper // C\n()// D\n. // E\n#a // F",
			"super().#a;",
			"// A\n\n// B\nsuper // C\n() // D\n. // E\n#a // F\n",
		},
		{ // 415
			"// A\n\n// B\na // C\n()// D\n",
			"a();",
			"// A\n\n// B\na // C\n() // D\n",
		},
		{ // 416
			"// A\n\n// B\nclass a{} // C\n",
			"class a {}",
			"// A\n\n// B\nclass a {} // C\n",
		},
		{ // 417
			"// A\n\n// B\nfunction a(){} // C\n",
			"function a() {}",
			"// A\n\n// B\nfunction a() {} // C\n",
		},
		{ // 418
			"// A\n\n// B\nconst a = 1; // B\n",
			"const a = 1;",
			"// A\n\n// B\nconst a = 1; // B\n",
		},
		{ // 419
			"// A\n\n// B\nlet a = 1; // B\n",
			"let a = 1;",
			"// A\n\n// B\nlet a = 1; // B\n",
		},
		{ // 420
			"a?. // A\n() // B\n",
			"a?.();",
			"a?. // A\n() // B\n",
		},
		{ // 421
			"a?. // A\n[ // B\n\n// C\nb // D\n\n// E\n] // F\n",
			"a?.[b];",
			"a?. // A\n[ // B\n\n\t// C\n\tb // D\n\n// E\n] // F\n",
		},
		{ // 422
			"a?. // A\n[\n// C\nb\n// E\n] // F\n",
			"a?.[b];",
			"a?. // A\n[\n\t// C\n\tb\n// E\n] // F\n",
		},
		{ // 423
			"a?. // A\nb // B\n",
			"a?.b;",
			"a?. // A\nb // B\n",
		},
		{ // 424
			"a?. // A\n`` // B\n",
			"a?.``;",
			"a?. // A\n`` // B\n",
		},
		{ // 425
			"a?. // A\n()// B\n`` // C\n",
			"a?.()``;",
			"a?. // A\n() // B\n`` // C\n",
		},
		{ // 426
			"a?. // A\n`` // B\n[ // C\n\n// D\nb // E\n\n// F\n] // G\n",
			"a?.``[b];",
			"a?. // A\n`` // B\n[ // C\n\n\t// D\n\tb // E\n\n// F\n] // G\n",
		},
		{ // 427
			"a?. //A\n[ // B\n\n// C\nb // D\n\n// E\n] // F\n. // G\nc// H\n",
			"a?.[b].c;",
			"a?. //A\n[ // B\n\n\t// C\n\tb // D\n\n// E\n] // F\n. // G\nc // H\n",
		},
		{ // 428
			"class a {\n// A\nstatic // B\n{} // C\n}",
			"class a {\n\tstatic {}\n}",
			"class a {\n\t// A\n\tstatic // B\n\t{} // C\n}",
		},
		{ // 429
			"class // A\na // B\n{ // C\n\n// D\n}",
			"class a {}",
			"class // A\na // B\n{ // C\n\n// D\n}",
		},
		{ // 430
			"class a { // A\n; // B\n; // C\n; // D\na(){} // E\n; // F\n;\n // G\n}",
			"class a {\n\ta() {}\n}",
			"class a { // A\n\n\t// B\n\t// C\n\t// D\n\ta() {} // E\n\n// F\n\n// G\n}",
		},
		{ // 431
			"class a { // A\n\n// B\na(){} // C\n// D\n\n// E\nb(){} // F\n\n// G\n}",
			"class a {\n\ta() {}\n\tb() {}\n}",
			"class a { // A\n\n\t// B\n\ta() {} // C\n\t       // D\n\n\t// E\n\tb() {} // F\n\n// G\n}",
		},
		{ // 432
			"let // A\na // B\n",
			"let a;",
			"let // A\na // B\n",
		},
		{ // 433
			"let // A\na // B\n= // C\nb // D\n",
			"let a = b;",
			"let // A\na // B\n= // C\nb // D\n",
		},
		{ // 434
			"let // A\n[]// B\n = // C\na // D",
			"let [] = a;",
			"let // A\n[] // B\n= // C\na // D\n",
		},
		{ // 435
			"let // A\n{}// B\n = // C\na // D",
			"let {} = a;",
			"let // A\n{} // B\n= // C\na // D\n",
		},
		{ // 436
			"continue /* A */ a // B\n",
			"continue a;",
			"continue /* A */ a // B\n",
		},
		{ // 437
			"continue /* A */ a // B\n;",
			"continue a;",
			"continue /* A */ a // B\n",
		},
		{ // 438
			"continue /* A */;",
			"continue;",
			"continue /* A */;",
		},
		{ // 439
			"break /* A */ a // B\n",
			"break a;",
			"break /* A */ a // B\n",
		},
		{ // 440
			"break /* A */ a // B\n;",
			"break a;",
			"break /* A */ a // B\n",
		},
		{ // 441
			"break /* A */;",
			"break;",
			"break /* A */;",
		},
		{ // 442
			"function a(){\nreturn /* A */ b // B\n}",
			"function a() {\n\treturn b;\n}",
			"function a() {\n\treturn /* A */ b // B\n}",
		},
		{ // 443
			"function a(){\nreturn /* A */ b // B\n;}",
			"function a() {\n\treturn b;\n}",
			"function a() {\n\treturn /* A */ b // B\n}",
		},
		{ // 444
			"function a(){\nreturn /* A */;\n}",
			"function a() {\n\treturn;\n}",
			"function a() {\n\treturn /* A */;\n}",
		},
		{ // 445
			"debugger // A\n",
			"debugger;",
			"debugger // A\n",
		},
		{ // 446
			"debugger /* A */;",
			"debugger;",
			"debugger /* A */;",
		},
		{ // 447
			"a /* A */: // B\ndebugger;",
			"a: debugger;",
			"a /* A */ : // B\ndebugger;",
		},
		{ // 448
			"if // A\n( // B\n\n// C\na // D\n\n// E\n) // F\nb // G",
			"if (a) b;",
			"if // A\n( // B\n\n\t// C\n\ta // D\n\n// E\n) // F\nb // G\n",
		},
		{ // 449
			"// A\n\n// B\nif // C\n(// D\n\n// E\na // F\n\n// G\n) // H\n{// I\n\n// J\n} // K\nelse // L\n{ // M\n\n// N\n} // O",
			"if (a) {} else {}",
			"// A\n\n// B\nif // C\n( // D\n\n\t// E\n\ta // F\n\n// G\n) // H\n{ // I\n\n// J\n} // K\nelse // L\n{ // M\n\n// N\n} // O\n",
		},
		{ // 450
			"// A\na /* B */ || // C\nb // D\n|| c /* E */",
			"a || b || c;",
			"// A\n\na /* B */ || // C\nb // D\n|| c /* E */;",
		},
		{ // 451
			"with // A\n( // B\n\n// C\na // D\n\n// E\n) // F\nb // G\n",
			"with (a) b;",
			"with // A\n( // B\n\n\t// C\n\ta // D\n\n// E\n) // F\nb // G\n",
		},
		{ // 452
			"while // A\n( // B\n\n// C\na // D\n\n// E\n) // F\nb // G\n",
			"while (a) b;",
			"while // A\n( // B\n\n\t// C\n\ta // D\n\n// E\n) // F\nb // G\n",
		},
		{ // 453
			"switch (a) {\n// A\ncase /* B */ b /* C */: // D\n\n// E\ncase c:\n// F\ncase /* G */ d /* H */: // I\n\n// J\n{} // F\ndefault:}",
			"switch (a) {\ncase b:\ncase c:\ncase d:\n\t{}\ndefault:\n}",
			"switch (a) {\n// A\ncase /* B */ b /* C */: // D\n\n// E\ncase c:\n// F\ncase /* G */ d /* H */: // I\n\n\t// J\n\t{} // F\n\ndefault:\n}",
		},
		{ // 454
			"switch // A\n( // B\n\n// C\na // D\n\n// E\n) // F\n{ // G\n\n// H\n}",
			"switch (a) {}",
			"switch // A\n( // B\n\n\t// C\n\ta // D\n\n// E\n) // F\n{ // G\n\n// H\n}",
		},
		{ // 455
			"switch // A\n( // B\n\n// C\na // D\n\n// E\n) // F\n{ // G\n\n// H\ncase /* I */ b /* J */ : // K\n\n// J\ncase c:// L\n\n// M\ndefault /* N */ : // O\n\td;// P\n\n// Q\n}",
			"switch (a) {\ncase b:\ncase c:\ndefault:\n\td;\n}",
			"switch // A\n( // B\n\n\t// C\n\ta // D\n\n// E\n) // F\n{ // G\n\n// H\ncase /* I */ b /* J */: // K\n\n// J\ncase c: // L\n\n// M\ndefault /* N */: // O\n\n\td; // P\n\n// Q\n}",
		},
		{ // 456
			"try // A\n{} // B\ncatch // C\n{}",
			"try {} catch {}",
			"try // A\n{} // B\ncatch // C\n{}",
		},
		{ // 457
			"try {}catch // A\n( // B\n\n// C\na // D\n\n// E\n) // F\n{} // G",
			"try {} catch (a) {}",
			"try {} catch // A\n( // B\n\n\t// C\n\ta // D\n\n// E\n) // F\n{} // G\n",
		},
		{ // 458
			"try{} // A\nfinally // B\n{} // C",
			"try {} finally {}",
			"try {} // A\nfinally // B\n{} // C\n",
		},
		{ // 459
			"try // A\n{}// B\ncatch /* C */ ( // D\n\n// E\na // F\n\n// G\n) // H\n{}// I\nfinally /* J */ {} // K",
			"try {} catch (a) {} finally {}",
			"try // A\n{} // B\ncatch /* C */ ( // D\n\n\t// E\n\ta // F\n\n// G\n) // H\n{} // I\nfinally /* J */ {} // K\n",
		},
		{ // 460
			"for // A\n( // B\n\n// C\n; // D\n; // E\n\n// F\n) // G\na",
			"for (;;) a;",
			"for // A\n( // B\n\n\t// C\n\t; // D\n\t; // E\n\n// F\n) // G\na;",
		},
		{ // 461
			"for ( // A\n\n// B\nvar // C\na // D\n, // E\nb // F\n; // G\nc // H\n; // I\n\n// J\n) // K\n{}",
			"for (var a, b; c;) {}",
			"for ( // A\n\n\t// B\n\tvar // C\n\ta // D\n\t, // E\n\tb // F\n\t; // G\n\tc // H\n\t; // I\n\n// J\n) // K\n{}",
		},
		{ // 462
			"for ( // A\n\n// B\n; // C\na // D\n; // E\nb // F\n\n// G\n) // H\nc",
			"for (; a; b) c;",
			"for ( // A\n\n\t// B\n\t; // C\n\ta // D\n\t; // E\n\tb // F\n\n// G\n) // H\nc;",
		},
		{ // 463
			"for ( // A\n\n// B\nlet // C\na // D\n; b; c) {}",
			"for (let a; b; c) {}",
			"for ( // A\n\n\t// B\n\tlet // C\n\ta // D\n\t; b; c) {}",
		},
		{ // 464
			"for ( // A\n\n// B\nvar // C\n{}// D\n= // E\na // F\n;;);",
			"for (var {} = a;;) ;",
			"for ( // A\n\n\t// B\n\tvar // C\n\t{} // D\n\t= // E\n\ta // F\n\t;;) ;",
		},
		{ // 465
			"for ( // A\n\n// B\nvar // C\n[]// D\n= a // E\n; // F\n; // G\n) // H\n;",
			"for (var [] = a;;) ;",
			"for ( // A\n\n\t// B\n\tvar // C\n\t[] // D\n\t= a // E\n\t; // F\n\t; // G\n) // H\n;",
		},
		{ // 466
			"for (\n// A\nvar // B\na // C\nin // D\nb // E\n\n// F\n);",
			"for (var a in b) ;",
			"for (\n\t// A\n\tvar // B\n\ta // C\n\tin // D\n\tb // E\n\n// F\n) ;",
		},
		{ // 467
			"[ // A\n\n// B\na.b // C\n, // D\na.c // E\n\n// F\n] = b",
			"[a.b, a.c] = b;",
			"[ // A\n\n\t// B\n\ta.b // C\n\t, // D\n\ta.c // E\n\n// F\n] = b;",
		},
		{ // 468
			"[ // A\n\n// B\n... // C\na // D\n\n// E\n] = b",
			"[...a] = b;",
			"[ // A\n\n\t// B\n\t... // C\n\ta // D\n\n// E\n] = b;",
		},
		{ // 469
			"({\n// A\na // B\n} = b)",
			"({a} = b);",
			"({\n\t// A\n\ta // B\n\t: a} = b);",
		},
		{ // 470
			"({\n// A\na // B\n= // C\nb // D\n} = c)",
			"({a = b} = c);",
			"({\n\t// A\n\ta // B\n\t: a = // C\n\tb // D\n} = c);",
		},
		{ // 471
			"({\n// A\na // B\n: // C\nb // D\n= // E\nc // F\n} = d)",
			"({a: b = c} = d);",
			"({\n\t// A\n\ta // B\n\t: // C\n\tb // D\n\t= // E\n\tc // F\n} = d);",
		},
		{ // 472
			"({ // A\n\n// B\na // C\n\n// D\n} = b)",
			"({a} = b);",
			"({ // A\n\n\t// B\n\ta // C\n\t: a\n// D\n} = b);",
		},
		{ // 473
			"({\n// A\n... // B\na // C\n} = b)",
			"({...a} = b);",
			"({\n\t// A\n\t... // B\n\ta // C\n} = b);",
		},
		{ // 474
			"(\n// A\n[a] // B\n= b)",
			"([a] = b);",
			"(\n\t// A\n\t[a] // B\n\t= b\n);",
		},
		{ // 475
			"(\n// A\n{a} // B\n= b)",
			"({a} = b);",
			"(\n\t// A\n\t{a: a} // B\n\t= b\n);",
		},
		{ // 476
			"(\n// A\na /* B */ => // C\n{} // D\n)",
			"(a => {});",
			"(\n\t// A\n\ta /* B */ => // C\n\t{} // D\n);",
		},
		{ // 477
			"(\n// A\n() /* B */ => /* C */ {} // D\n)",
			"(() => {});",
			"(\n\t// A\n\t() /* B */ => /* C */ {} // D\n);",
		},
		{ // 478
			"(\n// A\nasync /* B */ a /* C */ => // D\nb // E\n)",
			"(async a => b);",
			"(\n\t// A\n\tasync /* B */ a /* C */ => // D\n\tb // E\n);",
		},
		{ // 479
			"(// A\n\n// B\nasync /* C */ ()/* D */ => // E\n{}// F\n\n// G\n)",
			"(async () => {});",
			"( // A\n\n\t// B\n\tasync /* C */ () /* D */ => // E\n\t{} // F\n\n// G\n);",
		},
		{ // 480
			"function *a() {\n// A\nyield /* B */ a // C\n}",
			"function* a() {\n\tyield a;\n}",
			"function* a() {\n\t// A\n\tyield /* B */ a // C\n}",
		},
		{ // 481
			"function* a() {\n// A\nyield /* B */ * /* C */ a // D\n}",
			"function* a() {\n\tyield * a;\n}",
			"function* a() {\n\t// A\n\tyield /* B */ * /* C */ a // D\n}",
		},
		{ // 482
			"a // A\n?? b",
			"a ?? b;",
			"a // A\n?? b;",
		},
		{ // 483
			"a /* A */ ?? b",
			"a ?? b;",
			"a /* A */ ?? b;",
		},
		{ // 484
			"a /* A */ ? /* B */ b /* C */ : /* D */ c",
			"a ? b : c;",
			"a /* A */ ? /* B */ b /* C */ : /* D */ c;",
		},
		{ // 485
			"a // A\n? // B\nb // C\n: // D\nc",
			"a ? b : c;",
			"a // A\n? // B\nb // C\n: // D\nc;",
		},
		{ // 486
			"a /* A */ && /* B */ b",
			"a && b;",
			"a /* A */ && /* B */ b;",
		},
		{ // 487
			"a // A\n&& // B\nb",
			"a && b;",
			"a // A\n&& // B\nb;",
		},
		{ // 488
			"a /* A */ | /* B */ b",
			"a | b;",
			"a /* A */ | /* B */ b;",
		},
		{ // 489
			"a // A\n| // B\nb",
			"a | b;",
			"a // A\n| // B\nb;",
		},
		{ // 490
			"a /* A */ ^ /* B */ b",
			"a ^ b;",
			"a /* A */ ^ /* B */ b;",
		},
		{ // 491
			"a // A\n^ // B\nb",
			"a ^ b;",
			"a // A\n^ // B\nb;",
		},
		{ // 492
			"a /* A */ & /* B */ b",
			"a & b;",
			"a /* A */ & /* B */ b;",
		},
		{ // 493
			"a // A\n& // B\nb",
			"a & b;",
			"a // A\n& // B\nb;",
		},
		{ // 494
			"a /* A */ == /* B */ b",
			"a == b;",
			"a /* A */ == /* B */ b;",
		},
		{ // 495
			"a // A\n!== // B\nb",
			"a !== b;",
			"a // A\n!== // B\nb;",
		},
		{ // 496
			"a /* A */ < /* B */ b",
			"a < b;",
			"a /* A */ < /* B */ b;",
		},
		{ // 497
			"a // A\nin // B\nb",
			"a in b;",
			"a // A\nin // B\nb;",
		},
		{ // 498
			"a /* A */ << /* B */ b",
			"a << b;",
			"a /* A */ << /* B */ b;",
		},
		{ // 499
			"a // A\n>>> // B\nb",
			"a >>> b;",
			"a // A\n>>> // B\nb;",
		},
		{ // 500
			"a /* A */ + /* B */ b",
			"a + b;",
			"a /* A */ + /* B */ b;",
		},
		{ // 501
			"a // A\n- // B\nb",
			"a - b;",
			"a // A\n- // B\nb;",
		},
		{ // 502
			"a /* A */ * /* B */ b",
			"a * b;",
			"a /* A */ * /* B */ b;",
		},
		{ // 503
			"a // A\n% // B\nb",
			"a % b;",
			"a // A\n% // B\nb;",
		},
		{ // 504
			"a /* A */ ** /* B */ b",
			"a ** b;",
			"a /* A */ ** /* B */ b;",
		},
		{ // 505
			"a // A\n** // B\nb",
			"a ** b;",
			"a // A\n** // B\nb;",
		},
		{ // 506
			"do /* A */a /* B */; /* C */ while ( // D\n\n// E\nb // F\n\n// G\n) // H",
			"do a; while (b);",
			"do /* A */ a /* B */; /* C */ while ( // D\n\n\t// E\n\tb // F\n\n// G\n) // H\n",
		},
		{ // 507
			"do /* A */{}/* B */ while ( // C\n\n// D\nb // E\n\n// F\n) // G",
			"do {} while (b);",
			"do /* A */ {} /* B */ while ( // C\n\n\t// D\n\tb // E\n\n// F\n) // G\n",
		},
		{ // 508
			"([a // A\n= b] = c)",
			"([a = b] = c);",
			"([a // A\n\t= b] = c);",
		},
		{ // 509
			"// A\n\n// B\n#a /* C */ in /* D */ b // E",
			"#a in b;",
			"// A\n\n// B\n#a /* C */ in /* D */ b // E\n",
		},
		{ // 510
			"a({b: () => {\nc();\n}});",
			"a({b: () => {\n\tc();\n}});",
			"a({b: () => {\n\tc();\n}});",
		},
		{ // 511
			"a({b: // A\n() => {\nc();\n}});",
			"a({b: () => {\n\tc();\n}});",
			"a({\n\tb: // A\n\t() => {\n\t\tc();\n\t}\n});",
		},
		{ // 512
			"a?.[() => { // A\n}]",
			"a?.[() => {}];",
			"a?.[() => { // A\n}];",
		},
		{ // 513
			"a?.[ // A\n() => { // B\n}]",
			"a?.[() => {}];",
			"a?.[ // A\n\n\t() => { // B\n\t}\n];",
		},
		{ // 514
			"a?.[() => { // A\nreturn c;}]",
			"a?.[() => {\n\treturn c;\n}];",
			"a?.[() => { // A\n\n\treturn c;\n}];",
		},
		{ // 515
			"a?.[ // A\n() => { // B\nreturn c;}]",
			"a?.[() => {\n\treturn c;\n}];",
			"a?.[ // A\n\n\t() => { // B\n\n\t\treturn c;\n\t}\n];",
		},
		{ // 516
			"import(() => {})",
			"import(() => {});",
			"import(() => {});",
		},
		{ // 517
			"import(() => {return a})",
			"import(() => {\n\treturn a;\n});",
			"import(() => {\n\treturn a;\n});",
		},
		{ // 518
			"import(// A\n() => {return a})",
			"import(() => {\n\treturn a;\n});",
			"import( // A\n\n\t() => {\n\t\treturn a;\n\t});",
		},
		{ // 519
			"a()[() => {}]",
			"a()[() => {}];",
			"a()[() => {}];",
		},
		{ // 520
			"a()[() => {return b}]",
			"a()[() => {\n\treturn b;\n}];",
			"a()[() => {\n\treturn b;\n}];",
		},
		{ // 521
			"a()[// A\n() => {return b}]",
			"a()[() => {\n\treturn b;\n}];",
			"a()[ // A\n\n\t() => {\n\t\treturn b;\n\t}];",
		},
	} {
		for m, in := range [2]string{test.Input, test.VerboseOutput} {
			s, err := ParseScript(makeTokeniser(parser.NewStringTokeniser(in)))
			if err != nil {
				t.Errorf("test %d.%d.1: unexpected error: %s", n+1, m+1, err)
				continue
			}

			st.Verbose = false

			st.Buffer.Reset()
			s.Format(&st, 's')

			if str := st.String(); str != test.SimpleOutput {
				t.Errorf("test %d.%d.2: expecting %q, got %q\n%s", n+1, m+1, test.SimpleOutput, str, s)
			}

			st.Verbose = true

			st.Buffer.Reset()
			s.Format(&st, 's')

			if str := st.Buffer.String(); str != test.VerboseOutput {
				t.Errorf("test %d.%d.3: expecting %q, got %q\n%s", n+1, m+1, test.VerboseOutput, str, s)
			}
		}
	}
}

func TestPrintingModule(t *testing.T) {
	var st state

	for n, test := range [...]struct {
		Input, SimpleOutput, VerboseOutput string
	}{
		{ // 1
			"1",
			"1;",
			"1;",
		},
		{ // 2
			"1\n2\n3",
			"1;\n\n2;\n\n3;",
			"1;\n\n2;\n\n3;",
		},
		{ // 3
			"import\n'a'",
			"import 'a';",
			"import 'a';",
		},
		{ // 4
			"import\na\nfrom'b'",
			"import a from 'b';",
			"import a from 'b';",
		},
		{ // 5
			"export\n*\nfrom\n'a'",
			"export * from 'a';",
			"export * from 'a';",
		},
		{ // 6
			"export\n{\na\n}\nfrom\n'b'",
			"export {a} from 'b';",
			"export {a as a} from 'b';",
		},
		{ // 7
			"export\n{a}",
			"export {a};",
			"export {a as a};",
		},
		{ // 8
			"export\nvar\na",
			"export var a;",
			"export var a;",
		},
		{ // 9
			"export\nfunction\na(){}",
			"export function a() {}",
			"export function a() {}",
		},
		{ // 10
			"export\ndefault\nfunction(){}",
			"export default function () {}",
			"export default function () {}",
		},
		{ // 11
			"export\ndefault\nclass{}",
			"export default class {}",
			"export default class {}",
		},
		{ // 12
			"export\ndefault\na",
			"export default a;",
			"export default a;",
		},
		{ // 13
			"import\n{}\nfrom\n'z'",
			"import {} from 'z';",
			"import {} from 'z';",
		},
		{ // 14
			"import\na\nfrom'b'",
			"import a from 'b';",
			"import a from 'b';",
		},
		{ // 15
			"import*as\na\nfrom'b'",
			"import * as a from 'b';",
			"import * as a from 'b';",
		},
		{ // 16
			"import\na,*as\nb\nfrom'c'",
			"import a, * as b from 'c';",
			"import a, * as b from 'c';",
		},
		{ // 17
			"import\na,{b,c}from'd'",
			"import a, {b, c} from 'd';",
			"import a, {b as b, c as c} from 'd';",
		},
		{ // 18
			"export{}",
			"export {};",
			"export {};",
		},
		{ // 19
			"export{a}",
			"export {a};",
			"export {a as a};",
		},
		{ // 20
			"export{a,b}",
			"export {a, b};",
			"export {a as a, b as b};",
		},
		{ // 21
			"import{a\nas\nb}from'c'",
			"import {a as b} from 'c';",
			"import {a as b} from 'c';",
		},
		{ // 22
			"export{a\nas\nb}",
			"export {a as b};",
			"export {a as b};",
		},
		{ // 23
			"import . meta",
			"import.meta;",
			"import.meta;",
		},
		{ // 24
			"export\n*\nas\na\nfrom\n'b'",
			"export * as a from 'b';",
			"export * as a from 'b';",
		},
		{ // 25
			"import\na\nfrom'b' with {a:'b'}",
			"import a from 'b' with {a: 'b'};",
			"import a from 'b' with {a: 'b'};",
		},
		{ // 26
			"import\na\nfrom'b' with {'a':\"b\",}",
			"import a from 'b' with {'a': \"b\"};",
			"import a from 'b' with {'a': \"b\"};",
		},
		{ // 27
			"import\na\nfrom'b' with {'a':\"b\",c:'d'}",
			"import a from 'b' with {'a': \"b\", c: 'd'};",
			"import a from 'b' with {'a': \"b\", c: 'd'};",
		},
		{ // 28
			"const a = await b;",
			"const a = await b;",
			"const a = await b;",
		},
		{ // 29
			"export const a = await b;",
			"export const a = await b;",
			"export const a = await b;",
		},
		{ // 30
			"export var a = await b;",
			"export var a = await b;",
			"export var a = await b;",
		},
		{ // 31
			"export default await b;",
			"export default await b;",
			"export default await b;",
		},
		{ // 32
			"// A\n// B\n\nexport default await b;\n// C\n// D\n",
			"export default await b;",
			"// A\n// B\n\nexport default await b;\n// C\n// D\n",
		},
		{ // 33
			"/* A *//* B */\n// C\n\nexport default await b;\n// D\n/* E */   /* F */\n",
			"export default await b;",
			"/* A */ /* B */\n// C\n\nexport default await b;\n// D\n/* E */ /* F */",
		},
		{ // 34
			"/*\nA\n*//* B */\n// C\n\nexport default await b;\n// D\n/* E */   /*\n\nF\n\n*/\n",
			"export default await b;",
			"/*\nA\n*/ /* B */\n// C\n\nexport default await b;\n// D\n/* E */ /*\n\nF\n\n*/",
		},
		{ // 35
			"/*\nA\n*//* B */\n// C\n\n// D\nimport a from './b'; // E\n// F\n\n// G\n/* H */   /*\n\nI\n\n*/\n",
			"import a from './b';",
			"/*\nA\n*/ /* B */\n// C\n\n// D\nimport a from './b'; // E\n                     // F\n\n// G\n/* H */ /*\n\nI\n\n*/",
		},
		{ // 36
			"import {\n// A\na // B\n} from './b';",
			"import {a} from './b';",
			"import {\n\t// A\n\ta as a // B\n} from './b';",
		},
		{ // 37
			"import {\n// A\na /* B */ as // C\nb // D\n} from './b';",
			"import {a as b} from './b';",
			"import {\n\t// A\n\ta /* B */ as // C\n\tb // D\n} from './b';",
		},
		{ // 38
			"import { // A\n\n/* B */ a /* C */\n\n// D\n, // E\nb // F\n\n/* G */} from './c';",
			"import {a, b} from './c';",
			"import { // A\n\n\t/* B */ a as a /* C */\n\t// D\n\t, // E\n\tb as b // F\n\n/* G */} from './c';",
		},
		{ // 39
			"import // A\na /* B */ from './b';",
			"import a from './b';",
			"import // A\na /* B */ from './b';",
		},
		{ // 40
			"import // A\na, /* B */ * // C\nas /* D */ b /* E */ from './c';",
			"import a, * as b from './c';",
			"import // A\na, /* B */ * // C\nas /* D */ b /* E */ from './c';",
		},
		{ // 41
			"import // A\n{} /* B */ from './a';",
			"import {} from './a';",
			"import // A\n{} /* B */ from './a';",
		},
		{ // 42
			"import {} from './a' with {\n// A\nb/* B */:/* C */\"c\"// D\n\n// E\n, d:\"e\"};",
			"import {} from './a' with {b: \"c\", d: \"e\"};",
			"import {} from './a' with {\n\t// A\n\tb /* B */: /* C */ \"c\" // D\n\n\t// E\n\t, d: \"e\"};",
		},
		{ // 43
			"import {} from './a' with /* A */ {// B\n\n// C\nb/* D */:/* E */\"c\"// F\n\n// G\n, d:\"e\" // H\n\n// I\n};",
			"import {} from './a' with {b: \"c\", d: \"e\"};",
			"import {} from './a' with /* A */ { // B\n\n\t// C\n\tb /* D */: /* E */ \"c\" // F\n\n\t// G\n\t, d: \"e\" // H\n\n// I\n};",
		},
		{ // 44
			"import {} from './a' /* A */ with /* B */ {// C\n\n// D\nb/* E */:/* F */\"c\"// G\n\n// H\n, d:\"e\" // I\n\n// J\n};",
			"import {} from './a' with {b: \"c\", d: \"e\"};",
			"import {} from './a' /* A */ with /* B */ { // C\n\n\t// D\n\tb /* E */: /* F */ \"c\" // G\n\n\t// H\n\t, d: \"e\" // I\n\n// J\n};",
		},
		{ // 45
			"import a from /* A */ './b';",
			"import a from './b';",
			"import a from /* A */ './b';",
		},
		{ // 46
			"// A\n\n// B\nimport /* C */ a /* D */ from // E\n'b' /* F */ with /* G */ {c: 'd'} /* H */; // I\n",
			"import a from 'b' with {c: 'd'};",
			"// A\n\n// B\nimport /* C */ a /* D */ from // E\n'b' /* F */ with /* G */ {c: 'd'} /* H */; // I\n",
		},
		{ // 47
			"// A\n\n// B\nimport /* C */ \"\" /* D */; // E\n",
			"import \"\";",
			"// A\n\n// B\nimport /* C */ \"\" /* D */; // E\n",
		},
		{ // 48
			"export {\n// A\na /* B */ as /* C */ b // D\n\n// E\n, // F\nc // G\n\n};",
			"export {a as b, c};",
			"export {\n\t// A\n\ta /* B */ as /* C */ b // D\n\n\t// E\n\t, // F\n\tc as c // G\n};",
		},
		{ // 49
			"export { // A\n\n// B\na /* C */ as /* D */ b // E\n\n// F\n, // G\nc // H\n\n// I\n};",
			"export {a as b, c};",
			"export { // A\n\n\t// B\n\ta /* C */ as /* D */ b // E\n\n\t// F\n\t, // G\n\tc as c // H\n\n// I\n};",
		},
		{ // 50
			"// A\n\n// B\nexport/* C */{ // D\n\n// E\na // F\n\n// G\n, // H\nb // I\n\n// J\n}/* K */from/* L */''/* M */;",
			"export {a, b} from '';",
			"// A\n\n// B\nexport /* C */ { // D\n\n\t// E\n\ta as a // F\n\n\t// G\n\t, // H\n\tb as b // I\n\n// J\n} /* K */ from /* L */ '' /* M */;",
		},
		{ // 51
			"// A\n\n// B\nexport/* C */*/* D */as/* E */a/* F */from /* G */''/* H */",
			"export * as a from '';",
			"// A\n\n// B\nexport /* C */ * /* D */ as /* E */ a /* F */ from /* G */ '' /* H */;",
		},
		{ // 52
			"// A\n\n// B\nexport/* C */{ // D\n\n// E\na // F\n\n// G\n, // H\nb // I\n\n// J\n}/* K */;",
			"export {a, b};",
			"// A\n\n// B\nexport /* C */ { // D\n\n\t// E\n\ta as a // F\n\n\t// G\n\t, // H\n\tb as b // I\n\n// J\n} /* K */;",
		},
		{ // 53
			"// A\n\n// B\nexport/* C */var a;",
			"export var a;",
			"// A\n\n// B\nexport /* C */ var a;",
		},
		{ // 54
			"// A\n\n// B\nexport/* C */const a = 1;",
			"export const a = 1;",
			"// A\n\n// B\nexport /* C */ const a = 1;",
		},
		{ // 55
			"// A\n\n// B\nexport/* C */default/* D */function(){}",
			"export default function () {}",
			"// A\n\n// B\nexport /* C */ default /* D */ function () {}",
		},
		{ // 56
			"// A\n\n// B\nexport/* C */default/* D */class{}",
			"export default class {}",
			"// A\n\n// B\nexport /* C */ default /* D */ class {}",
		},
		{ // 57
			"// A\n\n// B\nexport/* C */default/* D */1/* E */;",
			"export default 1;",
			"// A\n\n// B\nexport /* C */ default /* D */ 1 /* E */;",
		},
		{ // 58
			"// A\n\n// B\nexport/* C */default/* D */1/* E */; // F\n\n// G",
			"export default 1;",
			"// A\n\n// B\nexport /* C */ default /* D */ 1 /* E */; // F\n\n// G\n",
		},
		{ // 59
			"for // A\nawait // B\n( // C\n\n// D\nconst // E\n[]// F\nof // G\nb // H\n\n// I\n) // J\n;",
			"for await (const [] of b) ;",
			"for // A\nawait // B\n( // C\n\n\t// D\n\tconst // E\n\t[] // F\n\tof // G\n\tb // H\n\n// I\n) // J\n;",
		},
		{ // 60
			"export * // A\nas a from 'b'",
			"export * as a from 'b';",
			"export * // A\nas a from 'b';",
		},
		{ // 61
			"export * // A\nas a // A\nfrom 'b'",
			"export * as a from 'b';",
			"export * // A\nas a // A\nfrom 'b';",
		},
		{ // 62
			"export * // A\nas a /* A */ from 'b'",
			"export * as a from 'b';",
			"export * // A\nas a /* A */ from 'b';",
		},
		{ // 63
			"import * as a // A\nfrom 'b'",
			"import * as a from 'b';",
			"import * as a // A\nfrom 'b';",
		},
		{ // 64
			"import * as a /* A */ from 'b'",
			"import * as a from 'b';",
			"import * as a /* A */ from 'b';",
		},
		{ // 65
			"export * // A\nas a from 'b'",
			"export * as a from 'b';",
			"export * // A\nas a from 'b';",
		},
		{ // 66
			"export {a // A\nas // B\nb}",
			"export {a as b};",
			"export {a // A\n\tas // B\n\tb};",
		},
		{ // 67
			"import {a // A\nas b} from 'c'",
			"import {a as b} from 'c';",
			"import {a // A\n\tas b} from 'c';",
		},
		{ // 68
			"import {a /* A */ as b} from 'c'",
			"import {a as b} from 'c';",
			"import {a /* A */ as b} from 'c';",
		},
	} {
		for m, in := range [2]string{test.Input, test.VerboseOutput} {
			s, err := ParseModule(makeTokeniser(parser.NewStringTokeniser(in)))
			if err != nil {
				t.Errorf("test %d.%d.1: unexpected error: %s", n+1, m+1, err)

				continue
			}

			st.Verbose = false

			st.Buffer.Reset()
			s.Format(&st, 's')

			if str := st.Buffer.String(); str != test.SimpleOutput {
				t.Errorf("test %d.%d.2: expecting %q, got %q\n%s", n+1, m+1, test.SimpleOutput, str, s)
			}

			st.Verbose = true

			st.Buffer.Reset()
			s.Format(&st, 's')

			if str := st.Buffer.String(); str != test.VerboseOutput {
				t.Errorf("test %d.%d.3: expecting %q, got %q\n%s", n+1, m+1, test.VerboseOutput, str, s)
			}
		}
	}
}
