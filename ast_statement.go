package javascript

import "vimagination.zapto.org/parser"

type StatementList struct {
	StatementListItems []StatementListItem
	Tokens             Tokens
}

func (j *jsParser) parseStatementList(yield, await, ret bool) (StatementList, error) {
	var sl StatementList
	for {
		g := j.NewGoal()
		g.AcceptRunWhitespace()
		h := g.NewGoal()
		si, err := h.parseStatementListItem(yield, await, ret)
		if err != nil {
			return sl, g.Error(err)
		}
		g.Score(h)
		j.Score(g)
		sl.StatementListItems = append(sl.StatementListItems, si)
	}
	sl.Tokens = j.ToTokens()
	return sl, nil
}

type Block StatementList

func (j *jsParser) parseBlock(yield, await, ret bool) (Block, error) {
	var b Block
	j.AcceptToken(parser.Token{TokenPunctuator, "{"})
	j.AcceptRunWhitespace()
	for {
		if j.Accept(TokenRightBracePunctuator) {
			break
		}
		g := j.NewGoal()
		si, err := g.parseStatementListItem(yield, await, ret)
		if err != nil {
			return b, j.Error(err)
		}
		j.Score(g)
		b.StatementListItems = append(b.StatementListItems, si)
		j.AcceptRunWhitespace()
	}
	b.Tokens = j.ToTokens()
	return b, nil
}

type StatementListItem struct {
	Statement   *Statement
	Declaration *Declaration
	Tokens      Tokens
}

func (j *jsParser) parseStatementListItem(yield, await, ret bool) (StatementListItem, error) {
	var si StatementListItem
	if err := j.FindGoal(
		func(j *jsParser) error {
			s, err := j.parseStatement(yield, await, ret)
			if err != nil {
				return err
			}
			si.Statement = &s
			return nil
		},
		func(j *jsParser) error {
			d, err := j.parseDeclaration(yield, ret)
			if err != nil {
				if err.(Error).Err == ErrInvalidDeclaration {
					return errNotApplicable
				}
				return err
			}
			si.Declaration = &d
			return nil
		},
	); err != nil {
		return si, err
	}
	if si.Statement == nil && si.Declaration == nil {
		return si, errNotApplicable
	}
	si.Tokens = j.ToTokens()
	return si, nil
}

type StatementType int

const (
	StatementNormal StatementType = iota
	StatementContinue
	StatementBreak
	StatementReturn
	StatementThrow
)

type Statement struct {
	Type                    StatementType
	BlockStatement          *StatementList
	VariableStatement       *VariableStatement
	ExpressionStatement     *Expression
	IfStatement             *IfStatement
	IterationStatementDo    *IterationStatementDo
	IterationStatementWhile *IterationStatementWhile
	IterationStatementFor   *IterationStatementFor
	SwitchStatement         *SwitchStatement
	ContinueStatement       *LabelIdentifier
	BreakStatement          *LabelIdentifier
	ReturnStatement         *Expression
	WithStatement           *WithStatement
	LabelIdentifier         *LabelIdentifier
	LabelledItemFunction    *FunctionDeclaration
	LabelledItemStatement   *Statement
	ThrowStatement          *Expression
	TryStatement            *TryStatement
	DebuggerStatement       *Token
	Tokens                  Tokens
}

func (j *jsParser) parseStatement(yield, await, ret bool) (Statement, error) {
	var s Statement
	g := j.NewGoal()
	switch g.Peek() {
	case parser.Token{TokenPunctuator, "{"}:
		g.Except()
		g.AcceptRunWhitespace()
		h := g.NewGoal()
		sl, err := h.parseStatementList(yield, await, ret)
		if err != nil {
			return s, g.Error(err)
		}
		g.Score(h)
		s.BlockStatement = &sl
	case parser.Token{TokenKeyword, "var"}:
		vs, err := g.parseVariableStatement(yield, await)
		if err != nil {
			return s, j.Error(err)
		}
		s.VariableStatement = &vs
	case parser.Token{TokenPunctuator, ";"}:
		g.Except()
	case parser.Token{TokenKeyword, "if"}:
		is, err := g.parseIfStatement(yield, await, ret)
		if err != nil {
			return s, j.Error(err)
		}
		s.IfStatement = &is
	case parser.Token{TokenKeyword, "do"}:
		ds, err := g.parseIterationStatementDo(yield, await, ret)
		if err != nil {
			return s, j.Error(err)
		}
		s.IterationStatementDo = &ds
	case parser.Token{TokenKeyword, "while"}:
		ws, err := g.parseIterationStatementWhile(yield, await, ret)
		if err != nil {
			return s, j.Error(err)
		}
		s.IterationStatementWhile = &ws
	case parser.Token{TokenKeyword, "for"}:
		fs, err := g.parseIterationStatementFor(yield, await, ret)
		if err != nil {
			return s, j.Error(err)
		}
		s.IterationStatementFor = &fs
	case parser.Token{TokenKeyword, "switch"}:
		ss, err := g.parseSwitchStatement(yield, await, ret)
		if err != nil {
			return s, j.Error(err)
		}
		s.SwitchStatement = &ss
	case parser.Token{TokenKeyword, "continue"}:
		g.Except()
		s.Type = StatementContinue
		g.AcceptRunWhitespaceNoNewLine()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			h := g.NewGoal()
			li, err := h.parseLabelIdentifier(yield, await)
			if err != nil {
				return s, g.Error(err)
			}
			s.ContinueStatement = &li
			g.Score(h)
			if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
				return s, g.Error(ErrMissingSemiColon)
			}
		}
	case parser.Token{TokenKeyword, "break"}:
		g.Except()
		s.Type = StatementBreak
		g.AcceptRunWhitespaceNoNewLine()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			h := g.NewGoal()
			li, err := h.parseLabelIdentifier(yield, await)
			if err != nil {
				return s, g.Error(err)
			}
			s.BreakStatement = &li
			g.Score(h)
			if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
				return s, g.Error(ErrMissingSemiColon)
			}
		}
	case parser.Token{TokenKeyword, "return"}:
		if !ret {
			return s, g.Error(ErrInvalidStatement)
		}
		g.Except()
		s.Type = StatementReturn
		g.AcceptRunWhitespaceNoNewLine()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			h := g.NewGoal()
			e, err := h.parseExpression(true, yield, await)
			if err != nil {
				return s, g.Error(err)
			}
			s.ReturnStatement = &e
			g.Score(h)
			if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
				return s, g.Error(ErrMissingSemiColon)
			}
		}
	case parser.Token{TokenKeyword, "with"}:
		ws, err := j.parseWithStatement(yield, await, ret)
		if err != nil {
			return s, j.Error(err)
		}
		s.WithStatement = &ws
	case parser.Token{TokenKeyword, "throw"}:
		g.Except()
		s.Type = StatementThrow
		g.AcceptRunWhitespaceNoNewLine()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			return s, g.Error(ErrMissingExpression)
		}
		h := g.NewGoal()
		e, err := h.parseExpression(true, yield, await)
		if err != nil {
			return s, g.Error(err)
		}
		s.ThrowStatement = &e
		g.Score(h)
		if !g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			return s, g.Error(ErrMissingSemiColon)
		}
	case parser.Token{TokenKeyword, "try"}:
		ts, err := g.parseTryStatement(yield, await, ret)
		if err != nil {
			return s, j.Error(err)
		}
		s.TryStatement = &ts
	case parser.Token{TokenKeyword, "debugger"}:
		g.Except()
		s.DebuggerStatement = g.GetLastToken()
	default:
		if err := g.FindGoal(
			func(j *jsParser) error {
				i, err := j.parseLabelIdentifier(yield, await)
				if err != nil {
					return err
				}
				j.AcceptRunWhitespace()
				if !j.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
					return ErrMissingColon
				}
				j.AcceptRunWhitespace()
				g := j.NewGoal()
				if g.Peek() == (parser.Token{TokenKeyword, "function"}) {
					fd, err := g.parseFunctionDeclaration(yield, await, false)
					if err != nil {
						return err
					}
					s.LabelledItemFunction = &fd
				} else {
					s, err := g.parseStatement(yield, await, ret)
					if err != nil {
						return err
					}
					s.LabelledItemStatement = &s
				}
				j.Score(g)
				s.LabelIdentifier = &i
				return nil
			},
			func(j *jsParser) error {
				e, err := j.parseExpression(true, yield, await)
				if err != nil {
					return err
				}
				s.ExpressionStatement = &e
				return nil
			},
		); err != nil {
			return s, err
		}
	}
	j.Score(g)
	s.Tokens = j.ToTokens()
	return s, nil
}

type IfStatement struct {
	Expression    Expression
	Statement     Statement
	ElseStatement *Statement
	Tokens        Tokens
}

func (j *jsParser) parseIfStatement(yield, await, ret bool) (IfStatement, error) {
	var (
		is  IfStatement
		err error
	)
	j.AcceptToken(parser.Token{TokenKeyword, "if"})
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return is, j.Error(ErrMissingOpeningParentheses)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	is.Expression, err = g.parseExpression(true, yield, await)
	if err != nil {
		return is, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return is, j.Error(ErrMissingClosingParentheses)
	}
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	is.Statement, err = g.parseStatement(yield, await, ret)
	if err != nil {
		return is, j.Error(err)
	}
	j.Score(g)
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenKeyword, "else"}) {
		h := g.NewGoal()
		s, err := h.parseStatement(yield, await, ret)
		if err != nil {
			return is, g.Error(err)
		}
		g.Score(h)
		j.Score(g)
		is.ElseStatement = &s
	}
	is.Tokens = j.ToTokens()
	return is, nil
}

type IterationStatementDo struct {
	Statement  Statement
	Expression Expression
	Tokens     Tokens
}

func (j *jsParser) parseIterationStatementDo(yield, await, ret bool) (IterationStatementDo, error) {
	var (
		is  IterationStatementDo
		err error
	)
	j.AcceptToken(parser.Token{TokenKeyword, "do"})
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	is.Statement, err = g.parseStatement(yield, await, ret)
	if err != nil {
		return is, j.Error(err)
	}
	j.Score(g)
	g = j.NewGoal()
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenKeyword, "while"}) {
		return is, j.Error(ErrInvalidIterationStatementDo)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return is, j.Error(ErrMissingOpeningParentheses)
	}
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	is.Expression, err = g.parseExpression(true, yield, await)
	if err != nil {
		return is, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return is, j.Error(ErrMissingClosingParentheses)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
		return is, j.Error(ErrMissingSemiColon)
	}
	is.Tokens = j.ToTokens()
	return is, nil
}

type IterationStatementWhile struct {
	Expression Expression
	Statement  Statement
	Tokens     Tokens
}

func (j *jsParser) parseIterationStatementWhile(yield, await, ret bool) (IterationStatementWhile, error) {
	var (
		is  IterationStatementWhile
		err error
	)
	j.AcceptToken(parser.Token{TokenKeyword, "while"})
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return is, j.Error(ErrMissingOpeningParentheses)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	is.Expression, err = g.parseExpression(true, await, ret)
	if err != nil {
		return is, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return is, j.Error(ErrMissingClosingParentheses)
	}
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	is.Statement, err = g.parseStatement(yield, await, ret)
	if err != nil {
		return is, j.Error(err)
	}
	j.Score(g)
	is.Tokens = j.ToTokens()
	return is, nil
}

type ForType int

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
)

type IterationStatementFor struct {
	Type ForType

	InitExpression *Expression
	InitVar        *VariableDeclaration
	InitLexical    *LexicalDeclaration
	Conditional    *Expression
	Afterthought   *Expression

	VariableDeclarationList *VariableDeclaration
	LeftHandSideExpression  *LeftHandSideExpression
	ForBindingIdentifier    *BindingIdentifier
	ForBindingPatternObject *ObjectBindingPattern
	ForBindingPatternArray  *ArrayBindingPattern
	In                      *Expression
	Of                      *AssignmentExpression

	Expression *Expression

	Statement Statement
	Tokens    Tokens
}

func (j *jsParser) parseIterationStatementFor(yield, await, ret bool) (IterationStatementFor, error) {
	var (
		is  IterationStatementFor
		err error
	)
	j.AcceptToken(parser.Token{TokenKeyword, "for"})
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return is, j.Error(ErrMissingOpeningParentheses)
	}
	j.AcceptRunWhitespace()

	if err = j.FindGoal(
		func(j *jsParser) error {
			if err = j.FindGoal(
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
					vd, err := g.parseVariableDeclaration(false, yield, await)
					if err != nil {
						return err
					}
					j.Score(g)
					is.InitVar = &vd
					is.Type = ForNormalVar
					return nil
				},
				func(j *jsParser) error {
					ld, err := j.parseLexicalDeclaration(false, yield, await)
					if err != nil {
						return err
					}
					is.InitLexical = &ld
					is.Type = ForNormalLexicalDeclaration
					return nil
				},
				func(j *jsParser) error {
					e, err := j.parseExpression(false, yield, await)
					if err != nil {
						return err
					}
					is.InitExpression = &e
					is.Type = ForNormalExpression
					return nil
				},
			); err != nil {
				return err
			}
			if j.GetLastToken().Token != (parser.Token{TokenPunctuator, ";"}) {
				j.AcceptRunWhitespace()
				if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
					is.InitVar = nil
					is.InitLexical = nil
					is.InitExpression = nil
					return ErrMissingSemiColon
				}
			}
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
				g := j.NewGoal()
				e, err := g.parseExpression(true, yield, await)
				if err != nil {
					return err
				}
				j.Score(g)
				is.Conditional = &e
				j.AcceptRunWhitespace()
				if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
					return ErrMissingSemiColon
				}
			}
			g := j.NewGoal()
			g.AcceptRunWhitespace()
			if g.Peek() != (parser.Token{TokenPunctuator, ")"}) {
				h := g.NewGoal()
				e, err := h.parseExpression(true, yield, await)
				if err != nil {
					return err
				}
				g.Score(h)
				is.Afterthought = &e
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
						ab, err := g.parseArrayBindingPattern(yield, await)
						if err != nil {
							return err
						}
						is.ForBindingPatternArray = &ab
					} else if tk == (parser.Token{TokenPunctuator, "{"}) {
						ob, err := g.parseObjectBindingPattern(yield, await)
						if err != nil {
							return err
						}
						is.ForBindingPatternObject = &ob
					} else {
						bi, err := g.parseBindingIdentifier(yield, await)
						if err != nil {
							return err
						}
						is.ForBindingIdentifier = &bi
					}
					j.Score(g)
					return nil
				},
				func(j *jsParser) error {
					g := j.NewGoal()
					lhs, err := g.parseLeftHandSideExpression(yield, await)
					if err != nil {
						return err
					}
					j.Score(g)
					is.LeftHandSideExpression = &lhs
					is.Type = ForInLeftHandSide
					return nil
				},
			); err != nil {
				return err
			}
			j.AcceptRunWhitespace()
			in := true
			if j.AcceptToken(parser.Token{TokenKeyword, "of"}) {
				in = false
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
			} else if !j.AcceptToken(parser.Token{TokenKeyword, "in"}) {
				return ErrInvalidForLoop
			}
			j.AcceptRunWhitespace()
			g := j.NewGoal()
			if in {
				e, err := g.parseExpression(true, yield, await)
				if err != nil {
					return err
				}
				is.In = &e
			} else {
				ae, err := j.parseAssignmentExpression(true, yield, await)
				if err != nil {
					return err
				}
				is.Of = &ae
			}
			j.Score(g)
			return nil
		},
	); err != nil {
		return is, err
	}

	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return is, j.Error(ErrMissingClosingParentheses)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	is.Statement, err = g.parseStatement(yield, await, ret)
	if err != nil {
		return is, j.Error(err)
	}
	j.Score(g)
	is.Tokens = j.ToTokens()
	return is, nil
}

type SwitchStatement struct {
	Expression             Expression
	CaseClauses            []CaseClause
	DefaultClause          *StatementList
	PostDefaultCaseClauses []CaseClause
	Tokens                 Tokens
}

func (j *jsParser) parseSwitchStatement(yield, await, ret bool) (SwitchStatement, error) {
	var (
		ss  SwitchStatement
		err error
	)
	j.AcceptToken(parser.Token{TokenKeyword, "switch"})
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return ss, j.Error(ErrMissingOpeningParentheses)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	ss.Expression, err = g.parseExpression(true, yield, await)
	if err != nil {
		return ss, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return ss, j.Error(ErrMissingClosingParentheses)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return ss, j.Error(ErrMissingOpeningBrace)
	}
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, "}"}) {
			break
		} else if j.AcceptToken(parser.Token{TokenKeyword, "default"}) {
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
				return ss, j.Error(ErrMissingColon)
			}
			g = j.NewGoal()
			sl, err := g.parseStatementList(yield, await, ret)
			if err != nil {
				return ss, j.Error(err)
			}
			j.Score(g)
			ss.DefaultClause = &sl
		} else {
			g := j.NewGoal()
			cc, err := g.parseCaseClause(yield, await, ret)
			if err != nil {
				return ss, j.Error(err)
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
	return ss, nil
}

type CaseClause struct {
	Expression    Expression
	StatementList *StatementList
	Tokens        Tokens
}

func (j *jsParser) parseCaseClause(yield, await, ret bool) (CaseClause, error) {
	var (
		cc  CaseClause
		err error
	)
	if !j.AcceptToken(parser.Token{TokenKeyword, "case"}) {
		return cc, ErrMissingCaseClause
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	cc.Expression, err = g.parseExpression(true, yield, await)
	if err != nil {
		return cc, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
		return cc, j.Error(ErrMissingColon)
	}
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if tk := g.Peek(); tk != (parser.Token{TokenKeyword, "case"}) && tk != (parser.Token{TokenKeyword, "default"}) && tk != (parser.Token{TokenPunctuator, "}"}) {
		h := g.NewGoal()
		sl, err := h.parseStatementList(yield, await, ret)
		if err != nil {
			return cc, g.Error(err)
		}
		g.Score(h)
		j.Score(g)
		cc.StatementList = &sl
	}
	cc.Tokens = j.ToTokens()
	return cc, nil
}

type WithStatement struct {
	Expression Expression
	Statement  Statement
	Tokens     Tokens
}

func (j *jsParser) parseWithStatement(yield, await, ret bool) (WithStatement, error) {
	var (
		ws  WithStatement
		err error
	)
	j.AcceptToken(parser.Token{TokenKeyword, "with"})
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return ws, j.Error(ErrMissingOpeningParentheses)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	ws.Expression, err = g.parseExpression(true, yield, await)
	if err != nil {
		return ws, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return ws, j.Error(ErrMissingClosingParentheses)
	}
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	ws.Statement, err = g.parseStatement(yield, await, ret)
	if err != nil {
		return ws, j.Error(err)
	}
	j.Score(g)
	ws.Tokens = j.ToTokens()
	return ws, nil
}

type TryStatement struct {
	TryBlock                           StatementList
	CatchParameterBindingIdentifier    *BindingIdentifier
	CatchParameterObjectBindingPattern *ObjectBindingPattern
	CatchParameterArrayBindingPattern  *ArrayBindingPattern
	CatchBlock                         *StatementList
	FinallyBlock                       *StatementList
	Tokens                             Tokens
}

func (j *jsParser) parseTryStatement(yield, await, ret bool) (TryStatement, error) {
	var (
		ts  TryStatement
		err error
	)
	j.AcceptToken(parser.Token{TokenKeyword, "try"})
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return ts, j.Error(ErrMissingOpeningBrace)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	ts.TryBlock, err = g.parseStatementList(yield, await, ret)
	if err != nil {
		return ts, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "}"}) {
		return ts, j.Error(ErrMissingClosingBrace)
	}
	j.AcceptRunWhitespace()
	if j.AcceptToken(parser.Token{TokenKeyword, "catch"}) {
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
			return ts, j.Error(ErrMissingOpeningParentheses)
		}
		j.AcceptRunWhitespace()
		g = j.NewGoal()
		switch g.Peek() {
		case parser.Token{TokenPunctuator, "{"}:
			ob, err := g.parseObjectBindingPattern(yield, await)
			if err != nil {
				return ts, j.Error(err)
			}
			ts.CatchParameterObjectBindingPattern = &ob
		case parser.Token{TokenPunctuator, "["}:
			ob, err := g.parseArrayBindingPattern(yield, await)
			if err != nil {
				return ts, j.Error(err)
			}
			ts.CatchParameterArrayBindingPattern = &ob
		default:
			bi, err := g.parseBindingIdentifier(yield, await)
			if err != nil {
				return ts, j.Error(err)
			}
			ts.CatchParameterBindingIdentifier = &bi
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
			return ts, j.Error(ErrMissingClosingParentheses)
		}
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
			return ts, j.Error(ErrMissingOpeningParentheses)
		}
		j.AcceptRunWhitespace()
		g = j.NewGoal()
		cb, err := g.parseStatementList(yield, await, ret)
		if err != nil {
			return ts, j.Error(err)
		}
		j.Score(g)
		ts.CatchBlock = &cb
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenPunctuator, "}"}) {
			return ts, j.Error(ErrMissingClosingBrace)
		}
	}
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenKeyword, "finally"}) {
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
			return ts, g.Error(ErrMissingOpeningBrace)
		}
		g.AcceptRunWhitespace()
		h := g.NewGoal()
		fb, err := h.parseStatementList(yield, await, ret)
		if err != nil {
			return ts, g.Error(err)
		}
		g.Score(h)
		ts.FinallyBlock = &fb
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "}"}) {
			return ts, g.Error(ErrMissingClosingBrace)
		}
		j.Score(g)
	}
	if ts.CatchBlock == nil && ts.FinallyBlock == nil {
		return ts, j.Error(ErrMissingCatchFinally)
	}
	ts.Tokens = j.ToTokens()
	return ts, nil
}

type VariableStatement struct {
	VariableDeclarationList []VariableDeclaration
	Tokens                  Tokens
}

func (j *jsParser) parseVariableStatement(yield, await bool) (VariableStatement, error) {
	var vs VariableStatement
	j.AcceptToken(parser.Token{TokenKeyword, "var"})
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			break
		}
		g := j.NewGoal()
		vd, err := g.parseVariableDeclaration(true, yield, await)
		if err != nil {
			return vs, j.Error(err)
		}
		j.Score(g)
		vs.VariableDeclarationList = append(vs.VariableDeclarationList, vd)
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return vs, j.Error(ErrMissingComma)
		}
	}
	vs.Tokens = j.ToTokens()
	return vs, nil
}
