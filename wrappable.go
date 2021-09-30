package javascript

// ConditionalWrappable is an interface that is implemented by all types that
// are accepted by the WrapConditional function, and is the returned type of
// the UnwrapConditional function
type ConditionalWrappable interface {
	Type
	conditionalWrapabble()
}

func (ConditionalExpression) conditionalWrapabble() {}

func (LogicalORExpression) conditionalWrapabble() {}

func (LogicalANDExpression) conditionalWrapabble() {}

func (BitwiseORExpression) conditionalWrapabble() {}

func (BitwiseXORExpression) conditionalWrapabble() {}

func (BitwiseANDExpression) conditionalWrapabble() {}

func (EqualityExpression) conditionalWrapabble() {}

func (RelationalExpression) conditionalWrapabble() {}

func (ShiftExpression) conditionalWrapabble() {}

func (AdditiveExpression) conditionalWrapabble() {}

func (MultiplicativeExpression) conditionalWrapabble() {}

func (ExponentiationExpression) conditionalWrapabble() {}

func (UnaryExpression) conditionalWrapabble() {}

func (UpdateExpression) conditionalWrapabble() {}

func (LeftHandSideExpression) conditionalWrapabble() {}

func (CallExpression) conditionalWrapabble() {}

func (NewExpression) conditionalWrapabble() {}

func (MemberExpression) conditionalWrapabble() {}

func (PrimaryExpression) conditionalWrapabble() {}

func (ArrayLiteral) conditionalWrapabble() {}

func (ObjectLiteral) conditionalWrapabble() {}

func (FunctionDeclaration) conditionalWrapabble() {}

func (ClassDeclaration) conditionalWrapabble() {}

func (TemplateLiteral) conditionalWrapabble() {}

func (ParenthesizedExpression) conditionalWrapabble() {}
