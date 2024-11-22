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
		io.WriteString(w, "[]")

		return
	}

	io.WriteString(w, "[")

	ipp := indentPrinter{w}

	for n, t := range t {
		ipp.Printf("\n%d: ", n)
		t.printType(w, v)
	}

	io.WriteString(w, "\n]")
}

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

// String implements the fmt.Stringer interface
func (ft FunctionType) String() string {
	switch ft {
	case FunctionNormal:
		return "Normal"
	case FunctionGenerator:
		return "Generator"
	case FunctionAsync:
		return "Async"
	case FunctionAsyncGenerator:
		return "Async Generator"
	default:
		return unknown
	}
}

func (ft FunctionType) printType(w io.Writer, _ bool) {
	io.WriteString(w, ft.String())
}

// String implements the fmt.Stringer interface
func (mt MethodType) String() string {
	switch mt {
	case MethodNormal:
		return "MethodNormal"
	case MethodGenerator:
		return "MethodGenerator"
	case MethodAsyncGenerator:
		return "MethodAsyncGenerator"
	case MethodAsync:
		return "MethodAsync"
	case MethodGetter:
		return "MethodGetter"
	case MethodSetter:
		return "MethodSetter"
	default:
		return unknown
	}
}

func (mt MethodType) printType(w io.Writer, _ bool) {
	io.WriteString(w, mt.String())
}

// String implements the fmt.Stringer interface
func (st StatementType) String() string {
	switch st {
	case StatementNormal:
		return "StatementNormal"
	case StatementContinue:
		return "StatementContinue"
	case StatementBreak:
		return "StatementBreak"
	case StatementReturn:
		return "StatementReturn"
	case StatementThrow:
		return "StatementThrow"
	default:
		return unknown
	}
}

func (st StatementType) printType(w io.Writer, _ bool) {
	io.WriteString(w, st.String())
}

// String implements the fmt.Stringer interface
func (ft ForType) String() string {
	switch ft {
	case ForNormal:
		return "ForNormal"
	case ForNormalVar:
		return "ForNormalVar"
	case ForNormalLexicalDeclaration:
		return "ForNormalLexicalDeclaration"
	case ForNormalExpression:
		return "ForNormalExpression"
	case ForInLeftHandSide:
		return "ForInLeftHandSide"
	case ForInVar:
		return "ForInVar"
	case ForInLet:
		return "ForInLet"
	case ForInConst:
		return "ForInConst"
	case ForOfLeftHandSide:
		return "ForOfLeftHandSide"
	case ForOfVar:
		return "ForOfVar"
	case ForOfLet:
		return "ForOfLet"
	case ForOfConst:
		return "ForOfConst"
	case ForAwaitOfLeftHandSide:
		return "ForAwaitOfLeftHandSide"
	case ForAwaitOfVar:
		return "ForAwaitOfVar"
	case ForAwaitOfLet:
		return "ForAwaitOfLet"
	case ForAwaitOfConst:
		return "ForAwaitOfConst"
	default:
		return unknown
	}
}

func (ft ForType) printType(w io.Writer, _ bool) {
	io.WriteString(w, ft.String())
}

// String implements the fmt.Stringer interface
func (e EqualityOperator) String() string {
	switch e {
	case EqualityNone:
		return ""
	case EqualityEqual:
		return "=="
	case EqualityNotEqual:
		return "!="
	case EqualityStrictEqual:
		return "==="
	case EqualityStrictNotEqual:
		return "!=="
	default:
		return unknown
	}
}

func (e EqualityOperator) printType(w io.Writer, _ bool) {
	io.WriteString(w, e.String())
}

// String implements the fmt.Stringer interface
func (r RelationshipOperator) String() string {
	switch r {
	case RelationshipNone:
		return ""
	case RelationshipLessThan:
		return "<"
	case RelationshipGreaterThan:
		return ">"
	case RelationshipLessThanEqual:
		return "<="
	case RelationshipGreaterThanEqual:
		return ">="
	case RelationshipInstanceOf:
		return "instanceof"
	case RelationshipIn:
		return "in"
	default:
		return unknown
	}
}

func (r RelationshipOperator) printType(w io.Writer, _ bool) {
	io.WriteString(w, r.String())
}

// String implements the fmt.Stringer interface
func (s ShiftOperator) String() string {
	switch s {
	case ShiftNone:
		return ""
	case ShiftLeft:
		return "<<"
	case ShiftRight:
		return ">>"
	case ShiftUnsignedRight:
		return ">>>"
	default:
		return unknown
	}
}

func (s ShiftOperator) printType(w io.Writer, _ bool) {
	io.WriteString(w, s.String())
}

// String implements the fmt.Stringer interface
func (a AdditiveOperator) String() string {
	switch a {
	case AdditiveNone:
		return ""
	case AdditiveAdd:
		return "+"
	case AdditiveMinus:
		return "-"
	default:
		return unknown
	}
}

func (a AdditiveOperator) printType(w io.Writer, _ bool) {
	io.WriteString(w, a.String())
}

// String implements the fmt.Stringer interface
func (m MultiplicativeOperator) String() string {
	switch m {
	case MultiplicativeNone:
		return ""
	case MultiplicativeMultiply:
		return "*"
	case MultiplicativeDivide:
		return "/"
	case MultiplicativeRemainder:
		return "%"
	default:
		return unknown
	}
}

func (m MultiplicativeOperator) printType(w io.Writer, _ bool) {
	io.WriteString(w, m.String())
}

// String implements the fmt.Stringer interface
func (u UnaryOperator) String() string {
	switch u {
	case UnaryNone:
		return ""
	case UnaryDelete:
		return "delete"
	case UnaryVoid:
		return "void"
	case UnaryTypeOf:
		return "typeof"
	case UnaryAdd:
		return "+"
	case UnaryMinus:
		return "-"
	case UnaryBitwiseNot:
		return "~"
	case UnaryLogicalNot:
		return "!"
	case UnaryAwait:
		return "await"
	default:
		return unknown
	}
}

func (u UnaryOperator) printType(w io.Writer, _ bool) {
	io.WriteString(w, u.String())
}

// String implements the fmt.Stringer interface
func (u UpdateOperator) String() string {
	switch u {
	case UpdateNone:
		return ""
	case UpdatePostIncrement:
		return " ++"
	case UpdatePostDecrement:
		return " --"
	case UpdatePreIncrement:
		return "++"
	case UpdatePreDecrement:
		return "--"
	default:
		return unknown
	}
}

func (u UpdateOperator) printType(w io.Writer, _ bool) {
	io.WriteString(w, u.String())
}

// String implements the fmt.Stringer interface
func (l LetOrConst) String() string {
	if l {
		return "Const"
	}

	return "Let"
}

func (l LetOrConst) printType(w io.Writer, _ bool) {
	io.WriteString(w, l.String())
}

// String implements the fmt.Stringer interface
func (a AssignmentOperator) String() string {
	switch a {
	case AssignmentNone:
		return ""
	case AssignmentAssign:
		return "="
	case AssignmentMultiply:
		return "*="
	case AssignmentDivide:
		return "/="
	case AssignmentRemainder:
		return "%="
	case AssignmentAdd:
		return "+="
	case AssignmentSubtract:
		return "-="
	case AssignmentLeftShift:
		return "<<="
	case AssignmentSignPropagatingRightShift:
		return ">>="
	case AssignmentZeroFillRightShift:
		return ">>>="
	case AssignmentBitwiseAND:
		return "&="
	case AssignmentBitwiseXOR:
		return "^="
	case AssignmentBitwiseOR:
		return "|="
	case AssignmentExponentiation:
		return "**="
	case AssignmentLogicalAnd:
		return "&&="
	case AssignmentLogicalOr:
		return "||="
	case AssignmentNullish:
		return "??="
	default:
		return unknown
	}
}

func (a AssignmentOperator) printType(w io.Writer, _ bool) {
	io.WriteString(w, a.String())
}

const unknown = "Unknown"
