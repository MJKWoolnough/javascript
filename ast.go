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
			if err.(Error).getLastPos() > errr.(Error).getLastPos() {
				return bp, j.Error(err)
			}
			return bp, j.Error(errr)
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

type AssignmentExpression struct {
	Tokens []TokenPos
}

func (j *jsParser) parseAssignmentExpression(in, yield, await bool) (AssignmentExpression, error) {
	var ae AssignmentExpression
	return ae, nil
}

type LeftHandSideExpression struct {
	Tokens []TokenPos
}

func (j *jsParser) parseLeftHandSideExpression(yield, await bool) (LeftHandSideExpression, error) {
	var lhs LeftHandSideExpression
	return lhs, nil
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
	ErrMissingClosingBracket      errors.Error = "missing closing bracket"
	ErrMissingComma               errors.Error = "missing comma"
	ErrInvalidFormalParameterList errors.Error = "invalid formal parameter list"
	ErrInvalidDeclaration         errors.Error = "invalid declaration"
	ErrInvalidLexicalDeclaration  errors.Error = "invalid lexical declaration"
)
