package javascript

import (
	"fmt"
	"io"
	"strings"
	"unsafe"

	"vimagination.zapto.org/parser"
)

var (
	indent    = []byte{'\t'}
	semiColon = []byte{';'}
)

type writer interface {
	io.Writer
	WriteString(string)
	Underlying() writer
	PrintSemiColon()
	LastChar() byte
	LastIsWhitespace() bool
	Pos() int
	Indent() writer
	Printf(string, ...any)
}

type indentPrinter struct {
	writer
	hadNewline bool
}

func (i *indentPrinter) Write(p []byte) (int, error) {
	var (
		total int
		last  int
	)

	for n, c := range p {
		if c == '\n' {
			if last != n {
				if err := i.printIndent(); err != nil {
					return total, err
				}
			}

			m, err := i.writer.Write(p[last : n+1])
			total += m

			if err != nil {
				return total, err
			}

			i.hadNewline = true
			last = n + 1
		}
	}

	if last != len(p) {
		if err := i.printIndent(); err != nil {
			return total, err
		}

		m, err := i.writer.Write(p[last:])
		total += m

		if err != nil {
			return total, err
		}
	}

	return total, nil
}

func (i *indentPrinter) printIndent() error {
	if i.hadNewline {
		if _, err := i.writer.Write(indent); err != nil {
			return err
		}

		i.hadNewline = false
	}

	return nil
}

func (i *indentPrinter) Printf(format string, args ...any) {
	fmt.Fprintf(i, format, args...)
}

func (i *indentPrinter) WriteString(s string) {
	i.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
}

func (i *indentPrinter) Indent() writer {
	return &indentPrinter{writer: i}
}

type underlyingWriter struct {
	io.Writer
	lastChar byte
	pos      int
}

func (u *underlyingWriter) Write(p []byte) (int, error) {
	for _, b := range p {
		if b == '\n' {
			u.pos = 0
		} else if b != '\t' || u.pos > 0 {
			u.pos++
		}

		u.lastChar = b
	}

	return u.Writer.Write(p)
}

func (u *underlyingWriter) WriteString(s string) {
	u.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
}

func (u *underlyingWriter) Underlying() writer {
	return u
}

func (u *underlyingWriter) PrintSemiColon() {
	if u.lastChar != '\n' {
		u.Writer.Write(semiColon)
		u.pos++
	}
}

func (u *underlyingWriter) LastChar() byte {
	return u.lastChar
}

func (u *underlyingWriter) LastIsWhitespace() bool {
	return u.lastChar == ' ' || u.lastChar == '\n' || u.lastChar == '\t'
}

func (u *underlyingWriter) Pos() int {
	return u.pos
}

func (u *underlyingWriter) Indent() writer {
	return &indentPrinter{writer: u}
}

func (u *underlyingWriter) Printf(format string, args ...any) {
	fmt.Fprintf(u, format, args...)
}

// Format implements the fmt.Formatter interface
func (t Token) Format(s fmt.State, v rune) {
	t.printType(s, s.Flag('+'))
}

func (t Token) printType(w io.Writer, v bool) {
	var typ string

	if t.Type&tokenTypescript != 0 {
		typ = "Typescript"
	}

	switch t.Type &^ tokenTypescript {
	case parser.TokenError:
		typ += "Error"
	case parser.TokenDone:
		typ += "Done"
	case TokenWhitespace:
		typ += "Whitespace"
	case TokenLineTerminator:
		typ += "LineTerminator"
	case TokenSingleLineComment:
		typ += "SingleLineComment"
	case TokenMultiLineComment:
		typ += "MultiLineComment"
	case TokenIdentifier:
		typ += "Identifier"
	case TokenBooleanLiteral:
		typ += "BooleanLiteral"
	case TokenKeyword:
		typ += "Keyword"
	case TokenPunctuator:
		typ += "Punctuator"
	case TokenNumericLiteral:
		typ += "NumericLiteral"
	case TokenStringLiteral:
		typ += "StringLiteral"
	case TokenNoSubstitutionTemplate:
		typ += "NoSubstitutionTemplate"
	case TokenTemplateHead:
		typ += "TemplateHead"
	case TokenTemplateMiddle:
		typ += "TemplateMiddle"
	case TokenTemplateTail:
		typ += "TemplateTail"
	case TokenDivPunctuator:
		typ += "DivPunctuator"
	case TokenRightBracePunctuator:
		typ += "RightBracePunctuator"
	case TokenRegularExpressionLiteral:
		typ += "RegulatExpressionLiteral"
	case TokenNullLiteral:
		typ += "NullLiteral"
	case TokenFutureReservedWord:
		typ += "FutureReservedWord"
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
	t.printType(&underlyingWriter{Writer: s}, s.Flag('+'))
}

func (t Tokens) printType(w writer, v bool) {
	if t == nil {
		w.WriteString("nil")

		return
	} else if len(t) == 0 {
		w.WriteString("[]")

		return
	}

	w.WriteString("[")

	ipp := w.Indent()

	for n, t := range t {
		ipp.Printf("\n%d: ", n)
		t.printType(w, v)
	}

	w.WriteString("\n]")
}

func (c Comments) printType(w writer, v bool) {
	if c == nil {
		w.WriteString("nil")

		return
	} else if len(c) == 0 {
		w.WriteString("[]")

		return
	}

	w.WriteString("[")

	ipp := w.Indent()

	for n, t := range c {
		ipp.Printf("\n%d: ", n)

		if t == nil {
			w.WriteString("nil")
		} else {
			t.printType(w, v)
		}
	}

	w.WriteString("\n]")
}

type commentPrinter bool

func (c Comments) printSource(w writer, postSpace, postNewline bool) {
	for len(c) > 0 && c[0] == nil {
		c = c[1:]
	}

	if len(c) > 0 {
		switch w.LastChar() {
		case 0, ' ', '\n', '\t':
		default:
			w.WriteString(" ")
		}

		line := c[0].Line + uint64(strings.Count(c[0].Data, "\n"))
		pos := w.Pos()

		var cp commentPrinter

		lastWasMulti := cp.print(w, *c[0], 0)

		if !lastWasMulti {
			line++
		}

		for _, c := range c[1:] {
			if c == nil {
				continue
			}

			if !cp {
				if line < c.Line {
					if !lastWasMulti {
						w.WriteString("\n")
					}

					w.WriteString("\n")

					line++
					pos = 0
				} else if lastWasMulti {
					w.WriteString(" ")
				} else {
					w.WriteString("\n")
				}
			}

			if lastWasMulti = cp.print(w, *c, pos); lastWasMulti {
				line += uint64(strings.Count(c.Data, "\n"))
			} else {
				line++
			}
		}

		if cp {
			w.WriteString("*/")
		}

		if postNewline || !lastWasMulti {
			w.WriteString("\n")
		} else if postSpace {
			w.WriteString(" ")
		}
	}
}

func (c Comments) LastIsMulti() bool {
	return len(c) > 0 && c[len(c)-1].Type == TokenMultiLineComment
}

func (cp *commentPrinter) print(w writer, c Token, pos int) bool {
	var multi bool

	if isSingleLine(c) {
		if *cp {
			w.WriteString("*/ ")

			pos = max(0, pos-3)
			*cp = false
		}

		w.WriteString(strings.Repeat(" ", pos))

		if !strings.HasPrefix(c.Data, "//") {
			w.WriteString("// ")
		}
	} else {
		multi = true

		if c.Type == TokenMultiLineComment {
			if *cp {
				w.WriteString("*/ ")
			}

			*cp = false
		} else if !bool(*cp) {
			w.WriteString("/*")
			*cp = true
		}
	}

	if *cp {
		w.WriteString(strings.ReplaceAll(c.Data, "*/", "* /"))
	} else {
		w.WriteString(strings.TrimSpace(c.Data))
	}

	return multi
}

func isSingleLine(c Token) bool {
	return c.Type == TokenSingleLineComment && !strings.Contains(c.Data, "\n")
}

type formatter interface {
	printType(writer, bool)
	printSource(writer, bool)
}

func format(f formatter, s fmt.State, v rune) {
	switch v {
	case 'v':
		f.printType(&underlyingWriter{Writer: s}, s.Flag('+'))
	case 's':
		f.printSource(&underlyingWriter{Writer: s}, s.Flag('+'))
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

func (ft FunctionType) printType(w writer, _ bool) {
	w.WriteString(ft.String())
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

func (mt MethodType) printType(w writer, _ bool) {
	w.WriteString(mt.String())
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
	case StatementDebugger:
		return "StatementDebugger"
	default:
		return unknown
	}
}

func (st StatementType) printType(w writer, _ bool) {
	w.WriteString(st.String())
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

func (ft ForType) printType(w writer, _ bool) {
	w.WriteString(ft.String())
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

func (e EqualityOperator) printType(w writer, _ bool) {
	w.WriteString(e.String())
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

func (r RelationshipOperator) printType(w writer, _ bool) {
	w.WriteString(r.String())
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

func (s ShiftOperator) printType(w writer, _ bool) {
	w.WriteString(s.String())
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

func (a AdditiveOperator) printType(w writer, _ bool) {
	w.WriteString(a.String())
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

func (m MultiplicativeOperator) printType(w writer, _ bool) {
	w.WriteString(m.String())
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

func (u UnaryOperator) printType(w writer, _ bool) {
	w.WriteString(u.String())
}

// String implements the fmt.Stringer interface
func (u UpdateOperator) String() string {
	switch u {
	case UpdateNone:
		return ""
	case UpdatePostIncrement:
		return "++"
	case UpdatePostDecrement:
		return "--"
	case UpdatePreIncrement:
		return "++"
	case UpdatePreDecrement:
		return "--"
	default:
		return unknown
	}
}

func (u UpdateOperator) printType(w writer, _ bool) {
	w.WriteString(u.String())
}

// String implements the fmt.Stringer interface
func (l LetOrConst) String() string {
	if l {
		return "Const"
	}

	return "Let"
}

func (l LetOrConst) printType(w writer, _ bool) {
	w.WriteString(l.String())
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

func (a AssignmentOperator) printType(w writer, _ bool) {
	w.WriteString(a.String())
}

const unknown = "Unknown"
