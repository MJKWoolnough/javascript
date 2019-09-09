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
			"let {a} = 1;",
		},
		{ // 117
			"const\n{\na\n}\n=\n1",
			"const {a} = 1;",
			"const {a} = 1;",
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
			"var {a} = 1;",
		},
		{ // 136
			"var { a , b, ...c } = 1",
			"var {a, b, ...c} = 1;",
			"var {a, b, ...c} = 1;",
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
			"var {a} = b;",
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
			"var [a = {b}] = c;",
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
			"a.b;",
		},
		{ // 196
			"a\n`b`",
			"a`b`;",
			"a`b`;",
		},
		{ // 197
			"new\nsuper\n[\na\n]\n[\nb\n]\n.\nc`d`\n(\nnew\n.\ntarget\n)\n",
			"new super[a][b].c`d`(new.target);",
			"new super[a][b].c`d`(new.target);",
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
			"a&&b",
			"a && b;",
			"a && b;",
		},
		{ // 204
			"this",
			"this;",
			"this;",
		},
		{ // 205
			"a",
			"a;",
			"a;",
		},
		{ // 206
			"1",
			"1;",
			"1;",
		},
		{ // 207
			"[\n]\n",
			"[];",
			"[];",
		},
		{ // 208
			"var a={}",
			"var a = {};",
			"var a = {};",
		},
		{ // 209
			"var a=function(){}",
			"var a = function () {};",
			"var a = function () {};",
		},
		{ // 210
			"var a=class{}",
			"var a = class {};",
			"var a = class {};",
		},
		{ // 211
			"`a`",
			"`a`;",
			"`a`;",
		},
		{ // 212
			"(a)",
			"(a);",
			"(a);",
		},
		{ // 213
			"a|b",
			"a | b;",
			"a | b;",
		},
		{ // 214
			"[a,b,...c]",
			"[a, b, ...c];",
			"[a, b, ...c];",
		},
		{ // 215
			"[...a]",
			"[...a];",
			"[...a];",
		},
		{ // 216
			"[a]",
			"[a];",
			"[a];",
		},
		{ // 217
			"var a={b:c}",
			"var a = {b: c};",
			"var a = {b: c};",
		},
		{ // 218
			"var a={b:c,d:e}",
			"var a = {b: c, d: e};",
			"var a = {b: c, d: e};",
		},
		{ // 219
			"var a={\nb\n:\nc\n}",
			"var a = {b: c};",
			"var a = {\n	b: c\n};",
		},
		{ // 220
			"var a={\nb\n:\nc\n,\nd\n:\ne\n}",
			"var a = {b: c, d: e};",
			"var a = {\n	b: c,\n	d: e\n};",
		},
		{ // 221
			"a^b",
			"a ^ b;",
			"a ^ b;",
		},
		{ // 222
			"var a\n=\n{\nb\n:\nc\n,\nd\n,\ne\n=\nf\n,\ng\n(\n)\n{\n}\n,\n...\nh\n}\n",
			"var a = {b: c, d, e = f, g() {}, ...h};",
			"var a = {\n	b: c,\n	d,\n	e = f,\n	g() {},\n	...h\n};",
		},
		{ // 223
			"a&b",
			"a & b;",
			"a & b;",
		},
		{ // 224
			"a==b",
			"a == b;",
			"a == b;",
		},
		{ // 225
			"a!=b",
			"a != b;",
			"a != b;",
		},
		{ // 226
			"a===b",
			"a === b;",
			"a === b;",
		},
		{ // 227
			"a!==b",
			"a !== b;",
			"a !== b;",
		},
		{ // 228
			"a<b",
			"a < b;",
			"a < b;",
		},
		{ // 229
			"a>b",
			"a > b;",
			"a > b;",
		},
		{ // 230
			"a<=b",
			"a <= b;",
			"a <= b;",
		},
		{ // 231
			"a>=b",
			"a >= b;",
			"a >= b;",
		},
		{ // 232
			"a instanceof b",
			"a instanceof b;",
			"a instanceof b;",
		},
		{ // 233
			"a in b",
			"a in b;",
			"a in b;",
		},
		{ // 234
			"a<<b",
			"a << b;",
			"a << b;",
		},
		{ // 235
			"a>>b",
			"a >> b;",
			"a >> b;",
		},
		{ // 236
			"a>>>b",
			"a >>> b;",
			"a >>> b;",
		},
		{ // 237
			"a+b",
			"a + b;",
			"a + b;",
		},
		{ // 238
			"a-b",
			"a - b;",
			"a - b;",
		},
		{ // 239
			"a*b",
			"a * b;",
			"a * b;",
		},
		{ // 240
			"a/b",
			"a / b;",
			"a / b;",
		},
		{ // 241
			"a%b",
			"a % b;",
			"a % b;",
		},
		{ // 242
			"a**b",
			"a ** b;",
			"a ** b;",
		},
		{ // 243
			"delete a",
			"delete a;",
			"delete a;",
		},
		{ // 244
			"void a",
			"void a;",
			"void a;",
		},
		{ // 245
			"typeof a",
			"typeof a;",
			"typeof a;",
		},
		{ // 246
			"+\na",
			"+a;",
			"+a;",
		},
		{ // 247
			"-\na",
			"-a;",
			"-a;",
		},
		{ // 248
			"~\na",
			"~a;",
			"~a;",
		},
		{ // 249
			"!\na",
			"!a;",
			"!a;",
		},
		{ // 250
			"async function a(){await b}",
			"async function a() {\n	await b;\n}",
			"async function a() { await b; }",
		},
		{ // 251
			"a ++",
			"a++;",
			"a++;",
		},
		{ // 252
			"a --",
			"a--;",
			"a--;",
		},
		{ // 253
			"++\na",
			"++a;",
			"++a;",
		},
		{ // 254
			"--\na",
			"--a;",
			"--a;",
		},
		{ // 255
			"a: function b(){}",
			"a: function b() {}",
			"a: function b() {}",
		},
		{ // 256
			"a: b",
			"a: b;",
			"a: b;",
		},
		{ // 257
			"continue a",
			"continue a;",
			"continue a;",
		},
		{ // 258
			"debugger",
			"debugger;",
			"debugger;",
		},
		{ // 259
			"for(var a,b,\nc;;){}",
			"for (var a, b, c;;) {}",
			"for (var a, b,\n	c;;) {}",
		},
		{ // 260
			"for(var{a}in b){}",
			"for (var {a} in b) {}",
			"for (var {a} in b) {}",
		},
		{ // 261
			"for(var[a]in b){}",
			"for (var [a] in b) {}",
			"for (var [a] in b) {}",
		},
		{ // 262
			"switch(a){default:b}",
			"switch (a) {\ndefault:\n	b;\n}",
			"switch (a) {\ndefault:\n	b;\n}",
		},
		{ // 263
			"function*a(){yield *b}",
			"function* a() {\n	yield * b;\n}",
			"function* a() { yield * b; }",
		},
		{ // 264
			"a*=b",
			"a *= b;",
			"a *= b;",
		},
		{ // 265
			"var[[a]]=b",
			"var [[a]] = b;",
			"var [[a]] = b;",
		},
		{ // 266
			"var[{a}]=b",
			"var [{a}] = b;",
			"var [{a}] = b;",
		},
		{ // 272
			"super\n.\na\n",
			"super.a;",
			"super.a;",
		},
	} {
		for m, in := range [2]string{test.Input, test.VerboseOutput} {
			s, err := ParseScript(parser.NewStringTokeniser(in))
			if err != nil {
				t.Errorf("test %d.%d.1: unexpected error: %s", n+1, m+1, err)
				continue
			}
			st.Verbose = false
			st.Buffer = st.Buffer[:0]
			s.Format(&st, 's')
			if str := string(st.Buffer); str != test.SimpleOutput {
				t.Errorf("test %d.%d.2: expecting %q, got %q\n%s", n+1, m+1, test.SimpleOutput, str, s)
			}
			st.Verbose = true
			st.Buffer = st.Buffer[:0]
			s.Format(&st, 's')
			if str := string(st.Buffer); str != test.VerboseOutput {
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
	} {
		for m, in := range [2]string{test.Input, test.VerboseOutput} {
			s, err := ParseModule(parser.NewStringTokeniser(in))
			if err != nil {
				t.Errorf("test %d.%d.1: unexpected error: %s", n+1, m+1, err)
				continue
			}
			st.Verbose = false
			st.Buffer = st.Buffer[:0]
			s.Format(&st, 's')
			if str := string(st.Buffer); str != test.SimpleOutput {
				t.Errorf("test %d.%d.2: expecting %q, got %q\n%s", n+1, m+1, test.SimpleOutput, str, s)
			}
			st.Verbose = true
			st.Buffer = st.Buffer[:0]
			s.Format(&st, 's')
			if str := string(st.Buffer); str != test.VerboseOutput {
				t.Errorf("test %d.%d.3: expecting %q, got %q\n%s", n+1, m+1, test.VerboseOutput, str, s)
			}
		}
	}
}
