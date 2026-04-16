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
	} {
		tk := parser.NewStringTokeniser(test.Input)

		m, err := javascript.ParseModule(&tk)
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
