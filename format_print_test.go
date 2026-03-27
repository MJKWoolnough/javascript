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
			"for (a;;) b;",
		},
		{ // 38
			"for(var a=b;c<d;e++){}",
			"for (var a = b; c < d; e++) {}",
			"for (var a = b; c < d; e++) {}",
		},
		{ // 39
			"for(\nvar a=b;\nc<d;\ne++){}",
			"for (var a = b; c < d; e++) {}",
			"for (var a = b; c < d; e++) {}",
		},
		{ // 40
			"for(let a=b;c<d;e++){}",
			"for (let a = b; c < d; e++) {}",
			"for (let a = b; c < d; e++) {}",
		},
		{ // 41
			"for(\nlet a=b;\nc<d;\ne++){}",
			"for (let a = b; c < d; e++) {}",
			"for (let a = b; c < d; e++) {}",
		},
		{ // 42
			"for(const a=b;c<d;e++){}",
			"for (const a = b; c < d; e++) {}",
			"for (const a = b; c < d; e++) {}",
		},
		{ // 43
			"for(\nconst a=b;\nc<d;\ne++){}",
			"for (const a = b; c < d; e++) {}",
			"for (const a = b; c < d; e++) {}",
		},
		{ // 44
			"for(a in b){}",
			"for (a in b) {}",
			"for (a in b) {}",
		},
		{ // 45
			"for\n(a\nin\nb\n)\n{}",
			"for (a in b) {}",
			"for (a in b) {}",
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
			"for (a of b) {}",
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
			"async () => {\n\tfor await (a of b) {}\n};",
		},
		{ // 62
			"switch(a) {}",
			"switch (a) {}",
			"switch (a) {}",
		},
		{ // 63
			"switch\n(\na\n)\n{\n}",
			"switch (a) {}",
			"switch (a) {}",
		},
		{ // 64
			"switch(a){case b:case c:default:case d:case e:}",
			"switch (a) {\ncase b:\ncase c:\ndefault:\ncase d:\ncase e:\n}",
			"switch (a) {\ncase b:\ncase c:\ndefault:\ncase d:\ncase e:\n}",
		},
		{ // 65
			"switch\n\n(\n\na\n\n)\n\n{\n\ncase\n\nb\n\n:\n\ncase\n\nc\n\n:\n\ndefault\n\n:\n\ncase\n\nd\n\n:\n\ncase\n\ne\n\n:\n\n}",
			"switch (a) {\ncase b:\ncase c:\ndefault:\ncase d:\ncase e:\n}",
			"switch (a) {\ncase b:\ncase c:\ndefault:\ncase d:\ncase e:\n}",
		},
		{ // 66
			"switch( // A\na) {}",
			"switch (a) {}",
			"switch ( // A\n\n\ta\n) {}",
		},
		{ // 67
			"with(a)b",
			"with (a) b;",
			"with (a) b;",
		},
		{ // 68
			"with\n(\na\n)\nb",
			"with (a) b;",
			"with (a) b;",
		},
		{ // 69
			"with( // A\na\n // B\n)b",
			"with (a) b;",
			"with ( // A\n\n\ta\n// B\n) b;",
		},
		{ // 70
			"function a(){}",
			"function a() {}",
			"function a() {}",
		},
		{ // 71
			"function a(b){}",
			"function a(b) {}",
			"function a(b) {}",
		},
		{ // 72
			"function a(b,c){}",
			"function a(b, c) {}",
			"function a(b, c) {}",
		},
		{ // 73
			"function\na(\nb\n,\nc\n){}",
			"function a(b, c) {}",
			"function a(b, c) {}",
		},
		{ // 74
			"function*a(){}",
			"function* a() {}",
			"function* a() {}",
		},
		{ // 75
			"function* a(b){}",
			"function* a(b) {}",
			"function* a(b) {}",
		},
		{ // 76
			"function *a(b,c){}",
			"function* a(b, c) {}",
			"function* a(b, c) {}",
		},
		{ // 77
			"function\n*a(\nb\n,\nc\n){}",
			"function* a(b, c) {}",
			"function* a(b, c) {}",
		},
		{ // 78
			"async function a(){}",
			"async function a() {}",
			"async function a() {}",
		},
		{ // 79
			"async function a(b){}",
			"async function a(b) {}",
			"async function a(b) {}",
		},
		{ // 80
			"async function a(b,c){}",
			"async function a(b, c) {}",
			"async function a(b, c) {}",
		},
		{ // 81
			"async function\na(\nb\n,\nc\n){}",
			"async function a(b, c) {}",
			"async function a(b, c) {}",
		},
		{ // 82
			"async function*a(){}",
			"async function* a() {}",
			"async function* a() {}",
		},
		{ // 83
			"async function* a(b){}",
			"async function* a(b) {}",
			"async function* a(b) {}",
		},
		{ // 84
			"async function *a(b,c){}",
			"async function* a(b, c) {}",
			"async function* a(b, c) {}",
		},
		{ // 85
			"async function\n*a(\nb\n,\nc\n){}",
			"async function* a(b, c) {}",
			"async function* a(b, c) {}",
		},
		{ // 86
			"a = function(){}",
			"a = function () {};",
			"a = function () {};",
		},
		{ // 87
			"a=function(b){}",
			"a = function (b) {};",
			"a = function (b) {};",
		},
		{ // 88
			"a=function *(b,c){}",
			"a = function* (b, c) {};",
			"a = function* (b, c) {};",
		},
		{ // 89
			"a=function\n(\nb\n,\nc\n){}",
			"a = function (b, c) {};",
			"a = function (b, c) {};",
		},
		{ // 90
			"try{}catch{}",
			"try {} catch {}",
			"try {} catch {}",
		},
		{ // 91
			"try\n{\n}\ncatch\n{\n}",
			"try {} catch {}",
			"try {} catch {}",
		},
		{ // 92
			"try{}catch(a){}",
			"try {} catch (a) {}",
			"try {} catch (a) {}",
		},
		{ // 93
			"try\n{\n}\ncatch\n(\na\n)\n{\n}",
			"try {} catch (a) {}",
			"try {} catch (a) {}",
		},
		{ // 94
			"try{}catch({}){}",
			"try {} catch ({}) {}",
			"try {} catch ({}) {}",
		},
		{ // 95
			"try{}catch([]){}",
			"try {} catch ([]) {}",
			"try {} catch ([]) {}",
		},
		{ // 96
			"try{}finally{}",
			"try {} finally {}",
			"try {} finally {}",
		},
		{ // 97
			"try\n{\n}\nfinally\n{\n}",
			"try {} finally {}",
			"try {} finally {}",
		},
		{ // 98
			"try{}catch{}finally{}",
			"try {} catch {} finally {}",
			"try {} catch {} finally {}",
		},
		{ // 99
			"try\n{\n}\ncatch\n{\n}\nfinally\n{\n}",
			"try {} catch {} finally {}",
			"try {} catch {} finally {}",
		},
		{ // 100
			"try{}catch(a){}finally{}",
			"try {} catch (a) {} finally {}",
			"try {} catch (a) {} finally {}",
		},
		{ // 101
			"try\n{\n}\ncatch\n(\na\n)\n{\n}\nfinally\n{\n}",
			"try {} catch (a) {} finally {}",
			"try {} catch (a) {} finally {}",
		},
		{ // 102
			"try{}catch( // A\na\n // B\n){}",
			"try {} catch (a) {}",
			"try {} catch ( // A\n\n\ta\n// B\n) {}",
		},
		{ // 103
			"class a{}",
			"class a {}",
			"class a {}",
		},
		{ // 104
			"class\na\n{\n}\n",
			"class a {}",
			"class a {}",
		},
		{ // 105
			"class a extends b {}",
			"class a extends b {}",
			"class a extends b {}",
		},
		{ // 106
			"class\na\nextends\nb\n{\n}",
			"class a extends b {}",
			"class a extends b {}",
		},
		{ // 107
			"a = class{}",
			"a = class {};",
			"a = class {};",
		},
		{ // 108
			"a\n=\nclass\nb\n{\n}",
			"a = class b {};",
			"a = class b {};",
		},
		{ // 109
			"a\n=\nclass\nextends\nb\n{\n}",
			"a = class extends b {};",
			"a = class extends b {};",
		},
		{ // 110
			"let a = 1",
			"let a = 1;",
			"let a = 1;",
		},
		{ // 111
			"let\na\n=\n1\n",
			"let a = 1;",
			"let a = 1;",
		},
		{ // 112
			"let a=1,b=2,c=3",
			"let a = 1, b = 2, c = 3;",
			"let a = 1,\nb = 2,\nc = 3;",
		},
		{ // 113
			"const a = 1",
			"const a = 1;",
			"const a = 1;",
		},
		{ // 114
			"const\na\n=\n1\n",
			"const a = 1;",
			"const a = 1;",
		},
		{ // 115
			"const a=1,b=2,c=3",
			"const a = 1, b = 2, c = 3;",
			"const a = 1,\nb = 2,\nc = 3;",
		},
		{ // 116
			"let a",
			"let a;",
			"let a;",
		},
		{ // 117
			"let\na\n,\nb\n=\n1\n,\nc\n",
			"let a, b = 1, c;",
			"let a,\nb = 1,\nc;",
		},
		{ // 118
			"const a",
			"const a;",
			"const a;",
		},
		{ // 119
			"const\na\n,\nb\n=\n1\n,\nc\n",
			"const a, b = 1, c;",
			"const a,\nb = 1,\nc;",
		},
		{ // 120
			"let [a]=1",
			"let [a] = 1;",
			"let [a] = 1;",
		},
		{ // 121
			"const\n[\na\n]\n=\n1",
			"const [a] = 1;",
			"const [a] = 1;",
		},
		{ // 122
			"let {a}=1",
			"let {a} = 1;",
			"let {a: a} = 1;",
		},
		{ // 123
			"const\n{\na\n}\n=\n1",
			"const {a} = 1;",
			"const {a: a} = 1;",
		},
		{ // 124
			"function* a() {yield a}",
			"function* a() {\n\tyield a;\n}",
			"function* a() {\n\tyield a;\n}",
		},
		{ // 125
			"() => {}",
			"() => {};",
			"() => {};",
		},
		{ // 126
			"a=b",
			"a = b;",
			"a = b;",
		},
		{ // 127
			"a/=b",
			"a /= b;",
			"a /= b;",
		},
		{ // 128
			"a%=b",
			"a %= b;",
			"a %= b;",
		},
		{ // 129
			"a+=b",
			"a += b;",
			"a += b;",
		},
		{ // 130
			"a-=b",
			"a -= b;",
			"a -= b;",
		},
		{ // 131
			"a<<=b",
			"a <<= b;",
			"a <<= b;",
		},
		{ // 132
			"a>>=b",
			"a >>= b;",
			"a >>= b;",
		},
		{ // 133
			"a>>>=b",
			"a >>>= b;",
			"a >>>= b;",
		},
		{ // 134
			"a&=b",
			"a &= b;",
			"a &= b;",
		},
		{ // 135
			"a^=b",
			"a ^= b;",
			"a ^= b;",
		},
		{ // 136
			"a|=b",
			"a |= b;",
			"a |= b;",
		},
		{ // 137
			"a**=b",
			"a **= b;",
			"a **= b;",
		},
		{ // 138
			"a?b:c",
			"a ? b : c;",
			"a ? b : c;",
		},
		{ // 139
			"new a",
			"new a;",
			"new a;",
		},
		{ // 140
			"a()",
			"a();",
			"a();",
		},
		{ // 141
			"var {a} = 1",
			"var {a} = 1;",
			"var {a: a} = 1;",
		},
		{ // 142
			"var { a , b, ...c } = 1",
			"var {a, b, ...c} = 1;",
			"var {a: a, b: b, ...c} = 1;",
		},
		{ // 143
			"var [a] = 1",
			"var [a] = 1;",
			"var [a] = 1;",
		},
		{ // 144
			"var [ a , b, ...c ] = 1",
			"var [a, b, ...c] = 1;",
			"var [a, b, ...c] = 1;",
		},
		{ // 145
			"switch (a) {case 1:}",
			"switch (a) {\ncase 1:\n}",
			"switch (a) {\ncase 1:\n}",
		},
		{ // 146
			"switch (a) {case 1:b;c;d}",
			"switch (a) {\ncase 1:\n\tb;\n\tc;\n\td;\n}",
			"switch (a) {\ncase 1:\n\tb;\n\tc;\n\td;\n}",
		},
		{ // 147
			"function a(b){}",
			"function a(b) {}",
			"function a(b) {}",
		},
		{ // 148
			"function\na(\nb\n)\n{\n}\n",
			"function a(b) {}",
			"function a(b) {}",
		},
		{ // 149
			"function a(b,c,...d){}",
			"function a(b, c, ...d) {}",
			"function a(b, c, ...d) {}",
		},
		{ // 150
			"class\na{b(){}c\n(){}}",
			"class a {\n\tb() {}\n\tc() {}\n}",
			"class a {\n\tb() {}\n\tc() {}\n}",
		},
		{ // 151
			"class\na{*b(){}\n*\nc\n(){}}",
			"class a {\n\t* b() {}\n\t* c() {}\n}",
			"class a {\n\t* b() {}\n\t* c() {}\n}",
		},
		{ // 152
			"class\na{async b(){}\nasync c\n(){}}",
			"class a {\n\tasync b() {}\n\tasync c() {}\n}",
			"class a {\n\tasync b() {}\n\tasync c() {}\n}",
		},
		{ // 153
			"class\na{async *b(){}\nasync *\nc\n(){}}",
			"class a {\n\tasync * b() {}\n\tasync * c() {}\n}",
			"class a {\n\tasync * b() {}\n\tasync * c() {}\n}",
		},
		{ // 154
			"class\na{get\nb(){}\nget c\n(){}}",
			"class a {\n\tget b() {}\n\tget c() {}\n}",
			"class a {\n\tget b() {}\n\tget c() {}\n}",
		},
		{ // 155
			"class\na{set\nb(c){}\nset d\n(e){}}",
			"class a {\n\tset b(c) {}\n\tset d(e) {}\n}",
			"class a {\n\tset b(c) {}\n\tset d(e) {}\n}",
		},
		{ // 156
			"class\na{static\nb(){}\nstatic c\n(){}}",
			"class a {\n\tstatic b() {}\n\tstatic c() {}\n}",
			"class a {\n\tstatic b() {}\n\tstatic c() {}\n}",
		},
		{ // 157
			"class\na{static\n*b(){}\nstatic *\nc\n(){}}",
			"class a {\n\tstatic * b() {}\n\tstatic * c() {}\n}",
			"class a {\n\tstatic * b() {}\n\tstatic * c() {}\n}",
		},
		{ // 158
			"class\na{static\nasync b(){}\nstatic async c\n(){}}",
			"class a {\n\tstatic async b() {}\n\tstatic async c() {}\n}",
			"class a {\n\tstatic async b() {}\n\tstatic async c() {}\n}",
		},
		{ // 159
			"class\na{static\nasync *b(){}\nstatic async *\nc(){}}",
			"class a {\n\tstatic async * b() {}\n\tstatic async * c() {}\n}",
			"class a {\n\tstatic async * b() {}\n\tstatic async * c() {}\n}",
		},
		{ // 160
			"class\na{static\nget\nb(){}static get c\n(){}}",
			"class a {\n\tstatic get b() {}\n\tstatic get c() {}\n}",
			"class a {\n\tstatic get b() {}\n\tstatic get c() {}\n}",
		},
		{ // 161
			"class\na{static\nset\nb(c){}static set d\n(e){}}",
			"class a {\n\tstatic set b(c) {}\n\tstatic set d(e) {}\n}",
			"class a {\n\tstatic set b(c) {}\n\tstatic set d(e) {}\n}",
		},
		{ // 162
			"a",
			"a;",
			"a;",
		},
		{ // 163
			"a?b:c",
			"a ? b : c;",
			"a ? b : c;",
		},
		{ // 164
			"a\n?\nb\n:\nc\n",
			"a ? b : c;",
			"a ? b : c;",
		},
		{ // 165
			"a=>b",
			"a => b;",
			"a => b;",
		},
		{ // 166
			"a =>\nb",
			"a => b;",
			"a => b;",
		},
		{ // 167
			"async a => b",
			"async a => b;",
			"async a => b;",
		},
		{ // 168
			"(a,b)=>b",
			"(a, b) => b;",
			"(a, b) => b;",
		},
		{ // 169
			"async (a,b)=>c",
			"async (a, b) => c;",
			"async (a, b) => c;",
		},
		{ // 170
			"a=>{}",
			"a => {};",
			"a => {};",
		},
		{ // 171
			"async a=>{}",
			"async a => {};",
			"async a => {};",
		},
		{ // 172
			"(a,b)=>{}",
			"(a, b) => {};",
			"(a, b) => {};",
		},
		{ // 173
			"async(a,b)=>{}",
			"async (a, b) => {};",
			"async (a, b) => {};",
		},
		{ // 174
			"new a",
			"new a;",
			"new a;",
		},
		{ // 175
			"new\nnew \n new	new\na",
			"new new new new a;",
			"new new new new a;",
		},
		{ // 176
			"super()",
			"super();",
			"super();",
		},
		{ // 177
			"super\n()",
			"super();",
			"super();",
		},
		{ // 178
			"import(a)",
			"import(a);",
			"import(a);",
		},
		{ // 179
			"import\n(a)",
			"import(a);",
			"import(a);",
		},
		{ // 180
			"a()",
			"a();",
			"a();",
		},
		{ // 181
			"a\n()",
			"a();",
			"a();",
		},
		{ // 182
			"a\n()\n()",
			"a()();",
			"a()();",
		},
		{ // 183
			"a()[b]",
			"a()[b];",
			"a()[b];",
		},
		{ // 184
			"a().b",
			"a().b;",
			"a().b;",
		},
		{ // 185
			"a()`b`",
			"a()`b`;",
			"a()`b`;",
		},
		{ // 186
			"var{a}=b",
			"var {a} = b;",
			"var {a: a} = b;",
		},
		{ // 187
			"var\n{\na\n:\nb\n}\n=\nc\n",
			"var {a: b} = c;",
			"var {a: b} = c;",
		},
		{ // 188
			"var[a=b]=c",
			"var [a = b] = c;",
			"var [a = b] = c;",
		},
		{ // 189
			"var[a=[b]] = c",
			"var [a = [b]] = c;",
			"var [a = [b]] = c;",
		},
		{ // 190
			"var[a={b}] = c",
			"var [a = {b}] = c;",
			"var [a = {b: b}] = c;",
		},
		{ // 191
			"var a={[\"b\"]:c}",
			"var a = {[\"b\"]: c};",
			"var a = {[\"b\"]: c};",
		},
		{ // 192
			"a||b",
			"a || b;",
			"a || b;",
		},
		{ // 193
			"(a,b,c)",
			"(a, b, c);",
			"(a, b, c);",
		},
		{ // 194
			"(\na\n,\nb\n,\nc\n)\n",
			"(a, b, c);",
			"(a, b, c);",
		},
		{ // 195
			"var a=(b,c,...d)=>{}",
			"var a = (b, c, ...d) => {};",
			"var a = (b, c, ...d) => {};",
		},
		{ // 196
			"var a=(b,c,...[e])=>{}",
			"var a = (b, c, ...[e]) => {};",
			"var a = (b, c, ...[e]) => {};",
		},
		{ // 197
			"var a=(b,c,...{...e})=>{}",
			"var a = (b, c, ...{...e}) => {};",
			"var a = (b, c, ...{...e}) => {};",
		},
		{ // 198
			"new a()",
			"new a();",
			"new a();",
		},
		{ // 199
			"new\nnew\na\n(\n)\n(\n)\n",
			"new new a()();",
			"new new a()();",
		},
		{ // 200
			"a\n[\n1\n]\n",
			"a[1];",
			"a[1];",
		},
		{ // 201
			"a\n.\nb\n",
			"a.b;",
			"a\n.b;",
		},
		{ // 202
			"a\n`b`",
			"a`b`;",
			"a`b`;",
		},
		{ // 203
			"new\nsuper\n[\na\n]\n[\nb\n]\n.\nc`d`\n(\nnew\n.\ntarget\n)\n",
			"new super[a][b].c`d`(new.target);",
			"new super[a][b]\n.c`d`(new.target);",
		},
		{ // 204
			"a(b,c,...d)",
			"a(b, c, ...d);",
			"a(b, c, ...d);",
		},
		{ // 205
			"a\n(\n...\nb\n)\n",
			"a(...b);",
			"a(...b);",
		},
		{ // 206
			"`a`",
			"`a`;",
			"`a`;",
		},
		{ // 207
			"`a${b}c`",
			"`a${b}c`;",
			"`a${b}c`;",
		},
		{ // 208
			"`a${\nb\n}c${\nd\n}e`",
			"`a${b}c${d}e`;",
			"`a${b}c${d}e`;",
		},
		{ // 209
			"{\n`a`\n}",
			"{\n\t`a`;\n}",
			"{\n\t`a`;\n}",
		},
		{ // 210
			"{\n`a\nb`\n}",
			"{\n\t`a\nb`;\n}",
			"{\n\t`a\nb`;\n}",
		},
		{ // 211
			"{\n`a\nb${c}d\ne${f}g\nh`\n}",
			"{\n\t`a\nb${c}d\ne${f}g\nh`;\n}",
			"{\n\t`a\nb${c}d\ne${f}g\nh`;\n}",
		},
		{ // 212
			"a&&b",
			"a && b;",
			"a && b;",
		},
		{ // 213
			"this",
			"this;",
			"this;",
		},
		{ // 214
			"a",
			"a;",
			"a;",
		},
		{ // 215
			"1",
			"1;",
			"1;",
		},
		{ // 216
			"[\n]\n",
			"[];",
			"[];",
		},
		{ // 217
			"var a={}",
			"var a = {};",
			"var a = {};",
		},
		{ // 218
			"var a=function(){}",
			"var a = function () {};",
			"var a = function () {};",
		},
		{ // 219
			"var a=class{}",
			"var a = class {};",
			"var a = class {};",
		},
		{ // 220
			"`a`",
			"`a`;",
			"`a`;",
		},
		{ // 221
			"(a)",
			"(a);",
			"(a);",
		},
		{ // 222
			"a|b",
			"a | b;",
			"a | b;",
		},
		{ // 223
			"[a,b,...c]",
			"[a, b, ...c];",
			"[a, b, ...c];",
		},
		{ // 224
			"[...a]",
			"[...a];",
			"[...a];",
		},
		{ // 225
			"[a]",
			"[a];",
			"[a];",
		},
		{ // 226
			"var a={b:c}",
			"var a = {b: c};",
			"var a = {b: c};",
		},
		{ // 227
			"var a={b:c,d:e}",
			"var a = {b: c, d: e};",
			"var a = {b: c, d: e};",
		},
		{ // 228
			"var a={\nb\n:\nc\n}",
			"var a = {b: c};",
			"var a = {b: c};",
		},
		{ // 229
			"var a={\nb\n:\nc\n,\nd\n:\ne\n}",
			"var a = {b: c, d: e};",
			"var a = {b: c, d: e};",
		},
		{ // 230
			"a^b",
			"a ^ b;",
			"a ^ b;",
		},
		{ // 231
			"var a\n=\n{\nb\n:\nc\n,\nd\n,\ne\n=\nf\n,\ng\n(\n)\n{\n}\n,\n...\nh\n}\n",
			"var a = {b: c, d, e = f, g() {}, ...h};",
			"var a = {b: c, d: d, e = f, g() {}, ...h};",
		},
		{ // 232
			"a&b",
			"a & b;",
			"a & b;",
		},
		{ // 233
			"a==b",
			"a == b;",
			"a == b;",
		},
		{ // 234
			"a!=b",
			"a != b;",
			"a != b;",
		},
		{ // 235
			"a===b",
			"a === b;",
			"a === b;",
		},
		{ // 236
			"a!==b",
			"a !== b;",
			"a !== b;",
		},
		{ // 237
			"a<b",
			"a < b;",
			"a < b;",
		},
		{ // 238
			"a>b",
			"a > b;",
			"a > b;",
		},
		{ // 239
			"a<=b",
			"a <= b;",
			"a <= b;",
		},
		{ // 240
			"a>=b",
			"a >= b;",
			"a >= b;",
		},
		{ // 241
			"a instanceof b",
			"a instanceof b;",
			"a instanceof b;",
		},
		{ // 242
			"a in b",
			"a in b;",
			"a in b;",
		},
		{ // 243
			"a<<b",
			"a << b;",
			"a << b;",
		},
		{ // 244
			"a>>b",
			"a >> b;",
			"a >> b;",
		},
		{ // 245
			"a>>>b",
			"a >>> b;",
			"a >>> b;",
		},
		{ // 246
			"a+b",
			"a + b;",
			"a + b;",
		},
		{ // 247
			"a-b",
			"a - b;",
			"a - b;",
		},
		{ // 248
			"a*b",
			"a * b;",
			"a * b;",
		},
		{ // 249
			"a/b",
			"a / b;",
			"a / b;",
		},
		{ // 250
			"a%b",
			"a % b;",
			"a % b;",
		},
		{ // 251
			"a**b",
			"a ** b;",
			"a ** b;",
		},
		{ // 252
			"delete a",
			"delete a;",
			"delete a;",
		},
		{ // 253
			"void a",
			"void a;",
			"void a;",
		},
		{ // 254
			"typeof a",
			"typeof a;",
			"typeof a;",
		},
		{ // 255
			"+\na",
			"+a;",
			"+a;",
		},
		{ // 256
			"-\na",
			"-a;",
			"-a;",
		},
		{ // 257
			"~\na",
			"~a;",
			"~a;",
		},
		{ // 258
			"!\na",
			"!a;",
			"!a;",
		},
		{ // 259
			"async function a(){await b}",
			"async function a() {\n\tawait b;\n}",
			"async function a() {\n\tawait b;\n}",
		},
		{ // 260
			"a ++",
			"a++;",
			"a++;",
		},
		{ // 261
			"a --",
			"a--;",
			"a--;",
		},
		{ // 262
			"++\na",
			"++a;",
			"++a;",
		},
		{ // 263
			"--\na",
			"--a;",
			"--a;",
		},
		{ // 264
			"a: function b(){}",
			"a: function b() {}",
			"a: function b() {}",
		},
		{ // 265
			"a: b",
			"a: b;",
			"a: b;",
		},
		{ // 266
			"continue a",
			"continue a;",
			"continue a;",
		},
		{ // 267
			"debugger",
			"debugger;",
			"debugger;",
		},
		{ // 268
			"for(var a,b,\nc;;){}",
			"for (var a, b, c;;) {}",
			"for (var a, b, c;;) {}",
		},
		{ // 269
			"for(var{a}in b){}",
			"for (var {a} in b) {}",
			"for (var {a: a} in b) {}",
		},
		{ // 270
			"for(var[a]in b){}",
			"for (var [a] in b) {}",
			"for (var [a] in b) {}",
		},
		{ // 271
			"switch(a){default:b}",
			"switch (a) {\ndefault:\n\tb;\n}",
			"switch (a) {\ndefault:\n\tb;\n}",
		},
		{ // 272
			"function*a(){yield *b}",
			"function* a() {\n\tyield * b;\n}",
			"function* a() {\n\tyield * b;\n}",
		},
		{ // 273
			"a*=b",
			"a *= b;",
			"a *= b;",
		},
		{ // 274
			"var[[a]]=b",
			"var [[a]] = b;",
			"var [[a]] = b;",
		},
		{ // 275
			"var[{a}]=b",
			"var [{a}] = b;",
			"var [{a: a}] = b;",
		},
		{ // 276
			"super\n.\na\n",
			"super.a;",
			"super.a;",
		},
		{ // 277
			"a\n?.\nb",
			"a?.b;",
			"a?.b;",
		},
		{ // 278
			"a\n??\nb",
			"a ?? b;",
			"a ?? b;",
		},
		{ // 279
			"a\n??\nb\n??\nc",
			"a ?? b ?? c;",
			"a ?? b ?? c;",
		},
		{ // 280
			"a = ([b]) => b",
			"a = ([b]) => b;",
			"a = ([b]) => b;",
		},
		{ // 281
			"a?.b().c",
			"a?.b().c;",
			"a?.b().c;",
		},
		{ // 282
			"a?.b()?.c",
			"a?.b()?.c;",
			"a?.b()?.c;",
		},
		{ // 283
			"a&&=1",
			"a &&= 1;",
			"a &&= 1;",
		},
		{ // 284
			"a||=1",
			"a ||= 1;",
			"a ||= 1;",
		},
		{ // 285
			"a??=1",
			"a ??= 1;",
			"a ??= 1;",
		},
		{ // 286
			"[a, b] = [b, a]",
			"[a, b] = [b, a];",
			"[a, b] = [b, a];",
		},
		{ // 287
			"[a.b, a.c] = [a.c, a.b]",
			"[a.b, a.c] = [a.c, a.b];",
			"[a.b, a.c] = [a.c, a.b];",
		},
		{ // 288
			"{a}",
			"{\n\ta;\n}",
			"{\n\ta;\n}",
		},
		{ // 289
			"{a;b}",
			"{\n\ta;\n\tb;\n}",
			"{\n\ta;\n\tb;\n}",
		},
		{ // 290
			"{a;\nb}",
			"{\n\ta;\n\tb;\n}",
			"{\n\ta;\n\tb;\n}",
		},
		{ // 291
			"{\na;\nb\n}",
			"{\n\ta;\n\tb;\n}",
			"{\n\ta;\n\tb;\n}",
		},
		{ // 292
			"({a, b} = {a: 1, b: 2})",
			"({a, b} = {a: 1, b: 2});",
			"({a: a, b: b} = {a: 1, b: 2});",
		},
		{ // 293
			"[a,b,...c] = [b, a]",
			"[a, b, ...c] = [b, a];",
			"[a, b, ...c] = [b, a];",
		},
		{ // 294
			"({a,b,...c}=d)",
			"({a, b, ...c} = d);",
			"({a: a, b: b, ...c} = d);",
		},
		{ // 295
			"({a:{b,c: d,...e},...f}=g)",
			"({a: {b, c: d, ...e}, ...f} = g);",
			"({a: {b: b, c: d, ...e}, ...f} = g);",
		},
		{ // 296
			"[a, ,[b,{c},,...d],,...e]=f",
			"[a, , [b, {c}, , ...d], , ...e] = f;",
			"[a, , [b, {c: c}, , ...d], , ...e] = f;",
		},
		{ // 297
			"a() ?.\nb",
			"a()?.b;",
			"a()?.b;",
		},
		{ // 298
			"a ?. [1]",
			"a?.[1];",
			"a?.[1];",
		},
		{ // 299
			"a ?. `1`",
			"a?.`1`;",
			"a?.`1`;",
		},
		{ // 300
			"a()\n.then()\n.catch()",
			"a().then().catch();",
			"a()\n.then()\n.catch();",
		},
		{ // 301
			"[a=b]=c",
			"[a = b] = c;",
			"[a = b] = c;",
		},
		{ // 302
			"({a=b}=c)",
			"({a = b} = c);",
			"({a: a = b} = c);",
		},
		{ // 303
			"a.#b",
			"a.#b;",
			"a.#b;",
		},
		{ // 304
			"a\n.#b",
			"a.#b;",
			"a\n.#b;",
		},
		{ // 305
			"a.#b.c",
			"a.#b.c;",
			"a.#b.c;",
		},
		{ // 306
			"a\n.#b\n.c",
			"a.#b.c;",
			"a\n.#b\n.c;",
		},
		{ // 307
			"class\na\n{\nb\n}",
			"class a {\n\tb;\n}",
			"class a {\n\tb;\n}",
		},
		{ // 308
			"class a { b () {} }",
			"class a {\n\tb() {}\n}",
			"class a {\n\tb() {}\n}",
		},
		{ // 309
			"class\na\n{\n#b\n}",
			"class a {\n\t#b;\n}",
			"class a {\n\t#b;\n}",
		},
		{ // 310
			"class a { #b () {} }",
			"class a {\n\t#b() {}\n}",
			"class a {\n\t#b() {}\n}",
		},
		{ // 311
			"class a { #b = 1 }",
			"class a {\n\t#b = 1;\n}",
			"class a {\n\t#b = 1;\n}",
		},
		{ // 312
			"class a { #b = 1; #c = 2 }",
			"class a {\n\t#b = 1;\n\t#c = 2;\n}",
			"class a {\n\t#b = 1;\n\t#c = 2;\n}",
		},
		{ // 313
			"class a { #b = 1\n#c = 2 }",
			"class a {\n\t#b = 1;\n\t#c = 2;\n}",
			"class a {\n\t#b = 1;\n\t#c = 2;\n}",
		},
		{ // 314
			"class a { #b(){}#c = 2 }",
			"class a {\n\t#b() {}\n\t#c = 2;\n}",
			"class a {\n\t#b() {}\n\t#c = 2;\n}",
		},
		{ // 315
			"class a { #b\n#c(){}}",
			"class a {\n\t#b;\n\t#c() {}\n}",
			"class a {\n\t#b;\n\t#c() {}\n}",
		},
		{ // 316
			"class a { #b = 1\n#c(){}}",
			"class a {\n\t#b = 1;\n\t#c() {}\n}",
			"class a {\n\t#b = 1;\n\t#c() {}\n}",
		},
		{ // 317
			"class a { #b = 1;#c(){}}",
			"class a {\n\t#b = 1;\n\t#c() {}\n}",
			"class a {\n\t#b = 1;\n\t#c() {}\n}",
		},
		{ // 318
			"class a { #b;#c(){}}",
			"class a {\n\t#b;\n\t#c() {}\n}",
			"class a {\n\t#b;\n\t#c() {}\n}",
		},
		{ // 319
			"class a {static a;static b\nstatic c = 2;static d(){};static{}static{e}static{e;f}}",
			"class a {\n\tstatic a;\n\tstatic b;\n\tstatic c = 2;\n\tstatic d() {}\n\tstatic {}\n\tstatic {\n\t\te;\n\t}\n\tstatic {\n\t\te;\n\t\tf;\n\t}\n}",
			"class a {\n\tstatic a;\n\tstatic b;\n\tstatic c = 2;\n\tstatic d() {}\n\tstatic {}\n\tstatic {\n\t\te;\n\t}\n\tstatic {\n\t\te;\n\t\tf;\n\t}\n}",
		},
		{ // 320
			"#a in b",
			"#a in b;",
			"#a in b;",
		},
		{ // 321
			"#a\nin\nb",
			"#a in b;",
			"#a in b;",
		},
		{ // 322
			"a().#b",
			"a().#b;",
			"a().#b;",
		},
		{ // 323
			"a\n(\n)\n.\n#b",
			"a().#b;",
			"a()\n.#b;",
		},
		{ // 324
			"a\n?.\n#b",
			"a?.#b;",
			"a?.#b;",
		},
		{ // 325
			"a?.c.#d",
			"a?.c.#d;",
			"a?.c.#d;",
		},
		{ // 326
			"// A\n// B\n\na();\n// C\n// D\n",
			"a();",
			"// A\n// B\n\na();\n// C\n// D\n",
		},
		{ // 327
			"/* A *//* B */\n// C\n\na();\n// D\n/* E */   /* F */\n",
			"a();",
			"/* A */ /* B */\n// C\n\na();\n// D\n/* E */ /* F */",
		},
		{ // 328
			"/*\nA\n*//* B */\n// C\n\na();\n// D\n/* E */   /*\n\nF\n\n*/\n",
			"a();",
			"/*\nA\n*/ /* B */\n// C\n\na();\n// D\n/* E */ /*\n\nF\n\n*/",
		},
		{ // 329
			"/*\nA\n*//* B */\n// C\n\n// D\na(); // E\n// F\n\n// G\n/* H */   /*\n\nI\n\n*/\n",
			"a();",
			"/*\nA\n*/ /* B */\n// C\n\n// D\na(); // E\n     // F\n\n// G\n/* H */ /*\n\nI\n\n*/",
		},
		{ // 330
			"// A\n\n// B\nsuper // C\n[ // D\n1\n // E\n]// F\n",
			"super[1];",
			"// A\n\n// B\nsuper // C\n[ // D\n\n\t1\n// E\n] // F\n",
		},
		{ // 331
			"// A\n\n// B\nsuper /* C */ . /* D */ a // E\n",
			"super.a;",
			"// A\n\n// B\nsuper /* C */ . /* D */ a // E\n",
		},
		{ // 332
			"// A\n\n// B\nnew /* C */./* D */target /* E */",
			"new.target;",
			"// A\n\n// B\nnew /* C */ . /* D */ target /* E */;",
		},
		{ // 333
			"// A\n\n/* B */import/* C */./* D */meta/* E */",
			"import.meta;",
			"// A\n\n/* B */ import /* C */ . /* D */ meta /* E */;",
		},
		{ // 334
			"// A\n\n// B\nnew/* C */1/* D */() // E\n",
			"new 1();",
			"// A\n\n// B\nnew /* C */ 1 /* D */ () // E\n",
		},
		{ // 335
			"// A\n\n// B\na // C\n",
			"a;",
			"// A\n\n// B\na // C\n",
		},
		{ // 336
			"// A\n\n// B\na /* C */``/* D */",
			"a``;",
			"// A\n\n// B\na /* C */ `` /* D */;",
		},
		{ // 337
			"// A\n\n/* B */a/* C */./* D */#b/* E */",
			"a.#b;",
			"// A\n\n/* B */ a /* C */ . /* D */ #b /* E */;",
		},
		{ // 338
			"// A\n\n/* B */a/* C */./* D */#b/* E */./* F */c // G\n",
			"a.#b.c;",
			"// A\n\n/* B */ a /* C */ . /* D */ #b /* E */ . /* F */ c // G\n",
		},
		{ // 339
			"// A\n\n// B\na /* C */ . /* D */ #b /* E */ [ // F\n\"c\"\n// G\n] /* H */",
			"a.#b[\"c\"];",
			"// A\n\n// B\na /* C */ . /* D */ #b /* E */ [ // F\n\n\t\"c\"\n// G\n] /* H */;",
		},
		{ // 340
			"super[ // C\n\n// D\n1 // E\n\n// F\n]",
			"super[1];",
			"super[ // C\n\n\t// D\n\t1 // E\n\n// F\n];",
		},
		{ // 341
			"// A\n\n// B\na /* C */ . /* D */ #b /* E */ [ // F\n\n// G\n\"c\" // H\n\n// I\n] /* J */",
			"a.#b[\"c\"];",
			"// A\n\n// B\na /* C */ . /* D */ #b /* E */ [ // F\n\n\t// G\n\t\"c\" // H\n\n// I\n] /* J */;",
		},
		{ // 342
			"a( // A\n\n// B\n)",
			"a();",
			"a( // A\n\n\n// B\n);",
		},
		{ // 343
			"a( // A\n\n// B\nb // C\n\n// D\n, // E\nc // F\n\n//G\n)",
			"a(b, c);",
			"a( // A\n\n\t// B\n\tb // C\n\n\t// D\n\t,\n\t// E\n\tc // F\n\n//G\n);",
		},
		{ // 344
			"a( // A\n\n// B\nb // C\n\n// D\n, // E\n...// F\nc // G\n\n//H\n)",
			"a(b, ...c);",
			"a( // A\n\n\t// B\n\tb // C\n\n\t// D\n\t,\n\t// E\n\t... // F\n\tc // G\n\n//H\n);",
		},
		{ // 345
			"( // A\n\n// B\na // C\n\n// D\n, // E\nb // F\n\n//G\n)",
			"(a, b);",
			"( // A\n\n\t// B\n\ta // C\n\n\t// D\n\t,\n\t// E\n\tb // F\n\n//G\n);",
		},
		{ // 346
			"[\n// A\n...// B\na // C\n]",
			"[...a];",
			"[\n\t// A\n\t... // B\n\ta // C\n];",
		},
		{ // 347
			"[ // A\n// B\n\n// C\na // D\n\n// E\n, // F\n\n// G\nb // H\n// I\n\n// J\n]",
			"[a, b];",
			"[ // A\n  // B\n\n\t// C\n\ta // D\n\n\t// E\n\t,\n\t// F\n\n\t// G\n\tb // H\n\t  // I\n\n// J\n];",
		},
		{ // 348
			"var a = {[ // A\n\n// B\nb // C\n\n// D\n]:c}",
			"var a = {[b]: c};",
			"var a = {[ // A\n\n\t// B\n\tb // C\n\n// D\n]: c};",
		},
		{ // 349
			"var a = {[ // A\n\nb]:c}",
			"var a = {[b]: c};",
			"var a = {[ // A\n\n\tb\n]: c};",
		},
		{ // 350
			"() => { // A\n\n// B\na // C\n\n// D\n}",
			"() => {\n\ta;\n};",
			"() => { // A\n\n\t// B\n\ta // C\n\n// D\n};",
		},
		{ // 351
			"let [\n// A\na// B\n] = b;",
			"let [a] = b;",
			"let [\n\t// A\n\ta // B\n] = b;",
		},
		{ // 352
			"let {\n// A\na// B\n} = b;",
			"let {a} = b;",
			"let {\n\ta: // A\n\ta // B\n} = b;",
		},
		{ // 353
			"let {\n// A\na// B\n: // C\nb // D\n} = c;",
			"let {a: b} = c;",
			"let {\n\t// A\n\ta // B\n\t: // C\n\tb // D\n} = c;",
		},
		{ // 354
			"let {\n// A\na// B\n= // C\nb // D\n} = c;",
			"let {a = b} = c;",
			"let {\n\ta: // A\n\ta // B\n\t= // C\n\tb // D\n} = c;",
		},
		{ // 355
			"let {\n// A\na// B\n: // C\nb // D\n= // E\nc // F\n} = d;",
			"let {a: b = c} = d;",
			"let {\n\t// A\n\ta // B\n\t: // C\n\tb // D\n\t= // E\n\tc // F\n} = d;",
		},
		{ // 356
			"let { // A\n\n// B\na // C\n\n// D\n} = b",
			"let {a} = b;",
			"let { // A\n\n\ta: // B\n\ta // C\n\n// D\n} = b;",
		},
		{ // 357
			"let { // A\n\n// B\n...// C\na // D\n// E\n\n//F\n} = b",
			"let {...a} = b;",
			"let { // A\n\n\t// B\n\t... // C\n\ta // D\n\t  // E\n\n//F\n} = b;",
		},
		{ // 358
			"let { // A\n\n// B\na // C\n, // D\n...// E\nb // F\n\n// G\n} = c",
			"let {a, ...b} = c;",
			"let { // A\n\n\ta: // B\n\ta // C\n\t,\n\t// D\n\t... // E\n\tb // F\n\n// G\n} = c;",
		},
		{ // 359
			"let [ // A\n\n// B\n] = a",
			"let [] = a;",
			"let [ // A\n\n\n// B\n] = a;",
		},
		{ // 360
			"let [ // A\n\n// B\n...// C\na // D\n\n// E\n] = b",
			"let [...a] = b;",
			"let [ // A\n\n\t// B\n\t... // C\n\ta // D\n\n// E\n] = b;",
		},
		{ // 361
			"let [ // A\n\n// B\na // C\n, // D\n\n// E\n, // F\n... // G\nb // H\n\n// I\n] = c",
			"let [a, , ...b] = c;",
			"let [ // A\n\n\t// B\n\ta // C\n\t,\n\t// D\n\n\t// E\n\t,\n\t// F\n\t... // G\n\tb // H\n\n// I\n] = c;",
		},
		{ // 362
			"function a( // A\n\n// B\n){}",
			"function a() {}",
			"function a( // A\n\n// B\n) {}",
		},
		{ // 363
			"function a( // A\n\n// B\n... // C\nb // D\n\n// E\n){}",
			"function a(...b) {}",
			"function a( // A\n\n\t// B\n\t... // C\n\tb // D\n\n// E\n) {}",
		},
		{ // 364
			"function a( // A\n\n// B\nb // C\n\n// D\n){}",
			"function a(b) {}",
			"function a( // A\n\n\t// B\n\tb // C\n\n// D\n) {}",
		},
		{ // 365
			"function a( // A\n\n// B\nb // C\n, // D\nc // E\n\n// F\n){}",
			"function a(b, c) {}",
			"function a( // A\n\n\t// B\n\tb // C\n\t,\n\t// D\n\tc // E\n\n// F\n) {}",
		},
		{ // 366
			"function a( // A\n\n// B\nb // C\n, // D\n... // E\nc // F\n\n// G\n){}",
			"function a(b, ...c) {}",
			"function a( // A\n\n\t// B\n\tb // C\n\t,\n\t// D\n\t... // E\n\tc // F\n\n// G\n) {}",
		},
		{ // 367
			"function a( // A\n\n// B\n... // C\n[]// D\n\n// E\n){}",
			"function a(...[]) {}",
			"function a( // A\n\n\t// B\n\t... // C\n\t[] // D\n\n// E\n) {}",
		},
		{ // 368
			"function a( // A\n\n// B\n... // C\n{}// D\n\n// E\n){}",
			"function a(...{}) {}",
			"function a( // A\n\n\t// B\n\t... // C\n\t{} // D\n\n// E\n) {}",
		},
		{ // 369
			"class a {\n// A\nb /* B */(){}\n}",
			"class a {\n\tb() {}\n}",
			"class a {\n\t// A\n\tb /* B */() {}\n}",
		},
		{ // 370
			"class a {\n// A\na /* B */(){}\n/* C */ b// D\n(){}\n}",
			"class a {\n\ta() {}\n\tb() {}\n}",
			"class a {\n\t// A\n\ta /* B */() {}\n\t/* C */ b // D\n\t() {}\n}",
		},
		{ // 371
			"class a {static //A\nb /* B */() {} }",
			"class a {\n\tstatic b() {}\n}",
			"class a {\n\tstatic //A\n\tb /* B */() {}\n}",
		},
		{ // 372
			"class a {static /* A */ [\"b\"]// B\n() {} }",
			"class a {\n\tstatic [\"b\"]() {}\n}",
			"class a {\n\tstatic /* A */ [\"b\"] // B\n\t() {}\n}",
		},
		{ // 373
			"class a {static // A\n#b// B\n() {} }",
			"class a {\n\tstatic #b() {}\n}",
			"class a {\n\tstatic // A\n\t#b // B\n\t() {}\n}",
		},
		{ // 374
			"class a {static // A\nb// B\n}",
			"class a {\n\tstatic b;\n}",
			"class a {\n\tstatic // A\n\tb // B\n}",
		},
		{ // 375
			"class a {static // A\nb // B\n= 1}",
			"class a {\n\tstatic b = 1;\n}",
			"class a {\n\tstatic // A\n\tb // B\n\t= 1;\n}",
		},
		{ // 376
			"class a {static // A\n[b] // B\n}",
			"class a {\n\tstatic [b];\n}",
			"class a {\n\tstatic // A\n\t[b] // B\n}",
		},
		{ // 377
			"// A\n\n// B\nnew // C\na // D\n",
			"new a;",
			"// A\n\n// B\nnew // C\na // D\n",
		},
		{ // 378
			"// A\n\n// B\nnew // C\na // D\n() // E\n",
			"new a();",
			"// A\n\n// B\nnew // C\na // D\n() // E\n",
		},
		{ // 379
			"// A\n\n// B\nnew // C\nnew // D\na // E\n() // F\n",
			"new new a();",
			"// A\n\n// B\nnew // C\nnew // D\na // E\n() // F\n",
		},
		{ // 380
			"({\n// A\nget // B\na // C\n( // D\n\n// E\n) // F\n{} // G\n})",
			"({get a() {}});",
			"({\n\t// A\n\tget // B\n\ta // C\n\t( // D\n\n\t// E\n\t) // F\n\t{} // G\n});",
		},
		{ // 381
			"({\n// A\nset // B\na // C\n( // D\n\n// E\nb // F\n\n// G\n) // H\n{} // I\n})",
			"({set a(b) {}});",
			"({\n\t// A\n\tset // B\n\ta // C\n\t( // D\n\n\t\t// E\n\t\tb // F\n\n\t// G\n\t) // H\n\t{} // I\n});",
		},
		{ // 382
			"({\n// A\na // B\n( // C\n\n// D\nb // E\n\n// F\n) // G\n{} // H\n})",
			"({a(b) {}});",
			"({\n\t// A\n\ta // B\n\t( // C\n\n\t\t// D\n\t\tb // E\n\n\t// F\n\t) // G\n\t{} // H\n});",
		},
		{ // 383
			"({\n// A\nasync /* B */ a // C\n( // D\n\n// E\n) // F\n{} // G\n})",
			"({async a() {}});",
			"({\n\t// A\n\tasync /* B */ a // C\n\t( // D\n\n\t// E\n\t) // F\n\t{} // G\n});",
		},
		{ // 384
			"({\n// A\n* // B\na // C\n() {}})",
			"({* a() {}});",
			"({\n\t// A\n\t* // B\n\ta // C\n\t() {}\n});",
		},
		{ // 385
			"({\n// A\nasync /* B*/ * // C\na(){}})",
			"({async * a() {}});",
			"({\n\t// A\n\tasync /* B*/ * // C\n\ta() {}\n});",
		},
		{ // 386
			"({\n// A\nasync // B\n(){}})",
			"({async() {}});",
			"({\n\t// A\n\tasync // B\n\t() {}\n});",
		},
		{ // 387
			"({\n// A\nget // B\n(){}})",
			"({get() {}});",
			"({\n\t// A\n\tget // B\n\t() {}\n});",
		},
		{ // 388
			"({\n// A\n... // B\na // C\n})",
			"({...a});",
			"({\n\t// A\n\t... // B\n\ta // C\n});",
		},
		{ // 389
			"({\n// A\na // B\n,})",
			"({a});",
			"({\n\t// A\n\ta // B\n\t: a\n});",
		},
		{ // 390
			"({\n// A\na // B\n= // C\nb // D\n})",
			"({a = b});",
			"({\n\t// A\n\ta // B\n\t= // C\n\tb // D\n});",
		},
		{ // 391
			"({\n// A\na // B\n: // C\nb // D\n})",
			"({a: b});",
			"({\n\t// A\n\ta // B\n\t: // C\n\tb // D\n});",
		},
		{ // 392
			"({\n// A\n[ // B\na\n// C\n] // D\n: // E\nb // F\n})",
			"({[a]: b});",
			"({\n\t// A\n\t[ // B\n\n\t\ta\n\t// C\n\t] // D\n\t: // E\n\tb // F\n});",
		},
		{ // 393
			"({ // A\n// B\n\n// C\n})",
			"({});",
			"({ // A\n   // B\n\n\n// C\n});",
		},
		{ // 394
			"({ // A\n// B\n\n// C\na // D\n// E\n\n// F\n})",
			"({a});",
			"({ // A\n   // B\n\n\t// C\n\ta // D\n\t  // E\n\t: a\n// F\n});",
		},
		{ // 395
			"({ // A\n\n// B\na // C\n, // D\nb // E\n\n// F\n})",
			"({a, b});",
			"({ // A\n\n\t// B\n\ta // C\n\t: a,\n\t// D\n\tb // E\n\t: b\n// F\n});",
		},
		{ // 396
			"function // A\na // B\n()// C\n{}",
			"function a() {}",
			"function // A\na // B\n() // C\n{}",
		},
		{ // 397
			"async /* A */ function // B\na // C\n()// D\n{}",
			"async function a() {}",
			"async /* A */ function // B\na // C\n() // D\n{}",
		},
		{ // 398
			"function // A\n* // B\na // C\n()// D\n{}",
			"function* a() {}",
			"function // A\n* // B\na // C\n() // D\n{}",
		},
		{ // 399
			"(function // A\n()// B\n{})",
			"(function () {});",
			"(function // A\n() // B\n{});",
		},
		{ // 400
			"(async /* A */ function // B\n()// C\n{})",
			"(async function () {});",
			"(async /* A */ function // B\n() // C\n{});",
		},
		{ // 401
			"(function // A\n* // B\n()// C\n{})",
			"(function* () {});",
			"(function // A\n* // B\n() // C\n{});",
		},
		{ // 402
			"(async /* A */ function // B\n* // C\n()// D\n{})",
			"(async function* () {});",
			"(async /* A */ function // B\n* // C\n() // D\n{});",
		},
		{ // 403
			"`${ // A\n\n// B\na // C\n\n// D\n}`",
			"`${a}`;",
			"`${ // A\n\n// B\na // C\n\n// D\n}`;",
		},
		{ // 404
			"`${ // A\na // B\n\n// C\n}${ // D\nb // E\n}`",
			"`${a}${b}`;",
			"`${ // A\na // B\n\n// C\n}${ // D\nb // E\n}`;",
		},
		{ // 405
			"// A\n\n// B\na /* C */ ++ // D\n",
			"a++;",
			"// A\n\n// B\na /* C */++ // D\n",
		},
		{ // 406
			"// A\n\n// B\na /* C */ -- // D\n",
			"a--;",
			"// A\n\n// B\na /* C */-- // D\n",
		},
		{ // 407
			"// A\n\n// B\n++ // C\na // D\n",
			"++a;",
			"// A\n\n// B\n++ // C\na // D\n",
		},
		{ // 408
			"// A\n\n// B\n-- // C\na // D\n",
			"--a;",
			"// A\n\n// B\n-- // C\na // D\n",
		},
		{ // 409
			"// A\n\n// B\ntypeof // C\na // D\n",
			"typeof a;",
			"// A\n\n// B\ntypeof // C\na // D\n",
		},
		{ // 410
			"// A\n\n// B\nvoid // C\n+ // D\na // E\n",
			"void +a;",
			"// A\n\n// B\nvoid // C\n+ // D\na // E\n",
		},
		{ // 411
			"// A\n\n// B\nsuper // C\n()// D\n",
			"super();",
			"// A\n\n// B\nsuper // C\n() // D\n",
		},
		{ // 412
			"// A\n\n// B\nimport // C\n( // D\n\n// E\na // F\n\n// G\n) // H\n",
			"import(a);",
			"// A\n\n// B\nimport // C\n( // D\n\n\t// E\n\ta // F\n\n// G\n) // H\n",
		},
		{ // 413
			"// A\n\n// B\nsuper // C\n()// D\n`` // E",
			"super()``;",
			"// A\n\n// B\nsuper // C\n() // D\n`` // E\n",
		},
		{ // 414
			"// A\n\n// B\nsuper // C\n()// D\n()// E\n",
			"super()();",
			"// A\n\n// B\nsuper // C\n() // D\n() // E\n",
		},
		{ // 415
			"// A\n\n// B\nsuper // C\n()// D\n[ // E\n\n// F\na // G\n\n// H\n]// I",
			"super()[a];",
			"// A\n\n// B\nsuper // C\n() // D\n[ // E\n\n\t// F\n\ta // G\n\n// H\n] // I\n",
		},
		{ // 416
			"// A\n\n// B\nsuper // C\n()// D\n. // E\na // F",
			"super().a;",
			"// A\n\n// B\nsuper // C\n() // D\n. // E\na // F\n",
		},
		{ // 417
			"// A\n\n// B\nsuper // C\n()// D\n. // E\n#a // F",
			"super().#a;",
			"// A\n\n// B\nsuper // C\n() // D\n. // E\n#a // F\n",
		},
		{ // 418
			"// A\n\n// B\na // C\n()// D\n",
			"a();",
			"// A\n\n// B\na // C\n() // D\n",
		},
		{ // 419
			"// A\n\n// B\nclass a{} // C\n",
			"class a {}",
			"// A\n\n// B\nclass a {} // C\n",
		},
		{ // 420
			"// A\n\n// B\nfunction a(){} // C\n",
			"function a() {}",
			"// A\n\n// B\nfunction a() {} // C\n",
		},
		{ // 421
			"// A\n\n// B\nconst a = 1; // B\n",
			"const a = 1;",
			"// A\n\n// B\nconst a = 1; // B\n",
		},
		{ // 422
			"// A\n\n// B\nlet a = 1; // B\n",
			"let a = 1;",
			"// A\n\n// B\nlet a = 1; // B\n",
		},
		{ // 423
			"a?. // A\n() // B\n",
			"a?.();",
			"a?. // A\n() // B\n",
		},
		{ // 424
			"a?. // A\n[ // B\n\n// C\nb // D\n\n// E\n] // F\n",
			"a?.[b];",
			"a?. // A\n[ // B\n\n\t// C\n\tb // D\n\n// E\n] // F\n",
		},
		{ // 425
			"a?. // A\n[\n// C\nb\n// E\n] // F\n",
			"a?.[b];",
			"a?. // A\n[\n\t// C\n\tb\n// E\n] // F\n",
		},
		{ // 426
			"a?. // A\nb // B\n",
			"a?.b;",
			"a?. // A\nb // B\n",
		},
		{ // 427
			"a?. // A\n`` // B\n",
			"a?.``;",
			"a?. // A\n`` // B\n",
		},
		{ // 428
			"a?. // A\n()// B\n`` // C\n",
			"a?.()``;",
			"a?. // A\n() // B\n`` // C\n",
		},
		{ // 429
			"a?. // A\n`` // B\n[ // C\n\n// D\nb // E\n\n// F\n] // G\n",
			"a?.``[b];",
			"a?. // A\n`` // B\n[ // C\n\n\t// D\n\tb // E\n\n// F\n] // G\n",
		},
		{ // 430
			"a?. //A\n[ // B\n\n// C\nb // D\n\n// E\n] // F\n. // G\nc// H\n",
			"a?.[b].c;",
			"a?. //A\n[ // B\n\n\t// C\n\tb // D\n\n// E\n] // F\n. // G\nc // H\n",
		},
		{ // 431
			"class a {\n// A\nstatic // B\n{} // C\n}",
			"class a {\n\tstatic {}\n}",
			"class a {\n\t// A\n\tstatic // B\n\t{} // C\n}",
		},
		{ // 432
			"class // A\na // B\n{ // C\n\n// D\n}",
			"class a {}",
			"class // A\na // B\n{ // C\n\n// D\n}",
		},
		{ // 433
			"class a { // A\n; // B\n; // C\n; // D\na(){} // E\n; // F\n;\n // G\n}",
			"class a {\n\ta() {}\n}",
			"class a { // A\n\n\t// B\n\t// C\n\t// D\n\ta() {} // E\n\n// F\n\n// G\n}",
		},
		{ // 434
			"class a { // A\n\n// B\na(){} // C\n// D\n\n// E\nb(){} // F\n\n// G\n}",
			"class a {\n\ta() {}\n\tb() {}\n}",
			"class a { // A\n\n\t// B\n\ta() {} // C\n\t       // D\n\n\t// E\n\tb() {} // F\n\n// G\n}",
		},
		{ // 435
			"let // A\na // B\n",
			"let a;",
			"let // A\na // B\n",
		},
		{ // 436
			"let // A\na // B\n= // C\nb // D\n",
			"let a = b;",
			"let // A\na // B\n= // C\nb // D\n",
		},
		{ // 437
			"let // A\n[]// B\n = // C\na // D",
			"let [] = a;",
			"let // A\n[] // B\n= // C\na // D\n",
		},
		{ // 438
			"let // A\n{}// B\n = // C\na // D",
			"let {} = a;",
			"let // A\n{} // B\n= // C\na // D\n",
		},
		{ // 439
			"continue /* A */ a // B\n",
			"continue a;",
			"continue /* A */ a // B\n",
		},
		{ // 440
			"continue /* A */ a // B\n;",
			"continue a;",
			"continue /* A */ a // B\n",
		},
		{ // 441
			"continue /* A */;",
			"continue;",
			"continue /* A */;",
		},
		{ // 442
			"break /* A */ a // B\n",
			"break a;",
			"break /* A */ a // B\n",
		},
		{ // 443
			"break /* A */ a // B\n;",
			"break a;",
			"break /* A */ a // B\n",
		},
		{ // 444
			"break /* A */;",
			"break;",
			"break /* A */;",
		},
		{ // 445
			"function a(){\nreturn /* A */ b // B\n}",
			"function a() {\n\treturn b;\n}",
			"function a() {\n\treturn /* A */ b // B\n}",
		},
		{ // 446
			"function a(){\nreturn /* A */ b // B\n;}",
			"function a() {\n\treturn b;\n}",
			"function a() {\n\treturn /* A */ b // B\n}",
		},
		{ // 447
			"function a(){\nreturn /* A */;\n}",
			"function a() {\n\treturn;\n}",
			"function a() {\n\treturn /* A */;\n}",
		},
		{ // 448
			"debugger // A\n",
			"debugger;",
			"debugger // A\n",
		},
		{ // 449
			"debugger /* A */;",
			"debugger;",
			"debugger /* A */;",
		},
		{ // 450
			"a /* A */: // B\ndebugger;",
			"a: debugger;",
			"a /* A */ : // B\ndebugger;",
		},
		{ // 451
			"if // A\n( // B\n\n// C\na // D\n\n// E\n) // F\nb // G",
			"if (a) b;",
			"if // A\n( // B\n\n\t// C\n\ta // D\n\n// E\n) // F\nb // G\n",
		},
		{ // 452
			"// A\n\n// B\nif // C\n(// D\n\n// E\na // F\n\n// G\n) // H\n{// I\n\n// J\n} // K\nelse // L\n{ // M\n\n// N\n} // O",
			"if (a) {} else {}",
			"// A\n\n// B\nif // C\n( // D\n\n\t// E\n\ta // F\n\n// G\n) // H\n{ // I\n\n// J\n} // K\nelse // L\n{ // M\n\n// N\n} // O\n",
		},
		{ // 453
			"// A\na /* B */ || // C\nb // D\n|| c /* E */",
			"a || b || c;",
			"// A\n\na /* B */ || // C\nb // D\n|| c /* E */;",
		},
		{ // 454
			"with // A\n( // B\n\n// C\na // D\n\n// E\n) // F\nb // G\n",
			"with (a) b;",
			"with // A\n( // B\n\n\t// C\n\ta // D\n\n// E\n) // F\nb // G\n",
		},
		{ // 455
			"while // A\n( // B\n\n// C\na // D\n\n// E\n) // F\nb // G\n",
			"while (a) b;",
			"while // A\n( // B\n\n\t// C\n\ta // D\n\n// E\n) // F\nb // G\n",
		},
		{ // 456
			"switch (a) {\n// A\ncase /* B */ b /* C */: // D\n\n// E\ncase c:\n// F\ncase /* G */ d /* H */: // I\n\n// J\n{} // F\ndefault:}",
			"switch (a) {\ncase b:\ncase c:\ncase d:\n\t{}\ndefault:\n}",
			"switch (a) {\n// A\ncase /* B */ b /* C */: // D\n\n// E\ncase c:\n// F\ncase /* G */ d /* H */: // I\n\n\t// J\n\t{} // F\n\ndefault:\n}",
		},
		{ // 457
			"switch // A\n( // B\n\n// C\na // D\n\n// E\n) // F\n{ // G\n\n// H\n}",
			"switch (a) {}",
			"switch // A\n( // B\n\n\t// C\n\ta // D\n\n// E\n) // F\n{ // G\n\n// H\n}",
		},
		{ // 458
			"switch // A\n( // B\n\n// C\na // D\n\n// E\n) // F\n{ // G\n\n// H\ncase /* I */ b /* J */ : // K\n\n// J\ncase c:// L\n\n// M\ndefault /* N */ : // O\n\td;// P\n\n// Q\n}",
			"switch (a) {\ncase b:\ncase c:\ndefault:\n\td;\n}",
			"switch // A\n( // B\n\n\t// C\n\ta // D\n\n// E\n) // F\n{ // G\n\n// H\ncase /* I */ b /* J */: // K\n\n// J\ncase c: // L\n\n// M\ndefault /* N */: // O\n\n\td; // P\n\n// Q\n}",
		},
		{ // 459
			"try // A\n{} // B\ncatch // C\n{}",
			"try {} catch {}",
			"try // A\n{} // B\ncatch // C\n{}",
		},
		{ // 460
			"try {}catch // A\n( // B\n\n// C\na // D\n\n// E\n) // F\n{} // G",
			"try {} catch (a) {}",
			"try {} catch // A\n( // B\n\n\t// C\n\ta // D\n\n// E\n) // F\n{} // G\n",
		},
		{ // 461
			"try{} // A\nfinally // B\n{} // C",
			"try {} finally {}",
			"try {} // A\nfinally // B\n{} // C\n",
		},
		{ // 462
			"try // A\n{}// B\ncatch /* C */ ( // D\n\n// E\na // F\n\n// G\n) // H\n{}// I\nfinally /* J */ {} // K",
			"try {} catch (a) {} finally {}",
			"try // A\n{} // B\ncatch /* C */ ( // D\n\n\t// E\n\ta // F\n\n// G\n) // H\n{} // I\nfinally /* J */ {} // K\n",
		},
		{ // 463
			"for // A\n( // B\n\n// C\n; // D\n; // E\n\n// F\n) // G\na",
			"for (;;) a;",
			"for // A\n( // B\n\n\t// C\n\t; // D\n\t; // E\n\n// F\n) // G\na;",
		},
		{ // 464
			"for ( // A\n\n// B\nvar // C\na // D\n, // E\nb // F\n; // G\nc // H\n; // I\n\n// J\n) // K\n{}",
			"for (var a, b; c;) {}",
			"for ( // A\n\n\t// B\n\tvar // C\n\ta // D\n\t,\n\t// E\n\tb // F\n\t; // G\n\tc // H\n\t; // I\n\n// J\n) // K\n{}",
		},
		{ // 465
			"for ( // A\n\n// B\n; // C\na // D\n; // E\nb // F\n\n// G\n) // H\nc",
			"for (; a; b) c;",
			"for ( // A\n\n\t// B\n\t; // C\n\ta // D\n\t; // E\n\tb // F\n\n// G\n) // H\nc;",
		},
		{ // 466
			"for ( // A\n\n// B\nlet // C\na // D\n; b; c) {}",
			"for (let a; b; c) {}",
			"for ( // A\n\n\t// B\n\tlet // C\n\ta // D\n\t; b; c\n) {}",
		},
		{ // 467
			"for ( // A\n\n// B\nvar // C\n{}// D\n= // E\na // F\n;;);",
			"for (var {} = a;;) ;",
			"for ( // A\n\n\t// B\n\tvar // C\n\t{} // D\n\t= // E\n\ta // F\n\t;;\n) ;",
		},
		{ // 468
			"for ( // A\n\n// B\nvar // C\n[]// D\n= a // E\n; // F\n; // G\n) // H\n;",
			"for (var [] = a;;) ;",
			"for ( // A\n\n\t// B\n\tvar // C\n\t[] // D\n\t= a // E\n\t; // F\n\t; // G\n) // H\n;",
		},
		{ // 469
			"for (\n// A\nvar // B\na // C\nin // D\nb // E\n\n// F\n);",
			"for (var a in b) ;",
			"for (\n\t// A\n\tvar // B\n\ta // C\n\tin // D\n\tb // E\n\n// F\n) ;",
		},
		{ // 470
			"[ // A\n\n// B\na.b // C\n, // D\na.c // E\n\n// F\n] = b",
			"[a.b, a.c] = b;",
			"[ // A\n\n\t// B\n\ta.b // C\n\t,\n\t// D\n\ta.c // E\n\n// F\n] = b;",
		},
		{ // 471
			"[ // A\n\n// B\n... // C\na // D\n\n// E\n] = b",
			"[...a] = b;",
			"[ // A\n\n\t// B\n\t... // C\n\ta // D\n\n// E\n] = b;",
		},
		{ // 472
			"({\n// A\na // B\n} = b)",
			"({a} = b);",
			"({\n\t// A\n\ta // B\n\t: a\n} = b);",
		},
		{ // 473
			"({\n// A\na // B\n= // C\nb // D\n} = c)",
			"({a = b} = c);",
			"({\n\t// A\n\ta // B\n\t: a = // C\n\tb // D\n} = c);",
		},
		{ // 474
			"({\n// A\na // B\n: // C\nb // D\n= // E\nc // F\n} = d)",
			"({a: b = c} = d);",
			"({\n\t// A\n\ta // B\n\t: // C\n\tb // D\n\t= // E\n\tc // F\n} = d);",
		},
		{ // 475
			"({ // A\n\n// B\na // C\n\n// D\n} = b)",
			"({a} = b);",
			"({ // A\n\n\t// B\n\ta // C\n\t: a\n// D\n} = b);",
		},
		{ // 476
			"({\n// A\n... // B\na // C\n} = b)",
			"({...a} = b);",
			"({\n\t// A\n\t... // B\n\ta // C\n} = b);",
		},
		{ // 477
			"(\n// A\n[a] // B\n= b)",
			"([a] = b);",
			"(\n\t// A\n\t[a] // B\n\t= b\n);",
		},
		{ // 478
			"(\n// A\n{a} // B\n= b)",
			"({a} = b);",
			"(\n\t// A\n\t{a: a} // B\n\t= b\n);",
		},
		{ // 479
			"(\n// A\na /* B */ => // C\n{} // D\n)",
			"(a => {});",
			"(\n\t// A\n\ta /* B */ => // C\n\t{} // D\n);",
		},
		{ // 480
			"(\n// A\n() /* B */ => /* C */ {} // D\n)",
			"(() => {});",
			"(\n\t// A\n\t() /* B */ => /* C */ {} // D\n);",
		},
		{ // 481
			"(\n// A\nasync /* B */ a /* C */ => // D\nb // E\n)",
			"(async a => b);",
			"(\n\t// A\n\tasync /* B */ a /* C */ => // D\n\tb // E\n);",
		},
		{ // 482
			"(// A\n\n// B\nasync /* C */ ()/* D */ => // E\n{}// F\n\n// G\n)",
			"(async () => {});",
			"( // A\n\n\t// B\n\tasync /* C */ () /* D */ => // E\n\t{} // F\n\n// G\n);",
		},
		{ // 483
			"function *a() {\n// A\nyield /* B */ a // C\n}",
			"function* a() {\n\tyield a;\n}",
			"function* a() {\n\t// A\n\tyield /* B */ a // C\n}",
		},
		{ // 484
			"function* a() {\n// A\nyield /* B */ * /* C */ a // D\n}",
			"function* a() {\n\tyield * a;\n}",
			"function* a() {\n\t// A\n\tyield /* B */ * /* C */ a // D\n}",
		},
		{ // 485
			"a // A\n?? b",
			"a ?? b;",
			"a // A\n?? b;",
		},
		{ // 486
			"a /* A */ ?? b",
			"a ?? b;",
			"a /* A */ ?? b;",
		},
		{ // 487
			"a /* A */ ? /* B */ b /* C */ : /* D */ c",
			"a ? b : c;",
			"a /* A */ ? /* B */ b /* C */ : /* D */ c;",
		},
		{ // 488
			"a // A\n? // B\nb // C\n: // D\nc",
			"a ? b : c;",
			"a // A\n? // B\nb // C\n: // D\nc;",
		},
		{ // 489
			"a /* A */ && /* B */ b",
			"a && b;",
			"a /* A */ && /* B */ b;",
		},
		{ // 490
			"a // A\n&& // B\nb",
			"a && b;",
			"a // A\n&& // B\nb;",
		},
		{ // 491
			"a /* A */ | /* B */ b",
			"a | b;",
			"a /* A */ | /* B */ b;",
		},
		{ // 492
			"a // A\n| // B\nb",
			"a | b;",
			"a // A\n| // B\nb;",
		},
		{ // 493
			"a /* A */ ^ /* B */ b",
			"a ^ b;",
			"a /* A */ ^ /* B */ b;",
		},
		{ // 494
			"a // A\n^ // B\nb",
			"a ^ b;",
			"a // A\n^ // B\nb;",
		},
		{ // 495
			"a /* A */ & /* B */ b",
			"a & b;",
			"a /* A */ & /* B */ b;",
		},
		{ // 496
			"a // A\n& // B\nb",
			"a & b;",
			"a // A\n& // B\nb;",
		},
		{ // 497
			"a /* A */ == /* B */ b",
			"a == b;",
			"a /* A */ == /* B */ b;",
		},
		{ // 498
			"a // A\n!== // B\nb",
			"a !== b;",
			"a // A\n!== // B\nb;",
		},
		{ // 499
			"a /* A */ < /* B */ b",
			"a < b;",
			"a /* A */ < /* B */ b;",
		},
		{ // 500
			"a // A\nin // B\nb",
			"a in b;",
			"a // A\nin // B\nb;",
		},
		{ // 501
			"a /* A */ << /* B */ b",
			"a << b;",
			"a /* A */ << /* B */ b;",
		},
		{ // 502
			"a // A\n>>> // B\nb",
			"a >>> b;",
			"a // A\n>>> // B\nb;",
		},
		{ // 503
			"a /* A */ + /* B */ b",
			"a + b;",
			"a /* A */ + /* B */ b;",
		},
		{ // 504
			"a // A\n- // B\nb",
			"a - b;",
			"a // A\n- // B\nb;",
		},
		{ // 505
			"a /* A */ * /* B */ b",
			"a * b;",
			"a /* A */ * /* B */ b;",
		},
		{ // 506
			"a // A\n% // B\nb",
			"a % b;",
			"a // A\n% // B\nb;",
		},
		{ // 507
			"a /* A */ ** /* B */ b",
			"a ** b;",
			"a /* A */ ** /* B */ b;",
		},
		{ // 508
			"a // A\n** // B\nb",
			"a ** b;",
			"a // A\n** // B\nb;",
		},
		{ // 509
			"do /* A */a /* B */; /* C */ while ( // D\n\n// E\nb // F\n\n// G\n) // H",
			"do a; while (b);",
			"do /* A */ a /* B */; /* C */ while ( // D\n\n\t// E\n\tb // F\n\n// G\n) // H\n",
		},
		{ // 510
			"do /* A */{}/* B */ while ( // C\n\n// D\nb // E\n\n// F\n) // G",
			"do {} while (b);",
			"do /* A */ {} /* B */ while ( // C\n\n\t// D\n\tb // E\n\n// F\n) // G\n",
		},
		{ // 511
			"([a // A\n= b] = c)",
			"([a = b] = c);",
			"([\n\ta // A\n\t= b\n] = c);",
		},
		{ // 512
			"// A\n\n// B\n#a /* C */ in /* D */ b // E",
			"#a in b;",
			"// A\n\n// B\n#a /* C */ in /* D */ b // E\n",
		},
		{ // 513
			"a({b: () => {\nc();\n}});",
			"a({b: () => {\n\tc();\n}});",
			"a({b: () => {\n\tc();\n}});",
		},
		{ // 514
			"a({b: // A\n() => {\nc();\n}});",
			"a({b: () => {\n\tc();\n}});",
			"a({\n\tb: // A\n\t() => {\n\t\tc();\n\t}\n});",
		},
		{ // 515
			"a?.[() => { // A\n}]",
			"a?.[() => {}];",
			"a?.[() => { // A\n}];",
		},
		{ // 516
			"a?.[ // A\n() => { // B\n}]",
			"a?.[() => {}];",
			"a?.[ // A\n\n\t() => { // B\n\t}\n];",
		},
		{ // 517
			"a?.[() => { // A\nreturn c;}]",
			"a?.[() => {\n\treturn c;\n}];",
			"a?.[() => { // A\n\n\treturn c;\n}];",
		},
		{ // 518
			"a?.[ // A\n() => { // B\nreturn c;}]",
			"a?.[() => {\n\treturn c;\n}];",
			"a?.[ // A\n\n\t() => { // B\n\n\t\treturn c;\n\t}\n];",
		},
		{ // 519
			"import(() => {})",
			"import(() => {});",
			"import(() => {});",
		},
		{ // 520
			"import(() => {return a})",
			"import(() => {\n\treturn a;\n});",
			"import(() => {\n\treturn a;\n});",
		},
		{ // 521
			"import(// A\n() => {return a})",
			"import(() => {\n\treturn a;\n});",
			"import( // A\n\n\t() => {\n\t\treturn a;\n\t}\n);",
		},
		{ // 522
			"a()[() => {}]",
			"a()[() => {}];",
			"a()[() => {}];",
		},
		{ // 523
			"a()[() => {return b}]",
			"a()[() => {\n\treturn b;\n}];",
			"a()[() => {\n\treturn b;\n}];",
		},
		{ // 524
			"a()[// A\n() => {return b}]",
			"a()[() => {\n\treturn b;\n}];",
			"a()[ // A\n\n\t() => {\n\t\treturn b;\n\t}\n];",
		},
		{ // 525
			"a[() => {}]",
			"a[() => {}];",
			"a[() => {}];",
		},
		{ // 526
			"a[() => {return b}]",
			"a[() => {\n\treturn b;\n}];",
			"a[() => {\n\treturn b;\n}];",
		},
		{ // 527
			"a[// A\n() => {return b}]",
			"a[() => {\n\treturn b;\n}];",
			"a[ // A\n\n\t() => {\n\t\treturn b;\n\t}\n];",
		},
		{ // 528
			"super[() => {}]",
			"super[() => {}];",
			"super[() => {}];",
		},
		{ // 529
			"super[() => {return b}]",
			"super[() => {\n\treturn b;\n}];",
			"super[() => {\n\treturn b;\n}];",
		},
		{ // 530
			"super[// A\n() => {return b}]",
			"super[() => {\n\treturn b;\n}];",
			"super[ // A\n\n\t() => {\n\t\treturn b;\n\t}\n];",
		},
		{ // 531
			"with(() => {}){}",
			"with (() => {}) {}",
			"with (() => {}) {}",
		},
		{ // 532
			"with(() => {return a}) {}",
			"with (() => {\n\treturn a;\n}) {}",
			"with (() => {\n\treturn a;\n}) {}",
		},
		{ // 533
			"with (// A\n() => {return a}) {}",
			"with (() => {\n\treturn a;\n}) {}",
			"with ( // A\n\n\t() => {\n\t\treturn a;\n\t}\n) {}",
		},
	} {
		for m, in := range [2]string{test.Input, test.VerboseOutput} {
			s, err := ParseScript(makeTokeniser(parser.NewStringTokeniser(in)))
			if err != nil {
				t.Errorf("test %d.%d.1: unexpected error: %s", n+1, m+1, err)
				continue
			}

			st.Verbose = false

			st.Reset()
			s.Format(&st, 's')

			if str := st.String(); str != test.SimpleOutput {
				t.Errorf("test %d.%d.2: expecting %q, got %q\n%s", n+1, m+1, test.SimpleOutput, str, s)
			}

			st.Verbose = true

			st.Reset()
			s.Format(&st, 's')

			if str := st.String(); str != test.VerboseOutput {
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
			"import { // A\n\n\t/* B */ a as a /* C */\n\t// D\n\t,\n\t// E\n\tb as b // F\n\n/* G */\n} from './c';",
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
			"import {} from './a' with {\n\t// A\n\tb /* B */: /* C */ \"c\" // D\n\n\t// E\n\t,\n\td: \"e\"\n};",
		},
		{ // 43
			"import {} from './a' with /* A */ {// B\n\n// C\nb/* D */:/* E */\"c\"// F\n\n// G\n, d:\"e\" // H\n\n// I\n};",
			"import {} from './a' with {b: \"c\", d: \"e\"};",
			"import {} from './a' with /* A */ { // B\n\n\t// C\n\tb /* D */: /* E */ \"c\" // F\n\n\t// G\n\t,\n\td: \"e\" // H\n\n// I\n};",
		},
		{ // 44
			"import {} from './a' /* A */ with /* B */ {// C\n\n// D\nb/* E */:/* F */\"c\"// G\n\n// H\n, d:\"e\" // I\n\n// J\n};",
			"import {} from './a' with {b: \"c\", d: \"e\"};",
			"import {} from './a' /* A */ with /* B */ { // C\n\n\t// D\n\tb /* E */: /* F */ \"c\" // G\n\n\t// H\n\t,\n\td: \"e\" // I\n\n// J\n};",
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
			"export {\n\t// A\n\ta /* B */ as /* C */ b // D\n\n\t// E\n\t,\n\t// F\n\tc as c // G\n};",
		},
		{ // 49
			"export { // A\n\n// B\na /* C */ as /* D */ b // E\n\n// F\n, // G\nc // H\n\n// I\n};",
			"export {a as b, c};",
			"export { // A\n\n\t// B\n\ta /* C */ as /* D */ b // E\n\n\t// F\n\t,\n\t// G\n\tc as c // H\n\n// I\n};",
		},
		{ // 50
			"// A\n\n// B\nexport/* C */{ // D\n\n// E\na // F\n\n// G\n, // H\nb // I\n\n// J\n}/* K */from/* L */''/* M */;",
			"export {a, b} from '';",
			"// A\n\n// B\nexport /* C */ { // D\n\n\t// E\n\ta as a // F\n\n\t// G\n\t,\n\t// H\n\tb as b // I\n\n// J\n} /* K */ from /* L */ '' /* M */;",
		},
		{ // 51
			"// A\n\n// B\nexport/* C */*/* D */as/* E */a/* F */from /* G */''/* H */",
			"export * as a from '';",
			"// A\n\n// B\nexport /* C */ * /* D */ as /* E */ a /* F */ from /* G */ '' /* H */;",
		},
		{ // 52
			"// A\n\n// B\nexport/* C */{ // D\n\n// E\na // F\n\n// G\n, // H\nb // I\n\n// J\n}/* K */;",
			"export {a, b};",
			"// A\n\n// B\nexport /* C */ { // D\n\n\t// E\n\ta as a // F\n\n\t// G\n\t,\n\t// H\n\tb as b // I\n\n// J\n} /* K */;",
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
			"export {\n\ta // A\n\tas // B\n\tb\n};",
		},
		{ // 67
			"import {a // A\nas b} from 'c'",
			"import {a as b} from 'c';",
			"import {\n\ta // A\n\tas b\n} from 'c';",
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

			st.Reset()
			s.Format(&st, 's')

			if str := st.String(); str != test.SimpleOutput {
				t.Errorf("test %d.%d.2: expecting %q, got %q\n%s", n+1, m+1, test.SimpleOutput, str, s)
			}

			st.Verbose = true

			st.Reset()
			s.Format(&st, 's')

			if str := st.String(); str != test.VerboseOutput {
				t.Errorf("test %d.%d.3: expecting %q, got %q\n%s", n+1, m+1, test.VerboseOutput, str, s)
			}
		}
	}
}

func TestPrintingJSX(t *testing.T) {
	var st state

	for n, test := range [...]struct {
		Input, SimpleOutput, VerboseOutput string
	}{
		{ // 1
			"<a/>",
			"<a />;",
			"<a />;",
		},
		{ // 2
			"<a.b/>",
			"<a.b />;",
			"<a.b />;",
		},
		{ // 3
			"<a.b.c/>",
			"<a.b.c />;",
			"<a.b.c />;",
		},
		{ // 4
			"<a:b/>",
			"<a:b />;",
			"<a:b />;",
		},
		{ // 5
			"<a>b</a>;",
			"<a>b</a>;",
			"<a>b</a>;",
		},
		{ // 6
			"<a>b\nc</a>;",
			"<a>\n\tb\n\tc\n</a>;",
			"<a>\n\tb\n\tc\n</a>;",
		},
		{ // 7
			"<a><b /></a>",
			"<a><b /></a>;",
			"<a><b /></a>;",
		},
		{ // 8
			"<a><></></a>",
			"<a><></></a>;",
			"<a><></></a>;",
		},
		{ // 9
			"<a>{b}</a>",
			"<a>{b}</a>;",
			"<a>{b}</a>;",
		},
		{ // 10
			"<a>{...b}</a>",
			"<a>{...b}</a>;",
			"<a>{...b}</a>;",
		},
		{ // 11
			"<a>b{...c}<d /></a>",
			"<a>\n\tb\n\t{...c}\n\t<d />\n</a>;",
			"<a>\n\tb\n\t{...c}\n\t<d />\n</a>;",
		},
		{ // 12
			"<a b/>",
			"<a b />;",
			"<a b />;",
		},
		{ // 13
			"<a b:c/>",
			"<a b:c />;",
			"<a b:c />;",
		},
		{ // 14
			"<a b:c='d'/>",
			"<a b:c='d' />;",
			"<a b:c='d' />;",
		},
		{ // 15
			"<a b=<c/>/>",
			"<a b=<c /> />;",
			"<a b=<c /> />;",
		},
		{ // 16
			"<a b=<></>/>",
			"<a b=<></> />;",
			"<a b=<></> />;",
		},
		{ // 17
			"<a b={c}/>",
			"<a b={c} />;",
			"<a b={c} />;",
		},
		{ // 18
			"<a {...b}/>",
			"<a {...b} />;",
			"<a {...b} />;",
		},
		{ // 19
			"<a b c='d' {...e}></a>;",
			"<a b c='d' {...e}></a>;",
			"<a b c='d' {...e}></a>;",
		},
		{ // 20
			"<>b</>;",
			"<>b</>;",
			"<>b</>;",
		},
		{ // 21
			"<>b\nc</>;",
			"<>\n\tb\n\tc\n</>;",
			"<>\n\tb\n\tc\n</>;",
		},
		{ // 22
			"<><b /></>",
			"<><b /></>;",
			"<><b /></>;",
		},
		{ // 23
			"<><></></>",
			"<><></></>;",
			"<><></></>;",
		},
		{ // 24
			"<>{b}</>",
			"<>{b}</>;",
			"<>{b}</>;",
		},
		{ // 25
			"<>{...b}</>",
			"<>{...b}</>;",
			"<>{...b}</>;",
		},
		{ // 26
			"<>b{...c}<d /></>",
			"<>\n\tb\n\t{...c}\n\t<d />\n</>;",
			"<>\n\tb\n\t{...c}\n\t<d />\n</>;",
		},
	} {
		for m, in := range [2]string{test.Input, test.VerboseOutput} {
			tk := parser.NewStringTokeniser(in)

			s, err := ParseModule(AsJSX(&tk))
			if err != nil {
				t.Errorf("test %d.%d.1: unexpected error: %s", n+1, m+1, err)

				continue
			}

			st.Verbose = false

			st.Reset()
			s.Format(&st, 's')

			if str := st.String(); str != test.SimpleOutput {
				t.Errorf("test %d.%d.2: expecting %q, got %q\n%s", n+1, m+1, test.SimpleOutput, str, s)
			}

			st.Verbose = true

			st.Reset()
			s.Format(&st, 's')

			if str := st.String(); str != test.VerboseOutput {
				t.Errorf("test %d.%d.3: expecting %q, got %q\n%s", n+1, m+1, test.VerboseOutput, str, s)
			}
		}
	}
}
