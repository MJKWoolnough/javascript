// Package walk provides a javascript type walker
package walk

import "vimagination.zapto.org/javascript"

// Handler is used to process javascript types.
type Handler interface {
	Handle(javascript.Type) error
}

// HandlerFunc wraps a func to implement Handler interface.
type HandlerFunc func(javascript.Type) error

// Handle implements the Handler interface.
func (h HandlerFunc) Handle(t javascript.Type) error {
	return h(t)
}

// Walk calls the Handle function on the given interface for each non-nil, non-Token field of the given javascript type.
func Walk(t javascript.Type, h Handler) error {
	switch t := t.(type) {
	case javascript.ClassDeclaration:
		return walkClassDeclaration(&t, h)
	case *javascript.ClassDeclaration:
		return walkClassDeclaration(t, h)
	case javascript.ClassElement:
		return walkClassElement(&t, h)
	case *javascript.ClassElement:
		return walkClassElement(t, h)
	case javascript.FieldDefinition:
		return walkFieldDefinition(&t, h)
	case *javascript.FieldDefinition:
		return walkFieldDefinition(t, h)
	case javascript.ClassElementName:
		return walkClassElementName(&t, h)
	case *javascript.ClassElementName:
		return walkClassElementName(t, h)
	case javascript.MethodDefinition:
		return walkMethodDefinition(&t, h)
	case *javascript.MethodDefinition:
		return walkMethodDefinition(t, h)
	case javascript.PropertyName:
		return walkPropertyName(&t, h)
	case *javascript.PropertyName:
		return walkPropertyName(t, h)
	case javascript.ConditionalExpression:
		return walkConditionalExpression(&t, h)
	case *javascript.ConditionalExpression:
		return walkConditionalExpression(t, h)
	case javascript.CoalesceExpression:
		return walkCoalesceExpression(&t, h)
	case *javascript.CoalesceExpression:
		return walkCoalesceExpression(t, h)
	case javascript.LogicalORExpression:
		return walkLogicalORExpression(&t, h)
	case *javascript.LogicalORExpression:
		return walkLogicalORExpression(t, h)
	case javascript.LogicalANDExpression:
		return walkLogicalANDExpression(&t, h)
	case *javascript.LogicalANDExpression:
		return walkLogicalANDExpression(t, h)
	case javascript.BitwiseORExpression:
		return walkBitwiseORExpression(&t, h)
	case *javascript.BitwiseORExpression:
		return walkBitwiseORExpression(t, h)
	case javascript.BitwiseXORExpression:
		return walkBitwiseXORExpression(&t, h)
	case *javascript.BitwiseXORExpression:
		return walkBitwiseXORExpression(t, h)
	case javascript.BitwiseANDExpression:
		return walkBitwiseANDExpression(&t, h)
	case *javascript.BitwiseANDExpression:
		return walkBitwiseANDExpression(t, h)
	case javascript.EqualityExpression:
		return walkEqualityExpression(&t, h)
	case *javascript.EqualityExpression:
		return walkEqualityExpression(t, h)
	case javascript.RelationalExpression:
		return walkRelationalExpression(&t, h)
	case *javascript.RelationalExpression:
		return walkRelationalExpression(t, h)
	case javascript.ShiftExpression:
		return walkShiftExpression(&t, h)
	case *javascript.ShiftExpression:
		return walkShiftExpression(t, h)
	case javascript.AdditiveExpression:
		return walkAdditiveExpression(&t, h)
	case *javascript.AdditiveExpression:
		return walkAdditiveExpression(t, h)
	case javascript.MultiplicativeExpression:
		return walkMultiplicativeExpression(&t, h)
	case *javascript.MultiplicativeExpression:
		return walkMultiplicativeExpression(t, h)
	case javascript.ExponentiationExpression:
		return walkExponentiationExpression(&t, h)
	case *javascript.ExponentiationExpression:
		return walkExponentiationExpression(t, h)
	case javascript.UnaryExpression:
		return walkUnaryExpression(&t, h)
	case *javascript.UnaryExpression:
		return walkUnaryExpression(t, h)
	case javascript.UpdateExpression:
		return walkUpdateExpression(&t, h)
	case *javascript.UpdateExpression:
		return walkUpdateExpression(t, h)
	case javascript.AssignmentExpression:
		return walkAssignmentExpression(&t, h)
	case *javascript.AssignmentExpression:
		return walkAssignmentExpression(t, h)
	case javascript.LeftHandSideExpression:
		return walkLeftHandSideExpression(&t, h)
	case *javascript.LeftHandSideExpression:
		return walkLeftHandSideExpression(t, h)
	case *javascript.AssignmentPattern:
		return walkAssignmentPattern(t, h)
	case javascript.AssignmentPattern:
		return walkAssignmentPattern(&t, h)
	case *javascript.ObjectAssignmentPattern:
		return walkObjectAssignmentPattern(t, h)
	case javascript.ObjectAssignmentPattern:
		return walkObjectAssignmentPattern(&t, h)
	case *javascript.AssignmentProperty:
		return walkAssignmentProperty(t, h)
	case javascript.AssignmentProperty:
		return walkAssignmentProperty(&t, h)
	case *javascript.DestructuringAssignmentTarget:
		return walkDestructuringAssignmentTarget(t, h)
	case javascript.DestructuringAssignmentTarget:
		return walkDestructuringAssignmentTarget(&t, h)
	case *javascript.AssignmentElement:
		return walkAssignmentElement(t, h)
	case javascript.AssignmentElement:
		return walkAssignmentElement(&t, h)
	case *javascript.ArrayAssignmentPattern:
		return walkArrayAssignmentPattern(t, h)
	case javascript.ArrayAssignmentPattern:
		return walkArrayAssignmentPattern(&t, h)
	case javascript.OptionalExpression:
		return walkOptionalExpression(&t, h)
	case *javascript.OptionalExpression:
		return walkOptionalExpression(t, h)
	case javascript.OptionalChain:
		return walkOptionalChain(&t, h)
	case *javascript.OptionalChain:
		return walkOptionalChain(t, h)
	case javascript.Expression:
		return walkExpression(&t, h)
	case *javascript.Expression:
		return walkExpression(t, h)
	case javascript.NewExpression:
		return walkNewExpression(&t, h)
	case *javascript.NewExpression:
		return walkNewExpression(t, h)
	case javascript.MemberExpression:
		return walkMemberExpression(&t, h)
	case *javascript.MemberExpression:
		return walkMemberExpression(t, h)
	case javascript.PrimaryExpression:
		return walkPrimaryExpression(&t, h)
	case *javascript.PrimaryExpression:
		return walkPrimaryExpression(t, h)
	case javascript.ParenthesizedExpression:
		return walkParenthesizedExpression(&t, h)
	case *javascript.ParenthesizedExpression:
		return walkParenthesizedExpression(t, h)
	case javascript.Arguments:
		return walkArguments(&t, h)
	case *javascript.Arguments:
		return walkArguments(t, h)
	case javascript.Argument:
		return walkArgument(&t, h)
	case *javascript.Argument:
		return walkArgument(t, h)
	case javascript.CallExpression:
		return walkCallExpression(&t, h)
	case *javascript.CallExpression:
		return walkCallExpression(t, h)
	case javascript.FunctionDeclaration:
		return walkFunctionDeclaration(&t, h)
	case *javascript.FunctionDeclaration:
		return walkFunctionDeclaration(t, h)
	case javascript.FormalParameters:
		return walkFormalParameters(&t, h)
	case *javascript.FormalParameters:
		return walkFormalParameters(t, h)
	case javascript.BindingElement:
		return walkBindingElement(&t, h)
	case *javascript.BindingElement:
		return walkBindingElement(t, h)
	case javascript.Script:
		return walkScript(&t, h)
	case *javascript.Script:
		return walkScript(t, h)
	case javascript.Declaration:
		return walkDeclaration(&t, h)
	case *javascript.Declaration:
		return walkDeclaration(t, h)
	case javascript.LexicalDeclaration:
		return walkLexicalDeclaration(&t, h)
	case *javascript.LexicalDeclaration:
		return walkLexicalDeclaration(t, h)
	case javascript.LexicalBinding:
		return walkLexicalBinding(&t, h)
	case *javascript.LexicalBinding:
		return walkLexicalBinding(t, h)
	case javascript.ArrayBindingPattern:
		return walkArrayBindingPattern(&t, h)
	case *javascript.ArrayBindingPattern:
		return walkArrayBindingPattern(t, h)
	case javascript.ObjectBindingPattern:
		return walkObjectBindingPattern(&t, h)
	case *javascript.ObjectBindingPattern:
		return walkObjectBindingPattern(t, h)
	case javascript.BindingProperty:
		return walkBindingProperty(&t, h)
	case *javascript.BindingProperty:
		return walkBindingProperty(t, h)
	case javascript.ArrayElement:
		return walkArrayElement(&t, h)
	case *javascript.ArrayElement:
		return walkArrayElement(t, h)
	case javascript.ArrayLiteral:
		return walkArrayLiteral(&t, h)
	case *javascript.ArrayLiteral:
		return walkArrayLiteral(t, h)
	case javascript.ObjectLiteral:
		return walkObjectLiteral(&t, h)
	case *javascript.ObjectLiteral:
		return walkObjectLiteral(t, h)
	case javascript.PropertyDefinition:
		return walkPropertyDefinition(&t, h)
	case *javascript.PropertyDefinition:
		return walkPropertyDefinition(t, h)
	case javascript.TemplateLiteral:
		return walkTemplateLiteral(&t, h)
	case *javascript.TemplateLiteral:
		return walkTemplateLiteral(t, h)
	case javascript.ArrowFunction:
		return walkArrowFunction(&t, h)
	case *javascript.ArrowFunction:
		return walkArrowFunction(t, h)
	case javascript.Module:
		return walkModule(&t, h)
	case *javascript.Module:
		return walkModule(t, h)
	case javascript.ModuleItem:
		return walkModuleItem(&t, h)
	case *javascript.ModuleItem:
		return walkModuleItem(t, h)
	case javascript.ImportDeclaration:
		return walkImportDeclaration(&t, h)
	case *javascript.ImportDeclaration:
		return walkImportDeclaration(t, h)
	case javascript.ImportClause:
		return walkImportClause(&t, h)
	case *javascript.ImportClause:
		return walkImportClause(t, h)
	case javascript.FromClause:
		return walkFromClause(&t, h)
	case *javascript.FromClause:
		return walkFromClause(t, h)
	case javascript.NamedImports:
		return walkNamedImports(&t, h)
	case *javascript.NamedImports:
		return walkNamedImports(t, h)
	case javascript.ImportSpecifier:
		return walkImportSpecifier(&t, h)
	case *javascript.ImportSpecifier:
		return walkImportSpecifier(t, h)
	case javascript.ExportDeclaration:
		return walkExportDeclaration(&t, h)
	case *javascript.ExportDeclaration:
		return walkExportDeclaration(t, h)
	case javascript.ExportClause:
		return walkExportClause(&t, h)
	case *javascript.ExportClause:
		return walkExportClause(t, h)
	case javascript.ExportSpecifier:
		return walkExportSpecifier(&t, h)
	case *javascript.ExportSpecifier:
		return walkExportSpecifier(t, h)
	case javascript.Block:
		return walkBlock(&t, h)
	case *javascript.Block:
		return walkBlock(t, h)
	case javascript.StatementListItem:
		return walkStatementListItem(&t, h)
	case *javascript.StatementListItem:
		return walkStatementListItem(t, h)
	case javascript.Statement:
		return walkStatement(&t, h)
	case *javascript.Statement:
		return walkStatement(t, h)
	case javascript.IfStatement:
		return walkIfStatement(&t, h)
	case *javascript.IfStatement:
		return walkIfStatement(t, h)
	case javascript.IterationStatementDo:
		return walkIterationStatementDo(&t, h)
	case *javascript.IterationStatementDo:
		return walkIterationStatementDo(t, h)
	case javascript.IterationStatementWhile:
		return walkIterationStatementWhile(&t, h)
	case *javascript.IterationStatementWhile:
		return walkIterationStatementWhile(t, h)
	case javascript.IterationStatementFor:
		return walkIterationStatementFor(&t, h)
	case *javascript.IterationStatementFor:
		return walkIterationStatementFor(t, h)
	case javascript.SwitchStatement:
		return walkSwitchStatement(&t, h)
	case *javascript.SwitchStatement:
		return walkSwitchStatement(t, h)
	case javascript.CaseClause:
		return walkCaseClause(&t, h)
	case *javascript.CaseClause:
		return walkCaseClause(t, h)
	case javascript.WithStatement:
		return walkWithStatement(&t, h)
	case *javascript.WithStatement:
		return walkWithStatement(t, h)
	case javascript.TryStatement:
		return walkTryStatement(&t, h)
	case *javascript.TryStatement:
		return walkTryStatement(t, h)
	case javascript.VariableStatement:
		return walkVariableStatement(&t, h)
	case *javascript.VariableStatement:
		return walkVariableStatement(t, h)
	case javascript.JSXElement:
		return walkJSXElement(&t, h)
	case *javascript.JSXElement:
		return walkJSXElement(t, h)
	case javascript.JSXFragment:
		return walkJSXFragment(&t, h)
	case *javascript.JSXFragment:
		return walkJSXFragment(t, h)
	case javascript.JSXElementName:
		return walkJSXElementName(&t, h)
	case *javascript.JSXElementName:
		return walkJSXElementName(t, h)
	case javascript.JSXAttribute:
		return walkJSXAttribute(&t, h)
	case *javascript.JSXAttribute:
		return walkJSXAttribute(t, h)
	case javascript.JSXChild:
		return walkJSXChild(&t, h)
	case *javascript.JSXChild:
		return walkJSXChild(t, h)
	}

	return nil
}

func walkClassDeclaration(t *javascript.ClassDeclaration, h Handler) error {
	if t.ClassHeritage != nil {
		if err := h.Handle(t.ClassHeritage); err != nil {
			return err
		}
	}

	for n := range t.ClassBody {
		if err := h.Handle(&t.ClassBody[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkClassElement(t *javascript.ClassElement, h Handler) error {
	if t.FieldDefinition != nil {
		if err := h.Handle(t.FieldDefinition); err != nil {
			return err
		}
	} else if t.MethodDefinition != nil {
		if err := h.Handle(t.MethodDefinition); err != nil {
			return err
		}
	} else if t.ClassStaticBlock != nil {
		if err := h.Handle(t.ClassStaticBlock); err != nil {
			return err
		}
	}

	return nil
}

func walkFieldDefinition(t *javascript.FieldDefinition, h Handler) error {
	if err := h.Handle(&t.ClassElementName); err != nil {
		return err
	} else if t.Initializer != nil {
		if err := h.Handle(t.Initializer); err != nil {
			return err
		}
	}

	return nil
}

func walkClassElementName(t *javascript.ClassElementName, h Handler) error {
	if t.PropertyName != nil {
		if err := h.Handle(t.PropertyName); err != nil {
			return err
		}
	}

	return nil
}

func walkMethodDefinition(t *javascript.MethodDefinition, h Handler) error {
	if err := h.Handle(&t.ClassElementName); err != nil {
		return err
	} else if err = h.Handle(&t.Params); err != nil {
		return err
	}

	return h.Handle(t.FunctionBody)
}

func walkPropertyName(t *javascript.PropertyName, h Handler) error {
	if t.ComputedPropertyName != nil {
		return h.Handle(t.ComputedPropertyName)
	}

	return nil
}

func walkConditionalExpression(t *javascript.ConditionalExpression, h Handler) error {
	if t.LogicalORExpression != nil {
		if err := h.Handle(t.LogicalORExpression); err != nil {
			return err
		}
	}

	if t.CoalesceExpression != nil {
		if err := h.Handle(t.CoalesceExpression); err != nil {
			return err
		}
	}

	if t.True != nil {
		if err := h.Handle(t.True); err != nil {
			return err
		}
	}

	if t.False != nil {
		if err := h.Handle(t.False); err != nil {
			return err
		}
	}

	return nil
}

func walkCoalesceExpression(t *javascript.CoalesceExpression, h Handler) error {
	if t.CoalesceExpressionHead != nil {
		if err := h.Handle(t.CoalesceExpressionHead); err != nil {
			return err
		}
	}

	return h.Handle(&t.BitwiseORExpression)
}

func walkLogicalORExpression(t *javascript.LogicalORExpression, h Handler) error {
	if t.LogicalORExpression != nil {
		if err := h.Handle(t.LogicalORExpression); err != nil {
			return err
		}
	}

	return h.Handle(&t.LogicalANDExpression)
}

func walkLogicalANDExpression(t *javascript.LogicalANDExpression, h Handler) error {
	if t.LogicalANDExpression != nil {
		if err := h.Handle(t.LogicalANDExpression); err != nil {
			return err
		}
	}

	return h.Handle(&t.BitwiseORExpression)
}

func walkBitwiseORExpression(t *javascript.BitwiseORExpression, h Handler) error {
	if t.BitwiseORExpression != nil {
		if err := h.Handle(t.BitwiseORExpression); err != nil {
			return err
		}
	}

	return h.Handle(&t.BitwiseXORExpression)
}

func walkBitwiseXORExpression(t *javascript.BitwiseXORExpression, h Handler) error {
	if t.BitwiseXORExpression != nil {
		if err := h.Handle(t.BitwiseXORExpression); err != nil {
			return err
		}
	}

	return h.Handle(&t.BitwiseANDExpression)
}

func walkBitwiseANDExpression(t *javascript.BitwiseANDExpression, h Handler) error {
	if t.BitwiseANDExpression != nil {
		if err := h.Handle(t.BitwiseANDExpression); err != nil {
			return err
		}
	}

	return h.Handle(&t.EqualityExpression)
}

func walkEqualityExpression(t *javascript.EqualityExpression, h Handler) error {
	if t.EqualityExpression != nil {
		if err := h.Handle(t.EqualityExpression); err != nil {
			return err
		}
	}

	return h.Handle(&t.RelationalExpression)
}

func walkRelationalExpression(t *javascript.RelationalExpression, h Handler) error {
	if t.RelationalExpression != nil {
		if err := h.Handle(t.RelationalExpression); err != nil {
			return err
		}
	}

	return h.Handle(&t.ShiftExpression)
}

func walkShiftExpression(t *javascript.ShiftExpression, h Handler) error {
	if t.ShiftExpression != nil {
		if err := h.Handle(t.ShiftExpression); err != nil {
			return err
		}
	}

	return h.Handle(&t.AdditiveExpression)
}

func walkAdditiveExpression(t *javascript.AdditiveExpression, h Handler) error {
	if t.AdditiveExpression != nil {
		if err := h.Handle(t.AdditiveExpression); err != nil {
			return err
		}
	}

	return h.Handle(&t.MultiplicativeExpression)
}

func walkMultiplicativeExpression(t *javascript.MultiplicativeExpression, h Handler) error {
	if t.MultiplicativeExpression != nil {
		if err := h.Handle(t.MultiplicativeExpression); err != nil {
			return err
		}
	}

	return h.Handle(&t.ExponentiationExpression)
}

func walkExponentiationExpression(t *javascript.ExponentiationExpression, h Handler) error {
	if t.ExponentiationExpression != nil {
		if err := h.Handle(t.ExponentiationExpression); err != nil {
			return err
		}
	}

	return h.Handle(&t.UnaryExpression)
}

func walkUnaryExpression(t *javascript.UnaryExpression, h Handler) error {
	return h.Handle(&t.UpdateExpression)
}

func walkUpdateExpression(t *javascript.UpdateExpression, h Handler) error {
	if t.LeftHandSideExpression != nil {
		if err := h.Handle(t.LeftHandSideExpression); err != nil {
			return err
		}
	}

	if t.UnaryExpression != nil {
		if err := h.Handle(t.UnaryExpression); err != nil {
			return err
		}
	}

	return nil
}

func walkAssignmentExpression(t *javascript.AssignmentExpression, h Handler) error {
	if t.ConditionalExpression != nil {
		if err := h.Handle(t.ConditionalExpression); err != nil {
			return err
		}
	}

	if t.ArrowFunction != nil {
		if err := h.Handle(t.ArrowFunction); err != nil {
			return err
		}
	}

	if t.LeftHandSideExpression != nil {
		if err := h.Handle(t.LeftHandSideExpression); err != nil {
			return err
		}
	}

	if t.AssignmentPattern != nil {
		if err := h.Handle(t.AssignmentPattern); err != nil {
			return err
		}
	}

	if t.AssignmentExpression != nil {
		if err := h.Handle(t.AssignmentExpression); err != nil {
			return err
		}
	}

	return nil
}

func walkLeftHandSideExpression(t *javascript.LeftHandSideExpression, h Handler) error {
	if t.NewExpression != nil {
		if err := h.Handle(t.NewExpression); err != nil {
			return err
		}
	}

	if t.CallExpression != nil {
		if err := h.Handle(t.CallExpression); err != nil {
			return err
		}
	}

	if t.OptionalExpression != nil {
		if err := h.Handle(t.OptionalExpression); err != nil {
			return err
		}
	}

	return nil
}

func walkAssignmentPattern(t *javascript.AssignmentPattern, h Handler) error {
	if t.ArrayAssignmentPattern != nil {
		if err := h.Handle(t.ArrayAssignmentPattern); err != nil {
			return err
		}
	}

	if t.ObjectAssignmentPattern != nil {
		if err := h.Handle(t.ObjectAssignmentPattern); err != nil {
			return err
		}
	}

	return nil
}

func walkObjectAssignmentPattern(t *javascript.ObjectAssignmentPattern, h Handler) error {
	for n := range t.AssignmentPropertyList {
		if err := h.Handle(&t.AssignmentPropertyList[n]); err != nil {
			return err
		}
	}

	if t.AssignmentRestElement != nil {
		if err := h.Handle(t.AssignmentRestElement); err != nil {
			return err
		}
	}

	return nil
}

func walkAssignmentProperty(t *javascript.AssignmentProperty, h Handler) error {
	if err := h.Handle(&t.PropertyName); err != nil {
		return err
	} else if t.DestructuringAssignmentTarget != nil {
		if err := h.Handle(t.DestructuringAssignmentTarget); err != nil {
			return err
		}
	}

	if t.Initializer != nil {
		if err := h.Handle(t.Initializer); err != nil {
			return err
		}
	}

	return nil
}

func walkDestructuringAssignmentTarget(t *javascript.DestructuringAssignmentTarget, h Handler) error {
	if t.LeftHandSideExpression != nil {
		if err := h.Handle(t.LeftHandSideExpression); err != nil {
			return err
		}
	}

	if t.AssignmentPattern != nil {
		if err := h.Handle(t.AssignmentPattern); err != nil {
			return err
		}
	}

	return nil
}

func walkAssignmentElement(t *javascript.AssignmentElement, h Handler) error {
	if err := h.Handle(&t.DestructuringAssignmentTarget); err != nil {
		return err
	} else if t.Initializer != nil {
		if err := h.Handle(t.Initializer); err != nil {
			return err
		}
	}

	return nil
}

func walkArrayAssignmentPattern(t *javascript.ArrayAssignmentPattern, h Handler) error {
	for n := range t.AssignmentElements {
		if err := h.Handle(&t.AssignmentElements[n]); err != nil {
			return err
		}
	}

	if t.AssignmentRestElement != nil {
		if err := h.Handle(t.AssignmentRestElement); err != nil {
			return err
		}
	}

	return nil
}

func walkOptionalExpression(t *javascript.OptionalExpression, h Handler) error {
	if t.MemberExpression != nil {
		if err := h.Handle(t.MemberExpression); err != nil {
			return err
		}
	}

	if t.CallExpression != nil {
		if err := h.Handle(t.CallExpression); err != nil {
			return err
		}
	}

	if t.OptionalExpression != nil {
		if err := h.Handle(t.OptionalExpression); err != nil {
			return err
		}
	}

	return h.Handle(&t.OptionalChain)
}

func walkOptionalChain(t *javascript.OptionalChain, h Handler) error {
	if t.OptionalChain != nil {
		if err := h.Handle(t.OptionalChain); err != nil {
			return err
		}
	}

	if t.Arguments != nil {
		if err := h.Handle(t.Arguments); err != nil {
			return err
		}
	}

	if t.Expression != nil {
		if err := h.Handle(t.Expression); err != nil {
			return err
		}
	}

	if t.TemplateLiteral != nil {
		if err := h.Handle(t.TemplateLiteral); err != nil {
			return err
		}
	}

	return nil
}

func walkExpression(t *javascript.Expression, h Handler) error {
	for n := range t.Expressions {
		if err := h.Handle(&t.Expressions[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkNewExpression(t *javascript.NewExpression, h Handler) error {
	return h.Handle(&t.MemberExpression)
}

func walkMemberExpression(t *javascript.MemberExpression, h Handler) error {
	if t.MemberExpression != nil {
		if err := h.Handle(t.MemberExpression); err != nil {
			return err
		}
	}

	if t.PrimaryExpression != nil {
		if err := h.Handle(t.PrimaryExpression); err != nil {
			return err
		}
	}

	if t.Expression != nil {
		if err := h.Handle(t.Expression); err != nil {
			return err
		}
	}

	if t.TemplateLiteral != nil {
		if err := h.Handle(t.TemplateLiteral); err != nil {
			return err
		}
	}

	if t.Arguments != nil {
		if err := h.Handle(t.Arguments); err != nil {
			return err
		}
	}

	return nil
}

func walkPrimaryExpression(t *javascript.PrimaryExpression, h Handler) error {
	if t.ArrayLiteral != nil {
		if err := h.Handle(t.ArrayLiteral); err != nil {
			return err
		}
	}

	if t.ObjectLiteral != nil {
		if err := h.Handle(t.ObjectLiteral); err != nil {
			return err
		}
	}

	if t.FunctionExpression != nil {
		if err := h.Handle(t.FunctionExpression); err != nil {
			return err
		}
	}

	if t.ClassExpression != nil {
		if err := h.Handle(t.ClassExpression); err != nil {
			return err
		}
	}

	if t.TemplateLiteral != nil {
		if err := h.Handle(t.TemplateLiteral); err != nil {
			return err
		}
	}

	if t.ParenthesizedExpression != nil {
		if err := h.Handle(t.ParenthesizedExpression); err != nil {
			return err
		}
	}

	if t.JSXElement != nil {
		if err := h.Handle(t.JSXElement); err != nil {
			return err
		}
	}

	if t.JSXFragment != nil {
		if err := h.Handle(t.JSXFragment); err != nil {
			return err
		}
	}

	return nil
}

func walkParenthesizedExpression(t *javascript.ParenthesizedExpression, h Handler) error {
	for n := range t.Expressions {
		if err := h.Handle(&t.Expressions[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkArguments(t *javascript.Arguments, h Handler) error {
	for n := range t.ArgumentList {
		if err := h.Handle(&t.ArgumentList[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkArgument(t *javascript.Argument, h Handler) error {
	return h.Handle(&t.AssignmentExpression)
}

func walkCallExpression(t *javascript.CallExpression, h Handler) error {
	if t.MemberExpression != nil {
		if err := h.Handle(t.MemberExpression); err != nil {
			return err
		}
	}

	if t.ImportCall != nil {
		if err := h.Handle(t.ImportCall); err != nil {
			return err
		}
	}

	if t.CallExpression != nil {
		if err := h.Handle(t.CallExpression); err != nil {
			return err
		}
	}

	if t.Arguments != nil {
		if err := h.Handle(t.Arguments); err != nil {
			return err
		}
	}

	if t.Expression != nil {
		if err := h.Handle(t.Expression); err != nil {
			return err
		}
	}

	if t.TemplateLiteral != nil {
		if err := h.Handle(t.TemplateLiteral); err != nil {
			return err
		}
	}

	return nil
}

func walkFunctionDeclaration(t *javascript.FunctionDeclaration, h Handler) error {
	if err := h.Handle(&t.FormalParameters); err != nil {
		return err
	}

	return h.Handle(&t.FunctionBody)
}

func walkFormalParameters(t *javascript.FormalParameters, h Handler) error {
	for n := range t.FormalParameterList {
		if err := h.Handle(&t.FormalParameterList[n]); err != nil {
			return err
		}
	}

	if t.ArrayBindingPattern != nil {
		if err := h.Handle(t.ArrayBindingPattern); err != nil {
			return err
		}
	}

	if t.ObjectBindingPattern != nil {
		if err := h.Handle(t.ObjectBindingPattern); err != nil {
			return err
		}
	}

	return nil
}

func walkBindingElement(t *javascript.BindingElement, h Handler) error {
	if t.ArrayBindingPattern != nil {
		if err := h.Handle(t.ArrayBindingPattern); err != nil {
			return err
		}
	}

	if t.ObjectBindingPattern != nil {
		if err := h.Handle(t.ObjectBindingPattern); err != nil {
			return err
		}
	}

	if t.Initializer != nil {
		if err := h.Handle(t.Initializer); err != nil {
			return err
		}
	}

	return nil
}

func walkScript(t *javascript.Script, h Handler) error {
	for n := range t.StatementList {
		if err := h.Handle(&t.StatementList[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkDeclaration(t *javascript.Declaration, h Handler) error {
	if t.ClassDeclaration != nil {
		if err := h.Handle(t.ClassDeclaration); err != nil {
			return err
		}
	}

	if t.FunctionDeclaration != nil {
		if err := h.Handle(t.FunctionDeclaration); err != nil {
			return err
		}
	}

	if t.LexicalDeclaration != nil {
		if err := h.Handle(t.LexicalDeclaration); err != nil {
			return err
		}
	}

	return nil
}

func walkLexicalDeclaration(t *javascript.LexicalDeclaration, h Handler) error {
	for n := range t.BindingList {
		if err := h.Handle(&t.BindingList[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkLexicalBinding(t *javascript.LexicalBinding, h Handler) error {
	if t.ArrayBindingPattern != nil {
		if err := h.Handle(t.ArrayBindingPattern); err != nil {
			return err
		}
	}

	if t.ObjectBindingPattern != nil {
		if err := h.Handle(t.ObjectBindingPattern); err != nil {
			return err
		}
	}

	if t.Initializer != nil {
		if err := h.Handle(t.Initializer); err != nil {
			return err
		}
	}

	return nil
}

func walkArrayBindingPattern(t *javascript.ArrayBindingPattern, h Handler) error {
	for n := range t.BindingElementList {
		if err := h.Handle(&t.BindingElementList[n]); err != nil {
			return err
		}
	}

	if t.BindingRestElement != nil {
		if err := h.Handle(t.BindingRestElement); err != nil {
			return err
		}
	}

	return nil
}

func walkObjectBindingPattern(t *javascript.ObjectBindingPattern, h Handler) error {
	for n := range t.BindingPropertyList {
		if err := h.Handle(&t.BindingPropertyList[n]); err != nil {
			return err
		}
	}

	if t.BindingRestProperty != nil {
		if err := h.Handle(t.BindingRestProperty); err != nil {
			return err
		}
	}

	return nil
}

func walkBindingProperty(t *javascript.BindingProperty, h Handler) error {
	if err := h.Handle(&t.PropertyName); err != nil {
		return err
	}

	return h.Handle(&t.BindingElement)
}

func walkArrayElement(t *javascript.ArrayElement, h Handler) error {
	return h.Handle(&t.AssignmentExpression)
}

func walkArrayLiteral(t *javascript.ArrayLiteral, h Handler) error {
	for n := range t.ElementList {
		if err := h.Handle(&t.ElementList[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkObjectLiteral(t *javascript.ObjectLiteral, h Handler) error {
	for n := range t.PropertyDefinitionList {
		if err := h.Handle(&t.PropertyDefinitionList[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkPropertyDefinition(t *javascript.PropertyDefinition, h Handler) error {
	if t.PropertyName != nil {
		if err := h.Handle(t.PropertyName); err != nil {
			return err
		}
	}

	if t.AssignmentExpression != nil {
		if err := h.Handle(t.AssignmentExpression); err != nil {
			return err
		}
	}

	if t.MethodDefinition != nil {
		if err := h.Handle(t.MethodDefinition); err != nil {
			return err
		}
	}

	return nil
}

func walkTemplateLiteral(t *javascript.TemplateLiteral, h Handler) error {
	for n := range t.Expressions {
		if err := h.Handle(&t.Expressions[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkArrowFunction(t *javascript.ArrowFunction, h Handler) error {
	if t.FormalParameters != nil {
		if err := h.Handle(t.FormalParameters); err != nil {
			return err
		}
	}

	if t.AssignmentExpression != nil {
		if err := h.Handle(t.AssignmentExpression); err != nil {
			return err
		}
	}

	if t.FunctionBody != nil {
		if err := h.Handle(t.FunctionBody); err != nil {
			return err
		}
	}

	return nil
}

func walkModule(t *javascript.Module, h Handler) error {
	for n := range t.ModuleListItems {
		if err := h.Handle(&t.ModuleListItems[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkModuleItem(t *javascript.ModuleItem, h Handler) error {
	if t.ImportDeclaration != nil {
		if err := h.Handle(t.ImportDeclaration); err != nil {
			return err
		}
	}

	if t.StatementListItem != nil {
		if err := h.Handle(t.StatementListItem); err != nil {
			return err
		}
	}

	if t.ExportDeclaration != nil {
		if err := h.Handle(t.ExportDeclaration); err != nil {
			return err
		}
	}

	return nil
}

func walkImportDeclaration(t *javascript.ImportDeclaration, h Handler) error {
	if t.ImportClause != nil {
		if err := h.Handle(t.ImportClause); err != nil {
			return err
		}
	}

	return h.Handle(&t.FromClause)
}

func walkImportClause(t *javascript.ImportClause, h Handler) error {
	if t.NamedImports != nil {
		if err := h.Handle(t.NamedImports); err != nil {
			return err
		}
	}

	return nil
}

func walkFromClause(_ *javascript.FromClause, _ Handler) error {
	return nil
}

func walkNamedImports(t *javascript.NamedImports, h Handler) error {
	for n := range t.ImportList {
		if err := h.Handle(&t.ImportList[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkImportSpecifier(_ *javascript.ImportSpecifier, _ Handler) error {
	return nil
}

func walkExportDeclaration(t *javascript.ExportDeclaration, h Handler) error {
	if t.ExportClause != nil {
		if err := h.Handle(t.ExportClause); err != nil {
			return err
		}
	}

	if t.FromClause != nil {
		if err := h.Handle(t.FromClause); err != nil {
			return err
		}
	}

	if t.VariableStatement != nil {
		if err := h.Handle(t.VariableStatement); err != nil {
			return err
		}
	}

	if t.Declaration != nil {
		if err := h.Handle(t.Declaration); err != nil {
			return err
		}
	}

	if t.DefaultFunction != nil {
		if err := h.Handle(t.DefaultFunction); err != nil {
			return err
		}
	}

	if t.DefaultClass != nil {
		if err := h.Handle(t.DefaultClass); err != nil {
			return err
		}
	}

	if t.DefaultAssignmentExpression != nil {
		if err := h.Handle(t.DefaultAssignmentExpression); err != nil {
			return err
		}
	}

	return nil
}

func walkExportClause(t *javascript.ExportClause, h Handler) error {
	for n := range t.ExportList {
		if err := h.Handle(&t.ExportList[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkExportSpecifier(_ *javascript.ExportSpecifier, _ Handler) error {
	return nil
}

func walkBlock(t *javascript.Block, h Handler) error {
	for n := range t.StatementList {
		if err := h.Handle(&t.StatementList[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkStatementListItem(t *javascript.StatementListItem, h Handler) error {
	if t.Statement != nil {
		if err := h.Handle(t.Statement); err != nil {
			return err
		}
	}

	if t.Declaration != nil {
		if err := h.Handle(t.Declaration); err != nil {
			return err
		}
	}

	return nil
}

func walkStatement(t *javascript.Statement, h Handler) error {
	if t.BlockStatement != nil {
		if err := h.Handle(t.BlockStatement); err != nil {
			return err
		}
	}

	if t.VariableStatement != nil {
		if err := h.Handle(t.VariableStatement); err != nil {
			return err
		}
	}

	if t.ExpressionStatement != nil {
		if err := h.Handle(t.ExpressionStatement); err != nil {
			return err
		}
	}

	if t.IfStatement != nil {
		if err := h.Handle(t.IfStatement); err != nil {
			return err
		}
	}

	if t.IterationStatementDo != nil {
		if err := h.Handle(t.IterationStatementDo); err != nil {
			return err
		}
	}

	if t.IterationStatementWhile != nil {
		if err := h.Handle(t.IterationStatementWhile); err != nil {
			return err
		}
	}

	if t.IterationStatementFor != nil {
		if err := h.Handle(t.IterationStatementFor); err != nil {
			return err
		}
	}

	if t.SwitchStatement != nil {
		if err := h.Handle(t.SwitchStatement); err != nil {
			return err
		}
	}

	if t.WithStatement != nil {
		if err := h.Handle(t.WithStatement); err != nil {
			return err
		}
	}

	if t.LabelledItemFunction != nil {
		if err := h.Handle(t.LabelledItemFunction); err != nil {
			return err
		}
	}

	if t.LabelledItemStatement != nil {
		if err := h.Handle(t.LabelledItemStatement); err != nil {
			return err
		}
	}

	if t.TryStatement != nil {
		if err := h.Handle(t.TryStatement); err != nil {
			return err
		}
	}

	return nil
}

func walkIfStatement(t *javascript.IfStatement, h Handler) error {
	if err := h.Handle(&t.Expression); err != nil {
		return err
	}

	if err := h.Handle(&t.Statement); err != nil {
		return err
	}

	if t.ElseStatement != nil {
		if err := h.Handle(t.ElseStatement); err != nil {
			return err
		}
	}

	return nil
}

func walkIterationStatementDo(t *javascript.IterationStatementDo, h Handler) error {
	if err := h.Handle(&t.Statement); err != nil {
		return err
	}

	return h.Handle(&t.Expression)
}

func walkIterationStatementWhile(t *javascript.IterationStatementWhile, h Handler) error {
	if err := h.Handle(&t.Expression); err != nil {
		return err
	}

	return h.Handle(&t.Statement)
}

func walkIterationStatementFor(t *javascript.IterationStatementFor, h Handler) error {
	if t.InitExpression != nil {
		if err := h.Handle(t.InitExpression); err != nil {
			return err
		}
	}

	for n := range t.InitVar {
		if err := h.Handle(&t.InitVar[n]); err != nil {
			return err
		}
	}

	if t.InitLexical != nil {
		if err := h.Handle(t.InitLexical); err != nil {
			return err
		}
	}

	if t.Conditional != nil {
		if err := h.Handle(t.Conditional); err != nil {
			return err
		}
	}

	if t.Afterthought != nil {
		if err := h.Handle(t.Afterthought); err != nil {
			return err
		}
	}

	if t.LeftHandSideExpression != nil {
		if err := h.Handle(t.LeftHandSideExpression); err != nil {
			return err
		}
	}

	if t.ForBindingPatternObject != nil {
		if err := h.Handle(t.ForBindingPatternObject); err != nil {
			return err
		}
	}

	if t.ForBindingPatternArray != nil {
		if err := h.Handle(t.ForBindingPatternArray); err != nil {
			return err
		}
	}

	if t.In != nil {
		if err := h.Handle(t.In); err != nil {
			return err
		}
	}

	if t.Of != nil {
		if err := h.Handle(t.Of); err != nil {
			return err
		}
	}

	return h.Handle(&t.Statement)
}

func walkSwitchStatement(t *javascript.SwitchStatement, h Handler) error {
	if err := h.Handle(&t.Expression); err != nil {
		return err
	}

	for n := range t.CaseClauses {
		if err := h.Handle(&t.CaseClauses[n]); err != nil {
			return err
		}
	}

	for n := range t.DefaultClause {
		if err := h.Handle(&t.DefaultClause[n]); err != nil {
			return err
		}
	}

	for n := range t.PostDefaultCaseClauses {
		if err := h.Handle(&t.PostDefaultCaseClauses[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkCaseClause(t *javascript.CaseClause, h Handler) error {
	if err := h.Handle(&t.Expression); err != nil {
		return err
	}

	for n := range t.StatementList {
		if err := h.Handle(&t.StatementList[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkWithStatement(t *javascript.WithStatement, h Handler) error {
	if err := h.Handle(&t.Expression); err != nil {
		return err
	}

	return h.Handle(&t.Statement)
}

func walkTryStatement(t *javascript.TryStatement, h Handler) error {
	if err := h.Handle(&t.TryBlock); err != nil {
		return err
	}

	if t.CatchParameterObjectBindingPattern != nil {
		if err := h.Handle(t.CatchParameterObjectBindingPattern); err != nil {
			return err
		}
	}

	if t.CatchParameterArrayBindingPattern != nil {
		if err := h.Handle(t.CatchParameterArrayBindingPattern); err != nil {
			return err
		}
	}

	if t.CatchBlock != nil {
		if err := h.Handle(t.CatchBlock); err != nil {
			return err
		}
	}

	if t.FinallyBlock != nil {
		if err := h.Handle(t.FinallyBlock); err != nil {
			return err
		}
	}

	return nil
}

func walkVariableStatement(t *javascript.VariableStatement, h Handler) error {
	for n := range t.VariableDeclarationList {
		if err := h.Handle(&t.VariableDeclarationList[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkJSXElement(t *javascript.JSXElement, h Handler) error {
	if err := h.Handle(&t.ElementName); err != nil {
		return err
	}

	for n := range t.Attributes {
		if err := h.Handle(&t.Attributes[n]); err != nil {
			return err
		}
	}

	for n := range t.Children {
		if err := h.Handle(&t.Children[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkJSXFragment(t *javascript.JSXFragment, h Handler) error {
	for n := range t.Children {
		if err := h.Handle(&t.Children[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkJSXElementName(_ *javascript.JSXElementName, _ Handler) error {
	return nil
}

func walkJSXAttribute(t *javascript.JSXAttribute, h Handler) error {
	if t.AssignmentExpression != nil {
		return h.Handle(t.AssignmentExpression)
	} else if t.JSXElement != nil {
		return h.Handle(t.JSXElement)
	} else if t.JSXFragment != nil {
		return h.Handle(t.JSXFragment)
	}

	return nil
}

func walkJSXChild(t *javascript.JSXChild, h Handler) error {
	if t.JSXChildExpression != nil {
		return h.Handle(t.JSXChildExpression)
	} else if t.JSXElement != nil {
		return h.Handle(t.JSXElement)
	} else if t.JSXFragment != nil {
		return h.Handle(t.JSXFragment)
	}

	return nil
}
