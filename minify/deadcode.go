package minify

import (
	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
	"vimagination.zapto.org/javascript/walk"
)

func removeDeadCode(m *javascript.Module) {
	for {
		s, err := scope.ModuleScope(m, nil)
		if err != nil {
			return
		}

		clearSinglesFromScope(s)

		walk.Walk(m, walk.HandlerFunc(deadWalker))
	}
}

func clearSinglesFromScope(s *scope.Scope) bool {
	changed := false
	for name, bindings := range s.Bindings {
		if name == "this" || name == "arguments" || len(bindings) != 1 {
			continue
		}
		bindings[0].Token.Data = ""
		changed = true
	}

	for _, cs := range s.Scopes {
		if clearSinglesFromScope(cs) {
			changed = true
		}
	}

	return changed
}

func deadWalker(t javascript.Type) error {
	switch t := t.(type) {
	case *javascript.Module:
		removeDeadCodeFromModule(t)
	case *javascript.Block:
		m := javascript.ScriptToModule(&javascript.Script{
			StatementList: t.StatementList,
		})
		removeDeadCodeFromModule(m)
		t.StatementList = make([]javascript.StatementListItem, len(m.ModuleListItems))
		for n, sli := range m.ModuleListItems {
			t.StatementList[n] = *sli.StatementListItem
		}
	default:
		deadWalker(t)
	}
	return nil
}

func removeDeadCodeFromModule(m *javascript.Module) {
	for i := 0; i < len(m.ModuleListItems); i++ {
		switch sliCLV(m.ModuleListItems[i].StatementListItem) {
		case clvConst, clvLet:
			bl := m.ModuleListItems[i].StatementListItem.Declaration.LexicalDeclaration.BindingList
			ls := make([]javascript.ModuleItem, len(bl), len(m.ModuleListItems)-i)
			for n := range ls {
				ls[n] = javascript.ModuleItem{
					StatementListItem: &javascript.StatementListItem{
						Declaration: &javascript.Declaration{
							LexicalDeclaration: &javascript.LexicalDeclaration{
								LetOrConst:  m.ModuleListItems[i].StatementListItem.Declaration.LexicalDeclaration.LetOrConst,
								BindingList: bl[n : n+1],
							},
						},
					},
				}
			}
			m.ModuleListItems = append(m.ModuleListItems[:i], append(ls, m.ModuleListItems[i+1:]...)...)
			i += len(bl)
		case clvVar:
			dl := m.ModuleListItems[i].StatementListItem.Statement.VariableStatement.VariableDeclarationList
			ls := make([]javascript.ModuleItem, len(dl), len(m.ModuleListItems)-i)
			for n := range ls {
				ls[n] = javascript.ModuleItem{
					StatementListItem: &javascript.StatementListItem{
						Statement: &javascript.Statement{
							VariableStatement: &javascript.VariableStatement{
								VariableDeclarationList: dl[n : n+1],
							},
						},
					},
				}
			}
			m.ModuleListItems = append(m.ModuleListItems[:i], append(ls, m.ModuleListItems[i+1:]...)...)
			i += len(dl)
		}
	}
	for i := 0; i < len(m.ModuleListItems); i++ {
		if removeDeadSLI(m.ModuleListItems[i].StatementListItem) {
			m.ModuleListItems = append(m.ModuleListItems[:i], m.ModuleListItems[i+1:]...)
			i--
		}
	}
	last := clvNone
	for i := 0; i < len(m.ModuleListItems); i++ {
		next := sliCLV(m.ModuleListItems[i].StatementListItem)
		if last == next {
			switch next {
			case clvConst, clvLet:
				ld := m.ModuleListItems[i-1].StatementListItem.Declaration.LexicalDeclaration
				ld.BindingList = append(ld.BindingList, m.ModuleListItems[i].StatementListItem.Declaration.LexicalDeclaration.BindingList[0])
			case clvVar:
				vd := m.ModuleListItems[i-1].StatementListItem.Statement.VariableStatement
				vd.VariableDeclarationList = append(vd.VariableDeclarationList, m.ModuleListItems[i].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0])
			}
			m.ModuleListItems = append(m.ModuleListItems[:i], m.ModuleListItems[i+1:]...)
			i--
		}
		last = next
	}
}

type clv byte

const (
	clvNone clv = iota
	clvConst
	clvLet
	clvVar
)

func sliCLV(sli *javascript.StatementListItem) clv {
	if sli != nil {
		if sli.Declaration != nil && sli.Declaration.LexicalDeclaration != nil {
			if sli.Declaration.LexicalDeclaration.LetOrConst == javascript.Const {
				return clvConst
			}
			return clvLet
		}
		if sli.Statement != nil && sli.Statement.VariableStatement != nil {
			return clvVar
		}
	}
	return clvNone
}

func removeDeadSLI(sli *javascript.StatementListItem) bool {
	switch sliCLV(sli) {
	case clvConst, clvLet:
		lb := sli.Declaration.LexicalDeclaration.BindingList[0]
		if lb.BindingIdentifier != nil {
			return lb.BindingIdentifier.Data == ""
		}
	case clvVar:
		vd := sli.Statement.VariableStatement.VariableDeclarationList[0]
		if vd.BindingIdentifier != nil {
			return vd.BindingIdentifier.Data == ""
		}
	}

	return false
}
