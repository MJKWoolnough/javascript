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

/*
ClassDeclaration (<>, implements)
AssignmentExpression (!, as)
FieldDefinition (private, protected)
MethodDefinition (<>)
FormalParameters (:TYPE)
FunctionDeclaration (<>, :TYPE)
ArrowFunction (<>, :TYPE)
StatementListItem (enum, type, interface)
LexicalBinding (!:TYPE)
TryStatement (:TYPE)
ModuleItem (import type)
LeftHandSideExpression (<>)
*/

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
	return false
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
	for _, fn := range [...]func(*jsParser) bool{
		(*jsParser).ReadUnionOrIntersectionOrPrimaryType,
		(*jsParser).ReadFunctionType,
		(*jsParser).ReadConstructorType,
	} {
		g := j.NewGoal()
		if fn(&g) {
			j.Score(g)
			return true
		}
	}
	return false
}

func (j *jsParser) ReadUnionOrIntersectionOrPrimaryType() bool {
	g := j.NewGoal()
	for {
		if !g.ReadPrimaryType() {
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

func (j *jsParser) ReadPrimaryType() bool {
	for _, fn := range [...]func(*jsParser) bool{
		(*jsParser).ReadParenthesizedType,
		(*jsParser).ReadPredefinedType,
		(*jsParser).ReadObjectType,
		(*jsParser).ReadArrayType,
		(*jsParser).ReadTupleType,
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
		case "any", "number", "boolean", "string", "symbol", "unknown":
			j.Skip()
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
			if g.Accept(TokenRightBracePunctuator) {
				break
			}
			if sep {
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
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) {
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
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return false
	}
	g.AcceptRunWhitespace()
	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		if !g.ReadParameterList() {
			return false
		}
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
			return false
		}
	}
	g.AcceptRunWhitespace()
	g.ReadTypeAnnotation()
	j.Score(g)
	return false
}

func (j *jsParser) ReadParameterList() bool {
	return false
}

func (j *jsParser) ReadConstructSignature() bool {
	return false
}

func (j *jsParser) ReadIndexSignature() bool {
	return false
}

func (j *jsParser) ReadMethodSignature() bool {
	return false
}

func (j *jsParser) ReadArrayType() bool {
	return false
}

func (j *jsParser) ReadTupleType() bool {
	return false
}

func (j *jsParser) ReadTypeQuery() bool {
	return false
}

func (j *jsParser) ReadTypeReference() bool {
	return false
}

func (j *jsParser) ReadFunctionType() bool {
	return false
}

func (j *jsParser) ReadConstructorType() bool {
	return false
}

func (j *jsParser) SkipGeneric() {}

func (j *jsParser) SkipAsType() {}

func (j *jsParser) SkipColonType() {}

func (j *jsParser) SkipType() {}

func (j *jsParser) SkipInterface() {}

func (j *jsParser) SkipEnum() {}

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

func (j *jsParser) SkipImportType() {
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
						j.AcceptRunWhitespace()
						j.Score(g)
					}
				}
			}
		}
	}
}
