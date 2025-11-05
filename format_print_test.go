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
		{ // 97
			"class a{}",
			"class a {}",
			"class a {}",
		},
		{ // 98
			"class\na\n{\n}\n",
			"class a {}",
			"class a {}",
		},
		{ // 99
			"class a extends b {}",
			"class a extends b {}",
			"class a extends b {}",
		},
		{ // 100
			"class\na\nextends\nb\n{\n}",
			"class a extends b {}",
			"class a extends b {}",
		},
		{ // 101
			"a = class{}",
			"a = class {};",
			"a = class {};",
		},
		{ // 102
			"a\n=\nclass\nb\n{\n}",
			"a = class b {};",
			"a = class b {};",
		},
		{ // 103
			"a\n=\nclass\nextends\nb\n{\n}",
			"a = class extends b {};",
			"a = class extends b {};",
		},
		{ // 104
			"let a = 1",
			"let a = 1;",
			"let a = 1;",
		},
		{ // 105
			"let\na\n=\n1\n",
			"let a = 1;",
			"let a = 1;",
		},
		{ // 106
			"let a=1,b=2,c=3",
			"let a = 1, b = 2, c = 3;",
			"let a = 1,\nb = 2,\nc = 3;",
		},
		{ // 107
			"const a = 1",
			"const a = 1;",
			"const a = 1;",
		},
		{ // 108
			"const\na\n=\n1\n",
			"const a = 1;",
			"const a = 1;",
		},
		{ // 109
			"const a=1,b=2,c=3",
			"const a = 1, b = 2, c = 3;",
			"const a = 1,\nb = 2,\nc = 3;",
		},
		{ // 110
			"let a",
			"let a;",
			"let a;",
		},
		{ // 111
			"let\na\n,\nb\n=\n1\n,\nc\n",
			"let a, b = 1, c;",
			"let a,\nb = 1,\nc;",
		},
		{ // 112
			"const a",
			"const a;",
			"const a;",
		},
		{ // 113
			"const\na\n,\nb\n=\n1\n,\nc\n",
			"const a, b = 1, c;",
			"const a,\nb = 1,\nc;",
		},
		{ // 114
			"let [a]=1",
			"let [a] = 1;",
			"let [a] = 1;",
		},
		{ // 115
			"const\n[\na\n]\n=\n1",
			"const [a] = 1;",
			"const [a] = 1;",
		},
		{ // 116
			"let {a}=1",
			"let {a} = 1;",
			"let {a: a} = 1;",
		},
		{ // 117
			"const\n{\na\n}\n=\n1",
			"const {a} = 1;",
			"const {a: a} = 1;",
		},
		{ // 118
			"function* a() {yield a}",
			"function* a() {\n	yield a;\n}",
			"function* a() { yield a; }",
		},
		{ // 119
			"() => {}",
			"() => {};",
			"() => {};",
		},
		{ // 120
			"a=b",
			"a = b;",
			"a = b;",
		},
		{ // 121
			"a/=b",
			"a /= b;",
			"a /= b;",
		},
		{ // 122
			"a%=b",
			"a %= b;",
			"a %= b;",
		},
		{ // 123
			"a+=b",
			"a += b;",
			"a += b;",
		},
		{ // 124
			"a-=b",
			"a -= b;",
			"a -= b;",
		},
		{ // 125
			"a<<=b",
			"a <<= b;",
			"a <<= b;",
		},
		{ // 126
			"a>>=b",
			"a >>= b;",
			"a >>= b;",
		},
		{ // 127
			"a>>>=b",
			"a >>>= b;",
			"a >>>= b;",
		},
		{ // 128
			"a&=b",
			"a &= b;",
			"a &= b;",
		},
		{ // 129
			"a^=b",
			"a ^= b;",
			"a ^= b;",
		},
		{ // 130
			"a|=b",
			"a |= b;",
			"a |= b;",
		},
		{ // 131
			"a**=b",
			"a **= b;",
			"a **= b;",
		},
		{ // 132
			"a?b:c",
			"a ? b : c;",
			"a ? b : c;",
		},
		{ // 133
			"new a",
			"new a;",
			"new a;",
		},
		{ // 134
			"a()",
			"a();",
			"a();",
		},
		{ // 135
			"var {a} = 1",
			"var {a} = 1;",
			"var {a: a} = 1;",
		},
		{ // 136
			"var { a , b, ...c } = 1",
			"var {a, b, ...c} = 1;",
			"var {a: a, b: b, ...c} = 1;",
		},
		{ // 137
			"var [a] = 1",
			"var [a] = 1;",
			"var [a] = 1;",
		},
		{ // 138
			"var [ a , b, ...c ] = 1",
			"var [a, b, ...c] = 1;",
			"var [a, b, ...c] = 1;",
		},
		{ // 139
			"switch (a) {case 1:}",
			"switch (a) {\ncase 1:\n}",
			"switch (a) {\ncase 1:\n}",
		},
		{ // 140
			"switch (a) {case 1:b;c;d}",
			"switch (a) {\ncase 1:\n	b;\n	c;\n	d;\n}",
			"switch (a) {\ncase 1:\n	b;\n	c;\n	d;\n}",
		},
		{ // 141
			"function a(b){}",
			"function a(b) {}",
			"function a(b) {}",
		},
		{ // 142
			"function\na(\nb\n)\n{\n}\n",
			"function a(b) {}",
			"function a(b) {}",
		},
		{ // 143
			"function a(b,c,...d){}",
			"function a(b, c, ...d) {}",
			"function a(b, c, ...d) {}",
		},
		{ // 144
			"class\na{b(){}c\n(){}}",
			"class a {\n	b() {}\n	c() {}\n}",
			"class a {\n	b() {}\n	c() {}\n}",
		},
		{ // 145
			"class\na{*b(){}\n*\nc\n(){}}",
			"class a {\n	* b() {}\n	* c() {}\n}",
			"class a {\n	* b() {}\n	* c() {}\n}",
		},
		{ // 146
			"class\na{async b(){}\nasync c\n(){}}",
			"class a {\n	async b() {}\n	async c() {}\n}",
			"class a {\n	async b() {}\n	async c() {}\n}",
		},
		{ // 147
			"class\na{async *b(){}\nasync *\nc\n(){}}",
			"class a {\n	async * b() {}\n	async * c() {}\n}",
			"class a {\n	async * b() {}\n	async * c() {}\n}",
		},
		{ // 148
			"class\na{get\nb(){}\nget c\n(){}}",
			"class a {\n	get b() {}\n	get c() {}\n}",
			"class a {\n	get b() {}\n	get c() {}\n}",
		},
		{ // 149
			"class\na{set\nb(c){}\nset d\n(e){}}",
			"class a {\n	set b(c) {}\n	set d(e) {}\n}",
			"class a {\n	set b(c) {}\n	set d(e) {}\n}",
		},
		{ // 150
			"class\na{static\nb(){}\nstatic c\n(){}}",
			"class a {\n	static b() {}\n	static c() {}\n}",
			"class a {\n	static b() {}\n	static c() {}\n}",
		},
		{ // 151
			"class\na{static\n*b(){}\nstatic *\nc\n(){}}",
			"class a {\n	static * b() {}\n	static * c() {}\n}",
			"class a {\n	static * b() {}\n	static * c() {}\n}",
		},
		{ // 152
			"class\na{static\nasync b(){}\nstatic async c\n(){}}",
			"class a {\n	static async b() {}\n	static async c() {}\n}",
			"class a {\n	static async b() {}\n	static async c() {}\n}",
		},
		{ // 153
			"class\na{static\nasync *b(){}\nstatic async *\nc(){}}",
			"class a {\n	static async * b() {}\n	static async * c() {}\n}",
			"class a {\n	static async * b() {}\n	static async * c() {}\n}",
		},
		{ // 154
			"class\na{static\nget\nb(){}static get c\n(){}}",
			"class a {\n	static get b() {}\n	static get c() {}\n}",
			"class a {\n	static get b() {}\n	static get c() {}\n}",
		},
		{ // 155
			"class\na{static\nset\nb(c){}static set d\n(e){}}",
			"class a {\n	static set b(c) {}\n	static set d(e) {}\n}",
			"class a {\n	static set b(c) {}\n	static set d(e) {}\n}",
		},
		{ // 156
			"a",
			"a;",
			"a;",
		},
		{ // 157
			"a?b:c",
			"a ? b : c;",
			"a ? b : c;",
		},
		{ // 158
			"a\n?\nb\n:\nc\n",
			"a ? b : c;",
			"a ? b : c;",
		},
		{ // 159
			"a=>b",
			"a => b;",
			"a => b;",
		},
		{ // 160
			"a =>\nb",
			"a => b;",
			"a => b;",
		},
		{ // 161
			"async a => b",
			"async a => b;",
			"async a => b;",
		},
		{ // 162
			"(a,b)=>b",
			"(a, b) => b;",
			"(a, b) => b;",
		},
		{ // 163
			"async (a,b)=>c",
			"async (a, b) => c;",
			"async (a, b) => c;",
		},
		{ // 164
			"a=>{}",
			"a => {};",
			"a => {};",
		},
		{ // 165
			"async a=>{}",
			"async a => {};",
			"async a => {};",
		},
		{ // 166
			"(a,b)=>{}",
			"(a, b) => {};",
			"(a, b) => {};",
		},
		{ // 167
			"async(a,b)=>{}",
			"async (a, b) => {};",
			"async (a, b) => {};",
		},
		{ // 168
			"new a",
			"new a;",
			"new a;",
		},
		{ // 169
			"new\nnew \n new	new\na",
			"new new new new a;",
			"new new new new a;",
		},
		{ // 170
			"super()",
			"super();",
			"super();",
		},
		{ // 171
			"super\n()",
			"super();",
			"super();",
		},
		{ // 172
			"import(a)",
			"import(a);",
			"import(a);",
		},
		{ // 173
			"import\n(a)",
			"import(a);",
			"import(a);",
		},
		{ // 174
			"a()",
			"a();",
			"a();",
		},
		{ // 175
			"a\n()",
			"a();",
			"a();",
		},
		{ // 176
			"a\n()\n()",
			"a()();",
			"a()();",
		},
		{ // 177
			"a()[b]",
			"a()[b];",
			"a()[b];",
		},
		{ // 178
			"a().b",
			"a().b;",
			"a().b;",
		},
		{ // 179
			"a()`b`",
			"a()`b`;",
			"a()`b`;",
		},
		{ // 180
			"var{a}=b",
			"var {a} = b;",
			"var {a: a} = b;",
		},
		{ // 181
			"var\n{\na\n:\nb\n}\n=\nc\n",
			"var {a: b} = c;",
			"var {a: b} = c;",
		},
		{ // 182
			"var[a=b]=c",
			"var [a = b] = c;",
			"var [a = b] = c;",
		},
		{ // 183
			"var[a=[b]] = c",
			"var [a = [b]] = c;",
			"var [a = [b]] = c;",
		},
		{ // 184
			"var[a={b}] = c",
			"var [a = {b}] = c;",
			"var [a = {b: b}] = c;",
		},
		{ // 185
			"var a={[\"b\"]:c}",
			"var a = {[\"b\"]: c};",
			"var a = {[\"b\"]: c};",
		},
		{ // 186
			"a||b",
			"a || b;",
			"a || b;",
		},
		{ // 187
			"(a,b,c)",
			"(a, b, c);",
			"(a, b, c);",
		},
		{ // 188
			"(\na\n,\nb\n,\nc\n)\n",
			"(a, b, c);",
			"(a, b, c);",
		},
		{ // 189
			"var a=(b,c,...d)=>{}",
			"var a = (b, c, ...d) => {};",
			"var a = (b, c, ...d) => {};",
		},
		{ // 190
			"var a=(b,c,...[e])=>{}",
			"var a = (b, c, ...[e]) => {};",
			"var a = (b, c, ...[e]) => {};",
		},
		{ // 191
			"var a=(b,c,...{...e})=>{}",
			"var a = (b, c, ...{...e}) => {};",
			"var a = (b, c, ...{...e}) => {};",
		},
		{ // 192
			"new a()",
			"new a();",
			"new a();",
		},
		{ // 193
			"new\nnew\na\n(\n)\n(\n)\n",
			"new new a()();",
			"new new a()();",
		},
		{ // 194
			"a\n[\n1\n]\n",
			"a[1];",
			"a[1];",
		},
		{ // 195
			"a\n.\nb\n",
			"a.b;",
			"a\n.b;",
		},
		{ // 196
			"a\n`b`",
			"a`b`;",
			"a`b`;",
		},
		{ // 197
			"new\nsuper\n[\na\n]\n[\nb\n]\n.\nc`d`\n(\nnew\n.\ntarget\n)\n",
			"new super[a][b].c`d`(new.target);",
			"new super[a][b]\n.c`d`(new.target);",
		},
		{ // 198
			"a(b,c,...d)",
			"a(b, c, ...d);",
			"a(b, c, ...d);",
		},
		{ // 199
			"a\n(\n...\nb\n)\n",
			"a(...b);",
			"a(...b);",
		},
		{ // 200
			"`a`",
			"`a`;",
			"`a`;",
		},
		{ // 201
			"`a${b}c`",
			"`a${b}c`;",
			"`a${b}c`;",
		},
		{ // 202
			"`a${\nb\n}c${\nd\n}e`",
			"`a${b}c${d}e`;",
			"`a${b}c${d}e`;",
		},
		{ // 203
			"{\n`a`\n}",
			"{\n\t`a`;\n}",
			"{\n\t`a`;\n}",
		},
		{ // 204
			"{\n`a\nb`\n}",
			"{\n\t`a\nb`;\n}",
			"{\n\t`a\nb`;\n}",
		},
		{ // 205
			"{\n`a\nb${c}d\ne${f}g\nh`\n}",
			"{\n\t`a\nb${c}d\ne${f}g\nh`;\n}",
			"{\n\t`a\nb${c}d\ne${f}g\nh`;\n}",
		},
		{ // 206
			"a&&b",
			"a && b;",
			"a && b;",
		},
		{ // 207
			"this",
			"this;",
			"this;",
		},
		{ // 208
			"a",
			"a;",
			"a;",
		},
		{ // 209
			"1",
			"1;",
			"1;",
		},
		{ // 210
			"[\n]\n",
			"[];",
			"[];",
		},
		{ // 211
			"var a={}",
			"var a = {};",
			"var a = {};",
		},
		{ // 212
			"var a=function(){}",
			"var a = function () {};",
			"var a = function () {};",
		},
		{ // 213
			"var a=class{}",
			"var a = class {};",
			"var a = class {};",
		},
		{ // 214
			"`a`",
			"`a`;",
			"`a`;",
		},
		{ // 215
			"(a)",
			"(a);",
			"(a);",
		},
		{ // 216
			"a|b",
			"a | b;",
			"a | b;",
		},
		{ // 217
			"[a,b,...c]",
			"[a, b, ...c];",
			"[a, b, ...c];",
		},
		{ // 218
			"[...a]",
			"[...a];",
			"[...a];",
		},
		{ // 219
			"[a]",
			"[a];",
			"[a];",
		},
		{ // 220
			"var a={b:c}",
			"var a = {b: c};",
			"var a = {b: c};",
		},
		{ // 221
			"var a={b:c,d:e}",
			"var a = {b: c, d: e};",
			"var a = {b: c, d: e};",
		},
		{ // 222
			"var a={\nb\n:\nc\n}",
			"var a = {b: c};",
			"var a = {\n	b: c\n};",
		},
		{ // 223
			"var a={\nb\n:\nc\n,\nd\n:\ne\n}",
			"var a = {b: c, d: e};",
			"var a = {\n	b: c,\n	d: e\n};",
		},
		{ // 224
			"a^b",
			"a ^ b;",
			"a ^ b;",
		},
		{ // 225
			"var a\n=\n{\nb\n:\nc\n,\nd\n,\ne\n=\nf\n,\ng\n(\n)\n{\n}\n,\n...\nh\n}\n",
			"var a = {b: c, d, e = f, g() {}, ...h};",
			"var a = {\n	b: c,\n	d: d,\n	e = f,\n	g() {},\n	...h\n};",
		},
		{ // 226
			"a&b",
			"a & b;",
			"a & b;",
		},
		{ // 227
			"a==b",
			"a == b;",
			"a == b;",
		},
		{ // 228
			"a!=b",
			"a != b;",
			"a != b;",
		},
		{ // 229
			"a===b",
			"a === b;",
			"a === b;",
		},
		{ // 230
			"a!==b",
			"a !== b;",
			"a !== b;",
		},
		{ // 231
			"a<b",
			"a < b;",
			"a < b;",
		},
		{ // 232
			"a>b",
			"a > b;",
			"a > b;",
		},
		{ // 233
			"a<=b",
			"a <= b;",
			"a <= b;",
		},
		{ // 234
			"a>=b",
			"a >= b;",
			"a >= b;",
		},
		{ // 235
			"a instanceof b",
			"a instanceof b;",
			"a instanceof b;",
		},
		{ // 236
			"a in b",
			"a in b;",
			"a in b;",
		},
		{ // 237
			"a<<b",
			"a << b;",
			"a << b;",
		},
		{ // 238
			"a>>b",
			"a >> b;",
			"a >> b;",
		},
		{ // 239
			"a>>>b",
			"a >>> b;",
			"a >>> b;",
		},
		{ // 240
			"a+b",
			"a + b;",
			"a + b;",
		},
		{ // 241
			"a-b",
			"a - b;",
			"a - b;",
		},
		{ // 242
			"a*b",
			"a * b;",
			"a * b;",
		},
		{ // 243
			"a/b",
			"a / b;",
			"a / b;",
		},
		{ // 244
			"a%b",
			"a % b;",
			"a % b;",
		},
		{ // 245
			"a**b",
			"a ** b;",
			"a ** b;",
		},
		{ // 246
			"delete a",
			"delete a;",
			"delete a;",
		},
		{ // 247
			"void a",
			"void a;",
			"void a;",
		},
		{ // 248
			"typeof a",
			"typeof a;",
			"typeof a;",
		},
		{ // 249
			"+\na",
			"+a;",
			"+a;",
		},
		{ // 250
			"-\na",
			"-a;",
			"-a;",
		},
		{ // 251
			"~\na",
			"~a;",
			"~a;",
		},
		{ // 252
			"!\na",
			"!a;",
			"!a;",
		},
		{ // 253
			"async function a(){await b}",
			"async function a() {\n	await b;\n}",
			"async function a() { await b; }",
		},
		{ // 254
			"a ++",
			"a++;",
			"a++;",
		},
		{ // 255
			"a --",
			"a--;",
			"a--;",
		},
		{ // 256
			"++\na",
			"++a;",
			"++a;",
		},
		{ // 257
			"--\na",
			"--a;",
			"--a;",
		},
		{ // 258
			"a: function b(){}",
			"a: function b() {}",
			"a: function b() {}",
		},
		{ // 259
			"a: b",
			"a: b;",
			"a: b;",
		},
		{ // 260
			"continue a",
			"continue a;",
			"continue a;",
		},
		{ // 261
			"debugger",
			"debugger;",
			"debugger;",
		},
		{ // 262
			"for(var a,b,\nc;;){}",
			"for (var a, b, c;;) {}",
			"for (var a, b,\n	c;;) {}",
		},
		{ // 263
			"for(var{a}in b){}",
			"for (var {a} in b) {}",
			"for (var {a: a} in b) {}",
		},
		{ // 264
			"for(var[a]in b){}",
			"for (var [a] in b) {}",
			"for (var [a] in b) {}",
		},
		{ // 265
			"switch(a){default:b}",
			"switch (a) {\ndefault:\n	b;\n}",
			"switch (a) {\ndefault:\n	b;\n}",
		},
		{ // 266
			"function*a(){yield *b}",
			"function* a() {\n	yield * b;\n}",
			"function* a() { yield * b; }",
		},
		{ // 267
			"a*=b",
			"a *= b;",
			"a *= b;",
		},
		{ // 268
			"var[[a]]=b",
			"var [[a]] = b;",
			"var [[a]] = b;",
		},
		{ // 269
			"var[{a}]=b",
			"var [{a}] = b;",
			"var [{a: a}] = b;",
		},
		{ // 270
			"super\n.\na\n",
			"super.a;",
			"super.a;",
		},
		{ // 271
			"a\n?.\nb",
			"a?.b;",
			"a?.b;",
		},
		{ // 272
			"a\n??\nb",
			"a ?? b;",
			"a ?? b;",
		},
		{ // 273
			"a\n??\nb\n??\nc",
			"a ?? b ?? c;",
			"a ?? b ?? c;",
		},
		{ // 274
			"a = ([b]) => b",
			"a = ([b]) => b;",
			"a = ([b]) => b;",
		},
		{ // 275
			"a?.b().c",
			"a?.b().c;",
			"a?.b().c;",
		},
		{ // 276
			"a?.b()?.c",
			"a?.b()?.c;",
			"a?.b()?.c;",
		},
		{ // 277
			"a&&=1",
			"a &&= 1;",
			"a &&= 1;",
		},
		{ // 278
			"a||=1",
			"a ||= 1;",
			"a ||= 1;",
		},
		{ // 279
			"a??=1",
			"a ??= 1;",
			"a ??= 1;",
		},
		{ // 280
			"[a, b] = [b, a]",
			"[a, b] = [b, a];",
			"[a, b] = [b, a];",
		},
		{ // 281
			"[a.b, a.c] = [a.c, a.b]",
			"[a.b, a.c] = [a.c, a.b];",
			"[a.b, a.c] = [a.c, a.b];",
		},
		{ // 282
			"{a}",
			"{\n	a;\n}",
			"{ a; }",
		},
		{ // 283
			"{a;b}",
			"{\n	a;\n	b;\n}",
			"{ a; b; }",
		},
		{ // 284
			"{a;\nb}",
			"{\n	a;\n	b;\n}",
			"{ a;\n	b; }",
		},
		{ // 285
			"{\na;\nb\n}",
			"{\n	a;\n	b;\n}",
			"{\n	a;\n	b;\n}",
		},
		{ // 286
			"({a, b} = {a: 1, b: 2})",
			"({a, b} = {a: 1, b: 2});",
			"({a: a, b: b} = {a: 1, b: 2});",
		},
		{ // 287
			"[a,b,...c] = [b, a]",
			"[a, b, ...c] = [b, a];",
			"[a, b, ...c] = [b, a];",
		},
		{ // 288
			"({a,b,...c}=d)",
			"({a, b, ...c} = d);",
			"({a: a, b: b, ...c} = d);",
		},
		{ // 289
			"({a:{b,c: d,...e},...f}=g)",
			"({a: {b, c: d, ...e}, ...f} = g);",
			"({a: {b: b, c: d, ...e}, ...f} = g);",
		},
		{ // 290
			"[a, ,[b,{c},,...d],,...e]=f",
			"[a, , [b, {c}, , ...d], , ...e] = f;",
			"[a, , [b, {c: c}, , ...d], , ...e] = f;",
		},
		{ // 291
			"a() ?.\nb",
			"a()?.b;",
			"a()?.b;",
		},
		{ // 292
			"a ?. [1]",
			"a?.[1];",
			"a?.[1];",
		},
		{ // 293
			"a ?. `1`",
			"a?.`1`;",
			"a?.`1`;",
		},
		{ // 294
			"a()\n.then()\n.catch()",
			"a().then().catch();",
			"a()\n.then()\n.catch();",
		},
		{ // 295
			"[a=b]=c",
			"[a = b] = c;",
			"[a = b] = c;",
		},
		{ // 296
			"({a=b}=c)",
			"({a = b} = c);",
			"({a = b} = c);",
		},
		{ // 297
			"a.#b",
			"a.#b;",
			"a.#b;",
		},
		{ // 298
			"a\n.#b",
			"a.#b;",
			"a\n.#b;",
		},
		{ // 299
			"a.#b.c",
			"a.#b.c;",
			"a.#b.c;",
		},
		{ // 300
			"a\n.#b\n.c",
			"a.#b.c;",
			"a\n.#b\n.c;",
		},
		{ // 301
			"class\na\n{\nb\n}",
			"class a {\n	b;\n}",
			"class a {\n	b;\n}",
		},
		{ // 302
			"class a { b () {} }",
			"class a {\n	b() {}\n}",
			"class a {\n	b() {}\n}",
		},
		{ // 303
			"class\na\n{\n#b\n}",
			"class a {\n	#b;\n}",
			"class a {\n	#b;\n}",
		},
		{ // 304
			"class a { #b () {} }",
			"class a {\n	#b() {}\n}",
			"class a {\n	#b() {}\n}",
		},
		{ // 305
			"class a { #b = 1 }",
			"class a {\n	#b = 1;\n}",
			"class a {\n	#b = 1;\n}",
		},
		{ // 306
			"class a { #b = 1; #c = 2 }",
			"class a {\n	#b = 1;\n	#c = 2;\n}",
			"class a {\n	#b = 1;\n	#c = 2;\n}",
		},
		{ // 307
			"class a { #b = 1\n#c = 2 }",
			"class a {\n	#b = 1;\n	#c = 2;\n}",
			"class a {\n	#b = 1;\n	#c = 2;\n}",
		},
		{ // 308
			"class a { #b(){}#c = 2 }",
			"class a {\n	#b() {}\n	#c = 2;\n}",
			"class a {\n	#b() {}\n	#c = 2;\n}",
		},
		{ // 309
			"class a { #b\n#c(){}}",
			"class a {\n	#b;\n	#c() {}\n}",
			"class a {\n	#b;\n	#c() {}\n}",
		},
		{ // 310
			"class a { #b = 1\n#c(){}}",
			"class a {\n	#b = 1;\n	#c() {}\n}",
			"class a {\n	#b = 1;\n	#c() {}\n}",
		},
		{ // 311
			"class a { #b = 1;#c(){}}",
			"class a {\n	#b = 1;\n	#c() {}\n}",
			"class a {\n	#b = 1;\n	#c() {}\n}",
		},
		{ // 312
			"class a { #b;#c(){}}",
			"class a {\n	#b;\n	#c() {}\n}",
			"class a {\n	#b;\n	#c() {}\n}",
		},
		{ // 313
			"class a {static a;static b\nstatic c = 2;static d(){};static{}static{e}static{e;f}}",
			"class a {\n\tstatic a;\n\tstatic b;\n\tstatic c = 2;\n\tstatic d() {}\n\tstatic {}\n\tstatic {\n\t\te;\n\t}\n\tstatic {\n\t\te;\n\t\tf;\n\t}\n}",
			"class a {\n\tstatic a;\n\tstatic b;\n\tstatic c = 2;\n\tstatic d() {}\n\tstatic {}\n\tstatic { e; }\n\tstatic { e; f; }\n}",
		},
		{ // 314
			"#a in b",
			"#a in b;",
			"#a in b;",
		},
		{ // 315
			"#a\nin\nb",
			"#a in b;",
			"#a in b;",
		},
		{ // 316
			"a().#b",
			"a().#b;",
			"a().#b;",
		},
		{ // 317
			"a\n(\n)\n.\n#b",
			"a().#b;",
			"a()\n.#b;",
		},
		{ // 318
			"a\n?.\n#b",
			"a?.#b;",
			"a?.#b;",
		},
		{ // 319
			"a?.c.#d",
			"a?.c.#d;",
			"a?.c.#d;",
		},
		{ // 320
			"// A\n// B\n\na();\n// C\n// D\n",
			"a();",
			"// A\n// B\n\na();\n// C\n// D\n",
		},
		{ // 321
			"/* A *//* B */\n// C\n\na();\n// D\n/* E */   /* F */\n",
			"a();",
			"/* A */ /* B */\n// C\n\na();\n// D\n/* E */ /* F */",
		},
		{ // 322
			"/*\nA\n*//* B */\n// C\n\na();\n// D\n/* E */   /*\n\nF\n\n*/\n",
			"a();",
			"/*\nA\n*/ /* B */\n// C\n\na();\n// D\n/* E */ /*\n\nF\n\n*/",
		},
		{ // 323
			"/*\nA\n*//* B */\n// C\n\n// D\na(); // E\n// F\n\n// G\n/* H */   /*\n\nI\n\n*/\n",
			"a();",
			"/*\nA\n*/ /* B */\n// C\n\n// D\na(); // E\n     // F\n\n// G\n/* H */ /*\n\nI\n\n*/",
		},
		{ // 324
			"// A\n\n// B\nsuper // C\n[ // D\n1\n // E\n]// F\n",
			"super[1];",
			"// A\n\n// B\nsuper // C\n[ // D\n\n\t1\n// E\n] // F\n",
		},
		{ // 325
			"// A\n\n// B\nsuper /* C */ . /* D */ a // E\n",
			"super.a;",
			"// A\n\n// B\nsuper /* C */ . /* D */ a // E\n",
		},
		{ // 326
			"// A\n\n// B\nnew /* C */./* D */target /* E */",
			"new.target;",
			"// A\n\n// B\nnew /* C */ . /* D */ target /* E */;",
		},
		{ // 327
			"// A\n\n/* B */import/* C */./* D */meta/* E */",
			"import.meta;",
			"// A\n\n/* B */ import /* C */ . /* D */ meta /* E */;",
		},
		{ // 328
			"// A\n\n// B\nnew/* C */1/* D */() // E\n",
			"new 1();",
			"// A\n\n// B\nnew /* C */ 1 /* D */ () // E\n",
		},
		{ // 329
			"// A\n\n// B\na // C\n",
			"a;",
			"// A\n\n// B\na // C\n",
		},
		{ // 330
			"// A\n\n// B\na /* C */``/* D */",
			"a``;",
			"// A\n\n// B\na /* C */ `` /* D */;",
		},
		{ // 331
			"// A\n\n/* B */a/* C */./* D */#b/* E */",
			"a.#b;",
			"// A\n\n/* B */ a /* C */ . /* D */ #b /* E */;",
		},
		{ // 332
			"// A\n\n/* B */a/* C */./* D */#b/* E */./* F */c // G\n",
			"a.#b.c;",
			"// A\n\n/* B */ a /* C */ . /* D */ #b /* E */ . /* F */ c // G\n",
		},
		{ // 333
			"// A\n\n// B\na /* C */ . /* D */ #b /* E */ [ // F\n\"c\"\n// G\n] /* H */",
			"a.#b[\"c\"];",
			"// A\n\n// B\na /* C */ . /* D */ #b /* E */ [ // F\n\n\t\"c\"\n// G\n] /* H */;",
		},
		{ // 334
			"super[ // C\n\n// D\n1 // E\n\n// F\n]",
			"super[1];",
			"super[ // C\n\n\t// D\n\t1 // E\n\n// F\n];",
		},
		{ // 335
			"// A\n\n// B\na /* C */ . /* D */ #b /* E */ [ // F\n\n// G\n\"c\" // H\n\n// I\n] /* J */",
			"a.#b[\"c\"];",
			"// A\n\n// B\na /* C */ . /* D */ #b /* E */ [ // F\n\n\t// G\n\t\"c\" // H\n\n// I\n] /* J */;",
		},
		{ // 336
			"a( // A\n\n// B\n)",
			"a();",
			"a( // A\n\n// B\n);",
		},
		{ // 337
			"a( // A\n\n// B\nb // C\n\n// D\n, // E\nc // F\n\n//G\n)",
			"a(b, c);",
			"a( // A\n\n\t// B\n\tb // C\n\n\t// D\n\t, // E\n\tc // F\n\n//G\n);",
		},
		{ // 338
			"a( // A\n\n// B\nb // C\n\n// D\n, // E\n...// F\nc // G\n\n//H\n)",
			"a(b, ...c);",
			"a( // A\n\n\t// B\n\tb // C\n\n\t// D\n\t, // E\n\t... // F\n\tc // G\n\n//H\n);",
		},
		{ // 339
			"( // A\n\n// B\na // C\n\n// D\n, // E\nb // F\n\n//G\n)",
			"(a, b);",
			"( // A\n\n\t// B\n\ta // C\n\n\t// D\n\t, // E\n\tb // F\n\n//G\n);",
		},
		{ // 340
			"[\n// A\n...// B\na // C\n]",
			"[...a];",
			"[\n\t// A\n\t... // B\n\ta // C\n];",
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
