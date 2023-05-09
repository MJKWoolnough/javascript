package minify

import (
	"fmt"
	"strings"
	"testing"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/parser"
)

func TestPrint(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output string
	}{
		{ // 1
			"var a = 1;",
			"var a=1",
		},
		{ // 2
			"var [a] = 1;",
			"var[a]=1",
		},
		{ // 3
			"async function a(){}",
			"async function a(){}",
		},
		{ // 4
			"typeof []",
			"typeof[]",
		},
		{ // 5
			"[] instanceof [].prototype",
			"[]instanceof[].prototype",
		},
		{ // 6
			"export * from 'a';",
			"export*from'a'",
		},
		{ // 7
			"export * as a from 'b';",
			"export*as a from'b'",
		},
		{ // 8
			"export {a, b as c, d} from 'f';",
			"export{a,b as c,d}from'f'",
		},
		{ // 9
			"import * as a from 'b';",
			"import*as a from'b'",
		},
		{ // 10
			"import {a, b as c, d} from 'e';",
			"import{a,b as c,d}from'e'",
		},
		{ // 11
			"import a from 'b';",
			"import a from'b'",
		},
		{ // 12
			"import a, {b, c} from 'e';",
			"import a,{b,c}from'e'",
		},
		{ // 13
			"a\nb\nc",
			"a;b;c",
		},
		{ // 14
			"a\n{}\nb",
			"a;{}b",
		},
		{ // 15
			"{a\nb\nc}",
			"{a;b;c}",
		},
		{ // 16
			"{a\n{}\nb}",
			"{a;{}b}",
		},
		{ // 17
			"if (a) b\nelse c",
			"if(a)b;else c",
		},
		{ // 18
			"if (a){\nb\n}\nelse{\nc\n}",
			"if(a){b}else{c}",
		},
		{ // 19
			"do a()\nwhile (1)",
			"do a();while(1)",
		},
		{ // 20
			"do{\na()\n}\nwhile (1)",
			"do{a()}while(1)",
		},
		{ // 21
			"switch(a){case a:\nb\ncase b:\n{}\ncase c: c}",
			"switch(a){case a:b;case b:{}case c:c}",
		},
		{ // 22
			"switch(a){case a:\nb\ndefault:\nc}",
			"switch(a){case a:b;default:c}",
		},
		{ // 23
			"switch(a){default:\na\ncase b: c}",
			"switch(a){default:a;case b:c}",
		},
		{ // 24
			"switch ( a ) { case []:\n1\ncase b: 2}",
			"switch(a){case[]:1;case b:2}",
		},
		{ // 25
			"switch ( a ) { case a:\na\nb\nc }",
			"switch(a){case a:a;b;c}",
		},
		{ // 26
			"class A {a\nb\nc\nd(){}\ne\n}",
			"class A{a;b;c;d(){}e}",
		},
		{ // 27
			"class A {static a = 1;static b(){}}",
			"class A{static a=1;static b(){}}",
		},
		{ // 28
			"class A {static [a] = 1;static [b](){}}",
			"class A{static[a]=1;static[b](){}}",
		},
		{ // 29
			"#a in b;",
			"#a in b",
		},
		{ // 30
			"#a in[b];",
			"#a in[b]",
		},
		{ // 31
			"import {a as b} from './c';",
			"import{a as b}from'./c'",
		},
		{ // 32
			"import * as a from './b';",
			"import*as a from'./b'",
		},
		{ // 33
			"var a = 1;",
			"var a=1",
		},
		{ // 34
			"var [a] = [1];",
			"var[a]=[1]",
		},
		{ // 35
			"function a(){}",
			"function a(){}",
		},
		{ // 36
			"(function (){})",
			"(function(){})",
		},
		{ // 37
			"async function a(){}",
			"async function a(){}",
		},
		{ // 38
			"class a extends (b) {}",
			"class a extends(b){}",
		},
		{ // 39
			"(class extends (a){})",
			"(class extends(a){})",
		},
		{ // 40
			"for (var a = 1;;){}",
			"for(var a=1;;){}",
		},
		{ // 41
			"for (var [a] = [1];;){}",
			"for(var[a]=[1];;){}",
		},
		{ // 42
			"new a();",
			"new a()",
		},
		{ // 43
			"new (a)();",
			"new(a)()",
		},
		{ // 44
			"new a;",
			"new a",
		},
		{ // 45
			"new (a);",
			"new(a)",
		},
		{ // 46
			"var a = 1, b = 2, c = 3;",
			"var a=1,b=2,c=3",
		},
		{ // 47
			"var [a, b, c, ...d] = [1, 2, 3, 4, 5], e = 4;",
			"var[a,b,c,...d]=[1,2,3,4,5],e=4",
		},
		{ // 48
			"var {a: b, c: c, d, ...e} = {a: 1, b: 2, c: 3, d: 4, e: 5, f: 6}, {g} = {g: 4}, h = 5;",
			"var{a:b,c,d,...e}={a:1,b:2,c:3,d:4,e:5,f:6},{g}={g:4},h=5",
		},
		{ // 49
			"let a = 1, b = 2, c = 3;",
			"let a=1,b=2,c=3",
		},
		{ // 50
			"let [a, b, c, ...d] = [1, 2, 3, 4, 5], e = 4;",
			"let[a,b,c,...d]=[1,2,3,4,5],e=4",
		},
		{ // 51
			"let {a: b, c: c, d, ...e} = {a: 1, b: 2, c: 3, d: 4, e: 5, f: 6}, {g} = {g: 4}, h = 5;",
			"let{a:b,c,d,...e}={a:1,b:2,c:3,d:4,e:5,f:6},{g}={g:4},h=5",
		},
		{ // 52
			"const a = 1, b = 2, c = 3;",
			"const a=1,b=2,c=3",
		},
		{ // 53
			"const [a, b, c, ...d] = [1, 2, 3, 4, 5], e = 4;",
			"const[a,b,c,...d]=[1,2,3,4,5],e=4",
		},
		{ // 54
			"const {a: b, c: c, d, ...e} = {a: 1, b: 2, c: 3, d: 4, e: 5, f: 6}, {g} = {g: 4}, h = 5;",
			"const{a:b,c,d,...e}={a:1,b:2,c:3,d:4,e:5,f:6},{g}={g:4},h=5",
		},
		{ // 55
			"do {aThing()} while (a == 1);",
			"do{aThing()}while(a==1)",
		},
		{ // 56
			"do aThing()\nwhile (a);",
			"do aThing();while(a)",
		},
		{ // 57
			"do [a] = next(); while(a);",
			"do[a]=next();while(a)",
		},
		{ // 58
			"while ( true ) run();",
			"while(true)run()",
		},
		{ // 59
			"while ( a = someThing()) {doAThing()}",
			"while(a=someThing()){doAThing()}",
		},
		{ // 60
			"while (a && b || c)[a]=runMe();",
			"while(a&&b||c)[a]=runMe()",
		},
		{ // 61
			"for (a = 1; b < 2; c++) {}",
			"for(a=1;b<2;c++){}",
		},
		{ // 62
			"for (a = 1, b = 2, [c] = [3]; b && c; c++) run();",
			"for(a=1,b=2,[c]=[3];b&&c;c++)run()",
		},
		{ // 63
			"for ( var a = 1, b = 2; b < 6; b++) {a(), b()}",
			"for(var a=1,b=2;b<6;b++){a(),b()}",
		},
		{ // 64
			"for ( var [a] = [1]; !a; a++) {a(); b()}",
			"for(var[a]=[1];!a;a++){a();b()}",
		},
		{ // 65
			"for ( let a = 1, b = 2; b < 6; b++) {a(), b()}",
			"for(let a=1,b=2;b<6;b++){a(),b()}",
		},
		{ // 66
			"for ( let [a] = [1]; !a; a++) {a(); b()}",
			"for(let[a]=[1];!a;a++){a();b()}",
		},
		{ // 67
			"for ( const a = 1, b = 2; b < 6; b++) {a(), b()}",
			"for(const a=1,b=2;b<6;b++){a(),b()}",
		},
		{ // 68
			"for ( const [a] = [1]; !a; a++) {a(); b()}",
			"for(const[a]=[1];!a;a++){a();b()}",
		},
		{ // 69
			"for ( a of b ){}",
			"for(a of b){}",
		},
		{ // 70
			"for ( [a, b] of c ){}",
			"for([a,b]of c){}",
		},
		{ // 71
			"for ( [a, b] of [c] ){}",
			"for([a,b]of[c]){}",
		},
		{ // 72
			"for ( a in b ){}",
			"for(a in b){}",
		},
		{ // 73
			"for ( {a: a, b} in c ){}",
			"for({a,b}in c){}",
		},
		{ // 74
			"for ( {a, b: d} of {c} ){}",
			"for({a,b:d}of{c}){}",
		},
		{ // 75
			"for ( var a of b ){}",
			"for(var a of b){}",
		},
		{ // 76
			"for ( var [a, b] of c ){}",
			"for(var[a,b]of c){}",
		},
		{ // 77
			"for ( var [a, b] of [c] ){}",
			"for(var[a,b]of[c]){}",
		},
		{ // 78
			"for ( var a in b ){}",
			"for(var a in b){}",
		},
		{ // 79
			"for ( var {a: a, b} in c ){}",
			"for(var{a,b}in c){}",
		},
		{ // 80
			"for ( var {a, b: d} of {c} ){}",
			"for(var{a,b:d}of{c}){}",
		},
		{ // 81
			"for ( let a of b ){}",
			"for(let a of b){}",
		},
		{ // 82
			"for ( let [a, b] of c ){}",
			"for(let[a,b]of c){}",
		},
		{ // 83
			"for ( let [a, b] of [c] ){}",
			"for(let[a,b]of[c]){}",
		},
		{ // 84
			"for ( let a in b ){}",
			"for(let a in b){}",
		},
		{ // 85
			"for ( let {a: a, b} in c ){}",
			"for(let{a,b}in c){}",
		},
		{ // 86
			"for ( let {a, b: d} of {c} ){}",
			"for(let{a,b:d}of{c}){}",
		},
		{ // 87
			"for ( const a of b ){}",
			"for(const a of b){}",
		},
		{ // 88
			"for ( const [a, b] of c ){}",
			"for(const[a,b]of c){}",
		},
		{ // 89
			"for ( const [a, b] of [c] ){}",
			"for(const[a,b]of[c]){}",
		},
		{ // 90
			"for ( const a in b ){}",
			"for(const a in b){}",
		},
		{ // 91
			"for ( const {a: a, b} in c ){}",
			"for(const{a,b}in c){}",
		},
		{ // 92
			"for ( const {a, b: d} of {c} ){}",
			"for(const{a,b:d}of{c}){}",
		},
		{ // 93
			"for await ( const a of b) {}",
			"for await(const a of b){}",
		},
		{ // 94
			"with ( a ) {}",
			"with(a){}",
		},
		{ // 95
			"with ( a ) b;",
			"with(a)b",
		},
		{ // 96
			"label: function a(){}",
			"label:function a(){}",
		},
		{ // 97
			"label: a++",
			"label:a++",
		},
		{ // 98
			"try { a(); b() } catch ( e ) {}",
			"try{a();b()}catch(e){}",
		},
		{ // 99
			"try { a(); } finally { something() }",
			"try{a()}finally{something()}",
		},
		{ // 100
			"try { a(); } catch ( e ) { e(); } finally { something(); }",
			"try{a()}catch(e){e()}finally{something()}",
		},
		{ // 101
			"continue;",
			"continue",
		},
		{ // 102
			"continue Label;",
			"continue Label",
		},
		{ // 103
			"break;",
			"break",
		},
		{ // 104
			"break Label;",
			"break Label",
		},
		{ // 105
			"() => {return;}",
			"()=>{return}",
		},
		{ // 106
			"() => {return a;}",
			"()=>{return a}",
		},
		{ // 107
			"() => {return [a];}",
			"()=>{return[a]}",
		},
		{ // 108
			"throw 1;",
			"throw 1",
		},
		{ // 109
			"throw [a];",
			"throw[a]",
		},
		{ // 110
			"debugger;",
			"debugger",
		},
		{ // 111
			"a, b, c;",
			"a,b,c",
		},
		{ // 112
			"function* a(){yield a;}",
			"function*a(){yield a}",
		},
		{ // 113
			"function* a(){yield [a];}",
			"function*a(){yield[a]}",
		},
		{ // 114
			"function* a(){yield * a;}",
			"function*a(){yield*a}",
		},
		{ // 115
			"function* a(){yield * [a];}",
			"function*a(){yield*[a]}",
		},
		{ // 116
			"(a) => b;",
			"(a)=>b",
		},
		{ // 117
			"(a, b) => c;",
			"(a,b)=>c",
		},
		{ // 118
			"(a, b, ...c) => d;",
			"(a,b,...c)=>d",
		},
		{ // 119
			"(a, b) =>{c;\nd;}",
			"(a,b)=>{c;d}",
		},
		{ // 120
			"(a, b) => c;",
			"(a,b)=>c",
		},
		{ // 121
			"a => b;",
			"a=>b",
		},
		{ // 122
			"a => {b;c;}",
			"a=>{b;c}",
		},
		{ // 123
			"async (a) => b;",
			"async(a)=>b",
		},
		{ // 124
			"async (a, b) => c;",
			"async(a,b)=>c",
		},
		{ // 125
			"async a => b;",
			"async a=>b",
		},
		{ // 126
			"a = 1",
			"a=1",
		},
		{ // 127
			"a *= 1",
			"a*=1",
		},
		{ // 128
			"a /= 1",
			"a/=1",
		},
		{ // 129
			"a %= 1",
			"a%=1",
		},
		{ // 130
			"a += 1",
			"a+=1",
		},
		{ // 131
			"a -= 1",
			"a-=1",
		},
		{ // 132
			"a <<= 1",
			"a<<=1",
		},
		{ // 133
			"a >>= 1",
			"a>>=1",
		},
		{ // 134
			"a >>>= 1",
			"a>>>=1",
		},
		{ // 135
			"a &= 1",
			"a&=1",
		},
		{ // 136
			"a ^= 1",
			"a^=1",
		},
		{ // 137
			"a |= 1",
			"a|=1",
		},
		{ // 138
			"a **= 1",
			"a**=1",
		},
		{ // 139
			"a &&= 1",
			"a&&=1",
		},
		{ // 140
			"a ||= 1",
			"a||=1",
		},
		{ // 141
			"a ??= 1",
			"a??=1",
		},
		{ // 142
			"[a] = b",
			"[a]=b",
		},
		{ // 143
			"[a, b] = c",
			"[a,b]=c",
		},
		{ // 144
			"[a, b, c] = d",
			"[a,b,c]=d",
		},
		{ // 145
			"[a, ...b] = c",
			"[a,...b]=c",
		},
		{ // 146
			"[...a] = b",
			"[...a]=b",
		},
		{ // 147
			"[a = b] = c",
			"[a=b]=c",
		},
		{ // 148
			"[a,b = c] = d",
			"[a,b=c]=d",
		},
		{ // 149
			"[a=b,c] = d",
			"[a=b,c]=d",
		},
		{ // 150
			"[a, b = c, d] = e",
			"[a,b=c,d]=e",
		},
		{ // 151
			"({a} = b)",
			"({a}=b)",
		},
		{ // 152
			"({a, b} = c)",
			"({a,b}=c)",
		},
		{ // 153
			"({a, b, c} = d)",
			"({a,b,c}=d)",
		},
		{ // 154
			"({a,...b} = c)",
			"({a,...b}=c)",
		},
		{ // 155
			"({a: b, c} = d)",
			"({a:b,c}=d)",
		},
		{ // 156
			"({a: b = c, d: e} = f)",
			"({a:b=c,d:e}=f)",
		},
		{ // 157
			"({a = b, c: d, e} = f)",
			"({a=b,c:d,e}=f)",
		},
		{ // 158
			"({a = b, c: d, e} = f)",
			"({a=b,c:d,e}=f)",
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)
		m, err := javascript.ParseModule(&tk)
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else {
			var sb strings.Builder
			if _, err := Print(&sb, m); err != nil {
				t.Errorf("test %d.1: unexpected error: %s", n+1, err)
			} else if str := sb.String(); str != test.Output {
				t.Errorf("test %d.1: expecting output %q, got %q", n+1, test.Output, str)
			} else {
				normalStr := fmt.Sprint(m)
				tk = parser.NewStringTokeniser(str)
				m, err := javascript.ParseModule(&tk)
				if err != nil {
					t.Errorf("test %d.2: unexpected error: %s", n+1, err)
				} else if otherStr := fmt.Sprint(m); normalStr != otherStr {
					t.Errorf("test %d.2: normal output not equal, expecting %s, got %s", n+1, normalStr, otherStr)
				}
			}
		}
	}
}
