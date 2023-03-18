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
