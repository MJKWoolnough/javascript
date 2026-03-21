package javascript

import (
	"testing"

	"vimagination.zapto.org/parser"
)

func TestTokeniser(t *testing.T) {
	for n, test := range [...]struct {
		Input   string
		Output  []parser.Token
		TS, JSX bool
	}{
		{ // 1
			Input: "",
			Output: []parser.Token{
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 2
			Input: " ",
			Output: []parser.Token{
				{Type: TokenWhitespace, Data: " "},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 3
			Input: " 	",
			Output: []parser.Token{
				{Type: TokenWhitespace, Data: " 	"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 4
			Input: "\t\v\f \xa0\ufeff",
			Output: []parser.Token{
				{Type: TokenWhitespace, Data: "\t\v\f \xa0\ufeff"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 5
			Input: " \n ",
			Output: []parser.Token{
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 6
			Input: "\n\r\u2028\u2029",
			Output: []parser.Token{
				{Type: TokenLineTerminator, Data: "\n\r\u2028\u2029"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 7
			Input: "\"\"",
			Output: []parser.Token{
				{Type: TokenStringLiteral, Data: "\"\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 8
			Input: "\"\\\"\"",
			Output: []parser.Token{
				{Type: TokenStringLiteral, Data: "\"\\\"\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 9
			Input: "\"\n\"",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "line terminator in string: \"\n"},
			},
		},
		{ // 10
			Input: "\"\\n\"",
			Output: []parser.Token{
				{Type: TokenStringLiteral, Data: "\"\\n\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 11
			Input: "\"\\0\"",
			Output: []parser.Token{
				{Type: TokenStringLiteral, Data: "\"\\0\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 12
			Input: "\"\\x20\"",
			Output: []parser.Token{
				{Type: TokenStringLiteral, Data: "\"\\x20\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 13
			Input: "\"\\u2020\"",
			Output: []parser.Token{
				{Type: TokenStringLiteral, Data: "\"\\u2020\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 14
			Input: "\"\\u\"",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid escape sequence: \"\\u\""},
			},
		},
		{ // 15
			Input: "\"\\up\"",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid escape sequence: \"\\up"},
			},
		},
		{ // 16
			Input: "\"\\u{20}\"",
			Output: []parser.Token{
				{Type: TokenStringLiteral, Data: "\"\\u{20}\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 17
			Input: "\"use strict\"",
			Output: []parser.Token{
				{Type: TokenStringLiteral, Data: "\"use strict\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 18
			Input: "\"use\\u{20}strict\\x65!\\0\"",
			Output: []parser.Token{
				{Type: TokenStringLiteral, Data: "\"use\\u{20}strict\\x65!\\0\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 19
			Input: "\"use strict\";",
			Output: []parser.Token{
				{Type: TokenStringLiteral, Data: "\"use strict\""},
				{Type: TokenPunctuator, Data: ";"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 20
			Input: "0",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 21
			Input: "0.1",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0.1"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 22
			Input: ".1",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: ".1"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 23
			Input: "0b0",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0b0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 24
			Input: "0b1",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0b1"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 25
			Input: "0b1001010101",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0b1001010101"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 26
			Input: "0",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 27
			Input: "1",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "1"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 28
			Input: "9",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "9"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 29
			Input: "12345678901",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "12345678901"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 30
			Input: "12345678.901",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "12345678.901"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 31
			Input: "12345678901E123",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "12345678901E123"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 32
			Input: "12345678901e+123",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "12345678901e+123"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 33
			Input: "12345678.901E-123",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "12345678.901E-123"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 34
			Input: "0x0",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0x0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 35
			Input: "0xa",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0xa"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 36
			Input: "0xf",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0xf"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 37
			Input: "0x0f",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0x0f"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 38
			Input: "0xaf",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0xaf"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 39
			Input: "0xDeAdBeEf",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0xDeAdBeEf"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 40
			Input: "0n",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0n"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 41
			Input: "1n",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "1n"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 42
			Input: "1234567890n",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "1234567890n"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 43
			Input: "0x1234567890n",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0x1234567890n"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 44
			Input: "Infinity",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "Infinity"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 45
			Input: "true",
			Output: []parser.Token{
				{Type: TokenBooleanLiteral, Data: "true"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 46
			Input: "false",
			Output: []parser.Token{
				{Type: TokenBooleanLiteral, Data: "false"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 47
			Input: "hello",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "hello"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 48
			Input: "this",
			Output: []parser.Token{
				{Type: TokenKeyword, Data: "this"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 49
			Input: "function",
			Output: []parser.Token{
				{Type: TokenKeyword, Data: "function"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 50
			Input: "/[a-z]+/g",
			Output: []parser.Token{
				{Type: TokenRegularExpressionLiteral, Data: "/[a-z]+/g"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 51
			Input: "/[\\n]/g",
			Output: []parser.Token{
				{Type: TokenRegularExpressionLiteral, Data: "/[\\n]/g"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 52
			Input: "var a =	 /^ab[cd]*$/ig;",
			Output: []parser.Token{
				{Type: TokenKeyword, Data: "var"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenWhitespace, Data: "	 "},
				{Type: TokenRegularExpressionLiteral, Data: "/^ab[cd]*$/ig"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 53
			Input: "num /= 4",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "num"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDivPunctuator, Data: "/="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumericLiteral, Data: "4"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 54
			Input: "const num = 8 / 4;",
			Output: []parser.Token{
				{Type: TokenKeyword, Data: "const"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "num"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumericLiteral, Data: "8"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDivPunctuator, Data: "/"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumericLiteral, Data: "4"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 55
			Input: "``",
			Output: []parser.Token{
				{Type: TokenNoSubstitutionTemplate, Data: "``"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 56
			Input: "`abc`",
			Output: []parser.Token{
				{Type: TokenNoSubstitutionTemplate, Data: "`abc`"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 57
			Input: "`ab${ (val.a / 2) + 1 }c${str}`",
			Output: []parser.Token{
				{Type: TokenTemplateHead, Data: "`ab${"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "("},
				{Type: TokenIdentifier, Data: "val"},
				{Type: TokenPunctuator, Data: "."},
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDivPunctuator, Data: "/"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumericLiteral, Data: "2"},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "+"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumericLiteral, Data: "1"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenTemplateMiddle, Data: "}c${"},
				{Type: TokenIdentifier, Data: "str"},
				{Type: TokenTemplateTail, Data: "}`"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 58
			Input: "const myFunc = function(aye, bee, cea) {\n	const num = [123, 4, lastNum(aye, \"beep\", () => window, val => val * 2, (myVar) => {myVar /= 2;return myVar;})], elm = document.getElementByID();\n	console.log(bee, num, elm);}",
			Output: []parser.Token{
				{Type: TokenKeyword, Data: "const"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "myFunc"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "function"},
				{Type: TokenPunctuator, Data: "("},
				{Type: TokenIdentifier, Data: "aye"},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "bee"},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "cea"},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "{"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenWhitespace, Data: "	"},
				{Type: TokenKeyword, Data: "const"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "num"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "["},
				{Type: TokenNumericLiteral, Data: "123"},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumericLiteral, Data: "4"},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "lastNum"},
				{Type: TokenPunctuator, Data: "("},
				{Type: TokenIdentifier, Data: "aye"},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenStringLiteral, Data: "\"beep\""},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "("},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "=>"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "window"},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "val"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "=>"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "val"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "*"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumericLiteral, Data: "2"},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "("},
				{Type: TokenIdentifier, Data: "myVar"},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "=>"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "{"},
				{Type: TokenIdentifier, Data: "myVar"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDivPunctuator, Data: "/="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumericLiteral, Data: "2"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenKeyword, Data: "return"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "myVar"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenRightBracePunctuator, Data: "}"},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenPunctuator, Data: "]"},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "elm"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "document"},
				{Type: TokenPunctuator, Data: "."},
				{Type: TokenIdentifier, Data: "getElementByID"},
				{Type: TokenPunctuator, Data: "("},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenWhitespace, Data: "	"},
				{Type: TokenIdentifier, Data: "console"},
				{Type: TokenPunctuator, Data: "."},
				{Type: TokenIdentifier, Data: "log"},
				{Type: TokenPunctuator, Data: "("},
				{Type: TokenIdentifier, Data: "bee"},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "num"},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "elm"},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: TokenRightBracePunctuator, Data: "}"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 59
			Input: "export {name1, name2};",
			Output: []parser.Token{
				{Type: TokenKeyword, Data: "export"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "{"},
				{Type: TokenIdentifier, Data: "name1"},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "name2"},
				{Type: TokenRightBracePunctuator, Data: "}"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 60
			Input: "export {var1 as name1, var2 as name2};",
			Output: []parser.Token{
				{Type: TokenKeyword, Data: "export"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "{"},
				{Type: TokenIdentifier, Data: "var1"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "as"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "name1"},
				{Type: TokenPunctuator, Data: ","},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "var2"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "as"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "name2"},
				{Type: TokenRightBracePunctuator, Data: "}"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 61
			Input: "export * from './other.js';",
			Output: []parser.Token{
				{Type: TokenKeyword, Data: "export"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "*"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "from"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenStringLiteral, Data: "'./other.js'"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 62
			Input: "import * as name from './module.js';",
			Output: []parser.Token{
				{Type: TokenKeyword, Data: "import"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "*"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "as"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "name"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "from"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenStringLiteral, Data: "'./module.js'"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 63
			Input: "$",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "$"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 64
			Input: "_",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "_"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 65
			Input: "\\u0061",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "\\u0061"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 66
			Input: "// Comment",
			Output: []parser.Token{
				{Type: TokenSingleLineComment, Data: "// Comment"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 67
			Input: "enum",
			Output: []parser.Token{
				{Type: TokenFutureReservedWord, Data: "enum"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 68
			Input: ".01234E56",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: ".01234E56"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 69
			Input: "0.01234E56",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0.01234E56"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 70
			Input: "0o1234567",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0o1234567"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 71
			Input: "`\\x60`",
			Output: []parser.Token{
				{Type: TokenNoSubstitutionTemplate, Data: "`\\x60`"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 72
			Input: "/\\(/",
			Output: []parser.Token{
				{Type: TokenRegularExpressionLiteral, Data: "/\\(/"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 73
			Input: "/a\\(/",
			Output: []parser.Token{
				{Type: TokenRegularExpressionLiteral, Data: "/a\\(/"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 74
			Input: "{",
			Output: []parser.Token{
				{Type: TokenPunctuator, Data: "{"},
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 75
			Input: "[",
			Output: []parser.Token{
				{Type: TokenPunctuator, Data: "["},
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 76
			Input: "(",
			Output: []parser.Token{
				{Type: TokenPunctuator, Data: "("},
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 77
			Input: "/*",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 78
			Input: "[}",
			Output: []parser.Token{
				{Type: TokenPunctuator, Data: "["},
				{Type: parser.TokenError, Data: "invalid character: }"},
			},
		},
		{ // 79
			Input: "(}",
			Output: []parser.Token{
				{Type: TokenPunctuator, Data: "("},
				{Type: parser.TokenError, Data: "invalid character: }"},
			},
		},
		{ // 80
			Input: "{)",
			Output: []parser.Token{
				{Type: TokenPunctuator, Data: "{"},
				{Type: parser.TokenError, Data: "invalid character: )"},
			},
		},
		{ // 81
			Input: "{]",
			Output: []parser.Token{
				{Type: TokenPunctuator, Data: "{"},
				{Type: parser.TokenError, Data: "invalid character: ]"},
			},
		},
		{ // 82
			Input: "(]",
			Output: []parser.Token{
				{Type: TokenPunctuator, Data: "("},
				{Type: parser.TokenError, Data: "invalid character: ]"},
			},
		},
		{ // 83
			Input: "[)",
			Output: []parser.Token{
				{Type: TokenPunctuator, Data: "["},
				{Type: parser.TokenError, Data: "invalid character: )"},
			},
		},
		{ // 84
			Input: "..",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 85
			Input: "..a",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid character sequence: ..a"},
			},
		},
		{ // 86
			Input: "/\\\n/",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid regexp sequence: /\\\n"},
			},
		},
		{ // 87
			Input: "/[",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 88
			Input: "/[\\",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 89
			Input: "/",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 90
			Input: "/\n",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid regexp character: \n"},
			},
		},
		{ // 91
			Input: "/a",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 92
			Input: "/a\\\n/",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid regexp sequence: /a\\\n"},
			},
		},
		{ // 93
			Input: "/a[",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 94
			Input: "/a\n",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid regexp character: \n"},
			},
		},
		{ // 95
			Input: "0B9",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 0B9"},
			},
		},
		{ // 96
			Input: "0O9",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 0O9"},
			},
		},
		{ // 97
			Input: "0XG",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 0XG"},
			},
		},
		{ // 98
			Input: "\\x60",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "unexpected backslash: \\x"},
			},
		},
		{ // 99
			Input: "\\ug",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid unicode escape sequence: \\ug"},
			},
		},
		{ // 100
			Input: "\\u{G}",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid unicode escape sequence: \\u{G"},
			},
		},
		{ // 101
			Input: "\\u{ffff",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid unicode escape sequence: \\u{ffff"},
			},
		},
		{ // 102
			Input: "}",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid character: }"},
			},
		},
		{ // 103
			Input: "`\\G`",
			Output: []parser.Token{
				{Type: TokenNoSubstitutionTemplate, Data: "`\\G`"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 104
			Input: "`",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 105
			Input: "1_234_567",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "1_234_567"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 106
			Input: "1_",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 1_"},
			},
		},
		{ // 107
			Input: "1__234_567",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 1__"},
			},
		},
		{ // 108
			Input: "123e-456_789",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "123e-456_789"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 109
			Input: "0.123_456",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0.123_456"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 110
			Input: "1.2_3_4_5_6_7",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "1.2_3_4_5_6_7"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 111
			Input: "0x1_2",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0x1_2"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 112
			Input: "0b1_0",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0b1_0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 113
			Input: "0o1_7",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: "0o1_7"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 114
			Input: "a.b",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenPunctuator, Data: "."},
				{Type: TokenIdentifier, Data: "b"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 115
			Input: "a?.b",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenPunctuator, Data: "?."},
				{Type: TokenIdentifier, Data: "b"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 116
			Input: "a??b",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenPunctuator, Data: "??"},
				{Type: TokenIdentifier, Data: "b"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 117
			Input: "0.",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 0."},
			},
		},
		{ // 118
			Input: "0.1e",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 0.1e"},
			},
		},
		{ // 119
			Input: "1.",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 1."},
			},
		},
		{ // 120
			Input: "1.1e",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 1.1e"},
			},
		},
		{ // 121
			Input: "import(a)",
			Output: []parser.Token{
				{Type: TokenKeyword, Data: "import"},
				{Type: TokenPunctuator, Data: "("},
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenPunctuator, Data: ")"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 122
			Input: "\\u0060",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid unicode escape sequence: \\u0060"},
			},
		},
		{ // 123
			Input: "\\u0024",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "\\u0024"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 124
			Input: "\\u{5f}",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "\\u{5f}"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 125
			Input: "\\u{41}",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "\\u{41}"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 126
			Input: "\\u{0}",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid unicode escape sequence: \\u{0}"},
			},
		},
		{ // 127
			Input: "\\u005C",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid unicode escape sequence: \\u005C"},
			},
		},
		{ // 128
			Input: "/a/g",
			Output: []parser.Token{
				{Type: TokenRegularExpressionLiteral, Data: "/a/g"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 129
			Input: "/a/\\u000A",
			Output: []parser.Token{
				{Type: TokenRegularExpressionLiteral, Data: "/a/"},
				{Type: parser.TokenError, Data: "invalid unicode escape sequence: \\u000A"},
			},
		},
		{ // 130
			Input: "a`b${f}c`",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenTemplateHead, Data: "`b${"},
				{Type: TokenIdentifier, Data: "f"},
				{Type: TokenTemplateTail, Data: "}c`"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 131
			Input: "#a",
			Output: []parser.Token{
				{Type: TokenPrivateIdentifier, Data: "#a"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 132
			Input: "a.#b",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenPunctuator, Data: "."},
				{Type: TokenPrivateIdentifier, Data: "#b"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 133
			Input: "#",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid character sequence: #"},
			},
		},
		{ // 134
			Input: "#.a",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid character sequence: #."},
			},
		},
		{ // 135
			Input: "Number(10000n * this.#numerator / this.#denominator) / 10000;",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "Number"},
				{Type: TokenPunctuator, Data: "("},
				{Type: TokenNumericLiteral, Data: "10000n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPunctuator, Data: "*"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "this"},
				{Type: TokenPunctuator, Data: "."},
				{Type: TokenPrivateIdentifier, Data: "#numerator"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDivPunctuator, Data: "/"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenKeyword, Data: "this"},
				{Type: TokenPunctuator, Data: "."},
				{Type: TokenPrivateIdentifier, Data: "#denominator"},
				{Type: TokenPunctuator, Data: ")"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDivPunctuator, Data: "/"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumericLiteral, Data: "10000"},
				{Type: TokenPunctuator, Data: ";"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 136
			Input: "`\\x`",
			Output: []parser.Token{
				{Type: parser.TokenError, Data: "invalid escape sequence: `\\x`"},
			},
		},
		{ // 137
			Input: ".123_456",
			Output: []parser.Token{
				{Type: TokenNumericLiteral, Data: ".123_456"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 138
			Input: "`a\\`b`",
			Output: []parser.Token{
				{Type: TokenNoSubstitutionTemplate, Data: "`a\\`b`"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 139
			Input: "<a></a>",
			Output: []parser.Token{
				{Type: TokenJSXElementStart, Data: "<"},
				{Type: TokenJSXIdentifier, Data: "a"},
				{Type: TokenJSXElementEnd, Data: ">"},
				{Type: TokenJSXElementStart, Data: "<"},
				{Type: TokenPunctuator, Data: "/"},
				{Type: TokenJSXIdentifier, Data: "a"},
				{Type: TokenJSXElementEnd, Data: ">"},
				{Type: parser.TokenDone, Data: ""},
			},
			JSX: true,
		},
		{ // 140
			Input: "<></>",
			Output: []parser.Token{
				{Type: TokenJSXElementStart, Data: "<"},
				{Type: TokenJSXElementEnd, Data: ">"},
				{Type: TokenJSXElementStart, Data: "<"},
				{Type: TokenPunctuator, Data: "/"},
				{Type: TokenJSXElementEnd, Data: ">"},
				{Type: parser.TokenDone, Data: ""},
			},
			JSX: true,
		},
		{ // 141
			Input: "<a/>",
			Output: []parser.Token{
				{Type: TokenJSXElementStart, Data: "<"},
				{Type: TokenJSXIdentifier, Data: "a"},
				{Type: TokenPunctuator, Data: "/"},
				{Type: TokenJSXElementEnd, Data: ">"},
				{Type: parser.TokenDone, Data: ""},
			},
			JSX: true,
		},
		{ // 142
			Input: "a=<b c=\"d\"></b>",
			Output: []parser.Token{
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenJSXElementStart, Data: "<"},
				{Type: TokenJSXIdentifier, Data: "b"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenJSXIdentifier, Data: "c"},
				{Type: TokenPunctuator, Data: "="},
				{Type: TokenJSXString, Data: "\"d\""},
				{Type: TokenJSXElementEnd, Data: ">"},
				{Type: TokenJSXElementStart, Data: "<"},
				{Type: TokenPunctuator, Data: "/"},
				{Type: TokenJSXIdentifier, Data: "b"},
				{Type: TokenJSXElementEnd, Data: ">"},
				{Type: parser.TokenDone, Data: ""},
			},
			JSX: true,
		},
		{ // 143
			Input: "<a/> / b",
			Output: []parser.Token{
				{Type: TokenJSXElementStart, Data: "<"},
				{Type: TokenJSXIdentifier, Data: "a"},
				{Type: TokenPunctuator, Data: "/"},
				{Type: TokenJSXElementEnd, Data: ">"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDivPunctuator, Data: "/"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdentifier, Data: "b"},
				{Type: parser.TokenDone, Data: ""},
			},
			JSX: true,
		},
	} {
		p := parser.NewStringTokeniser(test.Input)

		var tks Tokeniser = SetTokeniser(&p)

		if test.TS {
			tks = AsTypescript(tks)
		}

		if test.JSX {
			tks = AsJSX(tks)
		}

		for m, tkn := range test.Output {
			if tk, _ := tks.GetToken(); tk.Type != tkn.Type {
				if tk.Type == parser.TokenError {
					t.Errorf("test %d.%d: unexpected error: %s", n+1, m+1, tk.Data)
				} else {
					t.Errorf("test %d.%d: Incorrect type, expecting %d, got %d", n+1, m+1, tkn.Type, tk.Type)
				}

				break
			} else if tk.Data != tkn.Data {
				t.Errorf("test %d.%d: Incorrect data, expecting %q, got %q", n+1, m+1, tkn.Data, tk.Data)

				break
			}
		}
	}
}
