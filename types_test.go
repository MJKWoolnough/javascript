package javascript

import (
	"fmt"
	"reflect"
	"testing"
)

const disabled = true

var allTypes = []interface{}{
	AdditiveExpression{},
	Argument{},
	Arguments{},
	ArrayAssignmentPattern{},
	ArrayBindingPattern{},
	ArrayElement{},
	ArrayLiteral{},
	ArrowFunction{},
	AssignmentElement{},
	AssignmentExpression{},
	AssignmentPattern{},
	AssignmentProperty{},
	BindingElement{},
	BindingProperty{},
	BitwiseANDExpression{},
	BitwiseORExpression{},
	BitwiseXORExpression{},
	Block{},
	CallExpression{},
	CaseClause{},
	ClassDeclaration{},
	ClassElement{},
	ClassElementName{},
	CoalesceExpression{},
	ConditionalExpression{},
	ParenthesizedExpression{},
	Declaration{},
	DestructuringAssignmentTarget{},
	EqualityExpression{},
	ExponentiationExpression{},
	ExportClause{},
	ExportDeclaration{},
	ExportSpecifier{},
	Expression{},
	FieldDefinition{},
	FormalParameters{},
	FromClause{},
	FunctionDeclaration{},
	IfStatement{},
	ImportClause{},
	ImportDeclaration{},
	ImportSpecifier{},
	IterationStatementDo{},
	IterationStatementFor{},
	IterationStatementWhile{},
	LeftHandSideExpression{},
	LexicalBinding{},
	LexicalDeclaration{},
	LogicalANDExpression{},
	LogicalORExpression{},
	MemberExpression{},
	MethodDefinition{},
	Module{},
	ModuleItem{},
	MultiplicativeExpression{},
	NamedImports{},
	NewExpression{},
	ObjectAssignmentPattern{},
	ObjectBindingPattern{},
	ObjectLiteral{},
	OptionalChain{},
	OptionalExpression{},
	PrimaryExpression{},
	PropertyDefinition{},
	PropertyName{},
	RelationalExpression{},
	Script{},
	ShiftExpression{},
	Statement{},
	StatementListItem{},
	SwitchStatement{},
	TemplateLiteral{},
	TryStatement{},
	UnaryExpression{},
	UpdateExpression{},
	VariableStatement{},
	WithStatement{},
}

func TestTypeFormatting(*testing.T) {
	if disabled {
		return
	}
	for _, typ := range allTypes {
		t := reflect.TypeOf(typ)
		fmt.Printf("\n\n// Format implements the fmt.Formatter interface\nfunc (f %[1]s) Format(s fmt.State, v rune) {\n	if v == 'v' && s.Flag('#') {\n		type X = %[1]s\n		type %[1]s X\n		fmt.Fprintf(s, \"%%#v\", %[1]s(f))\n	} else {\n		format(&f, s, v)\n	}\n}", t.Name())
	}
}

func TestTypePrinting(*testing.T) {
	if disabled {
		return
	}
	stringer := reflect.TypeOf(struct{ S fmt.Stringer }{}).Field(0).Type
	fmt.Printf("package javascript\n\nimport \"io\"\n\nvar (")
	done := make(map[string]struct{})
	for _, typ := range allTypes {
		t := reflect.TypeOf(typ)
		name := t.Name()
		if _, ok := done[name]; !ok {
			fmt.Printf("\n	name%s = %v", name, toSlice(name))
			done[name] = struct{}{}
		}
		for i := 0; i < t.NumField(); i++ {
			if k := t.Field(i).Type.Kind(); k == reflect.Bool || k == reflect.Uint {
				continue
			}
			name := t.Field(i).Name
			if _, ok := done[name]; !ok {
				fmt.Printf("\n	name%s = %v", name, toSlice(name))
				done[name] = struct{}{}
			}
		}
	}
	fmt.Printf("\n)\n\n")
	for _, typ := range allTypes {
		t := reflect.TypeOf(typ)
		fmt.Printf("func (f *%s) printType(w io.Writer, v bool) {\n	w.Write(name%s[1:%d])\n	w.Write(objectOpen)\n	pp := indentPrinter{w}\n", t.Name(), t.Name(), len(t.Name())+1)
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Name == "Tokens" {
				fmt.Printf("	if v {\n		pp.Write(tokensTo)\n		f.Tokens.printType(&pp, v)\n	}\n")
			} else if f.Type.Kind() == reflect.Ptr {
				fmt.Printf("	if f.%s != nil {\n		pp.Write(name%s)\n		f.%s.printType(&pp, v)\n	} else if v {\n		pp.Write(name%s)\n		pp.Write(nilStr)\n	}\n", f.Name, f.Name, f.Name, f.Name)
			} else if f.Type.Kind() == reflect.Struct {
				fmt.Printf("	pp.Write(name%s)\n	f.%s.printType(&pp, v)\n", f.Name, f.Name)
			} else if f.Type.Implements(stringer) {
				fmt.Printf("	pp.Write(name%s)\n	io.WriteString(&pp, f.%s.String())\n", f.Name, f.Name)
			} else if f.Type.Kind() == reflect.Slice {
				fmt.Printf("	if f.%s == nil {\n		pp.Write(name%s)\n		pp.Write(nilStr)\n	} else if len(f.%s) > 0 {\n		pp.Write(name%s)\n		pp.Write(arrayOpen)\n		ipp := indentPrinter{&pp}\n		for n, e := range f.%s {\n			ipp.Printf(\"\\n%%d:\", n)\n			e.printType(&ipp, v)\n		}\n		pp.Write(arrayClose)\n	} else if v {\n		pp.Write(name%s)\n              pp.Write(arrayOpenClose)\n	}\n", f.Name, f.Name, f.Name, f.Name, f.Name, f.Name)
			} else if f.Type.Kind() == reflect.Bool {
				fmt.Printf("	if f.%s || v {\n		pp.Printf(\"\\n%s: %%v\", f.%s)\n	}\n", f.Name, f.Name, f.Name)
			} else if f.Type.Kind() == reflect.Uint {
				fmt.Printf("	pp.Printf(\"\\n%s: %%d\", f.%s)\n", f.Name, f.Name)
			}
		}
		fmt.Printf("	w.Write(objectClose)\n}\n\n")
	}
}

func toSlice(str string) string {
	s := "[]byte{'\\n'"
	for _, b := range []byte(str) {
		s += ", '" + string(rune(b)) + "'"
	}
	return s + ", ':', ' '}"
}
