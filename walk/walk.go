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
	if t.LogicalORExpression != nil {
		if err := fn(t.LogicalORExpression); err != nil {
			return err
		}
	}
	if t.CoalesceExpression != nil {
		if err := fn(t.CoalesceExpression); err != nil {
			return nil
		}
	}
	if t.True != nil {
		if err := fn(t.True); err != nil {
			return err
		}
	}
	if t.False != nil {
		if err := fn(t.False); err != nil {
			return err
		}
	}
	return nil
}

func walkCoalesceExpression(t *javascript.CoalesceExpression, fn func(javascript.Type) error) error {
	if t.CoalesceExpressionHead != nil {
		if err := fn(t.CoalesceExpressionHead); err != nil {
			return err
		}
	}
	return fn(&t.BitwiseORExpression)
}

func walkLogicalORExpression(t *javascript.LogicalORExpression, fn func(javascript.Type) error) error {
	if t.LogicalORExpression != nil {
		if err := fn(t.LogicalORExpression); err != nil {
			return err
		}
	}
	return fn(&t.LogicalANDExpression)
}

func walkLogicalANDExpression(t *javascript.LogicalANDExpression, fn func(javascript.Type) error) error {
	if t.LogicalANDExpression != nil {
		if err := fn(t.LogicalANDExpression); err != nil {
			return err
		}
	}
	return fn(&t.BitwiseORExpression)
}

func walkBitwiseORExpression(t *javascript.BitwiseORExpression, fn func(javascript.Type) error) error {
	if t.BitwiseORExpression != nil {
		if err := fn(t.BitwiseORExpression); err != nil {
			return err
		}
	}
	return fn(&t.BitwiseXORExpression)
}

func walkBitwiseXORExpression(t *javascript.BitwiseXORExpression, fn func(javascript.Type) error) error {
	if t.BitwiseXORExpression != nil {
		if err := fn(t.BitwiseXORExpression); err != nil {
			return err
		}
	}
	return fn(&t.BitwiseANDExpression)
}

func walkBitwiseANDExpression(t *javascript.BitwiseANDExpression, fn func(javascript.Type) error) error {
	if t.BitwiseANDExpression != nil {
		if err := fn(t.BitwiseANDExpression); err != nil {
			return err
		}
	}
	return fn(&t.EqualityExpression)
}

func walkEqualityExpression(t *javascript.EqualityExpression, fn func(javascript.Type) error) error {
	if t.EqualityExpression != nil {
		if err := fn(t.EqualityExpression); err != nil {
			return err
		}
	}
	return fn(&t.RelationalExpression)
}

func walkRelationalExpression(t *javascript.RelationalExpression, fn func(javascript.Type) error) error {
	if t.RelationalExpression != nil {
		if err := fn(t.RelationalExpression); err != nil {
			return err
		}
	}
	return fn(&t.ShiftExpression)
}

func walkShiftExpression(t *javascript.ShiftExpression, fn func(javascript.Type) error) error {
	if t.ShiftExpression != nil {
		if err := fn(t.ShiftExpression); err != nil {
			return err
		}
	}
	return fn(&t.AdditiveExpression)
}

func walkAdditiveExpression(t *javascript.AdditiveExpression, fn func(javascript.Type) error) error {
	if t.AdditiveExpression != nil {
		if err := fn(t.AdditiveExpression); err != nil {
			return err
		}
	}
	return fn(&t.MultiplicativeExpression)
}

func walkMultiplicativeExpression(t *javascript.MultiplicativeExpression, fn func(javascript.Type) error) error {
	if t.MultiplicativeExpression != nil {
		if err := fn(t.MultiplicativeExpression); err != nil {
			return err
		}
	}
	return fn(&t.ExponentiationExpression)
}

func walkExponentiationExpression(t *javascript.ExponentiationExpression, fn func(javascript.Type) error) error {
	if t.ExponentiationExpression != nil {
		if err := fn(t.ExponentiationExpression); err != nil {
			return err
		}
	}
	return fn(&t.UnaryExpression)
}

func walkUnaryExpression(t *javascript.UnaryExpression, fn func(javascript.Type) error) error {
	return fn(&t.UpdateExpression)
}

func walkUpdateExpression(t *javascript.UpdateExpression, fn func(javascript.Type) error) error {
	if t.LeftHandSideExpression != nil {
		if err := fn(t.LeftHandSideExpression); err != nil {
			return err
		}
	}
	if t.UnaryExpression != nil {
		if err := fn(t.UnaryExpression); err != nil {
			return err
		}
	}
	return nil
}

func walkAssignmentExpression(t *javascript.AssignmentExpression, fn func(javascript.Type) error) error {
	if t.ConditionalExpression != nil {
		if err := fn(t.ConditionalExpression); err != nil {
			return err
		}
	}
	if t.ArrowFunction != nil {
		if err := fn(t.ArrowFunction); err != nil {
			return err
		}
	}
	if t.LeftHandSideExpression != nil {
		if err := fn(t.LeftHandSideExpression); err != nil {
			return err
		}
	}
	if t.AssignmentExpression != nil {
		if err := fn(t.AssignmentExpression); err != nil {
			return err
		}
	}
	return nil
}

func walkLeftHandSideExpression(t *javascript.LeftHandSideExpression, fn func(javascript.Type) error) error {
	if t.NewExpression != nil {
		if err := fn(t.NewExpression); err != nil {
			return err
		}
	}
	if t.CallExpression != nil {
		if err := fn(t.CallExpression); err != nil {
			return err
		}
	}
	if t.OptionalExpression != nil {
		if err := fn(t.OptionalExpression); err != nil {
			return err
		}
	}
	return nil
}

func walkOptionalExpression(t *javascript.OptionalExpression, fn func(javascript.Type) error) error {
	if t.MemberExpression != nil {
		if err := fn(t.MemberExpression); err != nil {
			return err
		}
	}
	if t.CallExpression != nil {
		if err := fn(t.CallExpression); err != nil {
			return err
		}
	}
	if t.OptionalExpression != nil {
		if err := fn(t.OptionalExpression); err != nil {
			return err
		}
	}
	return fn(&t.OptionalChain)
}

func walkOptionalChain(t *javascript.OptionalChain, fn func(javascript.Type) error) error {
	if t.OptionalChain != nil {
		if err := fn(t.OptionalChain); err != nil {
			return err
		}
	}
	if t.Arguments != nil {
		if err := fn(t.Arguments); err != nil {
			return err
		}
	}
	if t.Expression != nil {
		if err := fn(t.Expression); err != nil {
			return err
		}
	}
	if t.TemplateLiteral != nil {
		if err := fn(t.TemplateLiteral); err != nil {
			return err
		}
	}
	return nil
}

func walkExpression(t *javascript.Expression, fn func(javascript.Type) error) error {
	for n := range t.Expressions {
		if err := fn(&t.Expressions[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkNewExpression(t *javascript.NewExpression, fn func(javascript.Type) error) error {
	return fn(&t.MemberExpression)
}

func walkMemberExpression(t *javascript.MemberExpression, fn func(javascript.Type) error) error {
	if t.MemberExpression != nil {
		if err := fn(t.MemberExpression); err != nil {
			return err
		}
	}
	if t.PrimaryExpression != nil {
		if err := fn(t.PrimaryExpression); err != nil {
			return err
		}
	}
	if t.Expression != nil {
		if err := fn(t.Expression); err != nil {
			return err
		}
	}
	if t.TemplateLiteral != nil {
		if err := fn(t.TemplateLiteral); err != nil {
			return err
		}
	}
	if t.Arguments != nil {
		if err := fn(t.Arguments); err != nil {
			return err
		}
	}
	return nil
}

func walkPrimaryExpression(t *javascript.PrimaryExpression, fn func(javascript.Type) error) error {
	if t.ArrayLiteral != nil {
		if err := fn(t.ArrayLiteral); err != nil {
			return err
		}
	}
	if t.ObjectLiteral != nil {
		if err := fn(t.ObjectLiteral); err != nil {
			return err
		}
	}
	if t.FunctionExpression != nil {
		if err := fn(t.FunctionExpression); err != nil {
			return err
		}
	}
	if t.ClassExpression != nil {
		if err := fn(t.ClassExpression); err != nil {
			return err
		}
	}
	if t.TemplateLiteral != nil {
		if err := fn(t.TemplateLiteral); err != nil {
			return err
		}
	}
	if t.CoverParenthesizedExpressionAndArrowParameterList != nil {
		if err := fn(t.CoverParenthesizedExpressionAndArrowParameterList); err != nil {
			return err
		}
	}
	return nil
}

func walkCoverParenthesizedExpressionAndArrowParameterList(t *javascript.CoverParenthesizedExpressionAndArrowParameterList, fn func(javascript.Type) error) error {
	for n := range t.Expressions {
		if err := fn(&t.Expressions[n]); err != nil {
			return err
		}
	}
	if t.ArrayBindingPattern != nil {
		if err := fn(t.ArrayBindingPattern); err != nil {
			return err
		}
	}
	if t.ObjectBindingPattern != nil {
		if err := fn(t.ObjectBindingPattern); err != nil {
			return err
		}
	}
	return nil
}

func walkArguments(t *javascript.Arguments, fn func(javascript.Type) error) error {
	for n := range t.ArgumentList {
		if err := fn(&t.ArgumentList[n]); err != nil {
			return err
		}
	}
	if t.SpreadArgument != nil {
		if err := fn(t.SpreadArgument); err != nil {
			return err
		}
	}
	return nil
}

func walkCallExpression(t *javascript.CallExpression, fn func(javascript.Type) error) error {
	if t.MemberExpression != nil {
		if err := fn(t.MemberExpression); err != nil {
			return err
		}
	}
	if t.ImportCall != nil {
		if err := fn(t.ImportCall); err != nil {
			return err
		}
	}
	if t.CallExpression != nil {
		if err := fn(t.CallExpression); err != nil {
			return err
		}
	}
	if t.Arguments != nil {
		if err := fn(t.Arguments); err != nil {
			return err
		}
	}
	if t.Expression != nil {
		if err := fn(t.Expression); err != nil {
			return err
		}
	}
	if t.TemplateLiteral != nil {
		if err := fn(t.TemplateLiteral); err != nil {
			return err
		}
	}
	return nil
}

func walkFunctionDeclaration(t *javascript.FunctionDeclaration, fn func(javascript.Type) error) error {
	if err := fn(&t.FormalParameters); err != nil {
		return err
	}
	return fn(&t.FunctionBody)
}

func walkFormalParameters(t *javascript.FormalParameters, fn func(javascript.Type) error) error {
	for n := range t.FormalParameterList {
		if err := fn(&t.FormalParameterList[n]); err != nil {
			return err
		}
	}
	if t.FunctionRestParameter != nil {
		if err := fn(t.FunctionRestParameter); err != nil {
			return err
		}
	}
	return nil
}

func walkBindingElement(t *javascript.BindingElement, fn func(javascript.Type) error) error {
	if t.ArrayBindingPattern != nil {
		if err := fn(t.ArrayBindingPattern); err != nil {
			return err
		}
	}
	if t.ObjectBindingPattern != nil {
		if err := fn(t.ObjectBindingPattern); err != nil {
			return err
		}
	}
	if t.Initializer != nil {
		if err := fn(t.Initializer); err != nil {
			return err
		}
	}
	return nil
}

func walkFunctionRestParameter(t *javascript.FunctionRestParameter, fn func(javascript.Type) error) error {
	if t.ArrayBindingPattern != nil {
		if err := fn(t.ArrayBindingPattern); err != nil {
			return err
		}
	}
	if t.ObjectBindingPattern != nil {
		if err := fn(t.ObjectBindingPattern); err != nil {
			return err
		}
	}
	return nil
}

func walkScript(t *javascript.Script, fn func(javascript.Type) error) error {
	for n := range t.StatementList {
		if err := fn(&t.StatementList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkDeclaration(t *javascript.Declaration, fn func(javascript.Type) error) error {
	if t.ClassDeclaration != nil {
		if err := fn(t.ClassDeclaration); err != nil {
			return err
		}
	}
	if t.FunctionDeclaration != nil {
		if err := fn(t.FunctionDeclaration); err != nil {
			return err
		}
	}
	if t.LexicalDeclaration != nil {
		if err := fn(t.LexicalDeclaration); err != nil {
			return err
		}
	}
	return nil
}

func walkLexicalDeclaration(t *javascript.LexicalDeclaration, fn func(javascript.Type) error) error {
	for n := range t.BindingList {
		if err := fn(&t.BindingList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkLexicalBinding(t *javascript.LexicalBinding, fn func(javascript.Type) error) error {
	if t.ArrayBindingPattern != nil {
		if err := fn(t.ArrayBindingPattern); err != nil {
			return err
		}
	}
	if t.ObjectBindingPattern != nil {
		if err := fn(t.ObjectBindingPattern); err != nil {
			return err
		}
	}
	if t.Initializer != nil {
		if err := fn(t.Initializer); err != nil {
			return err
		}
	}
	return nil
}

func walkArrayBindingPattern(t *javascript.ArrayBindingPattern, fn func(javascript.Type) error) error {
	for n := range t.BindingElementList {
		if err := fn(&t.BindingElementList[n]); err != nil {
			return err
		}
	}
	if t.BindingRestElement != nil {
		if err := fn(t.BindingRestElement); err != nil {
			return err
		}
	}
	return nil
}

func walkObjectBindingPattern(t *javascript.ObjectBindingPattern, fn func(javascript.Type) error) error {
	for n := range t.BindingPropertyList {
		if err := fn(&t.BindingPropertyList[n]); err != nil {
			return err
		}
	}
	if t.BindingRestProperty != nil {
		if err := fn(t.BindingRestProperty); err != nil {
			return err
		}
	}
	return nil
}

func walkBindingProperty(t *javascript.BindingProperty, fn func(javascript.Type) error) error {
	if err := fn(&t.PropertyName); err != nil {
		return err
	}
	return fn(&t.BindingElement)
}

func walkVariableDeclaration(t *javascript.VariableDeclaration, fn func(javascript.Type) error) error {
	return walkLexicalBinding((*javascript.LexicalBinding)(t), fn)
}

func walkArrayLiteral(t *javascript.ArrayLiteral, fn func(javascript.Type) error) error {
	for n := range t.ElementList {
		if err := fn(&t.ElementList[n]); err != nil {
			return err
		}
	}
	if t.SpreadElement != nil {
		if err := fn(t.SpreadElement); err != nil {
			return err
		}
	}
	return nil
}

func walkObjectLiteral(t *javascript.ObjectLiteral, fn func(javascript.Type) error) error {
	for n := range t.PropertyDefinitionList {
		if err := fn(&t.PropertyDefinitionList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkPropertyDefinition(t *javascript.PropertyDefinition, fn func(javascript.Type) error) error {
	if t.PropertyName != nil {
		if err := fn(t.PropertyName); err != nil {
			return err
		}
	}
	if t.AssignmentExpression != nil {
		if err := fn(t.AssignmentExpression); err != nil {
			return err
		}
	}
	if t.MethodDefinition != nil {
		if err := fn(t.MethodDefinition); err != nil {
			return err
		}
	}
	return nil
}

func walkTemplateLiteral(t *javascript.TemplateLiteral, fn func(javascript.Type) error) error {
	for n := range t.Expressions {
		if err := fn(&t.Expressions[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkArrowFunction(t *javascript.ArrowFunction, fn func(javascript.Type) error) error {
	if t.CoverParenthesizedExpressionAndArrowParameterList != nil {
		if err := fn(t.CoverParenthesizedExpressionAndArrowParameterList); err != nil {
			return err
		}
	}
	if t.FormalParameters != nil {
		if err := fn(t.FormalParameters); err != nil {
			return err
		}
	}
	if t.AssignmentExpression != nil {
		if err := fn(t.AssignmentExpression); err != nil {
			return err
		}
	}
	if t.FunctionBody != nil {
		if err := fn(t.FunctionBody); err != nil {
			return err
		}
	}
	return nil
}

func walkModule(t *javascript.Module, fn func(javascript.Type) error) error {
	for n := range t.ModuleListItems {
		if err := fn(&t.ModuleListItems[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkModuleItem(t *javascript.ModuleItem, fn func(javascript.Type) error) error {
	if t.ImportDeclaration != nil {
		if err := fn(t.ImportDeclaration); err != nil {
			return err
		}
	}
	if t.StatementListItem != nil {
		if err := fn(t.StatementListItem); err != nil {
			return err
		}
	}
	if t.ExportDeclaration != nil {
		if err := fn(t.ExportDeclaration); err != nil {
			return err
		}
	}
	return nil
}

func walkImportDeclaration(t *javascript.ImportDeclaration, fn func(javascript.Type) error) error {
	if t.ImportClause != nil {
		if err := fn(t.ImportClause); err != nil {
			return err
		}
	}
	return fn(&t.FromClause)
}

func walkImportClause(t *javascript.ImportClause, fn func(javascript.Type) error) error {
	if t.NamedImports != nil {
		if err := fn(t.NamedImports); err != nil {
			return err
		}
	}
	return nil
}

func walkFromClause(t *javascript.FromClause, fn func(javascript.Type) error) error {
	return nil
}

func walkNamedImports(t *javascript.NamedImports, fn func(javascript.Type) error) error {
	for n := range t.ImportList {
		if err := fn(&t.ImportList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkImportSpecifier(t *javascript.ImportSpecifier, fn func(javascript.Type) error) error {
	return nil
}

func walkExportDeclaration(t *javascript.ExportDeclaration, fn func(javascript.Type) error) error {
	if t.ExportClause != nil {
		if err := fn(t.ExportClause); err != nil {
			return err
		}
	}
	if t.FromClause != nil {
		if err := fn(t.FromClause); err != nil {
			return err
		}
	}
	if t.VariableStatement != nil {
		if err := fn(t.VariableStatement); err != nil {
			return err
		}
	}
	if t.Declaration != nil {
		if err := fn(t.Declaration); err != nil {
			return err
		}
	}
	if t.DefaultFunction != nil {
		if err := fn(t.DefaultFunction); err != nil {
			return err
		}
	}
	if t.DefaultClass != nil {
		if err := fn(t.DefaultClass); err != nil {
			return err
		}
	}
	if t.DefaultAssignmentExpression != nil {
		if err := fn(t.DefaultAssignmentExpression); err != nil {
			return err
		}
	}
	return nil
}

func walkExportClause(t *javascript.ExportClause, fn func(javascript.Type) error) error {
	for n := range t.ExportList {
		if err := fn(&t.ExportList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkExportSpecifier(t *javascript.ExportSpecifier, fn func(javascript.Type) error) error {
	return nil
}

func walkBlock(t *javascript.Block, fn func(javascript.Type) error) error {
	for n := range t.StatementList {
		if err := fn(&t.StatementList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkStatementListItem(t *javascript.StatementListItem, fn func(javascript.Type) error) error {
	if t.Statement != nil {
		if err := fn(t.Statement); err != nil {
			return err
		}
	}
	if t.Declaration != nil {
		if err := fn(t.Declaration); err != nil {
			return err
		}
	}
	return nil
}

func walkStatement(t *javascript.Statement, fn func(javascript.Type) error) error {

	if t.BlockStatement != nil {
		if err := fn(t.BlockStatement); err != nil {
			return err
		}
	}
	if t.VariableStatement != nil {
		if err := fn(t.VariableStatement); err != nil {
			return err
		}
	}
	if t.ExpressionStatement != nil {
		if err := fn(t.ExpressionStatement); err != nil {
			return err
		}
	}
	if t.IfStatement != nil {
		if err := fn(t.IfStatement); err != nil {
			return err
		}
	}
	if t.IterationStatementDo != nil {
		if err := fn(t.IterationStatementDo); err != nil {
			return err
		}
	}
	if t.IterationStatementWhile != nil {
		if err := fn(t.IterationStatementWhile); err != nil {
			return err
		}
	}
	if t.IterationStatementFor != nil {
		if err := fn(t.IterationStatementFor); err != nil {
			return err
		}
	}
	if t.SwitchStatement != nil {
		if err := fn(t.SwitchStatement); err != nil {
			return err
		}
	}
	if t.WithStatement != nil {
		if err := fn(t.WithStatement); err != nil {
			return err
		}
	}
	if t.LabelledItemFunction != nil {
		if err := fn(t.LabelledItemFunction); err != nil {
			return err
		}
	}
	if t.LabelledItemStatement != nil {
		if err := fn(t.LabelledItemStatement); err != nil {
			return err
		}
	}
	if t.TryStatement != nil {
		if err := fn(t.TryStatement); err != nil {
			return err
		}
	}
	return nil
}

func walkIfStatement(t *javascript.IfStatement, fn func(javascript.Type) error) error {
	if err := fn(&t.Expression); err != nil {
		return err
	}
	if err := fn(&t.Statement); err != nil {
		return err
	}
	if t.ElseStatement != nil {
		if err := fn(t.ElseStatement); err != nil {
			return err
		}
	}
	return nil
}

func walkIterationStatementDo(t *javascript.IterationStatementDo, fn func(javascript.Type) error) error {
	if err := fn(&t.Statement); err != nil {
		return err
	}
	return fn(&t.Expression)
}

func walkIterationStatementWhile(t *javascript.IterationStatementWhile, fn func(javascript.Type) error) error {
	if err := fn(&t.Expression); err != nil {
		return err
	}
	return fn(&t.Statement)
}

func walkIterationStatementFor(t *javascript.IterationStatementFor, fn func(javascript.Type) error) error {
	if t.InitExpression != nil {
		if err := fn(t.InitExpression); err != nil {
			return err
		}
	}
	for n := range t.InitVar {
		if err := fn(&t.InitVar[n]); err != nil {
			return err
		}
	}
	if t.InitLexical != nil {
		if err := fn(t.InitLexical); err != nil {
			return err
		}
	}
	if t.Conditional != nil {
		if err := fn(t.Conditional); err != nil {
			return err
		}
	}
	if t.Afterthought != nil {
		if err := fn(t.Afterthought); err != nil {
			return err
		}
	}
	if t.LeftHandSideExpression != nil {
		if err := fn(t.LeftHandSideExpression); err != nil {
			return err
		}
	}
	if t.ForBindingPatternObject != nil {
		if err := fn(t.ForBindingPatternObject); err != nil {
			return err
		}
	}
	if t.ForBindingPatternArray != nil {
		if err := fn(t.ForBindingPatternArray); err != nil {
			return err
		}
	}
	if t.In != nil {
		if err := fn(t.In); err != nil {
			return err
		}
	}
	if t.Of != nil {
		if err := fn(t.Of); err != nil {
			return err
		}
	}
	return fn(&t.Statement)
}

func walkSwitchStatement(t *javascript.SwitchStatement, fn func(javascript.Type) error) error {
	if err := fn(&t.Expression); err != nil {
		return err
	}
	for n := range t.CaseClauses {
		if err := fn(&t.CaseClauses[n]); err != nil {
			return err
		}
	}
	for n := range t.DefaultClause {
		if err := fn(&t.DefaultClause[n]); err != nil {
			return err
		}
	}
	for n := range t.PostDefaultCaseClauses {
		if err := fn(&t.PostDefaultCaseClauses[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkCaseClause(t *javascript.CaseClause, fn func(javascript.Type) error) error {
	if err := fn(&t.Expression); err != nil {
		return err
	}
	for n := range t.StatementList {
		if err := fn(&t.StatementList[n]); err != nil {
			return err
		}
	}
	return nil
}

func walkWithStatement(t *javascript.WithStatement, fn func(javascript.Type) error) error {
	if err := fn(&t.Expression); err != nil {
		return err
	}
	return fn(&t.Statement)
}

func walkTryStatement(t *javascript.TryStatement, fn func(javascript.Type) error) error {
	if err := fn(&t.TryBlock); err != nil {
		return err
	}
	if t.CatchParameterObjectBindingPattern != nil {
		if err := fn(t.CatchParameterObjectBindingPattern); err != nil {
			return err
		}
	}
	if t.CatchParameterArrayBindingPattern != nil {
		if err := fn(t.CatchParameterArrayBindingPattern); err != nil {
			return err
		}
	}
	if t.CatchBlock != nil {
		if err := fn(t.CatchBlock); err != nil {
			return err
		}
	}
	if t.FinallyBlock != nil {
		if err := fn(t.FinallyBlock); err != nil {
			return err
		}
	}
	return nil
}

func walkVariableStatement(t *javascript.VariableStatement, fn func(javascript.Type) error) error {
	for n := range t.VariableDeclarationList {
		if err := fn(&t.VariableDeclarationList[n]); err != nil {
			return err
		}
	}
	return nil
}
