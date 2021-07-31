package javascript

import "io"

var (
	nameAdditiveExpression                                = []byte{'\n', 'A', 'd', 'd', 'i', 't', 'i', 'v', 'e', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameAdditiveOperator                                  = []byte{'\n', 'A', 'd', 'd', 'i', 't', 'i', 'v', 'e', 'O', 'p', 'e', 'r', 'a', 't', 'o', 'r', ':', ' '}
	nameMultiplicativeExpression                          = []byte{'\n', 'M', 'u', 'l', 't', 'i', 'p', 'l', 'i', 'c', 'a', 't', 'i', 'v', 'e', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameArguments                                         = []byte{'\n', 'A', 'r', 'g', 'u', 'm', 'e', 'n', 't', 's', ':', ' '}
	nameArgumentList                                      = []byte{'\n', 'A', 'r', 'g', 'u', 'm', 'e', 'n', 't', 'L', 'i', 's', 't', ':', ' '}
	nameSpreadArgument                                    = []byte{'\n', 'S', 'p', 'r', 'e', 'a', 'd', 'A', 'r', 'g', 'u', 'm', 'e', 'n', 't', ':', ' '}
	nameArrayBindingPattern                               = []byte{'\n', 'A', 'r', 'r', 'a', 'y', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'P', 'a', 't', 't', 'e', 'r', 'n', ':', ' '}
	nameBindingElementList                                = []byte{'\n', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'E', 'l', 'e', 'm', 'e', 'n', 't', 'L', 'i', 's', 't', ':', ' '}
	nameBindingRestElement                                = []byte{'\n', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'R', 'e', 's', 't', 'E', 'l', 'e', 'm', 'e', 'n', 't', ':', ' '}
	nameArrayLiteral                                      = []byte{'\n', 'A', 'r', 'r', 'a', 'y', 'L', 'i', 't', 'e', 'r', 'a', 'l', ':', ' '}
	nameElementList                                       = []byte{'\n', 'E', 'l', 'e', 'm', 'e', 'n', 't', 'L', 'i', 's', 't', ':', ' '}
	nameSpreadElement                                     = []byte{'\n', 'S', 'p', 'r', 'e', 'a', 'd', 'E', 'l', 'e', 'm', 'e', 'n', 't', ':', ' '}
	nameArrowFunction                                     = []byte{'\n', 'A', 'r', 'r', 'o', 'w', 'F', 'u', 'n', 'c', 't', 'i', 'o', 'n', ':', ' '}
	nameBindingIdentifier                                 = []byte{'\n', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'I', 'd', 'e', 'n', 't', 'i', 'f', 'i', 'e', 'r', ':', ' '}
	nameCoverParenthesizedExpressionAndArrowParameterList = []byte{'\n', 'C', 'o', 'v', 'e', 'r', 'P', 'a', 'r', 'e', 'n', 't', 'h', 'e', 's', 'i', 'z', 'e', 'd', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', 'A', 'n', 'd', 'A', 'r', 'r', 'o', 'w', 'P', 'a', 'r', 'a', 'm', 'e', 't', 'e', 'r', 'L', 'i', 's', 't', ':', ' '}
	nameFormalParameters                                  = []byte{'\n', 'F', 'o', 'r', 'm', 'a', 'l', 'P', 'a', 'r', 'a', 'm', 'e', 't', 'e', 'r', 's', ':', ' '}
	nameAssignmentExpression                              = []byte{'\n', 'A', 's', 's', 'i', 'g', 'n', 'm', 'e', 'n', 't', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameFunctionBody                                      = []byte{'\n', 'F', 'u', 'n', 'c', 't', 'i', 'o', 'n', 'B', 'o', 'd', 'y', ':', ' '}
	nameConditionalExpression                             = []byte{'\n', 'C', 'o', 'n', 'd', 'i', 't', 'i', 'o', 'n', 'a', 'l', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameLeftHandSideExpression                            = []byte{'\n', 'L', 'e', 'f', 't', 'H', 'a', 'n', 'd', 'S', 'i', 'd', 'e', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameAssignmentOperator                                = []byte{'\n', 'A', 's', 's', 'i', 'g', 'n', 'm', 'e', 'n', 't', 'O', 'p', 'e', 'r', 'a', 't', 'o', 'r', ':', ' '}
	nameBindingElement                                    = []byte{'\n', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'E', 'l', 'e', 'm', 'e', 'n', 't', ':', ' '}
	nameSingleNameBinding                                 = []byte{'\n', 'S', 'i', 'n', 'g', 'l', 'e', 'N', 'a', 'm', 'e', 'B', 'i', 'n', 'd', 'i', 'n', 'g', ':', ' '}
	nameObjectBindingPattern                              = []byte{'\n', 'O', 'b', 'j', 'e', 'c', 't', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'P', 'a', 't', 't', 'e', 'r', 'n', ':', ' '}
	nameInitializer                                       = []byte{'\n', 'I', 'n', 'i', 't', 'i', 'a', 'l', 'i', 'z', 'e', 'r', ':', ' '}
	nameBindingProperty                                   = []byte{'\n', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'P', 'r', 'o', 'p', 'e', 'r', 't', 'y', ':', ' '}
	namePropertyName                                      = []byte{'\n', 'P', 'r', 'o', 'p', 'e', 'r', 't', 'y', 'N', 'a', 'm', 'e', ':', ' '}
	nameBitwiseANDExpression                              = []byte{'\n', 'B', 'i', 't', 'w', 'i', 's', 'e', 'A', 'N', 'D', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameEqualityExpression                                = []byte{'\n', 'E', 'q', 'u', 'a', 'l', 'i', 't', 'y', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameBitwiseORExpression                               = []byte{'\n', 'B', 'i', 't', 'w', 'i', 's', 'e', 'O', 'R', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameBitwiseXORExpression                              = []byte{'\n', 'B', 'i', 't', 'w', 'i', 's', 'e', 'X', 'O', 'R', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameBlock                                             = []byte{'\n', 'B', 'l', 'o', 'c', 'k', ':', ' '}
	nameStatementList                                     = []byte{'\n', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', 'L', 'i', 's', 't', ':', ' '}
	nameCallExpression                                    = []byte{'\n', 'C', 'a', 'l', 'l', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameMemberExpression                                  = []byte{'\n', 'M', 'e', 'm', 'b', 'e', 'r', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameImportCall                                        = []byte{'\n', 'I', 'm', 'p', 'o', 'r', 't', 'C', 'a', 'l', 'l', ':', ' '}
	nameExpression                                        = []byte{'\n', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameIdentifierName                                    = []byte{'\n', 'I', 'd', 'e', 'n', 't', 'i', 'f', 'i', 'e', 'r', 'N', 'a', 'm', 'e', ':', ' '}
	nameTemplateLiteral                                   = []byte{'\n', 'T', 'e', 'm', 'p', 'l', 'a', 't', 'e', 'L', 'i', 't', 'e', 'r', 'a', 'l', ':', ' '}
	nameCaseClause                                        = []byte{'\n', 'C', 'a', 's', 'e', 'C', 'l', 'a', 'u', 's', 'e', ':', ' '}
	nameClassDeclaration                                  = []byte{'\n', 'C', 'l', 'a', 's', 's', 'D', 'e', 'c', 'l', 'a', 'r', 'a', 't', 'i', 'o', 'n', ':', ' '}
	nameClassHeritage                                     = []byte{'\n', 'C', 'l', 'a', 's', 's', 'H', 'e', 'r', 'i', 't', 'a', 'g', 'e', ':', ' '}
	nameClassBody                                         = []byte{'\n', 'C', 'l', 'a', 's', 's', 'B', 'o', 'd', 'y', ':', ' '}
	nameLogicalORExpression                               = []byte{'\n', 'L', 'o', 'g', 'i', 'c', 'a', 'l', 'O', 'R', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameTrue                                              = []byte{'\n', 'T', 'r', 'u', 'e', ':', ' '}
	nameFalse                                             = []byte{'\n', 'F', 'a', 'l', 's', 'e', ':', ' '}
	nameExpressions                                       = []byte{'\n', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', 's', ':', ' '}
	nameDeclaration                                       = []byte{'\n', 'D', 'e', 'c', 'l', 'a', 'r', 'a', 't', 'i', 'o', 'n', ':', ' '}
	nameFunctionDeclaration                               = []byte{'\n', 'F', 'u', 'n', 'c', 't', 'i', 'o', 'n', 'D', 'e', 'c', 'l', 'a', 'r', 'a', 't', 'i', 'o', 'n', ':', ' '}
	nameLexicalDeclaration                                = []byte{'\n', 'L', 'e', 'x', 'i', 'c', 'a', 'l', 'D', 'e', 'c', 'l', 'a', 'r', 'a', 't', 'i', 'o', 'n', ':', ' '}
	nameEqualityOperator                                  = []byte{'\n', 'E', 'q', 'u', 'a', 'l', 'i', 't', 'y', 'O', 'p', 'e', 'r', 'a', 't', 'o', 'r', ':', ' '}
	nameRelationalExpression                              = []byte{'\n', 'R', 'e', 'l', 'a', 't', 'i', 'o', 'n', 'a', 'l', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameExponentiationExpression                          = []byte{'\n', 'E', 'x', 'p', 'o', 'n', 'e', 'n', 't', 'i', 'a', 't', 'i', 'o', 'n', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameUnaryExpression                                   = []byte{'\n', 'U', 'n', 'a', 'r', 'y', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameExportClause                                      = []byte{'\n', 'E', 'x', 'p', 'o', 'r', 't', 'C', 'l', 'a', 'u', 's', 'e', ':', ' '}
	nameExportFromClause                                  = []byte{'\n', 'E', 'x', 'p', 'o', 'r', 't', 'F', 'r', 'o', 'm', 'C', 'l', 'a', 'u', 's', 'e', ':', ' '}
	nameExportList                                        = []byte{'\n', 'E', 'x', 'p', 'o', 'r', 't', 'L', 'i', 's', 't', ':', ' '}
	nameExportDeclaration                                 = []byte{'\n', 'E', 'x', 'p', 'o', 'r', 't', 'D', 'e', 'c', 'l', 'a', 'r', 'a', 't', 'i', 'o', 'n', ':', ' '}
	nameFromClause                                        = []byte{'\n', 'F', 'r', 'o', 'm', 'C', 'l', 'a', 'u', 's', 'e', ':', ' '}
	nameVariableStatement                                 = []byte{'\n', 'V', 'a', 'r', 'i', 'a', 'b', 'l', 'e', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', ':', ' '}
	nameDefaultFunction                                   = []byte{'\n', 'D', 'e', 'f', 'a', 'u', 'l', 't', 'F', 'u', 'n', 'c', 't', 'i', 'o', 'n', ':', ' '}
	nameDefaultClass                                      = []byte{'\n', 'D', 'e', 'f', 'a', 'u', 'l', 't', 'C', 'l', 'a', 's', 's', ':', ' '}
	nameDefaultAssignmentExpression                       = []byte{'\n', 'D', 'e', 'f', 'a', 'u', 'l', 't', 'A', 's', 's', 'i', 'g', 'n', 'm', 'e', 'n', 't', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameExportSpecifier                                   = []byte{'\n', 'E', 'x', 'p', 'o', 'r', 't', 'S', 'p', 'e', 'c', 'i', 'f', 'i', 'e', 'r', ':', ' '}
	nameEIdentifierName                                   = []byte{'\n', 'E', 'I', 'd', 'e', 'n', 't', 'i', 'f', 'i', 'e', 'r', 'N', 'a', 'm', 'e', ':', ' '}
	nameFormalParameterList                               = []byte{'\n', 'F', 'o', 'r', 'm', 'a', 'l', 'P', 'a', 'r', 'a', 'm', 'e', 't', 'e', 'r', 'L', 'i', 's', 't', ':', ' '}
	nameModuleSpecifier                                   = []byte{'\n', 'M', 'o', 'd', 'u', 'l', 'e', 'S', 'p', 'e', 'c', 'i', 'f', 'i', 'e', 'r', ':', ' '}
	nameType                                              = []byte{'\n', 'T', 'y', 'p', 'e', ':', ' '}
	nameIfStatement                                       = []byte{'\n', 'I', 'f', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', ':', ' '}
	nameStatement                                         = []byte{'\n', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', ':', ' '}
	nameElseStatement                                     = []byte{'\n', 'E', 'l', 's', 'e', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', ':', ' '}
	nameImportClause                                      = []byte{'\n', 'I', 'm', 'p', 'o', 'r', 't', 'C', 'l', 'a', 'u', 's', 'e', ':', ' '}
	nameImportedDefaultBinding                            = []byte{'\n', 'I', 'm', 'p', 'o', 'r', 't', 'e', 'd', 'D', 'e', 'f', 'a', 'u', 'l', 't', 'B', 'i', 'n', 'd', 'i', 'n', 'g', ':', ' '}
	nameNameSpaceImport                                   = []byte{'\n', 'N', 'a', 'm', 'e', 'S', 'p', 'a', 'c', 'e', 'I', 'm', 'p', 'o', 'r', 't', ':', ' '}
	nameNamedImports                                      = []byte{'\n', 'N', 'a', 'm', 'e', 'd', 'I', 'm', 'p', 'o', 'r', 't', 's', ':', ' '}
	nameImportDeclaration                                 = []byte{'\n', 'I', 'm', 'p', 'o', 'r', 't', 'D', 'e', 'c', 'l', 'a', 'r', 'a', 't', 'i', 'o', 'n', ':', ' '}
	nameImportSpecifier                                   = []byte{'\n', 'I', 'm', 'p', 'o', 'r', 't', 'S', 'p', 'e', 'c', 'i', 'f', 'i', 'e', 'r', ':', ' '}
	nameImportedBinding                                   = []byte{'\n', 'I', 'm', 'p', 'o', 'r', 't', 'e', 'd', 'B', 'i', 'n', 'd', 'i', 'n', 'g', ':', ' '}
	nameIterationStatementDo                              = []byte{'\n', 'I', 't', 'e', 'r', 'a', 't', 'i', 'o', 'n', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', 'D', 'o', ':', ' '}
	nameIterationStatementFor                             = []byte{'\n', 'I', 't', 'e', 'r', 'a', 't', 'i', 'o', 'n', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', 'F', 'o', 'r', ':', ' '}
	nameInitExpression                                    = []byte{'\n', 'I', 'n', 'i', 't', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameInitVar                                           = []byte{'\n', 'I', 'n', 'i', 't', 'V', 'a', 'r', ':', ' '}
	nameInitLexical                                       = []byte{'\n', 'I', 'n', 'i', 't', 'L', 'e', 'x', 'i', 'c', 'a', 'l', ':', ' '}
	nameConditional                                       = []byte{'\n', 'C', 'o', 'n', 'd', 'i', 't', 'i', 'o', 'n', 'a', 'l', ':', ' '}
	nameAfterthought                                      = []byte{'\n', 'A', 'f', 't', 'e', 'r', 't', 'h', 'o', 'u', 'g', 'h', 't', ':', ' '}
	nameForBindingIdentifier                              = []byte{'\n', 'F', 'o', 'r', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'I', 'd', 'e', 'n', 't', 'i', 'f', 'i', 'e', 'r', ':', ' '}
	nameForBindingPatternObject                           = []byte{'\n', 'F', 'o', 'r', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'P', 'a', 't', 't', 'e', 'r', 'n', 'O', 'b', 'j', 'e', 'c', 't', ':', ' '}
	nameForBindingPatternArray                            = []byte{'\n', 'F', 'o', 'r', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'P', 'a', 't', 't', 'e', 'r', 'n', 'A', 'r', 'r', 'a', 'y', ':', ' '}
	nameIn                                                = []byte{'\n', 'I', 'n', ':', ' '}
	nameOf                                                = []byte{'\n', 'O', 'f', ':', ' '}
	nameIterationStatementWhile                           = []byte{'\n', 'I', 't', 'e', 'r', 'a', 't', 'i', 'o', 'n', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', 'W', 'h', 'i', 'l', 'e', ':', ' '}
	nameNewExpression                                     = []byte{'\n', 'N', 'e', 'w', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameLexicalBinding                                    = []byte{'\n', 'L', 'e', 'x', 'i', 'c', 'a', 'l', 'B', 'i', 'n', 'd', 'i', 'n', 'g', ':', ' '}
	nameLetOrConst                                        = []byte{'\n', 'L', 'e', 't', 'O', 'r', 'C', 'o', 'n', 's', 't', ':', ' '}
	nameBindingList                                       = []byte{'\n', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'L', 'i', 's', 't', ':', ' '}
	nameLogicalANDExpression                              = []byte{'\n', 'L', 'o', 'g', 'i', 'c', 'a', 'l', 'A', 'N', 'D', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	namePrimaryExpression                                 = []byte{'\n', 'P', 'r', 'i', 'm', 'a', 'r', 'y', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameMethodDefinition                                  = []byte{'\n', 'M', 'e', 't', 'h', 'o', 'd', 'D', 'e', 'f', 'i', 'n', 'i', 't', 'i', 'o', 'n', ':', ' '}
	nameParams                                            = []byte{'\n', 'P', 'a', 'r', 'a', 'm', 's', ':', ' '}
	nameModule                                            = []byte{'\n', 'M', 'o', 'd', 'u', 'l', 'e', ':', ' '}
	nameModuleListItems                                   = []byte{'\n', 'M', 'o', 'd', 'u', 'l', 'e', 'L', 'i', 's', 't', 'I', 't', 'e', 'm', 's', ':', ' '}
	nameModuleItem                                        = []byte{'\n', 'M', 'o', 'd', 'u', 'l', 'e', 'I', 't', 'e', 'm', ':', ' '}
	nameStatementListItem                                 = []byte{'\n', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', 'L', 'i', 's', 't', 'I', 't', 'e', 'm', ':', ' '}
	nameMultiplicativeOperator                            = []byte{'\n', 'M', 'u', 'l', 't', 'i', 'p', 'l', 'i', 'c', 'a', 't', 'i', 'v', 'e', 'O', 'p', 'e', 'r', 'a', 't', 'o', 'r', ':', ' '}
	nameImportList                                        = []byte{'\n', 'I', 'm', 'p', 'o', 'r', 't', 'L', 'i', 's', 't', ':', ' '}
	nameBindingPropertyList                               = []byte{'\n', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'P', 'r', 'o', 'p', 'e', 'r', 't', 'y', 'L', 'i', 's', 't', ':', ' '}
	nameBindingRestProperty                               = []byte{'\n', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'R', 'e', 's', 't', 'P', 'r', 'o', 'p', 'e', 'r', 't', 'y', ':', ' '}
	nameObjectLiteral                                     = []byte{'\n', 'O', 'b', 'j', 'e', 'c', 't', 'L', 'i', 't', 'e', 'r', 'a', 'l', ':', ' '}
	namePropertyDefinitionList                            = []byte{'\n', 'P', 'r', 'o', 'p', 'e', 'r', 't', 'y', 'D', 'e', 'f', 'i', 'n', 'i', 't', 'i', 'o', 'n', 'L', 'i', 's', 't', ':', ' '}
	nameIdentifierReference                               = []byte{'\n', 'I', 'd', 'e', 'n', 't', 'i', 'f', 'i', 'e', 'r', 'R', 'e', 'f', 'e', 'r', 'e', 'n', 'c', 'e', ':', ' '}
	nameLiteral                                           = []byte{'\n', 'L', 'i', 't', 'e', 'r', 'a', 'l', ':', ' '}
	nameFunctionExpression                                = []byte{'\n', 'F', 'u', 'n', 'c', 't', 'i', 'o', 'n', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameClassExpression                                   = []byte{'\n', 'C', 'l', 'a', 's', 's', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	namePropertyDefinition                                = []byte{'\n', 'P', 'r', 'o', 'p', 'e', 'r', 't', 'y', 'D', 'e', 'f', 'i', 'n', 'i', 't', 'i', 'o', 'n', ':', ' '}
	nameLiteralPropertyName                               = []byte{'\n', 'L', 'i', 't', 'e', 'r', 'a', 'l', 'P', 'r', 'o', 'p', 'e', 'r', 't', 'y', 'N', 'a', 'm', 'e', ':', ' '}
	nameComputedPropertyName                              = []byte{'\n', 'C', 'o', 'm', 'p', 'u', 't', 'e', 'd', 'P', 'r', 'o', 'p', 'e', 'r', 't', 'y', 'N', 'a', 'm', 'e', ':', ' '}
	nameRelationshipOperator                              = []byte{'\n', 'R', 'e', 'l', 'a', 't', 'i', 'o', 'n', 's', 'h', 'i', 'p', 'O', 'p', 'e', 'r', 'a', 't', 'o', 'r', ':', ' '}
	nameShiftExpression                                   = []byte{'\n', 'S', 'h', 'i', 'f', 't', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameScript                                            = []byte{'\n', 'S', 'c', 'r', 'i', 'p', 't', ':', ' '}
	nameShiftOperator                                     = []byte{'\n', 'S', 'h', 'i', 'f', 't', 'O', 'p', 'e', 'r', 'a', 't', 'o', 'r', ':', ' '}
	nameBlockStatement                                    = []byte{'\n', 'B', 'l', 'o', 'c', 'k', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', ':', ' '}
	nameExpressionStatement                               = []byte{'\n', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', ':', ' '}
	nameSwitchStatement                                   = []byte{'\n', 'S', 'w', 'i', 't', 'c', 'h', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', ':', ' '}
	nameWithStatement                                     = []byte{'\n', 'W', 'i', 't', 'h', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', ':', ' '}
	nameLabelIdentifier                                   = []byte{'\n', 'L', 'a', 'b', 'e', 'l', 'I', 'd', 'e', 'n', 't', 'i', 'f', 'i', 'e', 'r', ':', ' '}
	nameLabelledItemFunction                              = []byte{'\n', 'L', 'a', 'b', 'e', 'l', 'l', 'e', 'd', 'I', 't', 'e', 'm', 'F', 'u', 'n', 'c', 't', 'i', 'o', 'n', ':', ' '}
	nameLabelledItemStatement                             = []byte{'\n', 'L', 'a', 'b', 'e', 'l', 'l', 'e', 'd', 'I', 't', 'e', 'm', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', ':', ' '}
	nameTryStatement                                      = []byte{'\n', 'T', 'r', 'y', 'S', 't', 'a', 't', 'e', 'm', 'e', 'n', 't', ':', ' '}
	nameCaseClauses                                       = []byte{'\n', 'C', 'a', 's', 'e', 'C', 'l', 'a', 'u', 's', 'e', 's', ':', ' '}
	nameDefaultClause                                     = []byte{'\n', 'D', 'e', 'f', 'a', 'u', 'l', 't', 'C', 'l', 'a', 'u', 's', 'e', ':', ' '}
	namePostDefaultCaseClauses                            = []byte{'\n', 'P', 'o', 's', 't', 'D', 'e', 'f', 'a', 'u', 'l', 't', 'C', 'a', 's', 'e', 'C', 'l', 'a', 'u', 's', 'e', 's', ':', ' '}
	nameNoSubstitutionTemplate                            = []byte{'\n', 'N', 'o', 'S', 'u', 'b', 's', 't', 'i', 't', 'u', 't', 'i', 'o', 'n', 'T', 'e', 'm', 'p', 'l', 'a', 't', 'e', ':', ' '}
	nameTemplateHead                                      = []byte{'\n', 'T', 'e', 'm', 'p', 'l', 'a', 't', 'e', 'H', 'e', 'a', 'd', ':', ' '}
	nameTemplateMiddleList                                = []byte{'\n', 'T', 'e', 'm', 'p', 'l', 'a', 't', 'e', 'M', 'i', 'd', 'd', 'l', 'e', 'L', 'i', 's', 't', ':', ' '}
	nameTemplateTail                                      = []byte{'\n', 'T', 'e', 'm', 'p', 'l', 'a', 't', 'e', 'T', 'a', 'i', 'l', ':', ' '}
	nameTryBlock                                          = []byte{'\n', 'T', 'r', 'y', 'B', 'l', 'o', 'c', 'k', ':', ' '}
	nameCatchParameterBindingIdentifier                   = []byte{'\n', 'C', 'a', 't', 'c', 'h', 'P', 'a', 'r', 'a', 'm', 'e', 't', 'e', 'r', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'I', 'd', 'e', 'n', 't', 'i', 'f', 'i', 'e', 'r', ':', ' '}
	nameCatchParameterObjectBindingPattern                = []byte{'\n', 'C', 'a', 't', 'c', 'h', 'P', 'a', 'r', 'a', 'm', 'e', 't', 'e', 'r', 'O', 'b', 'j', 'e', 'c', 't', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'P', 'a', 't', 't', 'e', 'r', 'n', ':', ' '}
	nameCatchParameterArrayBindingPattern                 = []byte{'\n', 'C', 'a', 't', 'c', 'h', 'P', 'a', 'r', 'a', 'm', 'e', 't', 'e', 'r', 'A', 'r', 'r', 'a', 'y', 'B', 'i', 'n', 'd', 'i', 'n', 'g', 'P', 'a', 't', 't', 'e', 'r', 'n', ':', ' '}
	nameCatchBlock                                        = []byte{'\n', 'C', 'a', 't', 'c', 'h', 'B', 'l', 'o', 'c', 'k', ':', ' '}
	nameFinallyBlock                                      = []byte{'\n', 'F', 'i', 'n', 'a', 'l', 'l', 'y', 'B', 'l', 'o', 'c', 'k', ':', ' '}
	nameUnaryOperators                                    = []byte{'\n', 'U', 'n', 'a', 'r', 'y', 'O', 'p', 'e', 'r', 'a', 't', 'o', 'r', 's', ':', ' '}
	nameUpdateExpression                                  = []byte{'\n', 'U', 'p', 'd', 'a', 't', 'e', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameUpdateOperator                                    = []byte{'\n', 'U', 'p', 'd', 'a', 't', 'e', 'O', 'p', 'e', 'r', 'a', 't', 'o', 'r', ':', ' '}
	nameVariableDeclaration                               = []byte{'\n', 'V', 'a', 'r', 'i', 'a', 'b', 'l', 'e', 'D', 'e', 'c', 'l', 'a', 'r', 'a', 't', 'i', 'o', 'n', ':', ' '}
	nameVariableDeclarationList                           = []byte{'\n', 'V', 'a', 'r', 'i', 'a', 'b', 'l', 'e', 'D', 'e', 'c', 'l', 'a', 'r', 'a', 't', 'i', 'o', 'n', 'L', 'i', 's', 't', ':', ' '}
	nameOptionalExpression                                = []byte{'\n', 'O', 'p', 't', 'i', 'o', 'n', 'a', 'l', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameOptionalChain                                     = []byte{'\n', 'O', 'p', 't', 'i', 'o', 'n', 'a', 'l', 'C', 'h', 'a', 'i', 'n', ':', ' '}
	nameCoalesceExpression                                = []byte{'\n', 'C', 'o', 'a', 'l', 'e', 's', 'c', 'e', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', ':', ' '}
	nameCoalesceExpressionHead                            = []byte{'\n', 'C', 'o', 'a', 'l', 'e', 's', 'c', 'e', 'E', 'x', 'p', 'r', 'e', 's', 's', 'i', 'o', 'n', 'H', 'e', 'a', 'd', ':', ' '}
)

func (f *AdditiveExpression) printType(w io.Writer, v bool) {
	w.Write(nameAdditiveExpression[1:19])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.AdditiveExpression != nil {
		pp.Write(nameAdditiveExpression)
		f.AdditiveExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameAdditiveExpression)
		pp.Write(nilStr)
	}
	pp.Write(nameAdditiveOperator)
	io.WriteString(&pp, f.AdditiveOperator.String())
	pp.Write(nameMultiplicativeExpression)
	f.MultiplicativeExpression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *Arguments) printType(w io.Writer, v bool) {
	w.Write(nameArguments[1:10])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.ArgumentList) > 0 {
		pp.Write(nameArgumentList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.ArgumentList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameArgumentList)
		pp.Write(arrayOpenClose)
	}
	if f.SpreadArgument != nil {
		pp.Write(nameSpreadArgument)
		f.SpreadArgument.printType(&pp, v)
	} else if v {
		pp.Write(nameSpreadArgument)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ArrayBindingPattern) printType(w io.Writer, v bool) {
	w.Write(nameArrayBindingPattern[1:20])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.BindingElementList) > 0 {
		pp.Write(nameBindingElementList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.BindingElementList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameBindingElementList)
		pp.Write(arrayOpenClose)
	}
	if f.BindingRestElement != nil {
		pp.Write(nameBindingRestElement)
		f.BindingRestElement.printType(&pp, v)
	} else if v {
		pp.Write(nameBindingRestElement)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ArrayLiteral) printType(w io.Writer, v bool) {
	w.Write(nameArrayLiteral[1:13])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.ElementList) > 0 {
		pp.Write(nameElementList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.ElementList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameElementList)
		pp.Write(arrayOpenClose)
	}
	if f.SpreadElement != nil {
		pp.Write(nameSpreadElement)
		f.SpreadElement.printType(&pp, v)
	} else if v {
		pp.Write(nameSpreadElement)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ArrowFunction) printType(w io.Writer, v bool) {
	w.Write(nameArrowFunction[1:14])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.Async || v {
		pp.Printf("\nAsync: %v", f.Async)
	}
	if f.BindingIdentifier != nil {
		pp.Write(nameBindingIdentifier)
		f.BindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Write(nameBindingIdentifier)
		pp.Write(nilStr)
	}
	if f.FormalParameters != nil {
		pp.Write(nameFormalParameters)
		f.FormalParameters.printType(&pp, v)
	} else if v {
		pp.Write(nameFormalParameters)
		pp.Write(nilStr)
	}
	if f.AssignmentExpression != nil {
		pp.Write(nameAssignmentExpression)
		f.AssignmentExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameAssignmentExpression)
		pp.Write(nilStr)
	}
	if f.FunctionBody != nil {
		pp.Write(nameFunctionBody)
		f.FunctionBody.printType(&pp, v)
	} else if v {
		pp.Write(nameFunctionBody)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *AssignmentExpression) printType(w io.Writer, v bool) {
	w.Write(nameAssignmentExpression[1:21])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.ConditionalExpression != nil {
		pp.Write(nameConditionalExpression)
		f.ConditionalExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameConditionalExpression)
		pp.Write(nilStr)
	}
	if f.ArrowFunction != nil {
		pp.Write(nameArrowFunction)
		f.ArrowFunction.printType(&pp, v)
	} else if v {
		pp.Write(nameArrowFunction)
		pp.Write(nilStr)
	}
	if f.LeftHandSideExpression != nil {
		pp.Write(nameLeftHandSideExpression)
		f.LeftHandSideExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameLeftHandSideExpression)
		pp.Write(nilStr)
	}
	if f.Yield || v {
		pp.Printf("\nYield: %v", f.Yield)
	}
	if f.Delegate || v {
		pp.Printf("\nDelegate: %v", f.Delegate)
	}
	pp.Write(nameAssignmentOperator)
	io.WriteString(&pp, f.AssignmentOperator.String())
	if f.AssignmentExpression != nil {
		pp.Write(nameAssignmentExpression)
		f.AssignmentExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameAssignmentExpression)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *BindingElement) printType(w io.Writer, v bool) {
	w.Write(nameBindingElement[1:15])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.SingleNameBinding != nil {
		pp.Write(nameSingleNameBinding)
		f.SingleNameBinding.printType(&pp, v)
	} else if v {
		pp.Write(nameSingleNameBinding)
		pp.Write(nilStr)
	}
	if f.ArrayBindingPattern != nil {
		pp.Write(nameArrayBindingPattern)
		f.ArrayBindingPattern.printType(&pp, v)
	} else if v {
		pp.Write(nameArrayBindingPattern)
		pp.Write(nilStr)
	}
	if f.ObjectBindingPattern != nil {
		pp.Write(nameObjectBindingPattern)
		f.ObjectBindingPattern.printType(&pp, v)
	} else if v {
		pp.Write(nameObjectBindingPattern)
		pp.Write(nilStr)
	}
	if f.Initializer != nil {
		pp.Write(nameInitializer)
		f.Initializer.printType(&pp, v)
	} else if v {
		pp.Write(nameInitializer)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *BindingProperty) printType(w io.Writer, v bool) {
	w.Write(nameBindingProperty[1:16])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Write(namePropertyName)
	f.PropertyName.printType(&pp, v)
	pp.Write(nameBindingElement)
	f.BindingElement.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *BitwiseANDExpression) printType(w io.Writer, v bool) {
	w.Write(nameBitwiseANDExpression[1:21])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.BitwiseANDExpression != nil {
		pp.Write(nameBitwiseANDExpression)
		f.BitwiseANDExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameBitwiseANDExpression)
		pp.Write(nilStr)
	}
	pp.Write(nameEqualityExpression)
	f.EqualityExpression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *BitwiseORExpression) printType(w io.Writer, v bool) {
	w.Write(nameBitwiseORExpression[1:20])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.BitwiseORExpression != nil {
		pp.Write(nameBitwiseORExpression)
		f.BitwiseORExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameBitwiseORExpression)
		pp.Write(nilStr)
	}
	pp.Write(nameBitwiseXORExpression)
	f.BitwiseXORExpression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *BitwiseXORExpression) printType(w io.Writer, v bool) {
	w.Write(nameBitwiseXORExpression[1:21])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.BitwiseXORExpression != nil {
		pp.Write(nameBitwiseXORExpression)
		f.BitwiseXORExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameBitwiseXORExpression)
		pp.Write(nilStr)
	}
	pp.Write(nameBitwiseANDExpression)
	f.BitwiseANDExpression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *Block) printType(w io.Writer, v bool) {
	w.Write(nameBlock[1:6])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.StatementList) > 0 {
		pp.Write(nameStatementList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.StatementList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameStatementList)
		pp.Write(arrayOpenClose)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *CallExpression) printType(w io.Writer, v bool) {
	w.Write(nameCallExpression[1:15])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.MemberExpression != nil {
		pp.Write(nameMemberExpression)
		f.MemberExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameMemberExpression)
		pp.Write(nilStr)
	}
	if f.SuperCall || v {
		pp.Printf("\nSuperCall: %v", f.SuperCall)
	}
	if f.ImportCall != nil {
		pp.Write(nameImportCall)
		f.ImportCall.printType(&pp, v)
	} else if v {
		pp.Write(nameImportCall)
		pp.Write(nilStr)
	}
	if f.CallExpression != nil {
		pp.Write(nameCallExpression)
		f.CallExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameCallExpression)
		pp.Write(nilStr)
	}
	if f.Arguments != nil {
		pp.Write(nameArguments)
		f.Arguments.printType(&pp, v)
	} else if v {
		pp.Write(nameArguments)
		pp.Write(nilStr)
	}
	if f.Expression != nil {
		pp.Write(nameExpression)
		f.Expression.printType(&pp, v)
	} else if v {
		pp.Write(nameExpression)
		pp.Write(nilStr)
	}
	if f.IdentifierName != nil {
		pp.Write(nameIdentifierName)
		f.IdentifierName.printType(&pp, v)
	} else if v {
		pp.Write(nameIdentifierName)
		pp.Write(nilStr)
	}
	if f.TemplateLiteral != nil {
		pp.Write(nameTemplateLiteral)
		f.TemplateLiteral.printType(&pp, v)
	} else if v {
		pp.Write(nameTemplateLiteral)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *CaseClause) printType(w io.Writer, v bool) {
	w.Write(nameCaseClause[1:11])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Write(nameExpression)
	f.Expression.printType(&pp, v)
	if len(f.StatementList) > 0 {
		pp.Write(nameStatementList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.StatementList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameStatementList)
		pp.Write(arrayOpenClose)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ClassDeclaration) printType(w io.Writer, v bool) {
	w.Write(nameClassDeclaration[1:17])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.BindingIdentifier != nil {
		pp.Write(nameBindingIdentifier)
		f.BindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Write(nameBindingIdentifier)
		pp.Write(nilStr)
	}
	if f.ClassHeritage != nil {
		pp.Write(nameClassHeritage)
		f.ClassHeritage.printType(&pp, v)
	} else if v {
		pp.Write(nameClassHeritage)
		pp.Write(nilStr)
	}
	if len(f.ClassBody) > 0 {
		pp.Write(nameClassBody)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.ClassBody {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameClassBody)
		pp.Write(arrayOpenClose)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *CoalesceExpression) printType(w io.Writer, v bool) {
	w.Write(nameCoalesceExpression[1:19])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.CoalesceExpressionHead != nil {
		pp.Write(nameCoalesceExpressionHead)
		f.CoalesceExpressionHead.printType(&pp, v)
	} else if v {
		pp.Write(nameCoalesceExpressionHead)
		pp.Write(nilStr)
	}
	pp.Write(nameBitwiseORExpression)
	f.BitwiseORExpression.printType(w, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ConditionalExpression) printType(w io.Writer, v bool) {
	w.Write(nameConditionalExpression[1:22])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.LogicalORExpression != nil {
		pp.Write(nameLogicalORExpression)
		f.LogicalORExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameLogicalORExpression)
		pp.Write(nilStr)
	}
	if f.CoalesceExpression != nil {
		pp.Write(nameCoalesceExpression)
		f.CoalesceExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameCoalesceExpression)
		pp.Write(nilStr)
	}
	if f.True != nil {
		pp.Write(nameTrue)
		f.True.printType(&pp, v)
	} else if v {
		pp.Write(nameTrue)
		pp.Write(nilStr)
	}
	if f.False != nil {
		pp.Write(nameFalse)
		f.False.printType(&pp, v)
	} else if v {
		pp.Write(nameFalse)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *CoverParenthesizedExpressionAndArrowParameterList) printType(w io.Writer, v bool) {
	w.Write(nameCoverParenthesizedExpressionAndArrowParameterList[1:50])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.Expressions) > 0 {
		pp.Write(nameExpressions)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.Expressions {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameExpressions)
		pp.Write(arrayOpenClose)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *Declaration) printType(w io.Writer, v bool) {
	w.Write(nameDeclaration[1:12])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.ClassDeclaration != nil {
		pp.Write(nameClassDeclaration)
		f.ClassDeclaration.printType(&pp, v)
	} else if v {
		pp.Write(nameClassDeclaration)
		pp.Write(nilStr)
	}
	if f.FunctionDeclaration != nil {
		pp.Write(nameFunctionDeclaration)
		f.FunctionDeclaration.printType(&pp, v)
	} else if v {
		pp.Write(nameFunctionDeclaration)
		pp.Write(nilStr)
	}
	if f.LexicalDeclaration != nil {
		pp.Write(nameLexicalDeclaration)
		f.LexicalDeclaration.printType(&pp, v)
	} else if v {
		pp.Write(nameLexicalDeclaration)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *EqualityExpression) printType(w io.Writer, v bool) {
	w.Write(nameEqualityExpression[1:19])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.EqualityExpression != nil {
		pp.Write(nameEqualityExpression)
		f.EqualityExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameEqualityExpression)
		pp.Write(nilStr)
	}
	pp.Write(nameEqualityOperator)
	io.WriteString(&pp, f.EqualityOperator.String())
	pp.Write(nameRelationalExpression)
	f.RelationalExpression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ExponentiationExpression) printType(w io.Writer, v bool) {
	w.Write(nameExponentiationExpression[1:25])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.ExponentiationExpression != nil {
		pp.Write(nameExponentiationExpression)
		f.ExponentiationExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameExponentiationExpression)
		pp.Write(nilStr)
	}
	pp.Write(nameUnaryExpression)
	f.UnaryExpression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ExportClause) printType(w io.Writer, v bool) {
	w.Write(nameExportClause[1:13])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.ExportList) > 0 {
		pp.Write(nameExportList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.ExportList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameExportList)
		pp.Write(arrayOpenClose)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ExportDeclaration) printType(w io.Writer, v bool) {
	w.Write(nameExportDeclaration[1:18])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.ExportClause != nil {
		pp.Write(nameExportClause)
		f.ExportClause.printType(&pp, v)
	} else if v {
		pp.Write(nameExportClause)
		pp.Write(nilStr)
	}
	if f.ExportFromClause != nil {
		pp.Write(nameExportFromClause)
		f.ExportFromClause.printType(&pp, v)
	} else if v {
		pp.Write(nameExportFromClause)
		pp.Write(nilStr)
	}
	if f.FromClause != nil {
		pp.Write(nameFromClause)
		f.FromClause.printType(&pp, v)
	} else if v {
		pp.Write(nameFromClause)
		pp.Write(nilStr)
	}
	if f.VariableStatement != nil {
		pp.Write(nameVariableStatement)
		f.VariableStatement.printType(&pp, v)
	} else if v {
		pp.Write(nameVariableStatement)
		pp.Write(nilStr)
	}
	if f.Declaration != nil {
		pp.Write(nameDeclaration)
		f.Declaration.printType(&pp, v)
	} else if v {
		pp.Write(nameDeclaration)
		pp.Write(nilStr)
	}
	if f.DefaultFunction != nil {
		pp.Write(nameDefaultFunction)
		f.DefaultFunction.printType(&pp, v)
	} else if v {
		pp.Write(nameDefaultFunction)
		pp.Write(nilStr)
	}
	if f.DefaultClass != nil {
		pp.Write(nameDefaultClass)
		f.DefaultClass.printType(&pp, v)
	} else if v {
		pp.Write(nameDefaultClass)
		pp.Write(nilStr)
	}
	if f.DefaultAssignmentExpression != nil {
		pp.Write(nameDefaultAssignmentExpression)
		f.DefaultAssignmentExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameDefaultAssignmentExpression)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ExportSpecifier) printType(w io.Writer, v bool) {
	w.Write(nameExportSpecifier[1:16])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.IdentifierName != nil {
		pp.Write(nameIdentifierName)
		f.IdentifierName.printType(&pp, v)
	} else if v {
		pp.Write(nameIdentifierName)
		pp.Write(nilStr)
	}
	if f.EIdentifierName != nil {
		pp.Write(nameEIdentifierName)
		f.EIdentifierName.printType(&pp, v)
	} else if v {
		pp.Write(nameEIdentifierName)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *Expression) printType(w io.Writer, v bool) {
	w.Write(nameExpression[1:11])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.Expressions) > 0 {
		pp.Write(nameExpressions)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.Expressions {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameExpressions)
		pp.Write(arrayOpenClose)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *FormalParameters) printType(w io.Writer, v bool) {
	w.Write(nameFormalParameters[1:17])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.FormalParameterList) > 0 {
		pp.Write(nameFormalParameterList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.FormalParameterList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameFormalParameterList)
		pp.Write(arrayOpenClose)
	}
	if f.BindingIdentifier != nil {
		pp.Write(nameBindingIdentifier)
		f.BindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Write(nameBindingIdentifier)
		pp.Write(nilStr)
	}
	if f.ArrayBindingPattern != nil {
		pp.Write(nameArrayBindingPattern)
		f.ArrayBindingPattern.printType(&pp, v)
	} else if v {
		pp.Write(nameArrayBindingPattern)
		pp.Write(nilStr)
	}
	if f.ObjectBindingPattern != nil {
		pp.Write(nameObjectBindingPattern)
		f.ObjectBindingPattern.printType(&pp, v)
	} else if v {
		pp.Write(nameObjectBindingPattern)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *FromClause) printType(w io.Writer, v bool) {
	w.Write(nameFromClause[1:11])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.ModuleSpecifier != nil {
		pp.Write(nameModuleSpecifier)
		f.ModuleSpecifier.printType(&pp, v)
	} else if v {
		pp.Write(nameModuleSpecifier)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *FunctionDeclaration) printType(w io.Writer, v bool) {
	w.Write(nameFunctionDeclaration[1:20])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Write(nameType)
	io.WriteString(&pp, f.Type.String())
	if f.BindingIdentifier != nil {
		pp.Write(nameBindingIdentifier)
		f.BindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Write(nameBindingIdentifier)
		pp.Write(nilStr)
	}
	pp.Write(nameFormalParameters)
	f.FormalParameters.printType(&pp, v)
	pp.Write(nameFunctionBody)
	f.FunctionBody.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *IfStatement) printType(w io.Writer, v bool) {
	w.Write(nameIfStatement[1:12])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Write(nameExpression)
	f.Expression.printType(&pp, v)
	pp.Write(nameStatement)
	f.Statement.printType(&pp, v)
	if f.ElseStatement != nil {
		pp.Write(nameElseStatement)
		f.ElseStatement.printType(&pp, v)
	} else if v {
		pp.Write(nameElseStatement)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ImportClause) printType(w io.Writer, v bool) {
	w.Write(nameImportClause[1:13])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.ImportedDefaultBinding != nil {
		pp.Write(nameImportedDefaultBinding)
		f.ImportedDefaultBinding.printType(&pp, v)
	} else if v {
		pp.Write(nameImportedDefaultBinding)
		pp.Write(nilStr)
	}
	if f.NameSpaceImport != nil {
		pp.Write(nameNameSpaceImport)
		f.NameSpaceImport.printType(&pp, v)
	} else if v {
		pp.Write(nameNameSpaceImport)
		pp.Write(nilStr)
	}
	if f.NamedImports != nil {
		pp.Write(nameNamedImports)
		f.NamedImports.printType(&pp, v)
	} else if v {
		pp.Write(nameNamedImports)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ImportDeclaration) printType(w io.Writer, v bool) {
	w.Write(nameImportDeclaration[1:18])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.ImportClause != nil {
		pp.Write(nameImportClause)
		f.ImportClause.printType(&pp, v)
	} else if v {
		pp.Write(nameImportClause)
		pp.Write(nilStr)
	}
	pp.Write(nameFromClause)
	f.FromClause.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ImportSpecifier) printType(w io.Writer, v bool) {
	w.Write(nameImportSpecifier[1:16])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.IdentifierName != nil {
		pp.Write(nameIdentifierName)
		f.IdentifierName.printType(&pp, v)
	} else if v {
		pp.Write(nameIdentifierName)
		pp.Write(nilStr)
	}
	if f.ImportedBinding != nil {
		pp.Write(nameImportedBinding)
		f.ImportedBinding.printType(&pp, v)
	} else if v {
		pp.Write(nameImportedBinding)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *IterationStatementDo) printType(w io.Writer, v bool) {
	w.Write(nameIterationStatementDo[1:21])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Write(nameStatement)
	f.Statement.printType(&pp, v)
	pp.Write(nameExpression)
	f.Expression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *IterationStatementFor) printType(w io.Writer, v bool) {
	w.Write(nameIterationStatementFor[1:22])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Write(nameType)
	io.WriteString(&pp, f.Type.String())
	if f.InitExpression != nil {
		pp.Write(nameInitExpression)
		f.InitExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameInitExpression)
		pp.Write(nilStr)
	}
	if len(f.InitVar) > 0 {
		pp.Write(nameInitVar)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.InitVar {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameInitVar)
		pp.Write(arrayOpenClose)
	}
	if f.InitLexical != nil {
		pp.Write(nameInitLexical)
		f.InitLexical.printType(&pp, v)
	} else if v {
		pp.Write(nameInitLexical)
		pp.Write(nilStr)
	}
	if f.Conditional != nil {
		pp.Write(nameConditional)
		f.Conditional.printType(&pp, v)
	} else if v {
		pp.Write(nameConditional)
		pp.Write(nilStr)
	}
	if f.Afterthought != nil {
		pp.Write(nameAfterthought)
		f.Afterthought.printType(&pp, v)
	} else if v {
		pp.Write(nameAfterthought)
		pp.Write(nilStr)
	}
	if f.LeftHandSideExpression != nil {
		pp.Write(nameLeftHandSideExpression)
		f.LeftHandSideExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameLeftHandSideExpression)
		pp.Write(nilStr)
	}
	if f.ForBindingIdentifier != nil {
		pp.Write(nameForBindingIdentifier)
		f.ForBindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Write(nameForBindingIdentifier)
		pp.Write(nilStr)
	}
	if f.ForBindingPatternObject != nil {
		pp.Write(nameForBindingPatternObject)
		f.ForBindingPatternObject.printType(&pp, v)
	} else if v {
		pp.Write(nameForBindingPatternObject)
		pp.Write(nilStr)
	}
	if f.ForBindingPatternArray != nil {
		pp.Write(nameForBindingPatternArray)
		f.ForBindingPatternArray.printType(&pp, v)
	} else if v {
		pp.Write(nameForBindingPatternArray)
		pp.Write(nilStr)
	}
	if f.In != nil {
		pp.Write(nameIn)
		f.In.printType(&pp, v)
	} else if v {
		pp.Write(nameIn)
		pp.Write(nilStr)
	}
	if f.Of != nil {
		pp.Write(nameOf)
		f.Of.printType(&pp, v)
	} else if v {
		pp.Write(nameOf)
		pp.Write(nilStr)
	}
	pp.Write(nameStatement)
	f.Statement.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *IterationStatementWhile) printType(w io.Writer, v bool) {
	w.Write(nameIterationStatementWhile[1:24])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Write(nameExpression)
	f.Expression.printType(&pp, v)
	pp.Write(nameStatement)
	f.Statement.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *LeftHandSideExpression) printType(w io.Writer, v bool) {
	w.Write(nameLeftHandSideExpression[1:23])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.NewExpression != nil {
		pp.Write(nameNewExpression)
		f.NewExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameNewExpression)
		pp.Write(nilStr)
	}
	if f.CallExpression != nil {
		pp.Write(nameCallExpression)
		f.CallExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameCallExpression)
		pp.Write(nilStr)
	}
	if f.OptionalExpression != nil {
		pp.Write(nameOptionalExpression)
		f.OptionalExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameOptionalExpression)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *LexicalBinding) printType(w io.Writer, v bool) {
	w.Write(nameLexicalBinding[1:15])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.BindingIdentifier != nil {
		pp.Write(nameBindingIdentifier)
		f.BindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Write(nameBindingIdentifier)
		pp.Write(nilStr)
	}
	if f.ArrayBindingPattern != nil {
		pp.Write(nameArrayBindingPattern)
		f.ArrayBindingPattern.printType(&pp, v)
	} else if v {
		pp.Write(nameArrayBindingPattern)
		pp.Write(nilStr)
	}
	if f.ObjectBindingPattern != nil {
		pp.Write(nameObjectBindingPattern)
		f.ObjectBindingPattern.printType(&pp, v)
	} else if v {
		pp.Write(nameObjectBindingPattern)
		pp.Write(nilStr)
	}
	if f.Initializer != nil {
		pp.Write(nameInitializer)
		f.Initializer.printType(&pp, v)
	} else if v {
		pp.Write(nameInitializer)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *LexicalDeclaration) printType(w io.Writer, v bool) {
	w.Write(nameLexicalDeclaration[1:19])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Write(nameLetOrConst)
	io.WriteString(&pp, f.LetOrConst.String())
	if len(f.BindingList) > 0 {
		pp.Write(nameBindingList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.BindingList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameBindingList)
		pp.Write(arrayOpenClose)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *LogicalANDExpression) printType(w io.Writer, v bool) {
	w.Write(nameLogicalANDExpression[1:21])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.LogicalANDExpression != nil {
		pp.Write(nameLogicalANDExpression)
		f.LogicalANDExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameLogicalANDExpression)
		pp.Write(nilStr)
	}
	pp.Write(nameBitwiseORExpression)
	f.BitwiseORExpression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *LogicalORExpression) printType(w io.Writer, v bool) {
	w.Write(nameLogicalORExpression[1:20])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.LogicalORExpression != nil {
		pp.Write(nameLogicalORExpression)
		f.LogicalORExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameLogicalORExpression)
		pp.Write(nilStr)
	}
	pp.Write(nameLogicalANDExpression)
	f.LogicalANDExpression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *MemberExpression) printType(w io.Writer, v bool) {
	w.Write(nameMemberExpression[1:17])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.MemberExpression != nil {
		pp.Write(nameMemberExpression)
		f.MemberExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameMemberExpression)
		pp.Write(nilStr)
	}
	if f.PrimaryExpression != nil {
		pp.Write(namePrimaryExpression)
		f.PrimaryExpression.printType(&pp, v)
	} else if v {
		pp.Write(namePrimaryExpression)
		pp.Write(nilStr)
	}
	if f.Expression != nil {
		pp.Write(nameExpression)
		f.Expression.printType(&pp, v)
	} else if v {
		pp.Write(nameExpression)
		pp.Write(nilStr)
	}
	if f.IdentifierName != nil {
		pp.Write(nameIdentifierName)
		f.IdentifierName.printType(&pp, v)
	} else if v {
		pp.Write(nameIdentifierName)
		pp.Write(nilStr)
	}
	if f.TemplateLiteral != nil {
		pp.Write(nameTemplateLiteral)
		f.TemplateLiteral.printType(&pp, v)
	} else if v {
		pp.Write(nameTemplateLiteral)
		pp.Write(nilStr)
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
		pp.Write(nameArguments)
		f.Arguments.printType(&pp, v)
	} else if v {
		pp.Write(nameArguments)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *MethodDefinition) printType(w io.Writer, v bool) {
	w.Write(nameMethodDefinition[1:17])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Write(nameType)
	io.WriteString(&pp, f.Type.String())
	pp.Write(namePropertyName)
	f.PropertyName.printType(&pp, v)
	pp.Write(nameParams)
	f.Params.printType(&pp, v)
	pp.Write(nameFunctionBody)
	f.FunctionBody.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *Module) printType(w io.Writer, v bool) {
	w.Write(nameModule[1:7])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.ModuleListItems) > 0 {
		pp.Write(nameModuleListItems)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.ModuleListItems {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameModuleListItems)
		pp.Write(arrayOpenClose)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ModuleItem) printType(w io.Writer, v bool) {
	w.Write(nameModuleItem[1:11])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.ImportDeclaration != nil {
		pp.Write(nameImportDeclaration)
		f.ImportDeclaration.printType(&pp, v)
	} else if v {
		pp.Write(nameImportDeclaration)
		pp.Write(nilStr)
	}
	if f.StatementListItem != nil {
		pp.Write(nameStatementListItem)
		f.StatementListItem.printType(&pp, v)
	} else if v {
		pp.Write(nameStatementListItem)
		pp.Write(nilStr)
	}
	if f.ExportDeclaration != nil {
		pp.Write(nameExportDeclaration)
		f.ExportDeclaration.printType(&pp, v)
	} else if v {
		pp.Write(nameExportDeclaration)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *MultiplicativeExpression) printType(w io.Writer, v bool) {
	w.Write(nameMultiplicativeExpression[1:25])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.MultiplicativeExpression != nil {
		pp.Write(nameMultiplicativeExpression)
		f.MultiplicativeExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameMultiplicativeExpression)
		pp.Write(nilStr)
	}
	pp.Write(nameMultiplicativeOperator)
	io.WriteString(&pp, f.MultiplicativeOperator.String())
	pp.Write(nameExponentiationExpression)
	f.ExponentiationExpression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *NamedImports) printType(w io.Writer, v bool) {
	w.Write(nameNamedImports[1:13])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.ImportList) > 0 {
		pp.Write(nameImportList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.ImportList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameImportList)
		pp.Write(arrayOpenClose)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *NewExpression) printType(w io.Writer, v bool) {
	w.Write(nameNewExpression[1:14])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Printf("\nNews: %d", f.News)
	pp.Write(nameMemberExpression)
	f.MemberExpression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ObjectBindingPattern) printType(w io.Writer, v bool) {
	w.Write(nameObjectBindingPattern[1:21])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.BindingPropertyList) > 0 {
		pp.Write(nameBindingPropertyList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.BindingPropertyList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameBindingPropertyList)
		pp.Write(arrayOpenClose)
	}
	if f.BindingRestProperty != nil {
		pp.Write(nameBindingRestProperty)
		f.BindingRestProperty.printType(&pp, v)
	} else if v {
		pp.Write(nameBindingRestProperty)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ObjectLiteral) printType(w io.Writer, v bool) {
	w.Write(nameObjectLiteral[1:14])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.PropertyDefinitionList) > 0 {
		pp.Write(namePropertyDefinitionList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.PropertyDefinitionList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(namePropertyDefinitionList)
		pp.Write(arrayOpenClose)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *OptionalChain) printType(w io.Writer, v bool) {
	w.Write(nameOptionalChain[1:14])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.OptionalChain != nil {
		pp.Write(nameOptionalChain)
		f.OptionalChain.printType(&pp, v)
	} else if v {
		pp.Write(nameOptionalChain)
		pp.Write(nilStr)
	}
	if f.Arguments != nil {
		pp.Write(nameArguments)
		f.Arguments.printType(&pp, v)
	} else if v {
		pp.Write(nameArguments)
		pp.Write(nilStr)
	}
	if f.Expression != nil {
		pp.Write(nameExpression)
		f.Expression.printType(&pp, v)
	} else if v {
		pp.Write(nameExpression)
		pp.Write(nilStr)
	}
	if f.IdentifierName != nil {
		pp.Write(nameIdentifierName)
		f.IdentifierName.printType(&pp, v)
	} else if v {
		pp.Write(nameIdentifierName)
		pp.Write(nilStr)
	}
	if f.TemplateLiteral != nil {
		pp.Write(nameTemplateLiteral)
		f.TemplateLiteral.printType(&pp, v)
	} else if v {
		pp.Write(nameTemplateLiteral)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *OptionalExpression) printType(w io.Writer, v bool) {
	w.Write(nameOptionalExpression[1:20])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.MemberExpression != nil {
		pp.Write(nameMemberExpression)
		f.MemberExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameMemberExpression)
		pp.Write(nilStr)
	}
	if f.CallExpression != nil {
		pp.Write(nameCallExpression)
		f.CallExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameCallExpression)
		pp.Write(nilStr)
	}
	if f.OptionalExpression != nil {
		pp.Write(nameOptionalExpression)
		f.OptionalExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameOptionalExpression)
		pp.Write(nilStr)
	}
	pp.Write(nameOptionalChain)
	f.OptionalChain.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *PrimaryExpression) printType(w io.Writer, v bool) {
	w.Write(namePrimaryExpression[1:18])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.This != nil || v {
		pp.Printf("\nThis: %v", f.This)
	}
	if f.IdentifierReference != nil {
		pp.Write(nameIdentifierReference)
		f.IdentifierReference.printType(&pp, v)
	} else if v {
		pp.Write(nameIdentifierReference)
		pp.Write(nilStr)
	}
	if f.Literal != nil {
		pp.Write(nameLiteral)
		f.Literal.printType(&pp, v)
	} else if v {
		pp.Write(nameLiteral)
		pp.Write(nilStr)
	}
	if f.ArrayLiteral != nil {
		pp.Write(nameArrayLiteral)
		f.ArrayLiteral.printType(&pp, v)
	} else if v {
		pp.Write(nameArrayLiteral)
		pp.Write(nilStr)
	}
	if f.ObjectLiteral != nil {
		pp.Write(nameObjectLiteral)
		f.ObjectLiteral.printType(&pp, v)
	} else if v {
		pp.Write(nameObjectLiteral)
		pp.Write(nilStr)
	}
	if f.FunctionExpression != nil {
		pp.Write(nameFunctionExpression)
		f.FunctionExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameFunctionExpression)
		pp.Write(nilStr)
	}
	if f.ClassExpression != nil {
		pp.Write(nameClassExpression)
		f.ClassExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameClassExpression)
		pp.Write(nilStr)
	}
	if f.TemplateLiteral != nil {
		pp.Write(nameTemplateLiteral)
		f.TemplateLiteral.printType(&pp, v)
	} else if v {
		pp.Write(nameTemplateLiteral)
		pp.Write(nilStr)
	}
	if f.CoverParenthesizedExpressionAndArrowParameterList != nil {
		pp.Write(nameCoverParenthesizedExpressionAndArrowParameterList)
		f.CoverParenthesizedExpressionAndArrowParameterList.printType(&pp, v)
	} else if v {
		pp.Write(nameCoverParenthesizedExpressionAndArrowParameterList)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *PropertyDefinition) printType(w io.Writer, v bool) {
	w.Write(namePropertyDefinition[1:19])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.IsCoverInitializedName || v {
		pp.Printf("\nIsCoverInitializedName: %v", f.IsCoverInitializedName)
	}
	if f.PropertyName != nil {
		pp.Write(namePropertyName)
		f.PropertyName.printType(&pp, v)
	} else if v {
		pp.Write(namePropertyName)
		pp.Write(nilStr)
	}
	if f.AssignmentExpression != nil {
		pp.Write(nameAssignmentExpression)
		f.AssignmentExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameAssignmentExpression)
		pp.Write(nilStr)
	}
	if f.MethodDefinition != nil {
		pp.Write(nameMethodDefinition)
		f.MethodDefinition.printType(&pp, v)
	} else if v {
		pp.Write(nameMethodDefinition)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *PropertyName) printType(w io.Writer, v bool) {
	w.Write(namePropertyName[1:13])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.LiteralPropertyName != nil {
		pp.Write(nameLiteralPropertyName)
		f.LiteralPropertyName.printType(&pp, v)
	} else if v {
		pp.Write(nameLiteralPropertyName)
		pp.Write(nilStr)
	}
	if f.ComputedPropertyName != nil {
		pp.Write(nameComputedPropertyName)
		f.ComputedPropertyName.printType(&pp, v)
	} else if v {
		pp.Write(nameComputedPropertyName)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *RelationalExpression) printType(w io.Writer, v bool) {
	w.Write(nameRelationalExpression[1:21])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.RelationalExpression != nil {
		pp.Write(nameRelationalExpression)
		f.RelationalExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameRelationalExpression)
		pp.Write(nilStr)
	}
	pp.Write(nameRelationshipOperator)
	io.WriteString(&pp, f.RelationshipOperator.String())
	pp.Write(nameShiftExpression)
	f.ShiftExpression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *Script) printType(w io.Writer, v bool) {
	w.Write(nameScript[1:7])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.StatementList) > 0 {
		pp.Write(nameStatementList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.StatementList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameStatementList)
		pp.Write(arrayOpenClose)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *ShiftExpression) printType(w io.Writer, v bool) {
	w.Write(nameShiftExpression[1:16])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.ShiftExpression != nil {
		pp.Write(nameShiftExpression)
		f.ShiftExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameShiftExpression)
		pp.Write(nilStr)
	}
	pp.Write(nameShiftOperator)
	io.WriteString(&pp, f.ShiftOperator.String())
	pp.Write(nameAdditiveExpression)
	f.AdditiveExpression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *Statement) printType(w io.Writer, v bool) {
	w.Write(nameStatement[1:10])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Write(nameType)
	io.WriteString(&pp, f.Type.String())
	if f.BlockStatement != nil {
		pp.Write(nameBlockStatement)
		f.BlockStatement.printType(&pp, v)
	} else if v {
		pp.Write(nameBlockStatement)
		pp.Write(nilStr)
	}
	if f.VariableStatement != nil {
		pp.Write(nameVariableStatement)
		f.VariableStatement.printType(&pp, v)
	} else if v {
		pp.Write(nameVariableStatement)
		pp.Write(nilStr)
	}
	if f.ExpressionStatement != nil {
		pp.Write(nameExpressionStatement)
		f.ExpressionStatement.printType(&pp, v)
	} else if v {
		pp.Write(nameExpressionStatement)
		pp.Write(nilStr)
	}
	if f.IfStatement != nil {
		pp.Write(nameIfStatement)
		f.IfStatement.printType(&pp, v)
	} else if v {
		pp.Write(nameIfStatement)
		pp.Write(nilStr)
	}
	if f.IterationStatementDo != nil {
		pp.Write(nameIterationStatementDo)
		f.IterationStatementDo.printType(&pp, v)
	} else if v {
		pp.Write(nameIterationStatementDo)
		pp.Write(nilStr)
	}
	if f.IterationStatementWhile != nil {
		pp.Write(nameIterationStatementWhile)
		f.IterationStatementWhile.printType(&pp, v)
	} else if v {
		pp.Write(nameIterationStatementWhile)
		pp.Write(nilStr)
	}
	if f.IterationStatementFor != nil {
		pp.Write(nameIterationStatementFor)
		f.IterationStatementFor.printType(&pp, v)
	} else if v {
		pp.Write(nameIterationStatementFor)
		pp.Write(nilStr)
	}
	if f.SwitchStatement != nil {
		pp.Write(nameSwitchStatement)
		f.SwitchStatement.printType(&pp, v)
	} else if v {
		pp.Write(nameSwitchStatement)
		pp.Write(nilStr)
	}
	if f.WithStatement != nil {
		pp.Write(nameWithStatement)
		f.WithStatement.printType(&pp, v)
	} else if v {
		pp.Write(nameWithStatement)
		pp.Write(nilStr)
	}
	if f.LabelIdentifier != nil {
		pp.Write(nameLabelIdentifier)
		f.LabelIdentifier.printType(&pp, v)
	} else if v {
		pp.Write(nameLabelIdentifier)
		pp.Write(nilStr)
	}
	if f.LabelledItemFunction != nil {
		pp.Write(nameLabelledItemFunction)
		f.LabelledItemFunction.printType(&pp, v)
	} else if v {
		pp.Write(nameLabelledItemFunction)
		pp.Write(nilStr)
	}
	if f.LabelledItemStatement != nil {
		pp.Write(nameLabelledItemStatement)
		f.LabelledItemStatement.printType(&pp, v)
	} else if v {
		pp.Write(nameLabelledItemStatement)
		pp.Write(nilStr)
	}
	if f.TryStatement != nil {
		pp.Write(nameTryStatement)
		f.TryStatement.printType(&pp, v)
	} else if v {
		pp.Write(nameTryStatement)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *StatementListItem) printType(w io.Writer, v bool) {
	w.Write(nameStatementListItem[1:18])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.Statement != nil {
		pp.Write(nameStatement)
		f.Statement.printType(&pp, v)
	} else if v {
		pp.Write(nameStatement)
		pp.Write(nilStr)
	}
	if f.Declaration != nil {
		pp.Write(nameDeclaration)
		f.Declaration.printType(&pp, v)
	} else if v {
		pp.Write(nameDeclaration)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *SwitchStatement) printType(w io.Writer, v bool) {
	w.Write(nameSwitchStatement[1:16])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Write(nameExpression)
	f.Expression.printType(&pp, v)
	if len(f.CaseClauses) > 0 {
		pp.Write(nameCaseClauses)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.CaseClauses {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameCaseClauses)
		pp.Write(arrayOpenClose)
	}
	if len(f.DefaultClause) > 0 {
		pp.Write(nameDefaultClause)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.DefaultClause {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameDefaultClause)
		pp.Write(arrayOpenClose)
	}
	if len(f.PostDefaultCaseClauses) > 0 {
		pp.Write(namePostDefaultCaseClauses)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.PostDefaultCaseClauses {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(namePostDefaultCaseClauses)
		pp.Write(arrayOpenClose)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *TemplateLiteral) printType(w io.Writer, v bool) {
	w.Write(nameTemplateLiteral[1:16])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.NoSubstitutionTemplate != nil {
		pp.Write(nameNoSubstitutionTemplate)
		f.NoSubstitutionTemplate.printType(&pp, v)
	} else if v {
		pp.Write(nameNoSubstitutionTemplate)
		pp.Write(nilStr)
	}
	if f.TemplateHead != nil {
		pp.Write(nameTemplateHead)
		f.TemplateHead.printType(&pp, v)
	} else if v {
		pp.Write(nameTemplateHead)
		pp.Write(nilStr)
	}
	if len(f.Expressions) > 0 {
		pp.Write(nameExpressions)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.Expressions {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameExpressions)
		pp.Write(arrayOpenClose)
	}
	if len(f.TemplateMiddleList) > 0 {
		pp.Write(nameTemplateMiddleList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.TemplateMiddleList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameTemplateMiddleList)
		pp.Write(arrayOpenClose)
	}
	if f.TemplateTail != nil {
		pp.Write(nameTemplateTail)
		f.TemplateTail.printType(&pp, v)
	} else if v {
		pp.Write(nameTemplateTail)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *TryStatement) printType(w io.Writer, v bool) {
	w.Write(nameTryStatement[1:13])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Write(nameTryBlock)
	f.TryBlock.printType(&pp, v)
	if f.CatchParameterBindingIdentifier != nil {
		pp.Write(nameCatchParameterBindingIdentifier)
		f.CatchParameterBindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Write(nameCatchParameterBindingIdentifier)
		pp.Write(nilStr)
	}
	if f.CatchParameterObjectBindingPattern != nil {
		pp.Write(nameCatchParameterObjectBindingPattern)
		f.CatchParameterObjectBindingPattern.printType(&pp, v)
	} else if v {
		pp.Write(nameCatchParameterObjectBindingPattern)
		pp.Write(nilStr)
	}
	if f.CatchParameterArrayBindingPattern != nil {
		pp.Write(nameCatchParameterArrayBindingPattern)
		f.CatchParameterArrayBindingPattern.printType(&pp, v)
	} else if v {
		pp.Write(nameCatchParameterArrayBindingPattern)
		pp.Write(nilStr)
	}
	if f.CatchBlock != nil {
		pp.Write(nameCatchBlock)
		f.CatchBlock.printType(&pp, v)
	} else if v {
		pp.Write(nameCatchBlock)
		pp.Write(nilStr)
	}
	if f.FinallyBlock != nil {
		pp.Write(nameFinallyBlock)
		f.FinallyBlock.printType(&pp, v)
	} else if v {
		pp.Write(nameFinallyBlock)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *UnaryExpression) printType(w io.Writer, v bool) {
	w.Write(nameUnaryExpression[1:16])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.UnaryOperators) > 0 {
		pp.Write(nameUnaryOperators)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.UnaryOperators {
			ipp.Printf("\n%d:", n)
			ipp.WriteString(e.String())
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameUnaryOperators)
		pp.Write(arrayOpenClose)
	}
	pp.Write(nameUpdateExpression)
	f.UpdateExpression.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *UpdateExpression) printType(w io.Writer, v bool) {
	w.Write(nameUpdateExpression[1:17])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.LeftHandSideExpression != nil {
		pp.Write(nameLeftHandSideExpression)
		f.LeftHandSideExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameLeftHandSideExpression)
		pp.Write(nilStr)
	}
	pp.Write(nameUpdateOperator)
	io.WriteString(&pp, f.UpdateOperator.String())
	if f.UnaryExpression != nil {
		pp.Write(nameUnaryExpression)
		f.UnaryExpression.printType(&pp, v)
	} else if v {
		pp.Write(nameUnaryExpression)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *VariableDeclaration) printType(w io.Writer, v bool) {
	w.Write(nameVariableDeclaration[1:20])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if f.BindingIdentifier != nil {
		pp.Write(nameBindingIdentifier)
		f.BindingIdentifier.printType(&pp, v)
	} else if v {
		pp.Write(nameBindingIdentifier)
		pp.Write(nilStr)
	}
	if f.ArrayBindingPattern != nil {
		pp.Write(nameArrayBindingPattern)
		f.ArrayBindingPattern.printType(&pp, v)
	} else if v {
		pp.Write(nameArrayBindingPattern)
		pp.Write(nilStr)
	}
	if f.ObjectBindingPattern != nil {
		pp.Write(nameObjectBindingPattern)
		f.ObjectBindingPattern.printType(&pp, v)
	} else if v {
		pp.Write(nameObjectBindingPattern)
		pp.Write(nilStr)
	}
	if f.Initializer != nil {
		pp.Write(nameInitializer)
		f.Initializer.printType(&pp, v)
	} else if v {
		pp.Write(nameInitializer)
		pp.Write(nilStr)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *VariableStatement) printType(w io.Writer, v bool) {
	w.Write(nameVariableStatement[1:18])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	if len(f.VariableDeclarationList) > 0 {
		pp.Write(nameVariableDeclarationList)
		pp.Write(arrayOpen)
		ipp := indentPrinter{&pp}
		for n, e := range f.VariableDeclarationList {
			ipp.Printf("\n%d:", n)
			e.printType(&ipp, v)
		}
		pp.Write(arrayClose)
	} else if v {
		pp.Write(nameVariableDeclarationList)
		pp.Write(arrayOpenClose)
	}
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

func (f *WithStatement) printType(w io.Writer, v bool) {
	w.Write(nameWithStatement[1:14])
	w.Write(objectOpen)
	pp := indentPrinter{w}
	pp.Write(nameExpression)
	f.Expression.printType(&pp, v)
	pp.Write(nameStatement)
	f.Statement.printType(&pp, v)
	if v {
		pp.Write(tokensTo)
		f.Tokens.printType(&pp, v)
	}
	w.Write(objectClose)
}

// Format implements the fmt.Formatter interface
func (ft FunctionType) String() string {
	switch ft {
	case FunctionNormal:
		return "Normal"
	case FunctionGenerator:
		return "Generator"
	case FunctionAsync:
		return "Async"
	case FunctionAsyncGenerator:
		return "Async Generator"
	default:
		return unknown
	}
}

// String implements the fmt.Stringer interface
func (mt MethodType) String() string {
	switch mt {
	case MethodNormal:
		return "MethodNormal"
	case MethodGenerator:
		return "MethodGenerator"
	case MethodAsyncGenerator:
		return "MethodAsyncGenerator"
	case MethodAsync:
		return "MethodAsync"
	case MethodGetter:
		return "MethodGetter"
	case MethodSetter:
		return "MethodSetter"
	case MethodStatic:
		return "MethodStatic"
	case MethodStaticGenerator:
		return "MethodStaticGenerator"
	case MethodStaticAsync:
		return "MethodStaticAsync"
	case MethodStaticAsyncGenerator:
		return "MethodStaticAsyncGenerator"
	case MethodStaticGetter:
		return "MethodStaticGetter"
	case MethodStaticSetter:
		return "MethodStaticSetter"
	default:
		return unknown
	}
}

// String implements the fmt.Stringer interface
func (st StatementType) String() string {
	switch st {
	case StatementNormal:
		return "StatementNormal"
	case StatementContinue:
		return "StatementContinue"
	case StatementBreak:
		return "StatementBreak"
	case StatementReturn:
		return "StatementReturn"
	case StatementThrow:
		return "StatementThrow"
	default:
		return unknown
	}
}

// String implements the fmt.Stringer interface
func (ft ForType) String() string {
	switch ft {
	case ForNormal:
		return "ForNormal"
	case ForNormalVar:
		return "ForNormalVar"
	case ForNormalLexicalDeclaration:
		return "ForNormalLexicalDeclaration"
	case ForNormalExpression:
		return "ForNormalExpression"
	case ForInLeftHandSide:
		return "ForInLeftHandSide"
	case ForInVar:
		return "ForInVar"
	case ForInLet:
		return "ForInLet"
	case ForInConst:
		return "ForInConst"
	case ForOfLeftHandSide:
		return "ForOfLeftHandSide"
	case ForOfVar:
		return "ForOfVar"
	case ForOfLet:
		return "ForOfLet"
	case ForOfConst:
		return "ForOfConst"
	case ForAwaitOfLeftHandSide:
		return "ForAwaitOfLeftHandSide"
	case ForAwaitOfVar:
		return "ForAwaitOfVar"
	case ForAwaitOfLet:
		return "ForAwaitOfLet"
	case ForAwaitOfConst:
		return "ForAwaitOfConst"
	default:
		return unknown
	}
}

// String implements the fmt.Stringer interface
func (e EqualityOperator) String() string {
	switch e {
	case EqualityNone:
		return ""
	case EqualityEqual:
		return "=="
	case EqualityNotEqual:
		return "!="
	case EqualityStrictEqual:
		return "==="
	case EqualityStrictNotEqual:
		return "!=="
	default:
		return unknown
	}
}

// String implements the fmt.Stringer interface
func (r RelationshipOperator) String() string {
	switch r {
	case RelationshipNone:
		return ""
	case RelationshipLessThan:
		return "<"
	case RelationshipGreaterThan:
		return ">"
	case RelationshipLessThanEqual:
		return "<="
	case RelationshipGreaterThanEqual:
		return ">="
	case RelationshipInstanceOf:
		return "instanceof"
	case RelationshipIn:
		return "in"
	default:
		return unknown
	}
}

// String implements the fmt.Stringer interface
func (s ShiftOperator) String() string {
	switch s {
	case ShiftNone:
		return ""
	case ShiftLeft:
		return "<<"
	case ShiftRight:
		return ">>"
	case ShiftUnsignedRight:
		return ">>>"
	default:
		return unknown
	}
}

// String implements the fmt.Stringer interface
func (a AdditiveOperator) String() string {
	switch a {
	case AdditiveNone:
		return ""
	case AdditiveAdd:
		return "+"
	case AdditiveMinus:
		return "-"
	default:
		return unknown
	}
}

// String implements the fmt.Stringer interface
func (m MultiplicativeOperator) String() string {
	switch m {
	case MultiplicativeNone:
		return ""
	case MultiplicativeMultiply:
		return "*"
	case MultiplicativeDivide:
		return "/"
	case MultiplicativeRemainder:
		return "%"
	default:
		return unknown
	}
}

// String implements the fmt.Stringer interface
func (u UnaryOperator) String() string {
	switch u {
	case UnaryNone:
		return ""
	case UnaryDelete:
		return "delete"
	case UnaryVoid:
		return "void"
	case UnaryTypeOf:
		return "typeof"
	case UnaryAdd:
		return "+"
	case UnaryMinus:
		return "-"
	case UnaryBitwiseNot:
		return "~"
	case UnaryLogicalNot:
		return "!"
	case UnaryAwait:
		return "await"
	default:
		return unknown
	}
}

// String implements the fmt.Stringer interface
func (u UpdateOperator) String() string {
	switch u {
	case UpdateNone:
		return ""
	case UpdatePostIncrement:
		return " ++"
	case UpdatePostDecrement:
		return " --"
	case UpdatePreIncrement:
		return "++"
	case UpdatePreDecrement:
		return "--"
	default:
		return unknown
	}
}

// String implements the fmt.Stringer interface
func (l LetOrConst) String() string {
	if l {
		return "Const"
	}
	return "Let"
}

// String implements the fmt.Stringer interface
func (a AssignmentOperator) String() string {
	switch a {
	case AssignmentNone:
		return ""
	case AssignmentAssign:
		return "="
	case AssignmentMultiply:
		return "*="
	case AssignmentDivide:
		return "/="
	case AssignmentRemainder:
		return "%="
	case AssignmentAdd:
		return "+="
	case AssignmentSubtract:
		return "-="
	case AssignmentLeftShift:
		return "<<="
	case AssignmentSignPropagatinRightShift:
		return ">>="
	case AssignmentZeroFillRightShift:
		return ">>>="
	case AssignmentBitwiseAND:
		return "&="
	case AssignmentBitwiseXOR:
		return "^="
	case AssignmentBitwiseOR:
		return "|="
	case AssignmentExponentiation:
		return "**="
	case AssignmentLogicalAnd:
		return "&&="
	case AssignmentLogicalOr:
		return "||="
	case AssignmentNullish:
		return "??="
	default:
		return unknown
	}
}

const unknown = "Unknown"
