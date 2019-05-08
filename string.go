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
	format(s, v, t)
}

var (
	space       = []byte{' '}
	arrayOpen   = []byte{'['}
	arrayClose  = []byte{'\n', ']'}
	objectOpen  = []byte{'{'}
	objectClose = []byte{'\n', '}'}
)

func format(s fmt.State, v rune, f interface{}) {
	verbose := s.Flag('+')
	if v == 'v' {
		v := reflect.ValueOf(f)
		t := v.Type()
		name := t.Name()
		io.WriteString(s, name)
		ip := indentPrinter{s}
		for v.Kind() == reflect.Ptr && !v.IsNil() {
			v = v.Elem()
		}
		if k := v.Kind(); k == reflect.Slice || k == reflect.Array {
			if name != "" {
				s.Write(space)
			}
			s.Write(arrayOpen)
			ipp := indentPrinter{&ip}
			for i := 0; i < v.Len(); i++ {
				p := v.Index(i)
				if verbose {
					ipp.Printf("\n%d: %+v", i, p.Interface())
				} else {
					ipp.Printf("\n%d: %v", i, p.Interface())
				}
			}
			s.Write(arrayClose)
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
						ip.Printf("\n%+v", v.Field(i).Interface())
					}
				} else if verbose {
					ip.Printf("\n%s: %+v", f.Name, v.Field(i).Interface())
				} else {
					ip.Printf("\n%s: %v", f.Name, v.Field(i).Interface())
				}
			}
			s.Write(objectClose)
		}
	}
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
	case MethodAsync:
		return "MethodAsync"
	case MethodGetter:
		return "MethodGetter"
	case MethodSetter:
		return "MethodSetter"
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
		return "ForInConst:"
	case ForOfLeftHandSide:
		return "ForOfLeftHandSide"
	case ForOfVar:
		return "ForOfVar"
	case ForOfLet:
		return "ForOfLet"
	case ForOfConst:
		return "ForOfConst"
	default:
		return "unknown"
	}
}

func (f ClassDeclaration) Format(s fmt.State, v rune)         { format(s, v, f) }
func (f ClassBody) Format(s fmt.State, v rune)                { format(s, v, f) }
func (f MethodDefinition) Format(s fmt.State, v rune)         { format(s, v, f) }
func (f PropertyName) Format(s fmt.State, v rune)             { format(s, v, f) }
func (f ConditionalExpression) Format(s fmt.State, v rune)    { format(s, v, f) }
func (f LogicalORExpression) Format(s fmt.State, v rune)      { format(s, v, f) }
func (f LogicalANDExpression) Format(s fmt.State, v rune)     { format(s, v, f) }
func (f BitwiseORExpression) Format(s fmt.State, v rune)      { format(s, v, f) }
func (f BitwiseXORExpression) Format(s fmt.State, v rune)     { format(s, v, f) }
func (f BitwiseANDExpression) Format(s fmt.State, v rune)     { format(s, v, f) }
func (f EqualityOperator) Format(s fmt.State, v rune)         { format(s, v, f) }
func (f EqualityExpression) Format(s fmt.State, v rune)       { format(s, v, f) }
func (f RelationshipOperator) Format(s fmt.State, v rune)     { format(s, v, f) }
func (f RelationalExpression) Format(s fmt.State, v rune)     { format(s, v, f) }
func (f ShiftOperator) Format(s fmt.State, v rune)            { format(s, v, f) }
func (f ShiftExpression) Format(s fmt.State, v rune)          { format(s, v, f) }
func (f AdditiveOperator) Format(s fmt.State, v rune)         { format(s, v, f) }
func (f AdditiveExpression) Format(s fmt.State, v rune)       { format(s, v, f) }
func (f MultiplicativeOperator) Format(s fmt.State, v rune)   { format(s, v, f) }
func (f MultiplicativeExpression) Format(s fmt.State, v rune) { format(s, v, f) }
func (f ExponentiationExpression) Format(s fmt.State, v rune) { format(s, v, f) }
func (f UnaryOperator) Format(s fmt.State, v rune)            { format(s, v, f) }
func (f UnaryExpression) Format(s fmt.State, v rune)          { format(s, v, f) }
func (f UpdateOperator) Format(s fmt.State, v rune)           { format(s, v, f) }
func (f UpdateExpression) Format(s fmt.State, v rune)         { format(s, v, f) }
func (f AssignmentOperator) Format(s fmt.State, v rune)       { format(s, v, f) }
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
func (f IdentifierReference) Format(s fmt.State, v rune)      { format(s, v, f) }
func (f BindingIdentifier) Format(s fmt.State, v rune)        { format(s, v, f) }
func (f LabelIdentifier) Format(s fmt.State, v rune)          { format(s, v, f) }
func (f Identifier) Format(s fmt.State, v rune)               { format(s, v, f) }
func (f Declaration) Format(s fmt.State, v rune)              { format(s, v, f) }
func (f LetOrConst) Format(s fmt.State, v rune)               { format(s, v, f) }
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
func (f ImportDeclaration) Format(s fmt.State, v rune)        { format(s, v, f) }
func (f ImportClause) Format(s fmt.State, v rune)             { format(s, v, f) }
func (f ImportedBinding) Format(s fmt.State, v rune)          { format(s, v, f) }
func (f FromClause) Format(s fmt.State, v rune)               { format(s, v, f) }
func (f NamedImports) Format(s fmt.State, v rune)             { format(s, v, f) }
func (f ImportSpecifier) Format(s fmt.State, v rune)          { format(s, v, f) }
func (f ExportDeclaration) Format(s fmt.State, v rune)        { format(s, v, f) }
func (f ExportClause) Format(s fmt.State, v rune)             { format(s, v, f) }
func (f ExportSpecifier) Format(s fmt.State, v rune)          { format(s, v, f) }
func (f jsParser) Format(s fmt.State, v rune)                 { format(s, v, f) }
func (f StatementList) Format(s fmt.State, v rune)            { format(s, v, f) }
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
