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
			case clvConst, clvLet:
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
			case clvVar:
				dl := t.ModuleListItems[i].StatementListItem.Statement.VariableStatement.VariableDeclarationList
				ls := make([]javascript.ModuleItem, len(dl), len(t.ModuleListItems)-i)
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
				t.ModuleListItems = append(t.ModuleListItems[:i], append(ls, t.ModuleListItems[i+1:]...)...)
				i += len(dl)
			}
		}
		for i := 0; i < len(t.ModuleListItems); i++ {
			if removeDeadSLI(t.ModuleListItems[i].StatementListItem) {
				t.ModuleListItems = append(t.ModuleListItems[:i], t.ModuleListItems[i+1:]...)
				i--
			}
		}
		last := clvNone
		for i := 0; i < len(t.ModuleListItems); i++ {
			next := sliCLV(t.ModuleListItems[i].StatementListItem)
			if last == next {
				switch next {
				case clvConst, clvLet:
					ld := t.ModuleListItems[i-1].StatementListItem.Declaration.LexicalDeclaration
					ld.BindingList = append(ld.BindingList, t.ModuleListItems[i].StatementListItem.Declaration.LexicalDeclaration.BindingList[0])
				case clvVar:
					vd := t.ModuleListItems[i-1].StatementListItem.Statement.VariableStatement
					vd.VariableDeclarationList = append(vd.VariableDeclarationList, t.ModuleListItems[i].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0])
				}
				t.ModuleListItems = append(t.ModuleListItems[:i], t.ModuleListItems[i+1:]...)
				i--
			}
			last = next
		}
	case *javascript.Block:
		for i := 0; i < len(t.StatementList); i++ {
			switch sliCLV(&t.StatementList[i]) {
			case clvConst, clvLet:
				bl := t.StatementList[i].Declaration.LexicalDeclaration.BindingList
				ls := make([]javascript.StatementListItem, len(bl), len(t.StatementList)-i)
				for n := range ls {
					ls[n] = javascript.StatementListItem{
						Declaration: &javascript.Declaration{
							LexicalDeclaration: &javascript.LexicalDeclaration{
								LetOrConst:  t.StatementList[i].Declaration.LexicalDeclaration.LetOrConst,
								BindingList: bl[n : n+1],
							},
						},
					}
				}
				t.StatementList = append(t.StatementList[:i], append(ls, t.StatementList[i+1:]...)...)
				i += len(bl)
			case clvVar:
				dl := t.StatementList[i].Statement.VariableStatement.VariableDeclarationList
				ls := make([]javascript.StatementListItem, len(dl), len(t.StatementList)-i)
				for n := range ls {
					ls[n] = javascript.StatementListItem{
						Statement: &javascript.Statement{
							VariableStatement: &javascript.VariableStatement{
								VariableDeclarationList: dl[n : n+1],
							},
						},
					}
				}
				t.StatementList = append(t.StatementList[:i], append(ls, t.StatementList[i+1:]...)...)
				i += len(dl)
			}
		}
		for i := 0; i < len(t.StatementList); i++ {
			if removeDeadSLI(&t.StatementList[i]) {
				t.StatementList = append(t.StatementList[:i], t.StatementList[i+1:]...)
				i--
			}
		}
		last := clvNone
		for i := 0; i < len(t.StatementList); i++ {
			next := sliCLV(&t.StatementList[i])
			if last == next {
				switch next {
				case clvConst, clvLet:
					ld := t.StatementList[i-1].Declaration.LexicalDeclaration
					ld.BindingList = append(ld.BindingList, t.StatementList[i].Declaration.LexicalDeclaration.BindingList[0])
				case clvVar:
					vd := t.StatementList[i-1].Statement.VariableStatement
					vd.VariableDeclarationList = append(vd.VariableDeclarationList, t.StatementList[i].Statement.VariableStatement.VariableDeclarationList[0])
				}
				t.StatementList = append(t.StatementList[:i], t.StatementList[i+1:]...)
				i--
			}
			last = next
		}
	}
	return nil
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
	if sli == nil {
		return false
	}

	return false
}
