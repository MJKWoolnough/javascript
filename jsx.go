package javascript

import (
	"strings"

	"vimagination.zapto.org/parser"
)

const jsxMarker = "X"

type jsx struct {
	Tokeniser
}

func (j *jsx) Iter(fn func(parser.Token) bool) {
	for tk := range j.Tokeniser.Iter {
		if tk.Type == parser.TokenDone {
			tk.Data += jsxMarker
		}

		if !fn(tk) {
			break
		}
	}
}

func AsJSX(t Tokeniser) Tokeniser {
	return &jsx{Tokeniser: t}
}

func (j *jsParser) IsJSX() bool {
	return strings.HasSuffix((*j)[:cap(*j)][cap(*j)-1].Data, jsxMarker)
}
