package javascript

import (
	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

type Block struct {
	StatementListItems []StatementListItem
	Tokens             Tokens
}

func (b *Block) parse(j *jsParser, yield, await, ret bool) error {
	j.AcceptToken(parser.Token{TokenPunctuator, "{"})
	j.AcceptRunWhitespace()
	for {
		if j.Accept(TokenRightBracePunctuator) {
			break
		}
		g := j.NewGoal()
		var si StatementListItem
		if err := si.parse(&g, yield, await, ret); err != nil {
			si.clear()
			return j.Error("Block", err)
		}
		j.Score(g)
		b.StatementListItems = append(b.StatementListItems, si)
		j.AcceptRunWhitespace()
	}
	b.Tokens = j.ToTokens()
	return nil
}

type StatementListItem struct {
	Statement   *Statement
	Declaration *Declaration
	Tokens      Tokens
}

func (si *StatementListItem) parse(j *jsParser, yield, await, ret bool) error {
	if err := j.FindGoal(
		func(j *jsParser) error {
			g := j.NewGoal()
			d := newDeclaration()
			if err := d.parse(&g, yield, ret); err != nil {
				d.clear()
				poolDeclaration.Put(d)
				return j.Error("StatementListItem.Declaration", err)
			}
			j.Score(g)
			si.Declaration = d
			return nil
		},
		func(j *jsParser) error {
			g := j.NewGoal()
			s := newStatement()
			if err := s.parse(&g, yield, await, ret); err != nil {
				s.clear()
				poolStatement.Put(s)
				return j.Error("StatementListItem.Statement", err)
			}
			j.Score(g)
			si.Statement = s
			return nil
		},
	); err != nil {
		return err
	}
	si.Tokens = j.ToTokens()
	return nil
}

type StatementType uint8

const (
	StatementNormal StatementType = iota
	StatementContinue
	StatementBreak
	StatementReturn
	StatementThrow
)

type Statement struct {
	Type                    StatementType
	BlockStatement          *Block
	VariableStatement       *VariableStatement
	ExpressionStatement     *Expression
	IfStatement             *IfStatement
	IterationStatementDo    *IterationStatementDo
	IterationStatementWhile *IterationStatementWhile
	IterationStatementFor   *IterationStatementFor
	SwitchStatement         *SwitchStatement
	WithStatement           *WithStatement
	LabelIdentifier         *Token
	LabelledItemFunction    *FunctionDeclaration
	LabelledItemStatement   *Statement
	TryStatement            *TryStatement
	DebuggerStatement       *Token
	Tokens                  Tokens
}

func (s *Statement) parse(j *jsParser, yield, await, ret bool) error {
	g := j.NewGoal()
	switch g.Peek() {
	case parser.Token{TokenPunctuator, "{"}:
		s.BlockStatement = newBlock()
		if err := s.BlockStatement.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{TokenKeyword, "var"}:
		s.VariableStatement = newVariableStatement()
		if err := s.VariableStatement.parse(&g, yield, await); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{TokenPunctuator, ";"}:
		g.Except()
	case parser.Token{TokenKeyword, "if"}:
		s.IfStatement = newIfStatement()
		if err := s.IfStatement.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{TokenKeyword, "do"}:
		s.IterationStatementDo = newIterationStatementDo()
		if err := s.IterationStatementDo.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{TokenKeyword, "while"}:
		s.IterationStatementWhile = newIterationStatementWhile()
		if err := s.IterationStatementWhile.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{TokenKeyword, "for"}:
		s.IterationStatementFor = newIterationStatementFor()
		if err := s.IterationStatementFor.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{TokenKeyword, "switch"}:
		s.SwitchStatement = newSwitchStatement()
		if err := s.SwitchStatement.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{TokenKeyword, "continue"}:
		g.Except()
		s.Type = StatementContinue
		g.AcceptRunWhitespaceNoNewLine()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			h := g.NewGoal()
			li, err := h.parseIdentifier(yield, await)
			if err != nil {
				return g.Error("Statement", err)
			}
			s.LabelIdentifier = li
			g.Score(h)
			if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
				return g.Error("Statement", ErrMissingSemiColon)
			}
		}
	case parser.Token{TokenKeyword, "break"}:
		g.Except()
		s.Type = StatementBreak
		g.AcceptRunWhitespaceNoNewLine()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			h := g.NewGoal()
			li, err := h.parseIdentifier(yield, await)
			if err != nil {
				return g.Error("Statement", err)
			}
			s.LabelIdentifier = li
			g.Score(h)
			if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
				return g.Error("Statement", ErrMissingSemiColon)
			}
		}
	case parser.Token{TokenKeyword, "return"}:
		if !ret {
			return g.Error("Statement", ErrInvalidStatement)
		}
		g.Except()
		s.Type = StatementReturn
		g.AcceptRunWhitespaceNoNewLine()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			h := g.NewGoal()
			s.ExpressionStatement = newExpression()
			if err := s.ExpressionStatement.parse(&h, true, yield, await); err != nil {
				return g.Error("Statement", err)
			}
			g.Score(h)
			if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
				return g.Error("Statement", ErrMissingSemiColon)
			}
		}
	case parser.Token{TokenKeyword, "with"}:
		s.WithStatement = newWithStatement()
		if err := s.WithStatement.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{TokenKeyword, "throw"}:
		g.Except()
		s.Type = StatementThrow
		g.AcceptRunWhitespaceNoNewLine()
		h := g.NewGoal()
		s.ExpressionStatement = newExpression()
		if err := s.ExpressionStatement.parse(&h, true, yield, await); err != nil {
			return g.Error("Statement", err)
		}
		g.Score(h)
		if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			return g.Error("Statement", ErrMissingSemiColon)
		}
	case parser.Token{TokenKeyword, "try"}:
		s.TryStatement = newTryStatement()
		if err := s.TryStatement.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{TokenKeyword, "debugger"}:
		g.Except()
		s.DebuggerStatement = g.GetLastToken()
		h := g.NewGoal()
		h.AcceptRunWhitespace()
		if h.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			g.Score(h)
		}
	default:
		if err := g.FindGoal(
			func(j *jsParser) error {
				i, err := j.parseIdentifier(yield, await)
				if err != nil {
					return j.Error("Statement.LabelIdentifier", err)
				}
				j.AcceptRunWhitespace()
				if !j.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
					return j.Error("Statement.LabelIdentifier", ErrMissingColon)
				}
				j.AcceptRunWhitespace()
				g := j.NewGoal()
				if g.Peek() == (parser.Token{TokenKeyword, "function"}) {
					fd := newFunctionDeclaration()
					if err := fd.parse(&g, yield, await, false); err != nil {
						fd.clear()
						poolFunctionDeclaration.Put(fd)
						return j.Error("Statement.LabelIdentifier", err)
					}
					s.LabelledItemFunction = fd
				} else {
					ss := newStatement()
					if err := ss.parse(&g, yield, await, ret); err != nil {
						ss.clear()
						poolStatement.Put(ss)
						return j.Error("Statement.LabelIdentifier", err)
					}
					s.LabelledItemStatement = ss
				}
				j.Score(g)
				s.LabelIdentifier = i
				return nil
			},
			func(j *jsParser) error {
				switch j.Peek() {
				case parser.Token{TokenPunctuator, "{"}, parser.Token{TokenKeyword, "function"}, parser.Token{TokenKeyword, "class"}:
					return errNotApplicable
				case parser.Token{TokenIdentifier, "async"}:
					g := j.NewGoal()
					g.Except()
					g.AcceptRunWhitespaceNoNewLine()
					if g.AcceptToken(parser.Token{TokenKeyword, "function"}) {
						return j.Error("Statement.Expression", errNotApplicable)
					}
				}
				e := newExpression()
				if err := e.parse(j, true, yield, await); err != nil {
					e.clear()
					poolExpression.Put(e)
					return j.Error("Statement.Expression", err)
				}
				j.AcceptRunWhitespace()
				if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
					e.clear()
					poolExpression.Put(e)
					return j.Error("Statement.Expression", ErrMissingSemiColon)
				}
				s.ExpressionStatement = e
				return nil
			},
		); err != nil {
			return err
		}
	}
	j.Score(g)
	s.Tokens = j.ToTokens()
	return nil
}

type IfStatement struct {
	Expression    Expression
	Statement     Statement
	ElseStatement *Statement
	Tokens        Tokens
}

func (is *IfStatement) parse(j *jsParser, yield, await, ret bool) error {
	j.AcceptToken(parser.Token{TokenKeyword, "if"})
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return j.Error("IfStatement", ErrMissingOpeningParenthesis)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	if err := is.Expression.parse(&g, true, yield, await); err != nil {
		return j.Error("IfStatement", err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return j.Error("IfStatement", ErrMissingClosingParenthesis)
	}
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	if err := is.Statement.parse(&g, yield, await, ret); err != nil {
		return j.Error("IfStatement", err)
	}
	j.Score(g)
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenKeyword, "else"}) {
		g.AcceptRunWhitespace()
		h := g.NewGoal()
		is.ElseStatement = newStatement()
		if err := is.ElseStatement.parse(&h, yield, await, ret); err != nil {
			return g.Error("IfStatement", err)
		}
		g.Score(h)
		j.Score(g)
	}
	is.Tokens = j.ToTokens()
	return nil
}

type IterationStatementDo struct {
	Statement  Statement
	Expression Expression
	Tokens     Tokens
}

func (is *IterationStatementDo) parse(j *jsParser, yield, await, ret bool) error {
	j.AcceptToken(parser.Token{TokenKeyword, "do"})
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	if err := is.Statement.parse(&g, yield, await, ret); err != nil {
		return j.Error("IterationStatementDo", err)
	}
	j.Score(g)
	g = j.NewGoal()
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenKeyword, "while"}) {
		return j.Error("IterationStatementDo", ErrInvalidIterationStatementDo)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return j.Error("IterationStatementDo", ErrMissingOpeningParenthesis)
	}
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	if err := is.Expression.parse(&g, true, yield, await); err != nil {
		return j.Error("IterationStatementDo", err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return j.Error("IterationStatementDo", ErrMissingClosingParenthesis)
	}
	g = j.NewGoal()
	if g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
		j.Score(g)
	}
	is.Tokens = j.ToTokens()
	return nil
}

type IterationStatementWhile struct {
	Expression Expression
	Statement  Statement
	Tokens     Tokens
}

func (is *IterationStatementWhile) parse(j *jsParser, yield, await, ret bool) error {
	j.AcceptToken(parser.Token{TokenKeyword, "while"})
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return j.Error("IterationStatementWhile", ErrMissingOpeningParenthesis)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	if err := is.Expression.parse(&g, true, await, ret); err != nil {
		return j.Error("IterationStatementWhile", err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return j.Error("IterationStatementWhile", ErrMissingClosingParenthesis)
	}
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	if err := is.Statement.parse(&g, yield, await, ret); err != nil {
		return j.Error("IterationStatementWhile", err)
	}
	j.Score(g)
	is.Tokens = j.ToTokens()
	return nil
}

type ForType uint8

const (
	ForNormal ForType = iota
	ForNormalVar
	ForNormalLexicalDeclaration
	ForNormalExpression
	ForInLeftHandSide
	ForInVar
	ForInLet
	ForInConst
	ForOfLeftHandSide
	ForOfVar
	ForOfLet
	ForOfConst
	ForAwaitOfLeftHandSide
	ForAwaitOfVar
	ForAwaitOfLet
	ForAwaitOfConst
)

type IterationStatementFor struct {
	Type ForType

	InitExpression *Expression
	InitVar        *VariableDeclaration
	InitLexical    *LexicalDeclaration
	Conditional    *Expression
	Afterthought   *Expression

	LeftHandSideExpression  *LeftHandSideExpression
	ForBindingIdentifier    *Token
	ForBindingPatternObject *ObjectBindingPattern
	ForBindingPatternArray  *ArrayBindingPattern
	In                      *Expression
	Of                      *AssignmentExpression

	Statement Statement
	Tokens    Tokens
}

func (is *IterationStatementFor) parse(j *jsParser, yield, await, ret bool) error {
	j.AcceptToken(parser.Token{TokenKeyword, "for"})
	j.AcceptRunWhitespace()
	forAwait := j.AcceptToken(parser.Token{TokenKeyword, "await"})
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return j.Error("IterationStatementFor", ErrMissingOpeningParenthesis)
	}
	j.AcceptRunWhitespace()
	if err := j.FindGoal(
		func(j *jsParser) error {
			if forAwait {
				return errNotApplicable
			}
			if err := j.FindGoal(
				func(j *jsParser) error {
					if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
						return errNotApplicable
					}
					is.Type = ForNormal
					return nil
				},
				func(j *jsParser) error {
					if !j.AcceptToken(parser.Token{TokenKeyword, "var"}) {
						return errNotApplicable
					}
					j.AcceptRunWhitespace()
					g := j.NewGoal()
					vd := newVariableDeclaration()
					if err := vd.parse(&g, false, yield, await); err != nil {
						vd.clear()
						poolVariableDeclaration.Put(vd)
						return j.Error("IterationStatementFor.NormalVar", err)
					}
					j.Score(g)
					is.InitVar = vd
					is.Type = ForNormalVar
					return nil
				},
				func(j *jsParser) error {
					g := j.NewGoal()
					ld := newLexicalDeclaration()
					if err := ld.parse(&g, false, yield, await); err != nil {
						ld.clear()
						poolLexicalDeclaration.Put(ld)
						return j.Error("IterationStatementFor.NormalLexicalDeclaration", err)
					}
					j.Score(g)
					is.InitLexical = ld
					is.Type = ForNormalLexicalDeclaration
					return nil
				},
				func(j *jsParser) error {
					g := j.NewGoal()
					e := newExpression()
					if err := e.parse(&g, false, yield, await); err != nil {
						e.clear()
						poolExpression.Put(e)
						return j.Error("IterationStatementFor.NormalExpression", err)
					}
					j.Score(g)
					is.InitExpression = e
					is.Type = ForNormalExpression
					return nil
				},
			); err != nil {
				return err
			}
			if j.GetLastToken().Token != (parser.Token{TokenPunctuator, ";"}) {
				j.AcceptRunWhitespace()
				if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
					is.clear()
					return j.Error("IterationStatementFor.Normal", ErrMissingSemiColon)
				}
			}
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
				g := j.NewGoal()
				is.Conditional = newExpression()
				if err := is.Conditional.parse(&g, true, yield, await); err != nil {
					is.clear()
					return j.Error("IterationStatementFor.Normal", err)
				}
				j.Score(g)
				j.AcceptRunWhitespace()
				if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
					return j.Error("IterationStatementFor.Normal", ErrMissingSemiColon)
				}
			}
			g := j.NewGoal()
			g.AcceptRunWhitespace()
			if g.Peek() != (parser.Token{TokenPunctuator, ")"}) {
				h := g.NewGoal()
				is.Afterthought = newExpression()
				if err := is.Afterthought.parse(&h, true, yield, await); err != nil {
					is.clear()
					return j.Error("IterationStatementFor.Normal", err)
				}
				g.Score(h)
			}
			j.Score(g)
			return nil
		},
		func(j *jsParser) error {
			if err := j.FindGoal(
				func(j *jsParser) error {
					if j.AcceptToken(parser.Token{TokenKeyword, "var"}) {
						is.Type = ForInVar
					} else if j.AcceptToken(parser.Token{TokenKeyword, "const"}) {
						is.Type = ForInConst
					} else if j.AcceptToken(parser.Token{TokenIdentifier, "let"}) {
						is.Type = ForInLet
					} else {
						return errNotApplicable
					}
					j.AcceptRunWhitespace()
					g := j.NewGoal()
					if tk := g.Peek(); tk == (parser.Token{TokenPunctuator, "["}) {
						is.ForBindingPatternArray = newArrayBindingPattern()
						if err := is.ForBindingPatternArray.parse(&g, yield, await); err != nil {
							is.clear()
							return j.Error("IterationStatementFor.Other", err)
						}
					} else if tk == (parser.Token{TokenPunctuator, "{"}) {
						is.ForBindingPatternObject = newObjectBindingPattern()
						if err := is.ForBindingPatternObject.parse(&g, yield, await); err != nil {
							is.clear()
							return j.Error("IterationStatementFor.Other", err)
						}
					} else {
						bi, err := g.parseIdentifier(yield, await)
						if err != nil {
							is.clear()
							return j.Error("IterationStatementFor.Other", err)
						}
						is.ForBindingIdentifier = bi
					}
					j.Score(g)
					return nil
				},
				func(j *jsParser) error {
					g := j.NewGoal()
					is.LeftHandSideExpression = newLeftHandSideExpression()
					if err := is.LeftHandSideExpression.parse(&g, yield, await); err != nil {
						is.clear()
						return j.Error("IterationStatementFor.LeftHandSideExpression", err)
					}
					j.Score(g)
					is.Type = ForInLeftHandSide
					return nil
				},
			); err != nil {
				return err
			}
			j.AcceptRunWhitespace()
			in := true
			if j.AcceptToken(parser.Token{TokenIdentifier, "of"}) {
				in = false
				if forAwait {
					switch is.Type {
					case ForInVar:
						is.Type = ForAwaitOfVar
					case ForInConst:
						is.Type = ForAwaitOfConst
					case ForInLet:
						is.Type = ForAwaitOfLet
					case ForInLeftHandSide:
						is.Type = ForAwaitOfLeftHandSide
					}
				} else {
					switch is.Type {
					case ForInVar:
						is.Type = ForOfVar
					case ForInConst:
						is.Type = ForOfConst
					case ForInLet:
						is.Type = ForOfLet
					case ForInLeftHandSide:
						is.Type = ForOfLeftHandSide
					}
				}
			} else if forAwait || !j.AcceptToken(parser.Token{TokenKeyword, "in"}) {
				is.clear()
				return j.Error("IterationStatementFor.Other", ErrInvalidForLoop)
			}
			j.AcceptRunWhitespace()
			g := j.NewGoal()
			if in {
				is.In = newExpression()
				if err := is.In.parse(&g, true, yield, await); err != nil {
					is.clear()
					return j.Error("IterationStatementFor.Other", err)
				}
			} else {
				is.Of = newAssignmentExpression()
				if err := is.Of.parse(&g, true, yield, await); err != nil {
					is.clear()
					return j.Error("IterationStatementFor.Other", err)
				}
			}
			j.Score(g)
			return nil
		},
	); err != nil {
		return err
	}

	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return j.Error("IterationStatementFor", ErrMissingClosingParenthesis)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()

	if err := is.Statement.parse(&g, yield, await, ret); err != nil {
		return j.Error("IterationStatementFor", err)
	}
	j.Score(g)
	is.Tokens = j.ToTokens()
	return nil
}

type SwitchStatement struct {
	Expression             Expression
	CaseClauses            []CaseClause
	DefaultClause          []StatementListItem
	PostDefaultCaseClauses []CaseClause
	Tokens                 Tokens
}

func (ss *SwitchStatement) parse(j *jsParser, yield, await, ret bool) error {
	j.AcceptToken(parser.Token{TokenKeyword, "switch"})
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return j.Error("SwitchStatement", ErrMissingOpeningParenthesis)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	if err := ss.Expression.parse(&g, true, yield, await); err != nil {
		return j.Error("SwitchStatement", err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return j.Error("SwitchStatement", ErrMissingClosingParenthesis)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return j.Error("SwitchStatement", ErrMissingOpeningBrace)
	}
	for {
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		} else if j.AcceptToken(parser.Token{TokenKeyword, "default"}) {
			if ss.DefaultClause != nil {
				return j.Error("SwitchStatement", ErrDuplicateDefaultClause)
			}
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
				return j.Error("SwitchStatement", ErrMissingColon)
			}
			ss.DefaultClause = []StatementListItem{}
			for {
				j.AcceptRunWhitespace()
				if pt := j.Peek(); pt == (parser.Token{TokenKeyword, "case"}) || pt.Type == TokenRightBracePunctuator {
					break
				}
				g = j.NewGoal()
				var sl StatementListItem
				if err := sl.parse(&g, yield, await, ret); err != nil {
					sl.clear()
					return j.Error("SwitchStatement", err)
				}
				j.Score(g)
				ss.DefaultClause = append(ss.DefaultClause, sl)
			}
		} else {
			g := j.NewGoal()
			var cc CaseClause
			if err := cc.parse(&g, yield, await, ret); err != nil {
				cc.clear()
				return j.Error("SwitchStatement", err)
			}
			j.Score(g)
			if ss.DefaultClause == nil {
				ss.CaseClauses = append(ss.CaseClauses, cc)
			} else {
				ss.PostDefaultCaseClauses = append(ss.PostDefaultCaseClauses, cc)
			}
		}
	}
	ss.Tokens = j.ToTokens()
	return nil
}

type CaseClause struct {
	Expression    Expression
	StatementList []StatementListItem
	Tokens        Tokens
}

func (cc *CaseClause) parse(j *jsParser, yield, await, ret bool) error {
	if !j.AcceptToken(parser.Token{TokenKeyword, "case"}) {
		return j.Error("CaseClause", ErrMissingCaseClause)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	if err := cc.Expression.parse(&g, true, yield, await); err != nil {
		return j.Error("CaseClause", err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
		return j.Error("CaseClause", ErrMissingColon)
	}
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	for {
		if tk := g.Peek(); tk == (parser.Token{TokenKeyword, "case"}) || tk == (parser.Token{TokenKeyword, "default"}) || tk == (parser.Token{TokenRightBracePunctuator, "}"}) {
			break
		}
		h := g.NewGoal()
		var sl StatementListItem
		if err := sl.parse(&h, yield, await, ret); err != nil {
			sl.clear()
			return g.Error("CaseClause", err)
		}
		g.Score(h)
		cc.StatementList = append(cc.StatementList, sl)
	}
	j.Score(g)
	cc.Tokens = j.ToTokens()
	return nil
}

type WithStatement struct {
	Expression Expression
	Statement  Statement
	Tokens     Tokens
}

func (ws *WithStatement) parse(j *jsParser, yield, await, ret bool) error {
	j.AcceptToken(parser.Token{TokenKeyword, "with"})
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return j.Error("WithStatement", ErrMissingOpeningParenthesis)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	if err := ws.Expression.parse(&g, true, yield, await); err != nil {
		return j.Error("WithStatement", err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return j.Error("WithStatement", ErrMissingClosingParenthesis)
	}
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	if err := ws.Statement.parse(&g, yield, await, ret); err != nil {
		return j.Error("WithStatement", err)
	}
	j.Score(g)
	ws.Tokens = j.ToTokens()
	return nil
}

type TryStatement struct {
	TryBlock                           Block
	CatchParameterBindingIdentifier    *Token
	CatchParameterObjectBindingPattern *ObjectBindingPattern
	CatchParameterArrayBindingPattern  *ArrayBindingPattern
	CatchBlock                         *Block
	FinallyBlock                       *Block
	Tokens                             Tokens
}

func (ts *TryStatement) parse(j *jsParser, yield, await, ret bool) error {
	j.AcceptToken(parser.Token{TokenKeyword, "try"})
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	if err := ts.TryBlock.parse(&g, yield, await, ret); err != nil {
		return j.Error("TryStatement", err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if j.AcceptToken(parser.Token{TokenKeyword, "catch"}) {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
			j.AcceptRunWhitespace()
			g = j.NewGoal()
			switch g.Peek() {
			case parser.Token{TokenPunctuator, "{"}:
				ts.CatchParameterObjectBindingPattern = newObjectBindingPattern()
				if err := ts.CatchParameterObjectBindingPattern.parse(&g, yield, await); err != nil {
					return j.Error("TryStatement", err)
				}
			case parser.Token{TokenPunctuator, "["}:
				ts.CatchParameterArrayBindingPattern = newArrayBindingPattern()
				if err := ts.CatchParameterArrayBindingPattern.parse(&g, yield, await); err != nil {
					return j.Error("TryStatement", err)
				}
			default:
				bi, err := g.parseIdentifier(yield, await)
				if err != nil {
					return j.Error("TryStatement", err)
				}
				ts.CatchParameterBindingIdentifier = bi
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
				return j.Error("TryStatement", ErrMissingClosingParenthesis)
			}
			j.AcceptRunWhitespace()
		}
		g = j.NewGoal()
		ts.CatchBlock = newBlock()
		if err := ts.CatchBlock.parse(&g, yield, await, ret); err != nil {
			return j.Error("TryStatement", err)
		}
		j.Score(g)
	}
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenKeyword, "finally"}) {
		g.AcceptRunWhitespace()
		h := g.NewGoal()
		ts.FinallyBlock = newBlock()
		if err := ts.FinallyBlock.parse(&h, yield, await, ret); err != nil {
			return g.Error("TryStatement", err)
		}
		g.Score(h)
		j.Score(g)
	}
	if ts.CatchBlock == nil && ts.FinallyBlock == nil {
		return j.Error("TryStatement", ErrMissingCatchFinally)
	}
	ts.Tokens = j.ToTokens()
	return nil
}

type VariableStatement struct {
	VariableDeclarationList []VariableDeclaration
	Tokens                  Tokens
}

func (vs *VariableStatement) parse(j *jsParser, yield, await bool) error {
	j.AcceptToken(parser.Token{TokenKeyword, "var"})
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			break
		}
		g := j.NewGoal()
		var vd VariableDeclaration
		if err := vd.parse(&g, true, yield, await); err != nil {
			vd.clear()
			return j.Error("VariableStatement", err)
		}
		j.Score(g)
		vs.VariableDeclarationList = append(vs.VariableDeclarationList, vd)
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return j.Error("VariableStatement", ErrMissingComma)
		}
	}
	vs.Tokens = j.ToTokens()
	return nil
}

var (
	ErrDuplicateDefaultClause = errors.New("duplicate default clause")
)
