package minify

import "unicode"

var (
	idContinue = []*unicode.RangeTable{
		unicode.L,
		unicode.Nl,
		unicode.Other_ID_Start,
		unicode.Mn,
		unicode.Mc,
		unicode.Nd,
		unicode.Pc,
		unicode.Other_ID_Continue,
	}
	idStart = idContinue[:3]
	notID   = []*unicode.RangeTable{
		unicode.Pattern_Syntax,
		unicode.Pattern_White_Space,
	}
	maxSafeInt = "9007199254740991"
)

const (
	zwnj rune = 8204
	zwj  rune = 8205
)

func isIDStart(c rune) bool {
	if c == '$' || c == '_' || c == '\\' {
		return true
	}
	return unicode.In(c, idStart...) && !unicode.In(c, notID...)
}

func isIDContinue(c rune) bool {
	if c == '$' || c == '_' || c == '\\' || c == zwnj || c == zwj {
		return true
	}
	return unicode.In(c, idContinue...) && !unicode.In(c, notID...)
}

func isIdentifier(str string) bool {
	for n, r := range str {
		if (n == 0 && !isIDStart(r)) || (n > 0 && !isIDContinue(r)) {
			return false
		}
	}
	return true
}

func isSimpleNumber(str string) bool {
	if str == "0" {
		return true
	}
	if len(str) == 0 || len(str) > len(maxSafeInt) || str[0] < '1' || str[0] > '9' {
		return false
	}
	for _, r := range str[1:] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(str) < len(maxSafeInt) || str <= maxSafeInt
}
