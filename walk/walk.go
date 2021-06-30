// Package walk provides a javascript type walker
package walk

import "vimagination.zapto.org/javascript"

// Handler is used to process javascript types
type Handler interface {
	Handle(javascript.Type) error
}

// HandlerFunc wraps a func to implement Handler interface
type HandlerFunc func(javascript.Type) error

// Handle implements the Handler interface
func (h HandlerFunc) Handle(t javascript.Type) error {
	return h(t)
}

// Walk calls the Handle function on the given interface for each non-nil, non-Token field of the given javascript type
func Walk(t javascript.Type, fn Handler) error {
	switch t := t.(type) {
	case javascript.ClassDeclaration:
		return walkClassDeclaration(&t, fn)
	case *javascript.ClassDeclaration:
		return walkClassDeclaration(t, fn)
	case javascript.MethodDefinition:
		return walkMethodDefinition(&t, fn)
	case *javascript.MethodDefinition:
		return walkMethodDefinition(t, fn)
	case javascript.PropertyName:
		return walkPropertyName(&t, fn)
	case *javascript.PropertyName:
		return walkPropertyName(t, fn)
	case javascript.ConditionalExpression:
		return walkConditionalExpression(&t, fn)
	case *javascript.ConditionalExpression:
		return walkConditionalExpression(t, fn)
	case javascript.CoalesceExpression:
		return walkCoalesceExpression(&t, fn)
	case *javascript.CoalesceExpression:
		return walkCoalesceExpression(t, fn)
	case javascript.LogicalORExpression:
		return walkLogicalORExpression(&t, fn)
	case *javascript.LogicalORExpression:
		return walkLogicalORExpression(t, fn)
	case javascript.LogicalANDExpression:
		return walkLogicalANDExpression(&t, fn)
	case *javascript.LogicalANDExpression:
		return walkLogicalANDExpression(t, fn)
	case javascript.BitwiseORExpression:
		return walkBitwiseORExpression(&t, fn)
	case *javascript.BitwiseORExpression:
		return walkBitwiseORExpression(t, fn)
	case javascript.BitwiseXORExpression:
		return walkBitwiseXORExpression(&t, fn)
	case *javascript.BitwiseXORExpression:
		return walkBitwiseXORExpression(t, fn)
	case javascript.BitwiseANDExpression:
		return walkBitwiseANDExpression(&t, fn)
	case *javascript.BitwiseANDExpression:
		return walkBitwiseANDExpression(t, fn)
	case javascript.EqualityExpression:
		return walkEqualityExpression(&t, fn)
	case *javascript.EqualityExpression:
		return walkEqualityExpression(t, fn)
	case javascript.RelationalExpression:
		return walkRelationalExpression(&t, fn)
	case *javascript.RelationalExpression:
		return walkRelationalExpression(t, fn)
	case javascript.ShiftExpression:
		return walkShiftExpression(&t, fn)
	case *javascript.ShiftExpression:
		return walkShiftExpression(t, fn)
	case javascript.AdditiveExpression:
		return walkAdditiveExpression(&t, fn)
	case *javascript.AdditiveExpression:
		return walkAdditiveExpression(t, fn)
	case javascript.MultiplicativeExpression:
		return walkMultiplicativeExpression(&t, fn)
	case *javascript.MultiplicativeExpression:
		return walkMultiplicativeExpression(t, fn)
	case javascript.ExponentiationExpression:
		return walkExponentiationExpression(&t, fn)
	case *javascript.ExponentiationExpression:
		return walkExponentiationExpression(t, fn)
	case javascript.UnaryExpression:
		return walkUnaryExpression(&t, fn)
	case *javascript.UnaryExpression:
		return walkUnaryExpression(t, fn)
	case javascript.UpdateExpression:
		return walkUpdateExpression(&t, fn)
	case *javascript.UpdateExpression:
		return walkUpdateExpression(t, fn)
	case javascript.AssignmentExpression:
		return walkAssignmentExpression(&t, fn)
	case *javascript.AssignmentExpression:
		return walkAssignmentExpression(t, fn)
	case javascript.LeftHandSideExpression:
		return walkLeftHandSideExpression(&t, fn)
	case *javascript.LeftHandSideExpression:
		return walkLeftHandSideExpression(t, fn)
	case javascript.OptionalExpression:
		return walkOptionalExpression(&t, fn)
	case *javascript.OptionalExpression:
		return walkOptionalExpression(t, fn)
	case javascript.OptionalChain:
		return walkOptionalChain(&t, fn)
	case *javascript.OptionalChain:
		return walkOptionalChain(t, fn)
	case javascript.Expression:
		return walkExpression(&t, fn)
	case *javascript.Expression:
		return walkExpression(t, fn)
	case javascript.NewExpression:
		return walkNewExpression(&t, fn)
	case *javascript.NewExpression:
		return walkNewExpression(t, fn)
	case javascript.MemberExpression:
		return walkMemberExpression(&t, fn)
	case *javascript.MemberExpression:
		return walkMemberExpression(t, fn)
	case javascript.PrimaryExpression:
		return walkPrimaryExpression(&t, fn)
	case *javascript.PrimaryExpression:
		return walkPrimaryExpression(t, fn)
	case javascript.CoverParenthesizedExpressionAndArrowParameterList:
		return walkCoverParenthesizedExpressionAndArrowParameterList(&t, fn)
	case *javascript.CoverParenthesizedExpressionAndArrowParameterList:
		return walkCoverParenthesizedExpressionAndArrowParameterList(t, fn)
	case javascript.Arguments:
		return walkArguments(&t, fn)
	case *javascript.Arguments:
		return walkArguments(t, fn)
	case javascript.CallExpression:
		return walkCallExpression(&t, fn)
	case *javascript.CallExpression:
		return walkCallExpression(t, fn)
	case javascript.FunctionDeclaration:
		return walkFunctionDeclaration(&t, fn)
	case *javascript.FunctionDeclaration:
		return walkFunctionDeclaration(t, fn)
	case javascript.FormalParameters:
		return walkFormalParameters(&t, fn)
	case *javascript.FormalParameters:
		return walkFormalParameters(t, fn)
	case javascript.BindingElement:
		return walkBindingElement(&t, fn)
	case *javascript.BindingElement:
		return walkBindingElement(t, fn)
	case javascript.Script:
		return walkScript(&t, fn)
	case *javascript.Script:
		return walkScript(t, fn)
	case javascript.Declaration:
		return walkDeclaration(&t, fn)
	case *javascript.Declaration:
		return walkDeclaration(t, fn)
	case javascript.LexicalDeclaration:
		return walkLexicalDeclaration(&t, fn)
	case *javascript.LexicalDeclaration:
		return walkLexicalDeclaration(t, fn)
	case javascript.LexicalBinding:
		return walkLexicalBinding(&t, fn)
	case *javascript.LexicalBinding:
		return walkLexicalBinding(t, fn)
	case javascript.ArrayBindingPattern:
		return walkArrayBindingPattern(&t, fn)
	case *javascript.ArrayBindingPattern:
		return walkArrayBindingPattern(t, fn)
	case javascript.ObjectBindingPattern:
		return walkObjectBindingPattern(&t, fn)
	case *javascript.ObjectBindingPattern:
		return walkObjectBindingPattern(t, fn)
	case javascript.BindingProperty:
		return walkBindingProperty(&t, fn)
	case *javascript.BindingProperty:
		return walkBindingProperty(t, fn)
	case javascript.VariableDeclaration:
		return walkVariableDeclaration(&t, fn)
	case *javascript.VariableDeclaration:
		return walkVariableDeclaration(t, fn)
	case javascript.ArrayLiteral:
		return walkArrayLiteral(&t, fn)
	case *javascript.ArrayLiteral:
		return walkArrayLiteral(t, fn)
	case javascript.ObjectLiteral:
		return walkObjectLiteral(&t, fn)
	case *javascript.ObjectLiteral:
		return walkObjectLiteral(t, fn)
	case javascript.PropertyDefinition:
		return walkPropertyDefinition(&t, fn)
	case *javascript.PropertyDefinition:
		return walkPropertyDefinition(t, fn)
	case javascript.TemplateLiteral:
		return walkTemplateLiteral(&t, fn)
	case *javascript.TemplateLiteral:
		return walkTemplateLiteral(t, fn)
	case javascript.ArrowFunction:
		return walkArrowFunction(&t, fn)
	case *javascript.ArrowFunction:
		return walkArrowFunction(t, fn)
	case javascript.Module:
		return walkModule(&t, fn)
	case *javascript.Module:
		return walkModule(t, fn)
	case javascript.ModuleItem:
		return walkModuleItem(&t, fn)
	case *javascript.ModuleItem:
		return walkModuleItem(t, fn)
	case javascript.ImportDeclaration:
		return walkImportDeclaration(&t, fn)
	case *javascript.ImportDeclaration:
		return walkImportDeclaration(t, fn)
	case javascript.ImportClause:
		return walkImportClause(&t, fn)
	case *javascript.ImportClause:
		return walkImportClause(t, fn)
	case javascript.FromClause:
		return walkFromClause(&t, fn)
	case *javascript.FromClause:
		return walkFromClause(t, fn)
	case javascript.NamedImports:
		return walkNamedImports(&t, fn)
	case *javascript.NamedImports:
		return walkNamedImports(t, fn)
	case javascript.ImportSpecifier:
		return walkImportSpecifier(&t, fn)
	case *javascript.ImportSpecifier:
		return walkImportSpecifier(t, fn)
	case javascript.ExportDeclaration:
		return walkExportDeclaration(&t, fn)
	case *javascript.ExportDeclaration:
		return walkExportDeclaration(t, fn)
	case javascript.ExportClause:
		return walkExportClause(&t, fn)
	case *javascript.ExportClause:
		return walkExportClause(t, fn)
	case javascript.ExportSpecifier:
		return walkExportSpecifier(&t, fn)
	case *javascript.ExportSpecifier:
		return walkExportSpecifier(t, fn)
	case javascript.Block:
		return walkBlock(&t, fn)
	case *javascript.Block:
		return walkBlock(t, fn)
	case javascript.StatementListItem:
		return walkStatementListItem(&t, fn)
	case *javascript.StatementListItem:
		return walkStatementListItem(t, fn)
	case javascript.Statement:
		return walkStatement(&t, fn)
	case *javascript.Statement:
		return walkStatement(t, fn)
	case javascript.IfStatement:
		return walkIfStatement(&t, fn)
	case *javascript.IfStatement:
		return walkIfStatement(t, fn)
	case javascript.IterationStatementDo:
		return walkIterationStatementDo(&t, fn)
	case *javascript.IterationStatementDo:
		return walkIterationStatementDo(t, fn)
	case javascript.IterationStatementWhile:
		return walkIterationStatementWhile(&t, fn)
	case *javascript.IterationStatementWhile:
		return walkIterationStatementWhile(t, fn)
	case javascript.IterationStatementFor:
		return walkIterationStatementFor(&t, fn)
	case *javascript.IterationStatementFor:
		return walkIterationStatementFor(t, fn)
	case javascript.SwitchStatement:
		return walkSwitchStatement(&t, fn)
	case *javascript.SwitchStatement:
		return walkSwitchStatement(t, fn)
	case javascript.CaseClause:
		return walkCaseClause(&t, fn)
	case *javascript.CaseClause:
		return walkCaseClause(t, fn)
	case javascript.WithStatement:
		return walkWithStatement(&t, fn)
	case *javascript.WithStatement:
		return walkWithStatement(t, fn)
	case javascript.TryStatement:
		return walkTryStatement(&t, fn)
	case *javascript.TryStatement:
		return walkTryStatement(t, fn)
	case javascript.VariableStatement:
		return walkVariableStatement(&t, fn)
	case *javascript.VariableStatement:
		return walkVariableStatement(t, fn)
	}
	return nil
}

func walkClassDeclaration(t *javascript.ClassDeclaration, fn Handler) error {
	if t.ClassHeritage != nil {
		if err := fn.Handle(t.ClassHeritage); err != nil {
			return err
		}
	}
	for n := range t.ClassBody {
		if err := fn.Handle(&t.ClassBody[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkMethodDefinition(t *javascript.MethodDefinition, fn Handler) error {
	if err := fn.Handle(&t.PropertyName); err != nil {
		return err
	}
	if err := fn.Handle(&t.Params); err != nil {
		return err
	}
	return fn.Handle(t.FunctionBody)
}

func walkPropertyName(t *javascript.PropertyName, fn Handler) error {
	if t.ComputedPropertyName != nil {
		return fn.Handle(t.ComputedPropertyName)
	}
	return nil
}

func walkConditionalExpression(t *javascript.ConditionalExpression, fn Handler) error {
	if t.LogicalORExpression != nil {
		if err := fn.Handle(t.LogicalORExpression); err != nil {
			return err
		}
	}
	if t.CoalesceExpression != nil {
		if err := fn.Handle(t.CoalesceExpression); err != nil {
			return nil
		}
	}
	if t.True != nil {
		if err := fn.Handle(t.True); err != nil {
			return err
		}
	}
	if t.False != nil {
		if err := fn.Handle(t.False); err != nil {
			return err
		}
	}
	return nil
}

func walkCoalesceExpression(t *javascript.CoalesceExpression, fn Handler) error {
	if t.CoalesceExpressionHead != nil {
		if err := fn.Handle(t.CoalesceExpressionHead); err != nil {
			return err
		}
	}
	return fn.Handle(&t.BitwiseORExpression)
}

func walkLogicalORExpression(t *javascript.LogicalORExpression, fn Handler) error {
	if t.LogicalORExpression != nil {
		if err := fn.Handle(t.LogicalORExpression); err != nil {
			return err
		}
	}
	return fn.Handle(&t.LogicalANDExpression)
}

func walkLogicalANDExpression(t *javascript.LogicalANDExpression, fn Handler) error {
	if t.LogicalANDExpression != nil {
		if err := fn.Handle(t.LogicalANDExpression); err != nil {
			return err
		}
	}
	return fn.Handle(&t.BitwiseORExpression)
}

func walkBitwiseORExpression(t *javascript.BitwiseORExpression, fn Handler) error {
	if t.BitwiseORExpression != nil {
		if err := fn.Handle(t.BitwiseORExpression); err != nil {
			return err
		}
	}
	return fn.Handle(&t.BitwiseXORExpression)
}

func walkBitwiseXORExpression(t *javascript.BitwiseXORExpression, fn Handler) error {
	if t.BitwiseXORExpression != nil {
		if err := fn.Handle(t.BitwiseXORExpression); err != nil {
			return err
		}
	}
	return fn.Handle(&t.BitwiseANDExpression)
}

func walkBitwiseANDExpression(t *javascript.BitwiseANDExpression, fn Handler) error {
	if t.BitwiseANDExpression != nil {
		if err := fn.Handle(t.BitwiseANDExpression); err != nil {
			return err
		}
	}
	return fn.Handle(&t.EqualityExpression)
}

func walkEqualityExpression(t *javascript.EqualityExpression, fn Handler) error {
	if t.EqualityExpression != nil {
		if err := fn.Handle(t.EqualityExpression); err != nil {
			return err
		}
	}
	return fn.Handle(&t.RelationalExpression)
}

func walkRelationalExpression(t *javascript.RelationalExpression, fn Handler) error {
	if t.RelationalExpression != nil {
		if err := fn.Handle(t.RelationalExpression); err != nil {
			return err
		}
	}
	return fn.Handle(&t.ShiftExpression)
}

func walkShiftExpression(t *javascript.ShiftExpression, fn Handler) error {
	if t.ShiftExpression != nil {
		if err := fn.Handle(t.ShiftExpression); err != nil {
			return err
		}
	}
	return fn.Handle(&t.AdditiveExpression)
}

func walkAdditiveExpression(t *javascript.AdditiveExpression, fn Handler) error {
	if t.AdditiveExpression != nil {
		if err := fn.Handle(t.AdditiveExpression); err != nil {
			return err
		}
	}
	return fn.Handle(&t.MultiplicativeExpression)
}

func walkMultiplicativeExpression(t *javascript.MultiplicativeExpression, fn Handler) error {
	if t.MultiplicativeExpression != nil {
		if err := fn.Handle(t.MultiplicativeExpression); err != nil {
			return err
		}
	}
	return fn.Handle(&t.ExponentiationExpression)
}

func walkExponentiationExpression(t *javascript.ExponentiationExpression, fn Handler) error {
	if t.ExponentiationExpression != nil {
		if err := fn.Handle(t.ExponentiationExpression); err != nil {
			return err
		}
	}
	return fn.Handle(&t.UnaryExpression)
}

func walkUnaryExpression(t *javascript.UnaryExpression, fn Handler) error {
	return fn.Handle(&t.UpdateExpression)
}

func walkUpdateExpression(t *javascript.UpdateExpression, fn Handler) error {
	if t.LeftHandSideExpression != nil {
		if err := fn.Handle(t.LeftHandSideExpression); err != nil {
			return err
		}
	}
	if t.UnaryExpression != nil {
		if err := fn.Handle(t.UnaryExpression); err != nil {
			return err
		}
	}
	return nil
}

func walkAssignmentExpression(t *javascript.AssignmentExpression, fn Handler) error {
	if t.ConditionalExpression != nil {
		if err := fn.Handle(t.ConditionalExpression); err != nil {
			return err
		}
	}
	if t.ArrowFunction != nil {
		if err := fn.Handle(t.ArrowFunction); err != nil {
			return err
		}
	}
	if t.LeftHandSideExpression != nil {
		if err := fn.Handle(t.LeftHandSideExpression); err != nil {
			return err
		}
	}
	if t.AssignmentExpression != nil {
		if err := fn.Handle(t.AssignmentExpression); err != nil {
			return err
		}
	}
	return nil
}

func walkLeftHandSideExpression(t *javascript.LeftHandSideExpression, fn Handler) error {
	if t.NewExpression != nil {
		if err := fn.Handle(t.NewExpression); err != nil {
			return err
		}
	}
	if t.CallExpression != nil {
		if err := fn.Handle(t.CallExpression); err != nil {
			return err
		}
	}
	if t.OptionalExpression != nil {
		if err := fn.Handle(t.OptionalExpression); err != nil {
			return err
		}
	}
	return nil
}

func walkOptionalExpression(t *javascript.OptionalExpression, fn Handler) error {
	if t.MemberExpression != nil {
		if err := fn.Handle(t.MemberExpression); err != nil {
			return err
		}
	}
	if t.CallExpression != nil {
		if err := fn.Handle(t.CallExpression); err != nil {
			return err
		}
	}
	if t.OptionalExpression != nil {
		if err := fn.Handle(t.OptionalExpression); err != nil {
			return err
		}
	}
	return fn.Handle(&t.OptionalChain)
}

func walkOptionalChain(t *javascript.OptionalChain, fn Handler) error {
	if t.OptionalChain != nil {
		if err := fn.Handle(t.OptionalChain); err != nil {
			return err
		}
	}
	if t.Arguments != nil {
		if err := fn.Handle(t.Arguments); err != nil {
			return err
		}
	}
	if t.Expression != nil {
		if err := fn.Handle(t.Expression); err != nil {
			return err
		}
	}
	if t.TemplateLiteral != nil {
		if err := fn.Handle(t.TemplateLiteral); err != nil {
			return err
		}
	}
	return nil
}

func walkExpression(t *javascript.Expression, fn Handler) error {
	for n := range t.Expressions {
		if err := fn.Handle(&t.Expressions[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkNewExpression(t *javascript.NewExpression, fn Handler) error {
	return fn.Handle(&t.MemberExpression)
}

func walkMemberExpression(t *javascript.MemberExpression, fn Handler) error {
	if t.MemberExpression != nil {
		if err := fn.Handle(t.MemberExpression); err != nil {
			return err
		}
	}
	if t.PrimaryExpression != nil {
		if err := fn.Handle(t.PrimaryExpression); err != nil {
			return err
		}
	}
	if t.Expression != nil {
		if err := fn.Handle(t.Expression); err != nil {
			return err
		}
	}
	if t.TemplateLiteral != nil {
		if err := fn.Handle(t.TemplateLiteral); err != nil {
			return err
		}
	}
	if t.Arguments != nil {
		if err := fn.Handle(t.Arguments); err != nil {
			return err
		}
	}
	return nil
}

func walkPrimaryExpression(t *javascript.PrimaryExpression, fn Handler) error {
	if t.ArrayLiteral != nil {
		if err := fn.Handle(t.ArrayLiteral); err != nil {
			return err
		}
	}
	if t.ObjectLiteral != nil {
		if err := fn.Handle(t.ObjectLiteral); err != nil {
			return err
		}
	}
	if t.FunctionExpression != nil {
		if err := fn.Handle(t.FunctionExpression); err != nil {
			return err
		}
	}
	if t.ClassExpression != nil {
		if err := fn.Handle(t.ClassExpression); err != nil {
			return err
		}
	}
	if t.TemplateLiteral != nil {
		if err := fn.Handle(t.TemplateLiteral); err != nil {
			return err
		}
	}
	if t.CoverParenthesizedExpressionAndArrowParameterList != nil {
		if err := fn.Handle(t.CoverParenthesizedExpressionAndArrowParameterList); err != nil {
			return err
		}
	}
	return nil
}

func walkCoverParenthesizedExpressionAndArrowParameterList(t *javascript.CoverParenthesizedExpressionAndArrowParameterList, fn Handler) error {
	for n := range t.Expressions {
		if err := fn.Handle(&t.Expressions[n]); err != nil {
			return err
		}
	}
	if t.ArrayBindingPattern != nil {
		if err := fn.Handle(t.ArrayBindingPattern); err != nil {
			return err
		}
	}
	if t.ObjectBindingPattern != nil {
		if err := fn.Handle(t.ObjectBindingPattern); err != nil {
			return err
		}
	}
	return nil
}

func walkArguments(t *javascript.Arguments, fn Handler) error {
	for n := range t.ArgumentList {
		if err := fn.Handle(&t.ArgumentList[n]); err != nil {
			return err
		}
	}
	if t.SpreadArgument != nil {
		if err := fn.Handle(t.SpreadArgument); err != nil {
			return err
		}
	}
	return nil
}

func walkCallExpression(t *javascript.CallExpression, fn Handler) error {
	if t.MemberExpression != nil {
		if err := fn.Handle(t.MemberExpression); err != nil {
			return err
		}
	}
	if t.ImportCall != nil {
		if err := fn.Handle(t.ImportCall); err != nil {
			return err
		}
	}
	if t.CallExpression != nil {
		if err := fn.Handle(t.CallExpression); err != nil {
			return err
		}
	}
	if t.Arguments != nil {
		if err := fn.Handle(t.Arguments); err != nil {
			return err
		}
	}
	if t.Expression != nil {
		if err := fn.Handle(t.Expression); err != nil {
			return err
		}
	}
	if t.TemplateLiteral != nil {
		if err := fn.Handle(t.TemplateLiteral); err != nil {
			return err
		}
	}
	return nil
}

func walkFunctionDeclaration(t *javascript.FunctionDeclaration, fn Handler) error {
	if err := fn.Handle(&t.FormalParameters); err != nil {
		return err
	}
	return fn.Handle(&t.FunctionBody)
}

func walkFormalParameters(t *javascript.FormalParameters, fn Handler) error {
	for n := range t.FormalParameterList {
		if err := fn.Handle(&t.FormalParameterList[n]); err != nil {
			return err
		}
	}
	if t.ArrayBindingPattern != nil {
		if err := fn.Handle(t.ArrayBindingPattern); err != nil {
			return err
		}
	}
	if t.ObjectBindingPattern != nil {
		if err := fn.Handle(t.ObjectBindingPattern); err != nil {
			return err
		}
	}
	return nil
}

func walkBindingElement(t *javascript.BindingElement, fn Handler) error {
	if t.ArrayBindingPattern != nil {
		if err := fn.Handle(t.ArrayBindingPattern); err != nil {
			return err
		}
	}
	if t.ObjectBindingPattern != nil {
		if err := fn.Handle(t.ObjectBindingPattern); err != nil {
			return err
		}
	}
	if t.Initializer != nil {
		if err := fn.Handle(t.Initializer); err != nil {
			return err
		}
	}
	return nil
}

func walkScript(t *javascript.Script, fn Handler) error {
	for n := range t.StatementList {
		if err := fn.Handle(&t.StatementList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkDeclaration(t *javascript.Declaration, fn Handler) error {
	if t.ClassDeclaration != nil {
		if err := fn.Handle(t.ClassDeclaration); err != nil {
			return err
		}
	}
	if t.FunctionDeclaration != nil {
		if err := fn.Handle(t.FunctionDeclaration); err != nil {
			return err
		}
	}
	if t.LexicalDeclaration != nil {
		if err := fn.Handle(t.LexicalDeclaration); err != nil {
			return err
		}
	}
	return nil
}

func walkLexicalDeclaration(t *javascript.LexicalDeclaration, fn Handler) error {
	for n := range t.BindingList {
		if err := fn.Handle(&t.BindingList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkLexicalBinding(t *javascript.LexicalBinding, fn Handler) error {
	if t.ArrayBindingPattern != nil {
		if err := fn.Handle(t.ArrayBindingPattern); err != nil {
			return err
		}
	}
	if t.ObjectBindingPattern != nil {
		if err := fn.Handle(t.ObjectBindingPattern); err != nil {
			return err
		}
	}
	if t.Initializer != nil {
		if err := fn.Handle(t.Initializer); err != nil {
			return err
		}
	}
	return nil
}

func walkArrayBindingPattern(t *javascript.ArrayBindingPattern, fn Handler) error {
	for n := range t.BindingElementList {
		if err := fn.Handle(&t.BindingElementList[n]); err != nil {
			return err
		}
	}
	if t.BindingRestElement != nil {
		if err := fn.Handle(t.BindingRestElement); err != nil {
			return err
		}
	}
	return nil
}

func walkObjectBindingPattern(t *javascript.ObjectBindingPattern, fn Handler) error {
	for n := range t.BindingPropertyList {
		if err := fn.Handle(&t.BindingPropertyList[n]); err != nil {
			return err
		}
	}
	if t.BindingRestProperty != nil {
		if err := fn.Handle(t.BindingRestProperty); err != nil {
			return err
		}
	}
	return nil
}

func walkBindingProperty(t *javascript.BindingProperty, fn Handler) error {
	if err := fn.Handle(&t.PropertyName); err != nil {
		return err
	}
	return fn.Handle(&t.BindingElement)
}

func walkVariableDeclaration(t *javascript.VariableDeclaration, fn Handler) error {
	return walkLexicalBinding((*javascript.LexicalBinding)(t), fn)
}

func walkArrayLiteral(t *javascript.ArrayLiteral, fn Handler) error {
	for n := range t.ElementList {
		if err := fn.Handle(&t.ElementList[n]); err != nil {
			return err
		}
	}
	if t.SpreadElement != nil {
		if err := fn.Handle(t.SpreadElement); err != nil {
			return err
		}
	}
	return nil
}

func walkObjectLiteral(t *javascript.ObjectLiteral, fn Handler) error {
	for n := range t.PropertyDefinitionList {
		if err := fn.Handle(&t.PropertyDefinitionList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkPropertyDefinition(t *javascript.PropertyDefinition, fn Handler) error {
	if t.PropertyName != nil {
		if err := fn.Handle(t.PropertyName); err != nil {
			return err
		}
	}
	if t.AssignmentExpression != nil {
		if err := fn.Handle(t.AssignmentExpression); err != nil {
			return err
		}
	}
	if t.MethodDefinition != nil {
		if err := fn.Handle(t.MethodDefinition); err != nil {
			return err
		}
	}
	return nil
}

func walkTemplateLiteral(t *javascript.TemplateLiteral, fn Handler) error {
	for n := range t.Expressions {
		if err := fn.Handle(&t.Expressions[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkArrowFunction(t *javascript.ArrowFunction, fn Handler) error {
	if t.CoverParenthesizedExpressionAndArrowParameterList != nil {
		if err := fn.Handle(t.CoverParenthesizedExpressionAndArrowParameterList); err != nil {
			return err
		}
	}
	if t.FormalParameters != nil {
		if err := fn.Handle(t.FormalParameters); err != nil {
			return err
		}
	}
	if t.AssignmentExpression != nil {
		if err := fn.Handle(t.AssignmentExpression); err != nil {
			return err
		}
	}
	if t.FunctionBody != nil {
		if err := fn.Handle(t.FunctionBody); err != nil {
			return err
		}
	}
	return nil
}

func walkModule(t *javascript.Module, fn Handler) error {
	for n := range t.ModuleListItems {
		if err := fn.Handle(&t.ModuleListItems[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkModuleItem(t *javascript.ModuleItem, fn Handler) error {
	if t.ImportDeclaration != nil {
		if err := fn.Handle(t.ImportDeclaration); err != nil {
			return err
		}
	}
	if t.StatementListItem != nil {
		if err := fn.Handle(t.StatementListItem); err != nil {
			return err
		}
	}
	if t.ExportDeclaration != nil {
		if err := fn.Handle(t.ExportDeclaration); err != nil {
			return err
		}
	}
	return nil
}

func walkImportDeclaration(t *javascript.ImportDeclaration, fn Handler) error {
	if t.ImportClause != nil {
		if err := fn.Handle(t.ImportClause); err != nil {
			return err
		}
	}
	return fn.Handle(&t.FromClause)
}

func walkImportClause(t *javascript.ImportClause, fn Handler) error {
	if t.NamedImports != nil {
		if err := fn.Handle(t.NamedImports); err != nil {
			return err
		}
	}
	return nil
}

func walkFromClause(t *javascript.FromClause, fn Handler) error {
	return nil
}

func walkNamedImports(t *javascript.NamedImports, fn Handler) error {
	for n := range t.ImportList {
		if err := fn.Handle(&t.ImportList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkImportSpecifier(t *javascript.ImportSpecifier, fn Handler) error {
	return nil
}

func walkExportDeclaration(t *javascript.ExportDeclaration, fn Handler) error {
	if t.ExportClause != nil {
		if err := fn.Handle(t.ExportClause); err != nil {
			return err
		}
	}
	if t.FromClause != nil {
		if err := fn.Handle(t.FromClause); err != nil {
			return err
		}
	}
	if t.VariableStatement != nil {
		if err := fn.Handle(t.VariableStatement); err != nil {
			return err
		}
	}
	if t.Declaration != nil {
		if err := fn.Handle(t.Declaration); err != nil {
			return err
		}
	}
	if t.DefaultFunction != nil {
		if err := fn.Handle(t.DefaultFunction); err != nil {
			return err
		}
	}
	if t.DefaultClass != nil {
		if err := fn.Handle(t.DefaultClass); err != nil {
			return err
		}
	}
	if t.DefaultAssignmentExpression != nil {
		if err := fn.Handle(t.DefaultAssignmentExpression); err != nil {
			return err
		}
	}
	return nil
}

func walkExportClause(t *javascript.ExportClause, fn Handler) error {
	for n := range t.ExportList {
		if err := fn.Handle(&t.ExportList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkExportSpecifier(t *javascript.ExportSpecifier, fn Handler) error {
	return nil
}

func walkBlock(t *javascript.Block, fn Handler) error {
	for n := range t.StatementList {
		if err := fn.Handle(&t.StatementList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkStatementListItem(t *javascript.StatementListItem, fn Handler) error {
	if t.Statement != nil {
		if err := fn.Handle(t.Statement); err != nil {
			return err
		}
	}
	if t.Declaration != nil {
		if err := fn.Handle(t.Declaration); err != nil {
			return err
		}
	}
	return nil
}

func walkStatement(t *javascript.Statement, fn Handler) error {

	if t.BlockStatement != nil {
		if err := fn.Handle(t.BlockStatement); err != nil {
			return err
		}
	}
	if t.VariableStatement != nil {
		if err := fn.Handle(t.VariableStatement); err != nil {
			return err
		}
	}
	if t.ExpressionStatement != nil {
		if err := fn.Handle(t.ExpressionStatement); err != nil {
			return err
		}
	}
	if t.IfStatement != nil {
		if err := fn.Handle(t.IfStatement); err != nil {
			return err
		}
	}
	if t.IterationStatementDo != nil {
		if err := fn.Handle(t.IterationStatementDo); err != nil {
			return err
		}
	}
	if t.IterationStatementWhile != nil {
		if err := fn.Handle(t.IterationStatementWhile); err != nil {
			return err
		}
	}
	if t.IterationStatementFor != nil {
		if err := fn.Handle(t.IterationStatementFor); err != nil {
			return err
		}
	}
	if t.SwitchStatement != nil {
		if err := fn.Handle(t.SwitchStatement); err != nil {
			return err
		}
	}
	if t.WithStatement != nil {
		if err := fn.Handle(t.WithStatement); err != nil {
			return err
		}
	}
	if t.LabelledItemFunction != nil {
		if err := fn.Handle(t.LabelledItemFunction); err != nil {
			return err
		}
	}
	if t.LabelledItemStatement != nil {
		if err := fn.Handle(t.LabelledItemStatement); err != nil {
			return err
		}
	}
	if t.TryStatement != nil {
		if err := fn.Handle(t.TryStatement); err != nil {
			return err
		}
	}
	return nil
}

func walkIfStatement(t *javascript.IfStatement, fn Handler) error {
	if err := fn.Handle(&t.Expression); err != nil {
		return err
	}
	if err := fn.Handle(&t.Statement); err != nil {
		return err
	}
	if t.ElseStatement != nil {
		if err := fn.Handle(t.ElseStatement); err != nil {
			return err
		}
	}
	return nil
}

func walkIterationStatementDo(t *javascript.IterationStatementDo, fn Handler) error {
	if err := fn.Handle(&t.Statement); err != nil {
		return err
	}
	return fn.Handle(&t.Expression)
}

func walkIterationStatementWhile(t *javascript.IterationStatementWhile, fn Handler) error {
	if err := fn.Handle(&t.Expression); err != nil {
		return err
	}
	return fn.Handle(&t.Statement)
}

func walkIterationStatementFor(t *javascript.IterationStatementFor, fn Handler) error {
	if t.InitExpression != nil {
		if err := fn.Handle(t.InitExpression); err != nil {
			return err
		}
	}
	for n := range t.InitVar {
		if err := fn.Handle(&t.InitVar[n]); err != nil {
			return err
		}
	}
	if t.InitLexical != nil {
		if err := fn.Handle(t.InitLexical); err != nil {
			return err
		}
	}
	if t.Conditional != nil {
		if err := fn.Handle(t.Conditional); err != nil {
			return err
		}
	}
	if t.Afterthought != nil {
		if err := fn.Handle(t.Afterthought); err != nil {
			return err
		}
	}
	if t.LeftHandSideExpression != nil {
		if err := fn.Handle(t.LeftHandSideExpression); err != nil {
			return err
		}
	}
	if t.ForBindingPatternObject != nil {
		if err := fn.Handle(t.ForBindingPatternObject); err != nil {
			return err
		}
	}
	if t.ForBindingPatternArray != nil {
		if err := fn.Handle(t.ForBindingPatternArray); err != nil {
			return err
		}
	}
	if t.In != nil {
		if err := fn.Handle(t.In); err != nil {
			return err
		}
	}
	if t.Of != nil {
		if err := fn.Handle(t.Of); err != nil {
			return err
		}
	}
	return fn.Handle(&t.Statement)
}

func walkSwitchStatement(t *javascript.SwitchStatement, fn Handler) error {
	if err := fn.Handle(&t.Expression); err != nil {
		return err
	}
	for n := range t.CaseClauses {
		if err := fn.Handle(&t.CaseClauses[n]); err != nil {
			return err
		}
	}
	for n := range t.DefaultClause {
		if err := fn.Handle(&t.DefaultClause[n]); err != nil {
			return err
		}
	}
	for n := range t.PostDefaultCaseClauses {
		if err := fn.Handle(&t.PostDefaultCaseClauses[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkCaseClause(t *javascript.CaseClause, fn Handler) error {
	if err := fn.Handle(&t.Expression); err != nil {
		return err
	}
	for n := range t.StatementList {
		if err := fn.Handle(&t.StatementList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkWithStatement(t *javascript.WithStatement, fn Handler) error {
	if err := fn.Handle(&t.Expression); err != nil {
		return err
	}
	return fn.Handle(&t.Statement)
}

func walkTryStatement(t *javascript.TryStatement, fn Handler) error {
	if err := fn.Handle(&t.TryBlock); err != nil {
		return err
	}
	if t.CatchParameterObjectBindingPattern != nil {
		if err := fn.Handle(t.CatchParameterObjectBindingPattern); err != nil {
			return err
		}
	}
	if t.CatchParameterArrayBindingPattern != nil {
		if err := fn.Handle(t.CatchParameterArrayBindingPattern); err != nil {
			return err
		}
	}
	if t.CatchBlock != nil {
		if err := fn.Handle(t.CatchBlock); err != nil {
			return err
		}
	}
	if t.FinallyBlock != nil {
		if err := fn.Handle(t.FinallyBlock); err != nil {
			return err
		}
	}
	return nil
}

func walkVariableStatement(t *javascript.VariableStatement, fn Handler) error {
	for n := range t.VariableDeclarationList {
		if err := fn.Handle(&t.VariableDeclarationList[n]); err != nil {
			return err
		}
	}
	return nil
}
