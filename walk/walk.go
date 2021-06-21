package walk

import "vimagination.zapto.org/javascript"

func Walk(t javascript.Type, fn func(javascript.Type) error) error {
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
	case javascript.FunctionRestParameter:
		return walkFunctionRestParameter(&t, fn)
	case *javascript.FunctionRestParameter:
		return walkFunctionRestParameter(t, fn)
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

func walkClassDeclaration(t *javascript.ClassDeclaration, fn func(javascript.Type) error) error {
	if t.ClassHeritage != nil {
		if err := fn(t.ClassHeritage); err != nil {
			return err
		}
	}
	for n := range t.ClassBody {
		if err := fn(&t.ClassBody[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkMethodDefinition(t *javascript.MethodDefinition, fn func(javascript.Type) error) error {
	if err := fn(&t.PropertyName); err != nil {
		return err
	}
	if err := fn(&t.Params); err != nil {
		return err
	}
	return fn(t.FunctionBody)
}

func walkPropertyName(t *javascript.PropertyName, fn func(javascript.Type) error) error {
	if t.ComputedPropertyName != nil {
		return fn(t.ComputedPropertyName)
	}
	return nil
}

func walkConditionalExpression(t *javascript.ConditionalExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkCoalesceExpression(t *javascript.CoalesceExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkLogicalORExpression(t *javascript.LogicalORExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkLogicalANDExpression(t *javascript.LogicalANDExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkBitwiseORExpression(t *javascript.BitwiseORExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkBitwiseXORExpression(t *javascript.BitwiseXORExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkBitwiseANDExpression(t *javascript.BitwiseANDExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkEqualityExpression(t *javascript.EqualityExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkRelationalExpression(t *javascript.RelationalExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkShiftExpression(t *javascript.ShiftExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkAdditiveExpression(t *javascript.AdditiveExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkMultiplicativeExpression(t *javascript.MultiplicativeExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkExponentiationExpression(t *javascript.ExponentiationExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkUnaryExpression(t *javascript.UnaryExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkUpdateExpression(t *javascript.UpdateExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkAssignmentExpression(t *javascript.AssignmentExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkLeftHandSideExpression(t *javascript.LeftHandSideExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkOptionalExpression(t *javascript.OptionalExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkOptionalChain(t *javascript.OptionalChain, fn func(javascript.Type) error) error {
	return nil
}

func walkExpression(t *javascript.Expression, fn func(javascript.Type) error) error {
	return nil
}

func walkNewExpression(t *javascript.NewExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkMemberExpression(t *javascript.MemberExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkPrimaryExpression(t *javascript.PrimaryExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkCoverParenthesizedExpressionAndArrowParameterList(t *javascript.CoverParenthesizedExpressionAndArrowParameterList, fn func(javascript.Type) error) error {
	return nil
}

func walkArguments(t *javascript.Arguments, fn func(javascript.Type) error) error {
	return nil
}

func walkCallExpression(t *javascript.CallExpression, fn func(javascript.Type) error) error {
	return nil
}

func walkFunctionDeclaration(t *javascript.FunctionDeclaration, fn func(javascript.Type) error) error {
	return nil
}

func walkFormalParameters(t *javascript.FormalParameters, fn func(javascript.Type) error) error {
	return nil
}

func walkBindingElement(t *javascript.BindingElement, fn func(javascript.Type) error) error {
	return nil
}

func walkFunctionRestParameter(t *javascript.FunctionRestParameter, fn func(javascript.Type) error) error {
	return nil
}

func walkScript(t *javascript.Script, fn func(javascript.Type) error) error {
	return nil
}

func walkDeclaration(t *javascript.Declaration, fn func(javascript.Type) error) error {
	return nil
}

func walkLexicalDeclaration(t *javascript.LexicalDeclaration, fn func(javascript.Type) error) error {
	return nil
}

func walkLexicalBinding(t *javascript.LexicalBinding, fn func(javascript.Type) error) error {
	return nil
}

func walkArrayBindingPattern(t *javascript.ArrayBindingPattern, fn func(javascript.Type) error) error {
	return nil
}

func walkObjectBindingPattern(t *javascript.ObjectBindingPattern, fn func(javascript.Type) error) error {
	return nil
}

func walkBindingProperty(t *javascript.BindingProperty, fn func(javascript.Type) error) error {
	return nil
}

func walkVariableDeclaration(t *javascript.VariableDeclaration, fn func(javascript.Type) error) error {
	return nil
}

func walkArrayLiteral(t *javascript.ArrayLiteral, fn func(javascript.Type) error) error {
	return nil
}

func walkObjectLiteral(t *javascript.ObjectLiteral, fn func(javascript.Type) error) error {
	return nil
}

func walkPropertyDefinition(t *javascript.PropertyDefinition, fn func(javascript.Type) error) error {
	return nil
}

func walkTemplateLiteral(t *javascript.TemplateLiteral, fn func(javascript.Type) error) error {
	return nil
}

func walkArrowFunction(t *javascript.ArrowFunction, fn func(javascript.Type) error) error {
	return nil
}

func walkModule(t *javascript.Module, fn func(javascript.Type) error) error {
	return nil
}

func walkModuleItem(t *javascript.ModuleItem, fn func(javascript.Type) error) error {
	return nil
}

func walkImportDeclaration(t *javascript.ImportDeclaration, fn func(javascript.Type) error) error {
	return nil
}

func walkImportClause(t *javascript.ImportClause, fn func(javascript.Type) error) error {
	return nil
}

func walkFromClause(t *javascript.FromClause, fn func(javascript.Type) error) error {
	return nil
}

func walkNamedImports(t *javascript.NamedImports, fn func(javascript.Type) error) error {
	return nil
}

func walkImportSpecifier(t *javascript.ImportSpecifier, fn func(javascript.Type) error) error {
	return nil
}

func walkExportDeclaration(t *javascript.ExportDeclaration, fn func(javascript.Type) error) error {
	return nil
}

func walkExportClause(t *javascript.ExportClause, fn func(javascript.Type) error) error {
	return nil
}

func walkExportSpecifier(t *javascript.ExportSpecifier, fn func(javascript.Type) error) error {
	return nil
}

func walkBlock(t *javascript.Block, fn func(javascript.Type) error) error {
	return nil
}

func walkStatementListItem(t *javascript.StatementListItem, fn func(javascript.Type) error) error {
	return nil
}

func walkStatement(t *javascript.Statement, fn func(javascript.Type) error) error {
	return nil
}

func walkIfStatement(t *javascript.IfStatement, fn func(javascript.Type) error) error {
	return nil
}

func walkIterationStatementDo(t *javascript.IterationStatementDo, fn func(javascript.Type) error) error {
	return nil
}

func walkIterationStatementWhile(t *javascript.IterationStatementWhile, fn func(javascript.Type) error) error {
	return nil
}

func walkIterationStatementFor(t *javascript.IterationStatementFor, fn func(javascript.Type) error) error {
	return nil
}

func walkSwitchStatement(t *javascript.SwitchStatement, fn func(javascript.Type) error) error {
	return nil
}

func walkCaseClause(t *javascript.CaseClause, fn func(javascript.Type) error) error {
	return nil
}

func walkWithStatement(t *javascript.WithStatement, fn func(javascript.Type) error) error {
	return nil
}

func walkTryStatement(t *javascript.TryStatement, fn func(javascript.Type) error) error {
	return nil
}

func walkVariableStatement(t *javascript.VariableStatement, fn func(javascript.Type) error) error {
	return nil
}
