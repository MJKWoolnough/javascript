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
		for i := 0; i < len(t.ModuleListItems); i++ {
			switch sliCLV(t.ModuleListItems[i].StatementListItem) {
			case clvLexical:
				bl := t.ModuleListItems[i].StatementListItem.Declaration.LexicalDeclaration.BindingList
				ls := make([]javascript.ModuleItem, len(bl), len(t.ModuleListItems)-i)
				for n := range ls {
					ls[n] = javascript.ModuleItem{
						StatementListItem: &javascript.StatementListItem{
							Declaration: &javascript.Declaration{
								LexicalDeclaration: &javascript.LexicalDeclaration{
									LetOrConst:  t.ModuleListItems[i].StatementListItem.Declaration.LexicalDeclaration.LetOrConst,
									BindingList: bl[n : n+1],
								},
							},
						},
					}
				}
				t.ModuleListItems = append(t.ModuleListItems[:i], append(ls, t.ModuleListItems[i+1:]...)...)
				i += len(bl)
			}
		}
	}
	return nil
}

type clv byte

const (
	clvNone clv = iota
	clvLexical
	clvVar
)

func sliCLV(sli *javascript.StatementListItem) clv {
	if sli != nil {
		if sli.Declaration != nil && sli.Declaration.LexicalDeclaration != nil {
			return clvLexical
		}
		if sli.Statement != nil && sli.Statement.VariableStatement != nil {
			return clvVar
		}
	}
	return clvNone
}
