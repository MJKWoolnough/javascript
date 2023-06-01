package internal

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
)

const (
	zwnj rune = 8204
	zwj  rune = 8205
)

func IsIDStart(c rune) bool {
	if c == '$' || c == '_' || c == '\\' {
		return true
	}
	return unicode.In(c, idStart...) && !unicode.In(c, notID...)
}

func IsIDContinue(c rune) bool {
	if c == '$' || c == '_' || c == '\\' || c == zwnj || c == zwj {
		return true
	}
	return unicode.In(c, idContinue...) && !unicode.In(c, notID...)
}
