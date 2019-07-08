package javascript

import (
	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

type Script struct {
	StatementList []StatementListItem
	Tokens        Tokens
}

func ParseScript(t parser.Tokeniser) (Script, error) {
	j, err := newJSParser(t)
	if err != nil {
		return Script{}, err
	}
	return j.parseScript()
}

func (j *jsParser) parseScript() (Script, error) {
	var s Script
	for j.AcceptRunWhitespace() != parser.TokenDone {
		g := j.NewGoal()
		si, err := g.parseStatementListItem(false, false, false)
		if err != nil {
			return s, err
		}
		j.Score(g)
		s.StatementList = append(s.StatementList, si)
	}
	s.Tokens = j.ToTokens()
	return s, nil
}

func (j *jsParser) parseIdentifier(yield, await bool) (*Token, error) {
	if j.Accept(TokenIdentifier) || (!yield && j.AcceptToken(parser.Token{TokenKeyword, "yield"}) || (!await && j.AcceptToken(parser.Token{TokenKeyword, "await"}))) {
		return j.GetLastToken(), nil
	}
	return nil, j.Error(ErrNoIdentifier)
}

type Declaration struct {
	ClassDeclaration    *ClassDeclaration
	FunctionDeclaration *FunctionDeclaration
	LexicalDeclaration  *LexicalDeclaration
	Tokens              Tokens
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
	} else if tk == (parser.Token{TokenIdentifier, "async"}) || tk == (parser.Token{TokenKeyword, "function"}) {
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
	Tokens      Tokens
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
	BindingIdentifier    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Initializer          *AssignmentExpression
	Tokens               Tokens
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
		bi, err := g.parseIdentifier(yield, await)
		if err != nil {
			return lb, j.Error(err)
		}
		lb.BindingIdentifier = bi
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
	Tokens             Tokens
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
	ab.Tokens = j.ToTokens()
	return ab, nil
}

type ObjectBindingPattern struct {
	BindingPropertyList []BindingProperty
	BindingRestProperty *Token
	Tokens              Tokens
}

func (j *jsParser) parseObjectBindingPattern(yield, await bool) (ObjectBindingPattern, error) {
	var ob ObjectBindingPattern
	j.AcceptToken(parser.Token{TokenPunctuator, "{"})
	j.AcceptRunWhitespace()
	if !j.Accept(TokenRightBracePunctuator) {
		for {
			g := j.NewGoal()
			if g.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
				bi, err := g.parseIdentifier(yield, await)
				if err != nil {
					return ob, j.Error(err)
				}
				j.Score(g)
				ob.BindingRestProperty = bi
				j.AcceptRunWhitespace()
				if !j.Accept(TokenRightBracePunctuator) {
					return ob, j.Error(ErrMissingClosingBrace)
				}
				break
			}
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
			j.AcceptRunWhitespace()
		}
	}
	ob.Tokens = j.ToTokens()
	return ob, nil
}

type BindingProperty struct {
	SingleNameBinding *Token
	Initializer       *AssignmentExpression
	PropertyName      *PropertyName
	BindingElement    *BindingElement
	Tokens            Tokens
}

func (j *jsParser) parseBindingProperty(yield, await bool) (BindingProperty, error) {
	var bp BindingProperty
	if err := j.FindGoal(
		func(j *jsParser) error {
			pn, err := j.parsePropertyName(yield, await)
			if err != nil {
				return err
			}
			j.AcceptRunWhitespace()
			if j.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
				return ErrMissingColon
			}
			j.AcceptRunWhitespace()
			g := j.NewGoal()
			be, err := g.parseBindingElement(yield, await)
			if err != nil {
				return j.Error(err)
			}
			j.Score(g)
			bp.PropertyName = &pn
			bp.BindingElement = &be
			return nil
		},
		func(j *jsParser) error {
			bi, err := j.parseIdentifier(yield, await)
			if err != nil {
				return err
			}
			g := j.NewGoal()
			g.AcceptRunWhitespace()
			if g.AcceptToken(parser.Token{TokenPunctuator, "="}) {
				g.AcceptRunWhitespace()
				h := g.NewGoal()
				i, err := h.parseAssignmentExpression(true, yield, await)
				if err != nil {
					return g.Error(err)
				}
				g.Score(h)
				j.Score(g)
				bp.Initializer = &i
			}
			bp.SingleNameBinding = bi
			return nil
		},
	); err != nil {
		return bp, err
	}
	bp.Tokens = j.ToTokens()
	return bp, nil
}

type VariableDeclaration LexicalBinding

func (j *jsParser) parseVariableDeclaration(in, yield, await bool) (VariableDeclaration, error) {
	lb, err := j.parseLexicalBinding(in, yield, await)
	return VariableDeclaration(lb), err
}

type ArrayLiteral struct {
	ElementList   []AssignmentExpression
	SpreadElement *AssignmentExpression
	Tokens        Tokens
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
		ae, err := g.parseAssignmentExpression(true, yield, await)
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
		} else {
			al.ElementList = append(al.ElementList, ae)
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
	Tokens                 Tokens
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
	IdentifierReference  *Token
	PropertyName         *PropertyName
	Spread               bool
	AssignmentExpression *AssignmentExpression
	MethodDefinition     *MethodDefinition
	Tokens               Tokens
}

func (j *jsParser) parsePropertyDefinition(yield, await bool) (PropertyDefinition, error) {
	var pd PropertyDefinition
	if err := j.FindGoal(
		func(j *jsParser) error {
			pn, err := j.parsePropertyName(yield, await)
			if err != nil {
				return err
			}
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
				return ErrMissingColon
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
			ir, err := j.parseIdentifier(yield, await)
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
			pd.IdentifierReference = ir
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
		func(j *jsParser) error {
			if !j.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
				return errNotApplicable
			}
			ae, err := j.parseAssignmentExpression(true, yield, await)
			if err != nil {
				return err
			}
			pd.Spread = true
			pd.AssignmentExpression = &ae
			return nil
		},
	); err != nil {
		return pd, err
	}
	pd.Tokens = j.ToTokens()
	return pd, nil
}

type TemplateLiteral struct {
	NoSubstitutionTemplate *Token
	TemplateHead           *Token
	Expressions            []Expression
	TemplateMiddleList     []*Token
	TemplateTail           *Token
	Tokens                 Tokens
}

func (j *jsParser) parseTemplateLiteral(yield, await bool) (TemplateLiteral, error) {
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

type ArrowFunction struct {
	Async                                             bool
	BindingIdentifier                                 *Token
	CoverParenthesizedExpressionAndArrowParameterList *CoverParenthesizedExpressionAndArrowParameterList
	FormalParameters                                  *FormalParameters
	AssignmentExpression                              *AssignmentExpression
	FunctionBody                                      *Block
	Tokens                                            Tokens
}

func (j *jsParser) parseArrowFunction(in, yield, await bool) (ArrowFunction, error) {
	var af ArrowFunction
	if j.AcceptToken(parser.Token{TokenIdentifier, "async"}) {
		af.Async = true
		j.AcceptRunWhitespaceNoNewLine()
		if j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
			g := j.NewGoal()
			fp, err := g.parseFormalParameters(false, true)
			if err != nil {
				return af, j.Error(err)
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
				return af, j.Error(ErrMissingClosingParenthesis)
			}
			af.FormalParameters = &fp
		} else {
			g := j.NewGoal()
			bi, err := g.parseIdentifier(yield, true)
			if err != nil {
				return af, j.Error(err)
			}
			j.Score(g)
			af.BindingIdentifier = bi
		}
	} else if j.Peek() == (parser.Token{TokenPunctuator, "("}) {
		g := j.NewGoal()
		cp, err := g.parseCoverParenthesizedExpressionAndArrowParameterList(yield, await)
		if err != nil {
			return af, j.Error(err)
		}
		j.Score(g)
		af.CoverParenthesizedExpressionAndArrowParameterList = &cp
	} else {
		g := j.NewGoal()
		bi, err := g.parseIdentifier(yield, true)
		if err != nil {
			return af, j.Error(err)
		}
		j.Score(g)
		af.BindingIdentifier = bi
	}
	j.AcceptRunWhitespaceNoNewLine()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "=>"}) {
		return af, j.Error(ErrMissingArrow)
	}
	j.AcceptRunWhitespace()
	if j.Peek() == (parser.Token{TokenPunctuator, "{"}) {
		g := j.NewGoal()
		b, err := g.parseBlock(false, af.Async, true)
		if err != nil {
			return af, j.Error(err)
		}
		j.Score(g)
		af.FunctionBody = &b
	} else {
		g := j.NewGoal()
		ae, err := g.parseAssignmentExpression(in, false, af.Async)
		if err != nil {
			return af, j.Error(err)
		}
		j.Score(g)
		af.AssignmentExpression = &ae
	}
	af.Tokens = j.ToTokens()
	return af, nil
}

const (
	ErrReservedIdentifier          errors.Error = "reserved identifier"
	ErrNoIdentifier                errors.Error = "missing identifier"
	ErrMissingFunction             errors.Error = "missing function"
	ErrMissingOpeningParenthesis   errors.Error = "missing opening parenthesis"
	ErrMissingClosingParenthesis   errors.Error = "missing closing parenthesis"
	ErrMissingOpeningBrace         errors.Error = "missing opening brace"
	ErrMissingClosingBrace         errors.Error = "missing closing brace"
	ErrMissingOpeningBracket       errors.Error = "missing opening bracket"
	ErrMissingClosingBracket       errors.Error = "missing closing bracket"
	ErrMissingComma                errors.Error = "missing comma"
	ErrMissingArrow                errors.Error = "missing arrow"
	ErrMissingCaseClause           errors.Error = "missing case clause"
	ErrMissingExpression           errors.Error = "missing expression"
	ErrMissingCatchFinally         errors.Error = "missing catch/finally block"
	ErrMissingSemiColon            errors.Error = "missing semi-colon"
	ErrMissingColon                errors.Error = "missing colon"
	ErrInvalidStatementList        errors.Error = "invalid statement list"
	ErrInvalidStatement            errors.Error = "invalid statement"
	ErrInvalidFormalParameterList  errors.Error = "invalid formal parameter list"
	ErrInvalidDeclaration          errors.Error = "invalid declaration"
	ErrInvalidLexicalDeclaration   errors.Error = "invalid lexical declaration"
	ErrInvalidAssignment           errors.Error = "invalid assignment operator"
	ErrInvalidSuperProperty        errors.Error = "invalid super property"
	ErrInvalidMetaProperty         errors.Error = "invalid meta property"
	ErrInvalidTemplate             errors.Error = "invalid template"
	ErrInvalidIterationStatementDo errors.Error = "invalid do interation statement"
	ErrInvalidForLoop              errors.Error = "invalid for loop"
)
