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

type VariableDeclaration struct {
	Tokens []TokenPos
}

type Declaration struct {
	Tokens []TokenPos
}

func (j *jsParser) parseDeclaration(yield, await bool) (Declaration, error) {
	var d Declaration
	return d, nil
}

type ArrayBindingPattern struct {
	Token []TokenPos
}

func (j *jsParser) parseArrayBindingPattern(yield, await bool) (ArrayBindingPattern, error) {
	var ab ArrayBindingPattern
	return ab, nil
}

type ObjectBindingPattern struct {
	Token []TokenPos
}

func (j *jsParser) parseObjectBindingPattern(yield, await bool) (ObjectBindingPattern, error) {
	var ob ObjectBindingPattern
	return ob, nil
}

type AssignmentExpression struct {
	Tokens []TokenPos
}

func (j *jsParser) parseAssignmentExpression(in, yield, await bool) (AssignmentExpression, error) {
	var ae AssignmentExpression
	return ae, nil
}

type VariableStatement struct {
	Tokens []TokenPos
}

func (j *jsParser) parseVariableStatement(yield, await bool) (VariableStatement, error) {
	var vs VariableStatement
	return vs, nil
}

const (
	ErrInvalidStatementList       errors.Error = "invalid statement list"
	ErrMissingSemiColon           errors.Error = "missing semi-colon"
	ErrNoIdentifier               errors.Error = "missing identifier"
	ErrReservedIdentifier         errors.Error = "reserved identifier"
	ErrMissingFunction            errors.Error = "missing function"
	ErrMissingOpeningParentheses  errors.Error = "missing opening parentheses"
	ErrMissingClosingParentheses  errors.Error = "missing closing parentheses"
	ErrMissingOpeningBrace        errors.Error = "missing opening brace"
	ErrMissingClosingBrace        errors.Error = "missing closing brace"
	ErrInvalidFormalParameterList errors.Error = "invalid formal parameter list"
)
