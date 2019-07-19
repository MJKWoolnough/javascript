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
	pointer        = []byte{'*'}
	nilStr         = []byte{'<', 'n', 'i', 'l', '>'}
	to             = []byte{':', ' '}
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
func (f AdditiveExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f Arguments) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ArrayBindingPattern) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ArrayLiteral) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ArrowFunction) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f AssignmentExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f BindingElement) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f BindingProperty) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f BitwiseANDExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f BitwiseORExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f BitwiseXORExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f Block) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f CallExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f CaseClause) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ClassDeclaration) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ConditionalExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f CoverParenthesizedExpressionAndArrowParameterList) Format(s fmt.State, v rune) {
	format(&f, s, v)
}

// Format implements the fmt.Formatter interface
func (f Declaration) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f EqualityExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ExponentiationExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ExportClause) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ExportDeclaration) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ExportSpecifier) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f Expression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f FormalParameters) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f FromClause) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f FunctionDeclaration) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f FunctionRestParameter) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f IfStatement) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ImportClause) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ImportDeclaration) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ImportSpecifier) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f IterationStatementDo) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f IterationStatementFor) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f IterationStatementWhile) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f LeftHandSideExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f LexicalBinding) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f LexicalDeclaration) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f LogicalANDExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f LogicalORExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f MemberExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f MethodDefinition) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f Module) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ModuleItem) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f MultiplicativeExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f NamedImports) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f NewExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ObjectBindingPattern) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ObjectLiteral) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f PrimaryExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f PropertyDefinition) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f PropertyName) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f RelationalExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f Script) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f ShiftExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f Statement) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f StatementListItem) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f SwitchStatement) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f TemplateLiteral) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f TryStatement) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f UnaryExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f UpdateExpression) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f VariableDeclaration) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f VariableStatement) Format(s fmt.State, v rune) { format(&f, s, v) }

// Format implements the fmt.Formatter interface
func (f WithStatement) Format(s fmt.State, v rune) { format(&f, s, v) }
