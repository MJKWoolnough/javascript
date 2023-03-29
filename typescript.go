package javascript

import "vimagination.zapto.org/parser"

const marker = "TS"

func ParseTypescript(t Tokeniser) (*Script, error) {
	j, err := newJSParser(t)
	if err != nil {
		return nil, err
	}
	j[len(j)-1].Data = marker
	s := new(Script)
	if err := s.parse(&j); err != nil {
		return nil, err
	}
	return s, nil
}

func ParseTypescriptModule(t Tokeniser) (*Module, error) {
	j, err := newJSParser(t)
	if err != nil {
		return nil, err
	}
	j[:cap(j)][cap(j)-1].Data = marker
	m := new(Module)
	if err := m.parse(&j); err != nil {
		return nil, err
	}
	return m, nil
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
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "extends"}) {
		g.AcceptRunWhitespace()
		if !g.ReadType() {
			return false
		}
	}
	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
		g.AcceptRunWhitespace()
		if !g.ReadType() {
			return false
		}
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
	h.AcceptRunWhitespace()
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
		if !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
			return false
		}
		h.AcceptRunWhitespace()
		if !h.ReadType() {
			return false
		}
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
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "|"}) && !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "&"}) {
			break
		}
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
		(*jsParser).ReadObjectType,
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
		if !g.ReadType() {
			return false
		}
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
		(*jsParser).ReadMethodSignature,
		(*jsParser).ReadConstructSignature,
		(*jsParser).ReadIndexSignature,
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

func (j *jsParser) ReadCallSignature() bool {
	g := j.NewGoal()
	g.ReadTypeParameters()
	g.AcceptRunWhitespace()
	if !g.ReadParameterList() {
		return false
	}
	g.AcceptRunWhitespace()
	g.ReadTypeAnnotation()
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
				return false
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
	g.ReadTypeParameters()
	if !g.ReadParameterList() {
		return false
	}
	g.AcceptRunWhitespace()
	g.ReadTypeAnnotation()
	j.Score(g)
	return true
}

func (j *jsParser) ReadIndexSignature() bool {
	g := j.NewGoal()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		return false
	}
	g.AcceptRunWhitespace()
	if g.parseIdentifier(false, false) == nil {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "string"}) && !g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "number"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.ReadTypeAnnotation() {
		return false
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
	if !g.ReadType() {
		return false
	}
	g.AcceptRunWhitespace()
	for g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
		g.AcceptRunWhitespace()
		if !g.ReadType() {
			return false
		}
		g.AcceptRunWhitespace()
	}
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
		return false
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
	if !g.ReadTypeQueryExpression() {
		return false
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadTypeQueryExpression() bool {
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
		if !h.Accept(TokenIdentifier, TokenKeyword, TokenPrivateIdentifier) {
			return false
		}
		g.Score(h)
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadTypeReference() bool {
	g := j.NewGoal()
	if !g.ReadTypeName() {
		return false
	}
	h := g.NewGoal()
	h.AcceptRunWhitespaceNoNewLine()
	if h.ReadTypeArguments() {
		g.Score(h)
	}
	j.Score(g)
	return true
}

func (j *jsParser) ReadTypeName() bool {
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
	if !g.ReadType() {
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

func (j *jsParser) SkipOptionalColonType() bool {
	if j.IsTypescript() {
		g := j.NewGoal()
		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) {
			g.AcceptRunWhitespace()
		}
		if g.ReadTypeAnnotation() {
			j.Score(g)
			return true
		}
	}
	return false
}

func (j *jsParser) SkipType() bool {
	if j.IsTypescript() {
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
	}
	return false
}

func (j *jsParser) ParseEnum() {
}

func (j *jsParser) SkipParameterProperties() {
	if j.IsTypescript() {
		if tk := j.Peek(); tk == (parser.Token{Type: TokenIdentifier, Data: "private"}) || tk == (parser.Token{Type: TokenIdentifier, Data: "protected"}) || tk == (parser.Token{Type: TokenIdentifier, Data: "public"}) {
			g := j.NewGoal()
			g.Skip()
			g.AcceptRunWhitespaceNoNewLine()
			if tk := g.Peek(); tk.Type != TokenLineTerminator && tk != (parser.Token{Type: TokenPunctuator, Data: ";"}) {
				j.Score(g)
			}
		}
	}
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
	if h.ReadTypeAnnotation() {
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
