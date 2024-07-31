package javascript

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"vimagination.zapto.org/javascript/internal"
	"vimagination.zapto.org/parser"
)

const (
	whitespace      = "\t\v\f \xa0\ufeff" // Tab, Vertical Tab, Form Feed, Space, No-break space, ZeroWidth No-Break Space, https://262.ecma-international.org/11.0/#table-32
	lineTerminators = "\n\r\u2028\u2029"  // Line Feed, Carriage Return, Line Separator, Paragraph Separator, https://262.ecma-international.org/11.0/#table-33

	singleEscapeChar = "'\"\\bfnrtv"
	binaryDigit      = "01"
	octalDigit       = "01234567"
	decimalDigit     = "0123456789"
	hexDigit         = "0123456789abcdefABCDEF"
)

var keywords = [...]string{"await", "break", "case", "catch", "class", "const", "continue", "debugger", "default", "delete", "do", "else", "enum", "export", "extends", "finally", "for", "function", "if", "import", "in", "instanceof", "new", "return", "super", "switch", "this", "throw", "try", "typeof", "var", "void", "while", "with", "yield"}

// Javascript Token values
const (
	TokenWhitespace parser.TokenType = iota
	TokenLineTerminator
	TokenSingleLineComment
	TokenMultiLineComment
	TokenIdentifier
	TokenPrivateIdentifier
	TokenBooleanLiteral
	TokenKeyword
	TokenPunctuator
	TokenNumericLiteral
	TokenStringLiteral
	TokenNoSubstitutionTemplate
	TokenTemplateHead
	TokenTemplateMiddle
	TokenTemplateTail
	TokenDivPunctuator
	TokenRightBracePunctuator
	TokenRegularExpressionLiteral
	TokenNullLiteral
	TokenFutureReservedWord
)

type jsTokeniser struct {
	tokenDepth      []byte
	divisionAllowed bool
}

// SetTokeniser provides javascript parsing functions to a Tokeniser
func SetTokeniser(t *parser.Tokeniser) *parser.Tokeniser {
	t.TokeniserState(new(jsTokeniser).inputElement)

	return t
}

func (j *jsTokeniser) inputElement(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Accept(whitespace) {
		t.AcceptRun(whitespace)

		return t.Return(TokenWhitespace, j.inputElement)
	}

	if t.Accept(lineTerminators) {
		t.AcceptRun(lineTerminators)

		return t.Return(TokenLineTerminator, j.inputElement)
	}

	allowDivision := j.divisionAllowed
	j.divisionAllowed = false

	switch c := t.Peek(); c {
	case -1:
		if len(j.tokenDepth) == 0 {
			return t.Done()
		}

		t.Err = io.ErrUnexpectedEOF

		return t.Error()
	case '/':
		t.Next()

		if t.Accept("/") {
			t.ExceptRun(lineTerminators)

			return t.Return(TokenSingleLineComment, j.inputElement)
		}

		if t.Accept("*") {
			for {
				t.ExceptRun("*")
				t.Accept("*")

				if t.Accept("/") {
					j.divisionAllowed = allowDivision

					return t.Return(TokenMultiLineComment, j.inputElement)
				}
				if t.Peek() == -1 {
					t.Err = io.ErrUnexpectedEOF

					break
				}
			}

			return t.Error()
		}

		if allowDivision {
			t.Accept("=")

			return t.Return(TokenDivPunctuator, j.inputElement)
		}

		j.divisionAllowed = true

		return j.regexp(t)
	case '}':
		t.Next()

		switch j.lastDepth() {
		case '{':
			j.tokenDepth = j.tokenDepth[:len(j.tokenDepth)-1]

			return t.Return(TokenRightBracePunctuator, j.inputElement)
		case '$':
			j.tokenDepth = j.tokenDepth[:len(j.tokenDepth)-1]

			return j.template(t)
		}

		t.Err = fmt.Errorf("%w: %s", ErrInvalidCharacter, t.Get())

		return t.Error()
	case '\'', '"':
		j.divisionAllowed = true

		return j.stringToken(t)
	case '`':
		t.Next()

		return j.template(t)
	case '#':
		t.Next()

		if !internal.IsIDStart(t.Peek()) {
			t.Next()

			t.Err = fmt.Errorf("%w: %s", ErrInvalidSequence, t.Get())

			return t.Error()
		}

		tk, tf := j.identifier(t)

		if tk.Type == TokenIdentifier {
			tk.Type = TokenPrivateIdentifier
			j.divisionAllowed = true
		}

		return tk, tf
	default:
		if strings.ContainsRune(decimalDigit, c) {
			j.divisionAllowed = true

			return j.number(t)
		}

		if internal.IsIDStart(c) {
			tk, tf := j.identifier(t)

			if tk.Type == TokenIdentifier {
				if tk.Data == "true" || tk.Data == "false" {
					j.divisionAllowed = true
					tk.Type = TokenBooleanLiteral
				} else if tk.Data == "null" {
					j.divisionAllowed = true
					tk.Type = TokenNullLiteral
				} else if tk.Data == "enum" {
					j.divisionAllowed = true
					tk.Type = TokenFutureReservedWord
				} else if tk.Data == "Infinity" {
					j.divisionAllowed = true
					tk.Type = TokenNumericLiteral
				} else {
					for _, kw := range keywords {
						if kw == tk.Data {
							tk.Type = TokenKeyword

							if tk.Data == "this" {
								j.divisionAllowed = true
							}

							break
						}
					}

					if tk.Type == TokenIdentifier {
						if tk.Data[0] == '\\' {
							code := ""

							if tk.Data[2] == '{' {
								n := 3

								for ; tk.Data[n] != '}'; n++ {
								}

								code = tk.Data[3:n]
							} else {
								code = tk.Data[2:6]
							}

							r, err := strconv.ParseInt(code, 16, 64)
							if err != nil || r == 92 || !internal.IsIDStart(rune(r)) {
								t.Err = fmt.Errorf("%w: %s", ErrInvalidUnicode, tk.Data)

								return t.Error()
							}
						}

						j.divisionAllowed = true
					}
				}
			}

			return tk, tf
		}

		t.Next()

		switch c {
		case '{', '(', '[':
			j.tokenDepth = append(j.tokenDepth, byte(c))
		case '?':
			if t.Accept("?") {
				t.Accept("=")
			} else {
				t.Accept(".")
			}
		case ';', ',', ':', '~', '>':
		case ')', ']':
			if ld := j.lastDepth(); !(ld == '(' && c == ')') && !(ld == '[' && c == ']') {
				t.Err = fmt.Errorf("%w: %s", ErrInvalidCharacter, t.Get())

				return t.Error()
			}

			j.divisionAllowed = true
			j.tokenDepth = j.tokenDepth[:len(j.tokenDepth)-1]
		case '.':
			if t.Accept(".") {
				if !t.Accept(".") { // ...
					if t.Next() == -1 {
						t.Err = io.ErrUnexpectedEOF
					} else {
						t.Err = fmt.Errorf("%w: %s", ErrInvalidSequence, t.Get())
					}

					return t.Error()
				}
			} else if t.Accept(decimalDigit) {
				numberRun(t, decimalDigit)

				if t.Accept("eE") {
					t.Accept("+-")
					numberRun(t, decimalDigit)
				}

				j.divisionAllowed = true

				return t.Return(TokenNumericLiteral, j.inputElement)
			}
		case '<', '*':
			if !t.Accept("=") { // <=, *=
				if t.Peek() == c { // <<, **
					t.Next()
					t.Accept("=") // <<=, **=
				}
			}
		case '=':
			if t.Accept("=") { // ==
				t.Accept("=") // ===
			} else {
				t.Accept(">") // =>
			}
		case '!':
			if t.Accept("=") { // !=
				t.Accept("=") // !==
			}
		case '+', '-', '&', '|':
			if t.Peek() == c {
				t.Next() // ++, --, &&, ||

				if c == '&' || c == '|' {
					t.Accept("=")
				}
			} else {
				t.Accept("=") // +=, -=, &=, |=
			}
		case '%', '^':
			t.Accept("=") // %=, ^=
		default:
			t.Err = fmt.Errorf("%w: %s", ErrInvalidCharacter, t.Get())

			return t.Error()
		}

		return t.Return(TokenPunctuator, j.inputElement)
	}
}

func (j *jsTokeniser) regexpBackslashSequence(t *parser.Tokeniser) bool {
	t.Next()

	if !t.Except(lineTerminators) {
		if t.Peek() == -1 {
			t.Err = io.ErrUnexpectedEOF
		} else {
			t.Next()

			t.Err = fmt.Errorf("%w: %s", ErrInvalidRegexpSequence, t.Get())
		}

		return false
	}

	return true
}

func (j *jsTokeniser) regexpExpressionClass(t *parser.Tokeniser) bool {
	t.Next()

	for {
		switch t.ExceptRun("]\\") {
		case ']':
			t.Next()

			return true
		case '\\':
			if !j.regexpBackslashSequence(t) {
				return false
			}
		default:
			t.Err = io.ErrUnexpectedEOF

			return false
		}
	}
}

func (j *jsTokeniser) regexp(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	switch c := t.Peek(); c {
	case -1:
		t.Err = io.ErrUnexpectedEOF

		return t.Error()
	case '\\':
		if !j.regexpBackslashSequence(t) {
			return t.Error()
		}
	case '[':
		if !j.regexpExpressionClass(t) {
			return t.Error()
		}
	default:
		if strings.ContainsRune(lineTerminators, c) {
			t.Get()
			t.Next()

			t.Err = fmt.Errorf("%w: %s", ErrInvalidRegexpCharacter, t.Get())

			return t.Error()
		}

		t.Next()
	}

Loop:
	for {
		switch c := t.ExceptRun(lineTerminators + "\\[/"); c {
		case -1:
			t.Err = io.ErrUnexpectedEOF

			return t.Error()
		case '\\':
			if !j.regexpBackslashSequence(t) {
				return t.Error()
			}
		case '[':
			if !j.regexpExpressionClass(t) {
				return t.Error()
			}
		case '/':
			t.Next()

			break Loop
		default:
			t.Get()
			t.Next()

			t.Err = fmt.Errorf("%w: %s", ErrInvalidRegexpCharacter, t.Get())

			return t.Error()
		}
	}
	for {
		if c := t.Peek(); !internal.IsIDContinue(c) || c == '\\' {
			break
		}

		t.Next()
	}

	return t.Return(TokenRegularExpressionLiteral, j.inputElement)
}

func numberRun(t *parser.Tokeniser, digits string) bool {
	for {
		if !t.Accept(digits) {
			return false
		}

		t.AcceptRun(digits)

		if !t.Accept("_") {
			break
		}
	}

	return true
}

func (j *jsTokeniser) number(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Accept("0") {
		if t.Accept("bB") {
			if !numberRun(t, binaryDigit) {
				t.Next()

				t.Err = fmt.Errorf("%w: %s", ErrInvalidNumber, t.Get())

				return t.Error()
			}

			t.Accept("n")
		} else if t.Accept("oO") {
			if !numberRun(t, octalDigit) {
				t.Next()

				t.Err = fmt.Errorf("%w: %s", ErrInvalidNumber, t.Get())

				return t.Error()
			}

			t.Accept("n")
		} else if t.Accept("xX") {
			if !numberRun(t, hexDigit) {
				t.Next()

				t.Err = fmt.Errorf("%w: %s", ErrInvalidNumber, t.Get())

				return t.Error()
			}

			t.Accept("n")
		} else if t.Accept(".") {
			if !numberRun(t, decimalDigit) {
				t.Next()

				t.Err = fmt.Errorf("%w: %s", ErrInvalidNumber, t.Get())

				return t.Error()
			}

			if t.Accept("eE") {
				t.Accept("+-")

				if !numberRun(t, decimalDigit) {
					t.Next()

					t.Err = fmt.Errorf("%w: %s", ErrInvalidNumber, t.Get())

					return t.Error()
				}
			}
		} else {
			t.Accept("n")
		}
	} else {
		if !numberRun(t, decimalDigit) {
			t.Next()

			t.Err = fmt.Errorf("%w: %s", ErrInvalidNumber, t.Get())

			return t.Error()
		}

		if !t.Accept("n") {
			if t.Accept(".") {
				if !numberRun(t, decimalDigit) {
					t.Next()

					t.Err = fmt.Errorf("%w: %s", ErrInvalidNumber, t.Get())

					return t.Error()
				}
			}

			if t.Accept("eE") {
				t.Accept("+-")

				if !numberRun(t, decimalDigit) {
					t.Next()

					t.Err = fmt.Errorf("%w: %s", ErrInvalidNumber, t.Get())

					return t.Error()
				}
			}
		}
	}

	return t.Return(TokenNumericLiteral, j.inputElement)
}

func (j *jsTokeniser) identifier(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	c := t.Next()

	if c == '\\' {
		if !t.Accept("u") {
			t.Next()

			t.Err = fmt.Errorf("%w: %s", ErrUnexpectedBackslash, t.Get())

			return t.Error()
		}
		if !j.unicodeEscapeSequence(t) {
			t.Err = fmt.Errorf("%w: %s", ErrInvalidUnicode, t.Get())

			return t.Error()
		}
	}

	for {
		c = t.Peek()

		if internal.IsIDContinue(c) {
			t.Next()

			continue
		}

		break
	}

	return t.Return(TokenIdentifier, j.inputElement)
}

func (j *jsTokeniser) unicodeEscapeSequence(t *parser.Tokeniser) bool {
	if t.Accept("{") {
		if !t.Accept(hexDigit) {
			t.Next()

			return false
		}

		t.AcceptRun(hexDigit)

		if !t.Accept("}") {
			return false
		}
	} else if !t.Accept(hexDigit) || !t.Accept(hexDigit) || !t.Accept(hexDigit) || !t.Accept(hexDigit) {
		t.Next()

		return false
	}

	return true
}

func (j *jsTokeniser) lastDepth() rune {
	if len(j.tokenDepth) == 0 {
		return -1
	}

	return rune(j.tokenDepth[len(j.tokenDepth)-1])
}

func (j *jsTokeniser) escapeSequence(t *parser.Tokeniser) bool {
	t.Accept("\\")

	if t.Accept("x") {
		return t.Accept(hexDigit) && t.Accept(hexDigit)
	} else if t.Accept("u") {
		return j.unicodeEscapeSequence(t)
	} else if t.Accept("0") {
		return !t.Accept(decimalDigit)
	}

	t.Except(lineTerminators)

	return true
}

var (
	stringChars       = "'\\" + lineTerminators + "\""
	doubleStringChars = stringChars[1:]
	singleStringChars = stringChars[:len(stringChars)-1]
)

func (j *jsTokeniser) stringToken(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	var chars string

	if t.Next() == '"' {
		chars = doubleStringChars
	} else {
		chars = singleStringChars
	}

Loop:
	for {
		switch c := t.ExceptRun(chars); c {
		case '"', '\'':
			t.Next()

			break Loop
		case '\\':
			if j.escapeSequence(t) {
				continue
			}

			if t.Err == nil {
				t.Err = fmt.Errorf("%w: %s", ErrInvalidEscapeSequence, t.Get())
			}
		default:
			t.Err = io.ErrUnexpectedEOF

			if strings.ContainsRune(lineTerminators, c) {
				t.Next()

				t.Err = fmt.Errorf("%w: %s", ErrUnexpectedLineTerminator, t.Get())
			}
		}

		return t.Error()
	}

	return t.Return(TokenStringLiteral, j.inputElement)
}

func (j *jsTokeniser) template(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
Loop:
	for {
		switch t.ExceptRun("`\\$") {
		case '`':
			t.Next()

			break Loop
		case '\\':
			if j.escapeSequence(t) {
				continue
			}

			t.Next()

			t.Err = fmt.Errorf("%w: %s", ErrInvalidEscapeSequence, t.Get())

			return t.Error()
		case '$':
			t.Next()

			if t.Accept("{") {
				j.tokenDepth = append(j.tokenDepth, '$')
				v := t.Get()

				var typ parser.TokenType

				if v[0] == '`' {
					typ = TokenTemplateHead
				} else {
					typ = TokenTemplateMiddle
				}

				return parser.Token{
					Type: typ,
					Data: v,
				}, j.inputElement
			}
		default:
			t.Err = io.ErrUnexpectedEOF

			return t.Error()
		}
	}

	j.divisionAllowed = true
	v := t.Get()

	var typ parser.TokenType

	if v[0] == '`' {
		typ = TokenNoSubstitutionTemplate
	} else {
		typ = TokenTemplateTail
	}

	return parser.Token{
		Type: typ,
		Data: v,
	}, j.inputElement
}
