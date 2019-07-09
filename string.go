package javascript

import (
	"fmt"
	"io"
	"reflect"

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

func (t Token) Format(s fmt.State, v rune) {
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
	fmt.Fprintf(s, "Type: %s - Data: %q", typ, t.Data)
	if s.Flag('+') {
		fmt.Fprintf(s, " - Position: %d (%d: %d)", t.Pos, t.Line, t.LinePos)
	}
}

func (t Tokens) Format(s fmt.State, v rune) {
	formatArray(s, s.Flag('+'), reflect.ValueOf(t))
}

var (
	space          = []byte{' '}
	arrayOpen      = []byte{'['}
	arrayClose     = []byte{'\n', ']'}
	arrayOpenClose = []byte{'[', ']'}
	objectOpen     = []byte{'{'}
	objectClose    = []byte{'\n', '}'}
	pointer        = []byte{'*'}
)

func format(s fmt.State, v rune, f interface{}) {
	verbose := s.Flag('+')
	switch v {
	case 'v':
		v := reflect.ValueOf(f)
		t := v.Type()
		name := t.Name()
		io.WriteString(s, name)
		for v.Kind() == reflect.Ptr && !v.IsNil() {
			s.Write(pointer)
			v = v.Elem()
		}
		ip := indentPrinter{s}
		if k := v.Kind(); k == reflect.Slice || k == reflect.Array {
			if name != "" {
				s.Write(space)
			}
			formatArray(&ip, verbose, v)
		} else {
			t := v.Type()
			s.Write(objectOpen)
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				if f.PkgPath != "" {
					continue
				}
				if f.Name == "Tokens" {
					if verbose {
						ip.Printf("\nTokens: %+v", v.Field(i).Interface())
					}
				} else if k := f.Type.Kind(); k == reflect.Slice || k == reflect.Array {
					ip.Printf("\n%s: ", f.Name)
					formatArray(&ip, verbose, v.Field(i))
				} else if verbose {
					ip.Printf("\n%s: %+v", f.Name, v.Field(i).Interface())
				} else {
					vf := v.Field(i)
					switch k {
					case reflect.Map, reflect.Ptr, reflect.Slice:
						if vf.IsNil() {
							continue
						}
					case reflect.Bool:
						if !vf.Bool() {
							continue
						}
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						if vf.Uint() == 0 {
							continue
						}
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						if vf.Int() == 0 {
							continue
						}
					}
					ip.Printf("\n%s: %v", f.Name, v.Field(i).Interface())
				}
			}
			s.Write(objectClose)
		}
	case 's':
		if ps, ok := f.(interface{ printSource(io.Writer, bool) }); ok {
			ps.printSource(s, verbose)
		}
	}
}

func formatArray(ip io.Writer, verbose bool, v reflect.Value) {
	if v.Len() == 0 {
		ip.Write(arrayOpenClose)
		return
	}
	ip.Write(arrayOpen)
	ipp := indentPrinter{ip}
	for i := 0; i < v.Len(); i++ {
		p := v.Index(i)
		if verbose {
			ipp.Printf("\n%d: %+v", i, p.Interface())
		} else {
			ipp.Printf("\n%d: %v", i, p.Interface())
		}
	}
	ip.Write(arrayClose)
}

func (ft FunctionType) Format(s fmt.State, _ rune) {
	if s.Flag('+') {
		fmt.Fprintf(s, "%s (%d)", ft.String(), uint8(ft))
	} else {
		io.WriteString(s, ft.String())
	}
}

func (ft FunctionType) String() string {
	switch ft {
	case FunctionNormal:
		return "Normal"
	case FunctionGenerator:
		return "Generator"
	case FunctionAsync:
		return "Async"
	default:
		return "Unknown"
	}
}

func (mt MethodType) Format(s fmt.State, _ rune) {
	if s.Flag('+') {
		fmt.Fprintf(s, "%s (%d)", mt.String(), uint8(mt))
	} else {
		io.WriteString(s, mt.String())
	}
}

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
	case MethodStatic:
		return "MethodStatic"
	case MethodStaticGenerator:
		return "MethodStaticGenerator"
	case MethodStaticAsync:
		return "MethodStaticAsync"
	case MethodStaticAsyncGenerator:
		return "MethodStaticAsyncGenerator"
	case MethodStaticGetter:
		return "MethodStaticGetter"
	case MethodStaticSetter:
		return "MethodStaticSetter"
	default:
		return "Unknown"
	}
}

func (st StatementType) Format(s fmt.State, _ rune) {
	if s.Flag('+') {
		fmt.Fprintf(s, "%s (%d)", st.String(), uint8(st))
	} else {
		io.WriteString(s, st.String())
	}
}

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
		return "Unknown"
	}
}

func (ft ForType) Format(s fmt.State, _ rune) {
	if s.Flag('+') {
		fmt.Fprintf(s, "%s (%d)", ft.String(), uint8(ft))
	} else {
		io.WriteString(s, ft.String())
	}
}

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
		return "unknown"
	}
}

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
		return "unknown"
	}
}

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
		return "unknown"
	}
}

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
		return "unknown"
	}
}

func (a AdditiveOperator) String() string {
	switch a {
	case AdditiveNone:
		return ""
	case AdditiveAdd:
		return "+"
	case AdditiveMinus:
		return "-"
	default:
		return "unknown"
	}
}

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
		return "unknown"
	}
}

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
		return "unknown"
	}
}

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
		return "unknown"
	}
}

func (l LetOrConst) String() string {
	if l {
		return "Const"
	}
	return "Let"
}

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
	case AssignmentSignPropagatinRightShift:
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
	default:
		return "unknown"
	}
}

func (f ClassDeclaration) Format(s fmt.State, v rune)         { format(s, v, f) }
func (f MethodDefinition) Format(s fmt.State, v rune)         { format(s, v, f) }
func (f PropertyName) Format(s fmt.State, v rune)             { format(s, v, f) }
func (f ConditionalExpression) Format(s fmt.State, v rune)    { format(s, v, f) }
func (f LogicalORExpression) Format(s fmt.State, v rune)      { format(s, v, f) }
func (f LogicalANDExpression) Format(s fmt.State, v rune)     { format(s, v, f) }
func (f BitwiseORExpression) Format(s fmt.State, v rune)      { format(s, v, f) }
func (f BitwiseXORExpression) Format(s fmt.State, v rune)     { format(s, v, f) }
func (f BitwiseANDExpression) Format(s fmt.State, v rune)     { format(s, v, f) }
func (f EqualityExpression) Format(s fmt.State, v rune)       { format(s, v, f) }
func (f RelationalExpression) Format(s fmt.State, v rune)     { format(s, v, f) }
func (f ShiftExpression) Format(s fmt.State, v rune)          { format(s, v, f) }
func (f AdditiveExpression) Format(s fmt.State, v rune)       { format(s, v, f) }
func (f MultiplicativeExpression) Format(s fmt.State, v rune) { format(s, v, f) }
func (f ExponentiationExpression) Format(s fmt.State, v rune) { format(s, v, f) }
func (f UnaryExpression) Format(s fmt.State, v rune)          { format(s, v, f) }
func (f UpdateExpression) Format(s fmt.State, v rune)         { format(s, v, f) }
func (f AssignmentExpression) Format(s fmt.State, v rune)     { format(s, v, f) }
func (f LeftHandSideExpression) Format(s fmt.State, v rune)   { format(s, v, f) }
func (f Expression) Format(s fmt.State, v rune)               { format(s, v, f) }
func (f NewExpression) Format(s fmt.State, v rune)            { format(s, v, f) }
func (f MemberExpression) Format(s fmt.State, v rune)         { format(s, v, f) }
func (f PrimaryExpression) Format(s fmt.State, v rune)        { format(s, v, f) }
func (f Arguments) Format(s fmt.State, v rune)                { format(s, v, f) }
func (f CallExpression) Format(s fmt.State, v rune)           { format(s, v, f) }
func (f FunctionDeclaration) Format(s fmt.State, v rune)      { format(s, v, f) }
func (f FormalParameters) Format(s fmt.State, v rune)         { format(s, v, f) }
func (f BindingElement) Format(s fmt.State, v rune)           { format(s, v, f) }
func (f FunctionRestParameter) Format(s fmt.State, v rune)    { format(s, v, f) }
func (f Script) Format(s fmt.State, v rune)                   { format(s, v, f) }
func (f Declaration) Format(s fmt.State, v rune)              { format(s, v, f) }
func (f LexicalDeclaration) Format(s fmt.State, v rune)       { format(s, v, f) }
func (f LexicalBinding) Format(s fmt.State, v rune)           { format(s, v, f) }
func (f ArrayBindingPattern) Format(s fmt.State, v rune)      { format(s, v, f) }
func (f ObjectBindingPattern) Format(s fmt.State, v rune)     { format(s, v, f) }
func (f BindingProperty) Format(s fmt.State, v rune)          { format(s, v, f) }
func (f VariableDeclaration) Format(s fmt.State, v rune)      { format(s, v, f) }
func (f ArrayLiteral) Format(s fmt.State, v rune)             { format(s, v, f) }
func (f ObjectLiteral) Format(s fmt.State, v rune)            { format(s, v, f) }
func (f PropertyDefinition) Format(s fmt.State, v rune)       { format(s, v, f) }
func (f TemplateLiteral) Format(s fmt.State, v rune)          { format(s, v, f) }
func (f ArrowFunction) Format(s fmt.State, v rune)            { format(s, v, f) }
func (f Module) Format(s fmt.State, v rune)                   { format(s, v, f) }
func (f ModuleListItem) Format(s fmt.State, v rune)           { format(s, v, f) }
func (f ImportDeclaration) Format(s fmt.State, v rune)        { format(s, v, f) }
func (f ImportClause) Format(s fmt.State, v rune)             { format(s, v, f) }
func (f FromClause) Format(s fmt.State, v rune)               { format(s, v, f) }
func (f NamedImports) Format(s fmt.State, v rune)             { format(s, v, f) }
func (f ImportSpecifier) Format(s fmt.State, v rune)          { format(s, v, f) }
func (f ExportDeclaration) Format(s fmt.State, v rune)        { format(s, v, f) }
func (f ExportClause) Format(s fmt.State, v rune)             { format(s, v, f) }
func (f ExportSpecifier) Format(s fmt.State, v rune)          { format(s, v, f) }
func (f Block) Format(s fmt.State, v rune)                    { format(s, v, f) }
func (f StatementListItem) Format(s fmt.State, v rune)        { format(s, v, f) }
func (f Statement) Format(s fmt.State, v rune)                { format(s, v, f) }
func (f IfStatement) Format(s fmt.State, v rune)              { format(s, v, f) }
func (f IterationStatementDo) Format(s fmt.State, v rune)     { format(s, v, f) }
func (f IterationStatementWhile) Format(s fmt.State, v rune)  { format(s, v, f) }
func (f IterationStatementFor) Format(s fmt.State, v rune)    { format(s, v, f) }
func (f SwitchStatement) Format(s fmt.State, v rune)          { format(s, v, f) }
func (f CaseClause) Format(s fmt.State, v rune)               { format(s, v, f) }
func (f WithStatement) Format(s fmt.State, v rune)            { format(s, v, f) }
func (f TryStatement) Format(s fmt.State, v rune)             { format(s, v, f) }
func (f VariableStatement) Format(s fmt.State, v rune)        { format(s, v, f) }
func (f CoverParenthesizedExpressionAndArrowParameterList) Format(s fmt.State, v rune) {
	format(s, v, f)
}

func (t *Token) String() string {
	if t == nil {
		return ""
	}
	return t.Data
}
