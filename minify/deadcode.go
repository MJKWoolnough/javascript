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
		blockAsModule(t, removeDeadCodeFromModule)
	case *javascript.Expression:
		expressionsAsModule(t, removeDeadCodeFromModule)
	default:
		deadWalker(t)
	}
	return nil
}

func removeDeadCodeFromModule(m *javascript.Module) {
	for i := 0; i < len(m.ModuleListItems); i++ {
		switch sliBindable(m.ModuleListItems[i].StatementListItem) {
		case bindableConst, bindableLet:
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
		case bindableVar:
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
	last := bindableNone
	for i := 0; i < len(m.ModuleListItems); i++ {
		next := sliBindable(m.ModuleListItems[i].StatementListItem)
		if last == next {
			switch next {
			case bindableConst, bindableLet:
				ld := m.ModuleListItems[i-1].StatementListItem.Declaration.LexicalDeclaration
				ld.BindingList = append(ld.BindingList, m.ModuleListItems[i].StatementListItem.Declaration.LexicalDeclaration.BindingList[0])
			case bindableVar:
				vd := m.ModuleListItems[i-1].StatementListItem.Statement.VariableStatement
				vd.VariableDeclarationList = append(vd.VariableDeclarationList, m.ModuleListItems[i].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0])
			}
			m.ModuleListItems = append(m.ModuleListItems[:i], m.ModuleListItems[i+1:]...)
			i--
		}
		last = next
	}
}

type bindable byte

const (
	bindableNone bindable = iota
	bindableConst
	bindableLet
	bindableVar
	bindableClass
	bindableFunction
	bindableBare
)

func sliBindable(sli *javascript.StatementListItem) bindable {
	if sli != nil {
		if sli.Declaration != nil {
			if sli.Declaration.LexicalDeclaration != nil {
				if sli.Declaration.LexicalDeclaration.LetOrConst == javascript.Const {
					return bindableConst
				}
				return bindableLet
			} else if sli.Declaration.ClassDeclaration != nil {
				return bindableClass
			} else if sli.Declaration.FunctionDeclaration != nil {
				return bindableFunction
			}
		}
		if sli.Statement != nil && sli.Statement.VariableStatement != nil {
			return bindableVar
		}
		if isStatementExpression(sli.Statement) {
			if sli.Statement.ExpressionStatement.Expressions[0].AssignmentOperator == javascript.AssignmentAssign {
				return bindableBare
			}
		}
	}
	return bindableNone
}

func removeDeadSLI(sli *javascript.StatementListItem) bool {
	switch sliBindable(sli) {
	case bindableConst, bindableLet:
		lb := sli.Declaration.LexicalDeclaration.BindingList[0]
		if lb.BindingIdentifier != nil {
			return lb.BindingIdentifier.Data == ""
		}
	case bindableVar:
		vd := sli.Statement.VariableStatement.VariableDeclarationList[0]
		if vd.BindingIdentifier != nil {
			return vd.BindingIdentifier.Data == ""
		}
	case bindableClass:
		return sli.Declaration.ClassDeclaration.BindingIdentifier == nil || sli.Declaration.ClassDeclaration.BindingIdentifier.Data == ""
	case bindableFunction:
		return sli.Declaration.FunctionDeclaration.BindingIdentifier == nil || sli.Declaration.FunctionDeclaration.BindingIdentifier.Data == ""
	}

	return false
}
