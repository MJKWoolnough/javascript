package javascript

import (
	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

type Token interface {
}

type Script StatementList

func ParseScript(t parser.Tokeniser) (Script, error) {
	j, err := newJSParser(t)
	if err != nil {
		return Script{}, err
	}
	s, err := j.parseStatementList(false, false, false)
	if err != nil {
		return Script{}, err
	}
	return Script(s), nil
}

type StatementList struct {
	Tokens []TokenPos
}

func (j *jsParser) parseStatementList(yield, await, ret bool) (StatementList, error) {
	var sl StatementList
	return sl, nil
}

type StatementListItem struct {
	Tokens []TokenPos
}

func (j *jsParser) parseStatementListItem(yield, await, ret bool) (StatementListItem, error) {
	var si StatementListItem
	return si, nil
}

type IdentifierReference Identifier

func (j *jsParser) parseIdentifierReference(yield, await bool) (IdentifierReference, error) {
	i, err := j.parseIdentifier(yield, await)
	return IdentifierReference(i), err
}

type BindingIdentifier Identifier

func (j *jsParser) parseBindingIdentifier(yield, await bool) (BindingIdentifier, error) {
	i, err := j.parseIdentifier(yield, await)
	return BindingIdentifier(i), err
}

type LabelIdentifier Identifier

func (j *jsParser) parseLabelIdentifier(yield, await bool) (LabelIdentifier, error) {
	i, err := j.parseIdentifier(yield, await)
	return LabelIdentifier(i), err
}

type Identifier struct {
	Identifier *TokenPos
}

func (j *jsParser) parseIdentifier(yield, await bool) (Identifier, error) {
	if j.Accept(TokenIdentifier) || (!yield && j.AcceptToken(parser.Token{TokenKeyword, "yield"}) || !await && j.AcceptToken(parser.Token{TokenKeyword, "await"})) {
		return Identifier{j.GetLastToken()}, nil
	}
	return Identifier{}, j.Error(ErrNoIdentifier)
}

type Declaration struct {
	ClassDeclaration    *ClassDeclaration
	FunctionDeclaration *FunctionDeclaration
	LexicalDeclaration  *LexicalDeclaration
	Tokens              []TokenPos
}

func (j *jsParser) parseDeclaration(yield, await bool) (Declaration, error) {
	var d Declaration
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{TokenKeyword, "class"}) {
		cd, err := g.parseClassDeclaration(yield, await, false)
		if err != nil {
			return d, j.Error(err)
		}
		d.ClassDeclaration = &cd
	} else if tk := g.Peek(); tk == (parser.Token{TokenKeyword, "const"}) || tk == (parser.Token{TokenIdentifier, "let"}) {
		ld, err := g.parseLexicalDeclaration(true, yield, await)
		if err != nil {
			return d, j.Error(err)
		}
		d.LexicalDeclaration = &ld
	} else if tk == (parser.Token{TokenKeyword, "async"}) || tk == (parser.Token{TokenKeyword, "function"}) {
		fd, err := g.parseFunctionDeclaration(yield, await, false)
		if err != nil {
			return d, j.Error(err)
		}
		d.FunctionDeclaration = &fd
	} else {
		return d, j.Error(ErrInvalidDeclaration)
	}
	j.Score(g)
	d.Tokens = j.ToTokens()
	return d, nil
}

type LetOrConst bool

const (
	Let   LetOrConst = false
	Const LetOrConst = true
)

type LexicalDeclaration struct {
	LetOrConst
	BindingList []LexicalBinding
	Tokens      []TokenPos
}

func (j *jsParser) parseLexicalDeclaration(in, yield, await bool) (LexicalDeclaration, error) {
	var ld LexicalDeclaration
	if !j.AcceptToken(parser.Token{TokenIdentifier, "let"}) {
		if !j.AcceptToken(parser.Token{TokenKeyword, "const"}) {
			return ld, j.Error(ErrInvalidLexicalDeclaration)
		}
		ld.LetOrConst = Const
	}
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			break
		}
		g := j.NewGoal()
		lb, err := g.parseLexicalBinding(in, yield, await)
		if err != nil {
			return ld, j.Error(err)
		}
		j.Score(g)
		ld.BindingList = append(ld.BindingList, lb)
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return ld, j.Error(ErrInvalidLexicalDeclaration)
		}
	}
	ld.Tokens = j.ToTokens()
	return ld, nil
}

type LexicalBinding struct {
	BindingIdentifier    *BindingIdentifier
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Initializer          *AssignmentExpression
	Tokens               []TokenPos
}

func (j *jsParser) parseLexicalBinding(in, yield, await bool) (LexicalBinding, error) {
	var lb LexicalBinding
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{TokenPunctuator, "["}) {
		ab, err := g.parseArrayBindingPattern(yield, await)
		if err != nil {
			return lb, j.Error(err)
		}
		lb.ArrayBindingPattern = &ab
	} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		ob, err := g.parseObjectBindingPattern(yield, await)
		if err != nil {
			return lb, j.Error(err)
		}
		lb.ObjectBindingPattern = &ob
	} else {
		bi, err := g.parseBindingIdentifier(yield, await)
		if err != nil {
			return lb, j.Error(err)
		}
		lb.BindingIdentifier = &bi
	}
	j.Score(g)
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenPunctuator, "="}) {
		g.AcceptRunWhitespace()
		j.Score(g)
		g = j.NewGoal()
		ae, err := g.parseAssignmentExpression(in, yield, await)
		if err != nil {
			return lb, j.Error(err)
		}
		j.Score(g)
		lb.Initializer = &ae
	}
	lb.Tokens = j.ToTokens()
	return lb, nil
}

type ArrayBindingPattern struct {
	BindingElementList []BindingElement
	BindingRestElement *BindingElement
	Token              []TokenPos
}

func (j *jsParser) parseArrayBindingPattern(yield, await bool) (ArrayBindingPattern, error) {
	var ab ArrayBindingPattern
	j.AcceptToken(parser.Token{TokenPunctuator, "["})
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
			break
		} else if j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			ab.BindingElementList = append(ab.BindingElementList, BindingElement{})
			continue
		}
		g := j.NewGoal()
		rest := g.AcceptToken(parser.Token{TokenPunctuator, "..."})
		g.AcceptRunWhitespace()
		be, err := g.parseBindingElement(yield, await)
		if err != nil {
			return ab, j.Error(err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		if rest {
			ab.BindingRestElement = &be
			if !j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
				return ab, j.Error(ErrMissingClosingBracket)
			}
			break
		}
		ab.BindingElementList = append(ab.BindingElementList, be)
		if j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return ab, j.Error(ErrMissingComma)
		}
	}
	ab.Token = j.ToTokens()
	return ab, nil
}

type ObjectBindingPattern struct {
	BindingPropertyList []BindingProperty
	Token               []TokenPos
}

func (j *jsParser) parseObjectBindingPattern(yield, await bool) (ObjectBindingPattern, error) {
	var ob ObjectBindingPattern
	j.AcceptToken(parser.Token{TokenPunctuator, "{"})
	for {
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		}
		g := j.NewGoal()
		bp, err := g.parseBindingProperty(yield, await)
		if err != nil {
			return ob, j.Error(err)
		}
		j.Score(g)
		ob.BindingPropertyList = append(ob.BindingPropertyList, bp)
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return ob, j.Error(ErrMissingComma)
		}
	}
	return ob, nil
}

type BindingProperty struct {
	SingleNameBinding *BindingIdentifier
	Initializer       *AssignmentExpression
	PropertyName      *PropertyName
	BindingElement    *BindingElement
	Tokens            []TokenPos
}

func (j *jsParser) parseBindingProperty(yield, await bool) (BindingProperty, error) {
	var bp BindingProperty
	g := j.NewGoal()
	pn, err := g.parsePropertyName(yield, await)
	if err == nil {
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
			err = j.Error(ErrMissingColon)
			g = j.NewGoal()
		} else {
			g.AcceptRunWhitespace()
			var be BindingElement
			h := g.NewGoal()
			be, err = h.parseBindingElement(yield, await)
			if err == nil {
				bp.PropertyName = &pn
				bp.BindingElement = &be
			}
		}
	}
	if err != nil {
		bi, errr := g.parseBindingIdentifier(yield, await)
		if errr != nil {
			return bp, j.Error(farthestError(err, errr))
		}
		g.AcceptRunWhitespace()
		if g.AcceptToken(parser.Token{TokenPunctuator, "="}) {
			g.AcceptRunWhitespace()
			h := g.NewGoal()
			i, err := h.parseAssignmentExpression(true, yield, await)
			if err != nil {
				return bp, g.Error(err)
			}
			g.Score(h)
			bp.Initializer = &i
		}
		bp.SingleNameBinding = &bi
	}
	j.Score(g)
	bp.Tokens = j.ToTokens()
	return bp, nil
}

type VariableStatement struct {
	VariableDeclarationList []VariableDeclaration
	Tokens                  []TokenPos
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

type VariableDeclaration LexicalBinding

func (j *jsParser) parseVariableDeclaration(in, yield, await bool) (VariableDeclaration, error) {
	lb, err := j.parseLexicalBinding(in, yield, await)
	return VariableDeclaration(lb), err
}

type AssignmentOperator uint8

const (
	AssignmentNone AssignmentOperator = iota
	AssignmentAssign
	AssignmentMultiply
	AssignmentDivide
	AssignmentRemainder
	AssignmentAdd
	AssignmentSubtract
	AssignmentLeftShift
	AssignmentSignPropagatinRightShift
	AssignmentZeroFillRightShift
	AssignmentBitwiseAND
	AssignmentBitwiseXOR
	AssignmentBitwiseOR
	AssignmentExponentiation
)

func (j *jsParser) parseAssignmentOperator() (AssignmentOperator, error) {
	var ao AssignmentOperator
	switch j.Peek() {
	case parser.Token{TokenPunctuator, "="}:
		ao = AssignmentAssign
	case parser.Token{TokenPunctuator, "*="}:
		ao = AssignmentMultiply
	case parser.Token{TokenPunctuator, "/="}:
		ao = AssignmentDivide
	case parser.Token{TokenPunctuator, "%="}:
		ao = AssignmentRemainder
	case parser.Token{TokenPunctuator, "+="}:
		ao = AssignmentAdd
	case parser.Token{TokenPunctuator, "-="}:
		ao = AssignmentSubtract
	case parser.Token{TokenPunctuator, "<<="}:
		ao = AssignmentLeftShift
	case parser.Token{TokenPunctuator, ">>="}:
		ao = AssignmentSignPropagatinRightShift
	case parser.Token{TokenPunctuator, ">>>="}:
		ao = AssignmentZeroFillRightShift
	case parser.Token{TokenPunctuator, "&="}:
		ao = AssignmentBitwiseAND
	case parser.Token{TokenPunctuator, "^="}:
		ao = AssignmentBitwiseXOR
	case parser.Token{TokenPunctuator, "|="}:
		ao = AssignmentBitwiseOR
	case parser.Token{TokenPunctuator, "**="}:
		ao = AssignmentExponentiation
	default:
		return 0, ErrInvalidAssignment
	}
	j.Except()
	return ao, nil
}

type AssignmentExpression struct {
	ConditionalExpression  *ConditionalExpression
	ArrowFunction          *ArrowFunction
	LeftHandSideExpression *LeftHandSideExpression
	Yield                  bool
	Delegate               bool
	AssignmentOperator     AssignmentOperator
	AssignmentExpression   *AssignmentExpression
	Tokens                 []TokenPos
}

func (j *jsParser) parseAssignmentExpression(in, yield, await bool) (AssignmentExpression, error) {
	var ae AssignmentExpression
	if err := j.findGoal(func(j *jsParser) error {
		if !yield || !j.AcceptToken(parser.Token{TokenKeyword, "yield"}) {
			return errNotApplicable
		}
		ae.Yield = true
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
			ae.Delegate = true
			j.AcceptRunWhitespace()
		}
		g := j.NewGoal()
		nae, err := g.parseAssignmentExpression(in, true, await)
		if err != nil {
			return err
		}
		j.Score(g)
		ae.AssignmentExpression = &nae
		return nil
	}, func(j *jsParser) error {
		af, err := j.parseArrowFunction(in, yield, await)
		if err != nil {
			return err
		}
		ae.ArrowFunction = &af
		return nil
	}, func(j *jsParser) error {
		ce, err := j.parseConditionalExpression(in, yield, await)
		if err != nil {
			return err
		}
		ae.ConditionalExpression = &ce
		return nil
	}, func(j *jsParser) error {
		lhs, err := j.parseLeftHandSideExpression(yield, await)
		if err != nil {
			return err
		}
		ae.LeftHandSideExpression = &lhs
		j.AcceptRunWhitespace()
		ae.AssignmentOperator, err = j.parseAssignmentOperator()
		if err != nil {
			return err
		}
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		nae, err := g.parseAssignmentExpression(in, yield, await)
		if err != nil {
			return err
		}
		j.Score(g)
		ae.AssignmentExpression = &nae
		return nil
	}); err != nil {
		return ae, err
	}
	ae.Tokens = j.ToTokens()
	return ae, nil
}

type LeftHandSideExpression struct {
	NewExpression  *NewExpression
	CallExpression *CallExpression
	Tokens         []TokenPos
}

func (j *jsParser) parseLeftHandSideExpression(yield, await bool) (LeftHandSideExpression, error) {
	var lhs LeftHandSideExpression
	if err := j.findGoal(func(j *jsParser) error {
		ne, err := j.parseNewExpression(yield, await)
		if err != nil {
			return err
		}
		lhs.NewExpression = &ne
		return nil
	}, func(j *jsParser) error {
		ce, err := j.parseCallExpression(yield, await)
		if err != nil {
			return err
		}
		lhs.CallExpression = &ce
		return nil
	}); err != nil {
		return lhs, err
	}
	lhs.Tokens = j.ToTokens()
	return lhs, nil
}

type Expression struct {
	Expressions []AssignmentExpression
	Tokens      []TokenPos
}

func (j *jsParser) parseExpression(in, yield, await bool) (Expression, error) {
	var e Expression
	for {
		g := j.NewGoal()
		ae, err := g.parseAssignmentExpression(in, yield, await)
		if err != nil {
			return e, j.Error(err)
		}
		j.Score(g)
		e.Expressions = append(e.Expressions, ae)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			break
		}
		j.Score(g)
	}
	e.Tokens = j.ToTokens()
	return e, nil
}

type NewExpression struct {
	NewExpression    *NewExpression
	MemberExpression *MemberExpression
	Tokens           []TokenPos
}

func (j *jsParser) parseNewExpression(yield, await bool) (NewExpression, error) {
	var ne NewExpression
	ne.Tokens = j.ToTokens()
	return ne, nil
}

type MemberExpression struct {
	MemberExpression    *MemberExpression
	PrimaryExpression   *PrimaryExpression
	Expression          *Expression
	IdentifierName      *TokenPos
	TemplateLiteral     *TemplateLiteral
	SuperProperty       bool
	MetaProperty        bool
	NewTarget           bool
	NewMemberExpression *MemberExpression
	Arguments           *Arguments
	Tokens              []TokenPos
}

func (j *jsParser) parseMemberExpression(yield, await bool) (MemberExpression, error) {
	var me MemberExpression
	if err := j.findGoal(
		func(j *jsParser) error {
			pe, err := j.parserPrimaryExpression(yield, await)
			if err != nil {
				return err
			}
			me.PrimaryExpression = &pe
			return nil
		},
		func(j *jsParser) error {
			if !j.AcceptToken(parser.Token{TokenKeyword, "super"}) {
				return errNotApplicable
			}
			j.AcceptRunWhitespace()
			if j.AcceptToken(parser.Token{TokenPunctuator, "["}) {
				g := j.NewGoal()
				e, err := g.parseExpression(true, yield, await)
				if err != nil {
					return err
				}
				j.Score(g)
				me.Expression = &e
			} else if j.AcceptToken(parser.Token{TokenPunctuator, "."}) {
				if !j.Accept(TokenIdentifier) {
					return ErrNoIdentifier
				}
				me.IdentifierName = j.GetLastToken()
			} else {
				return ErrInvalidSuperProperty
			}
			me.SuperProperty = true
			return nil
		},
		func(j *jsParser) error {
			if !j.AcceptToken(parser.Token{TokenKeyword, "new"}) {
				return errNotApplicable
			}
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, "."}) {
				return errNotApplicable
			}
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenIdentifier, "target"}) {
				return ErrInvalidMetaProperty
			}
			me.MetaProperty = true
			return nil
		},
		func(j *jsParser) error {
			if !j.AcceptToken(parser.Token{TokenKeyword, "new"}) {
				return errNotApplicable
			}
			j.AcceptRunWhitespace()
			g := j.NewGoal()
			nme, err := g.parseMemberExpression(yield, await)
			if err != nil {
				return err
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			g = j.NewGoal()
			a, err := g.parseArguments(yield, await)
			if err != nil {
				return err
			}
			j.Score(g)
			me.NewMemberExpression = &nme
			me.Arguments = &a
			return nil
		},
	); err != nil {
		return me, err
	}
Loop:
	for {
		var nme MemberExpression
		g := *j
		g.AcceptRunWhitespace()
		h := g.NewGoal()
		switch tk := h.Peek(); tk.Type {
		case TokenNoSubstitutionTemplate, TokenTemplateHead:
			tl, err := h.parserTemplateLiteral(yield, await)
			if err != nil {
				return me, g.Error(err)
			}
			nme.TemplateLiteral = &tl
		case TokenIdentifier:
			h.Except()
			nme.IdentifierName = h.GetLastToken()
		case TokenPunctuator:
			if tk.Data != "[" {
				break Loop
			}
			h.Except()
			h.AcceptRunWhitespace()
			i := h.NewGoal()
			e, err := i.parseExpression(true, yield, await)
			if err != nil {
				return me, g.Error(err)
			}
			h.Score(i)
			h.AcceptRunWhitespace()
			if !h.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
				return me, g.Error(ErrMissingClosingBracket)
			}
			nme.Expression = &e
		default:
			break Loop
		}
		g.Score(h)
		ome := me
		ome.Tokens = j.ToTokens()
		nme.MemberExpression = &ome
		me = nme
		*j = g
	}
	me.Tokens = j.ToTokens()
	return me, nil
}

type PrimaryExpression struct {
	This                                              bool
	IdentifierReference                               *IdentifierReference
	Literal                                           *TokenPos
	ArrayLiteral                                      *ArrayLiteral
	ObjectLiteral                                     *ObjectLiteral
	FunctionExpression                                *FunctionDeclaration
	ClassExpression                                   *ClassDeclaration
	TemplateLiteral                                   *TemplateLiteral
	CoverParenthesizedExpressionAndArrowParameterList *CoverParenthesizedExpressionAndArrowParameterList
	Tokens                                            []TokenPos
}

func (j *jsParser) parserPrimaryExpression(yield, await bool) (PrimaryExpression, error) {
	var pe PrimaryExpression
	if err := j.findGoal(
		func(j *jsParser) error {
			if !j.AcceptToken(parser.Token{TokenKeyword, "this"}) {
				return errNotApplicable
			}
			pe.This = true
			return nil
		},
		func(j *jsParser) error {
			i, err := j.parseIdentifierReference(yield, await)
			if err != nil {
				return err
			}
			pe.IdentifierReference = &i
			return nil
		},
		func(j *jsParser) error {
			if !j.Accept(TokenNullLiteral, TokenBooleanLiteral, TokenNumericLiteral, TokenStringLiteral, TokenRegularExpressionLiteral) {
				return errNotApplicable
			}
			pe.Literal = j.GetLastToken()
			return nil
		},
		func(j *jsParser) error {
			al, err := j.parseArrayLiteral(yield, await)
			if err != nil {
				return err
			}
			pe.ArrayLiteral = &al
			return nil
		},
		func(j *jsParser) error {
			ol, err := j.parseObjectLiteral(yield, await)
			if err != nil {
				return err
			}
			pe.ObjectLiteral = &ol
			return nil
		},
		func(j *jsParser) error {
			fe, err := j.parseFunctionDeclaration(false, false, true)
			if err != nil {
				return err
			}
			pe.FunctionExpression = &fe
			return nil
		},
		func(j *jsParser) error {
			ce, err := j.parseClassDeclaration(yield, await, true)
			if err != nil {
				return err
			}
			pe.ClassExpression = &ce
			return nil
		},
		func(j *jsParser) error {
			tl, err := j.parserTemplateLiteral(yield, await)
			if err != nil {
				return err
			}
			pe.TemplateLiteral = &tl
			return nil
		},
		func(j *jsParser) error {
			cp, err := j.parseCoverParenthesizedExpressionAndArrowParameterList(yield, await)
			if err != nil {
				return err
			}
			pe.CoverParenthesizedExpressionAndArrowParameterList = &cp
			return nil
		},
	); err != nil {
		return pe, err
	}
	pe.Tokens = j.ToTokens()
	return pe, nil
}

type ArrayLiteral struct {
	ElementList   []AssignmentExpression
	SpreadElement *AssignmentExpression
	Tokens        []TokenPos
}

func (j *jsParser) parseArrayLiteral(yield, await bool) (ArrayLiteral, error) {
	var al ArrayLiteral
	if !j.AcceptToken(parser.Token{TokenPunctuator, "["}) {
		return al, j.Error(ErrMissingOpeningBracket)
	}
	for {
		var spread bool
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
			break
		} else if j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			al.ElementList = append(al.ElementList, AssignmentExpression{})
			continue
		} else if j.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
			spread = true
		}
		g := j.NewGoal()
		ae, err := j.parseAssignmentExpression(true, yield, await)
		if err != nil {
			return al, j.Error(err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		if spread {
			al.SpreadElement = &ae
			if !j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
				return al, j.Error(ErrMissingClosingBracket)
			}
			break
		}
		if j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return al, j.Error(ErrMissingComma)
		}
	}
	al.Tokens = j.ToTokens()
	return al, nil
}

type ObjectLiteral struct {
	PropertyDefinitionList []PropertyDefinition
	Tokens                 []TokenPos
}

func (j *jsParser) parseObjectLiteral(yield, await bool) (ObjectLiteral, error) {
	var ol ObjectLiteral
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return ol, j.Error(ErrMissingOpeningBrace)
	}
	for {
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		}
		g := j.NewGoal()
		pd, err := g.parsePropertyDefinition(yield, await)
		if err != nil {
			return ol, j.Error(err)
		}
		j.Score(g)
		ol.PropertyDefinitionList = append(ol.PropertyDefinitionList, pd)
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return ol, j.Error(ErrMissingComma)
		}
	}
	ol.Tokens = j.ToTokens()
	return ol, nil
}

type PropertyDefinition struct {
	IdentifierReference  *IdentifierReference
	PropertyName         *PropertyName
	AssignmentExpression *AssignmentExpression
	MethodDefinition     *MethodDefinition
	Tokens               []TokenPos
}

func (j *jsParser) parsePropertyDefinition(yield, await bool) (PropertyDefinition, error) {
	var pd PropertyDefinition
	if err := j.findGoal(
		func(j *jsParser) error {
			ir, err := j.parseIdentifierReference(yield, await)
			if err != nil {
				return err
			}
			j.AcceptRunWhitespace()
			if j.AcceptToken(parser.Token{TokenPunctuator, "="}) {
				g := j.NewGoal()
				ae, err := g.parseAssignmentExpression(true, yield, await)
				if err != nil {
					return err
				}
				j.Score(g)
				pd.AssignmentExpression = &ae
			}
			pd.IdentifierReference = &ir
			return nil
		},
		func(j *jsParser) error {
			pn, err := j.parsePropertyName(yield, await)
			if err != nil {
				return err
			}
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
				return j.Error(ErrMissingColon)
			}
			j.AcceptRunWhitespace()
			g := j.NewGoal()
			ae, err := j.parseAssignmentExpression(true, yield, await)
			if err != nil {
				return err
			}
			j.Score(g)
			pd.PropertyName = &pn
			pd.AssignmentExpression = &ae
			return nil
		},
		func(j *jsParser) error {
			md, err := j.parseMethodDefinition(yield, await)
			if err != nil {
				return err
			}
			pd.MethodDefinition = &md
			return nil
		},
	); err != nil {
		return pd, err
	}
	return pd, nil
}

type CoverParenthesizedExpressionAndArrowParameterList struct {
	Expressions          []Expression
	BindingIdentifier    *BindingIdentifier
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Tokens               []TokenPos
}

func (j *jsParser) parseCoverParenthesizedExpressionAndArrowParameterList(yield, await bool) (CoverParenthesizedExpressionAndArrowParameterList, error) {
	var cp CoverParenthesizedExpressionAndArrowParameterList
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return cp, j.Error(ErrMissingOpeningParentheses)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		for {
			if j.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
				j.AcceptRunWhitespace()
				g := j.NewGoal()
				if g.AcceptToken(parser.Token{TokenPunctuator, "["}) {
					ab, err := j.parseArrayBindingPattern(yield, await)
					if err != nil {
						return cp, j.Error(err)
					}
					cp.ArrayBindingPattern = &ab

				} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
					ob, err := j.parseObjectBindingPattern(yield, await)
					if err != nil {
						return cp, j.Error(err)
					}
					cp.ObjectBindingPattern = &ob
				} else {
					bi, err := g.parseBindingIdentifier(yield, await)
					if err != nil {
						return cp, j.Error(err)
					}
					cp.BindingIdentifier = &bi
				}
				j.Score(g)
				j.AcceptRunWhitespace()
				if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
					return cp, j.Error(ErrMissingClosingParentheses)
				}
				break
			}
			g := j.NewGoal()
			e, err := g.parseExpression(true, yield, await)
			if err != nil {
				return cp, j.Error(err)
			}
			j.Score(g)
			cp.Expressions = append(cp.Expressions, e)
			j.AcceptRunWhitespace()
			if j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
				break
			} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
				return cp, j.Error(ErrMissingComma)
			}
			j.AcceptRunWhitespace()
		}
	}
	cp.Tokens = j.ToTokens()
	return cp, nil
}

type TemplateLiteral struct {
	NoSubstitutionTemplate *TokenPos
	TemplateHead           *TokenPos
	Expressions            []Expression
	TemplateMiddleList     []*TokenPos
	TemplateTail           *TokenPos
	Tokens                 []TokenPos
}

func (j *jsParser) parserTemplateLiteral(yield, await bool) (TemplateLiteral, error) {
	var tl TemplateLiteral
	if j.Accept(TokenNoSubstitutionTemplate) {
		tl.NoSubstitutionTemplate = j.GetLastToken()
	} else if !j.Accept(TokenTemplateHead) {
		return tl, j.Error(ErrInvalidTemplate)
	} else {
		tl.TemplateHead = j.GetLastToken()
		for {
			j.AcceptRunWhitespace()
			g := j.NewGoal()
			e, err := j.parseExpression(true, yield, await)
			if err != nil {
				return tl, j.Error(err)
			}
			j.Score(g)
			tl.Expressions = append(tl.Expressions, e)
			j.AcceptRunWhitespace()
			if j.Accept(TokenTemplateTail) {
				tl.TemplateTail = j.GetLastToken()
				break
			} else if !j.Accept(TokenTemplateMiddle) {
				return tl, j.Error(ErrInvalidTemplate)
			}
			tl.TemplateMiddleList = append(tl.TemplateMiddleList, j.GetLastToken())
		}
	}
	tl.Tokens = j.ToTokens()
	return tl, nil
}

type Arguments struct {
	ArgumentList   []AssignmentExpression
	SpreadArgument *AssignmentExpression
	Tokens         []TokenPos
}

func (j *jsParser) parseArguments(yield, await bool) (Arguments, error) {
	var a Arguments
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return a, j.Error(ErrMissingOpeningParentheses)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		for {
			var spread bool
			if j.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
				spread = true
				j.AcceptRunWhitespace()
			}
			g := j.NewGoal()
			ae, err := g.parseAssignmentExpression(true, yield, await)
			if err != nil {
				return a, j.Error(err)
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			if spread {
				a.SpreadArgument = &ae
				if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
					return a, j.Error(err)
				}
				break
			}
			a.ArgumentList = append(a.ArgumentList, ae)
			if j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
				break
			} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
				return a, j.Error(ErrMissingComma)
			}
			j.AcceptRunWhitespace()
		}
	}
	a.Tokens = j.ToTokens()
	return a, nil
}

type CallExpression struct {
	CoverCallExpressionAndAsyncArrowHead *MemberExpression
	SuperCall                            bool
	CallExpression                       *CallExpression
	Arguments                            *Arguments
	Expression                           *Expression
	IdentifierName                       *TokenPos
	TemplateLiteral                      *TemplateLiteral
	Tokens                               []TokenPos
}

func (j *jsParser) parseCallExpression(yield, await bool) (CallExpression, error) {
	var ce CallExpression
	if j.AcceptToken(parser.Token{TokenKeyword, "super"}) {
		ce.SuperCall = true
	} else {
		g := j.NewGoal()
		me, err := g.parseMemberExpression(yield, await)
		if err != nil {
			return ce, j.Error(err)
		}
		ce.CoverCallExpressionAndAsyncArrowHead = &me
	}
	j.AcceptRunWhitespace()
	a, err := j.parseArguments(yield, await)
	if err != nil {
		return ce, err
	}
	ce.Arguments = &a
	for {
		oce := ce
		nce := CallExpression{
			CallExpression: &oce,
		}
		g := j.NewGoal()
		g.AcceptRunWhitespace()
		if tk := g.Peek(); tk == (parser.Token{TokenPunctuator, "("}) {
			h := g.NewGoal()
			a, err := h.parseArguments(yield, await)
			if err != nil {
				return ce, j.Error(err)
			}
			g.Score(h)
			nce.Arguments = &a
		} else if tk == (parser.Token{TokenPunctuator, "["}) {
			g.Except()
			g.AcceptRunWhitespace()
			h := g.NewGoal()
			e, err := h.parseExpression(true, yield, await)
			if err != nil {
				return ce, j.Error(err)
			}
			g.Score(h)
			if !g.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
				return ce, j.Error(ErrMissingClosingBracket)
			}
			nce.Expression = &e
		} else if tk == (parser.Token{TokenPunctuator, "."}) {
			g.Except()
			g.AcceptRunWhitespace()
			if !g.Accept(TokenIdentifier, TokenKeyword) {
				return ce, j.Error(ErrNoIdentifier)
			}
			nce.IdentifierName = g.GetLastToken()
		} else if tk.Type == TokenTemplateHead || tk.Type == TokenNoSubstitutionTemplate {
			h := g.NewGoal()
			tl, err := h.parserTemplateLiteral(yield, await)
			if err != nil {
				return ce, j.Error(err)
			}
			g.Score(h)
			nce.TemplateLiteral = &tl
		} else {
			break
		}
		j.Score(g)
		nce.Tokens = j.ToTokens()
		ce = nce
	}
	ce.Tokens = j.ToTokens()
	return ce, nil
}

type ConditionalExpression struct {
	Tokens []TokenPos
}

func (j *jsParser) parseConditionalExpression(in, yield, await bool) (ConditionalExpression, error) {
	var ce ConditionalExpression
	return ce, nil
}

type ArrowFunction struct {
	Tokens []TokenPos
}

func (j *jsParser) parseArrowFunction(in, yield, await bool) (ArrowFunction, error) {
	var af ArrowFunction
	return af, nil
}

const (
	ErrInvalidStatementList       errors.Error = "invalid statement list"
	ErrMissingSemiColon           errors.Error = "missing semi-colon"
	ErrMissingColon               errors.Error = "missing colon"
	ErrNoIdentifier               errors.Error = "missing identifier"
	ErrReservedIdentifier         errors.Error = "reserved identifier"
	ErrMissingFunction            errors.Error = "missing function"
	ErrMissingOpeningParentheses  errors.Error = "missing opening parentheses"
	ErrMissingClosingParentheses  errors.Error = "missing closing parentheses"
	ErrMissingOpeningBrace        errors.Error = "missing opening brace"
	ErrMissingClosingBrace        errors.Error = "missing closing brace"
	ErrMissingOpeningBracket      errors.Error = "missing opening bracket"
	ErrMissingClosingBracket      errors.Error = "missing closing bracket"
	ErrMissingComma               errors.Error = "missing comma"
	ErrInvalidFormalParameterList errors.Error = "invalid formal parameter list"
	ErrInvalidDeclaration         errors.Error = "invalid declaration"
	ErrInvalidLexicalDeclaration  errors.Error = "invalid lexical declaration"
	ErrInvalidAssignment          errors.Error = "invalid assignment operator"
	ErrInvalidSuperProperty       errors.Error = "invalid super property"
	ErrInvalidMetaProperty        errors.Error = "invalid meta property"
	ErrInvalidTemplate            errors.Error = "invalid template"
)
