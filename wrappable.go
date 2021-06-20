package javascript

type ConditionalWrappable interface {
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
