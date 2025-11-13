package javascript

// File automatically generated with format.sh.

import "fmt"

// Type is an interface satisfied by all javascript structural types.
type Type interface {
	fmt.Formatter
	javascriptType()
}

func (Tokens) javascriptType() {}

func (Token) javascriptType() {}

func (AdditiveExpression) javascriptType() {}

func (Argument) javascriptType() {}

func (Arguments) javascriptType() {}

func (ArrayAssignmentPattern) javascriptType() {}

func (ArrayBindingPattern) javascriptType() {}

func (ArrayElement) javascriptType() {}

func (ArrayLiteral) javascriptType() {}

func (ArrowFunction) javascriptType() {}

func (AssignmentElement) javascriptType() {}

func (AssignmentExpression) javascriptType() {}

func (AssignmentPattern) javascriptType() {}

func (AssignmentProperty) javascriptType() {}

func (BindingElement) javascriptType() {}

func (BindingProperty) javascriptType() {}

func (BitwiseANDExpression) javascriptType() {}

func (BitwiseORExpression) javascriptType() {}

func (BitwiseXORExpression) javascriptType() {}

func (Block) javascriptType() {}

func (CallExpression) javascriptType() {}

func (CaseClause) javascriptType() {}

func (ClassDeclaration) javascriptType() {}

func (ClassElement) javascriptType() {}

func (ClassElementName) javascriptType() {}

func (CoalesceExpression) javascriptType() {}

func (ConditionalExpression) javascriptType() {}

func (Declaration) javascriptType() {}

func (DestructuringAssignmentTarget) javascriptType() {}

func (EqualityExpression) javascriptType() {}

func (ExponentiationExpression) javascriptType() {}

func (ExportClause) javascriptType() {}

func (ExportDeclaration) javascriptType() {}

func (ExportSpecifier) javascriptType() {}

func (Expression) javascriptType() {}

func (FieldDefinition) javascriptType() {}

func (FormalParameters) javascriptType() {}

func (FromClause) javascriptType() {}

func (FunctionDeclaration) javascriptType() {}

func (IfStatement) javascriptType() {}

func (ImportClause) javascriptType() {}

func (ImportDeclaration) javascriptType() {}

func (ImportSpecifier) javascriptType() {}

func (IterationStatementDo) javascriptType() {}

func (IterationStatementFor) javascriptType() {}

func (IterationStatementWhile) javascriptType() {}

func (LeftHandSideExpression) javascriptType() {}

func (LexicalBinding) javascriptType() {}

func (LexicalDeclaration) javascriptType() {}

func (LogicalANDExpression) javascriptType() {}

func (LogicalORExpression) javascriptType() {}

func (MemberExpression) javascriptType() {}

func (MethodDefinition) javascriptType() {}

func (Module) javascriptType() {}

func (ModuleItem) javascriptType() {}

func (MultiplicativeExpression) javascriptType() {}

func (NamedImports) javascriptType() {}

func (NewExpression) javascriptType() {}

func (ObjectAssignmentPattern) javascriptType() {}

func (ObjectBindingPattern) javascriptType() {}

func (ObjectLiteral) javascriptType() {}

func (OptionalChain) javascriptType() {}

func (OptionalExpression) javascriptType() {}

func (ParenthesizedExpression) javascriptType() {}

func (PrimaryExpression) javascriptType() {}

func (PropertyDefinition) javascriptType() {}

func (PropertyName) javascriptType() {}

func (RelationalExpression) javascriptType() {}

func (Script) javascriptType() {}

func (ShiftExpression) javascriptType() {}

func (Statement) javascriptType() {}

func (StatementListItem) javascriptType() {}

func (SwitchStatement) javascriptType() {}

func (TemplateLiteral) javascriptType() {}

func (TryStatement) javascriptType() {}

func (UnaryExpression) javascriptType() {}

func (UnaryOperatorComments) javascriptType() {}

func (UpdateExpression) javascriptType() {}

func (VariableStatement) javascriptType() {}

func (WithClause) javascriptType() {}

func (WithEntry) javascriptType() {}

func (WithStatement) javascriptType() {}
