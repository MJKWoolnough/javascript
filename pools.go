package javascript

import "sync"

var (
	poolClassDeclaration = sync.Pool{
		New: func() interface{} {
			return new(ClassDeclaration)
		},
	}
	poolMethodDefinition = sync.Pool{
		New: func() interface{} {
			return new(MethodDefinition)
		},
	}
	poolPropertyName = sync.Pool{
		New: func() interface{} {
			return new(PropertyName)
		},
	}
	poolConditionalExpression = sync.Pool{
		New: func() interface{} {
			return new(ConditionalExpression)
		},
	}
	poolLogicalORExpression = sync.Pool{
		New: func() interface{} {
			return new(LogicalORExpression)
		},
	}
	poolLogicalANDExpression = sync.Pool{
		New: func() interface{} {
			return new(LogicalANDExpression)
		},
	}
	poolBitwiseORExpression = sync.Pool{
		New: func() interface{} {
			return new(BitwiseORExpression)
		},
	}
	poolBitwiseXORExpression = sync.Pool{
		New: func() interface{} {
			return new(BitwiseXORExpression)
		},
	}
	poolBitwiseANDExpression = sync.Pool{
		New: func() interface{} {
			return new(BitwiseANDExpression)
		},
	}
	poolEqualityExpression = sync.Pool{
		New: func() interface{} {
			return new(EqualityExpression)
		},
	}
	poolRelationalExpression = sync.Pool{
		New: func() interface{} {
			return new(RelationalExpression)
		},
	}
	poolShiftExpression = sync.Pool{
		New: func() interface{} {
			return new(ShiftExpression)
		},
	}
	poolAdditiveExpression = sync.Pool{
		New: func() interface{} {
			return new(AdditiveExpression)
		},
	}
	poolMultiplicativeExpression = sync.Pool{
		New: func() interface{} {
			return new(MultiplicativeExpression)
		},
	}
	poolExponentiationExpression = sync.Pool{
		New: func() interface{} {
			return new(ExponentiationExpression)
		},
	}
	poolUnaryExpression = sync.Pool{
		New: func() interface{} {
			return new(UnaryExpression)
		},
	}
	poolUpdateExpression = sync.Pool{
		New: func() interface{} {
			return new(UpdateExpression)
		},
	}
	poolAssignmentExpression = sync.Pool{
		New: func() interface{} {
			return new(AssignmentExpression)
		},
	}
	poolLeftHandSideExpression = sync.Pool{
		New: func() interface{} {
			return new(LeftHandSideExpression)
		},
	}
	poolExpression = sync.Pool{
		New: func() interface{} {
			return new(Expression)
		},
	}
	poolNewExpression = sync.Pool{
		New: func() interface{} {
			return new(NewExpression)
		},
	}
	poolMemberExpression = sync.Pool{
		New: func() interface{} {
			return new(MemberExpression)
		},
	}
	poolPrimaryExpression = sync.Pool{
		New: func() interface{} {
			return new(PrimaryExpression)
		},
	}
	poolArguments = sync.Pool{
		New: func() interface{} {
			return new(Arguments)
		},
	}
	poolCallExpression = sync.Pool{
		New: func() interface{} {
			return new(CallExpression)
		},
	}
	poolFunctionDeclaration = sync.Pool{
		New: func() interface{} {
			return new(FunctionDeclaration)
		},
	}
	poolFormalParameters = sync.Pool{
		New: func() interface{} {
			return new(FormalParameters)
		},
	}
	poolBindingElement = sync.Pool{
		New: func() interface{} {
			return new(BindingElement)
		},
	}
	poolFunctionRestParameter = sync.Pool{
		New: func() interface{} {
			return new(FunctionRestParameter)
		},
	}
	poolScript = sync.Pool{
		New: func() interface{} {
			return new(Script)
		},
	}
	poolDeclaration = sync.Pool{
		New: func() interface{} {
			return new(Declaration)
		},
	}
	poolLexicalDeclaration = sync.Pool{
		New: func() interface{} {
			return new(LexicalDeclaration)
		},
	}
	poolLexicalBinding = sync.Pool{
		New: func() interface{} {
			return new(LexicalBinding)
		},
	}
	poolArrayBindingPattern = sync.Pool{
		New: func() interface{} {
			return new(ArrayBindingPattern)
		},
	}
	poolObjectBindingPattern = sync.Pool{
		New: func() interface{} {
			return new(ObjectBindingPattern)
		},
	}
	poolBindingProperty = sync.Pool{
		New: func() interface{} {
			return new(BindingProperty)
		},
	}
	poolVariableDeclaration = sync.Pool{
		New: func() interface{} {
			return new(VariableDeclaration)
		},
	}
	poolArrayLiteral = sync.Pool{
		New: func() interface{} {
			return new(ArrayLiteral)
		},
	}
	poolObjectLiteral = sync.Pool{
		New: func() interface{} {
			return new(ObjectLiteral)
		},
	}
	poolPropertyDefinition = sync.Pool{
		New: func() interface{} {
			return new(PropertyDefinition)
		},
	}
	poolTemplateLiteral = sync.Pool{
		New: func() interface{} {
			return new(TemplateLiteral)
		},
	}
	poolArrowFunction = sync.Pool{
		New: func() interface{} {
			return new(ArrowFunction)
		},
	}
	poolModule = sync.Pool{
		New: func() interface{} {
			return new(Module)
		},
	}
	poolModuleListItem = sync.Pool{
		New: func() interface{} {
			return new(ModuleListItem)
		},
	}
	poolImportDeclaration = sync.Pool{
		New: func() interface{} {
			return new(ImportDeclaration)
		},
	}
	poolImportClause = sync.Pool{
		New: func() interface{} {
			return new(ImportClause)
		},
	}
	poolFromClause = sync.Pool{
		New: func() interface{} {
			return new(FromClause)
		},
	}
	poolNamedImports = sync.Pool{
		New: func() interface{} {
			return new(NamedImports)
		},
	}
	poolImportSpecifier = sync.Pool{
		New: func() interface{} {
			return new(ImportSpecifier)
		},
	}
	poolExportDeclaration = sync.Pool{
		New: func() interface{} {
			return new(ExportDeclaration)
		},
	}
	poolExportClause = sync.Pool{
		New: func() interface{} {
			return new(ExportClause)
		},
	}
	poolExportSpecifier = sync.Pool{
		New: func() interface{} {
			return new(ExportSpecifier)
		},
	}
	poolBlock = sync.Pool{
		New: func() interface{} {
			return new(Block)
		},
	}
	poolStatementListItem = sync.Pool{
		New: func() interface{} {
			return new(StatementListItem)
		},
	}
	poolStatement = sync.Pool{
		New: func() interface{} {
			return new(Statement)
		},
	}
	poolIfStatement = sync.Pool{
		New: func() interface{} {
			return new(IfStatement)
		},
	}
	poolIterationStatementDo = sync.Pool{
		New: func() interface{} {
			return new(IterationStatementDo)
		},
	}
	poolIterationStatementWhile = sync.Pool{
		New: func() interface{} {
			return new(IterationStatementWhile)
		},
	}
	poolIterationStatementFor = sync.Pool{
		New: func() interface{} {
			return new(IterationStatementFor)
		},
	}
	poolSwitchStatement = sync.Pool{
		New: func() interface{} {
			return new(SwitchStatement)
		},
	}
	poolCaseClause = sync.Pool{
		New: func() interface{} {
			return new(CaseClause)
		},
	}
	poolWithStatement = sync.Pool{
		New: func() interface{} {
			return new(WithStatement)
		},
	}
	poolTryStatement = sync.Pool{
		New: func() interface{} {
			return new(TryStatement)
		},
	}
	poolVariableStatement = sync.Pool{
		New: func() interface{} {
			return new(VariableStatement)
		},
	}
	poolCoverParenthesizedExpressionAndArrowParameterList = sync.Pool{
		New: func() interface{} {
			return new(CoverParenthesizedExpressionAndArrowParameterList)
		},
	}
)

func newClassDeclaration() *ClassDeclaration { return poolClassDeclaration.Get().(*ClassDeclaration) }
func newMethodDefinition() *MethodDefinition { return poolMethodDefinition.Get().(*MethodDefinition) }
func newPropertyName() *PropertyName         { return poolPropertyName.Get().(*PropertyName) }
func newConditionalExpression() *ConditionalExpression {
	return poolConditionalExpression.Get().(*ConditionalExpression)
}
func newLogicalORExpression() *LogicalORExpression {
	return poolLogicalORExpression.Get().(*LogicalORExpression)
}
func newLogicalANDExpression() *LogicalANDExpression {
	return poolLogicalANDExpression.Get().(*LogicalANDExpression)
}
func newBitwiseORExpression() *BitwiseORExpression {
	return poolBitwiseORExpression.Get().(*BitwiseORExpression)
}
func newBitwiseXORExpression() *BitwiseXORExpression {
	return poolBitwiseXORExpression.Get().(*BitwiseXORExpression)
}
func newBitwiseANDExpression() *BitwiseANDExpression {
	return poolBitwiseANDExpression.Get().(*BitwiseANDExpression)
}
func newEqualityExpression() *EqualityExpression {
	return poolEqualityExpression.Get().(*EqualityExpression)
}
func newRelationalExpression() *RelationalExpression {
	return poolRelationalExpression.Get().(*RelationalExpression)
}
func newShiftExpression() *ShiftExpression { return poolShiftExpression.Get().(*ShiftExpression) }
func newAdditiveExpression() *AdditiveExpression {
	return poolAdditiveExpression.Get().(*AdditiveExpression)
}
func newMultiplicativeExpression() *MultiplicativeExpression {
	return poolMultiplicativeExpression.Get().(*MultiplicativeExpression)
}
func newExponentiationExpression() *ExponentiationExpression {
	return poolExponentiationExpression.Get().(*ExponentiationExpression)
}
func newUnaryExpression() *UnaryExpression   { return poolUnaryExpression.Get().(*UnaryExpression) }
func newUpdateExpression() *UpdateExpression { return poolUpdateExpression.Get().(*UpdateExpression) }
func newAssignmentExpression() *AssignmentExpression {
	return poolAssignmentExpression.Get().(*AssignmentExpression)
}
func newLeftHandSideExpression() *LeftHandSideExpression {
	return poolLeftHandSideExpression.Get().(*LeftHandSideExpression)
}
func newExpression() *Expression             { return poolExpression.Get().(*Expression) }
func newNewExpression() *NewExpression       { return poolNewExpression.Get().(*NewExpression) }
func newMemberExpression() *MemberExpression { return poolMemberExpression.Get().(*MemberExpression) }
func newPrimaryExpression() *PrimaryExpression {
	return poolPrimaryExpression.Get().(*PrimaryExpression)
}
func newArguments() *Arguments           { return poolArguments.Get().(*Arguments) }
func newCallExpression() *CallExpression { return poolCallExpression.Get().(*CallExpression) }
func newFunctionDeclaration() *FunctionDeclaration {
	return poolFunctionDeclaration.Get().(*FunctionDeclaration)
}
func newFormalParameters() *FormalParameters { return poolFormalParameters.Get().(*FormalParameters) }
func newBindingElement() *BindingElement     { return poolBindingElement.Get().(*BindingElement) }
func newFunctionRestParameter() *FunctionRestParameter {
	return poolFunctionRestParameter.Get().(*FunctionRestParameter)
}
func newScript() *Script           { return poolScript.Get().(*Script) }
func newDeclaration() *Declaration { return poolDeclaration.Get().(*Declaration) }
func newLexicalDeclaration() *LexicalDeclaration {
	return poolLexicalDeclaration.Get().(*LexicalDeclaration)
}
func newLexicalBinding() *LexicalBinding { return poolLexicalBinding.Get().(*LexicalBinding) }
func newArrayBindingPattern() *ArrayBindingPattern {
	return poolArrayBindingPattern.Get().(*ArrayBindingPattern)
}
func newObjectBindingPattern() *ObjectBindingPattern {
	return poolObjectBindingPattern.Get().(*ObjectBindingPattern)
}
func newBindingProperty() *BindingProperty { return poolBindingProperty.Get().(*BindingProperty) }
func newVariableDeclaration() *VariableDeclaration {
	return poolVariableDeclaration.Get().(*VariableDeclaration)
}
func newArrayLiteral() *ArrayLiteral   { return poolArrayLiteral.Get().(*ArrayLiteral) }
func newObjectLiteral() *ObjectLiteral { return poolObjectLiteral.Get().(*ObjectLiteral) }
func newPropertyDefinition() *PropertyDefinition {
	return poolPropertyDefinition.Get().(*PropertyDefinition)
}
func newTemplateLiteral() *TemplateLiteral { return poolTemplateLiteral.Get().(*TemplateLiteral) }
func newArrowFunction() *ArrowFunction     { return poolArrowFunction.Get().(*ArrowFunction) }
func newModule() *Module                   { return poolModule.Get().(*Module) }
func newModuleListItem() *ModuleListItem   { return poolModuleListItem.Get().(*ModuleListItem) }
func newImportDeclaration() *ImportDeclaration {
	return poolImportDeclaration.Get().(*ImportDeclaration)
}
func newImportClause() *ImportClause       { return poolImportClause.Get().(*ImportClause) }
func newFromClause() *FromClause           { return poolFromClause.Get().(*FromClause) }
func newNamedImports() *NamedImports       { return poolNamedImports.Get().(*NamedImports) }
func newImportSpecifier() *ImportSpecifier { return poolImportSpecifier.Get().(*ImportSpecifier) }
func newExportDeclaration() *ExportDeclaration {
	return poolExportDeclaration.Get().(*ExportDeclaration)
}
func newExportClause() *ExportClause       { return poolExportClause.Get().(*ExportClause) }
func newExportSpecifier() *ExportSpecifier { return poolExportSpecifier.Get().(*ExportSpecifier) }
func newBlock() *Block                     { return poolBlock.Get().(*Block) }
func newStatementListItem() *StatementListItem {
	return poolStatementListItem.Get().(*StatementListItem)
}
func newStatement() *Statement     { return poolStatement.Get().(*Statement) }
func newIfStatement() *IfStatement { return poolIfStatement.Get().(*IfStatement) }
func newIterationStatementDo() *IterationStatementDo {
	return poolIterationStatementDo.Get().(*IterationStatementDo)
}
func newIterationStatementWhile() *IterationStatementWhile {
	return poolIterationStatementWhile.Get().(*IterationStatementWhile)
}
func newIterationStatementFor() *IterationStatementFor {
	return poolIterationStatementFor.Get().(*IterationStatementFor)
}
func newSwitchStatement() *SwitchStatement { return poolSwitchStatement.Get().(*SwitchStatement) }
func newCaseClause() *CaseClause           { return poolCaseClause.Get().(*CaseClause) }
func newWithStatement() *WithStatement     { return poolWithStatement.Get().(*WithStatement) }
func newTryStatement() *TryStatement       { return poolTryStatement.Get().(*TryStatement) }
func newVariableStatement() *VariableStatement {
	return poolVariableStatement.Get().(*VariableStatement)
}
func newCoverParenthesizedExpressionAndArrowParameterList() *CoverParenthesizedExpressionAndArrowParameterList {
	return poolCoverParenthesizedExpressionAndArrowParameterList.Get().(*CoverParenthesizedExpressionAndArrowParameterList)
}

func (t *ClassDeclaration) clear() {
	t.BindingIdentifier = nil
	if t.ClassHeritage != nil {
		t.ClassHeritage.clear()
		poolLeftHandSideExpression.Put(t.ClassHeritage)

	}
	for _, e := range t.ClassBody {
		e.clear()
	}
	*t = ClassDeclaration{}
}

func (t *MethodDefinition) clear() {
	t.PropertyName.clear()
	t.Params.clear()
	t.FunctionBody.clear()
	*t = MethodDefinition{}
}

func (t *PropertyName) clear() {
	t.LiteralPropertyName = nil
	if t.ComputedPropertyName != nil {
		t.ComputedPropertyName.clear()
		poolAssignmentExpression.Put(t.ComputedPropertyName)

	}
	*t = PropertyName{}
}

func (t *ConditionalExpression) clear() {
	t.LogicalORExpression.clear()
	if t.True != nil {
		t.True.clear()
		poolAssignmentExpression.Put(t.True)

	}
	if t.False != nil {
		t.False.clear()
		poolAssignmentExpression.Put(t.False)

	}
	*t = ConditionalExpression{}
}

func (t *LogicalORExpression) clear() {
	if t.LogicalORExpression != nil {
		t.LogicalORExpression.clear()
		poolLogicalORExpression.Put(t.LogicalORExpression)

	}
	t.LogicalANDExpression.clear()
	*t = LogicalORExpression{}
}

func (t *LogicalANDExpression) clear() {
	if t.LogicalANDExpression != nil {
		t.LogicalANDExpression.clear()
		poolLogicalANDExpression.Put(t.LogicalANDExpression)

	}
	t.BitwiseORExpression.clear()
	*t = LogicalANDExpression{}
}

func (t *BitwiseORExpression) clear() {
	if t.BitwiseORExpression != nil {
		t.BitwiseORExpression.clear()
		poolBitwiseORExpression.Put(t.BitwiseORExpression)

	}
	t.BitwiseXORExpression.clear()
	*t = BitwiseORExpression{}
}

func (t *BitwiseXORExpression) clear() {
	if t.BitwiseXORExpression != nil {
		t.BitwiseXORExpression.clear()
		poolBitwiseXORExpression.Put(t.BitwiseXORExpression)

	}
	t.BitwiseANDExpression.clear()
	*t = BitwiseXORExpression{}
}

func (t *BitwiseANDExpression) clear() {
	if t.BitwiseANDExpression != nil {
		t.BitwiseANDExpression.clear()
		poolBitwiseANDExpression.Put(t.BitwiseANDExpression)

	}
	t.EqualityExpression.clear()
	*t = BitwiseANDExpression{}
}

func (t *EqualityExpression) clear() {
	if t.EqualityExpression != nil {
		t.EqualityExpression.clear()
		poolEqualityExpression.Put(t.EqualityExpression)

	}
	t.RelationalExpression.clear()
	*t = EqualityExpression{}
}

func (t *RelationalExpression) clear() {
	if t.RelationalExpression != nil {
		t.RelationalExpression.clear()
		poolRelationalExpression.Put(t.RelationalExpression)

	}
	t.ShiftExpression.clear()
	*t = RelationalExpression{}
}

func (t *ShiftExpression) clear() {
	if t.ShiftExpression != nil {
		t.ShiftExpression.clear()
		poolShiftExpression.Put(t.ShiftExpression)

	}
	t.AdditiveExpression.clear()
	*t = ShiftExpression{}
}

func (t *AdditiveExpression) clear() {
	if t.AdditiveExpression != nil {
		t.AdditiveExpression.clear()
		poolAdditiveExpression.Put(t.AdditiveExpression)

	}
	t.MultiplicativeExpression.clear()
	*t = AdditiveExpression{}
}

func (t *MultiplicativeExpression) clear() {
	if t.MultiplicativeExpression != nil {
		t.MultiplicativeExpression.clear()
		poolMultiplicativeExpression.Put(t.MultiplicativeExpression)

	}
	t.ExponentiationExpression.clear()
	*t = MultiplicativeExpression{}
}

func (t *ExponentiationExpression) clear() {
	if t.ExponentiationExpression != nil {
		t.ExponentiationExpression.clear()
		poolExponentiationExpression.Put(t.ExponentiationExpression)

	}
	t.UnaryExpression.clear()
	*t = ExponentiationExpression{}
}

func (t *UnaryExpression) clear() {
	t.UpdateExpression.clear()
	*t = UnaryExpression{}
}

func (t *UpdateExpression) clear() {
	if t.LeftHandSideExpression != nil {
		t.LeftHandSideExpression.clear()
		poolLeftHandSideExpression.Put(t.LeftHandSideExpression)

	}
	if t.UnaryExpression != nil {
		t.UnaryExpression.clear()
		poolUnaryExpression.Put(t.UnaryExpression)

	}
	*t = UpdateExpression{}
}

func (t *AssignmentExpression) clear() {
	if t.ConditionalExpression != nil {
		t.ConditionalExpression.clear()
		poolConditionalExpression.Put(t.ConditionalExpression)

	}
	if t.ArrowFunction != nil {
		t.ArrowFunction.clear()
		poolArrowFunction.Put(t.ArrowFunction)

	}
	if t.LeftHandSideExpression != nil {
		t.LeftHandSideExpression.clear()
		poolLeftHandSideExpression.Put(t.LeftHandSideExpression)

	}
	if t.AssignmentExpression != nil {
		t.AssignmentExpression.clear()
		poolAssignmentExpression.Put(t.AssignmentExpression)

	}
	*t = AssignmentExpression{}
}

func (t *LeftHandSideExpression) clear() {
	if t.NewExpression != nil {
		t.NewExpression.clear()
		poolNewExpression.Put(t.NewExpression)

	}
	if t.CallExpression != nil {
		t.CallExpression.clear()
		poolCallExpression.Put(t.CallExpression)

	}
	*t = LeftHandSideExpression{}
}

func (t *Expression) clear() {
	for _, e := range t.Expressions {
		e.clear()
	}
	*t = Expression{}
}

func (t *NewExpression) clear() {
	t.MemberExpression.clear()
	*t = NewExpression{}
}

func (t *MemberExpression) clear() {
	if t.MemberExpression != nil {
		t.MemberExpression.clear()
		poolMemberExpression.Put(t.MemberExpression)

	}
	if t.PrimaryExpression != nil {
		t.PrimaryExpression.clear()
		poolPrimaryExpression.Put(t.PrimaryExpression)

	}
	if t.Expression != nil {
		t.Expression.clear()
		poolExpression.Put(t.Expression)

	}
	t.IdentifierName = nil
	if t.TemplateLiteral != nil {
		t.TemplateLiteral.clear()
		poolTemplateLiteral.Put(t.TemplateLiteral)

	}
	if t.Arguments != nil {
		t.Arguments.clear()
		poolArguments.Put(t.Arguments)

	}
	*t = MemberExpression{}
}

func (t *PrimaryExpression) clear() {
	t.IdentifierReference = nil
	t.Literal = nil
	if t.ArrayLiteral != nil {
		t.ArrayLiteral.clear()
		poolArrayLiteral.Put(t.ArrayLiteral)

	}
	if t.ObjectLiteral != nil {
		t.ObjectLiteral.clear()
		poolObjectLiteral.Put(t.ObjectLiteral)

	}
	if t.FunctionExpression != nil {
		t.FunctionExpression.clear()
		poolFunctionDeclaration.Put(t.FunctionExpression)

	}
	if t.ClassExpression != nil {
		t.ClassExpression.clear()
		poolClassDeclaration.Put(t.ClassExpression)

	}
	if t.TemplateLiteral != nil {
		t.TemplateLiteral.clear()
		poolTemplateLiteral.Put(t.TemplateLiteral)

	}
	if t.CoverParenthesizedExpressionAndArrowParameterList != nil {
		t.CoverParenthesizedExpressionAndArrowParameterList.clear()
		poolCoverParenthesizedExpressionAndArrowParameterList.Put(t.CoverParenthesizedExpressionAndArrowParameterList)

	}
	*t = PrimaryExpression{}
}

func (t *Arguments) clear() {
	for _, e := range t.ArgumentList {
		e.clear()
	}
	if t.SpreadArgument != nil {
		t.SpreadArgument.clear()
		poolAssignmentExpression.Put(t.SpreadArgument)

	}
	*t = Arguments{}
}

func (t *CallExpression) clear() {
	if t.MemberExpression != nil {
		t.MemberExpression.clear()
		poolMemberExpression.Put(t.MemberExpression)

	}
	if t.ImportCall != nil {
		t.ImportCall.clear()
		poolAssignmentExpression.Put(t.ImportCall)

	}
	if t.CallExpression != nil {
		t.CallExpression.clear()
		poolCallExpression.Put(t.CallExpression)

	}
	if t.Arguments != nil {
		t.Arguments.clear()
		poolArguments.Put(t.Arguments)

	}
	if t.Expression != nil {
		t.Expression.clear()
		poolExpression.Put(t.Expression)

	}
	t.IdentifierName = nil
	if t.TemplateLiteral != nil {
		t.TemplateLiteral.clear()
		poolTemplateLiteral.Put(t.TemplateLiteral)

	}
	*t = CallExpression{}
}

func (t *FunctionDeclaration) clear() {
	t.BindingIdentifier = nil
	t.FormalParameters.clear()
	t.FunctionBody.clear()
	*t = FunctionDeclaration{}
}

func (t *FormalParameters) clear() {
	for _, e := range t.FormalParameterList {
		e.clear()
	}
	if t.FunctionRestParameter != nil {
		t.FunctionRestParameter.clear()
		poolFunctionRestParameter.Put(t.FunctionRestParameter)

	}
	*t = FormalParameters{}
}

func (t *BindingElement) clear() {
	t.SingleNameBinding = nil
	if t.ArrayBindingPattern != nil {
		t.ArrayBindingPattern.clear()
		poolArrayBindingPattern.Put(t.ArrayBindingPattern)

	}
	if t.ObjectBindingPattern != nil {
		t.ObjectBindingPattern.clear()
		poolObjectBindingPattern.Put(t.ObjectBindingPattern)

	}
	if t.Initializer != nil {
		t.Initializer.clear()
		poolAssignmentExpression.Put(t.Initializer)

	}
	*t = BindingElement{}
}

func (t *FunctionRestParameter) clear() {
	t.BindingIdentifier = nil
	if t.ArrayBindingPattern != nil {
		t.ArrayBindingPattern.clear()
		poolArrayBindingPattern.Put(t.ArrayBindingPattern)

	}
	if t.ObjectBindingPattern != nil {
		t.ObjectBindingPattern.clear()
		poolObjectBindingPattern.Put(t.ObjectBindingPattern)

	}
	*t = FunctionRestParameter{}
}

func (t *Script) clear() {
	for _, e := range t.StatementList {
		e.clear()
	}
	*t = Script{}
}

func (t *Declaration) clear() {
	if t.ClassDeclaration != nil {
		t.ClassDeclaration.clear()
		poolClassDeclaration.Put(t.ClassDeclaration)

	}
	if t.FunctionDeclaration != nil {
		t.FunctionDeclaration.clear()
		poolFunctionDeclaration.Put(t.FunctionDeclaration)

	}
	if t.LexicalDeclaration != nil {
		t.LexicalDeclaration.clear()
		poolLexicalDeclaration.Put(t.LexicalDeclaration)

	}
	*t = Declaration{}
}

func (t *LexicalDeclaration) clear() {
	for _, e := range t.BindingList {
		e.clear()
	}
	*t = LexicalDeclaration{}
}

func (t *LexicalBinding) clear() {
	t.BindingIdentifier = nil
	if t.ArrayBindingPattern != nil {
		t.ArrayBindingPattern.clear()
		poolArrayBindingPattern.Put(t.ArrayBindingPattern)

	}
	if t.ObjectBindingPattern != nil {
		t.ObjectBindingPattern.clear()
		poolObjectBindingPattern.Put(t.ObjectBindingPattern)

	}
	if t.Initializer != nil {
		t.Initializer.clear()
		poolAssignmentExpression.Put(t.Initializer)

	}
	*t = LexicalBinding{}
}

func (t *ArrayBindingPattern) clear() {
	for _, e := range t.BindingElementList {
		e.clear()
	}
	if t.BindingRestElement != nil {
		t.BindingRestElement.clear()
		poolBindingElement.Put(t.BindingRestElement)

	}
	*t = ArrayBindingPattern{}
}

func (t *ObjectBindingPattern) clear() {
	for _, e := range t.BindingPropertyList {
		e.clear()
	}
	t.BindingRestProperty = nil
	*t = ObjectBindingPattern{}
}

func (t *BindingProperty) clear() {
	t.SingleNameBinding = nil
	if t.Initializer != nil {
		t.Initializer.clear()
		poolAssignmentExpression.Put(t.Initializer)

	}
	if t.PropertyName != nil {
		t.PropertyName.clear()
		poolPropertyName.Put(t.PropertyName)

	}
	if t.BindingElement != nil {
		t.BindingElement.clear()
		poolBindingElement.Put(t.BindingElement)

	}
	*t = BindingProperty{}
}

func (t *VariableDeclaration) clear() {
	t.BindingIdentifier = nil
	if t.ArrayBindingPattern != nil {
		t.ArrayBindingPattern.clear()
		poolArrayBindingPattern.Put(t.ArrayBindingPattern)

	}
	if t.ObjectBindingPattern != nil {
		t.ObjectBindingPattern.clear()
		poolObjectBindingPattern.Put(t.ObjectBindingPattern)

	}
	if t.Initializer != nil {
		t.Initializer.clear()
		poolAssignmentExpression.Put(t.Initializer)

	}
	*t = VariableDeclaration{}
}

func (t *ArrayLiteral) clear() {
	for _, e := range t.ElementList {
		e.clear()
	}
	if t.SpreadElement != nil {
		t.SpreadElement.clear()
		poolAssignmentExpression.Put(t.SpreadElement)

	}
	*t = ArrayLiteral{}
}

func (t *ObjectLiteral) clear() {
	for _, e := range t.PropertyDefinitionList {
		e.clear()
	}
	*t = ObjectLiteral{}
}

func (t *PropertyDefinition) clear() {
	t.IdentifierReference = nil
	if t.PropertyName != nil {
		t.PropertyName.clear()
		poolPropertyName.Put(t.PropertyName)

	}
	if t.AssignmentExpression != nil {
		t.AssignmentExpression.clear()
		poolAssignmentExpression.Put(t.AssignmentExpression)

	}
	if t.MethodDefinition != nil {
		t.MethodDefinition.clear()
		poolMethodDefinition.Put(t.MethodDefinition)

	}
	*t = PropertyDefinition{}
}

func (t *TemplateLiteral) clear() {
	t.NoSubstitutionTemplate = nil
	t.TemplateHead = nil
	for _, e := range t.Expressions {
		e.clear()
	}
	t.TemplateTail = nil
	*t = TemplateLiteral{}
}

func (t *ArrowFunction) clear() {
	t.BindingIdentifier = nil
	if t.CoverParenthesizedExpressionAndArrowParameterList != nil {
		t.CoverParenthesizedExpressionAndArrowParameterList.clear()
		poolCoverParenthesizedExpressionAndArrowParameterList.Put(t.CoverParenthesizedExpressionAndArrowParameterList)

	}
	if t.FormalParameters != nil {
		t.FormalParameters.clear()
		poolFormalParameters.Put(t.FormalParameters)

	}
	if t.AssignmentExpression != nil {
		t.AssignmentExpression.clear()
		poolAssignmentExpression.Put(t.AssignmentExpression)

	}
	if t.FunctionBody != nil {
		t.FunctionBody.clear()
		poolBlock.Put(t.FunctionBody)

	}
	*t = ArrowFunction{}
}

func (t *Module) clear() {
	for _, e := range t.ModuleListItems {
		e.clear()
	}
	*t = Module{}
}

func (t *ModuleListItem) clear() {
	if t.ImportDeclaration != nil {
		t.ImportDeclaration.clear()
		poolImportDeclaration.Put(t.ImportDeclaration)

	}
	if t.StatementListItem != nil {
		t.StatementListItem.clear()
		poolStatementListItem.Put(t.StatementListItem)

	}
	if t.ExportDeclaration != nil {
		t.ExportDeclaration.clear()
		poolExportDeclaration.Put(t.ExportDeclaration)

	}
	*t = ModuleListItem{}
}

func (t *ImportDeclaration) clear() {
	if t.ImportClause != nil {
		t.ImportClause.clear()
		poolImportClause.Put(t.ImportClause)

	}
	t.FromClause.clear()
	*t = ImportDeclaration{}
}

func (t *ImportClause) clear() {
	t.ImportedDefaultBinding = nil
	t.NameSpaceImport = nil
	if t.NamedImports != nil {
		t.NamedImports.clear()
		poolNamedImports.Put(t.NamedImports)

	}
	*t = ImportClause{}
}

func (t *FromClause) clear() {
	t.ModuleSpecifier = nil
	*t = FromClause{}
}

func (t *NamedImports) clear() {
	for _, e := range t.ImportList {
		e.clear()
	}
	*t = NamedImports{}
}

func (t *ImportSpecifier) clear() {
	t.IdentifierName = nil
	t.ImportedBinding = nil
	*t = ImportSpecifier{}
}

func (t *ExportDeclaration) clear() {
	if t.ExportClause != nil {
		t.ExportClause.clear()
		poolExportClause.Put(t.ExportClause)

	}
	if t.FromClause != nil {
		t.FromClause.clear()
		poolFromClause.Put(t.FromClause)

	}
	if t.VariableStatement != nil {
		t.VariableStatement.clear()
		poolVariableStatement.Put(t.VariableStatement)

	}
	if t.Declaration != nil {
		t.Declaration.clear()
		poolDeclaration.Put(t.Declaration)

	}
	if t.DefaultFunction != nil {
		t.DefaultFunction.clear()
		poolFunctionDeclaration.Put(t.DefaultFunction)

	}
	if t.DefaultClass != nil {
		t.DefaultClass.clear()
		poolClassDeclaration.Put(t.DefaultClass)

	}
	if t.DefaultAssignmentExpression != nil {
		t.DefaultAssignmentExpression.clear()
		poolAssignmentExpression.Put(t.DefaultAssignmentExpression)

	}
	*t = ExportDeclaration{}
}

func (t *ExportClause) clear() {
	for _, e := range t.ExportList {
		e.clear()
	}
	*t = ExportClause{}
}

func (t *ExportSpecifier) clear() {
	t.IdentifierName = nil
	t.EIdentifierName = nil
	*t = ExportSpecifier{}
}

func (t *Block) clear() {
	for _, e := range t.StatementListItems {
		e.clear()
	}
	*t = Block{}
}

func (t *StatementListItem) clear() {
	if t.Statement != nil {
		t.Statement.clear()
		poolStatement.Put(t.Statement)

	}
	if t.Declaration != nil {
		t.Declaration.clear()
		poolDeclaration.Put(t.Declaration)

	}
	*t = StatementListItem{}
}

func (t *Statement) clear() {
	if t.BlockStatement != nil {
		t.BlockStatement.clear()
		poolBlock.Put(t.BlockStatement)

	}
	if t.VariableStatement != nil {
		t.VariableStatement.clear()
		poolVariableStatement.Put(t.VariableStatement)

	}
	if t.ExpressionStatement != nil {
		t.ExpressionStatement.clear()
		poolExpression.Put(t.ExpressionStatement)

	}
	if t.IfStatement != nil {
		t.IfStatement.clear()
		poolIfStatement.Put(t.IfStatement)

	}
	if t.IterationStatementDo != nil {
		t.IterationStatementDo.clear()
		poolIterationStatementDo.Put(t.IterationStatementDo)

	}
	if t.IterationStatementWhile != nil {
		t.IterationStatementWhile.clear()
		poolIterationStatementWhile.Put(t.IterationStatementWhile)

	}
	if t.IterationStatementFor != nil {
		t.IterationStatementFor.clear()
		poolIterationStatementFor.Put(t.IterationStatementFor)

	}
	if t.SwitchStatement != nil {
		t.SwitchStatement.clear()
		poolSwitchStatement.Put(t.SwitchStatement)

	}
	if t.WithStatement != nil {
		t.WithStatement.clear()
		poolWithStatement.Put(t.WithStatement)

	}
	t.LabelIdentifier = nil
	if t.LabelledItemFunction != nil {
		t.LabelledItemFunction.clear()
		poolFunctionDeclaration.Put(t.LabelledItemFunction)

	}
	if t.LabelledItemStatement != nil {
		t.LabelledItemStatement.clear()
		poolStatement.Put(t.LabelledItemStatement)

	}
	if t.TryStatement != nil {
		t.TryStatement.clear()
		poolTryStatement.Put(t.TryStatement)

	}
	t.DebuggerStatement = nil
	*t = Statement{}
}

func (t *IfStatement) clear() {
	t.Expression.clear()
	t.Statement.clear()
	if t.ElseStatement != nil {
		t.ElseStatement.clear()
		poolStatement.Put(t.ElseStatement)

	}
	*t = IfStatement{}
}

func (t *IterationStatementDo) clear() {
	t.Statement.clear()
	t.Expression.clear()
	*t = IterationStatementDo{}
}

func (t *IterationStatementWhile) clear() {
	t.Expression.clear()
	t.Statement.clear()
	*t = IterationStatementWhile{}
}

func (t *IterationStatementFor) clear() {
	if t.InitExpression != nil {
		t.InitExpression.clear()
		poolExpression.Put(t.InitExpression)

	}
	for _, e := range t.InitVar {
		e.clear()
	}
	if t.InitLexical != nil {
		t.InitLexical.clear()
		poolLexicalDeclaration.Put(t.InitLexical)

	}
	if t.Conditional != nil {
		t.Conditional.clear()
		poolExpression.Put(t.Conditional)

	}
	if t.Afterthought != nil {
		t.Afterthought.clear()
		poolExpression.Put(t.Afterthought)

	}
	if t.LeftHandSideExpression != nil {
		t.LeftHandSideExpression.clear()
		poolLeftHandSideExpression.Put(t.LeftHandSideExpression)

	}
	t.ForBindingIdentifier = nil
	if t.ForBindingPatternObject != nil {
		t.ForBindingPatternObject.clear()
		poolObjectBindingPattern.Put(t.ForBindingPatternObject)

	}
	if t.ForBindingPatternArray != nil {
		t.ForBindingPatternArray.clear()
		poolArrayBindingPattern.Put(t.ForBindingPatternArray)

	}
	if t.In != nil {
		t.In.clear()
		poolExpression.Put(t.In)

	}
	if t.Of != nil {
		t.Of.clear()
		poolAssignmentExpression.Put(t.Of)

	}
	t.Statement.clear()
	*t = IterationStatementFor{}
}

func (t *SwitchStatement) clear() {
	t.Expression.clear()
	for _, e := range t.CaseClauses {
		e.clear()
	}
	for _, e := range t.DefaultClause {
		e.clear()
	}
	for _, e := range t.PostDefaultCaseClauses {
		e.clear()
	}
	*t = SwitchStatement{}
}

func (t *CaseClause) clear() {
	t.Expression.clear()
	for _, e := range t.StatementList {
		e.clear()
	}
	*t = CaseClause{}
}

func (t *WithStatement) clear() {
	t.Expression.clear()
	t.Statement.clear()
	*t = WithStatement{}
}

func (t *TryStatement) clear() {
	t.TryBlock.clear()
	t.CatchParameterBindingIdentifier = nil
	if t.CatchParameterObjectBindingPattern != nil {
		t.CatchParameterObjectBindingPattern.clear()
		poolObjectBindingPattern.Put(t.CatchParameterObjectBindingPattern)

	}
	if t.CatchParameterArrayBindingPattern != nil {
		t.CatchParameterArrayBindingPattern.clear()
		poolArrayBindingPattern.Put(t.CatchParameterArrayBindingPattern)

	}
	if t.CatchBlock != nil {
		t.CatchBlock.clear()
		poolBlock.Put(t.CatchBlock)

	}
	if t.FinallyBlock != nil {
		t.FinallyBlock.clear()
		poolBlock.Put(t.FinallyBlock)

	}
	*t = TryStatement{}
}

func (t *VariableStatement) clear() {
	for _, e := range t.VariableDeclarationList {
		e.clear()
	}
	*t = VariableStatement{}
}

func (t *CoverParenthesizedExpressionAndArrowParameterList) clear() {
	for _, e := range t.Expressions {
		e.clear()
	}
	t.BindingIdentifier = nil
	if t.ArrayBindingPattern != nil {
		t.ArrayBindingPattern.clear()
		poolArrayBindingPattern.Put(t.ArrayBindingPattern)

	}
	if t.ObjectBindingPattern != nil {
		t.ObjectBindingPattern.clear()
		poolObjectBindingPattern.Put(t.ObjectBindingPattern)

	}
	*t = CoverParenthesizedExpressionAndArrowParameterList{}
}
