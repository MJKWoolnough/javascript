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
		type X = AdditiveExpression
		type AdditiveExpression X
		fmt.Fprintf(s, "%#v", AdditiveExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Argument) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Argument
		type Argument X
		fmt.Fprintf(s, "%#v", Argument(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Arguments) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Arguments
		type Arguments X
		fmt.Fprintf(s, "%#v", Arguments(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ArrayAssignmentPattern) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ArrayAssignmentPattern
		type ArrayAssignmentPattern X
		fmt.Fprintf(s, "%#v", ArrayAssignmentPattern(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ArrayBindingPattern) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ArrayBindingPattern
		type ArrayBindingPattern X
		fmt.Fprintf(s, "%#v", ArrayBindingPattern(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ArrayElement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ArrayElement
		type ArrayElement X
		fmt.Fprintf(s, "%#v", ArrayElement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ArrayLiteral) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ArrayLiteral
		type ArrayLiteral X
		fmt.Fprintf(s, "%#v", ArrayLiteral(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ArrowFunction) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ArrowFunction
		type ArrowFunction X
		fmt.Fprintf(s, "%#v", ArrowFunction(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f AssignmentElement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = AssignmentElement
		type AssignmentElement X
		fmt.Fprintf(s, "%#v", AssignmentElement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f AssignmentExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = AssignmentExpression
		type AssignmentExpression X
		fmt.Fprintf(s, "%#v", AssignmentExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f AssignmentPattern) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = AssignmentPattern
		type AssignmentPattern X
		fmt.Fprintf(s, "%#v", AssignmentPattern(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f AssignmentProperty) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = AssignmentProperty
		type AssignmentProperty X
		fmt.Fprintf(s, "%#v", AssignmentProperty(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f BindingElement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = BindingElement
		type BindingElement X
		fmt.Fprintf(s, "%#v", BindingElement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f BindingProperty) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = BindingProperty
		type BindingProperty X
		fmt.Fprintf(s, "%#v", BindingProperty(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f BitwiseANDExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = BitwiseANDExpression
		type BitwiseANDExpression X
		fmt.Fprintf(s, "%#v", BitwiseANDExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f BitwiseORExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = BitwiseORExpression
		type BitwiseORExpression X
		fmt.Fprintf(s, "%#v", BitwiseORExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f BitwiseXORExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = BitwiseXORExpression
		type BitwiseXORExpression X
		fmt.Fprintf(s, "%#v", BitwiseXORExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Block) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Block
		type Block X
		fmt.Fprintf(s, "%#v", Block(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f CallExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = CallExpression
		type CallExpression X
		fmt.Fprintf(s, "%#v", CallExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f CaseClause) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = CaseClause
		type CaseClause X
		fmt.Fprintf(s, "%#v", CaseClause(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ClassDeclaration) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ClassDeclaration
		type ClassDeclaration X
		fmt.Fprintf(s, "%#v", ClassDeclaration(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ClassElement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ClassElement
		type ClassElement X
		fmt.Fprintf(s, "%#v", ClassElement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ClassElementName) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ClassElementName
		type ClassElementName X
		fmt.Fprintf(s, "%#v", ClassElementName(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f CoalesceExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = CoalesceExpression
		type CoalesceExpression X
		fmt.Fprintf(s, "%#v", CoalesceExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ConditionalExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ConditionalExpression
		type ConditionalExpression X
		fmt.Fprintf(s, "%#v", ConditionalExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ParenthesizedExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ParenthesizedExpression
		type ParenthesizedExpression X
		fmt.Fprintf(s, "%#v", ParenthesizedExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Declaration) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Declaration
		type Declaration X
		fmt.Fprintf(s, "%#v", Declaration(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f DestructuringAssignmentTarget) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = DestructuringAssignmentTarget
		type DestructuringAssignmentTarget X
		fmt.Fprintf(s, "%#v", DestructuringAssignmentTarget(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f EqualityExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = EqualityExpression
		type EqualityExpression X
		fmt.Fprintf(s, "%#v", EqualityExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ExponentiationExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ExponentiationExpression
		type ExponentiationExpression X
		fmt.Fprintf(s, "%#v", ExponentiationExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ExportClause) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ExportClause
		type ExportClause X
		fmt.Fprintf(s, "%#v", ExportClause(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ExportDeclaration) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ExportDeclaration
		type ExportDeclaration X
		fmt.Fprintf(s, "%#v", ExportDeclaration(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ExportSpecifier) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ExportSpecifier
		type ExportSpecifier X
		fmt.Fprintf(s, "%#v", ExportSpecifier(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Expression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Expression
		type Expression X
		fmt.Fprintf(s, "%#v", Expression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f FieldDefinition) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = FieldDefinition
		type FieldDefinition X
		fmt.Fprintf(s, "%#v", FieldDefinition(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f FormalParameters) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = FormalParameters
		type FormalParameters X
		fmt.Fprintf(s, "%#v", FormalParameters(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f FromClause) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = FromClause
		type FromClause X
		fmt.Fprintf(s, "%#v", FromClause(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f FunctionDeclaration) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = FunctionDeclaration
		type FunctionDeclaration X
		fmt.Fprintf(s, "%#v", FunctionDeclaration(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f IfStatement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = IfStatement
		type IfStatement X
		fmt.Fprintf(s, "%#v", IfStatement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ImportClause) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ImportClause
		type ImportClause X
		fmt.Fprintf(s, "%#v", ImportClause(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ImportDeclaration) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ImportDeclaration
		type ImportDeclaration X
		fmt.Fprintf(s, "%#v", ImportDeclaration(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ImportSpecifier) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ImportSpecifier
		type ImportSpecifier X
		fmt.Fprintf(s, "%#v", ImportSpecifier(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f IterationStatementDo) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = IterationStatementDo
		type IterationStatementDo X
		fmt.Fprintf(s, "%#v", IterationStatementDo(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f IterationStatementFor) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = IterationStatementFor
		type IterationStatementFor X
		fmt.Fprintf(s, "%#v", IterationStatementFor(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f IterationStatementWhile) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = IterationStatementWhile
		type IterationStatementWhile X
		fmt.Fprintf(s, "%#v", IterationStatementWhile(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f LeftHandSideExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = LeftHandSideExpression
		type LeftHandSideExpression X
		fmt.Fprintf(s, "%#v", LeftHandSideExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f LexicalBinding) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = LexicalBinding
		type LexicalBinding X
		fmt.Fprintf(s, "%#v", LexicalBinding(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f LexicalDeclaration) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = LexicalDeclaration
		type LexicalDeclaration X
		fmt.Fprintf(s, "%#v", LexicalDeclaration(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f LogicalANDExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = LogicalANDExpression
		type LogicalANDExpression X
		fmt.Fprintf(s, "%#v", LogicalANDExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f LogicalORExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = LogicalORExpression
		type LogicalORExpression X
		fmt.Fprintf(s, "%#v", LogicalORExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f MemberExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = MemberExpression
		type MemberExpression X
		fmt.Fprintf(s, "%#v", MemberExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f MethodDefinition) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = MethodDefinition
		type MethodDefinition X
		fmt.Fprintf(s, "%#v", MethodDefinition(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Module) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Module
		type Module X
		fmt.Fprintf(s, "%#v", Module(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ModuleItem) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ModuleItem
		type ModuleItem X
		fmt.Fprintf(s, "%#v", ModuleItem(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f MultiplicativeExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = MultiplicativeExpression
		type MultiplicativeExpression X
		fmt.Fprintf(s, "%#v", MultiplicativeExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f NamedImports) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = NamedImports
		type NamedImports X
		fmt.Fprintf(s, "%#v", NamedImports(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f NewExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = NewExpression
		type NewExpression X
		fmt.Fprintf(s, "%#v", NewExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ObjectAssignmentPattern) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ObjectAssignmentPattern
		type ObjectAssignmentPattern X
		fmt.Fprintf(s, "%#v", ObjectAssignmentPattern(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ObjectBindingPattern) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ObjectBindingPattern
		type ObjectBindingPattern X
		fmt.Fprintf(s, "%#v", ObjectBindingPattern(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ObjectLiteral) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ObjectLiteral
		type ObjectLiteral X
		fmt.Fprintf(s, "%#v", ObjectLiteral(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f OptionalChain) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = OptionalChain
		type OptionalChain X
		fmt.Fprintf(s, "%#v", OptionalChain(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f OptionalExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = OptionalExpression
		type OptionalExpression X
		fmt.Fprintf(s, "%#v", OptionalExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f PrimaryExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = PrimaryExpression
		type PrimaryExpression X
		fmt.Fprintf(s, "%#v", PrimaryExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f PropertyDefinition) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = PropertyDefinition
		type PropertyDefinition X
		fmt.Fprintf(s, "%#v", PropertyDefinition(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f PropertyName) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = PropertyName
		type PropertyName X
		fmt.Fprintf(s, "%#v", PropertyName(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f RelationalExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = RelationalExpression
		type RelationalExpression X
		fmt.Fprintf(s, "%#v", RelationalExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Script) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Script
		type Script X
		fmt.Fprintf(s, "%#v", Script(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ShiftExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ShiftExpression
		type ShiftExpression X
		fmt.Fprintf(s, "%#v", ShiftExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Statement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Statement
		type Statement X
		fmt.Fprintf(s, "%#v", Statement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f StatementListItem) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = StatementListItem
		type StatementListItem X
		fmt.Fprintf(s, "%#v", StatementListItem(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f SwitchStatement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = SwitchStatement
		type SwitchStatement X
		fmt.Fprintf(s, "%#v", SwitchStatement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f TemplateLiteral) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = TemplateLiteral
		type TemplateLiteral X
		fmt.Fprintf(s, "%#v", TemplateLiteral(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f TryStatement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = TryStatement
		type TryStatement X
		fmt.Fprintf(s, "%#v", TryStatement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f UnaryExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = UnaryExpression
		type UnaryExpression X
		fmt.Fprintf(s, "%#v", UnaryExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f UpdateExpression) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = UpdateExpression
		type UpdateExpression X
		fmt.Fprintf(s, "%#v", UpdateExpression(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f VariableStatement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = VariableStatement
		type VariableStatement X
		fmt.Fprintf(s, "%#v", VariableStatement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f WithStatement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = WithStatement
		type WithStatement X
		fmt.Fprintf(s, "%#v", WithStatement(f))
	} else {
		format(&f, s, v)
	}
}
