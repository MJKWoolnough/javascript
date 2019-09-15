package javascript

import (
	"strconv"

	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

func Unquote(str string) (string, error) {
	s := parser.NewStringTokeniser(str)
	var chars string
	if s.Accept("\"") {
		chars = doubleStringChars
	} else if s.Accept("'") {
		chars = singleStringChars
	} else {
		return "", ErrInvalidQuoted
	}
	s.Get()
	var ret string
Loop:
	for {
		switch s.ExceptRun(chars) {
		case '"', '\'':
			ret += s.Get()
			return ret, nil
		case '\\':
			ret += s.Get()
			s.Accept("\\")
			s.Get()
			if s.Accept("x") {
				s.Get()
				if !s.Accept(hexDigit) || !s.Accept(hexDigit) {
					break Loop
				}
				c, _ := strconv.ParseUint(s.Get(), 16, 8)
				ret += string(rune(c))
			} else if s.Accept("u") {
				s.Get()
				if s.Accept("{") {
					s.Get()
					if !s.Accept(hexDigit) {
						break Loop
					}
					s.AcceptRun(hexDigit)
					c, _ := strconv.ParseUint(s.Get(), 16, 8)
					ret += string(rune(c))
					if !s.Accept("}") {
						break Loop
					}
				} else if !s.Accept(hexDigit) || !s.Accept(hexDigit) || !s.Accept(hexDigit) || !s.Accept(hexDigit) {
					break Loop
				} else {
					c, _ := strconv.ParseUint(s.Get(), 16, 8)
					ret += string(rune(c))
				}
			} else if s.Accept("0") {
				if s.Accept(decimalDigit) {
					break Loop
				}
				s.Get()
				ret += "\000"
			} else if s.Accept(singleEscapeChar) {
				switch s.Get() {
				case "'":
					ret += singleEscapeChar[0:1]
				case "\"":
					ret += singleEscapeChar[1:2]
				case "\\":
					ret += singleEscapeChar[2:3]
				case "b":
					ret += "\b"
				case "f":
					ret += "\f"
				case "n":
					ret += "\n"
				case "r":
					ret += "\r"
				case "t":
					ret += "\t"
				case "v":
					ret += "\v"
				default:
					break Loop
				}
			} else {
				break Loop
			}
		default:
			break Loop
		}
	}
	return "", ErrInvalidQuoted
}

var (
	ErrInvalidQuoted = errors.New("invalid quoted string")
)
