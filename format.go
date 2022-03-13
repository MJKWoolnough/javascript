package javascript

import (
	"fmt"
	"io"

	"vimagination.zapto.org/parser"
)

type indentPrinter struct {
	io.Writer
}

var indent = []byte{'	'}

func (i *indentPrinter) Write(p []byte) (int, error) {
	var (
		total int
		last  int
	)
	for n, c := range p {
		if c == '\n' {
			m, err := i.Writer.Write(p[last : n+1])
			total += m
			if err != nil {
				return total, err
			}
			_, err = i.Writer.Write(indent)
			if err != nil {
				return total, err
			}
			last = n + 1
		}
	}
	if last != len(p) {
		m, err := i.Writer.Write(p[last:])
		total += m
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func (i *indentPrinter) Print(args ...interface{}) {
	fmt.Fprint(i, args...)
}

func (i *indentPrinter) Printf(format string, args ...interface{}) {
	fmt.Fprintf(i, format, args...)
}

func (i *indentPrinter) WriteString(s string) (int, error) {
	return i.Write([]byte(s))
}

// Format implements the fmt.Formatter interface
func (t Token) Format(s fmt.State, v rune) {
	t.printType(s, s.Flag('+'))
}

func (t Token) printType(w io.Writer, v bool) {
	var typ string
	switch t.Type {
	case parser.TokenError:
		typ = "Error"
	case parser.TokenDone:
		typ = "Done"
	case TokenWhitespace:
		typ = "Whitespace"
	case TokenLineTerminator:
		typ = "LineTerminator"
	case TokenSingleLineComment:
		typ = "SingleLineComment"
	case TokenMultiLineComment:
		typ = "MultiLineComment"
	case TokenIdentifier:
		typ = "Identifier"
	case TokenBooleanLiteral:
		typ = "BooleanLiteral"
	case TokenKeyword:
		typ = "Keyword"
	case TokenPunctuator:
		typ = "Punctuator"
	case TokenNumericLiteral:
		typ = "NumericLiteral"
	case TokenStringLiteral:
		typ = "StringLiteral"
	case TokenNoSubstitutionTemplate:
		typ = "NoSubstitutionTemplate"
	case TokenTemplateHead:
		typ = "TemplateHead"
	case TokenTemplateMiddle:
		typ = "TemplateMiddle"
	case TokenTemplateTail:
		typ = "TemplateTail"
	case TokenDivPunctuator:
		typ = "DivPunctuator"
	case TokenRightBracePunctuator:
		typ = "RightBracePunctuator"
	case TokenRegularExpressionLiteral:
		typ = "RegulatExpressionLiteral"
	case TokenNullLiteral:
		typ = "NullLiteral"
	case TokenFutureReservedWord:
		typ = "FutureReservedWord"
	default:
		typ = fmt.Sprintf("%d", t.Type)
	}
	fmt.Fprintf(w, "Type: %s - Data: %q", typ, t.Data)
	if v {
		fmt.Fprintf(w, " - Position: %d (%d: %d)", t.Pos, t.Line, t.LinePos)
	}
}

// Format implements the fmt.Formatter interface
func (t Tokens) Format(s fmt.State, v rune) {
	t.printType(s, s.Flag('+'))
}

func (t Tokens) printType(w io.Writer, v bool) {
	if len(t) == 0 {
		w.Write(arrayOpenClose)
		return
	}
	w.Write(arrayOpen)
	ipp := indentPrinter{w}
	for n, t := range t {
		ipp.Printf("\n%d: ", n)
		t.printType(w, v)
	}
	w.Write(arrayClose)
}

var (
	space          = []byte{' '}
	arrayOpen      = []byte{'['}
	arrayClose     = []byte{'\n', ']'}
	arrayOpenClose = []byte{'[', ']'}
	objectOpen     = []byte{'{'}
	objectClose    = []byte{'\n', '}'}
	nilStr         = []byte{'<', 'n', 'i', 'l', '>'}
	tokensTo       = []byte{'\n', 'T', 'o', 'k', 'e', 'n', 's', ':', ' '}
)

type formatter interface {
	printType(io.Writer, bool)
	printSource(io.Writer, bool)
}

func format(f formatter, s fmt.State, v rune) {
	switch v {
	case 'v':
		f.printType(s, s.Flag('+'))
	case 's':
		f.printSource(s, s.Flag('+'))
	}
}

// Format implements the fmt.Formatter interface
func (f AdditiveExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _AdditiveExpression AdditiveExpression
		fmt.Fprintf(s, "%#v", _AdditiveExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Arguments) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _Arguments Arguments
		fmt.Fprintf(s, "%#v", _Arguments(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ArrayAssignmentPattern) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ArrayAssignmentPattern ArrayAssignmentPattern
		fmt.Fprintf(s, "%#v", _ArrayAssignmentPattern(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ArrayBindingPattern) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ArrayBindingPattern ArrayBindingPattern
		fmt.Fprintf(s, "%#v", _ArrayBindingPattern(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ArrayLiteral) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ArrayLiteral ArrayLiteral
		fmt.Fprintf(s, "%#v", _ArrayLiteral(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ArrowFunction) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ArrowFunction ArrowFunction
		fmt.Fprintf(s, "%#v", _ArrowFunction(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f AssignmentElement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _AssignmentElement AssignmentElement
		fmt.Fprintf(s, "%#v", _AssignmentElement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f AssignmentExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _AssignmentExpression AssignmentExpression
		fmt.Fprintf(s, "%#v", _AssignmentExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f AssignmentPattern) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _AssignmentPattern AssignmentPattern
		fmt.Fprintf(s, "%#v", _AssignmentPattern(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f AssignmentProperty) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _AssignmentProperty AssignmentProperty
		fmt.Fprintf(s, "%#v", _AssignmentProperty(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f BindingElement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _BindingElement BindingElement
		fmt.Fprintf(s, "%#v", _BindingElement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f BindingProperty) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _BindingProperty BindingProperty
		fmt.Fprintf(s, "%#v", _BindingProperty(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f BitwiseANDExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _BitwiseANDExpression BitwiseANDExpression
		fmt.Fprintf(s, "%#v", _BitwiseANDExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f BitwiseORExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _BitwiseORExpression BitwiseORExpression
		fmt.Fprintf(s, "%#v", _BitwiseORExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f BitwiseXORExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _BitwiseXORExpression BitwiseXORExpression
		fmt.Fprintf(s, "%#v", _BitwiseXORExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Block) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _Block Block
		fmt.Fprintf(s, "%#v", _Block(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f CallExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _CallExpression CallExpression
		fmt.Fprintf(s, "%#v", _CallExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f CaseClause) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _CaseClause CaseClause
		fmt.Fprintf(s, "%#v", _CaseClause(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ClassDeclaration) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ClassDeclaration ClassDeclaration
		fmt.Fprintf(s, "%#v", _ClassDeclaration(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ClassElement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ClassElement ClassElement
		fmt.Fprintf(s, "%#v", _ClassElement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ClassElementName) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ClassElementName ClassElementName
		fmt.Fprintf(s, "%#v", _ClassElementName(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f CoalesceExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _CoalesceExpression CoalesceExpression
		fmt.Fprintf(s, "%#v", _CoalesceExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ConditionalExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ConditionalExpression ConditionalExpression
		fmt.Fprintf(s, "%#v", _ConditionalExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ParenthesizedExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ParenthesizedExpression ParenthesizedExpression
		fmt.Fprintf(s, "%#v", _ParenthesizedExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Declaration) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _Declaration Declaration
		fmt.Fprintf(s, "%#v", _Declaration(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f DestructuringAssignmentTarget) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _DestructuringAssignmentTarget DestructuringAssignmentTarget
		fmt.Fprintf(s, "%#v", _DestructuringAssignmentTarget(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f EqualityExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _EqualityExpression EqualityExpression
		fmt.Fprintf(s, "%#v", _EqualityExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ExponentiationExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ExponentiationExpression ExponentiationExpression
		fmt.Fprintf(s, "%#v", _ExponentiationExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ExportClause) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ExportClause ExportClause
		fmt.Fprintf(s, "%#v", _ExportClause(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ExportDeclaration) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ExportDeclaration ExportDeclaration
		fmt.Fprintf(s, "%#v", _ExportDeclaration(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ExportSpecifier) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ExportSpecifier ExportSpecifier
		fmt.Fprintf(s, "%#v", _ExportSpecifier(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Expression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _Expression Expression
		fmt.Fprintf(s, "%#v", _Expression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f FieldDefinition) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _FieldDefinition FieldDefinition
		fmt.Fprintf(s, "%#v", _FieldDefinition(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f FormalParameters) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _FormalParameters FormalParameters
		fmt.Fprintf(s, "%#v", _FormalParameters(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f FromClause) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _FromClause FromClause
		fmt.Fprintf(s, "%#v", _FromClause(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f FunctionDeclaration) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _FunctionDeclaration FunctionDeclaration
		fmt.Fprintf(s, "%#v", _FunctionDeclaration(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f IfStatement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _IfStatement IfStatement
		fmt.Fprintf(s, "%#v", _IfStatement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ImportClause) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ImportClause ImportClause
		fmt.Fprintf(s, "%#v", _ImportClause(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ImportDeclaration) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ImportDeclaration ImportDeclaration
		fmt.Fprintf(s, "%#v", _ImportDeclaration(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ImportSpecifier) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ImportSpecifier ImportSpecifier
		fmt.Fprintf(s, "%#v", _ImportSpecifier(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f IterationStatementDo) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _IterationStatementDo IterationStatementDo
		fmt.Fprintf(s, "%#v", _IterationStatementDo(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f IterationStatementFor) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _IterationStatementFor IterationStatementFor
		fmt.Fprintf(s, "%#v", _IterationStatementFor(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f IterationStatementWhile) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _IterationStatementWhile IterationStatementWhile
		fmt.Fprintf(s, "%#v", _IterationStatementWhile(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f LeftHandSideExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _LeftHandSideExpression LeftHandSideExpression
		fmt.Fprintf(s, "%#v", _LeftHandSideExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f LexicalBinding) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _LexicalBinding LexicalBinding
		fmt.Fprintf(s, "%#v", _LexicalBinding(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f LexicalDeclaration) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _LexicalDeclaration LexicalDeclaration
		fmt.Fprintf(s, "%#v", _LexicalDeclaration(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f LogicalANDExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _LogicalANDExpression LogicalANDExpression
		fmt.Fprintf(s, "%#v", _LogicalANDExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f LogicalORExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _LogicalORExpression LogicalORExpression
		fmt.Fprintf(s, "%#v", _LogicalORExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f MemberExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _MemberExpression MemberExpression
		fmt.Fprintf(s, "%#v", _MemberExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f MethodDefinition) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _MethodDefinition MethodDefinition
		fmt.Fprintf(s, "%#v", _MethodDefinition(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Module) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _Module Module
		fmt.Fprintf(s, "%#v", _Module(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ModuleItem) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ModuleItem ModuleItem
		fmt.Fprintf(s, "%#v", _ModuleItem(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f MultiplicativeExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _MultiplicativeExpression MultiplicativeExpression
		fmt.Fprintf(s, "%#v", _MultiplicativeExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f NamedImports) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _NamedImports NamedImports
		fmt.Fprintf(s, "%#v", _NamedImports(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f NewExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _NewExpression NewExpression
		fmt.Fprintf(s, "%#v", _NewExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ObjectAssignmentPattern) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ObjectAssignmentPattern ObjectAssignmentPattern
		fmt.Fprintf(s, "%#v", _ObjectAssignmentPattern(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ObjectBindingPattern) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ObjectBindingPattern ObjectBindingPattern
		fmt.Fprintf(s, "%#v", _ObjectBindingPattern(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ObjectLiteral) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ObjectLiteral ObjectLiteral
		fmt.Fprintf(s, "%#v", _ObjectLiteral(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f OptionalChain) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _OptionalChain OptionalChain
		fmt.Fprintf(s, "%#v", _OptionalChain(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f OptionalExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _OptionalExpression OptionalExpression
		fmt.Fprintf(s, "%#v", _OptionalExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f PrimaryExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _PrimaryExpression PrimaryExpression
		fmt.Fprintf(s, "%#v", _PrimaryExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f PropertyDefinition) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _PropertyDefinition PropertyDefinition
		fmt.Fprintf(s, "%#v", _PropertyDefinition(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f PropertyName) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _PropertyName PropertyName
		fmt.Fprintf(s, "%#v", _PropertyName(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f RelationalExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _RelationalExpression RelationalExpression
		fmt.Fprintf(s, "%#v", _RelationalExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Script) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _Script Script
		fmt.Fprintf(s, "%#v", _Script(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ShiftExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _ShiftExpression ShiftExpression
		fmt.Fprintf(s, "%#v", _ShiftExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Statement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _Statement Statement
		fmt.Fprintf(s, "%#v", _Statement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f StatementListItem) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _StatementListItem StatementListItem
		fmt.Fprintf(s, "%#v", _StatementListItem(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f SwitchStatement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _SwitchStatement SwitchStatement
		fmt.Fprintf(s, "%#v", _SwitchStatement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f TemplateLiteral) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _TemplateLiteral TemplateLiteral
		fmt.Fprintf(s, "%#v", _TemplateLiteral(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f TryStatement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _TryStatement TryStatement
		fmt.Fprintf(s, "%#v", _TryStatement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f UnaryExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _UnaryExpression UnaryExpression
		fmt.Fprintf(s, "%#v", _UnaryExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f UpdateExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _UpdateExpression UpdateExpression
		fmt.Fprintf(s, "%#v", _UpdateExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f VariableDeclaration) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _VariableDeclaration VariableDeclaration
		fmt.Fprintf(s, "%#v", _VariableDeclaration(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f VariableStatement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _VariableStatement VariableStatement
		fmt.Fprintf(s, "%#v", _VariableStatement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f WithStatement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type _WithStatement WithStatement
		fmt.Fprintf(s, "%#v", _WithStatement(f))
	} else {
		format(&f, s, v)
	}
}
