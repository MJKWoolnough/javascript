package minify

import (
	"errors"
	"io"
	"unicode"
	"unicode/utf8"

	"vimagination.zapto.org/javascript"
)

var (
	idContinue = []*unicode.RangeTable{
		unicode.L,
		unicode.Nl,
		unicode.Other_ID_Start,
		unicode.Mn,
		unicode.Mc,
		unicode.Nd,
		unicode.Pc,
		unicode.Other_ID_Continue,
	}
	notID = []*unicode.RangeTable{
		unicode.Pattern_Syntax,
		unicode.Pattern_White_Space,
	}
)

const (
	zwnj rune = 8204
	zwj  rune = 8205
)

func isIDContinue(c rune) bool {
	if c == '$' || c == '_' || c == '\\' || c == zwnj || c == zwj {
		return true
	}
	return unicode.In(c, idContinue...) && !unicode.In(c, notID...)
}

type writer struct {
	io.Writer
	count    int64
	err      error
	lastChar rune
}

func (w *writer) WriteString(str string) {
	if w.err == nil {
		var n int
		if isIDContinue(w.lastChar) {
			r, _ := utf8.DecodeRuneInString(str)
			if isIDContinue(r) {
				n, w.err = io.WriteString(w.Writer, " ")
				w.count += int64(n)
			}
		}
		n, w.err = io.WriteString(w.Writer, str)
		w.count += int64(n)
		if len(str) > 0 {
			w.lastChar, _ = utf8.DecodeLastRuneInString(str)
		}
	}
}

func (w *writer) WriteEOS() {
	if w.err == nil && w.lastChar != '}' {
		w.WriteString(";")
	}
}

func Print(w io.Writer, m *javascript.Module) (int64, error) {
	wr := writer{Writer: w}
	for n := range m.ModuleListItems {
		if n > 0 {
			wr.WriteEOS()
		}
		wr.WriteModuleListItem(&m.ModuleListItems[n])
	}
	return wr.count, wr.err
}

func (w *writer) WriteModuleListItem(mi *javascript.ModuleItem) {
	if mi.ExportDeclaration != nil {
		w.WriteExportDeclaration(mi.ExportDeclaration)
	} else if mi.ImportDeclaration != nil {
		w.WriteImportDeclaration(mi.ImportDeclaration)
	} else if mi.StatementListItem != nil {
		w.WriteStatementListItem(mi.StatementListItem)
	}
}

func (w *writer) WriteExportDeclaration(ed *javascript.ExportDeclaration) {
	w.WriteString("export")
	if ed.FromClause != nil {
		if ed.ExportClause != nil {
			w.WriteExportClause(ed.ExportClause)
		} else {
			w.WriteString("*")
			if ed.ExportFromClause != nil {
				w.WriteString("as")
				w.WriteString(ed.ExportFromClause.Data)
			}
		}
		w.WriteFromClause(ed.FromClause)
	} else if ed.ExportClause != nil {
		w.WriteExportClause(ed.ExportClause)
	} else if ed.VariableStatement != nil {
		w.WriteVariableStatement(ed.VariableStatement)
	} else if ed.Declaration != nil {
		w.WriteDeclaration(ed.Declaration)
	} else if ed.DefaultFunction != nil {
		w.WriteString("default")
		w.WriteFunctionDeclaration(ed.DefaultFunction)
	} else if ed.DefaultClass != nil {
		w.WriteString("default")
		w.WriteClassDeclaration(ed.DefaultClass)
	} else if ed.DefaultAssignmentExpression != nil {
		w.WriteString("default")
		w.WriteAssignmentExpression(ed.DefaultAssignmentExpression)
	}
}

func (w *writer) WriteExportClause(ec *javascript.ExportClause) {
	w.WriteString("{")
	for n := range ec.ExportList {
		if n > 0 {
			w.WriteString(",")
		}
		w.WriteExportSpecifier(&ec.ExportList[n])
	}
	w.WriteString("}")
}

func (w *writer) WriteExportSpecifier(es *javascript.ExportSpecifier) {
	if es.IdentifierName == nil {
		w.err = ErrInvalidAST
		return
	}
	w.WriteString(es.IdentifierName.Data)
	if es.EIdentifierName != nil && es.EIdentifierName.Data != es.IdentifierName.Data {
		w.WriteString("as")
		w.WriteString(es.EIdentifierName.Data)
	}
}

func (w *writer) WriteFromClause(fc *javascript.FromClause) {
	w.WriteString("from")
	w.WriteString(fc.ModuleSpecifier.Data)
}

func (w *writer) WriteImportDeclaration(id *javascript.ImportDeclaration) {
	if id.ImportClause == nil && id.FromClause.ModuleSpecifier == nil {
		w.err = ErrInvalidAST
		return
	}
	w.WriteString("import")
	if id.ImportClause != nil {
		w.WriteImportClause(id.ImportClause)
		w.WriteFromClause(&id.FromClause)
	} else if id.FromClause.ModuleSpecifier != nil {
		w.WriteString(id.FromClause.ModuleSpecifier.Data)
	}
}

func (w *writer) WriteImportClause(ic *javascript.ImportClause) {
	if ic.ImportedDefaultBinding != nil {
		w.WriteString(ic.ImportedDefaultBinding.Data)
		if ic.NameSpaceImport != nil || ic.NamedImports != nil {
			w.WriteString(",")
		}
	}
	if ic.NameSpaceImport != nil {
		w.WriteString("*as")
		w.WriteString(ic.NameSpaceImport.Data)
	} else if ic.NamedImports != nil {
		w.WriteNamedImports(ic.NamedImports)
	}
}

func (w *writer) WriteNamedImports(ni *javascript.NamedImports) {
	w.WriteString("{")
	for n := range ni.ImportList {
		if n > 0 {
			w.WriteString(",")
		}
		w.WriteImportSpecifier(&ni.ImportList[n])
	}
	w.WriteString("}")
}

func (w *writer) WriteImportSpecifier(is *javascript.ImportSpecifier) {
	if is.IdentifierName == nil {
		w.err = ErrInvalidAST
		return
	}
	if is.IdentifierName != nil && is.IdentifierName.Data != is.ImportedBinding.Data {
		w.WriteString(is.IdentifierName.Data)
		w.WriteString("as")
	}
	w.WriteString(is.ImportedBinding.Data)
}

func (w *writer) WriteStatementListItem(si *javascript.StatementListItem) {
	if si.Statement != nil {
		w.WriteStatement(si.Statement)
	} else if si.Declaration != nil {
		w.WriteDeclaration(si.Declaration)
	}
}

func (w *writer) WriteStatement(s *javascript.Statement) {
	switch s.Type {
	case javascript.StatementNormal:
		if s.BlockStatement != nil {
			w.WriteBlock(s.BlockStatement)
		} else if s.VariableStatement != nil {
			w.WriteVariableStatement(s.VariableStatement)
		} else if s.ExpressionStatement != nil {
			w.WriteExpressionStatement(s.ExpressionStatement)
		} else if s.IfStatement != nil {
			w.WriteIfStatement(s.IfStatement)
		} else if s.IterationStatementDo != nil {
			w.WriteIterationStatementDo(s.IterationStatementDo)
		} else if s.IterationStatementWhile != nil {
			w.WriteIterationStatementWhile(s.IterationStatementWhile)
		} else if s.IterationStatementFor != nil {
			w.WriteIterationStatementFor(s.IterationStatementFor)
		} else if s.SwitchStatement != nil {
			w.WriteSwitchStatement(s.SwitchStatement)
		} else if s.WithStatement != nil {
			w.WriteWithStatement(s.WithStatement)
		} else if s.LabelIdentifier != nil {
			w.WriteString(s.LabelIdentifier.Data)
			w.WriteString(":")
			if s.LabelledItemFunction != nil {
				w.WriteFunctionDeclaration(s.LabelledItemFunction)
			} else if s.LabelledItemStatement != nil {
				w.WriteStatement(s.LabelledItemStatement)
			}
		} else if s.TryStatement != nil {
			w.WriteTryStatement(s.TryStatement)
		}
	case javascript.StatementContinue:
		w.WriteString("continue")
		if s.LabelIdentifier != nil {
			w.WriteString(s.LabelIdentifier.Data)
		}
	case javascript.StatementBreak:
		w.WriteString("break")
		if s.LabelIdentifier != nil {
			w.WriteString(s.LabelIdentifier.Data)
		}
	case javascript.StatementReturn:
		w.WriteString("return")
		if s.ExpressionStatement != nil {
			w.WriteExpressionStatement(s.ExpressionStatement)
		}
	case javascript.StatementThrow:
		if s.ExpressionStatement != nil {
			w.WriteString("throw")
			w.WriteExpressionStatement(s.ExpressionStatement)
		}
	case javascript.StatementDebugger:
		w.WriteString("debugger")
	}
}

func (w *writer) WriteExpressionStatement(e *javascript.Expression) {
	for n := range e.Expressions {
		if n > 0 {
			w.WriteString(",")
		}
		w.WriteAssignmentExpression(&e.Expressions[n])
	}
}

func (w *writer) WriteIfStatement(i *javascript.IfStatement) {
	w.WriteString("if(")
	w.WriteExpressionStatement(&i.Expression)
	w.WriteString(")")
	w.WriteStatement(&i.Statement)
	if i.ElseStatement != nil {
		w.WriteEOS()
		w.WriteString("else")
		w.WriteStatement(i.ElseStatement)
	}
}

func (w *writer) WriteIterationStatementDo(i *javascript.IterationStatementDo) {
	w.WriteString("do")
	w.WriteStatement(&i.Statement)
	w.WriteEOS()
	w.WriteString("while(")
	w.WriteExpressionStatement(&i.Expression)
	w.WriteString(")")
}

func (w *writer) WriteIterationStatementWhile(i *javascript.IterationStatementWhile) {
	w.WriteString("while(")
	w.WriteExpressionStatement(&i.Expression)
	w.WriteString(")")
	w.WriteStatement(&i.Statement)
}

func (w *writer) WriteIterationStatementFor(i *javascript.IterationStatementFor) {
	switch i.Type {
	case javascript.ForNormal:
		if i.InitVar != nil || i.InitLexical != nil || i.InitExpression != nil {
			w.err = ErrInvalidAST
			return
		}
	case javascript.ForNormalVar:
		if len(i.InitVar) == 0 {
			w.err = ErrInvalidAST
			return
		}
	case javascript.ForNormalLexicalDeclaration:
		if i.InitLexical == nil {
			w.err = ErrInvalidAST
			return
		}
	case javascript.ForNormalExpression:
		if i.InitExpression == nil {
			w.err = ErrInvalidAST
			return
		}
	case javascript.ForInLeftHandSide, javascript.ForOfLeftHandSide, javascript.ForAwaitOfLeftHandSide:
		if i.LeftHandSideExpression == nil {
			w.err = ErrInvalidAST
			return
		}
	case javascript.ForInVar, javascript.ForOfVar, javascript.ForAwaitOfVar, javascript.ForInLet, javascript.ForOfLet, javascript.ForAwaitOfLet, javascript.ForInConst, javascript.ForOfConst, javascript.ForAwaitOfConst:
		if i.ForBindingIdentifier == nil && i.ForBindingPatternObject == nil && i.ForBindingPatternArray == nil {
			w.err = ErrInvalidAST
			return
		}
	default:
		w.err = ErrInvalidAST
		return
	}
	switch i.Type {
	case javascript.ForInLeftHandSide, javascript.ForInVar, javascript.ForInLet, javascript.ForInConst:
		if i.In == nil {
			w.err = ErrInvalidAST
			return
		}
	case javascript.ForOfLeftHandSide, javascript.ForOfVar, javascript.ForOfLet, javascript.ForOfConst, javascript.ForAwaitOfLeftHandSide, javascript.ForAwaitOfVar, javascript.ForAwaitOfLet, javascript.ForAwaitOfConst:
		if i.Of == nil {
			w.err = ErrInvalidAST
			return
		}
	}
	switch i.Type {
	case javascript.ForAwaitOfLeftHandSide, javascript.ForAwaitOfVar, javascript.ForAwaitOfLet, javascript.ForAwaitOfConst:
		w.WriteString("for await(")
	default:
		w.WriteString("for(")
	}
	switch i.Type {
	case javascript.ForNormal:
		w.WriteString(";")
	case javascript.ForNormalVar:
		w.WriteString("var")
		for n := range i.InitVar {
			if n > 0 {
				w.WriteString(",")
			}
			w.WriteLexicalBinding((*javascript.LexicalBinding)(&i.InitVar[n]))
		}
		w.WriteString(";")
	case javascript.ForNormalLexicalDeclaration:
		w.WriteLexicalDeclaration(i.InitLexical)
		w.WriteString(";")
	case javascript.ForNormalExpression:
		w.WriteExpressionStatement(i.InitExpression)
		w.WriteString(";")
	case javascript.ForInLeftHandSide, javascript.ForOfLeftHandSide, javascript.ForAwaitOfLeftHandSide:
		w.WriteLeftHandSideExpression(i.LeftHandSideExpression)
	default:
		switch i.Type {
		case javascript.ForInVar, javascript.ForOfVar, javascript.ForAwaitOfVar:
			w.WriteString("var")
		case javascript.ForInLet, javascript.ForOfLet, javascript.ForAwaitOfLet:
			w.WriteString("let")
		case javascript.ForInConst, javascript.ForOfConst, javascript.ForAwaitOfConst:
			w.WriteString("const")
		}
		if i.ForBindingIdentifier != nil {
			w.WriteString(i.ForBindingIdentifier.Data)
		} else if i.ForBindingPatternObject != nil {
			w.WriteObjectBindingPattern(i.ForBindingPatternObject)
		} else {
			w.WriteArrayBindingPattern(i.ForBindingPatternArray)
		}
	}
	switch i.Type {
	case javascript.ForNormal, javascript.ForNormalVar, javascript.ForNormalLexicalDeclaration, javascript.ForNormalExpression:
		if i.Conditional != nil {
			w.WriteExpressionStatement(i.Conditional)
		}
		w.WriteString(";")
		if i.Afterthought != nil {
			w.WriteExpressionStatement(i.Afterthought)
		}
	case javascript.ForInLeftHandSide, javascript.ForInVar, javascript.ForInLet, javascript.ForInConst:
		w.WriteString("in")
		w.WriteExpressionStatement(i.In)
	case javascript.ForOfLeftHandSide, javascript.ForOfVar, javascript.ForOfLet, javascript.ForOfConst, javascript.ForAwaitOfLeftHandSide, javascript.ForAwaitOfVar, javascript.ForAwaitOfLet, javascript.ForAwaitOfConst:
		w.WriteString("of")
		w.WriteAssignmentExpression(i.Of)
	}
	w.WriteString(")")
	w.WriteStatement(&i.Statement)
}

func (w *writer) WriteLexicalBinding(lb *javascript.LexicalBinding) {
	if lb.BindingIdentifier != nil {
		w.WriteString(lb.BindingIdentifier.Data)
	} else if lb.ArrayBindingPattern != nil {
		w.WriteArrayBindingPattern(lb.ArrayBindingPattern)
	} else if lb.ObjectBindingPattern != nil {
		w.WriteObjectBindingPattern(lb.ObjectBindingPattern)
	} else {
		w.err = ErrInvalidAST
		return
	}
	if lb.Initializer != nil {
		w.WriteString("=")
		w.WriteAssignmentExpression(lb.Initializer)
	}
}

func (w *writer) WriteLexicalDeclaration(ld *javascript.LexicalDeclaration) {
	if len(ld.BindingList) == 0 {
		w.err = ErrInvalidAST
		return
	}
	if ld.LetOrConst == javascript.Let {
		w.WriteString("let")
	} else {
		w.WriteString("const")
	}
	for n := range ld.BindingList {
		if n > 0 {
			w.WriteString(",")
		}
		w.WriteLexicalBinding(&ld.BindingList[n])
	}
}

func (w *writer) WriteLeftHandSideExpression(lhs *javascript.LeftHandSideExpression) {
	if lhs.NewExpression != nil {
		w.WriteNewExpression(lhs.NewExpression)
	} else if lhs.CallExpression != nil {
		w.WriteCallExpression(lhs.CallExpression)
	} else if lhs.OptionalExpression != nil {
		w.WriteOptionalExpression(lhs.OptionalExpression)
	}
}

func (w *writer) WriteNewExpression(ne *javascript.NewExpression) {
	for i := uint(0); i < ne.News; i++ {
		w.WriteString("new")
	}
	w.WriteMemberExpression(&ne.MemberExpression)
}

func (w *writer) WriteMemberExpression(me *javascript.MemberExpression) {
	if me.MemberExpression != nil {
		if me.Arguments != nil {
			w.WriteString("new")
			w.WriteMemberExpression(me.MemberExpression)
			w.WriteArguments(me.Arguments)
		} else if me.Expression != nil {
			w.WriteMemberExpression(me.MemberExpression)
			w.WriteString("[")
			w.WriteExpressionStatement(me.Expression)
			w.WriteString("]")
		} else if me.IdentifierName != nil {
			w.WriteMemberExpression(me.MemberExpression)
			w.WriteString(".")
			w.WriteString(me.IdentifierName.Data)
		} else if me.PrivateIdentifier != nil {
			w.WriteMemberExpression(me.MemberExpression)
			w.WriteString(".")
			w.WriteString(me.PrivateIdentifier.Data)
		} else if me.TemplateLiteral != nil {
			w.WriteMemberExpression(me.MemberExpression)
			w.WriteTemplateLiteral(me.TemplateLiteral)
		}
	} else if me.PrimaryExpression != nil {
		w.WritePrimaryExpression(me.PrimaryExpression)
	} else if me.SuperProperty {
		if me.Expression != nil {
			w.WriteString("super[")
			w.WriteExpressionStatement(me.Expression)
			w.WriteString("]")
		} else if me.IdentifierName != nil {
			w.WriteString("super.")
			w.WriteString(me.IdentifierName.Data)
		}
	} else if me.NewTarget {
		w.WriteString("new.target")
	} else if me.ImportMeta {
		w.WriteString("import.meta")
	}
}

func (w *writer) WriteArguments(a *javascript.Arguments) {
	w.WriteString("(")
	for n := range a.ArgumentList {
		if n > 0 {
			w.WriteString(",")
		}
		w.WriteArgument(&a.ArgumentList[n])
	}
	w.WriteString(")")
}

func (w *writer) WriteArgument(a *javascript.Argument) {
	if a.Spread {
		w.WriteString("...")
	}
	w.WriteAssignmentExpression(&a.AssignmentExpression)
}

func (w *writer) WriteTemplateLiteral(tl *javascript.TemplateLiteral) {
	if tl.NoSubstitutionTemplate != nil {
		w.WriteString(tl.NoSubstitutionTemplate.Data)
	} else if tl.TemplateHead != nil && tl.TemplateTail != nil && len(tl.Expressions) == len(tl.TemplateMiddleList)+1 {
		w.WriteString(tl.TemplateHead.Data)
		w.WriteExpressionStatement(&tl.Expressions[0])
		for n := range tl.TemplateMiddleList {
			w.WriteString(tl.TemplateMiddleList[n].Data)
			w.WriteExpressionStatement(&tl.Expressions[n+1])
		}
		w.WriteString(tl.TemplateTail.Data)
	}
}

func (w *writer) WritePrimaryExpression(pe *javascript.PrimaryExpression) {
	if pe.This != nil {
		w.WriteString("this")
	} else if pe.IdentifierReference != nil {
		w.WriteString(pe.IdentifierReference.Data)
	} else if pe.Literal != nil {
		w.WriteString(pe.Literal.Data)
	} else if pe.ArrayLiteral != nil {
		w.WriteArrayLiteral(pe.ArrayLiteral)
	} else if pe.ObjectLiteral != nil {
		w.WriteObjectLiteral(pe.ObjectLiteral)
	} else if pe.FunctionExpression != nil {
		w.WriteFunctionDeclaration(pe.FunctionExpression)
	} else if pe.ClassExpression != nil {
		w.WriteClassDeclaration(pe.ClassExpression)
	} else if pe.TemplateLiteral != nil {
		w.WriteTemplateLiteral(pe.TemplateLiteral)
	} else if pe.ParenthesizedExpression != nil {
		w.WriteParenthesizedExpression(pe.ParenthesizedExpression)
	}
}

func (w *writer) WriteArrayLiteral(al *javascript.ArrayLiteral) {
	w.WriteString("[")
	for n := range al.ElementList {
		if n > 0 {
			w.WriteString(",")
		}
		w.WriteArrayElement(&al.ElementList[n])
	}
	w.WriteString("]")
}

func (w *writer) WriteArrayElement(ae *javascript.ArrayElement) {
	if ae.Spread {
		w.WriteString("...")
	}
	w.WriteAssignmentExpression(&ae.AssignmentExpression)
}

func (w *writer) WriteObjectLiteral(ol *javascript.ObjectLiteral) {
	w.WriteString("{")
	for n := range ol.PropertyDefinitionList {
		if n > 0 {
			w.WriteString(",")
		}
		w.WritePropertyDefinition(&ol.PropertyDefinitionList[n])
	}
	w.WriteString("}")
}

func (w *writer) WritePropertyDefinition(pd *javascript.PropertyDefinition) {
	if pd.AssignmentExpression != nil {
		if pd.PropertyName != nil {
			w.WritePropertyName(pd.PropertyName)
			var done bool
			if !pd.IsCoverInitializedName && pd.PropertyName.LiteralPropertyName != nil && pd.AssignmentExpression.ConditionalExpression != nil {
				c := javascript.UnwrapConditional(pd.AssignmentExpression.ConditionalExpression)
				if pe, ok := c.(*javascript.PrimaryExpression); ok && pe.IdentifierReference != nil {
					done = pe.IdentifierReference.Type == pd.PropertyName.LiteralPropertyName.Type && pe.IdentifierReference.Data == pd.PropertyName.LiteralPropertyName.Data
				}
			}
			if !done {
				if pd.IsCoverInitializedName {
					w.WriteString("=")
				} else {
					w.WriteString(":")
				}
				w.WriteAssignmentExpression(pd.AssignmentExpression)
			}
		} else {
			w.WriteString("...")
			w.WriteAssignmentExpression(pd.AssignmentExpression)
		}
	} else if pd.MethodDefinition != nil {
		w.WriteMethodDefinition(pd.MethodDefinition)
	}
}

func (w *writer) WritePropertyName(pn *javascript.PropertyName) {
	if pn.LiteralPropertyName != nil {
		w.WriteString(pn.LiteralPropertyName.Data)
	} else if pn.ComputedPropertyName != nil {
		w.WriteString("[")
		w.WriteAssignmentExpression(pn.ComputedPropertyName)
		w.WriteString("]")
	}
}

func (w *writer) WriteMethodDefinition(md *javascript.MethodDefinition) {
	switch md.Type {
	case javascript.MethodNormal:
	case javascript.MethodGenerator:
		w.WriteString("*")
	case javascript.MethodAsync:
		w.WriteString("async")
	case javascript.MethodAsyncGenerator:
		w.WriteString("async*")
	case javascript.MethodGetter:
		w.WriteString("get")
	case javascript.MethodSetter:
		w.WriteString("set")
	default:
		w.err = ErrInvalidAST
		return
	}
	w.WriteClassElementName(&md.ClassElementName)
	w.WriteFormalParameters(&md.Params)
	w.WriteBlock(&md.FunctionBody)
}

func (w *writer) WriteClassElementName(cem *javascript.ClassElementName) {
	if cem.PropertyName != nil {
		w.WritePropertyName(cem.PropertyName)
	} else if cem.PrivateIdentifier != nil {
		w.WriteString(cem.PrivateIdentifier.Data)
	}
}

func (w *writer) WriteFormalParameters(fp *javascript.FormalParameters) {
	w.WriteString("(")
	for n := range fp.FormalParameterList {
		if n > 0 {
			w.WriteString(",")
		}
		w.WriteBindingElement(&fp.FormalParameterList[n])
	}
	if fp.BindingIdentifier != nil || fp.ArrayBindingPattern != nil || fp.ObjectBindingPattern != nil {
		if len(fp.FormalParameterList) > 0 {
			w.WriteString(",")
		}
		if fp.BindingIdentifier != nil {
			w.WriteString("...")
			w.WriteString(fp.BindingIdentifier.Data)
		} else if fp.ArrayBindingPattern != nil {
			w.WriteString("...")
			w.WriteArrayBindingPattern(fp.ArrayBindingPattern)
		} else if fp.ObjectBindingPattern != nil {
			w.WriteString("...")
			w.WriteObjectBindingPattern(fp.ObjectBindingPattern)
		}
	}
	w.WriteString(")")
}

func (w *writer) WriteBindingElement(be *javascript.BindingElement) {
	if be.SingleNameBinding != nil {
		w.WriteString(be.SingleNameBinding.Data)
	} else if be.ArrayBindingPattern != nil {
		w.WriteArrayBindingPattern(be.ArrayBindingPattern)
	} else if be.ObjectBindingPattern != nil {
		w.WriteObjectBindingPattern(be.ObjectBindingPattern)
	} else {
		w.err = ErrInvalidAST
		return
	}
	if be.Initializer != nil {
		w.WriteString("=")
		w.WriteAssignmentExpression(be.Initializer)
	}
}

func (w *writer) WriteBlock(b *javascript.Block) {
	w.WriteString("{")
	for n := range b.StatementList {
		if n > 0 {
			w.WriteEOS()
		}
		w.WriteStatementListItem(&b.StatementList[n])
	}
	w.WriteString("}")
}

func (w *writer) WriteParenthesizedExpression(pe *javascript.ParenthesizedExpression) {
	w.WriteString("(")
	for n := range pe.Expressions {
		if n > 0 {
			w.WriteString(",")
		}
		w.WriteAssignmentExpression(&pe.Expressions[n])
	}
	w.WriteString(")")
}

func (w *writer) WriteCallExpression(ce *javascript.CallExpression) {
	if ce.SuperCall && ce.Arguments != nil {
		w.WriteString("super")
		w.WriteArguments(ce.Arguments)
	} else if ce.ImportCall != nil {
		w.WriteString("import(")
		w.WriteAssignmentExpression(ce.ImportCall)
		w.WriteString(")")
	} else if ce.MemberExpression != nil && ce.Arguments != nil {
		w.WriteMemberExpression(ce.MemberExpression)
		w.WriteArguments(ce.Arguments)
	} else if ce.CallExpression != nil {
		w.WriteCallExpression(ce.CallExpression)
		if ce.Arguments != nil {
			w.WriteArguments(ce.Arguments)
		} else if ce.Expression != nil {
			w.WriteString("[")
			w.WriteExpressionStatement(ce.Expression)
			w.WriteString("]")
		} else if ce.IdentifierName != nil {
			w.WriteString(".")
			w.WriteString(ce.IdentifierName.Data)
		} else if ce.TemplateLiteral != nil {
			w.WriteTemplateLiteral(ce.TemplateLiteral)
		} else if ce.PrivateIdentifier != nil {
			w.WriteString(".")
			w.WriteString(ce.PrivateIdentifier.Data)
		}
	}
}

func (w *writer) WriteOptionalExpression(oe *javascript.OptionalExpression) {
	if oe.MemberExpression != nil {
		w.WriteMemberExpression(oe.MemberExpression)
	} else if oe.CallExpression != nil {
		w.WriteCallExpression(oe.CallExpression)
	} else if oe.OptionalExpression != nil {
		w.WriteOptionalExpression(oe.OptionalExpression)
	}
	w.WriteOptionalChain(&oe.OptionalChain)
}

func (w *writer) WriteOptionalChain(oc *javascript.OptionalChain) {
	if oc.OptionalChain != nil {
		w.WriteOptionalChain(oc.OptionalChain)
	} else {
		w.WriteString("?.")
	}
	if oc.Arguments != nil {
		w.WriteArguments(oc.Arguments)
	} else if oc.Expression != nil {
		w.WriteString("[")
		w.WriteExpressionStatement(oc.Expression)
		w.WriteString("]")
	} else if oc.IdentifierName != nil {
		if oc.OptionalChain != nil {
			w.WriteString(".")
		}
		w.WriteString(oc.IdentifierName.Data)
	} else if oc.TemplateLiteral != nil {
		w.WriteTemplateLiteral(oc.TemplateLiteral)
	} else if oc.PrivateIdentifier != nil {
		if oc.OptionalChain != nil {
			w.WriteString(".")
		}
		w.WriteString(oc.PrivateIdentifier.Data)
	}
}

func (w *writer) WriteObjectBindingPattern(ob *javascript.ObjectBindingPattern) {
	w.WriteString("{")
	for n := range ob.BindingPropertyList {
		if n > 0 {
			w.WriteString(",")
		}
		w.WriteBindingProperty(&ob.BindingPropertyList[n])
	}
	if ob.BindingRestProperty != nil {
		if len(ob.BindingPropertyList) > 0 {
			w.WriteString(",")
		}
		w.WriteString("...")
		w.WriteString(ob.BindingRestProperty.Data)
	}
	w.WriteString("}")
}

func (w *writer) WriteBindingProperty(bp *javascript.BindingProperty) {
	if bp.PropertyName.LiteralPropertyName != nil && bp.BindingElement.SingleNameBinding != nil && bp.PropertyName.LiteralPropertyName.Data == bp.BindingElement.SingleNameBinding.Data {
		w.WriteBindingElement(&bp.BindingElement)
	} else {
		w.WritePropertyName(&bp.PropertyName)
		w.WriteString(":")
		w.WriteBindingElement(&bp.BindingElement)
	}
}

func (w *writer) WriteArrayBindingPattern(ab *javascript.ArrayBindingPattern) {
	w.WriteString("[")
	for n := range ab.BindingElementList {
		if n > 0 {
			w.WriteString(",")
		}
		w.WriteBindingElement(&ab.BindingElementList[n])
	}
	if ab.BindingRestElement != nil {
		if len(ab.BindingElementList) > 0 {
			w.WriteString(",")
		}
		w.WriteString("...")
		w.WriteBindingElement(ab.BindingRestElement)
	}
	w.WriteString("]")
}

func (w *writer) WriteSwitchStatement(s *javascript.SwitchStatement) {
	w.WriteString("switch(")
	w.WriteExpressionStatement(&s.Expression)
	w.WriteString("){")
	for n := range s.CaseClauses {
		if n > 0 {
			w.WriteEOS()
		}
		w.WriteCaseClause(&s.CaseClauses[n])
	}
	if len(s.DefaultClause) > 0 {
		if len(s.CaseClauses) > 0 {
			w.WriteEOS()
		}
		w.WriteString("default:")
		for n := range s.DefaultClause {
			if n > 0 {
				w.WriteEOS()
			}
			w.WriteStatementListItem(&s.DefaultClause[n])
		}
	}
	for n := range s.PostDefaultCaseClauses {
		if n > 0 || len(s.CaseClauses) > 0 || len(s.DefaultClause) > 0 {
			w.WriteEOS()
		}
		w.WriteCaseClause(&s.PostDefaultCaseClauses[n])
	}
	w.WriteString("}")
}

func (w *writer) WriteCaseClause(cc *javascript.CaseClause) {
	w.WriteString("case")
	w.WriteExpressionStatement(&cc.Expression)
	w.WriteString(":")
	for n := range cc.StatementList {
		if n > 0 {
			w.WriteEOS()
		}
		w.WriteStatementListItem(&cc.StatementList[n])
	}
}

func (w *writer) WriteWithStatement(ws *javascript.WithStatement) {
	w.WriteString("with(")
	w.WriteExpressionStatement(&ws.Expression)
	w.WriteString(")")
	w.WriteStatement(&ws.Statement)
}

func (w *writer) WriteTryStatement(t *javascript.TryStatement) {
	w.WriteString("try")
	w.WriteBlock(&t.TryBlock)
	if t.CatchBlock != nil {
		w.WriteString("catch(")
		if t.CatchParameterBindingIdentifier != nil {
			w.WriteString(t.CatchParameterBindingIdentifier.Data)
		} else if t.CatchParameterArrayBindingPattern != nil {
			w.WriteArrayBindingPattern(t.CatchParameterArrayBindingPattern)
		} else if t.CatchParameterObjectBindingPattern != nil {
			w.WriteObjectBindingPattern(t.CatchParameterObjectBindingPattern)
		}
		w.WriteString(")")
		w.WriteBlock(t.CatchBlock)
	}
	if t.FinallyBlock != nil {
		w.WriteString("finally")
		w.WriteBlock(t.FinallyBlock)
	}
}

func (w *writer) WriteVariableStatement(vs *javascript.VariableStatement) {
	if len(vs.VariableDeclarationList) > 0 {
		w.WriteString("var")
		for n := range vs.VariableDeclarationList {
			if n > 0 {
				w.WriteString(",")
			}
			w.WriteLexicalBinding((*javascript.LexicalBinding)(&vs.VariableDeclarationList[n]))
		}
	}
}

func (w *writer) WriteDeclaration(d *javascript.Declaration) {
	if d.ClassDeclaration != nil {
		w.WriteClassDeclaration(d.ClassDeclaration)
	} else if d.FunctionDeclaration != nil {
		w.WriteFunctionDeclaration(d.FunctionDeclaration)
	} else if d.LexicalDeclaration != nil {
		w.WriteLexicalDeclaration(d.LexicalDeclaration)
	}
}

func (w *writer) WriteFunctionDeclaration(f *javascript.FunctionDeclaration) {
	if f.Type == javascript.FunctionAsync || f.Type == javascript.FunctionAsyncGenerator {
		w.WriteString("async")
	}
	w.WriteString("function")
	if f.Type == javascript.FunctionGenerator || f.Type == javascript.FunctionAsyncGenerator {
		w.WriteString("*")
	} else if f.BindingIdentifier != nil {
		w.WriteString(f.BindingIdentifier.Data)
	}
	w.WriteFormalParameters(&f.FormalParameters)
	w.WriteBlock(&f.FunctionBody)
}

func (w *writer) WriteClassDeclaration(cd *javascript.ClassDeclaration) {
	w.WriteString("class")
	if cd.BindingIdentifier != nil {
		w.WriteString(cd.BindingIdentifier.Data)
	}
	if cd.ClassHeritage != nil {
		w.WriteString("extends")
		w.WriteLeftHandSideExpression(cd.ClassHeritage)
	}
	w.WriteString("{")
	for n := range cd.ClassBody {
		if n > 0 {
			w.WriteEOS()
		}
		w.WriteClassElement(&cd.ClassBody[n])
	}
	w.WriteString("}")
}

func (w *writer) WriteClassElement(ce *javascript.ClassElement) {
	if ce.Static {
		w.WriteString("static")
	}
	if ce.MethodDefinition != nil {
		w.WriteMethodDefinition(ce.MethodDefinition)
	} else if ce.FieldDefinition != nil {
		w.WriteFieldDefinition(ce.FieldDefinition)
	} else if ce.ClassStaticBlock != nil {
		w.WriteBlock(ce.ClassStaticBlock)
	}
}

func (w *writer) WriteFieldDefinition(fd *javascript.FieldDefinition) {
	w.WriteClassElementName(&fd.ClassElementName)
	if fd.Initializer != nil {
		w.WriteString("=")
		w.WriteAssignmentExpression(fd.Initializer)
	}
}

func (w *writer) WriteAssignmentExpression(ae *javascript.AssignmentExpression) {
	if ae.Yield && ae.AssignmentExpression != nil {
		w.WriteString("yield")
		if ae.Delegate {
			w.WriteString("*")
		}
		w.WriteAssignmentExpression(ae.AssignmentExpression)
	} else if ae.ArrowFunction != nil {
		w.WriteArrowFunction(ae.ArrowFunction)
	} else if ae.LeftHandSideExpression != nil && ae.AssignmentExpression != nil {
		var ao string
		switch ae.AssignmentOperator {
		case javascript.AssignmentAssign:
			ao = "="
		case javascript.AssignmentMultiply:
			ao = "*="
		case javascript.AssignmentDivide:
			ao = "/="
		case javascript.AssignmentRemainder:
			ao = "%="
		case javascript.AssignmentAdd:
			ao = "+="
		case javascript.AssignmentSubtract:
			ao = "-="
		case javascript.AssignmentLeftShift:
			ao = "<<="
		case javascript.AssignmentSignPropagatingRightShift:
			ao = ">>="
		case javascript.AssignmentZeroFillRightShift:
			ao = ">>>="
		case javascript.AssignmentBitwiseAND:
			ao = "&="
		case javascript.AssignmentBitwiseXOR:
			ao = "^="
		case javascript.AssignmentBitwiseOR:
			ao = "|="
		case javascript.AssignmentExponentiation:
			ao = "**="
		case javascript.AssignmentLogicalAnd:
			ao = "&&="
		case javascript.AssignmentLogicalOr:
			ao = "||="
		case javascript.AssignmentNullish:
			ao = "??="
		default:
			w.err = ErrInvalidAST
			return
		}
		w.WriteLeftHandSideExpression(ae.LeftHandSideExpression)
		w.WriteString(ao)
		w.WriteAssignmentExpression(ae.AssignmentExpression)
	} else if ae.AssignmentPattern != nil && ae.AssignmentExpression != nil && ae.AssignmentOperator == javascript.AssignmentAssign {
		w.WriteAssignmentPattern(ae.AssignmentPattern)
		w.WriteString("=")
		w.WriteAssignmentExpression(ae.AssignmentExpression)
	} else if ae.ConditionalExpression != nil {
		w.WriteConditionalExpression(ae.ConditionalExpression)
	}
}

func (w *writer) WriteArrowFunction(af *javascript.ArrowFunction) {
	if af.FunctionBody == nil && af.AssignmentExpression == nil || (af.BindingIdentifier == nil && af.FormalParameters == nil) {
		w.err = ErrInvalidAST
		return
	}
	if af.Async {
		w.WriteString("async")
	}
	if af.BindingIdentifier != nil {
		w.WriteString(af.BindingIdentifier.Data)
	} else if af.FormalParameters != nil {
		w.WriteFormalParameters(af.FormalParameters)
	}
	w.WriteString("=>")
	if af.FunctionBody != nil {
		w.WriteBlock(af.FunctionBody)
	} else {
		w.WriteAssignmentExpression(af.AssignmentExpression)
	}
}

func (w *writer) WriteAssignmentPattern(ap *javascript.AssignmentPattern) {
	if ap.ArrayAssignmentPattern != nil {
		w.WriteArrayAssignmentPattern(ap.ArrayAssignmentPattern)
	} else if ap.ObjectAssignmentPattern != nil {
		w.WriteObjectAssignmentPattern(ap.ObjectAssignmentPattern)
	}
}

func (w *writer) WriteArrayAssignmentPattern(aa *javascript.ArrayAssignmentPattern) {
	w.WriteString("[")
	for n := range aa.AssignmentElements {
		if n > 0 {
			w.WriteString(",")
		}
		w.WriteAssignmentElement(&aa.AssignmentElements[n])
	}
	if aa.AssignmentRestElement != nil {
		if len(aa.AssignmentElements) > 0 {
			w.WriteString(",")
		}
		w.WriteString("...")
		w.WriteLeftHandSideExpression(aa.AssignmentRestElement)
	}
	w.WriteString("]")
}

func (w *writer) WriteAssignmentElement(ae *javascript.AssignmentElement) {
	w.WriteDestructuringAssignmentTarget(&ae.DestructuringAssignmentTarget)
	if ae.Initializer != nil {
		w.WriteString("=")
		w.WriteAssignmentExpression(ae.Initializer)
	}
}

func (w *writer) WriteDestructuringAssignmentTarget(da *javascript.DestructuringAssignmentTarget) {
	if da.LeftHandSideExpression != nil {
		w.WriteLeftHandSideExpression(da.LeftHandSideExpression)
	} else if da.AssignmentPattern != nil {
		w.WriteAssignmentPattern(da.AssignmentPattern)
	}
}

func (w *writer) WriteObjectAssignmentPattern(oa *javascript.ObjectAssignmentPattern) {
	w.WriteString("{")
	for n := range oa.AssignmentPropertyList {
		if n > 0 {
			w.WriteString(",")
		}
		w.WriteAssignmentProperty(&oa.AssignmentPropertyList[n])
	}
	if oa.AssignmentRestElement != nil {
		if len(oa.AssignmentPropertyList) > 0 {
			w.WriteString(",")
		}
		w.WriteLeftHandSideExpression(oa.AssignmentRestElement)
	}
	w.WriteString("}")
}

func (w *writer) WriteAssignmentProperty(ap *javascript.AssignmentProperty) {
	w.WritePropertyName(&ap.PropertyName)
	if ap.DestructuringAssignmentTarget != nil {
		if ap.PropertyName.LiteralPropertyName != nil && ap.DestructuringAssignmentTarget.LeftHandSideExpression != nil && ap.DestructuringAssignmentTarget.LeftHandSideExpression.CallExpression == nil && ap.DestructuringAssignmentTarget.LeftHandSideExpression.OptionalExpression == nil && ap.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression != nil && ap.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.News == 0 && ap.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression != nil && ap.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference != nil && ap.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference.Data == ap.PropertyName.LiteralPropertyName.Data {
			w.err = ErrInvalidAST
			return
		}
		w.WriteString(":")
		w.WriteDestructuringAssignmentTarget(ap.DestructuringAssignmentTarget)
	}
	if ap.Initializer != nil {
		w.WriteString("=")
		w.WriteAssignmentExpression(ap.Initializer)
	}
}

func (w *writer) WriteConditionalExpression(ce *javascript.ConditionalExpression) {
	if ce.LogicalORExpression != nil {
		w.WriteLogicalORExpression(ce.LogicalORExpression)
	} else if ce.CoalesceExpression != nil {
		w.WriteCoalesceExpression(ce.CoalesceExpression)
	} else if ce.True != nil && ce.False != nil {
		w.WriteString("?")
		w.WriteAssignmentExpression(ce.True)
		w.WriteString(":")
		w.WriteAssignmentExpression(ce.False)
	}
}

func (w *writer) WriteLogicalORExpression(lo *javascript.LogicalORExpression) {
	if lo.LogicalORExpression != nil {
		w.WriteLogicalORExpression(lo.LogicalORExpression)
		w.WriteString("||")
	}
	w.WriteLogicalANDExpression(&lo.LogicalANDExpression)
}

func (w *writer) WriteLogicalANDExpression(la *javascript.LogicalANDExpression) {
	if la.LogicalANDExpression != nil {
		w.WriteLogicalANDExpression(la.LogicalANDExpression)
		w.WriteString("&&")
	}
	w.WriteBitwiseORExpression(&la.BitwiseORExpression)
}

func (w *writer) WriteBitwiseORExpression(bo *javascript.BitwiseORExpression) {
	if bo.BitwiseORExpression != nil {
		w.WriteBitwiseORExpression(bo.BitwiseORExpression)
		w.WriteString("|")
	}
	w.WriteBitwiseXORExpression(&bo.BitwiseXORExpression)
}

func (w *writer) WriteBitwiseXORExpression(bx *javascript.BitwiseXORExpression) {
	if bx.BitwiseXORExpression != nil {
		w.WriteBitwiseXORExpression(bx.BitwiseXORExpression)
		w.WriteString("^")
	}
	w.WriteBitwiseANDExpression(&bx.BitwiseANDExpression)
}

func (w *writer) WriteBitwiseANDExpression(ba *javascript.BitwiseANDExpression) {
	if ba.BitwiseANDExpression != nil {
		w.WriteBitwiseANDExpression(ba.BitwiseANDExpression)
		w.WriteString("&")
	}
	w.WriteEqualityExpression(&ba.EqualityExpression)
}

func (w *writer) WriteEqualityExpression(ee *javascript.EqualityExpression) {
	if ee.EqualityExpression != nil {
		var eo string
		switch ee.EqualityOperator {
		case javascript.EqualityEqual:
			eo = "=="
		case javascript.EqualityNotEqual:
			eo = "!="
		case javascript.EqualityStrictEqual:
			eo = "==="
		case javascript.EqualityStrictNotEqual:
			eo = "!=="
		default:
			w.err = ErrInvalidAST
			return
		}
		w.WriteEqualityExpression(ee.EqualityExpression)
		w.WriteString(eo)
	}
	w.WriteRelationalExpression(&ee.RelationalExpression)
}

func (w *writer) WriteRelationalExpression(re *javascript.RelationalExpression) {
	if re.PrivateIdentifier != nil {
		w.WriteString(re.PrivateIdentifier.Data)
		w.WriteString("in")
	} else if re.RelationalExpression != nil {
		var ro string
		switch re.RelationshipOperator {
		case javascript.RelationshipLessThan:
			ro = "<"
		case javascript.RelationshipGreaterThan:
			ro = ">"
		case javascript.RelationshipLessThanEqual:
			ro = "<="
		case javascript.RelationshipGreaterThanEqual:
			ro = ">="
		case javascript.RelationshipInstanceOf:
			ro = "instanceof"
		case javascript.RelationshipIn:
			ro = "in"
		default:
			w.err = ErrInvalidAST
			return
		}
		w.WriteRelationalExpression(re.RelationalExpression)
		w.WriteString(ro)
	}
	w.WriteShiftExpression(&re.ShiftExpression)
}

func (w *writer) WriteShiftExpression(se *javascript.ShiftExpression) {
	if se.ShiftExpression != nil {
		var so string
		switch se.ShiftOperator {
		case javascript.ShiftLeft:
			so = "<<"
		case javascript.ShiftRight:
			so = ">>"
		case javascript.ShiftUnsignedRight:
			so = ">>>"
		default:
			w.err = ErrInvalidAST
			return
		}
		w.WriteShiftExpression(se.ShiftExpression)
		w.WriteString(so)
	}
	w.WriteAdditiveExpression(&se.AdditiveExpression)
}

func (w *writer) WriteAdditiveExpression(ae *javascript.AdditiveExpression) {
	if ae.AdditiveExpression != nil {
		var ao string
		switch ae.AdditiveOperator {
		case javascript.AdditiveAdd:
			ao = "+"
		case javascript.AdditiveMinus:
			ao = "-"
		default:
			w.err = ErrInvalidAST
			return
		}
		w.WriteAdditiveExpression(ae.AdditiveExpression)
		w.WriteString(ao)
	}
	w.WriteMultiplicativeExpression(&ae.MultiplicativeExpression)
}

func (w *writer) WriteMultiplicativeExpression(me *javascript.MultiplicativeExpression) {
	if me.MultiplicativeExpression != nil {
		var mo string
		switch me.MultiplicativeOperator {
		case javascript.MultiplicativeMultiply:
			mo = "*"
		case javascript.MultiplicativeDivide:
			mo = "/"
		case javascript.MultiplicativeRemainder:
			mo = "%"
		default:
			w.err = ErrInvalidAST
			return
		}
		w.WriteMultiplicativeExpression(me.MultiplicativeExpression)
		w.WriteString(mo)
	}
	w.WriteExponentiationExpression(&me.ExponentiationExpression)
}

func (w *writer) WriteExponentiationExpression(ee *javascript.ExponentiationExpression) {
	if ee.ExponentiationExpression != nil {
		w.WriteExponentiationExpression(ee.ExponentiationExpression)
		w.WriteString("**")
	}
	w.WriteUnaryExpression(&ee.UnaryExpression)
}

func (w *writer) WriteUnaryExpression(ue *javascript.UnaryExpression) {
	for _, uo := range ue.UnaryOperators {
		switch uo {
		case javascript.UnaryDelete:
			w.WriteString("delete")
		case javascript.UnaryVoid:
			w.WriteString("void")
		case javascript.UnaryTypeOf:
			w.WriteString("typeof")
		case javascript.UnaryAdd:
			w.WriteString("+")
		case javascript.UnaryMinus:
			w.WriteString("-")
		case javascript.UnaryBitwiseNot:
			w.WriteString("~")
		case javascript.UnaryLogicalNot:
			w.WriteString("!")
		case javascript.UnaryAwait:
			w.WriteString("await")
		}
	}
	w.WriteUpdateExpression(&ue.UpdateExpression)
}

func (w *writer) WriteUpdateExpression(ue *javascript.UpdateExpression) {
	if ue.LeftHandSideExpression != nil {
		var uo string
		switch ue.UpdateOperator {
		case javascript.UpdatePostIncrement:
			uo = "++"
		case javascript.UpdatePostDecrement:
			uo = "--"
		case javascript.UpdatePreIncrement, javascript.UpdatePreDecrement:
			w.err = ErrInvalidAST
			return
		default:
		}
		w.WriteLeftHandSideExpression(ue.LeftHandSideExpression)
		if len(uo) > 0 {
			w.WriteString(uo)
		}
	} else if ue.UnaryExpression != nil {
		switch ue.UpdateOperator {
		case javascript.UpdatePreIncrement:
			w.WriteString("++")
		case javascript.UpdatePreDecrement:
			w.WriteString("--")
		default:
			w.err = ErrInvalidAST
			return
		}
		w.WriteUnaryExpression(ue.UnaryExpression)
	}
}

func (w *writer) WriteCoalesceExpression(ce *javascript.CoalesceExpression) {
	if ce.CoalesceExpressionHead != nil {
		w.WriteCoalesceExpression(ce.CoalesceExpressionHead)
		w.WriteString("??")
	}
	w.WriteBitwiseORExpression(&ce.BitwiseORExpression)
}

var ErrInvalidAST = errors.New("invalid AST")
