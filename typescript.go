package javascript

import "vimagination.zapto.org/parser"

const marker = "TS"

type typescript struct {
	Tokeniser
}

func (t *typescript) GetToken() (parser.Token, error) {
	tk, err := t.Tokeniser.GetToken()
	if tk.Type == parser.TokenDone {
		tk.Data = marker
	}
	return tk, err
}

func AsTypescript(t Tokeniser) Tokeniser {
	return &typescript{Tokeniser: t}
}

func (j *jsParser) IsTypescript() bool {
	return (*j)[:cap(*j)][cap(*j)-1].Data == marker
}

func (j *jsParser) ReadTypeParameters() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "<"}) {
		return false
	}
	for {
		g.AcceptRunWhitespace()
		if !g.ReadTypeParameter() {
			return false
		}
		g.AcceptRunWhitespace()
		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">"}) {
			break
		} else if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			return false
		}
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadTypeParameter() bool {
	g := j.NewGoal()
	if g.parseIdentifier(false, false) == nil {
		return false
	}
	h := g.NewGoal()
	h.AcceptRunWhitespace()
	if h.AcceptToken(parser.Token{Type: TokenKeyword, Data: "extends"}) {
		h.AcceptRunWhitespace()
		if !h.ReadType() {
			return false
		}
		g.Score(h)
		h = g.NewGoal()
		h.AcceptRunWhitespace()
	}
	if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
		h.AcceptRunWhitespace()
		if !h.ReadType() {
			return false
		}
		g.Score(h)
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadTypeArguments() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "<"}) {
		return false
	}
	for {
		g.AcceptRunWhitespace()
		if !g.ReadType() {
			return false
		}
		g.AcceptRunWhitespace()
		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">"}) {
			break
		} else if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			return false
		}
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadType() bool {
	if g := j.NewGoal(); g.ReadFunctionType() {
		j.Score(g)
		return true
	}
	g := j.NewGoal()
	if !g.ReadUnionOrIntersectionOrPrimaryType() {
		return false
	}
	h := g.NewGoal()
	h.AcceptRunWhitespaceNoNewLine()
	if h.AcceptToken(parser.Token{Type: TokenKeyword, Data: "extends"}) {
		h.AcceptRunWhitespace()
		if !h.ReadType() {
			return false
		}
		h.AcceptRunWhitespace()
		if !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) {
			return false
		}
		h.AcceptRunWhitespace()
		if !h.ReadType() {
			return false
		}
		h.AcceptRunWhitespace()
		if !h.ReadTypeAnnotation() {
			return false
		}
		g.Score(h)
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadUnionOrIntersectionOrPrimaryType() bool {
	g := j.NewGoal()
	for {
		if !g.ReadTypeOperator() {
			return false
		}
		h := g.NewGoal()
		h.AcceptRunWhitespace()
		if !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "|"}) && !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "&"}) {
			break
		}
		g.Score(h)
		g.AcceptRunWhitespace()
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadTypeOperator() bool {
	g := j.NewGoal()
	if tk := g.Peek(); tk.Type == TokenIdentifier && (tk.Data == "keyof" || tk.Data == "unique" || tk.Data == "readonly") {
		g.Skip()
		g.AcceptRunWhitespace()
		if !g.ReadTypeOperator() {
			return false
		}
	} else if tk == (parser.Token{Type: TokenIdentifier, Data: "infer"}) {
		g.Skip()
		g.AcceptRunWhitespace()
		if g.parseIdentifier(false, false) == nil {
			return false
		}
	} else if !g.ReadPostfixType() {
		return false
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadPostfixType() bool {
	g := j.NewGoal()
	if !g.ReadPrimaryType() {
		return false
	}
	for {
		h := g.NewGoal()
		h.AcceptRunWhitespaceNoNewLine()
		if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
			h.AcceptRunWhitespace()
			i := h.NewGoal()
			if i.ReadType() {
				h.Score(i)
				h.AcceptRunWhitespace()
			}
			if !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
				return false
			}
		} else if !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "!"}) {
			break
		}
		g.Score(h)
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadPrimaryType() bool {
	for _, fn := range [...]func(*jsParser) bool{
		(*jsParser).ReadLiteralType,
		(*jsParser).ReadTemplateType,
		(*jsParser).ReadParenthesizedType,
		(*jsParser).ReadPredefinedType,
		(*jsParser).ReadObjectOrMappedType,
		(*jsParser).ReadTupleType,
		(*jsParser).ReadThisType,
		(*jsParser).ReadTypeQuery,
		(*jsParser).ReadTypeReference,
	} {
		g := j.NewGoal()
		if fn(&g) {
			j.Score(g)
			return true
		}
	}
	return false
}

func (j *jsParser) ReadLiteralType() bool {
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "-"}) {
		g.AcceptRunWhitespace()
		if g.Accept(TokenNumericLiteral) {
			j.Score(g)
			return true
		}
	}
	return j.Accept(TokenNullLiteral, TokenBooleanLiteral, TokenNumericLiteral, TokenStringLiteral, TokenNoSubstitutionTemplate)
}

func (j *jsParser) ReadTemplateType() bool {
	g := j.NewGoal()
	if !g.Accept(TokenTemplateHead) {
		return false
	}
	for {
		g.AcceptRunWhitespace()
		if !g.ReadType() {
			return false
		}
		g.AcceptRunWhitespace()
		if g.Accept(TokenTemplateTail) {
			break
		}
		if !g.Accept(TokenTemplateMiddle) {
			return false
		}
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadParenthesizedType() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.ReadType() {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		return false
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadPredefinedType() bool {
	if j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "void"}) {
		return true
	}
	if tk := j.Peek(); tk.Type == TokenIdentifier {
		switch tk.Data {
		case "any", "number", "boolean", "string", "symbol", "unknown", "bigint", "undefined", "never", "object":
			j.Skip()
			g := j.NewGoal()
			g.AcceptRunWhitespace()
			if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "."}) {
				g.AcceptRunWhitespace()
				if !g.ReadTypeReference() {
					return false
				}
				j.Score(g)
			}
			return true
		}
	}
	return false
}

func (j *jsParser) ReadObjectOrMappedType() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return false
	}
	g.AcceptRunWhitespace()
	h := g.NewGoal()
	if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "+"}) || h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "-"}) {
		h.AcceptRunWhitespace()
		if h.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "readonly"}) {
			return j.ReadMappedType()
		}
	}
	h = g.NewGoal()
	if h.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "readonly"}) {
		h.AcceptRunWhitespace()
	}
	if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		h.AcceptRunWhitespace()
		if h.parseIdentifier(false, false) != nil {
			h.AcceptRunWhitespace()
			if h.AcceptToken(parser.Token{Type: TokenKeyword, Data: "in"}) {
				return j.ReadMappedType()
			}
		}
	}
	return j.ReadObjectType()
}

func (j *jsParser) ReadMappedType() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "readonly"}) && (g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "+"}) || g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "-"})) {
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "readonly"}) {
			return false
		}
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		return false
	}
	if g.parseIdentifier(false, false) == nil {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "in"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.ReadType() {
		return false
	}
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "as"}) {
		g.AcceptRunWhitespace()
		if !g.ReadType() {
			return false
		}
		g.AcceptRunWhitespace()
	}
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) && g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "+"}) || g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "-"}) {
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) {
			return false
		}
		g.AcceptRunWhitespace()
	}
	if g.ReadTypeAnnotation() {
		g.AcceptRunWhitespace()
	}
	if !g.Accept(TokenRightBracePunctuator) {
		for {
			if !g.ReadTypeMember() {
				return false
			}
			g.AcceptRunWhitespace()
			sep := g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) || g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","})
			g.AcceptRunWhitespace()
			if g.Accept(TokenRightBracePunctuator) {
				break
			}
			if !sep {
				return false
			}
		}
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadObjectType() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.Accept(TokenRightBracePunctuator) {
		for {
			if !g.ReadTypeMember() {
				return false
			}
			g.AcceptRunWhitespace()
			sep := g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) || g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","})
			g.AcceptRunWhitespace()
			if g.Accept(TokenRightBracePunctuator) {
				break
			}
			if !sep {
				return false
			}
		}
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadTypeMember() bool {
	for _, fn := range [...]func(*jsParser) bool{
		(*jsParser).ReadCallSignature,
		(*jsParser).ReadConstructSignature,
		(*jsParser).ReadAccessorDeclaration,
		(*jsParser).ReadIndexSignature,
		(*jsParser).ReadMethodSignature,
		(*jsParser).ReadPropertySignature,
	} {
		g := j.NewGoal()
		if fn(&g) {
			j.Score(g)
			return true
		}
	}
	return false
}

func (j *jsParser) ReadAccessorDeclaration() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "get"}) && !g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "set"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.ReadPropertyName() {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.ReadCallSignature() {
		return false
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadPropertyName() bool {
	if j.Accept(TokenStringLiteral, TokenNumericLiteral, TokenPrivateIdentifier) || j.parseIdentifier(false, false) != nil {
		return true
	}
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		return false
	}
	g.AcceptRunWhitespace()
	var e Expression
	if e.parse(&g, false, false, false) != nil {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
		return false
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadPropertySignature() bool {
	g := j.NewGoal()
	if !g.Accept(TokenIdentifier, TokenStringLiteral, TokenNumericLiteral) {
		return false
	}
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) {
		g.AcceptRunWhitespace()
	}
	g.ReadTypeAnnotation()
	j.Score(g)
	return true
}

func (j *jsParser) ReadTypeAnnotation() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.ReadType() {
		return false
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadColonReturnType() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.ReadReturnType() {
		return false
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadReturnType() bool {
	g := j.NewGoal()
	h := g.NewGoal()
	if h.parseIdentifier(false, false) != nil {
		h.AcceptRunWhitespaceNoNewLine()
		if h.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "is"}) {
			h.AcceptRunWhitespace()
			g.Score(h)
		}
	}
	if !g.ReadType() {
		return false
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadCallSignature() bool {
	g := j.NewGoal()
	if g.ReadTypeParameters() {
		g.AcceptRunWhitespace()
	}
	if !g.ReadParameterList() {
		return false
	}
	h := g.NewGoal()
	h.AcceptRunWhitespace()
	if h.ReadColonReturnType() {
		g.Score(h)
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadParameterList() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		optional := false
		first := true
		for {
			g.AcceptRunWhitespace()
			if g.Peek() == (parser.Token{Type: TokenPunctuator, Data: "..."}) {
				g.Skip()
				g.AcceptRunWhitespace()
				if g.parseIdentifier(false, false) == nil {
					return false
				}
				g.AcceptRunWhitespace()
				g.ReadTypeAnnotation()
				g.AcceptRunWhitespace()
				if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
					return false
				}
				break
			}
			if bi := g.parseIdentifier(false, false); bi == nil {
				if !first || !g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "this"}) {
					return false
				}
			} else if bi.Data == "public" || bi.Data == "private" || bi.Data == "protected" {
				g.AcceptRunWhitespace()
				g.parseIdentifier(false, false)
			}
			g.AcceptRunWhitespace()
			opt := g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"})
			g.AcceptRunWhitespace()
			g.ReadTypeAnnotation()
			g.AcceptRunWhitespace()
			init := g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="})
			if optional && (!opt || !init) {
				return false
			}
			if init {
				g.AcceptRunWhitespace()
				var ae AssignmentExpression
				if err := ae.parse(&g, false, false, false); err != nil {
					return false
				}
				g.AcceptRunWhitespace()
			}
			optional = optional || opt || init
			if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
				break
			}
			if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
				return false
			}
			first = false
		}
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadConstructSignature() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "new"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.ReadCallSignature() {
		return false
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadIndexSignature() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.ReadParameter() {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
		return false
	}
	h := g.NewGoal()
	h.AcceptRunWhitespace()
	if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) {
		g.Score(h)
		h = g.NewGoal()
		h.AcceptRunWhitespace()
	}
	if h.ReadTypeAnnotation() {
		g.Score(h)
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadParameter() bool {
	g := j.NewGoal()
	var seenConst, seenStatic bool
	for {
		if !seenConst && g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "const"}) {
			seenConst = true
		} else if !seenStatic && g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "static"}) {
			seenStatic = true
		} else {
			break
		}
		g.AcceptRunWhitespace()
	}
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "this"}) {
		j.Score(g)
		return true
	}
	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
		g.AcceptRunWhitespace()
	}
	if g.parseIdentifier(false, false) == nil {
		return false
	}
	h := g.NewGoal()
	h.AcceptRunWhitespace()
	if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) {
		g.Score(h)
		h = g.NewGoal()
		h.AcceptRunWhitespace()
	}
	if h.ReadTypeAnnotation() {
		g.Score(h)
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadMethodSignature() bool {
	g := j.NewGoal()
	if !g.Accept(TokenIdentifier, TokenStringLiteral, TokenNumericLiteral) {
		return false
	}
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) {
		g.AcceptRunWhitespace()
	}
	if !g.ReadCallSignature() {
		return false
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadTupleType() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
		for {
			if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
				g.AcceptRunWhitespace()
				if !g.ReadType() {
					return false
				}
				g.AcceptRunWhitespace()
				break
			}
			if !g.ReadType() {
				return false
			}
			g.AcceptRunWhitespace()
			if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
				break
			}
			g.AcceptRunWhitespace()
		}
		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			return false
		}
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadThisType() bool {
	return j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "this"})
}

func (j *jsParser) ReadTypeQuery() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "typeof"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if g.parseIdentifier(false, false) == nil {
		return false
	}
	for {
		h := g.NewGoal()
		h.AcceptRunWhitespace()
		if !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "."}) {
			break
		}
		h.AcceptRunWhitespace()
		if !h.Accept(TokenIdentifier, TokenKeyword, TokenPrivateIdentifier) {
			return false
		}
		g.Score(h)
	}
	h := g.NewGoal()
	h.AcceptRunWhitespaceNoNewLine()
	if h.ReadTypeParameters() {
		g.Score(h)
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadTypeReference() bool {
	g := j.NewGoal()
	if g.parseIdentifier(false, false) == nil {
		return false
	}
	for {
		h := g.NewGoal()
		h.AcceptRunWhitespace()
		if !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "."}) {
			break
		}
		h.AcceptRunWhitespace()
		if h.parseIdentifier(false, false) == nil {
			return false
		}
		g.Score(h)
	}
	h := g.NewGoal()
	h.AcceptRunWhitespaceNoNewLine()
	if h.ReadTypeArguments() {
		g.Score(h)
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadFunctionType() bool {
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "new"}) {
		g.AcceptRunWhitespace()
	}
	if g.ReadTypeParameters() {
		g.AcceptRunWhitespace()
	}
	if !g.ReadParameterList() {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "=>"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.ReadReturnType() {
		return false
	}
	j.Score(g)
	return true
}

func (j *jsParser) SkipHeritage() bool {
	if j.IsTypescript() && j.Peek() == (parser.Token{Type: TokenIdentifier, Data: "implements"}) {
		g := j.NewGoal()
		if g.ReadHeritage() {
			j.Score(g)
			return true
		}
	}
	return false
}

func (j *jsParser) SkipGeneric() bool {
	return j.IsTypescript() && j.ReadTypeParameters()
}

func (j *jsParser) SkipAsType() bool {
	if j.IsTypescript() {
		g := j.NewGoal()
		if !g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "satisfies"}) && !g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "as"}) {
			return false
		}
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "const"}) && !g.ReadType() {
			return false
		}
		j.Score(g)
		return true
	}
	return false
}

func (j *jsParser) SkipColonType() bool {
	if j.IsTypescript() {
		g := j.NewGoal()
		if g.ReadTypeAnnotation() {
			j.Score(g)
			return true
		}
	}
	return false
}

func (j *jsParser) SkipReturnType() bool {
	if j.IsTypescript() {
		g := j.NewGoal()
		if g.ReadColonReturnType() {
			j.Score(g)
			return true
		}
	}
	return false
}

func (j *jsParser) SkipOptionalColonType() bool {
	ret := false
	if j.IsTypescript() {
		g := j.NewGoal()
		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) {
			j.Score(g)
			g = j.NewGoal()
			g.AcceptRunWhitespace()
			ret = true
		}
		if g.ReadTypeAnnotation() {
			j.Score(g)
			ret = true
		}
	}
	return ret
}

func (j *jsParser) SkipType() bool {
	return j.IsTypescript() && j.ReadTypeDeclaration()
}

func (j *jsParser) ReadTypeDeclaration() bool {
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "type"}) {
		g.AcceptRunWhitespace()
		if g.parseIdentifier(false, false) == nil {
			return false
		}
		g.AcceptRunWhitespace()
		if g.ReadTypeParameters() {
			g.AcceptRunWhitespace()
		}
		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
			return false
		}
		g.AcceptRunWhitespace()
		// TODO: instrinsic check?
		if !g.ReadType() {
			return false
		}
		h := g.NewGoal()
		h.AcceptRunWhitespace()
		if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
			g.Score(h)
		}
		j.Score(g)
		return true
	}
	return false
}

func (j *jsParser) ReadHeritage() bool {
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "extends"}) || g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "implements"}) {
		for {
			g.AcceptRunWhitespace()
			var lhs LeftHandSideExpression
			if lhs.parse(&g, false, false) != nil {
				return false
			}
			g.AcceptRunWhitespace()
			g.ReadTypeParameters()
			g.AcceptRunWhitespace()
			if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
				break
			}
		}
		j.Score(g)
		return true
	}
	return false
}

func (j *jsParser) SkipInterface() bool {
	if j.IsTypescript() {
		g := j.NewGoal()
		if g.ReadInterface() {
			j.Score(g)
			return true
		}
	}
	return false
}

func (j *jsParser) ReadInterface() bool {
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "interface"}) {
		g.AcceptRunWhitespace()
		if g.parseIdentifier(false, false) == nil {
			return false
		}
		g.AcceptRunWhitespace()
		if g.ReadTypeParameters() {
			g.AcceptRunWhitespace()
		}
		if g.ReadHeritage() {
			g.AcceptRunWhitespace()
			return false
		}
		if !g.ReadObjectType() {
			return false
		}
		j.Score(g)
		return true
	}
	return false
}

func (j *jsParser) SkipParameterProperties() bool {
	if j.IsTypescript() {
		if tk := j.Peek(); tk == (parser.Token{Type: TokenIdentifier, Data: "private"}) || tk == (parser.Token{Type: TokenIdentifier, Data: "protected"}) || tk == (parser.Token{Type: TokenIdentifier, Data: "public"}) {
			g := j.NewGoal()
			g.Skip()
			g.AcceptRunWhitespaceNoNewLine()
			if tk := g.Peek(); tk.Type != TokenLineTerminator && tk != (parser.Token{Type: TokenPunctuator, Data: ";"}) {
				j.Score(g)
				return true
			}
		}
	}
	return false
}

func (j *jsParser) SkipTypeArguments() bool {
	return j.IsTypescript() && j.ReadTypeArguments()
}

func (j *jsParser) SkipReadOnly() {
	if j.IsTypescript() {
		if tk := j.Peek(); tk == (parser.Token{Type: TokenIdentifier, Data: "readonly"}) {
			g := j.NewGoal()
			g.Skip()
			g.AcceptRunWhitespaceNoNewLine()
			if tk := g.Peek(); tk.Type != TokenLineTerminator && tk != (parser.Token{Type: TokenPunctuator, Data: ";"}) {
				j.Score(g)
			}
		}
	}
}

func (j *jsParser) SkipImportType() bool {
	if j.IsTypescript() && j.Peek() == (parser.Token{Type: TokenKeyword, Data: "import"}) {
		g := j.NewGoal()
		g.Skip()
		g.AcceptRunWhitespace()
		if g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "type"}) {
			g.AcceptRunWhitespace()
			if tk := g.Peek(); tk != (parser.Token{Type: TokenPunctuator, Data: ","}) && tk != (parser.Token{Type: TokenIdentifier, Data: "from"}) {
				var ic ImportClause
				err := ic.parse(&g)
				if err == nil {
					g.AcceptRunWhitespace()
					var fc FromClause
					err := fc.parse(&g)
					if err == nil {
						h := g.NewGoal()
						h.AcceptRunWhitespace()
						if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
							g.Score(h)
						}
						j.Score(g)
						return true
					}
				}
			}
		}
	}
	return false
}

func (j *jsParser) SkipExportType() bool {
	if j.IsTypescript() && j.Peek() == (parser.Token{Type: TokenKeyword, Data: "export"}) {
		g := j.NewGoal()
		g.Skip()
		g.AcceptRunWhitespace()
		if g.ReadTypeDeclaration() {
			j.Score(g)
			return true
		} else if g.ReadInterface() {
			j.Score(g)
			return true
		} else if g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "type"}) {
			g.AcceptRunWhitespace()
			if tk := g.Peek(); tk != (parser.Token{Type: TokenPunctuator, Data: ","}) && tk != (parser.Token{Type: TokenIdentifier, Data: "from"}) {
				var ec ExportClause
				err := ec.parse(&g)
				if err == nil {
					h := g.NewGoal()
					h.AcceptRunWhitespace()
					if h.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "from"}) {
						h.AcceptRunWhitespace()
						if h.Accept(TokenStringLiteral) {
							g.Score(h)
						}
					}
					h = g.NewGoal()
					h.AcceptRunWhitespace()
					if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
						g.Score(h)
					}
					j.Score(g)
					return true
				}
			}
		}
	}
	return false
}

func (j *jsParser) SkipTypeImport() bool {
	g := j.NewGoal()
	if g.IsTypescript() && g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "type"}) {
		h := g.NewGoal()
		h.AcceptRunWhitespace()
		if h.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "as"}) {
			i := h.NewGoal()
			i.AcceptRunWhitespace()
			if i.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "as"}) {
				i.AcceptRunWhitespace()
				if i.parseIdentifier(false, false) != nil {
					h.Score(i)
					g.Score(h)
					j.Score(g)
					return true
				}
			} else {
				g.Score(h)
				j.Score(g)
				return true
			}
		} else if h.Accept(TokenIdentifier, TokenKeyword) {
			i := h.NewGoal()
			i.AcceptRunWhitespace()
			if i.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "as"}) {
				i.AcceptRunWhitespace()
				if i.parseIdentifier(false, false) != nil {
					h.Score(h)
					g.Score(h)
					j.Score(g)
					return true
				}
			} else {
				g.Score(h)
				j.Score(g)
				return true
			}
		}
	}
	return false
}

func (j *jsParser) SkipThisParam() bool {
	g := j.NewGoal()
	if g.IsTypescript() && g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "this"}) {
		g.AcceptRunWhitespace()
		g.SkipColonType()
		j.Score(g)
		return true
	}
	return false
}

func (j *jsParser) SkipDeclare() bool {
	g := j.NewGoal()
	if g.IsTypescript() && g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "declare"}) {
		g.AcceptRunWhitespace()
		switch g.Peek() {
		case parser.Token{Type: TokenKeyword, Data: "var"}:
			var vd VariableStatement
			if vd.parse(&g, false, false) == nil {
				return true
			}
		case parser.Token{Type: TokenKeyword, Data: "const"}, parser.Token{Type: TokenIdentifier, Data: "let"}:
			var ld LexicalDeclaration
			if ld.parse(&g, true, false, false) == nil {
				return true
			}
		case parser.Token{Type: TokenKeyword, Data: "async"}:
			g.Skip()
			g.AcceptRunWhitespace()
			fallthrough
		case parser.Token{Type: TokenKeyword, Data: "function"}:
			if g.ReadFunctionDeclaration() {
				return true
			}
		case parser.Token{Type: TokenKeyword, Data: "class"}:
			if g.ReadClassDeclaration() {
				return true
			}
		}
	}
	return false
}

func (j *jsParser) SkipForce() bool {
	return j.IsTypescript() && j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "!"})
}

func (j *jsParser) SkipAbstract() bool {
	return j.IsTypescript() && j.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "abstract"})
}

func (j *jsParser) SkipAbstractField() bool {
	if j.IsTypescript() {
		g := j.NewGoal()
		if g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "abstract"}) {
			g.AcceptRunWhitespace()
			if g.ReadTypeMember() {
				j.Score(g)
				return true
			}
		}
	}
	return false
}

func (j *jsParser) ReadFunctionDeclaration() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "function"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "*"}) {
		g.AcceptRunWhitespace()
	}
	if g.parseIdentifier(false, false) == nil {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.ReadTypeParameters() {
		return false
	}
	h := g.NewGoal()
	h.AcceptRunWhitespace()
	if h.ReadColonReturnType() {
		g.Score(h)
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadClassDeclaration() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "class"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if g.parseIdentifier(false, false) == nil {
		return false
	}
	g.AcceptRunWhitespace()
	for g.ReadHeritage() {
		g.AcceptRunWhitespace()
	}
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return false
	}
	g.AcceptRunWhitespace()
	for g.ReadTypeMember() {
		g.AcceptRunWhitespace()
	}
	if !g.Accept(TokenRightBracePunctuator) {
		return false
	}
	j.Score(g)
	return true
}

func (j *jsParser) OnOptionalType() bool {
	if j.IsTypescript() {
		g := j.NewGoal()
		g.AcceptRunWhitespace()
		if g.Accept(TokenPunctuator) {
			tk := g.GetLastToken()
			return tk.Data == ")" || tk.Data == ":"
		}
	}
	return false
}
