package javascript

import (
	"fmt"
)

// Type is an interface satisfied by all javascript structural types
type Type interface {
	fmt.Formatter
	javascriptType()
}

func (Token) javascriptType() {}

func (ClassDeclaration) javascriptType() {}

func (ClassElement) javascriptType() {}

func (FieldDefinition) javascriptType() {}

func (ClassElementName) javascriptType() {}

func (MethodDefinition) javascriptType() {}

func (PropertyName) javascriptType() {}

func (ConditionalExpression) javascriptType() {}

func (CoalesceExpression) javascriptType() {}

func (LogicalORExpression) javascriptType() {}

func (LogicalANDExpression) javascriptType() {}

func (BitwiseORExpression) javascriptType() {}

func (BitwiseXORExpression) javascriptType() {}

func (BitwiseANDExpression) javascriptType() {}

func (EqualityExpression) javascriptType() {}

func (RelationalExpression) javascriptType() {}

func (ShiftExpression) javascriptType() {}

func (AdditiveExpression) javascriptType() {}

func (MultiplicativeExpression) javascriptType() {}

func (ExponentiationExpression) javascriptType() {}

func (UnaryExpression) javascriptType() {}

func (UpdateExpression) javascriptType() {}

func (AssignmentExpression) javascriptType() {}

func (LeftHandSideExpression) javascriptType() {}

func (AssignmentPattern) javascriptType() {}

func (ObjectAssignmentPattern) javascriptType() {}

func (AssignmentProperty) javascriptType() {}

func (DestructuringAssignmentTarget) javascriptType() {}

func (AssignmentElement) javascriptType() {}

func (ArrayAssignmentPattern) javascriptType() {}

func (OptionalExpression) javascriptType() {}

func (OptionalChain) javascriptType() {}

func (Expression) javascriptType() {}

func (NewExpression) javascriptType() {}

func (MemberExpression) javascriptType() {}

func (PrimaryExpression) javascriptType() {}

func (ParenthesizedExpression) javascriptType() {}

func (Argument) javascriptType() {}

func (Arguments) javascriptType() {}

func (CallExpression) javascriptType() {}

func (FunctionDeclaration) javascriptType() {}

func (FormalParameters) javascriptType() {}

func (BindingElement) javascriptType() {}

func (Script) javascriptType() {}

func (Declaration) javascriptType() {}

func (LexicalDeclaration) javascriptType() {}

func (LexicalBinding) javascriptType() {}

func (ArrayBindingPattern) javascriptType() {}

func (ObjectBindingPattern) javascriptType() {}

func (BindingProperty) javascriptType() {}

func (ArrayElement) javascriptType() {}

func (ArrayLiteral) javascriptType() {}

func (ObjectLiteral) javascriptType() {}

func (PropertyDefinition) javascriptType() {}

func (TemplateLiteral) javascriptType() {}

func (ArrowFunction) javascriptType() {}

func (Module) javascriptType() {}

func (ModuleItem) javascriptType() {}

func (ImportDeclaration) javascriptType() {}

func (ImportClause) javascriptType() {}

func (FromClause) javascriptType() {}

func (NamedImports) javascriptType() {}

func (ImportSpecifier) javascriptType() {}

func (ExportDeclaration) javascriptType() {}

func (ExportClause) javascriptType() {}

func (ExportSpecifier) javascriptType() {}

func (Block) javascriptType() {}

func (StatementListItem) javascriptType() {}

func (Statement) javascriptType() {}

func (IfStatement) javascriptType() {}

func (IterationStatementDo) javascriptType() {}

func (IterationStatementWhile) javascriptType() {}

func (IterationStatementFor) javascriptType() {}

func (SwitchStatement) javascriptType() {}

func (CaseClause) javascriptType() {}

func (WithStatement) javascriptType() {}

func (TryStatement) javascriptType() {}

func (VariableStatement) javascriptType() {}
