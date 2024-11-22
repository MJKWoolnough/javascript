package javascript

// File automatically generated with format.sh.

import "io"

func (f *AdditiveExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("AdditiveExpression {")

	if f.AdditiveExpression != nil {
		pp.Print("\nAdditiveExpression: ")
		f.AdditiveExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nAdditiveExpression: nil")
	}

	pp.Print("\nAdditiveOperator: ")
	f.AdditiveOperator.printType(&pp, v)

	pp.Print("\nMultiplicativeExpression: ")
	f.MultiplicativeExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Argument) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Argument {")

	if f.Spread || v {
		pp.Printf("\nSpread: %v", f.Spread)
	}

	pp.Print("\nAssignmentExpression: ")
	f.AssignmentExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Arguments) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Arguments {")

	if f.ArgumentList == nil {
		pp.Print("\nArgumentList: nil")
	} else if len(f.ArgumentList) > 0 {
		pp.Print("\nArgumentList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.ArgumentList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nArgumentList: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ArrayAssignmentPattern) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ArrayAssignmentPattern {")

	if f.AssignmentElements == nil {
		pp.Print("\nAssignmentElements: nil")
	} else if len(f.AssignmentElements) > 0 {
		pp.Print("\nAssignmentElements: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.AssignmentElements {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nAssignmentElements: []")
	}

	if f.AssignmentRestElement != nil {
		pp.Print("\nAssignmentRestElement: ")
		f.AssignmentRestElement.printType(&pp, v)
	} else if v {
		pp.Print("\nAssignmentRestElement: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ArrayBindingPattern) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ArrayBindingPattern {")

	if f.BindingElementList == nil {
		pp.Print("\nBindingElementList: nil")
	} else if len(f.BindingElementList) > 0 {
		pp.Print("\nBindingElementList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.BindingElementList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nBindingElementList: []")
	}

	if f.BindingRestElement != nil {
		pp.Print("\nBindingRestElement: ")
		f.BindingRestElement.printType(&pp, v)
	} else if v {
		pp.Print("\nBindingRestElement: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ArrayElement) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ArrayElement {")

	if f.Spread || v {
		pp.Printf("\nSpread: %v", f.Spread)
	}

	pp.Print("\nAssignmentExpression: ")
	f.AssignmentExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ArrayLiteral) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ArrayLiteral {")

	if f.ElementList == nil {
		pp.Print("\nElementList: nil")
	} else if len(f.ElementList) > 0 {
		pp.Print("\nElementList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.ElementList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nElementList: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ArrowFunction) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ArrowFunction {")

	if f.Async || v {
		pp.Printf("\nAsync: %v", f.Async)
	}

	if f.BindingIdentifier != nil {
		pp.Print("\nBindingIdentifier: ")
		f.BindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nBindingIdentifier: nil")
	}

	if f.FormalParameters != nil {
		pp.Print("\nFormalParameters: ")
		f.FormalParameters.printType(&pp, v)
	} else if v {
		pp.Print("\nFormalParameters: nil")
	}

	if f.AssignmentExpression != nil {
		pp.Print("\nAssignmentExpression: ")
		f.AssignmentExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nAssignmentExpression: nil")
	}

	if f.FunctionBody != nil {
		pp.Print("\nFunctionBody: ")
		f.FunctionBody.printType(&pp, v)
	} else if v {
		pp.Print("\nFunctionBody: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *AssignmentElement) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("AssignmentElement {")

	pp.Print("\nDestructuringAssignmentTarget: ")
	f.DestructuringAssignmentTarget.printType(&pp, v)

	if f.Initializer != nil {
		pp.Print("\nInitializer: ")
		f.Initializer.printType(&pp, v)
	} else if v {
		pp.Print("\nInitializer: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *AssignmentExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("AssignmentExpression {")

	if f.ConditionalExpression != nil {
		pp.Print("\nConditionalExpression: ")
		f.ConditionalExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nConditionalExpression: nil")
	}

	if f.ArrowFunction != nil {
		pp.Print("\nArrowFunction: ")
		f.ArrowFunction.printType(&pp, v)
	} else if v {
		pp.Print("\nArrowFunction: nil")
	}

	if f.LeftHandSideExpression != nil {
		pp.Print("\nLeftHandSideExpression: ")
		f.LeftHandSideExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nLeftHandSideExpression: nil")
	}

	if f.AssignmentPattern != nil {
		pp.Print("\nAssignmentPattern: ")
		f.AssignmentPattern.printType(&pp, v)
	} else if v {
		pp.Print("\nAssignmentPattern: nil")
	}

	if f.Yield || v {
		pp.Printf("\nYield: %v", f.Yield)
	}

	if f.Delegate || v {
		pp.Printf("\nDelegate: %v", f.Delegate)
	}

	pp.Print("\nAssignmentOperator: ")
	f.AssignmentOperator.printType(&pp, v)

	if f.AssignmentExpression != nil {
		pp.Print("\nAssignmentExpression: ")
		f.AssignmentExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nAssignmentExpression: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *AssignmentPattern) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("AssignmentPattern {")

	if f.ObjectAssignmentPattern != nil {
		pp.Print("\nObjectAssignmentPattern: ")
		f.ObjectAssignmentPattern.printType(&pp, v)
	} else if v {
		pp.Print("\nObjectAssignmentPattern: nil")
	}

	if f.ArrayAssignmentPattern != nil {
		pp.Print("\nArrayAssignmentPattern: ")
		f.ArrayAssignmentPattern.printType(&pp, v)
	} else if v {
		pp.Print("\nArrayAssignmentPattern: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *AssignmentProperty) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("AssignmentProperty {")

	pp.Print("\nPropertyName: ")
	f.PropertyName.printType(&pp, v)

	if f.DestructuringAssignmentTarget != nil {
		pp.Print("\nDestructuringAssignmentTarget: ")
		f.DestructuringAssignmentTarget.printType(&pp, v)
	} else if v {
		pp.Print("\nDestructuringAssignmentTarget: nil")
	}

	if f.Initializer != nil {
		pp.Print("\nInitializer: ")
		f.Initializer.printType(&pp, v)
	} else if v {
		pp.Print("\nInitializer: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *BindingElement) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("BindingElement {")

	if f.SingleNameBinding != nil {
		pp.Print("\nSingleNameBinding: ")
		f.SingleNameBinding.printType(&pp, v)
	} else if v {
		pp.Print("\nSingleNameBinding: nil")
	}

	if f.ArrayBindingPattern != nil {
		pp.Print("\nArrayBindingPattern: ")
		f.ArrayBindingPattern.printType(&pp, v)
	} else if v {
		pp.Print("\nArrayBindingPattern: nil")
	}

	if f.ObjectBindingPattern != nil {
		pp.Print("\nObjectBindingPattern: ")
		f.ObjectBindingPattern.printType(&pp, v)
	} else if v {
		pp.Print("\nObjectBindingPattern: nil")
	}

	if f.Initializer != nil {
		pp.Print("\nInitializer: ")
		f.Initializer.printType(&pp, v)
	} else if v {
		pp.Print("\nInitializer: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *BindingProperty) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("BindingProperty {")

	pp.Print("\nPropertyName: ")
	f.PropertyName.printType(&pp, v)

	pp.Print("\nBindingElement: ")
	f.BindingElement.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *BitwiseANDExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("BitwiseANDExpression {")

	if f.BitwiseANDExpression != nil {
		pp.Print("\nBitwiseANDExpression: ")
		f.BitwiseANDExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nBitwiseANDExpression: nil")
	}

	pp.Print("\nEqualityExpression: ")
	f.EqualityExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *BitwiseORExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("BitwiseORExpression {")

	if f.BitwiseORExpression != nil {
		pp.Print("\nBitwiseORExpression: ")
		f.BitwiseORExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nBitwiseORExpression: nil")
	}

	pp.Print("\nBitwiseXORExpression: ")
	f.BitwiseXORExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *BitwiseXORExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("BitwiseXORExpression {")

	if f.BitwiseXORExpression != nil {
		pp.Print("\nBitwiseXORExpression: ")
		f.BitwiseXORExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nBitwiseXORExpression: nil")
	}

	pp.Print("\nBitwiseANDExpression: ")
	f.BitwiseANDExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Block) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Block {")

	if f.StatementList == nil {
		pp.Print("\nStatementList: nil")
	} else if len(f.StatementList) > 0 {
		pp.Print("\nStatementList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.StatementList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nStatementList: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *CallExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("CallExpression {")

	if f.MemberExpression != nil {
		pp.Print("\nMemberExpression: ")
		f.MemberExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nMemberExpression: nil")
	}

	if f.SuperCall || v {
		pp.Printf("\nSuperCall: %v", f.SuperCall)
	}

	if f.ImportCall != nil {
		pp.Print("\nImportCall: ")
		f.ImportCall.printType(&pp, v)
	} else if v {
		pp.Print("\nImportCall: nil")
	}

	if f.CallExpression != nil {
		pp.Print("\nCallExpression: ")
		f.CallExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nCallExpression: nil")
	}

	if f.Arguments != nil {
		pp.Print("\nArguments: ")
		f.Arguments.printType(&pp, v)
	} else if v {
		pp.Print("\nArguments: nil")
	}

	if f.Expression != nil {
		pp.Print("\nExpression: ")
		f.Expression.printType(&pp, v)
	} else if v {
		pp.Print("\nExpression: nil")
	}

	if f.IdentifierName != nil {
		pp.Print("\nIdentifierName: ")
		f.IdentifierName.printType(&pp, v)
	} else if v {
		pp.Print("\nIdentifierName: nil")
	}

	if f.TemplateLiteral != nil {
		pp.Print("\nTemplateLiteral: ")
		f.TemplateLiteral.printType(&pp, v)
	} else if v {
		pp.Print("\nTemplateLiteral: nil")
	}

	if f.PrivateIdentifier != nil {
		pp.Print("\nPrivateIdentifier: ")
		f.PrivateIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nPrivateIdentifier: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *CaseClause) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("CaseClause {")

	pp.Print("\nExpression: ")
	f.Expression.printType(&pp, v)

	if f.StatementList == nil {
		pp.Print("\nStatementList: nil")
	} else if len(f.StatementList) > 0 {
		pp.Print("\nStatementList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.StatementList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nStatementList: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ClassDeclaration) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ClassDeclaration {")

	if f.BindingIdentifier != nil {
		pp.Print("\nBindingIdentifier: ")
		f.BindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nBindingIdentifier: nil")
	}

	if f.ClassHeritage != nil {
		pp.Print("\nClassHeritage: ")
		f.ClassHeritage.printType(&pp, v)
	} else if v {
		pp.Print("\nClassHeritage: nil")
	}

	if f.ClassBody == nil {
		pp.Print("\nClassBody: nil")
	} else if len(f.ClassBody) > 0 {
		pp.Print("\nClassBody: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.ClassBody {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nClassBody: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ClassElement) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ClassElement {")

	if f.Static || v {
		pp.Printf("\nStatic: %v", f.Static)
	}

	if f.MethodDefinition != nil {
		pp.Print("\nMethodDefinition: ")
		f.MethodDefinition.printType(&pp, v)
	} else if v {
		pp.Print("\nMethodDefinition: nil")
	}

	if f.FieldDefinition != nil {
		pp.Print("\nFieldDefinition: ")
		f.FieldDefinition.printType(&pp, v)
	} else if v {
		pp.Print("\nFieldDefinition: nil")
	}

	if f.ClassStaticBlock != nil {
		pp.Print("\nClassStaticBlock: ")
		f.ClassStaticBlock.printType(&pp, v)
	} else if v {
		pp.Print("\nClassStaticBlock: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ClassElementName) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ClassElementName {")

	if f.PropertyName != nil {
		pp.Print("\nPropertyName: ")
		f.PropertyName.printType(&pp, v)
	} else if v {
		pp.Print("\nPropertyName: nil")
	}

	if f.PrivateIdentifier != nil {
		pp.Print("\nPrivateIdentifier: ")
		f.PrivateIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nPrivateIdentifier: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *CoalesceExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("CoalesceExpression {")

	if f.CoalesceExpressionHead != nil {
		pp.Print("\nCoalesceExpressionHead: ")
		f.CoalesceExpressionHead.printType(&pp, v)
	} else if v {
		pp.Print("\nCoalesceExpressionHead: nil")
	}

	pp.Print("\nBitwiseORExpression: ")
	f.BitwiseORExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ConditionalExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ConditionalExpression {")

	if f.LogicalORExpression != nil {
		pp.Print("\nLogicalORExpression: ")
		f.LogicalORExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nLogicalORExpression: nil")
	}

	if f.CoalesceExpression != nil {
		pp.Print("\nCoalesceExpression: ")
		f.CoalesceExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nCoalesceExpression: nil")
	}

	if f.True != nil {
		pp.Print("\nTrue: ")
		f.True.printType(&pp, v)
	} else if v {
		pp.Print("\nTrue: nil")
	}

	if f.False != nil {
		pp.Print("\nFalse: ")
		f.False.printType(&pp, v)
	} else if v {
		pp.Print("\nFalse: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Declaration) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Declaration {")

	if f.ClassDeclaration != nil {
		pp.Print("\nClassDeclaration: ")
		f.ClassDeclaration.printType(&pp, v)
	} else if v {
		pp.Print("\nClassDeclaration: nil")
	}

	if f.FunctionDeclaration != nil {
		pp.Print("\nFunctionDeclaration: ")
		f.FunctionDeclaration.printType(&pp, v)
	} else if v {
		pp.Print("\nFunctionDeclaration: nil")
	}

	if f.LexicalDeclaration != nil {
		pp.Print("\nLexicalDeclaration: ")
		f.LexicalDeclaration.printType(&pp, v)
	} else if v {
		pp.Print("\nLexicalDeclaration: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *DestructuringAssignmentTarget) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("DestructuringAssignmentTarget {")

	if f.LeftHandSideExpression != nil {
		pp.Print("\nLeftHandSideExpression: ")
		f.LeftHandSideExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nLeftHandSideExpression: nil")
	}

	if f.AssignmentPattern != nil {
		pp.Print("\nAssignmentPattern: ")
		f.AssignmentPattern.printType(&pp, v)
	} else if v {
		pp.Print("\nAssignmentPattern: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *EqualityExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("EqualityExpression {")

	if f.EqualityExpression != nil {
		pp.Print("\nEqualityExpression: ")
		f.EqualityExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nEqualityExpression: nil")
	}

	pp.Print("\nEqualityOperator: ")
	f.EqualityOperator.printType(&pp, v)

	pp.Print("\nRelationalExpression: ")
	f.RelationalExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ExponentiationExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ExponentiationExpression {")

	if f.ExponentiationExpression != nil {
		pp.Print("\nExponentiationExpression: ")
		f.ExponentiationExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nExponentiationExpression: nil")
	}

	pp.Print("\nUnaryExpression: ")
	f.UnaryExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ExportClause) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ExportClause {")

	if f.ExportList == nil {
		pp.Print("\nExportList: nil")
	} else if len(f.ExportList) > 0 {
		pp.Print("\nExportList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.ExportList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nExportList: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ExportDeclaration) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ExportDeclaration {")

	if f.ExportClause != nil {
		pp.Print("\nExportClause: ")
		f.ExportClause.printType(&pp, v)
	} else if v {
		pp.Print("\nExportClause: nil")
	}

	if f.ExportFromClause != nil {
		pp.Print("\nExportFromClause: ")
		f.ExportFromClause.printType(&pp, v)
	} else if v {
		pp.Print("\nExportFromClause: nil")
	}

	if f.FromClause != nil {
		pp.Print("\nFromClause: ")
		f.FromClause.printType(&pp, v)
	} else if v {
		pp.Print("\nFromClause: nil")
	}

	if f.VariableStatement != nil {
		pp.Print("\nVariableStatement: ")
		f.VariableStatement.printType(&pp, v)
	} else if v {
		pp.Print("\nVariableStatement: nil")
	}

	if f.Declaration != nil {
		pp.Print("\nDeclaration: ")
		f.Declaration.printType(&pp, v)
	} else if v {
		pp.Print("\nDeclaration: nil")
	}

	if f.DefaultFunction != nil {
		pp.Print("\nDefaultFunction: ")
		f.DefaultFunction.printType(&pp, v)
	} else if v {
		pp.Print("\nDefaultFunction: nil")
	}

	if f.DefaultClass != nil {
		pp.Print("\nDefaultClass: ")
		f.DefaultClass.printType(&pp, v)
	} else if v {
		pp.Print("\nDefaultClass: nil")
	}

	if f.DefaultAssignmentExpression != nil {
		pp.Print("\nDefaultAssignmentExpression: ")
		f.DefaultAssignmentExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nDefaultAssignmentExpression: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ExportSpecifier) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ExportSpecifier {")

	if f.IdentifierName != nil {
		pp.Print("\nIdentifierName: ")
		f.IdentifierName.printType(&pp, v)
	} else if v {
		pp.Print("\nIdentifierName: nil")
	}

	if f.EIdentifierName != nil {
		pp.Print("\nEIdentifierName: ")
		f.EIdentifierName.printType(&pp, v)
	} else if v {
		pp.Print("\nEIdentifierName: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Expression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Expression {")

	if f.Expressions == nil {
		pp.Print("\nExpressions: nil")
	} else if len(f.Expressions) > 0 {
		pp.Print("\nExpressions: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Expressions {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nExpressions: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *FieldDefinition) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("FieldDefinition {")

	pp.Print("\nClassElementName: ")
	f.ClassElementName.printType(&pp, v)

	if f.Initializer != nil {
		pp.Print("\nInitializer: ")
		f.Initializer.printType(&pp, v)
	} else if v {
		pp.Print("\nInitializer: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *FormalParameters) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("FormalParameters {")

	if f.FormalParameterList == nil {
		pp.Print("\nFormalParameterList: nil")
	} else if len(f.FormalParameterList) > 0 {
		pp.Print("\nFormalParameterList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.FormalParameterList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nFormalParameterList: []")
	}

	if f.BindingIdentifier != nil {
		pp.Print("\nBindingIdentifier: ")
		f.BindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nBindingIdentifier: nil")
	}

	if f.ArrayBindingPattern != nil {
		pp.Print("\nArrayBindingPattern: ")
		f.ArrayBindingPattern.printType(&pp, v)
	} else if v {
		pp.Print("\nArrayBindingPattern: nil")
	}

	if f.ObjectBindingPattern != nil {
		pp.Print("\nObjectBindingPattern: ")
		f.ObjectBindingPattern.printType(&pp, v)
	} else if v {
		pp.Print("\nObjectBindingPattern: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *FromClause) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("FromClause {")

	if f.ModuleSpecifier != nil {
		pp.Print("\nModuleSpecifier: ")
		f.ModuleSpecifier.printType(&pp, v)
	} else if v {
		pp.Print("\nModuleSpecifier: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *FunctionDeclaration) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("FunctionDeclaration {")

	pp.Print("\nType: ")
	f.Type.printType(&pp, v)

	if f.BindingIdentifier != nil {
		pp.Print("\nBindingIdentifier: ")
		f.BindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nBindingIdentifier: nil")
	}

	pp.Print("\nFormalParameters: ")
	f.FormalParameters.printType(&pp, v)

	pp.Print("\nFunctionBody: ")
	f.FunctionBody.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *IfStatement) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("IfStatement {")

	pp.Print("\nExpression: ")
	f.Expression.printType(&pp, v)

	pp.Print("\nStatement: ")
	f.Statement.printType(&pp, v)

	if f.ElseStatement != nil {
		pp.Print("\nElseStatement: ")
		f.ElseStatement.printType(&pp, v)
	} else if v {
		pp.Print("\nElseStatement: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ImportClause) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ImportClause {")

	if f.ImportedDefaultBinding != nil {
		pp.Print("\nImportedDefaultBinding: ")
		f.ImportedDefaultBinding.printType(&pp, v)
	} else if v {
		pp.Print("\nImportedDefaultBinding: nil")
	}

	if f.NameSpaceImport != nil {
		pp.Print("\nNameSpaceImport: ")
		f.NameSpaceImport.printType(&pp, v)
	} else if v {
		pp.Print("\nNameSpaceImport: nil")
	}

	if f.NamedImports != nil {
		pp.Print("\nNamedImports: ")
		f.NamedImports.printType(&pp, v)
	} else if v {
		pp.Print("\nNamedImports: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ImportDeclaration) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ImportDeclaration {")

	if f.ImportClause != nil {
		pp.Print("\nImportClause: ")
		f.ImportClause.printType(&pp, v)
	} else if v {
		pp.Print("\nImportClause: nil")
	}

	pp.Print("\nFromClause: ")
	f.FromClause.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ImportSpecifier) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ImportSpecifier {")

	if f.IdentifierName != nil {
		pp.Print("\nIdentifierName: ")
		f.IdentifierName.printType(&pp, v)
	} else if v {
		pp.Print("\nIdentifierName: nil")
	}

	if f.ImportedBinding != nil {
		pp.Print("\nImportedBinding: ")
		f.ImportedBinding.printType(&pp, v)
	} else if v {
		pp.Print("\nImportedBinding: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *IterationStatementDo) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("IterationStatementDo {")

	pp.Print("\nStatement: ")
	f.Statement.printType(&pp, v)

	pp.Print("\nExpression: ")
	f.Expression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *IterationStatementFor) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("IterationStatementFor {")

	pp.Print("\nType: ")
	f.Type.printType(&pp, v)

	if f.InitExpression != nil {
		pp.Print("\nInitExpression: ")
		f.InitExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nInitExpression: nil")
	}

	if f.InitVar == nil {
		pp.Print("\nInitVar: nil")
	} else if len(f.InitVar) > 0 {
		pp.Print("\nInitVar: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.InitVar {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nInitVar: []")
	}

	if f.InitLexical != nil {
		pp.Print("\nInitLexical: ")
		f.InitLexical.printType(&pp, v)
	} else if v {
		pp.Print("\nInitLexical: nil")
	}

	if f.Conditional != nil {
		pp.Print("\nConditional: ")
		f.Conditional.printType(&pp, v)
	} else if v {
		pp.Print("\nConditional: nil")
	}

	if f.Afterthought != nil {
		pp.Print("\nAfterthought: ")
		f.Afterthought.printType(&pp, v)
	} else if v {
		pp.Print("\nAfterthought: nil")
	}

	if f.LeftHandSideExpression != nil {
		pp.Print("\nLeftHandSideExpression: ")
		f.LeftHandSideExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nLeftHandSideExpression: nil")
	}

	if f.ForBindingIdentifier != nil {
		pp.Print("\nForBindingIdentifier: ")
		f.ForBindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nForBindingIdentifier: nil")
	}

	if f.ForBindingPatternObject != nil {
		pp.Print("\nForBindingPatternObject: ")
		f.ForBindingPatternObject.printType(&pp, v)
	} else if v {
		pp.Print("\nForBindingPatternObject: nil")
	}

	if f.ForBindingPatternArray != nil {
		pp.Print("\nForBindingPatternArray: ")
		f.ForBindingPatternArray.printType(&pp, v)
	} else if v {
		pp.Print("\nForBindingPatternArray: nil")
	}

	if f.In != nil {
		pp.Print("\nIn: ")
		f.In.printType(&pp, v)
	} else if v {
		pp.Print("\nIn: nil")
	}

	if f.Of != nil {
		pp.Print("\nOf: ")
		f.Of.printType(&pp, v)
	} else if v {
		pp.Print("\nOf: nil")
	}

	pp.Print("\nStatement: ")
	f.Statement.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *IterationStatementWhile) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("IterationStatementWhile {")

	pp.Print("\nExpression: ")
	f.Expression.printType(&pp, v)

	pp.Print("\nStatement: ")
	f.Statement.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *LeftHandSideExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("LeftHandSideExpression {")

	if f.NewExpression != nil {
		pp.Print("\nNewExpression: ")
		f.NewExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nNewExpression: nil")
	}

	if f.CallExpression != nil {
		pp.Print("\nCallExpression: ")
		f.CallExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nCallExpression: nil")
	}

	if f.OptionalExpression != nil {
		pp.Print("\nOptionalExpression: ")
		f.OptionalExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nOptionalExpression: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *LexicalBinding) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("LexicalBinding {")

	if f.BindingIdentifier != nil {
		pp.Print("\nBindingIdentifier: ")
		f.BindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nBindingIdentifier: nil")
	}

	if f.ArrayBindingPattern != nil {
		pp.Print("\nArrayBindingPattern: ")
		f.ArrayBindingPattern.printType(&pp, v)
	} else if v {
		pp.Print("\nArrayBindingPattern: nil")
	}

	if f.ObjectBindingPattern != nil {
		pp.Print("\nObjectBindingPattern: ")
		f.ObjectBindingPattern.printType(&pp, v)
	} else if v {
		pp.Print("\nObjectBindingPattern: nil")
	}

	if f.Initializer != nil {
		pp.Print("\nInitializer: ")
		f.Initializer.printType(&pp, v)
	} else if v {
		pp.Print("\nInitializer: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *LexicalDeclaration) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("LexicalDeclaration {")

	pp.Print("\nLetOrConst: ")
	f.LetOrConst.printType(&pp, v)

	if f.BindingList == nil {
		pp.Print("\nBindingList: nil")
	} else if len(f.BindingList) > 0 {
		pp.Print("\nBindingList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.BindingList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nBindingList: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *LogicalANDExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("LogicalANDExpression {")

	if f.LogicalANDExpression != nil {
		pp.Print("\nLogicalANDExpression: ")
		f.LogicalANDExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nLogicalANDExpression: nil")
	}

	pp.Print("\nBitwiseORExpression: ")
	f.BitwiseORExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *LogicalORExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("LogicalORExpression {")

	if f.LogicalORExpression != nil {
		pp.Print("\nLogicalORExpression: ")
		f.LogicalORExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nLogicalORExpression: nil")
	}

	pp.Print("\nLogicalANDExpression: ")
	f.LogicalANDExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *MemberExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("MemberExpression {")

	if f.MemberExpression != nil {
		pp.Print("\nMemberExpression: ")
		f.MemberExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nMemberExpression: nil")
	}

	if f.PrimaryExpression != nil {
		pp.Print("\nPrimaryExpression: ")
		f.PrimaryExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nPrimaryExpression: nil")
	}

	if f.Expression != nil {
		pp.Print("\nExpression: ")
		f.Expression.printType(&pp, v)
	} else if v {
		pp.Print("\nExpression: nil")
	}

	if f.IdentifierName != nil {
		pp.Print("\nIdentifierName: ")
		f.IdentifierName.printType(&pp, v)
	} else if v {
		pp.Print("\nIdentifierName: nil")
	}

	if f.TemplateLiteral != nil {
		pp.Print("\nTemplateLiteral: ")
		f.TemplateLiteral.printType(&pp, v)
	} else if v {
		pp.Print("\nTemplateLiteral: nil")
	}

	if f.SuperProperty || v {
		pp.Printf("\nSuperProperty: %v", f.SuperProperty)
	}

	if f.NewTarget || v {
		pp.Printf("\nNewTarget: %v", f.NewTarget)
	}

	if f.ImportMeta || v {
		pp.Printf("\nImportMeta: %v", f.ImportMeta)
	}

	if f.Arguments != nil {
		pp.Print("\nArguments: ")
		f.Arguments.printType(&pp, v)
	} else if v {
		pp.Print("\nArguments: nil")
	}

	if f.PrivateIdentifier != nil {
		pp.Print("\nPrivateIdentifier: ")
		f.PrivateIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nPrivateIdentifier: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *MethodDefinition) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("MethodDefinition {")

	pp.Print("\nType: ")
	f.Type.printType(&pp, v)

	pp.Print("\nClassElementName: ")
	f.ClassElementName.printType(&pp, v)

	pp.Print("\nParams: ")
	f.Params.printType(&pp, v)

	pp.Print("\nFunctionBody: ")
	f.FunctionBody.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Module) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Module {")

	if f.ModuleListItems == nil {
		pp.Print("\nModuleListItems: nil")
	} else if len(f.ModuleListItems) > 0 {
		pp.Print("\nModuleListItems: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.ModuleListItems {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nModuleListItems: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ModuleItem) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ModuleItem {")

	if f.ImportDeclaration != nil {
		pp.Print("\nImportDeclaration: ")
		f.ImportDeclaration.printType(&pp, v)
	} else if v {
		pp.Print("\nImportDeclaration: nil")
	}

	if f.StatementListItem != nil {
		pp.Print("\nStatementListItem: ")
		f.StatementListItem.printType(&pp, v)
	} else if v {
		pp.Print("\nStatementListItem: nil")
	}

	if f.ExportDeclaration != nil {
		pp.Print("\nExportDeclaration: ")
		f.ExportDeclaration.printType(&pp, v)
	} else if v {
		pp.Print("\nExportDeclaration: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *MultiplicativeExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("MultiplicativeExpression {")

	if f.MultiplicativeExpression != nil {
		pp.Print("\nMultiplicativeExpression: ")
		f.MultiplicativeExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nMultiplicativeExpression: nil")
	}

	pp.Print("\nMultiplicativeOperator: ")
	f.MultiplicativeOperator.printType(&pp, v)

	pp.Print("\nExponentiationExpression: ")
	f.ExponentiationExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *NamedImports) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("NamedImports {")

	if f.ImportList == nil {
		pp.Print("\nImportList: nil")
	} else if len(f.ImportList) > 0 {
		pp.Print("\nImportList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.ImportList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nImportList: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *NewExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("NewExpression {")

	if f.News != 0 || v {
		pp.Printf("\nNews: %v", f.News)
	}

	pp.Print("\nMemberExpression: ")
	f.MemberExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ObjectAssignmentPattern) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ObjectAssignmentPattern {")

	if f.AssignmentPropertyList == nil {
		pp.Print("\nAssignmentPropertyList: nil")
	} else if len(f.AssignmentPropertyList) > 0 {
		pp.Print("\nAssignmentPropertyList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.AssignmentPropertyList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nAssignmentPropertyList: []")
	}

	if f.AssignmentRestElement != nil {
		pp.Print("\nAssignmentRestElement: ")
		f.AssignmentRestElement.printType(&pp, v)
	} else if v {
		pp.Print("\nAssignmentRestElement: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ObjectBindingPattern) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ObjectBindingPattern {")

	if f.BindingPropertyList == nil {
		pp.Print("\nBindingPropertyList: nil")
	} else if len(f.BindingPropertyList) > 0 {
		pp.Print("\nBindingPropertyList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.BindingPropertyList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nBindingPropertyList: []")
	}

	if f.BindingRestProperty != nil {
		pp.Print("\nBindingRestProperty: ")
		f.BindingRestProperty.printType(&pp, v)
	} else if v {
		pp.Print("\nBindingRestProperty: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ObjectLiteral) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ObjectLiteral {")

	if f.PropertyDefinitionList == nil {
		pp.Print("\nPropertyDefinitionList: nil")
	} else if len(f.PropertyDefinitionList) > 0 {
		pp.Print("\nPropertyDefinitionList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.PropertyDefinitionList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nPropertyDefinitionList: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *OptionalChain) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("OptionalChain {")

	if f.OptionalChain != nil {
		pp.Print("\nOptionalChain: ")
		f.OptionalChain.printType(&pp, v)
	} else if v {
		pp.Print("\nOptionalChain: nil")
	}

	if f.Arguments != nil {
		pp.Print("\nArguments: ")
		f.Arguments.printType(&pp, v)
	} else if v {
		pp.Print("\nArguments: nil")
	}

	if f.Expression != nil {
		pp.Print("\nExpression: ")
		f.Expression.printType(&pp, v)
	} else if v {
		pp.Print("\nExpression: nil")
	}

	if f.IdentifierName != nil {
		pp.Print("\nIdentifierName: ")
		f.IdentifierName.printType(&pp, v)
	} else if v {
		pp.Print("\nIdentifierName: nil")
	}

	if f.TemplateLiteral != nil {
		pp.Print("\nTemplateLiteral: ")
		f.TemplateLiteral.printType(&pp, v)
	} else if v {
		pp.Print("\nTemplateLiteral: nil")
	}

	if f.PrivateIdentifier != nil {
		pp.Print("\nPrivateIdentifier: ")
		f.PrivateIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nPrivateIdentifier: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *OptionalExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("OptionalExpression {")

	if f.MemberExpression != nil {
		pp.Print("\nMemberExpression: ")
		f.MemberExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nMemberExpression: nil")
	}

	if f.CallExpression != nil {
		pp.Print("\nCallExpression: ")
		f.CallExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nCallExpression: nil")
	}

	if f.OptionalExpression != nil {
		pp.Print("\nOptionalExpression: ")
		f.OptionalExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nOptionalExpression: nil")
	}

	pp.Print("\nOptionalChain: ")
	f.OptionalChain.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ParenthesizedExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ParenthesizedExpression {")

	if f.Expressions == nil {
		pp.Print("\nExpressions: nil")
	} else if len(f.Expressions) > 0 {
		pp.Print("\nExpressions: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Expressions {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nExpressions: []")
	}

	if f.bindingIdentifier != nil {
		pp.Print("\nbindingIdentifier: ")
		f.bindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nbindingIdentifier: nil")
	}

	if f.arrayBindingPattern != nil {
		pp.Print("\narrayBindingPattern: ")
		f.arrayBindingPattern.printType(&pp, v)
	} else if v {
		pp.Print("\narrayBindingPattern: nil")
	}

	if f.objectBindingPattern != nil {
		pp.Print("\nobjectBindingPattern: ")
		f.objectBindingPattern.printType(&pp, v)
	} else if v {
		pp.Print("\nobjectBindingPattern: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *PrimaryExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("PrimaryExpression {")

	if f.This != nil {
		pp.Print("\nThis: ")
		f.This.printType(&pp, v)
	} else if v {
		pp.Print("\nThis: nil")
	}

	if f.IdentifierReference != nil {
		pp.Print("\nIdentifierReference: ")
		f.IdentifierReference.printType(&pp, v)
	} else if v {
		pp.Print("\nIdentifierReference: nil")
	}

	if f.Literal != nil {
		pp.Print("\nLiteral: ")
		f.Literal.printType(&pp, v)
	} else if v {
		pp.Print("\nLiteral: nil")
	}

	if f.ArrayLiteral != nil {
		pp.Print("\nArrayLiteral: ")
		f.ArrayLiteral.printType(&pp, v)
	} else if v {
		pp.Print("\nArrayLiteral: nil")
	}

	if f.ObjectLiteral != nil {
		pp.Print("\nObjectLiteral: ")
		f.ObjectLiteral.printType(&pp, v)
	} else if v {
		pp.Print("\nObjectLiteral: nil")
	}

	if f.FunctionExpression != nil {
		pp.Print("\nFunctionExpression: ")
		f.FunctionExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nFunctionExpression: nil")
	}

	if f.ClassExpression != nil {
		pp.Print("\nClassExpression: ")
		f.ClassExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nClassExpression: nil")
	}

	if f.TemplateLiteral != nil {
		pp.Print("\nTemplateLiteral: ")
		f.TemplateLiteral.printType(&pp, v)
	} else if v {
		pp.Print("\nTemplateLiteral: nil")
	}

	if f.ParenthesizedExpression != nil {
		pp.Print("\nParenthesizedExpression: ")
		f.ParenthesizedExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nParenthesizedExpression: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *PropertyDefinition) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("PropertyDefinition {")

	if f.IsCoverInitializedName || v {
		pp.Printf("\nIsCoverInitializedName: %v", f.IsCoverInitializedName)
	}

	if f.PropertyName != nil {
		pp.Print("\nPropertyName: ")
		f.PropertyName.printType(&pp, v)
	} else if v {
		pp.Print("\nPropertyName: nil")
	}

	if f.AssignmentExpression != nil {
		pp.Print("\nAssignmentExpression: ")
		f.AssignmentExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nAssignmentExpression: nil")
	}

	if f.MethodDefinition != nil {
		pp.Print("\nMethodDefinition: ")
		f.MethodDefinition.printType(&pp, v)
	} else if v {
		pp.Print("\nMethodDefinition: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *PropertyName) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("PropertyName {")

	if f.LiteralPropertyName != nil {
		pp.Print("\nLiteralPropertyName: ")
		f.LiteralPropertyName.printType(&pp, v)
	} else if v {
		pp.Print("\nLiteralPropertyName: nil")
	}

	if f.ComputedPropertyName != nil {
		pp.Print("\nComputedPropertyName: ")
		f.ComputedPropertyName.printType(&pp, v)
	} else if v {
		pp.Print("\nComputedPropertyName: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *RelationalExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("RelationalExpression {")

	if f.PrivateIdentifier != nil {
		pp.Print("\nPrivateIdentifier: ")
		f.PrivateIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nPrivateIdentifier: nil")
	}

	if f.RelationalExpression != nil {
		pp.Print("\nRelationalExpression: ")
		f.RelationalExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nRelationalExpression: nil")
	}

	pp.Print("\nRelationshipOperator: ")
	f.RelationshipOperator.printType(&pp, v)

	pp.Print("\nShiftExpression: ")
	f.ShiftExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Script) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Script {")

	if f.StatementList == nil {
		pp.Print("\nStatementList: nil")
	} else if len(f.StatementList) > 0 {
		pp.Print("\nStatementList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.StatementList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nStatementList: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ShiftExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ShiftExpression {")

	if f.ShiftExpression != nil {
		pp.Print("\nShiftExpression: ")
		f.ShiftExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nShiftExpression: nil")
	}

	pp.Print("\nShiftOperator: ")
	f.ShiftOperator.printType(&pp, v)

	pp.Print("\nAdditiveExpression: ")
	f.AdditiveExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Statement) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Statement {")

	pp.Print("\nType: ")
	f.Type.printType(&pp, v)

	if f.BlockStatement != nil {
		pp.Print("\nBlockStatement: ")
		f.BlockStatement.printType(&pp, v)
	} else if v {
		pp.Print("\nBlockStatement: nil")
	}

	if f.VariableStatement != nil {
		pp.Print("\nVariableStatement: ")
		f.VariableStatement.printType(&pp, v)
	} else if v {
		pp.Print("\nVariableStatement: nil")
	}

	if f.ExpressionStatement != nil {
		pp.Print("\nExpressionStatement: ")
		f.ExpressionStatement.printType(&pp, v)
	} else if v {
		pp.Print("\nExpressionStatement: nil")
	}

	if f.IfStatement != nil {
		pp.Print("\nIfStatement: ")
		f.IfStatement.printType(&pp, v)
	} else if v {
		pp.Print("\nIfStatement: nil")
	}

	if f.IterationStatementDo != nil {
		pp.Print("\nIterationStatementDo: ")
		f.IterationStatementDo.printType(&pp, v)
	} else if v {
		pp.Print("\nIterationStatementDo: nil")
	}

	if f.IterationStatementWhile != nil {
		pp.Print("\nIterationStatementWhile: ")
		f.IterationStatementWhile.printType(&pp, v)
	} else if v {
		pp.Print("\nIterationStatementWhile: nil")
	}

	if f.IterationStatementFor != nil {
		pp.Print("\nIterationStatementFor: ")
		f.IterationStatementFor.printType(&pp, v)
	} else if v {
		pp.Print("\nIterationStatementFor: nil")
	}

	if f.SwitchStatement != nil {
		pp.Print("\nSwitchStatement: ")
		f.SwitchStatement.printType(&pp, v)
	} else if v {
		pp.Print("\nSwitchStatement: nil")
	}

	if f.WithStatement != nil {
		pp.Print("\nWithStatement: ")
		f.WithStatement.printType(&pp, v)
	} else if v {
		pp.Print("\nWithStatement: nil")
	}

	if f.LabelIdentifier != nil {
		pp.Print("\nLabelIdentifier: ")
		f.LabelIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nLabelIdentifier: nil")
	}

	if f.LabelledItemFunction != nil {
		pp.Print("\nLabelledItemFunction: ")
		f.LabelledItemFunction.printType(&pp, v)
	} else if v {
		pp.Print("\nLabelledItemFunction: nil")
	}

	if f.LabelledItemStatement != nil {
		pp.Print("\nLabelledItemStatement: ")
		f.LabelledItemStatement.printType(&pp, v)
	} else if v {
		pp.Print("\nLabelledItemStatement: nil")
	}

	if f.TryStatement != nil {
		pp.Print("\nTryStatement: ")
		f.TryStatement.printType(&pp, v)
	} else if v {
		pp.Print("\nTryStatement: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *StatementListItem) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("StatementListItem {")

	if f.Statement != nil {
		pp.Print("\nStatement: ")
		f.Statement.printType(&pp, v)
	} else if v {
		pp.Print("\nStatement: nil")
	}

	if f.Declaration != nil {
		pp.Print("\nDeclaration: ")
		f.Declaration.printType(&pp, v)
	} else if v {
		pp.Print("\nDeclaration: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *SwitchStatement) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("SwitchStatement {")

	pp.Print("\nExpression: ")
	f.Expression.printType(&pp, v)

	if f.CaseClauses == nil {
		pp.Print("\nCaseClauses: nil")
	} else if len(f.CaseClauses) > 0 {
		pp.Print("\nCaseClauses: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.CaseClauses {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nCaseClauses: []")
	}

	if f.DefaultClause == nil {
		pp.Print("\nDefaultClause: nil")
	} else if len(f.DefaultClause) > 0 {
		pp.Print("\nDefaultClause: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.DefaultClause {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nDefaultClause: []")
	}

	if f.PostDefaultCaseClauses == nil {
		pp.Print("\nPostDefaultCaseClauses: nil")
	} else if len(f.PostDefaultCaseClauses) > 0 {
		pp.Print("\nPostDefaultCaseClauses: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.PostDefaultCaseClauses {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nPostDefaultCaseClauses: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *TemplateLiteral) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("TemplateLiteral {")

	if f.NoSubstitutionTemplate != nil {
		pp.Print("\nNoSubstitutionTemplate: ")
		f.NoSubstitutionTemplate.printType(&pp, v)
	} else if v {
		pp.Print("\nNoSubstitutionTemplate: nil")
	}

	if f.TemplateHead != nil {
		pp.Print("\nTemplateHead: ")
		f.TemplateHead.printType(&pp, v)
	} else if v {
		pp.Print("\nTemplateHead: nil")
	}

	if f.Expressions == nil {
		pp.Print("\nExpressions: nil")
	} else if len(f.Expressions) > 0 {
		pp.Print("\nExpressions: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Expressions {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nExpressions: []")
	}

	if f.TemplateMiddleList == nil {
		pp.Print("\nTemplateMiddleList: nil")
	} else if len(f.TemplateMiddleList) > 0 {
		pp.Print("\nTemplateMiddleList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.TemplateMiddleList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nTemplateMiddleList: []")
	}

	if f.TemplateTail != nil {
		pp.Print("\nTemplateTail: ")
		f.TemplateTail.printType(&pp, v)
	} else if v {
		pp.Print("\nTemplateTail: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *TryStatement) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("TryStatement {")

	pp.Print("\nTryBlock: ")
	f.TryBlock.printType(&pp, v)

	if f.CatchParameterBindingIdentifier != nil {
		pp.Print("\nCatchParameterBindingIdentifier: ")
		f.CatchParameterBindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nCatchParameterBindingIdentifier: nil")
	}

	if f.CatchParameterObjectBindingPattern != nil {
		pp.Print("\nCatchParameterObjectBindingPattern: ")
		f.CatchParameterObjectBindingPattern.printType(&pp, v)
	} else if v {
		pp.Print("\nCatchParameterObjectBindingPattern: nil")
	}

	if f.CatchParameterArrayBindingPattern != nil {
		pp.Print("\nCatchParameterArrayBindingPattern: ")
		f.CatchParameterArrayBindingPattern.printType(&pp, v)
	} else if v {
		pp.Print("\nCatchParameterArrayBindingPattern: nil")
	}

	if f.CatchBlock != nil {
		pp.Print("\nCatchBlock: ")
		f.CatchBlock.printType(&pp, v)
	} else if v {
		pp.Print("\nCatchBlock: nil")
	}

	if f.FinallyBlock != nil {
		pp.Print("\nFinallyBlock: ")
		f.FinallyBlock.printType(&pp, v)
	} else if v {
		pp.Print("\nFinallyBlock: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *UnaryExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("UnaryExpression {")

	if f.UnaryOperators == nil {
		pp.Print("\nUnaryOperators: nil")
	} else if len(f.UnaryOperators) > 0 {
		pp.Print("\nUnaryOperators: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.UnaryOperators {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nUnaryOperators: []")
	}

	pp.Print("\nUpdateExpression: ")
	f.UpdateExpression.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *UpdateExpression) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("UpdateExpression {")

	if f.LeftHandSideExpression != nil {
		pp.Print("\nLeftHandSideExpression: ")
		f.LeftHandSideExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nLeftHandSideExpression: nil")
	}

	pp.Print("\nUpdateOperator: ")
	f.UpdateOperator.printType(&pp, v)

	if f.UnaryExpression != nil {
		pp.Print("\nUnaryExpression: ")
		f.UnaryExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nUnaryExpression: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *VariableStatement) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("VariableStatement {")

	if f.VariableDeclarationList == nil {
		pp.Print("\nVariableDeclarationList: nil")
	} else if len(f.VariableDeclarationList) > 0 {
		pp.Print("\nVariableDeclarationList: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.VariableDeclarationList {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nVariableDeclarationList: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *WithStatement) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("WithStatement {")

	pp.Print("\nExpression: ")
	f.Expression.printType(&pp, v)

	pp.Print("\nStatement: ")
	f.Statement.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}
