package scope

import (
	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/walk"
)

type scoper struct {
	bt    BindingType
	scope *Scope
	set   bool
}

func (s *scoper) newFunctionScope(t javascript.Type) *scoper {
	return s.newScoper(s.scope.newFunctionScope(t), BindingRef)
}

func (s *scoper) newScoper(t *Scope, bt BindingType) *scoper {
	return &scoper{
		bt:    bt,
		scope: t,
		set:   s.set,
	}
}

func (s *scoper) newArrowFunctionScope(t javascript.Type) *scoper {
	return s.newScoper(s.scope.newArrowFunctionScope(t), s.bt)
}

func (s *scoper) newLexicalScope(t javascript.Type) *scoper {
	return s.newScoper(s.scope.newLexicalScope(t), s.bt)
}

func (s *scoper) setBindingType(bt BindingType) *scoper {
	return s.newScoper(s.scope, bt)
}

func (s *scoper) Handle(t javascript.Type) error {
	switch t := t.(type) {
	case *javascript.ImportDeclaration:
		return s.processImportDeclaration(t)
	case *javascript.ExportDeclaration:
		return s.processExportDeclaration(t)
	case *javascript.ExportSpecifier:
		return s.processExportSpecifier(t)
	case *javascript.Statement:
		return s.processStatement(t)
	case *javascript.Declaration:
		return s.processDeclaration(t)
	case *javascript.IterationStatementFor:
		return s.processIterationStatementFor(t)
	case *javascript.SwitchStatement:
		return s.processSwitchStatement(t)
	case *javascript.FunctionDeclaration:
		return s.processFunctionDeclaration(t)
	case *javascript.TryStatement:
		return s.processTryStatement(t)
	case *javascript.ClassDeclaration:
		return s.processClassDeclaration(t)
	case *javascript.ClassElement:
		return s.processClassElement(t)
	case *javascript.LexicalDeclaration:
		return s.processLexicalDeclaration(t)
	case *javascript.VariableStatement:
		return s.processVariableStatment(t)
	case *javascript.LexicalBinding:
		return s.processLexicalBinding(t)
	case *javascript.AssignmentExpression:
		return s.processAssignmentExpression(t)
	case *javascript.ObjectAssignmentPattern:
		return s.processObjectAssignmentPattern(t)
	case *javascript.DestructuringAssignmentTarget:
		return s.processDestructuringAssignmentTarget(t)
	case *javascript.ObjectBindingPattern:
		return s.processObjectBindingPattern(t)
	case *javascript.FormalParameters:
		return s.processFormalParameters(t)
	case *javascript.MethodDefinition:
		return s.processMethodDefinition(t)
	case *javascript.ArrowFunction:
		return s.processArrowFunction(t)
	case *javascript.BindingElement:
		return s.processBindingElement(t)
	case *javascript.PrimaryExpression:
		return s.processPrimaryExpression(t)
	case *javascript.JSXElement:
		return s.processJSXElement(t)
	}

	return walk.Walk(t, s)
}

func (s *scoper) processImportDeclaration(t *javascript.ImportDeclaration) error {
	if !s.set || t.ImportClause == nil {
		return nil
	}

	if t.ImportedDefaultBinding != nil {
		if err := s.scope.setBinding(t.ImportedDefaultBinding, BindingImport); err != nil {
			return err
		}
	}

	if t.NameSpaceImport != nil {
		if err := s.scope.setBinding(t.NameSpaceImport, BindingImport); err != nil {
			return err
		}
	}

	if t.NamedImports != nil {
		for _, is := range t.NamedImports.ImportList {
			if is.ImportedBinding != nil {
				if err := s.scope.setBinding(is.ImportedBinding, BindingImport); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *scoper) processExportDeclaration(t *javascript.ExportDeclaration) error {
	if t.DefaultFunction != nil || t.DefaultClass != nil {
		return walk.Walk(t, s.setBindingType(BindingHoistable))
	}

	return walk.Walk(t, s)
}

func (s *scoper) processExportSpecifier(t *javascript.ExportSpecifier) error {
	if !s.set && t.IdentifierName != nil {
		s.scope.addBinding(t.IdentifierName, BindingRef)
	}

	return nil
}

func (s *scoper) processStatement(t *javascript.Statement) error {
	if t.BlockStatement != nil {
		s = s.newLexicalScope(t.BlockStatement)
	} else if t.LabelledItemFunction != nil {
		s = s.setBindingType(BindingHoistable)
	}

	return walk.Walk(t, s)
}

func (s *scoper) processDeclaration(t *javascript.Declaration) error {
	if t.FunctionDeclaration != nil || t.ClassDeclaration != nil {
		s = s.setBindingType(BindingHoistable)
	}

	return walk.Walk(t, s)
}

func (s *scoper) processIterationStatementFor(t *javascript.IterationStatementFor) error {
	s = s.newLexicalScope(t)

	switch t.Type {
	case javascript.ForInLet, javascript.ForOfLet, javascript.ForAwaitOfLet:
		s = s.setBindingType(BindingLexicalLet)
	case javascript.ForInConst, javascript.ForOfConst, javascript.ForAwaitOfConst:
		s = s.setBindingType(BindingLexicalConst)
	case javascript.ForNormalVar, javascript.ForInVar, javascript.ForOfVar, javascript.ForAwaitOfVar:
		s = s.setBindingType(BindingVar)
	}

	if s.set && t.ForBindingIdentifier != nil {
		if err := s.scope.setBinding(t.ForBindingIdentifier, s.bt); err != nil {
			return err
		}
	}

	return walk.Walk(t, s)
}

func (s *scoper) processSwitchStatement(t *javascript.SwitchStatement) error {
	if err := walk.Walk(&t.Expression, s); err != nil {
		return err
	}

	s = s.newLexicalScope(t)

	for n := range t.CaseClauses {
		if err := walk.Walk(&t.CaseClauses[n], s); err != nil {
			return err
		}
	}

	for n := range t.DefaultClause {
		if err := walk.Walk(&t.DefaultClause[n], s); err != nil {
			return err
		}
	}

	for n := range t.PostDefaultCaseClauses {
		if err := walk.Walk(&t.PostDefaultCaseClauses[n], s); err != nil {
			return err
		}
	}

	return nil
}

func (s *scoper) processFunctionDeclaration(t *javascript.FunctionDeclaration) error {
	if s.bt == BindingHoistable && s.set && t.BindingIdentifier != nil {
		if err := s.scope.setBinding(t.BindingIdentifier, BindingHoistable); err != nil {
			return err
		}
	}

	return walk.Walk(t, s.newFunctionScope(t))
}

func (s *scoper) processTryStatement(t *javascript.TryStatement) error {
	if err := walk.Walk(&t.TryBlock, s.newLexicalScope(&t.TryBlock)); err != nil {
		return err
	}

	if t.CatchBlock != nil {
		s = s.newLexicalScope(t.CatchBlock)

		if t.CatchParameterArrayBindingPattern != nil {
			if err := walk.Walk(t.CatchParameterArrayBindingPattern, s.setBindingType(BindingCatch)); err != nil {
				return err
			}
		} else if t.CatchParameterObjectBindingPattern != nil {
			if err := walk.Walk(t.CatchParameterObjectBindingPattern, s.setBindingType(BindingCatch)); err != nil {
				return err
			}
		} else if s.set && t.CatchParameterBindingIdentifier != nil {
			s.scope.setBinding(t.CatchParameterBindingIdentifier, BindingCatch)
		}

		if err := walk.Walk(t.CatchBlock, s); err != nil {
			return err
		}
	}

	if t.FinallyBlock != nil {
		return walk.Walk(t.FinallyBlock, s.newLexicalScope(t.FinallyBlock))
	}

	return nil
}

func (s *scoper) processClassDeclaration(t *javascript.ClassDeclaration) error {
	if s.bt == BindingHoistable && s.set && t.BindingIdentifier != nil {
		if err := s.scope.setBinding(t.BindingIdentifier, BindingHoistable); err != nil {
			return err
		}
	}

	return walk.Walk(t, s.setBindingType(BindingRef))
}

func (s *scoper) processClassElement(t *javascript.ClassElement) error {
	if t.ClassStaticBlock != nil {
		s = s.newLexicalScope(t.ClassStaticBlock)
	}

	return walk.Walk(t, s)
}

func (s *scoper) processLexicalDeclaration(t *javascript.LexicalDeclaration) error {
	typ := BindingLexicalLet

	if t.LetOrConst == javascript.Const {
		typ = BindingLexicalConst
	}

	return walk.Walk(t, s.setBindingType(typ))
}

func (s *scoper) processVariableStatment(t *javascript.VariableStatement) error {
	return walk.Walk(t, s.setBindingType(BindingVar))
}

func (s *scoper) processLexicalBinding(t *javascript.LexicalBinding) error {
	if s.set && t.BindingIdentifier != nil {
		if err := s.scope.setBinding(t.BindingIdentifier, s.bt); err != nil {
			return err
		}
	}

	return walk.Walk(t, s)
}

func (s *scoper) processAssignmentExpression(t *javascript.AssignmentExpression) error {
	if t.LeftHandSideExpression != nil {
		if err := s.processLeftHandSideExpressionAsAssignment(t.LeftHandSideExpression); err != nil {
			return err
		} else if t.AssignmentExpression == nil {
			return nil
		}

		return walk.Walk(t.AssignmentExpression, s)
	}

	return walk.Walk(t, s)
}

func (s *scoper) processLeftHandSideExpressionAsAssignment(t *javascript.LeftHandSideExpression) error {
	if t.NewExpression != nil && len(t.NewExpression.News) == 0 && t.NewExpression.MemberExpression.PrimaryExpression != nil && t.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference != nil {
		if !s.set {
			s.scope.addBinding(t.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference, BindingBare)
		}

		return nil
	}

	return walk.Walk(t, s)
}

func (s *scoper) processObjectAssignmentPattern(t *javascript.ObjectAssignmentPattern) error {
	for n := range t.AssignmentPropertyList {
		if err := walk.Walk(&t.AssignmentPropertyList[n], s); err != nil {
			return err
		}
	}

	if t.AssignmentRestElement != nil {
		return s.processLeftHandSideExpressionAsAssignment(t.AssignmentRestElement)
	}

	return nil
}

func (s *scoper) processDestructuringAssignmentTarget(t *javascript.DestructuringAssignmentTarget) error {
	if t.LeftHandSideExpression != nil {
		return s.processLeftHandSideExpressionAsAssignment(t.LeftHandSideExpression)
	} else if t.AssignmentPattern != nil {
		return walk.Walk(t.AssignmentPattern, s)
	}

	return nil
}

func (s *scoper) processObjectBindingPattern(t *javascript.ObjectBindingPattern) error {
	if err := walk.Walk(t, s); err != nil {
		return err
	} else if s.set && t.BindingRestProperty != nil {
		return s.scope.setBinding(t.BindingRestProperty, s.bt)
	}

	return nil
}

func (s *scoper) processFormalParameters(t *javascript.FormalParameters) error {
	if err := walk.Walk(t, s.setBindingType(BindingFunctionParam)); err != nil {
		return err
	} else if s.set && t.BindingIdentifier != nil {
		return s.scope.setBinding(t.BindingIdentifier, BindingFunctionParam)
	}

	return nil
}

func (s *scoper) processMethodDefinition(t *javascript.MethodDefinition) error {
	if err := walk.Walk(t.ClassElementName, s); err != nil {
		return err
	}

	s = s.newFunctionScope(t)

	if err := walk.Walk(&t.Params, s); err != nil {
		return err
	}

	return walk.Walk(&t.FunctionBody, s)
}

func (s *scoper) processArrowFunction(t *javascript.ArrowFunction) error {
	s = s.newArrowFunctionScope(t)

	if s.set && t.BindingIdentifier != nil {
		s.scope.setBinding(t.BindingIdentifier, BindingFunctionParam)
	}

	return walk.Walk(t, s)
}

func (s *scoper) processBindingElement(t *javascript.BindingElement) error {
	if s.set && t.SingleNameBinding != nil {
		if err := s.scope.setBinding(t.SingleNameBinding, s.bt); err != nil {
			return err
		}
	}

	return walk.Walk(t, s)
}

func (s *scoper) processPrimaryExpression(t *javascript.PrimaryExpression) error {
	if t.This != nil {
		if !s.set {
			s.scope.addBinding(t.This, BindingRef)
		}
	} else if t.IdentifierReference != nil {
		if !s.set {
			s.scope.addBinding(t.IdentifierReference, BindingRef)
		}
	} else {
		return walk.Walk(t, s)
	}

	return nil
}

func (s *scoper) processJSXElement(t *javascript.JSXElement) error {
	if !s.set && t.ElementName.Identifier != nil {
		s.scope.addBinding(t.ElementName.Identifier, BindingRef)
	}

	return walk.Walk(t, s)
}
