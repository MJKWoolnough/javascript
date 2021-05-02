package scope

import "vimagination.zapto.org/javascript"

type Binding struct {
	*Scope
	*Token
}

type Scope struct {
	IsBlockScope bool
	Parent       *Scope
	Scopes       []Scope
	Bindings     map[string]Binding
}

func ModuleScope(m *javascript.Module, global *Scope) *Scope {
	if global == nil {
		global = new(Scope)
	}
	for _, i := range m.ModuleListItems {
		if i.ImportDeclaration != nil {

		} else if i.StatementListItem != nil {
			processStatementListItem(i.StatementListItem, global)
		} else if i.ExportDeclaration != nil {

		}
	}
	return global
}

func ScriptScope(s *javascript.Script, global *Scope) *Scope {
	if global == nil {
		global = new(Scope)
	}
	for _, i := range s.StatementList {
		processStatementListItem(i, global)
	}
	return global
}

func processStatementListItem(s *javascript.StatementListItem, scope *Scope) {
	if s.Statement != nil {
		processStatement(s.Statement, scope)
	} else if s.Declaration != nil {
		processDeclaration(s.Declaration, scope)
	}
}

func processStatement(s *javascript.Statement, scope *Scope) {

}

func processDeclaration(d *javascript.Declaration, scope *Scope) {

}
