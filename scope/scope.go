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

		} else if i.ExportDeclaration != nil {

		}
	}
	return global
}

func ScriptScope(s *javascript.Script, global *Scope) *Scope {
	if global == nil {
		global = new(Scope)
	}

	return global
}
