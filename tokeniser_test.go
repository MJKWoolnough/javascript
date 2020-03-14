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
				{parser.TokenDone, ""},
			},
		},
		{ // 2
			" ",
			[]parser.Token{
				{TokenWhitespace, " "},
				{parser.TokenDone, ""},
			},
		},
		{ // 3
			" 	",
			[]parser.Token{
				{TokenWhitespace, " 	"},
				{parser.TokenDone, ""},
			},
		},
		{ // 4
			" \n ",
			[]parser.Token{
				{TokenWhitespace, " "},
				{TokenLineTerminator, "\n"},
				{TokenWhitespace, " "},
				{parser.TokenDone, ""},
			},
		},
		{ // 5
			"\"\"",
			[]parser.Token{
				{TokenStringLiteral, "\"\""},
				{parser.TokenDone, ""},
			},
		},
		{ // 6
			"\"\\\"\"",
			[]parser.Token{
				{TokenStringLiteral, "\"\\\"\""},
				{parser.TokenDone, ""},
			},
		},
		{ // 7
			"\"\n\"",
			[]parser.Token{
				{parser.TokenError, "line terminator in string: \"\n"},
			},
		},
		{ // 8
			"\"\\n\"",
			[]parser.Token{
				{TokenStringLiteral, "\"\\n\""},
				{parser.TokenDone, ""},
			},
		},
		{ // 9
			"\"\\0\"",
			[]parser.Token{
				{TokenStringLiteral, "\"\\0\""},
				{parser.TokenDone, ""},
			},
		},
		{ // 10
			"\"\\x20\"",
			[]parser.Token{
				{TokenStringLiteral, "\"\\x20\""},
				{parser.TokenDone, ""},
			},
		},
		{ // 11
			"\"\\u2020\"",
			[]parser.Token{
				{TokenStringLiteral, "\"\\u2020\""},
				{parser.TokenDone, ""},
			},
		},
		{ // 12
			"\"\\u\"",
			[]parser.Token{
				{parser.TokenError, "invalid escape sequence: \"\\u\""},
			},
		},
		{ // 13
			"\"\\up\"",
			[]parser.Token{
				{parser.TokenError, "invalid escape sequence: \"\\up"},
			},
		},
		{ // 14
			"\"\\u{20}\"",
			[]parser.Token{
				{TokenStringLiteral, "\"\\u{20}\""},
				{parser.TokenDone, ""},
			},
		},
		{ // 15
			"\"use strict\"",
			[]parser.Token{
				{TokenStringLiteral, "\"use strict\""},
				{parser.TokenDone, ""},
			},
		},
		{ // 16
			"\"use\\u{20}strict\\x65!\\0\"",
			[]parser.Token{
				{TokenStringLiteral, "\"use\\u{20}strict\\x65!\\0\""},
				{parser.TokenDone, ""},
			},
		},
		{ // 17
			"\"use strict\";",
			[]parser.Token{
				{TokenStringLiteral, "\"use strict\""},
				{TokenPunctuator, ";"},
				{parser.TokenDone, ""},
			},
		},
		{ // 18
			"0",
			[]parser.Token{
				{TokenNumericLiteral, "0"},
				{parser.TokenDone, ""},
			},
		},
		{ // 19
			"0.1",
			[]parser.Token{
				{TokenNumericLiteral, "0.1"},
				{parser.TokenDone, ""},
			},
		},
		{ // 20
			".1",
			[]parser.Token{
				{TokenNumericLiteral, ".1"},
				{parser.TokenDone, ""},
			},
		},
		{ // 21
			"0b0",
			[]parser.Token{
				{TokenNumericLiteral, "0b0"},
				{parser.TokenDone, ""},
			},
		},
		{ // 22
			"0b1",
			[]parser.Token{
				{TokenNumericLiteral, "0b1"},
				{parser.TokenDone, ""},
			},
		},
		{ // 23
			"0b1001010101",
			[]parser.Token{
				{TokenNumericLiteral, "0b1001010101"},
				{parser.TokenDone, ""},
			},
		},
		{ // 24
			"0",
			[]parser.Token{
				{TokenNumericLiteral, "0"},
				{parser.TokenDone, ""},
			},
		},
		{ // 25
			"1",
			[]parser.Token{
				{TokenNumericLiteral, "1"},
				{parser.TokenDone, ""},
			},
		},
		{ // 26
			"9",
			[]parser.Token{
				{TokenNumericLiteral, "9"},
				{parser.TokenDone, ""},
			},
		},
		{ // 27
			"12345678901",
			[]parser.Token{
				{TokenNumericLiteral, "12345678901"},
				{parser.TokenDone, ""},
			},
		},
		{ // 28
			"12345678.901",
			[]parser.Token{
				{TokenNumericLiteral, "12345678.901"},
				{parser.TokenDone, ""},
			},
		},
		{ // 29
			"12345678901E123",
			[]parser.Token{
				{TokenNumericLiteral, "12345678901E123"},
				{parser.TokenDone, ""},
			},
		},
		{ // 30
			"12345678901e+123",
			[]parser.Token{
				{TokenNumericLiteral, "12345678901e+123"},
				{parser.TokenDone, ""},
			},
		},
		{ // 31
			"12345678.901E-123",
			[]parser.Token{
				{TokenNumericLiteral, "12345678.901E-123"},
				{parser.TokenDone, ""},
			},
		},
		{ // 32
			"0x0",
			[]parser.Token{
				{TokenNumericLiteral, "0x0"},
				{parser.TokenDone, ""},
			},
		},
		{ // 33
			"0xa",
			[]parser.Token{
				{TokenNumericLiteral, "0xa"},
				{parser.TokenDone, ""},
			},
		},
		{ // 34
			"0xf",
			[]parser.Token{
				{TokenNumericLiteral, "0xf"},
				{parser.TokenDone, ""},
			},
		},
		{ // 35
			"0x0f",
			[]parser.Token{
				{TokenNumericLiteral, "0x0f"},
				{parser.TokenDone, ""},
			},
		},
		{ // 36
			"0xaf",
			[]parser.Token{
				{TokenNumericLiteral, "0xaf"},
				{parser.TokenDone, ""},
			},
		},
		{ // 37
			"0xDeAdBeEf",
			[]parser.Token{
				{TokenNumericLiteral, "0xDeAdBeEf"},
				{parser.TokenDone, ""},
			},
		},
		{ // 38
			"0n",
			[]parser.Token{
				{TokenNumericLiteral, "0n"},
				{parser.TokenDone, ""},
			},
		},
		{ // 39
			"1n",
			[]parser.Token{
				{TokenNumericLiteral, "1n"},
				{parser.TokenDone, ""},
			},
		},
		{ // 40
			"1234567890n",
			[]parser.Token{
				{TokenNumericLiteral, "1234567890n"},
				{parser.TokenDone, ""},
			},
		},
		{ // 41
			"0x1234567890n",
			[]parser.Token{
				{TokenNumericLiteral, "0x1234567890n"},
				{parser.TokenDone, ""},
			},
		},
		{ // 42
			"Infinity",
			[]parser.Token{
				{TokenNumericLiteral, "Infinity"},
				{parser.TokenDone, ""},
			},
		},
		{ // 43
			"true",
			[]parser.Token{
				{TokenBooleanLiteral, "true"},
				{parser.TokenDone, ""},
			},
		},
		{ // 44
			"false",
			[]parser.Token{
				{TokenBooleanLiteral, "false"},
				{parser.TokenDone, ""},
			},
		},
		{ // 45
			"hello",
			[]parser.Token{
				{TokenIdentifier, "hello"},
				{parser.TokenDone, ""},
			},
		},
		{ // 46
			"this",
			[]parser.Token{
				{TokenKeyword, "this"},
				{parser.TokenDone, ""},
			},
		},
		{ // 47
			"function",
			[]parser.Token{
				{TokenKeyword, "function"},
				{parser.TokenDone, ""},
			},
		},
		{ // 48
			"/[a-z]+/g",
			[]parser.Token{
				{TokenRegularExpressionLiteral, "/[a-z]+/g"},
				{parser.TokenDone, ""},
			},
		},
		{ // 49
			"/[\\n]/g",
			[]parser.Token{
				{TokenRegularExpressionLiteral, "/[\\n]/g"},
				{parser.TokenDone, ""},
			},
		},
		{ // 50
			"var a =	 /^ab[cd]*$/ig;",
			[]parser.Token{
				{TokenKeyword, "var"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "a"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "="},
				{TokenWhitespace, "	 "},
				{TokenRegularExpressionLiteral, "/^ab[cd]*$/ig"},
				{TokenPunctuator, ";"},
				{parser.TokenDone, ""},
			},
		},
		{ // 51
			"num /= 4",
			[]parser.Token{
				{TokenIdentifier, "num"},
				{TokenWhitespace, " "},
				{TokenDivPunctuator, "/="},
				{TokenWhitespace, " "},
				{TokenNumericLiteral, "4"},
				{parser.TokenDone, ""},
			},
		},
		{ // 52
			"const num = 8 / 4;",
			[]parser.Token{
				{TokenKeyword, "const"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "num"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "="},
				{TokenWhitespace, " "},
				{TokenNumericLiteral, "8"},
				{TokenWhitespace, " "},
				{TokenDivPunctuator, "/"},
				{TokenWhitespace, " "},
				{TokenNumericLiteral, "4"},
				{TokenPunctuator, ";"},
				{parser.TokenDone, ""},
			},
		},
		{ // 53
			"``",
			[]parser.Token{
				{TokenNoSubstitutionTemplate, "``"},
				{parser.TokenDone, ""},
			},
		},
		{ // 54
			"`abc`",
			[]parser.Token{
				{TokenNoSubstitutionTemplate, "`abc`"},
				{parser.TokenDone, ""},
			},
		},
		{ // 55
			"`ab${ (val.a / 2) + 1 }c${str}`",
			[]parser.Token{
				{TokenTemplateHead, "`ab${"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "("},
				{TokenIdentifier, "val"},
				{TokenPunctuator, "."},
				{TokenIdentifier, "a"},
				{TokenWhitespace, " "},
				{TokenDivPunctuator, "/"},
				{TokenWhitespace, " "},
				{TokenNumericLiteral, "2"},
				{TokenPunctuator, ")"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "+"},
				{TokenWhitespace, " "},
				{TokenNumericLiteral, "1"},
				{TokenWhitespace, " "},
				{TokenTemplateMiddle, "}c${"},
				{TokenIdentifier, "str"},
				{TokenTemplateTail, "}`"},
				{parser.TokenDone, ""},
			},
		},
		{ // 56
			"const myFunc = function(aye, bee, cea) {\n	const num = [123, 4, lastNum(aye, \"beep\", () => window, val => val * 2, (myVar) => {myVar /= 2;return myVar;})], elm = document.getElementByID();\n	console.log(bee, num, elm);}",
			[]parser.Token{
				{TokenKeyword, "const"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "myFunc"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "="},
				{TokenWhitespace, " "},
				{TokenKeyword, "function"},
				{TokenPunctuator, "("},
				{TokenIdentifier, "aye"},
				{TokenPunctuator, ","},
				{TokenWhitespace, " "},
				{TokenIdentifier, "bee"},
				{TokenPunctuator, ","},
				{TokenWhitespace, " "},
				{TokenIdentifier, "cea"},
				{TokenPunctuator, ")"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "{"},
				{TokenLineTerminator, "\n"},
				{TokenWhitespace, "	"},
				{TokenKeyword, "const"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "num"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "="},
				{TokenWhitespace, " "},
				{TokenPunctuator, "["},
				{TokenNumericLiteral, "123"},
				{TokenPunctuator, ","},
				{TokenWhitespace, " "},
				{TokenNumericLiteral, "4"},
				{TokenPunctuator, ","},
				{TokenWhitespace, " "},
				{TokenIdentifier, "lastNum"},
				{TokenPunctuator, "("},
				{TokenIdentifier, "aye"},
				{TokenPunctuator, ","},
				{TokenWhitespace, " "},
				{TokenStringLiteral, "\"beep\""},
				{TokenPunctuator, ","},
				{TokenWhitespace, " "},
				{TokenPunctuator, "("},
				{TokenPunctuator, ")"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "=>"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "window"},
				{TokenPunctuator, ","},
				{TokenWhitespace, " "},
				{TokenIdentifier, "val"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "=>"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "val"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "*"},
				{TokenWhitespace, " "},
				{TokenNumericLiteral, "2"},
				{TokenPunctuator, ","},
				{TokenWhitespace, " "},
				{TokenPunctuator, "("},
				{TokenIdentifier, "myVar"},
				{TokenPunctuator, ")"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "=>"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "{"},
				{TokenIdentifier, "myVar"},
				{TokenWhitespace, " "},
				{TokenDivPunctuator, "/="},
				{TokenWhitespace, " "},
				{TokenNumericLiteral, "2"},
				{TokenPunctuator, ";"},
				{TokenKeyword, "return"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "myVar"},
				{TokenPunctuator, ";"},
				{TokenRightBracePunctuator, "}"},
				{TokenPunctuator, ")"},
				{TokenPunctuator, "]"},
				{TokenPunctuator, ","},
				{TokenWhitespace, " "},
				{TokenIdentifier, "elm"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "="},
				{TokenWhitespace, " "},
				{TokenIdentifier, "document"},
				{TokenPunctuator, "."},
				{TokenIdentifier, "getElementByID"},
				{TokenPunctuator, "("},
				{TokenPunctuator, ")"},
				{TokenPunctuator, ";"},
				{TokenLineTerminator, "\n"},
				{TokenWhitespace, "	"},
				{TokenIdentifier, "console"},
				{TokenPunctuator, "."},
				{TokenIdentifier, "log"},
				{TokenPunctuator, "("},
				{TokenIdentifier, "bee"},
				{TokenPunctuator, ","},
				{TokenWhitespace, " "},
				{TokenIdentifier, "num"},
				{TokenPunctuator, ","},
				{TokenWhitespace, " "},
				{TokenIdentifier, "elm"},
				{TokenPunctuator, ")"},
				{TokenPunctuator, ";"},
				{TokenRightBracePunctuator, "}"},
				{parser.TokenDone, ""},
			},
		},
		{ // 57
			"export {name1, name2};",
			[]parser.Token{
				{TokenKeyword, "export"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "{"},
				{TokenIdentifier, "name1"},
				{TokenPunctuator, ","},
				{TokenWhitespace, " "},
				{TokenIdentifier, "name2"},
				{TokenRightBracePunctuator, "}"},
				{TokenPunctuator, ";"},
				{parser.TokenDone, ""},
			},
		},
		{ // 58
			"export {var1 as name1, var2 as name2};",
			[]parser.Token{
				{TokenKeyword, "export"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "{"},
				{TokenIdentifier, "var1"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "as"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "name1"},
				{TokenPunctuator, ","},
				{TokenWhitespace, " "},
				{TokenIdentifier, "var2"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "as"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "name2"},
				{TokenRightBracePunctuator, "}"},
				{TokenPunctuator, ";"},
				{parser.TokenDone, ""},
			},
		},
		{ // 59
			"export * from './other.js';",
			[]parser.Token{
				{TokenKeyword, "export"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "*"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "from"},
				{TokenWhitespace, " "},
				{TokenStringLiteral, "'./other.js'"},
				{TokenPunctuator, ";"},
				{parser.TokenDone, ""},
			},
		},
		{ // 60
			"import * as name from './module.js';",
			[]parser.Token{
				{TokenKeyword, "import"},
				{TokenWhitespace, " "},
				{TokenPunctuator, "*"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "as"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "name"},
				{TokenWhitespace, " "},
				{TokenIdentifier, "from"},
				{TokenWhitespace, " "},
				{TokenStringLiteral, "'./module.js'"},
				{TokenPunctuator, ";"},
				{parser.TokenDone, ""},
			},
		},
		{ // 61
			"$",
			[]parser.Token{
				{TokenIdentifier, "$"},
				{parser.TokenDone, ""},
			},
		},
		{ // 62
			"_",
			[]parser.Token{
				{TokenIdentifier, "_"},
				{parser.TokenDone, ""},
			},
		},
		{ // 63
			"\\u0060",
			[]parser.Token{
				{TokenIdentifier, "\\u0060"},
				{parser.TokenDone, ""},
			},
		},
		{ // 64
			"// Comment",
			[]parser.Token{
				{TokenSingleLineComment, "// Comment"},
				{parser.TokenDone, ""},
			},
		},
		{ // 65
			"enum",
			[]parser.Token{
				{TokenFutureReservedWord, "enum"},
				{parser.TokenDone, ""},
			},
		},
		{ // 66
			".01234E56",
			[]parser.Token{
				{TokenNumericLiteral, ".01234E56"},
				{parser.TokenDone, ""},
			},
		},
		{ // 67
			"0.01234E56",
			[]parser.Token{
				{TokenNumericLiteral, "0.01234E56"},
				{parser.TokenDone, ""},
			},
		},
		{ // 68
			"0o1234567",
			[]parser.Token{
				{TokenNumericLiteral, "0o1234567"},
				{parser.TokenDone, ""},
			},
		},
		{ // 69
			"`\\x60`",
			[]parser.Token{
				{TokenNoSubstitutionTemplate, "`\\x60`"},
				{parser.TokenDone, ""},
			},
		},
		{ // 70
			"/\\(/",
			[]parser.Token{
				{TokenRegularExpressionLiteral, "/\\(/"},
				{parser.TokenDone, ""},
			},
		},
		{ // 71
			"/a\\(/",
			[]parser.Token{
				{TokenRegularExpressionLiteral, "/a\\(/"},
				{parser.TokenDone, ""},
			},
		},
		{ // 72
			"{",
			[]parser.Token{
				{TokenPunctuator, "{"},
				{parser.TokenError, "unexpected EOF"},
			},
		},
		{ // 73
			"[",
			[]parser.Token{
				{TokenPunctuator, "["},
				{parser.TokenError, "unexpected EOF"},
			},
		},
		{ // 74
			"(",
			[]parser.Token{
				{TokenPunctuator, "("},
				{parser.TokenError, "unexpected EOF"},
			},
		},
		{ // 75
			"/*",
			[]parser.Token{
				{parser.TokenError, "unexpected EOF"},
			},
		},
		{ // 76
			"[}",
			[]parser.Token{
				{TokenPunctuator, "["},
				{parser.TokenError, "invalid character: }"},
			},
		},
		{ // 77
			"(}",
			[]parser.Token{
				{TokenPunctuator, "("},
				{parser.TokenError, "invalid character: }"},
			},
		},
		{ // 78
			"{)",
			[]parser.Token{
				{TokenPunctuator, "{"},
				{parser.TokenError, "invalid character: )"},
			},
		},
		{ // 79
			"{]",
			[]parser.Token{
				{TokenPunctuator, "{"},
				{parser.TokenError, "invalid character: ]"},
			},
		},
		{ // 80
			"(]",
			[]parser.Token{
				{TokenPunctuator, "("},
				{parser.TokenError, "invalid character: ]"},
			},
		},
		{ // 81
			"[)",
			[]parser.Token{
				{TokenPunctuator, "["},
				{parser.TokenError, "invalid character: )"},
			},
		},
		{ // 82
			"..",
			[]parser.Token{
				{parser.TokenError, "unexpected EOF"},
			},
		},
		{ // 83
			"..a",
			[]parser.Token{
				{parser.TokenError, "invalid character sequence: ..a"},
			},
		},
		{ // 84
			"/\\\n/",
			[]parser.Token{
				{parser.TokenError, "invalid regexp sequence: /\\\n"},
			},
		},
		{ // 85
			"/[",
			[]parser.Token{
				{parser.TokenError, "unexpected EOF"},
			},
		},
		{ // 86
			"/[\\",
			[]parser.Token{
				{parser.TokenError, "unexpected EOF"},
			},
		},
		{ // 87
			"/",
			[]parser.Token{
				{parser.TokenError, "unexpected EOF"},
			},
		},
		{ // 88
			"/\n",
			[]parser.Token{
				{parser.TokenError, "invalid regexp character: \n"},
			},
		},
		{ // 89
			"/a",
			[]parser.Token{
				{parser.TokenError, "unexpected EOF"},
			},
		},
		{ // 90
			"/a\\\n/",
			[]parser.Token{
				{parser.TokenError, "invalid regexp sequence: /a\\\n"},
			},
		},
		{ // 91
			"/a[",
			[]parser.Token{
				{parser.TokenError, "unexpected EOF"},
			},
		},
		{ // 92
			"/a\n",
			[]parser.Token{
				{parser.TokenError, "invalid regexp character: \n"},
			},
		},
		{ // 93
			"0B9",
			[]parser.Token{
				{parser.TokenError, "invalid number: 0B9"},
			},
		},
		{ // 94
			"0O9",
			[]parser.Token{
				{parser.TokenError, "invalid number: 0O9"},
			},
		},
		{ // 95
			"0XG",
			[]parser.Token{
				{parser.TokenError, "invalid number: 0XG"},
			},
		},
		{ // 96
			"\\x60",
			[]parser.Token{
				{parser.TokenError, "unexpected backslash: \\x"},
			},
		},
		{ // 97
			"\\ug",
			[]parser.Token{
				{parser.TokenError, "invalid unicode escape sequence: \\ug"},
			},
		},
		{ // 98
			"\\u{G}",
			[]parser.Token{
				{parser.TokenError, "invalid unicode escape sequence: \\u{G"},
			},
		},
		{ // 99
			"\\u{ffff",
			[]parser.Token{
				{parser.TokenError, "invalid unicode escape sequence: \\u{ffff"},
			},
		},
		{ // 100
			"}",
			[]parser.Token{
				{parser.TokenError, "invalid character: }"},
			},
		},
		{ // 101
			"`\\G`",
			[]parser.Token{
				{parser.TokenError, "invalid escape sequence: `\\G"},
			},
		},
		{ // 102
			"`",
			[]parser.Token{
				{parser.TokenError, "unexpected EOF"},
			},
		},
		{ // 103
			"1_234_567",
			[]parser.Token{
				{TokenNumericLiteral, "1_234_567"},
				{parser.TokenDone, ""},
			},
		},
		{ // 104
			"1_",
			[]parser.Token{
				{parser.TokenError, "invalid number: 1_"},
			},
		},
		{ // 105
			"1__234_567",
			[]parser.Token{
				{parser.TokenError, "invalid number: 1__"},
			},
		},
		{ // 106
			"123e-456_789",
			[]parser.Token{
				{TokenNumericLiteral, "123e-456_789"},
				{parser.TokenDone, ""},
			},
		},
		{ // 107
			"0.123_456",
			[]parser.Token{
				{TokenNumericLiteral, "0.123_456"},
				{parser.TokenDone, ""},
			},
		},
		{ // 108
			"1.2_3_4_5_6_7",
			[]parser.Token{
				{TokenNumericLiteral, "1.2_3_4_5_6_7"},
				{parser.TokenDone, ""},
			},
		},
		{ // 109
			"0x1_2",
			[]parser.Token{
				{TokenNumericLiteral, "0x1_2"},
				{parser.TokenDone, ""},
			},
		},
		{ // 110
			"0b1_0",
			[]parser.Token{
				{TokenNumericLiteral, "0b1_0"},
				{parser.TokenDone, ""},
			},
		},
		{ // 111
			"0o1_7",
			[]parser.Token{
				{TokenNumericLiteral, "0o1_7"},
				{parser.TokenDone, ""},
			},
		},
		{ // 112
			"a.b",
			[]parser.Token{
				{TokenIdentifier, "a"},
				{TokenPunctuator, "."},
				{TokenIdentifier, "b"},
			},
		},
		{ // 113
			"a?.b",
			[]parser.Token{
				{TokenIdentifier, "a"},
				{TokenPunctuator, "?."},
				{TokenIdentifier, "b"},
			},
		},
		{ // 114
			"a??b",
			[]parser.Token{
				{TokenIdentifier, "a"},
				{TokenPunctuator, "??"},
				{TokenIdentifier, "b"},
			},
		},
		{ // 115
			"0.",
			[]parser.Token{
				{parser.TokenError, "invalid number: 0."},
			},
		},
		{ // 116
			"0.1e",
			[]parser.Token{
				{parser.TokenError, "invalid number: 0.1e"},
			},
		},
		{ // 117
			"1.",
			[]parser.Token{
				{parser.TokenError, "invalid number: 1."},
			},
		},
		{ // 118
			"1.1e",
			[]parser.Token{
				{parser.TokenError, "invalid number: 1.1e"},
			},
		},
	} {
		p := parser.NewStringTokeniser(test.Input)
		SetTokeniser(&p)
		for m, tkn := range test.Output {
			tk, _ := p.GetToken()
			if tk.Type != tkn.Type {
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
