package javascript

// File automatically generated with format.sh.

func (f *AdditiveExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("AdditiveExpression {")

	if f.AdditiveExpression != nil {
		pp.WriteString("\nAdditiveExpression: ")
		f.AdditiveExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nAdditiveExpression: nil")
	}

	pp.WriteString("\nAdditiveOperator: ")
	f.AdditiveOperator.printType(pp, v)

	pp.WriteString("\nMultiplicativeExpression: ")
	f.MultiplicativeExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Argument) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Argument {")

	if f.Spread || v {
		pp.Printf("\nSpread: %v", f.Spread)
	}

	pp.WriteString("\nAssignmentExpression: ")
	f.AssignmentExpression.printType(pp, v)

	pp.WriteString("\nComments: ")
	f.Comments.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Arguments) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Arguments {")

	if f.ArgumentList == nil {
		pp.WriteString("\nArgumentList: nil")
	} else if len(f.ArgumentList) > 0 {
		pp.WriteString("\nArgumentList: [")

		ipp := pp.Indent()

		for n, e := range f.ArgumentList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nArgumentList: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ArrayAssignmentPattern) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ArrayAssignmentPattern {")

	if f.AssignmentElements == nil {
		pp.WriteString("\nAssignmentElements: nil")
	} else if len(f.AssignmentElements) > 0 {
		pp.WriteString("\nAssignmentElements: [")

		ipp := pp.Indent()

		for n, e := range f.AssignmentElements {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nAssignmentElements: []")
	}

	if f.AssignmentRestElement != nil {
		pp.WriteString("\nAssignmentRestElement: ")
		f.AssignmentRestElement.printType(pp, v)
	} else if v {
		pp.WriteString("\nAssignmentRestElement: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ArrayBindingPattern) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ArrayBindingPattern {")

	if f.BindingElementList == nil {
		pp.WriteString("\nBindingElementList: nil")
	} else if len(f.BindingElementList) > 0 {
		pp.WriteString("\nBindingElementList: [")

		ipp := pp.Indent()

		for n, e := range f.BindingElementList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nBindingElementList: []")
	}

	if f.BindingRestElement != nil {
		pp.WriteString("\nBindingRestElement: ")
		f.BindingRestElement.printType(pp, v)
	} else if v {
		pp.WriteString("\nBindingRestElement: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ArrayElement) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ArrayElement {")

	if f.Spread || v {
		pp.Printf("\nSpread: %v", f.Spread)
	}

	pp.WriteString("\nAssignmentExpression: ")
	f.AssignmentExpression.printType(pp, v)

	pp.WriteString("\nComments: ")
	f.Comments.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ArrayLiteral) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ArrayLiteral {")

	if f.ElementList == nil {
		pp.WriteString("\nElementList: nil")
	} else if len(f.ElementList) > 0 {
		pp.WriteString("\nElementList: [")

		ipp := pp.Indent()

		for n, e := range f.ElementList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nElementList: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ArrowFunction) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ArrowFunction {")

	if f.Async || v {
		pp.Printf("\nAsync: %v", f.Async)
	}

	if f.BindingIdentifier != nil {
		pp.WriteString("\nBindingIdentifier: ")
		f.BindingIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nBindingIdentifier: nil")
	}

	if f.FormalParameters != nil {
		pp.WriteString("\nFormalParameters: ")
		f.FormalParameters.printType(pp, v)
	} else if v {
		pp.WriteString("\nFormalParameters: nil")
	}

	if f.AssignmentExpression != nil {
		pp.WriteString("\nAssignmentExpression: ")
		f.AssignmentExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nAssignmentExpression: nil")
	}

	if f.FunctionBody != nil {
		pp.WriteString("\nFunctionBody: ")
		f.FunctionBody.printType(pp, v)
	} else if v {
		pp.WriteString("\nFunctionBody: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *AssignmentElement) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("AssignmentElement {")

	pp.WriteString("\nDestructuringAssignmentTarget: ")
	f.DestructuringAssignmentTarget.printType(pp, v)

	if f.Initializer != nil {
		pp.WriteString("\nInitializer: ")
		f.Initializer.printType(pp, v)
	} else if v {
		pp.WriteString("\nInitializer: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *AssignmentExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("AssignmentExpression {")

	if f.ConditionalExpression != nil {
		pp.WriteString("\nConditionalExpression: ")
		f.ConditionalExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nConditionalExpression: nil")
	}

	if f.ArrowFunction != nil {
		pp.WriteString("\nArrowFunction: ")
		f.ArrowFunction.printType(pp, v)
	} else if v {
		pp.WriteString("\nArrowFunction: nil")
	}

	if f.LeftHandSideExpression != nil {
		pp.WriteString("\nLeftHandSideExpression: ")
		f.LeftHandSideExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nLeftHandSideExpression: nil")
	}

	if f.AssignmentPattern != nil {
		pp.WriteString("\nAssignmentPattern: ")
		f.AssignmentPattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nAssignmentPattern: nil")
	}

	if f.Yield || v {
		pp.Printf("\nYield: %v", f.Yield)
	}

	if f.Delegate || v {
		pp.Printf("\nDelegate: %v", f.Delegate)
	}

	pp.WriteString("\nAssignmentOperator: ")
	f.AssignmentOperator.printType(pp, v)

	if f.AssignmentExpression != nil {
		pp.WriteString("\nAssignmentExpression: ")
		f.AssignmentExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nAssignmentExpression: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *AssignmentPattern) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("AssignmentPattern {")

	if f.ObjectAssignmentPattern != nil {
		pp.WriteString("\nObjectAssignmentPattern: ")
		f.ObjectAssignmentPattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nObjectAssignmentPattern: nil")
	}

	if f.ArrayAssignmentPattern != nil {
		pp.WriteString("\nArrayAssignmentPattern: ")
		f.ArrayAssignmentPattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nArrayAssignmentPattern: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *AssignmentProperty) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("AssignmentProperty {")

	pp.WriteString("\nPropertyName: ")
	f.PropertyName.printType(pp, v)

	if f.DestructuringAssignmentTarget != nil {
		pp.WriteString("\nDestructuringAssignmentTarget: ")
		f.DestructuringAssignmentTarget.printType(pp, v)
	} else if v {
		pp.WriteString("\nDestructuringAssignmentTarget: nil")
	}

	if f.Initializer != nil {
		pp.WriteString("\nInitializer: ")
		f.Initializer.printType(pp, v)
	} else if v {
		pp.WriteString("\nInitializer: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *BindingElement) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("BindingElement {")

	if f.SingleNameBinding != nil {
		pp.WriteString("\nSingleNameBinding: ")
		f.SingleNameBinding.printType(pp, v)
	} else if v {
		pp.WriteString("\nSingleNameBinding: nil")
	}

	if f.ArrayBindingPattern != nil {
		pp.WriteString("\nArrayBindingPattern: ")
		f.ArrayBindingPattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nArrayBindingPattern: nil")
	}

	if f.ObjectBindingPattern != nil {
		pp.WriteString("\nObjectBindingPattern: ")
		f.ObjectBindingPattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nObjectBindingPattern: nil")
	}

	if f.Initializer != nil {
		pp.WriteString("\nInitializer: ")
		f.Initializer.printType(pp, v)
	} else if v {
		pp.WriteString("\nInitializer: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *BindingProperty) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("BindingProperty {")

	pp.WriteString("\nPropertyName: ")
	f.PropertyName.printType(pp, v)

	pp.WriteString("\nBindingElement: ")
	f.BindingElement.printType(pp, v)

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *BitwiseANDExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("BitwiseANDExpression {")

	if f.BitwiseANDExpression != nil {
		pp.WriteString("\nBitwiseANDExpression: ")
		f.BitwiseANDExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nBitwiseANDExpression: nil")
	}

	pp.WriteString("\nEqualityExpression: ")
	f.EqualityExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *BitwiseORExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("BitwiseORExpression {")

	if f.BitwiseORExpression != nil {
		pp.WriteString("\nBitwiseORExpression: ")
		f.BitwiseORExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nBitwiseORExpression: nil")
	}

	pp.WriteString("\nBitwiseXORExpression: ")
	f.BitwiseXORExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *BitwiseXORExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("BitwiseXORExpression {")

	if f.BitwiseXORExpression != nil {
		pp.WriteString("\nBitwiseXORExpression: ")
		f.BitwiseXORExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nBitwiseXORExpression: nil")
	}

	pp.WriteString("\nBitwiseANDExpression: ")
	f.BitwiseANDExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Block) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Block {")

	if f.StatementList == nil {
		pp.WriteString("\nStatementList: nil")
	} else if len(f.StatementList) > 0 {
		pp.WriteString("\nStatementList: [")

		ipp := pp.Indent()

		for n, e := range f.StatementList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nStatementList: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *CallExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("CallExpression {")

	if f.MemberExpression != nil {
		pp.WriteString("\nMemberExpression: ")
		f.MemberExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nMemberExpression: nil")
	}

	if f.SuperCall || v {
		pp.Printf("\nSuperCall: %v", f.SuperCall)
	}

	if f.ImportCall != nil {
		pp.WriteString("\nImportCall: ")
		f.ImportCall.printType(pp, v)
	} else if v {
		pp.WriteString("\nImportCall: nil")
	}

	if f.CallExpression != nil {
		pp.WriteString("\nCallExpression: ")
		f.CallExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nCallExpression: nil")
	}

	if f.Arguments != nil {
		pp.WriteString("\nArguments: ")
		f.Arguments.printType(pp, v)
	} else if v {
		pp.WriteString("\nArguments: nil")
	}

	if f.Expression != nil {
		pp.WriteString("\nExpression: ")
		f.Expression.printType(pp, v)
	} else if v {
		pp.WriteString("\nExpression: nil")
	}

	if f.IdentifierName != nil {
		pp.WriteString("\nIdentifierName: ")
		f.IdentifierName.printType(pp, v)
	} else if v {
		pp.WriteString("\nIdentifierName: nil")
	}

	if f.TemplateLiteral != nil {
		pp.WriteString("\nTemplateLiteral: ")
		f.TemplateLiteral.printType(pp, v)
	} else if v {
		pp.WriteString("\nTemplateLiteral: nil")
	}

	if f.PrivateIdentifier != nil {
		pp.WriteString("\nPrivateIdentifier: ")
		f.PrivateIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nPrivateIdentifier: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *CaseClause) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("CaseClause {")

	pp.WriteString("\nExpression: ")
	f.Expression.printType(pp, v)

	if f.StatementList == nil {
		pp.WriteString("\nStatementList: nil")
	} else if len(f.StatementList) > 0 {
		pp.WriteString("\nStatementList: [")

		ipp := pp.Indent()

		for n, e := range f.StatementList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nStatementList: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ClassDeclaration) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ClassDeclaration {")

	if f.BindingIdentifier != nil {
		pp.WriteString("\nBindingIdentifier: ")
		f.BindingIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nBindingIdentifier: nil")
	}

	if f.ClassHeritage != nil {
		pp.WriteString("\nClassHeritage: ")
		f.ClassHeritage.printType(pp, v)
	} else if v {
		pp.WriteString("\nClassHeritage: nil")
	}

	if f.ClassBody == nil {
		pp.WriteString("\nClassBody: nil")
	} else if len(f.ClassBody) > 0 {
		pp.WriteString("\nClassBody: [")

		ipp := pp.Indent()

		for n, e := range f.ClassBody {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nClassBody: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ClassElement) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ClassElement {")

	if f.Static || v {
		pp.Printf("\nStatic: %v", f.Static)
	}

	if f.MethodDefinition != nil {
		pp.WriteString("\nMethodDefinition: ")
		f.MethodDefinition.printType(pp, v)
	} else if v {
		pp.WriteString("\nMethodDefinition: nil")
	}

	if f.FieldDefinition != nil {
		pp.WriteString("\nFieldDefinition: ")
		f.FieldDefinition.printType(pp, v)
	} else if v {
		pp.WriteString("\nFieldDefinition: nil")
	}

	if f.ClassStaticBlock != nil {
		pp.WriteString("\nClassStaticBlock: ")
		f.ClassStaticBlock.printType(pp, v)
	} else if v {
		pp.WriteString("\nClassStaticBlock: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ClassElementName) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ClassElementName {")

	if f.PropertyName != nil {
		pp.WriteString("\nPropertyName: ")
		f.PropertyName.printType(pp, v)
	} else if v {
		pp.WriteString("\nPropertyName: nil")
	}

	if f.PrivateIdentifier != nil {
		pp.WriteString("\nPrivateIdentifier: ")
		f.PrivateIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nPrivateIdentifier: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *CoalesceExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("CoalesceExpression {")

	if f.CoalesceExpressionHead != nil {
		pp.WriteString("\nCoalesceExpressionHead: ")
		f.CoalesceExpressionHead.printType(pp, v)
	} else if v {
		pp.WriteString("\nCoalesceExpressionHead: nil")
	}

	pp.WriteString("\nBitwiseORExpression: ")
	f.BitwiseORExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ConditionalExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ConditionalExpression {")

	if f.LogicalORExpression != nil {
		pp.WriteString("\nLogicalORExpression: ")
		f.LogicalORExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nLogicalORExpression: nil")
	}

	if f.CoalesceExpression != nil {
		pp.WriteString("\nCoalesceExpression: ")
		f.CoalesceExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nCoalesceExpression: nil")
	}

	if f.True != nil {
		pp.WriteString("\nTrue: ")
		f.True.printType(pp, v)
	} else if v {
		pp.WriteString("\nTrue: nil")
	}

	if f.False != nil {
		pp.WriteString("\nFalse: ")
		f.False.printType(pp, v)
	} else if v {
		pp.WriteString("\nFalse: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Declaration) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Declaration {")

	if f.ClassDeclaration != nil {
		pp.WriteString("\nClassDeclaration: ")
		f.ClassDeclaration.printType(pp, v)
	} else if v {
		pp.WriteString("\nClassDeclaration: nil")
	}

	if f.FunctionDeclaration != nil {
		pp.WriteString("\nFunctionDeclaration: ")
		f.FunctionDeclaration.printType(pp, v)
	} else if v {
		pp.WriteString("\nFunctionDeclaration: nil")
	}

	if f.LexicalDeclaration != nil {
		pp.WriteString("\nLexicalDeclaration: ")
		f.LexicalDeclaration.printType(pp, v)
	} else if v {
		pp.WriteString("\nLexicalDeclaration: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *DestructuringAssignmentTarget) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("DestructuringAssignmentTarget {")

	if f.LeftHandSideExpression != nil {
		pp.WriteString("\nLeftHandSideExpression: ")
		f.LeftHandSideExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nLeftHandSideExpression: nil")
	}

	if f.AssignmentPattern != nil {
		pp.WriteString("\nAssignmentPattern: ")
		f.AssignmentPattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nAssignmentPattern: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *EqualityExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("EqualityExpression {")

	if f.EqualityExpression != nil {
		pp.WriteString("\nEqualityExpression: ")
		f.EqualityExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nEqualityExpression: nil")
	}

	pp.WriteString("\nEqualityOperator: ")
	f.EqualityOperator.printType(pp, v)

	pp.WriteString("\nRelationalExpression: ")
	f.RelationalExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ExponentiationExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ExponentiationExpression {")

	if f.ExponentiationExpression != nil {
		pp.WriteString("\nExponentiationExpression: ")
		f.ExponentiationExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nExponentiationExpression: nil")
	}

	pp.WriteString("\nUnaryExpression: ")
	f.UnaryExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ExportClause) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ExportClause {")

	if f.ExportList == nil {
		pp.WriteString("\nExportList: nil")
	} else if len(f.ExportList) > 0 {
		pp.WriteString("\nExportList: [")

		ipp := pp.Indent()

		for n, e := range f.ExportList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nExportList: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ExportDeclaration) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ExportDeclaration {")

	if f.ExportClause != nil {
		pp.WriteString("\nExportClause: ")
		f.ExportClause.printType(pp, v)
	} else if v {
		pp.WriteString("\nExportClause: nil")
	}

	if f.ExportFromClause != nil {
		pp.WriteString("\nExportFromClause: ")
		f.ExportFromClause.printType(pp, v)
	} else if v {
		pp.WriteString("\nExportFromClause: nil")
	}

	if f.FromClause != nil {
		pp.WriteString("\nFromClause: ")
		f.FromClause.printType(pp, v)
	} else if v {
		pp.WriteString("\nFromClause: nil")
	}

	if f.VariableStatement != nil {
		pp.WriteString("\nVariableStatement: ")
		f.VariableStatement.printType(pp, v)
	} else if v {
		pp.WriteString("\nVariableStatement: nil")
	}

	if f.Declaration != nil {
		pp.WriteString("\nDeclaration: ")
		f.Declaration.printType(pp, v)
	} else if v {
		pp.WriteString("\nDeclaration: nil")
	}

	if f.DefaultFunction != nil {
		pp.WriteString("\nDefaultFunction: ")
		f.DefaultFunction.printType(pp, v)
	} else if v {
		pp.WriteString("\nDefaultFunction: nil")
	}

	if f.DefaultClass != nil {
		pp.WriteString("\nDefaultClass: ")
		f.DefaultClass.printType(pp, v)
	} else if v {
		pp.WriteString("\nDefaultClass: nil")
	}

	if f.DefaultAssignmentExpression != nil {
		pp.WriteString("\nDefaultAssignmentExpression: ")
		f.DefaultAssignmentExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nDefaultAssignmentExpression: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ExportSpecifier) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ExportSpecifier {")

	if f.IdentifierName != nil {
		pp.WriteString("\nIdentifierName: ")
		f.IdentifierName.printType(pp, v)
	} else if v {
		pp.WriteString("\nIdentifierName: nil")
	}

	if f.EIdentifierName != nil {
		pp.WriteString("\nEIdentifierName: ")
		f.EIdentifierName.printType(pp, v)
	} else if v {
		pp.WriteString("\nEIdentifierName: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Expression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Expression {")

	if f.Expressions == nil {
		pp.WriteString("\nExpressions: nil")
	} else if len(f.Expressions) > 0 {
		pp.WriteString("\nExpressions: [")

		ipp := pp.Indent()

		for n, e := range f.Expressions {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nExpressions: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *FieldDefinition) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("FieldDefinition {")

	pp.WriteString("\nClassElementName: ")
	f.ClassElementName.printType(pp, v)

	if f.Initializer != nil {
		pp.WriteString("\nInitializer: ")
		f.Initializer.printType(pp, v)
	} else if v {
		pp.WriteString("\nInitializer: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *FormalParameters) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("FormalParameters {")

	if f.FormalParameterList == nil {
		pp.WriteString("\nFormalParameterList: nil")
	} else if len(f.FormalParameterList) > 0 {
		pp.WriteString("\nFormalParameterList: [")

		ipp := pp.Indent()

		for n, e := range f.FormalParameterList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nFormalParameterList: []")
	}

	if f.BindingIdentifier != nil {
		pp.WriteString("\nBindingIdentifier: ")
		f.BindingIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nBindingIdentifier: nil")
	}

	if f.ArrayBindingPattern != nil {
		pp.WriteString("\nArrayBindingPattern: ")
		f.ArrayBindingPattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nArrayBindingPattern: nil")
	}

	if f.ObjectBindingPattern != nil {
		pp.WriteString("\nObjectBindingPattern: ")
		f.ObjectBindingPattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nObjectBindingPattern: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *FromClause) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("FromClause {")

	if f.ModuleSpecifier != nil {
		pp.WriteString("\nModuleSpecifier: ")
		f.ModuleSpecifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nModuleSpecifier: nil")
	}

	pp.WriteString("\nComments: ")
	f.Comments.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *FunctionDeclaration) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("FunctionDeclaration {")

	pp.WriteString("\nType: ")
	f.Type.printType(pp, v)

	if f.BindingIdentifier != nil {
		pp.WriteString("\nBindingIdentifier: ")
		f.BindingIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nBindingIdentifier: nil")
	}

	pp.WriteString("\nFormalParameters: ")
	f.FormalParameters.printType(pp, v)

	pp.WriteString("\nFunctionBody: ")
	f.FunctionBody.printType(pp, v)

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *IfStatement) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("IfStatement {")

	pp.WriteString("\nExpression: ")
	f.Expression.printType(pp, v)

	pp.WriteString("\nStatement: ")
	f.Statement.printType(pp, v)

	if f.ElseStatement != nil {
		pp.WriteString("\nElseStatement: ")
		f.ElseStatement.printType(pp, v)
	} else if v {
		pp.WriteString("\nElseStatement: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ImportClause) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ImportClause {")

	if f.ImportedDefaultBinding != nil {
		pp.WriteString("\nImportedDefaultBinding: ")
		f.ImportedDefaultBinding.printType(pp, v)
	} else if v {
		pp.WriteString("\nImportedDefaultBinding: nil")
	}

	if f.NameSpaceImport != nil {
		pp.WriteString("\nNameSpaceImport: ")
		f.NameSpaceImport.printType(pp, v)
	} else if v {
		pp.WriteString("\nNameSpaceImport: nil")
	}

	if f.NamedImports != nil {
		pp.WriteString("\nNamedImports: ")
		f.NamedImports.printType(pp, v)
	} else if v {
		pp.WriteString("\nNamedImports: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ImportDeclaration) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ImportDeclaration {")

	if f.ImportClause != nil {
		pp.WriteString("\nImportClause: ")
		f.ImportClause.printType(pp, v)
	} else if v {
		pp.WriteString("\nImportClause: nil")
	}

	pp.WriteString("\nFromClause: ")
	f.FromClause.printType(pp, v)

	if f.WithClause != nil {
		pp.WriteString("\nWithClause: ")
		f.WithClause.printType(pp, v)
	} else if v {
		pp.WriteString("\nWithClause: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ImportSpecifier) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ImportSpecifier {")

	if f.IdentifierName != nil {
		pp.WriteString("\nIdentifierName: ")
		f.IdentifierName.printType(pp, v)
	} else if v {
		pp.WriteString("\nIdentifierName: nil")
	}

	if f.ImportedBinding != nil {
		pp.WriteString("\nImportedBinding: ")
		f.ImportedBinding.printType(pp, v)
	} else if v {
		pp.WriteString("\nImportedBinding: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *IterationStatementDo) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("IterationStatementDo {")

	pp.WriteString("\nStatement: ")
	f.Statement.printType(pp, v)

	pp.WriteString("\nExpression: ")
	f.Expression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *IterationStatementFor) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("IterationStatementFor {")

	pp.WriteString("\nType: ")
	f.Type.printType(pp, v)

	if f.InitExpression != nil {
		pp.WriteString("\nInitExpression: ")
		f.InitExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nInitExpression: nil")
	}

	if f.InitVar == nil {
		pp.WriteString("\nInitVar: nil")
	} else if len(f.InitVar) > 0 {
		pp.WriteString("\nInitVar: [")

		ipp := pp.Indent()

		for n, e := range f.InitVar {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nInitVar: []")
	}

	if f.InitLexical != nil {
		pp.WriteString("\nInitLexical: ")
		f.InitLexical.printType(pp, v)
	} else if v {
		pp.WriteString("\nInitLexical: nil")
	}

	if f.Conditional != nil {
		pp.WriteString("\nConditional: ")
		f.Conditional.printType(pp, v)
	} else if v {
		pp.WriteString("\nConditional: nil")
	}

	if f.Afterthought != nil {
		pp.WriteString("\nAfterthought: ")
		f.Afterthought.printType(pp, v)
	} else if v {
		pp.WriteString("\nAfterthought: nil")
	}

	if f.LeftHandSideExpression != nil {
		pp.WriteString("\nLeftHandSideExpression: ")
		f.LeftHandSideExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nLeftHandSideExpression: nil")
	}

	if f.ForBindingIdentifier != nil {
		pp.WriteString("\nForBindingIdentifier: ")
		f.ForBindingIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nForBindingIdentifier: nil")
	}

	if f.ForBindingPatternObject != nil {
		pp.WriteString("\nForBindingPatternObject: ")
		f.ForBindingPatternObject.printType(pp, v)
	} else if v {
		pp.WriteString("\nForBindingPatternObject: nil")
	}

	if f.ForBindingPatternArray != nil {
		pp.WriteString("\nForBindingPatternArray: ")
		f.ForBindingPatternArray.printType(pp, v)
	} else if v {
		pp.WriteString("\nForBindingPatternArray: nil")
	}

	if f.In != nil {
		pp.WriteString("\nIn: ")
		f.In.printType(pp, v)
	} else if v {
		pp.WriteString("\nIn: nil")
	}

	if f.Of != nil {
		pp.WriteString("\nOf: ")
		f.Of.printType(pp, v)
	} else if v {
		pp.WriteString("\nOf: nil")
	}

	pp.WriteString("\nStatement: ")
	f.Statement.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *IterationStatementWhile) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("IterationStatementWhile {")

	pp.WriteString("\nExpression: ")
	f.Expression.printType(pp, v)

	pp.WriteString("\nStatement: ")
	f.Statement.printType(pp, v)

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *LeftHandSideExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("LeftHandSideExpression {")

	if f.NewExpression != nil {
		pp.WriteString("\nNewExpression: ")
		f.NewExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nNewExpression: nil")
	}

	if f.CallExpression != nil {
		pp.WriteString("\nCallExpression: ")
		f.CallExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nCallExpression: nil")
	}

	if f.OptionalExpression != nil {
		pp.WriteString("\nOptionalExpression: ")
		f.OptionalExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nOptionalExpression: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *LexicalBinding) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("LexicalBinding {")

	if f.BindingIdentifier != nil {
		pp.WriteString("\nBindingIdentifier: ")
		f.BindingIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nBindingIdentifier: nil")
	}

	if f.ArrayBindingPattern != nil {
		pp.WriteString("\nArrayBindingPattern: ")
		f.ArrayBindingPattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nArrayBindingPattern: nil")
	}

	if f.ObjectBindingPattern != nil {
		pp.WriteString("\nObjectBindingPattern: ")
		f.ObjectBindingPattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nObjectBindingPattern: nil")
	}

	if f.Initializer != nil {
		pp.WriteString("\nInitializer: ")
		f.Initializer.printType(pp, v)
	} else if v {
		pp.WriteString("\nInitializer: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *LexicalDeclaration) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("LexicalDeclaration {")

	pp.WriteString("\nLetOrConst: ")
	f.LetOrConst.printType(pp, v)

	if f.BindingList == nil {
		pp.WriteString("\nBindingList: nil")
	} else if len(f.BindingList) > 0 {
		pp.WriteString("\nBindingList: [")

		ipp := pp.Indent()

		for n, e := range f.BindingList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nBindingList: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *LogicalANDExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("LogicalANDExpression {")

	if f.LogicalANDExpression != nil {
		pp.WriteString("\nLogicalANDExpression: ")
		f.LogicalANDExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nLogicalANDExpression: nil")
	}

	pp.WriteString("\nBitwiseORExpression: ")
	f.BitwiseORExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *LogicalORExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("LogicalORExpression {")

	if f.LogicalORExpression != nil {
		pp.WriteString("\nLogicalORExpression: ")
		f.LogicalORExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nLogicalORExpression: nil")
	}

	pp.WriteString("\nLogicalANDExpression: ")
	f.LogicalANDExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *MemberExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("MemberExpression {")

	if f.MemberExpression != nil {
		pp.WriteString("\nMemberExpression: ")
		f.MemberExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nMemberExpression: nil")
	}

	if f.PrimaryExpression != nil {
		pp.WriteString("\nPrimaryExpression: ")
		f.PrimaryExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nPrimaryExpression: nil")
	}

	if f.Expression != nil {
		pp.WriteString("\nExpression: ")
		f.Expression.printType(pp, v)
	} else if v {
		pp.WriteString("\nExpression: nil")
	}

	if f.IdentifierName != nil {
		pp.WriteString("\nIdentifierName: ")
		f.IdentifierName.printType(pp, v)
	} else if v {
		pp.WriteString("\nIdentifierName: nil")
	}

	if f.TemplateLiteral != nil {
		pp.WriteString("\nTemplateLiteral: ")
		f.TemplateLiteral.printType(pp, v)
	} else if v {
		pp.WriteString("\nTemplateLiteral: nil")
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
		pp.WriteString("\nArguments: ")
		f.Arguments.printType(pp, v)
	} else if v {
		pp.WriteString("\nArguments: nil")
	}

	if f.PrivateIdentifier != nil {
		pp.WriteString("\nPrivateIdentifier: ")
		f.PrivateIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nPrivateIdentifier: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *MethodDefinition) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("MethodDefinition {")

	pp.WriteString("\nType: ")
	f.Type.printType(pp, v)

	pp.WriteString("\nClassElementName: ")
	f.ClassElementName.printType(pp, v)

	pp.WriteString("\nParams: ")
	f.Params.printType(pp, v)

	pp.WriteString("\nFunctionBody: ")
	f.FunctionBody.printType(pp, v)

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Module) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Module {")

	if f.ModuleListItems == nil {
		pp.WriteString("\nModuleListItems: nil")
	} else if len(f.ModuleListItems) > 0 {
		pp.WriteString("\nModuleListItems: [")

		ipp := pp.Indent()

		for n, e := range f.ModuleListItems {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nModuleListItems: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ModuleItem) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ModuleItem {")

	if f.ImportDeclaration != nil {
		pp.WriteString("\nImportDeclaration: ")
		f.ImportDeclaration.printType(pp, v)
	} else if v {
		pp.WriteString("\nImportDeclaration: nil")
	}

	if f.StatementListItem != nil {
		pp.WriteString("\nStatementListItem: ")
		f.StatementListItem.printType(pp, v)
	} else if v {
		pp.WriteString("\nStatementListItem: nil")
	}

	if f.ExportDeclaration != nil {
		pp.WriteString("\nExportDeclaration: ")
		f.ExportDeclaration.printType(pp, v)
	} else if v {
		pp.WriteString("\nExportDeclaration: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *MultiplicativeExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("MultiplicativeExpression {")

	if f.MultiplicativeExpression != nil {
		pp.WriteString("\nMultiplicativeExpression: ")
		f.MultiplicativeExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nMultiplicativeExpression: nil")
	}

	pp.WriteString("\nMultiplicativeOperator: ")
	f.MultiplicativeOperator.printType(pp, v)

	pp.WriteString("\nExponentiationExpression: ")
	f.ExponentiationExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *NamedImports) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("NamedImports {")

	if f.ImportList == nil {
		pp.WriteString("\nImportList: nil")
	} else if len(f.ImportList) > 0 {
		pp.WriteString("\nImportList: [")

		ipp := pp.Indent()

		for n, e := range f.ImportList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nImportList: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *NewExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("NewExpression {")

	if f.News == nil {
		pp.WriteString("\nNews: nil")
	} else if len(f.News) > 0 {
		pp.WriteString("\nNews: [")

		ipp := pp.Indent()

		for n, e := range f.News {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nNews: []")
	}

	pp.WriteString("\nMemberExpression: ")
	f.MemberExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ObjectAssignmentPattern) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ObjectAssignmentPattern {")

	if f.AssignmentPropertyList == nil {
		pp.WriteString("\nAssignmentPropertyList: nil")
	} else if len(f.AssignmentPropertyList) > 0 {
		pp.WriteString("\nAssignmentPropertyList: [")

		ipp := pp.Indent()

		for n, e := range f.AssignmentPropertyList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nAssignmentPropertyList: []")
	}

	if f.AssignmentRestElement != nil {
		pp.WriteString("\nAssignmentRestElement: ")
		f.AssignmentRestElement.printType(pp, v)
	} else if v {
		pp.WriteString("\nAssignmentRestElement: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ObjectBindingPattern) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ObjectBindingPattern {")

	if f.BindingPropertyList == nil {
		pp.WriteString("\nBindingPropertyList: nil")
	} else if len(f.BindingPropertyList) > 0 {
		pp.WriteString("\nBindingPropertyList: [")

		ipp := pp.Indent()

		for n, e := range f.BindingPropertyList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nBindingPropertyList: []")
	}

	if f.BindingRestProperty != nil {
		pp.WriteString("\nBindingRestProperty: ")
		f.BindingRestProperty.printType(pp, v)
	} else if v {
		pp.WriteString("\nBindingRestProperty: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ObjectLiteral) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ObjectLiteral {")

	if f.PropertyDefinitionList == nil {
		pp.WriteString("\nPropertyDefinitionList: nil")
	} else if len(f.PropertyDefinitionList) > 0 {
		pp.WriteString("\nPropertyDefinitionList: [")

		ipp := pp.Indent()

		for n, e := range f.PropertyDefinitionList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nPropertyDefinitionList: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *OptionalChain) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("OptionalChain {")

	if f.OptionalChain != nil {
		pp.WriteString("\nOptionalChain: ")
		f.OptionalChain.printType(pp, v)
	} else if v {
		pp.WriteString("\nOptionalChain: nil")
	}

	if f.Arguments != nil {
		pp.WriteString("\nArguments: ")
		f.Arguments.printType(pp, v)
	} else if v {
		pp.WriteString("\nArguments: nil")
	}

	if f.Expression != nil {
		pp.WriteString("\nExpression: ")
		f.Expression.printType(pp, v)
	} else if v {
		pp.WriteString("\nExpression: nil")
	}

	if f.IdentifierName != nil {
		pp.WriteString("\nIdentifierName: ")
		f.IdentifierName.printType(pp, v)
	} else if v {
		pp.WriteString("\nIdentifierName: nil")
	}

	if f.TemplateLiteral != nil {
		pp.WriteString("\nTemplateLiteral: ")
		f.TemplateLiteral.printType(pp, v)
	} else if v {
		pp.WriteString("\nTemplateLiteral: nil")
	}

	if f.PrivateIdentifier != nil {
		pp.WriteString("\nPrivateIdentifier: ")
		f.PrivateIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nPrivateIdentifier: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *OptionalExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("OptionalExpression {")

	if f.MemberExpression != nil {
		pp.WriteString("\nMemberExpression: ")
		f.MemberExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nMemberExpression: nil")
	}

	if f.CallExpression != nil {
		pp.WriteString("\nCallExpression: ")
		f.CallExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nCallExpression: nil")
	}

	if f.OptionalExpression != nil {
		pp.WriteString("\nOptionalExpression: ")
		f.OptionalExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nOptionalExpression: nil")
	}

	pp.WriteString("\nOptionalChain: ")
	f.OptionalChain.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ParenthesizedExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ParenthesizedExpression {")

	if f.Expressions == nil {
		pp.WriteString("\nExpressions: nil")
	} else if len(f.Expressions) > 0 {
		pp.WriteString("\nExpressions: [")

		ipp := pp.Indent()

		for n, e := range f.Expressions {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nExpressions: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *PrimaryExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("PrimaryExpression {")

	if f.This != nil {
		pp.WriteString("\nThis: ")
		f.This.printType(pp, v)
	} else if v {
		pp.WriteString("\nThis: nil")
	}

	if f.IdentifierReference != nil {
		pp.WriteString("\nIdentifierReference: ")
		f.IdentifierReference.printType(pp, v)
	} else if v {
		pp.WriteString("\nIdentifierReference: nil")
	}

	if f.Literal != nil {
		pp.WriteString("\nLiteral: ")
		f.Literal.printType(pp, v)
	} else if v {
		pp.WriteString("\nLiteral: nil")
	}

	if f.ArrayLiteral != nil {
		pp.WriteString("\nArrayLiteral: ")
		f.ArrayLiteral.printType(pp, v)
	} else if v {
		pp.WriteString("\nArrayLiteral: nil")
	}

	if f.ObjectLiteral != nil {
		pp.WriteString("\nObjectLiteral: ")
		f.ObjectLiteral.printType(pp, v)
	} else if v {
		pp.WriteString("\nObjectLiteral: nil")
	}

	if f.FunctionExpression != nil {
		pp.WriteString("\nFunctionExpression: ")
		f.FunctionExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nFunctionExpression: nil")
	}

	if f.ClassExpression != nil {
		pp.WriteString("\nClassExpression: ")
		f.ClassExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nClassExpression: nil")
	}

	if f.TemplateLiteral != nil {
		pp.WriteString("\nTemplateLiteral: ")
		f.TemplateLiteral.printType(pp, v)
	} else if v {
		pp.WriteString("\nTemplateLiteral: nil")
	}

	if f.ParenthesizedExpression != nil {
		pp.WriteString("\nParenthesizedExpression: ")
		f.ParenthesizedExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nParenthesizedExpression: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *PropertyDefinition) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("PropertyDefinition {")

	if f.IsCoverInitializedName || v {
		pp.Printf("\nIsCoverInitializedName: %v", f.IsCoverInitializedName)
	}

	if f.PropertyName != nil {
		pp.WriteString("\nPropertyName: ")
		f.PropertyName.printType(pp, v)
	} else if v {
		pp.WriteString("\nPropertyName: nil")
	}

	if f.AssignmentExpression != nil {
		pp.WriteString("\nAssignmentExpression: ")
		f.AssignmentExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nAssignmentExpression: nil")
	}

	if f.MethodDefinition != nil {
		pp.WriteString("\nMethodDefinition: ")
		f.MethodDefinition.printType(pp, v)
	} else if v {
		pp.WriteString("\nMethodDefinition: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *PropertyName) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("PropertyName {")

	if f.LiteralPropertyName != nil {
		pp.WriteString("\nLiteralPropertyName: ")
		f.LiteralPropertyName.printType(pp, v)
	} else if v {
		pp.WriteString("\nLiteralPropertyName: nil")
	}

	if f.ComputedPropertyName != nil {
		pp.WriteString("\nComputedPropertyName: ")
		f.ComputedPropertyName.printType(pp, v)
	} else if v {
		pp.WriteString("\nComputedPropertyName: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *RelationalExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("RelationalExpression {")

	if f.PrivateIdentifier != nil {
		pp.WriteString("\nPrivateIdentifier: ")
		f.PrivateIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nPrivateIdentifier: nil")
	}

	if f.RelationalExpression != nil {
		pp.WriteString("\nRelationalExpression: ")
		f.RelationalExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nRelationalExpression: nil")
	}

	pp.WriteString("\nRelationshipOperator: ")
	f.RelationshipOperator.printType(pp, v)

	pp.WriteString("\nShiftExpression: ")
	f.ShiftExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Script) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Script {")

	if f.StatementList == nil {
		pp.WriteString("\nStatementList: nil")
	} else if len(f.StatementList) > 0 {
		pp.WriteString("\nStatementList: [")

		ipp := pp.Indent()

		for n, e := range f.StatementList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nStatementList: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ShiftExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ShiftExpression {")

	if f.ShiftExpression != nil {
		pp.WriteString("\nShiftExpression: ")
		f.ShiftExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nShiftExpression: nil")
	}

	pp.WriteString("\nShiftOperator: ")
	f.ShiftOperator.printType(pp, v)

	pp.WriteString("\nAdditiveExpression: ")
	f.AdditiveExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Statement) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Statement {")

	pp.WriteString("\nType: ")
	f.Type.printType(pp, v)

	if f.BlockStatement != nil {
		pp.WriteString("\nBlockStatement: ")
		f.BlockStatement.printType(pp, v)
	} else if v {
		pp.WriteString("\nBlockStatement: nil")
	}

	if f.VariableStatement != nil {
		pp.WriteString("\nVariableStatement: ")
		f.VariableStatement.printType(pp, v)
	} else if v {
		pp.WriteString("\nVariableStatement: nil")
	}

	if f.ExpressionStatement != nil {
		pp.WriteString("\nExpressionStatement: ")
		f.ExpressionStatement.printType(pp, v)
	} else if v {
		pp.WriteString("\nExpressionStatement: nil")
	}

	if f.IfStatement != nil {
		pp.WriteString("\nIfStatement: ")
		f.IfStatement.printType(pp, v)
	} else if v {
		pp.WriteString("\nIfStatement: nil")
	}

	if f.IterationStatementDo != nil {
		pp.WriteString("\nIterationStatementDo: ")
		f.IterationStatementDo.printType(pp, v)
	} else if v {
		pp.WriteString("\nIterationStatementDo: nil")
	}

	if f.IterationStatementWhile != nil {
		pp.WriteString("\nIterationStatementWhile: ")
		f.IterationStatementWhile.printType(pp, v)
	} else if v {
		pp.WriteString("\nIterationStatementWhile: nil")
	}

	if f.IterationStatementFor != nil {
		pp.WriteString("\nIterationStatementFor: ")
		f.IterationStatementFor.printType(pp, v)
	} else if v {
		pp.WriteString("\nIterationStatementFor: nil")
	}

	if f.SwitchStatement != nil {
		pp.WriteString("\nSwitchStatement: ")
		f.SwitchStatement.printType(pp, v)
	} else if v {
		pp.WriteString("\nSwitchStatement: nil")
	}

	if f.WithStatement != nil {
		pp.WriteString("\nWithStatement: ")
		f.WithStatement.printType(pp, v)
	} else if v {
		pp.WriteString("\nWithStatement: nil")
	}

	if f.LabelIdentifier != nil {
		pp.WriteString("\nLabelIdentifier: ")
		f.LabelIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nLabelIdentifier: nil")
	}

	if f.LabelledItemFunction != nil {
		pp.WriteString("\nLabelledItemFunction: ")
		f.LabelledItemFunction.printType(pp, v)
	} else if v {
		pp.WriteString("\nLabelledItemFunction: nil")
	}

	if f.LabelledItemStatement != nil {
		pp.WriteString("\nLabelledItemStatement: ")
		f.LabelledItemStatement.printType(pp, v)
	} else if v {
		pp.WriteString("\nLabelledItemStatement: nil")
	}

	if f.TryStatement != nil {
		pp.WriteString("\nTryStatement: ")
		f.TryStatement.printType(pp, v)
	} else if v {
		pp.WriteString("\nTryStatement: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *StatementListItem) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("StatementListItem {")

	if f.Statement != nil {
		pp.WriteString("\nStatement: ")
		f.Statement.printType(pp, v)
	} else if v {
		pp.WriteString("\nStatement: nil")
	}

	if f.Declaration != nil {
		pp.WriteString("\nDeclaration: ")
		f.Declaration.printType(pp, v)
	} else if v {
		pp.WriteString("\nDeclaration: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *SwitchStatement) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("SwitchStatement {")

	pp.WriteString("\nExpression: ")
	f.Expression.printType(pp, v)

	if f.CaseClauses == nil {
		pp.WriteString("\nCaseClauses: nil")
	} else if len(f.CaseClauses) > 0 {
		pp.WriteString("\nCaseClauses: [")

		ipp := pp.Indent()

		for n, e := range f.CaseClauses {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nCaseClauses: []")
	}

	if f.DefaultClause == nil {
		pp.WriteString("\nDefaultClause: nil")
	} else if len(f.DefaultClause) > 0 {
		pp.WriteString("\nDefaultClause: [")

		ipp := pp.Indent()

		for n, e := range f.DefaultClause {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nDefaultClause: []")
	}

	if f.PostDefaultCaseClauses == nil {
		pp.WriteString("\nPostDefaultCaseClauses: nil")
	} else if len(f.PostDefaultCaseClauses) > 0 {
		pp.WriteString("\nPostDefaultCaseClauses: [")

		ipp := pp.Indent()

		for n, e := range f.PostDefaultCaseClauses {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nPostDefaultCaseClauses: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *TemplateLiteral) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("TemplateLiteral {")

	if f.NoSubstitutionTemplate != nil {
		pp.WriteString("\nNoSubstitutionTemplate: ")
		f.NoSubstitutionTemplate.printType(pp, v)
	} else if v {
		pp.WriteString("\nNoSubstitutionTemplate: nil")
	}

	if f.TemplateHead != nil {
		pp.WriteString("\nTemplateHead: ")
		f.TemplateHead.printType(pp, v)
	} else if v {
		pp.WriteString("\nTemplateHead: nil")
	}

	if f.Expressions == nil {
		pp.WriteString("\nExpressions: nil")
	} else if len(f.Expressions) > 0 {
		pp.WriteString("\nExpressions: [")

		ipp := pp.Indent()

		for n, e := range f.Expressions {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nExpressions: []")
	}

	if f.TemplateMiddleList == nil {
		pp.WriteString("\nTemplateMiddleList: nil")
	} else if len(f.TemplateMiddleList) > 0 {
		pp.WriteString("\nTemplateMiddleList: [")

		ipp := pp.Indent()

		for n, e := range f.TemplateMiddleList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nTemplateMiddleList: []")
	}

	if f.TemplateTail != nil {
		pp.WriteString("\nTemplateTail: ")
		f.TemplateTail.printType(pp, v)
	} else if v {
		pp.WriteString("\nTemplateTail: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *TryStatement) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("TryStatement {")

	pp.WriteString("\nTryBlock: ")
	f.TryBlock.printType(pp, v)

	if f.CatchParameterBindingIdentifier != nil {
		pp.WriteString("\nCatchParameterBindingIdentifier: ")
		f.CatchParameterBindingIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nCatchParameterBindingIdentifier: nil")
	}

	if f.CatchParameterObjectBindingPattern != nil {
		pp.WriteString("\nCatchParameterObjectBindingPattern: ")
		f.CatchParameterObjectBindingPattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nCatchParameterObjectBindingPattern: nil")
	}

	if f.CatchParameterArrayBindingPattern != nil {
		pp.WriteString("\nCatchParameterArrayBindingPattern: ")
		f.CatchParameterArrayBindingPattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nCatchParameterArrayBindingPattern: nil")
	}

	if f.CatchBlock != nil {
		pp.WriteString("\nCatchBlock: ")
		f.CatchBlock.printType(pp, v)
	} else if v {
		pp.WriteString("\nCatchBlock: nil")
	}

	if f.FinallyBlock != nil {
		pp.WriteString("\nFinallyBlock: ")
		f.FinallyBlock.printType(pp, v)
	} else if v {
		pp.WriteString("\nFinallyBlock: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *UnaryExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("UnaryExpression {")

	if f.UnaryOperators == nil {
		pp.WriteString("\nUnaryOperators: nil")
	} else if len(f.UnaryOperators) > 0 {
		pp.WriteString("\nUnaryOperators: [")

		ipp := pp.Indent()

		for n, e := range f.UnaryOperators {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nUnaryOperators: []")
	}

	pp.WriteString("\nUpdateExpression: ")
	f.UpdateExpression.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *UnaryOperatorComments) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("UnaryOperatorComments {")

	pp.WriteString("\nUnaryOperator: ")
	f.UnaryOperator.printType(pp, v)

	pp.WriteString("\nComments: ")
	f.Comments.printType(pp, v)

	w.WriteString("\n}")
}

func (f *UpdateExpression) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("UpdateExpression {")

	if f.LeftHandSideExpression != nil {
		pp.WriteString("\nLeftHandSideExpression: ")
		f.LeftHandSideExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nLeftHandSideExpression: nil")
	}

	pp.WriteString("\nUpdateOperator: ")
	f.UpdateOperator.printType(pp, v)

	if f.UnaryExpression != nil {
		pp.WriteString("\nUnaryExpression: ")
		f.UnaryExpression.printType(pp, v)
	} else if v {
		pp.WriteString("\nUnaryExpression: nil")
	}

	pp.WriteString("\nComments: ")
	f.Comments.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *VariableStatement) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("VariableStatement {")

	if f.VariableDeclarationList == nil {
		pp.WriteString("\nVariableDeclarationList: nil")
	} else if len(f.VariableDeclarationList) > 0 {
		pp.WriteString("\nVariableDeclarationList: [")

		ipp := pp.Indent()

		for n, e := range f.VariableDeclarationList {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nVariableDeclarationList: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *WithClause) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("WithClause {")

	if f.WithEntries == nil {
		pp.WriteString("\nWithEntries: nil")
	} else if len(f.WithEntries) > 0 {
		pp.WriteString("\nWithEntries: [")

		ipp := pp.Indent()

		for n, e := range f.WithEntries {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nWithEntries: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *WithEntry) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("WithEntry {")

	if f.AttributeKey != nil {
		pp.WriteString("\nAttributeKey: ")
		f.AttributeKey.printType(pp, v)
	} else if v {
		pp.WriteString("\nAttributeKey: nil")
	}

	if f.Value != nil {
		pp.WriteString("\nValue: ")
		f.Value.printType(pp, v)
	} else if v {
		pp.WriteString("\nValue: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *WithStatement) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("WithStatement {")

	pp.WriteString("\nExpression: ")
	f.Expression.printType(pp, v)

	pp.WriteString("\nStatement: ")
	f.Statement.printType(pp, v)

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}
