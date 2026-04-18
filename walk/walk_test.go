package walk

import (
	"errors"
	"reflect"
	"testing"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/parser"
)

var (
	sentinel = errors.New("")
	nilErr   = errors.New("nil received")
	nilRet   = func(_ *javascript.Module) javascript.Type { return nil }
)

type walker struct {
	end   javascript.Type
	level []string
}

func (w *walker) Handle(t javascript.Type) error {
	if reflect.ValueOf(t).IsNil() {
		return nilErr
	}

	if t == w.end {
		w.level = append(w.level, reflect.TypeOf(t).Elem().Name())

		return sentinel
	}

	err := Walk(t, w)
	if err != nil {
		w.level = append(w.level, reflect.TypeOf(t).Elem().Name())
	}

	return err
}

func TestWalk(t *testing.T) {
	for n, test := range [...]struct {
		Input string
		End   func(m *javascript.Module) javascript.Type
		Level []string
	}{
		{ // 1
			"",
			nilRet,
			nil,
		},
		{ // 2
			"import a from './b';",
			func(m *javascript.Module) javascript.Type { return &m.ModuleListItems[0] },
			[]string{"Module", "ModuleItem"},
		},
		{ // 3
			"import a from './b';",
			func(m *javascript.Module) javascript.Type { return m.ModuleListItems[0].ImportDeclaration },
			[]string{"Module", "ModuleItem", "ImportDeclaration"},
		},
		{ // 4
			"import a from './b';",
			nilRet,
			nil,
		},
		{ // 5
			"import a from './b';",
			func(m *javascript.Module) javascript.Type { return m.ModuleListItems[0].ImportDeclaration.ImportClause },
			[]string{"Module", "ModuleItem", "ImportDeclaration", "ImportClause"},
		},
		{ // 6
			"import a from './b';",
			func(m *javascript.Module) javascript.Type { return &m.ModuleListItems[0].ImportDeclaration.FromClause },
			[]string{"Module", "ModuleItem", "ImportDeclaration", "FromClause"},
		},
		{ // 7
			"import a from './b' with {};",
			func(m *javascript.Module) javascript.Type { return m.ModuleListItems[0].ImportDeclaration.WithClause },
			[]string{"Module", "ModuleItem", "ImportDeclaration", "WithClause"},
		},
		{ // 8
			"import {} from './b';",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].ImportDeclaration.ImportClause.NamedImports
			},
			[]string{"Module", "ModuleItem", "ImportDeclaration", "ImportClause", "NamedImports"},
		},
		{ // 9
			"import {} from './b';",
			nilRet,
			nil,
		},
		{ // 10
			"import {a} from './b';",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].ImportDeclaration.ImportClause.NamedImports.ImportList[0]
			},
			[]string{"Module", "ModuleItem", "ImportDeclaration", "ImportClause", "NamedImports", "ImportSpecifier"},
		},
		{ // 11
			"import {a, b} from './c';",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].ImportDeclaration.ImportClause.NamedImports.ImportList[1]
			},
			[]string{"Module", "ModuleItem", "ImportDeclaration", "ImportClause", "NamedImports", "ImportSpecifier"},
		},
		{ // 12
			"export {};",
			nilRet,
			nil,
		},
		{ // 13
			"export default a;",
			func(m *javascript.Module) javascript.Type { return m.ModuleListItems[0].ExportDeclaration },
			[]string{"Module", "ModuleItem", "ExportDeclaration"},
		},
		{ // 14
			"export {a};",
			func(m *javascript.Module) javascript.Type { return m.ModuleListItems[0].ExportDeclaration.ExportClause },
			[]string{"Module", "ModuleItem", "ExportDeclaration", "ExportClause"},
		},
		{ // 15
			"export {a} from './b';",
			func(m *javascript.Module) javascript.Type { return m.ModuleListItems[0].ExportDeclaration.FromClause },
			[]string{"Module", "ModuleItem", "ExportDeclaration", "FromClause"},
		},
		{ // 16
			"export var a;",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].ExportDeclaration.VariableStatement
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "VariableStatement"},
		},
		{ // 17
			"export let a;",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].ExportDeclaration.Declaration
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "Declaration"},
		},
		{ // 18
			"export default function a(){}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].ExportDeclaration.DefaultFunction
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "FunctionDeclaration"},
		},
		{ // 19
			"export default class a{}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].ExportDeclaration.DefaultClass
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "ClassDeclaration"},
		},
		{ // 20
			"export default () => {}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].ExportDeclaration.DefaultAssignmentExpression
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "AssignmentExpression"},
		},
		{ // 21
			"export {a, b}",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].ExportDeclaration.ExportClause.ExportList[0]
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "ExportClause", "ExportSpecifier"},
		},
		{ // 22
			"export {a, b}",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].ExportDeclaration.ExportClause.ExportList[1]
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "ExportClause", "ExportSpecifier"},
		},
		{ // 23
			"a",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem
			},
			[]string{"Module", "ModuleItem", "StatementListItem"},
		},
		{ // 24
			"a",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement"},
		},
		{ // 25
			"let a",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Declaration
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Declaration"},
		},
		{ // 26
			"{}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.BlockStatement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Block"},
		},
		{ // 27
			"var a",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.VariableStatement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement"},
		},
		{ // 28
			"a",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression"},
		},
		{ // 29
			"if (a) {}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.IfStatement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "IfStatement"},
		},
		{ // 30
			"do; while (a)",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.IterationStatementDo
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "IterationStatementDo"},
		},
		{ // 31
			"while (a){}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.IterationStatementWhile
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "IterationStatementWhile"},
		},
		{ // 32
			"for (;;) {}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "IterationStatementFor"},
		},
		{ // 33
			"switch (a){}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.SwitchStatement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "SwitchStatement"},
		},
		{ // 34
			"with (a){}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.WithStatement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "WithStatement"},
		},
		{ // 35
			"a: function b (){}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.LabelledItemFunction
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "FunctionDeclaration"},
		},
		{ // 36
			"a: b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.LabelledItemStatement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Statement"},
		},
		{ // 37
			"try {} finally{}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.TryStatement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "TryStatement"},
		},
		{ // 38
			"{}",
			nilRet,
			nil,
		},
		{ // 39
			"{a}",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.BlockStatement.StatementList[0]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Block", "StatementListItem"},
		},
		{ // 40
			"{a; b}",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.BlockStatement.StatementList[1]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Block", "StatementListItem"},
		},
		{ // 41
			"var a",
			nilRet,
			nil,
		},
		{ // 42
			"var a",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding"},
		},
		{ // 43
			"var a, b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[1]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding"},
		},
		{ // 44
			"var [a] = b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ArrayBindingPattern
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "ArrayBindingPattern"},
		},
		{ // 45
			"var {a} = b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ObjectBindingPattern
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "ObjectBindingPattern"},
		},
		{ // 46
			"var a = b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].Initializer
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "AssignmentExpression"},
		},
		{ // 47
			"var [a] = []",
			nilRet,
			nil,
		},
		{ // 48
			"var [a, b, ...c] = []",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ArrayBindingPattern.BindingElementList[0]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "ArrayBindingPattern", "BindingElement"},
		},
		{ // 49
			"var [a, b, ...c] = []",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ArrayBindingPattern.BindingElementList[1]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "ArrayBindingPattern", "BindingElement"},
		},
		{ // 50
			"var [a, b, ...c] = []",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ArrayBindingPattern.BindingRestElement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "ArrayBindingPattern", "BindingElement"},
		},
		{ // 51
			"var [[]] = []",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ArrayBindingPattern.BindingElementList[0].ArrayBindingPattern
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "ArrayBindingPattern", "BindingElement", "ArrayBindingPattern"},
		},
		{ // 52
			"var [{}] = []",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ArrayBindingPattern.BindingElementList[0].ObjectBindingPattern
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "ArrayBindingPattern", "BindingElement", "ObjectBindingPattern"},
		},
		{ // 53
			"var [a = b] = []",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ArrayBindingPattern.BindingElementList[0].Initializer
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "ArrayBindingPattern", "BindingElement", "AssignmentExpression"},
		},
		{ // 54
			"var {a} = {}",
			nilRet,
			nil,
		},
		{ // 55
			"var {a, b, ...c} = []",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ObjectBindingPattern.BindingPropertyList[0]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "ObjectBindingPattern", "BindingProperty"},
		},
		{ // 56
			"var {a, b, ...c} = []",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ObjectBindingPattern.BindingPropertyList[1]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "ObjectBindingPattern", "BindingProperty"},
		},
		{ // 57
			"var {a: b} = []",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ObjectBindingPattern.BindingPropertyList[0].PropertyName
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "ObjectBindingPattern", "BindingProperty", "PropertyName"},
		},
		{ // 58
			"var {a: b} = []",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ObjectBindingPattern.BindingPropertyList[0].BindingElement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "ObjectBindingPattern", "BindingProperty", "BindingElement"},
		},
		{ // 59
			"var {[a]: b} = []",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ObjectBindingPattern.BindingPropertyList[0].PropertyName.ComputedPropertyName
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "VariableStatement", "LexicalBinding", "ObjectBindingPattern", "BindingProperty", "PropertyName", "AssignmentExpression"},
		},
		{ // 60
			"a, b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression"},
		},
		{ // 61
			"a, b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression"},
		},
		{ // 62
			"a = b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].LeftHandSideExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "LeftHandSideExpression"},
		},
		{ // 63
			"a = b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "AssignmentExpression"},
		},
		{ // 64
			"[a] = b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "AssignmentPattern"},
		},
		{ // 65
			"[a] = b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "AssignmentExpression"},
		},
		{ // 66
			"a",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression"},
		},
		{ // 67
			"() => {}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ArrowFunction
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ArrowFunction"},
		},
		{ // 68
			"[a] = b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "AssignmentPattern", "ArrayAssignmentPattern"},
		},
		{ // 69
			"({a} = b)",
			func(m *javascript.Module) javascript.Type {
				return javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "AssignmentPattern", "ObjectAssignmentPattern"},
		},
		{ // 70
			"[a, b, ...c] = d",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "AssignmentPattern", "ArrayAssignmentPattern", "AssignmentElement"},
		},
		{ // 71
			"[a, b, ...c] = d",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[1]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "AssignmentPattern", "ArrayAssignmentPattern", "AssignmentElement"},
		},
		{ // 72
			"[a, b, ...c] = d",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentRestElement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "AssignmentPattern", "ArrayAssignmentPattern", "LeftHandSideExpression"},
		},
		{ // 73
			"[a = b] = c",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].DestructuringAssignmentTarget
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "AssignmentPattern", "ArrayAssignmentPattern", "AssignmentElement", "DestructuringAssignmentTarget"},
		},
		{ // 74
			"[a = b] = c",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].Initializer
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "AssignmentPattern", "ArrayAssignmentPattern", "AssignmentElement", "AssignmentExpression"},
		},
		{ // 75
			"[a = b] = c",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].DestructuringAssignmentTarget.LeftHandSideExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "AssignmentPattern", "ArrayAssignmentPattern", "AssignmentElement", "DestructuringAssignmentTarget", "LeftHandSideExpression"},
		},
		{ // 76
			"[[a] = b] = c",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].DestructuringAssignmentTarget.AssignmentPattern
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "AssignmentPattern", "ArrayAssignmentPattern", "AssignmentElement", "DestructuringAssignmentTarget", "AssignmentPattern"},
		},
		{ // 77
			"({a, b, ...c} = d)",
			func(m *javascript.Module) javascript.Type {
				return &javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[0]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "AssignmentPattern", "ObjectAssignmentPattern", "AssignmentProperty"},
		},
		{ // 78
			"({a, b, ...c} = d)",
			func(m *javascript.Module) javascript.Type {
				return &javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[1]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "AssignmentPattern", "ObjectAssignmentPattern", "AssignmentProperty"},
		},
		{ // 79
			"({a, b, ...c} = d)",
			func(m *javascript.Module) javascript.Type {
				return javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentRestElement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "AssignmentPattern", "ObjectAssignmentPattern", "LeftHandSideExpression"},
		},
		{ // 80
			"({a} = b)",
			func(m *javascript.Module) javascript.Type {
				return &javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[0].PropertyName
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "AssignmentPattern", "ObjectAssignmentPattern", "AssignmentProperty", "PropertyName"},
		},
		{ // 81
			"({a: {b} = c} = d)",
			func(m *javascript.Module) javascript.Type {
				return javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[0].DestructuringAssignmentTarget
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "AssignmentPattern", "ObjectAssignmentPattern", "AssignmentProperty", "DestructuringAssignmentTarget"},
		},
		{ // 82
			"({a: {b} = c} = d)",
			func(m *javascript.Module) javascript.Type {
				return javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[0].Initializer
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "AssignmentPattern", "ObjectAssignmentPattern", "AssignmentProperty", "AssignmentExpression"},
		},
		{ // 83
			"a",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression"},
		},
		{ // 84
			"a",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression"},
		},
		{ // 85
			"a ?? b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.CoalesceExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "CoalesceExpression"},
		},
		{ // 86
			"a ? b : c",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression"},
		},
		{ // 87
			"a ? b : c",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.True
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "AssignmentExpression"},
		},
		{ // 88
			"a ? b : c",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.False
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "AssignmentExpression"},
		},
		{ // 89
			"a ?? b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.CoalesceExpression.BitwiseORExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "CoalesceExpression", "BitwiseORExpression"},
		},
		{ // 90
			"a ?? b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.CoalesceExpression.CoalesceExpressionHead
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "CoalesceExpression", "CoalesceExpression"},
		},
		{ // 91
			"a || b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalORExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalORExpression"},
		},
		{ // 92
			"a || b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression"},
		},
		{ // 93
			"a && b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.LogicalANDExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "LogicalANDExpression"},
		},
		{ // 94
			"a && b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression"},
		},
		{ // 95
			"a | b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseORExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseORExpression"},
		},
		{ // 96
			"a | b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression"},
		},
		{ // 97
			"a ^ b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseXORExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseXORExpression"},
		},
		{ // 98
			"a ^ b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression"},
		},
		{ // 99
			"a & b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.BitwiseANDExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "BitwiseANDExpression"},
		},
		{ // 100
			"a & b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression"},
		},
		{ // 101
			"a == b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.EqualityExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "EqualityExpression"},
		},
		{ // 102
			"a == b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression"},
		},
		{ // 103
			"a <= b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.RelationalExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "RelationalExpression"},
		},
		{ // 104
			"a <= b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression"},
		},
		{ // 105
			"a << b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.ShiftExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "ShiftExpression"},
		},
		{ // 106
			"a << b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression"},
		},
		{ // 107
			"a + b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.AdditiveExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "AdditiveExpression"},
		},
		{ // 108
			"a + b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression"},
		},
		{ // 109
			"a * b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.MultiplicativeExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "MultiplicativeExpression"},
		},
		{ // 110
			"a * b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression"},
		},
		{ // 111
			"a ** b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.ExponentiationExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "ExponentiationExpression"},
		},
		{ // 112
			"a ** b",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression"},
		},
		{ // 113
			"+a",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression"},
		},
		{ // 114
			"++a",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.UnaryExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "UnaryExpression"},
		},
		{ // 115
			"a++",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression"},
		},
		{ // 116
			"a",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression"},
		},
		{ // 117
			"a()",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.CallExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "CallExpression"},
		},
		{ // 118
			"a?.b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.OptionalExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "OptionalExpression"},
		},
		{ // 119
			"a",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression"},
		},
		{ // 120
			"a.b",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.MemberExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "MemberExpression"},
		},
		{ // 121
			"a",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression"},
		},
		{ // 122
			"a[b]",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.Expression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "Expression"},
		},
		{ // 123
			"a``",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.TemplateLiteral
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "TemplateLiteral"},
		},
		{ // 124
			"new a()",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.Arguments
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "Arguments"},
		},
		{ // 125
			"[]",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ArrayLiteral
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ArrayLiteral"},
		},
		{ // 126
			"(a)",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression"},
		},
		{ // 127
			"({})",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ObjectLiteral
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ObjectLiteral"},
		},
		{ // 128
			"(function(){})",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.FunctionExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "FunctionDeclaration"},
		},
		{ // 129
			"(class{})",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ClassExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ClassDeclaration"},
		},
		{ // 130
			"``",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.TemplateLiteral
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "TemplateLiteral"},
		},
		{ // 131
			"<a/>",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.JSXElement
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "JSXElement"},
		},
		{ // 132
			"<></>",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.JSXFragment
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "JSXFragment"},
		},
		{ // 133
			"[a, b]",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ArrayLiteral.ElementList[0]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ArrayLiteral", "ArrayElement"},
		},
		{ // 134
			"[a, b]",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ArrayLiteral.ElementList[1]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ArrayLiteral", "ArrayElement"},
		},
		{ // 135
			"[a]",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ArrayLiteral.ElementList[0].AssignmentExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ArrayLiteral", "ArrayElement", "AssignmentExpression"},
		},
		{ // 136
			"({a: b, c: d})",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ObjectLiteral.PropertyDefinitionList[0]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ObjectLiteral", "PropertyDefinition"},
		},
		{ // 137
			"({a: b, c: d})",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ObjectLiteral.PropertyDefinitionList[1]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ObjectLiteral", "PropertyDefinition"},
		},
		{ // 138
			"({a: b})",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ObjectLiteral.PropertyDefinitionList[0].PropertyName
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ObjectLiteral", "PropertyDefinition", "PropertyName"},
		},
		{ // 139
			"({a: b})",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ObjectLiteral.PropertyDefinitionList[0].AssignmentExpression
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ObjectLiteral", "PropertyDefinition", "AssignmentExpression"},
		},
		{ // 140
			"({a(){}})",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ObjectLiteral.PropertyDefinitionList[0].MethodDefinition
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ParenthesizedExpression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "ObjectLiteral", "PropertyDefinition", "MethodDefinition"},
		},
		{ // 141
			"`a${b}c${d}e`",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.TemplateLiteral.Expressions[0]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "TemplateLiteral", "Expression"},
		},
		{ // 142
			"`a${b}c${d}e`",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.TemplateLiteral.Expressions[1]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "TemplateLiteral", "Expression"},
		},
		{ // 143
			"<a/>",
			nilRet,
			nil,
		},
		{ // 144
			"<a b='c' d='e'>d<></></a>",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.JSXElement.ElementName
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "JSXElement", "JSXElementName"},
		},
		{ // 145
			"<a b='c' d='e'>d<></></a>",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.JSXElement.Attributes[0]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "JSXElement", "JSXAttribute"},
		},
		{ // 146
			"<a b='c' d='e'>d<></></a>",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.JSXElement.Attributes[1]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "JSXElement", "JSXAttribute"},
		},
		{ // 147
			"<a b='c' d='e'>d<></></a>",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.JSXElement.Children[0]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "JSXElement", "JSXChild"},
		},
		{ // 148
			"<a b='c' d='e'>d<></></a>",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.JSXElement.Children[1]
			},
			[]string{"Module", "ModuleItem", "StatementListItem", "Statement", "Expression", "AssignmentExpression", "ConditionalExpression", "LogicalORExpression", "LogicalANDExpression", "BitwiseORExpression", "BitwiseXORExpression", "BitwiseANDExpression", "EqualityExpression", "RelationalExpression", "ShiftExpression", "AdditiveExpression", "MultiplicativeExpression", "ExponentiationExpression", "UnaryExpression", "UpdateExpression", "LeftHandSideExpression", "NewExpression", "MemberExpression", "PrimaryExpression", "JSXElement", "JSXChild"},
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)

		m, err := javascript.ParseModule(javascript.AsJSX(&tk))
		if err != nil {
			t.Errorf("test %d: unexpected error parsing script: %s", n+1, err)
		} else {
			w := walker{end: test.End(m)}

			if err := w.Handle(m); err == nil && test.Level != nil {
				t.Errorf("test %d: expected to recieve sentinel error, but didn't", n+1)
			} else if err != nil && test.Level == nil {
				t.Errorf("test %d: expected no error, but recieved %v", n+1, err)
			} else if len(w.level) != len(test.Level) {
				t.Errorf("test %d: expected to have %d levels, got %d", n+1, len(test.Level), len(w.level))
			} else {
				for m, l := range w.level {
					if e := test.Level[len(test.Level)-m-1]; e != l {
						t.Errorf("test %d.%d: expected to read level %s, got %s", n+1, m+1, e, l)
					}
				}
			}
		}
	}
}
