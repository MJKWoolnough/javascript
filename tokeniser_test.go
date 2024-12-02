package javascript

import (
	"testing"

	"vimagination.zapto.org/parser"
)

func TestTokeniser(t *testing.T) {
	for n, test := range [...]struct {
		Input  string
		Output []parser.Token
	}{
		{ // 1
			"",
			[]parser.Token{
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 2
			" ",
			[]parser.Token{
				{Type: TokenWhitespace, Data: " "},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 3
			" 	",
			[]parser.Token{
				{Type: TokenWhitespace, Data: " 	"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 4
			"\t\v\f \xa0\ufeff",
			[]parser.Token{
				{Type: TokenWhitespace, Data: "\t\v\f \xa0\ufeff"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 5
			" \n ",
			[]parser.Token{
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenLineTerminator, Data: "\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 6
			"\"\"",
			[]parser.Token{
				{Type: TokenStringLiteral, Data: "\"\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 7
			"\"\\\"\"",
			[]parser.Token{
				{Type: TokenStringLiteral, Data: "\"\\\"\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 8
			"\"\n\"",
			[]parser.Token{
				{Type: parser.TokenError, Data: "line terminator in string: \"\n"},
			},
		},
		{ // 9
			"\"\\n\"",
			[]parser.Token{
				{Type: TokenStringLiteral, Data: "\"\\n\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 10
			"\"\\0\"",
			[]parser.Token{
				{Type: TokenStringLiteral, Data: "\"\\0\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 11
			"\"\\x20\"",
			[]parser.Token{
				{Type: TokenStringLiteral, Data: "\"\\x20\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 12
			"\"\\u2020\"",
			[]parser.Token{
				{Type: TokenStringLiteral, Data: "\"\\u2020\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 13
			"\"\\u\"",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid escape sequence: \"\\u\""},
			},
		},
		{ // 14
			"\"\\up\"",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid escape sequence: \"\\up"},
			},
		},
		{ // 15
			"\"\\u{20}\"",
			[]parser.Token{
				{Type: TokenStringLiteral, Data: "\"\\u{20}\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 16
			"\"use strict\"",
			[]parser.Token{
				{Type: TokenStringLiteral, Data: "\"use strict\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 17
			"\"use\\u{20}strict\\x65!\\0\"",
			[]parser.Token{
				{Type: TokenStringLiteral, Data: "\"use\\u{20}strict\\x65!\\0\""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 18
			"\"use strict\";",
			[]parser.Token{
				{Type: TokenStringLiteral, Data: "\"use strict\""},
				{Type: TokenPunctuator, Data: ";"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 19
			"0",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 20
			"0.1",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0.1"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 21
			".1",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: ".1"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 22
			"0b0",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0b0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 23
			"0b1",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0b1"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 24
			"0b1001010101",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0b1001010101"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 25
			"0",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 26
			"1",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "1"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 27
			"9",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "9"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 28
			"12345678901",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "12345678901"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 29
			"12345678.901",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "12345678.901"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 30
			"12345678901E123",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "12345678901E123"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 31
			"12345678901e+123",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "12345678901e+123"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 32
			"12345678.901E-123",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "12345678.901E-123"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 33
			"0x0",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0x0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 34
			"0xa",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0xa"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 35
			"0xf",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0xf"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 36
			"0x0f",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0x0f"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 37
			"0xaf",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0xaf"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 38
			"0xDeAdBeEf",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0xDeAdBeEf"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 39
			"0n",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0n"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 40
			"1n",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "1n"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 41
			"1234567890n",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "1234567890n"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 42
			"0x1234567890n",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0x1234567890n"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 43
			"Infinity",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "Infinity"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 44
			"true",
			[]parser.Token{
				{Type: TokenBooleanLiteral, Data: "true"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 45
			"false",
			[]parser.Token{
				{Type: TokenBooleanLiteral, Data: "false"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 46
			"hello",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "hello"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 47
			"this",
			[]parser.Token{
				{Type: TokenKeyword, Data: "this"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 48
			"function",
			[]parser.Token{
				{Type: TokenKeyword, Data: "function"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 49
			"/[a-z]+/g",
			[]parser.Token{
				{Type: TokenRegularExpressionLiteral, Data: "/[a-z]+/g"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 50
			"/[\\n]/g",
			[]parser.Token{
				{Type: TokenRegularExpressionLiteral, Data: "/[\\n]/g"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 51
			"var a =	 /^ab[cd]*$/ig;",
			[]parser.Token{
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
		{ // 52
			"num /= 4",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "num"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDivPunctuator, Data: "/="},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumericLiteral, Data: "4"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 53
			"const num = 8 / 4;",
			[]parser.Token{
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
		{ // 54
			"``",
			[]parser.Token{
				{Type: TokenNoSubstitutionTemplate, Data: "``"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 55
			"`abc`",
			[]parser.Token{
				{Type: TokenNoSubstitutionTemplate, Data: "`abc`"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 56
			"`ab${ (val.a / 2) + 1 }c${str}`",
			[]parser.Token{
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
		{ // 57
			"const myFunc = function(aye, bee, cea) {\n	const num = [123, 4, lastNum(aye, \"beep\", () => window, val => val * 2, (myVar) => {myVar /= 2;return myVar;})], elm = document.getElementByID();\n	console.log(bee, num, elm);}",
			[]parser.Token{
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
		{ // 58
			"export {name1, name2};",
			[]parser.Token{
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
		{ // 59
			"export {var1 as name1, var2 as name2};",
			[]parser.Token{
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
		{ // 60
			"export * from './other.js';",
			[]parser.Token{
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
		{ // 61
			"import * as name from './module.js';",
			[]parser.Token{
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
		{ // 62
			"$",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "$"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 63
			"_",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "_"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 64
			"\\u0061",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "\\u0061"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 65
			"// Comment",
			[]parser.Token{
				{Type: TokenSingleLineComment, Data: "// Comment"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 66
			"enum",
			[]parser.Token{
				{Type: TokenFutureReservedWord, Data: "enum"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 67
			".01234E56",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: ".01234E56"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 68
			"0.01234E56",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0.01234E56"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 69
			"0o1234567",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0o1234567"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 70
			"`\\x60`",
			[]parser.Token{
				{Type: TokenNoSubstitutionTemplate, Data: "`\\x60`"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 71
			"/\\(/",
			[]parser.Token{
				{Type: TokenRegularExpressionLiteral, Data: "/\\(/"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 72
			"/a\\(/",
			[]parser.Token{
				{Type: TokenRegularExpressionLiteral, Data: "/a\\(/"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 73
			"{",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "{"},
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 74
			"[",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "["},
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 75
			"(",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "("},
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 76
			"/*",
			[]parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 77
			"[}",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "["},
				{Type: parser.TokenError, Data: "invalid character: }"},
			},
		},
		{ // 78
			"(}",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "("},
				{Type: parser.TokenError, Data: "invalid character: }"},
			},
		},
		{ // 79
			"{)",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "{"},
				{Type: parser.TokenError, Data: "invalid character: )"},
			},
		},
		{ // 80
			"{]",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "{"},
				{Type: parser.TokenError, Data: "invalid character: ]"},
			},
		},
		{ // 81
			"(]",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "("},
				{Type: parser.TokenError, Data: "invalid character: ]"},
			},
		},
		{ // 82
			"[)",
			[]parser.Token{
				{Type: TokenPunctuator, Data: "["},
				{Type: parser.TokenError, Data: "invalid character: )"},
			},
		},
		{ // 83
			"..",
			[]parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 84
			"..a",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid character sequence: ..a"},
			},
		},
		{ // 85
			"/\\\n/",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid regexp sequence: /\\\n"},
			},
		},
		{ // 86
			"/[",
			[]parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 87
			"/[\\",
			[]parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 88
			"/",
			[]parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 89
			"/\n",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid regexp character: \n"},
			},
		},
		{ // 90
			"/a",
			[]parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 91
			"/a\\\n/",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid regexp sequence: /a\\\n"},
			},
		},
		{ // 92
			"/a[",
			[]parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 93
			"/a\n",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid regexp character: \n"},
			},
		},
		{ // 94
			"0B9",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 0B9"},
			},
		},
		{ // 95
			"0O9",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 0O9"},
			},
		},
		{ // 96
			"0XG",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 0XG"},
			},
		},
		{ // 97
			"\\x60",
			[]parser.Token{
				{Type: parser.TokenError, Data: "unexpected backslash: \\x"},
			},
		},
		{ // 98
			"\\ug",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid unicode escape sequence: \\ug"},
			},
		},
		{ // 99
			"\\u{G}",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid unicode escape sequence: \\u{G"},
			},
		},
		{ // 100
			"\\u{ffff",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid unicode escape sequence: \\u{ffff"},
			},
		},
		{ // 101
			"}",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid character: }"},
			},
		},
		{ // 102
			"`\\G`",
			[]parser.Token{
				{Type: TokenNoSubstitutionTemplate, Data: "`\\G`"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 103
			"`",
			[]parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 104
			"1_234_567",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "1_234_567"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 105
			"1_",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 1_"},
			},
		},
		{ // 106
			"1__234_567",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 1__"},
			},
		},
		{ // 107
			"123e-456_789",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "123e-456_789"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 108
			"0.123_456",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0.123_456"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 109
			"1.2_3_4_5_6_7",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "1.2_3_4_5_6_7"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 110
			"0x1_2",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0x1_2"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 111
			"0b1_0",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0b1_0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 112
			"0o1_7",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: "0o1_7"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 113
			"a.b",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenPunctuator, Data: "."},
				{Type: TokenIdentifier, Data: "b"},
			},
		},
		{ // 114
			"a?.b",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenPunctuator, Data: "?."},
				{Type: TokenIdentifier, Data: "b"},
			},
		},
		{ // 115
			"a??b",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenPunctuator, Data: "??"},
				{Type: TokenIdentifier, Data: "b"},
			},
		},
		{ // 116
			"0.",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 0."},
			},
		},
		{ // 117
			"0.1e",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 0.1e"},
			},
		},
		{ // 118
			"1.",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 1."},
			},
		},
		{ // 119
			"1.1e",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid number: 1.1e"},
			},
		},
		{ // 120
			"import(a)",
			[]parser.Token{
				{Type: TokenKeyword, Data: "import"},
				{Type: TokenPunctuator, Data: "("},
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenPunctuator, Data: ")"},
			},
		},
		{ // 121
			"\\u0060",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid unicode escape sequence: \\u0060"},
			},
		},
		{ // 122
			"\\u0024",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "\\u0024"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 123
			"\\u{5f}",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "\\u{5f}"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 124
			"\\u{41}",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "\\u{41}"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{ // 125
			"\\u{0}",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid unicode escape sequence: \\u{0}"},
			},
		},
		{ // 126
			"\\u005C",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid unicode escape sequence: \\u005C"},
			},
		},
		{ // 127
			"/a/g",
			[]parser.Token{
				{Type: TokenRegularExpressionLiteral, Data: "/a/g"},
			},
		},
		{ // 128
			"/a/\\u000A",
			[]parser.Token{
				{Type: TokenRegularExpressionLiteral, Data: "/a/"},
			},
		},
		{ // 129
			"a`b${f}c`",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenTemplateHead, Data: "`b${"},
				{Type: TokenIdentifier, Data: "f"},
				{Type: TokenTemplateTail, Data: "}c`"},
			},
		},
		{ // 130
			"#a",
			[]parser.Token{
				{Type: TokenPrivateIdentifier, Data: "#a"},
			},
		},
		{ // 131
			"a.#b",
			[]parser.Token{
				{Type: TokenIdentifier, Data: "a"},
				{Type: TokenPunctuator, Data: "."},
				{Type: TokenPrivateIdentifier, Data: "#b"},
			},
		},
		{ // 132
			"#",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid character sequence: #"},
			},
		},
		{ // 133
			"#.a",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid character sequence: #."},
			},
		},
		{ // 134
			"Number(10000n * this.#numerator / this.#denominator) / 10000;",
			[]parser.Token{
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
			},
		},
		{ // 135
			"`\\x`",
			[]parser.Token{
				{Type: parser.TokenError, Data: "invalid escape sequence: `\\x`"},
			},
		},
		{ // 136
			".123_456",
			[]parser.Token{
				{Type: TokenNumericLiteral, Data: ".123_456"},
			},
		},
		{ // 137
			"`a\\`b`",
			[]parser.Token{
				{Type: TokenNoSubstitutionTemplate, Data: "`a\\`b`"},
			},
		},
	} {
		p := parser.NewStringTokeniser(test.Input)

		SetTokeniser(&p)

		for m, tkn := range test.Output {
			if tk, _ := p.GetToken(); tk.Type != tkn.Type {
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
