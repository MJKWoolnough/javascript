package walk

import (
	"errors"
	"reflect"
	"testing"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/parser"
)

var sentinel = errors.New("")

type walker struct {
	end   javascript.Type
	level []string
}

func (w *walker) Handle(t javascript.Type) error {
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
		{
			"",
			func(m *javascript.Module) javascript.Type { return nil },
			nil,
		},
		{
			"import a from './b';",
			func(m *javascript.Module) javascript.Type { return &m.ModuleListItems[0] },
			[]string{"Module", "ModuleItem"},
		},
		{
			"import a from './b';",
			func(m *javascript.Module) javascript.Type { return m.ModuleListItems[0].ImportDeclaration },
			[]string{"Module", "ModuleItem", "ImportDeclaration"},
		},
		{
			"import a from './b';",
			func(m *javascript.Module) javascript.Type { return nil },
			nil,
		},
		{
			"import a from './b';",
			func(m *javascript.Module) javascript.Type { return m.ModuleListItems[0].ImportDeclaration.ImportClause },
			[]string{"Module", "ModuleItem", "ImportDeclaration", "ImportClause"},
		},
		{
			"import a from './b';",
			func(m *javascript.Module) javascript.Type { return &m.ModuleListItems[0].ImportDeclaration.FromClause },
			[]string{"Module", "ModuleItem", "ImportDeclaration", "FromClause"},
		},
		{
			"import a from './b' with {};",
			func(m *javascript.Module) javascript.Type { return m.ModuleListItems[0].ImportDeclaration.WithClause },
			[]string{"Module", "ModuleItem", "ImportDeclaration", "WithClause"},
		},
		{
			"import {} from './b';",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].ImportDeclaration.ImportClause.NamedImports
			},
			[]string{"Module", "ModuleItem", "ImportDeclaration", "ImportClause", "NamedImports"},
		},
		{
			"import {} from './b';",
			func(m *javascript.Module) javascript.Type { return nil },
			nil,
		},
		{
			"import {a} from './b';",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].ImportDeclaration.ImportClause.NamedImports.ImportList[0]
			},
			[]string{"Module", "ModuleItem", "ImportDeclaration", "ImportClause", "NamedImports", "ImportSpecifier"},
		},
		{
			"import {a, b} from './c';",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].ImportDeclaration.ImportClause.NamedImports.ImportList[1]
			},
			[]string{"Module", "ModuleItem", "ImportDeclaration", "ImportClause", "NamedImports", "ImportSpecifier"},
		},
		{
			"export {};",
			func(m *javascript.Module) javascript.Type { return nil },
			nil,
		},
		{
			"export default a;",
			func(m *javascript.Module) javascript.Type { return m.ModuleListItems[0].ExportDeclaration },
			[]string{"Module", "ModuleItem", "ExportDeclaration"},
		},
		{
			"export {a};",
			func(m *javascript.Module) javascript.Type { return m.ModuleListItems[0].ExportDeclaration.ExportClause },
			[]string{"Module", "ModuleItem", "ExportDeclaration", "ExportClause"},
		},
		{
			"export {a} from './b';",
			func(m *javascript.Module) javascript.Type { return m.ModuleListItems[0].ExportDeclaration.FromClause },
			[]string{"Module", "ModuleItem", "ExportDeclaration", "FromClause"},
		},
		{
			"export var a;",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].ExportDeclaration.VariableStatement
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "VariableStatement"},
		},
		{
			"export let a;",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].ExportDeclaration.Declaration
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "Declaration"},
		},
		{
			"export default function a(){}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].ExportDeclaration.DefaultFunction
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "FunctionDeclaration"},
		},
		{
			"export default class a{}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].ExportDeclaration.DefaultClass
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "ClassDeclaration"},
		},
		{
			"export default () => {}",
			func(m *javascript.Module) javascript.Type {
				return m.ModuleListItems[0].ExportDeclaration.DefaultAssignmentExpression
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "AssignmentExpression"},
		},
		{
			"export {a, b}",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].ExportDeclaration.ExportClause.ExportList[0]
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "ExportClause", "ExportSpecifier"},
		},
		{
			"export {a, b}",
			func(m *javascript.Module) javascript.Type {
				return &m.ModuleListItems[0].ExportDeclaration.ExportClause.ExportList[1]
			},
			[]string{"Module", "ModuleItem", "ExportDeclaration", "ExportClause", "ExportSpecifier"},
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
