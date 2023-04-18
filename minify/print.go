package minify

import (
	"errors"
	"io"

	"vimagination.zapto.org/javascript"
)

type writer struct {
	io.Writer
	count    int64
	err      error
	lastChar byte
}

func (w *writer) WriteString(str string) {
	if w.err == nil {
		var n int
		n, w.err = io.WriteString(w.Writer, str)
		w.count += int64(n)
		if len(str) > 0 {
			w.lastChar = str[len(str)-1]
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
			if ed.ExportClause != nil {
				w.WriteString("as ")
				w.WriteString(ed.ExportFromClause.Data)
				w.WriteString(" ")
			}
		}
		w.WriteFromClause(ed.FromClause)
	} else if ed.ExportClause != nil {
		w.WriteExportClause(ed.ExportClause)
	} else if ed.VariableStatement != nil {
		w.WriteString(" ")
		w.WriteVariableStatement(ed.VariableStatement)
	} else if ed.Declaration != nil {
		w.WriteString(" ")
		w.WriteDeclaration(ed.Declaration)
	} else if ed.DefaultFunction != nil {
		w.WriteString(" default ")
		w.WriteFunctionDeclaration(ed.DefaultFunction)
	} else if ed.DefaultClass != nil {
		w.WriteString(" default ")
		w.WriteClassDeclaration(ed.DefaultClass)
	} else if ed.DefaultAssignmentExpression != nil {
		w.WriteString(" default ")
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
		w.WriteString(" as ")
		w.WriteString(es.EIdentifierName.Data)
	}
}

func (w *writer) WriteFromClause(fc *javascript.FromClause) {
	w.WriteString("from ")
	w.WriteString(fc.ModuleSpecifier.Data)
}

func (w *writer) WriteImportDeclaration(id *javascript.ImportDeclaration) {
	if id.ImportClause == nil && id.FromClause.ModuleSpecifier == nil {
		return
	}
	w.WriteString("import")
	if id.ImportClause != nil {
		w.WriteImportClause(id.ImportClause)
		w.WriteFromClause(&id.FromClause)
	} else if id.FromClause.ModuleSpecifier != nil {
		w.WriteString(" ")
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
		w.WriteString("*as ")
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
		w.WriteString(" as ")
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
			w.WriteBlockStatement(s.BlockStatement)
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
		if s.LabelIdentifier == nil {
			w.WriteString("continue")
		} else {
			w.WriteString("continue ")
			w.WriteString(s.LabelIdentifier.Data)
		}
	case javascript.StatementBreak:
		if s.LabelIdentifier == nil {
			w.WriteString("break")
		} else {
			w.WriteString("break ")
			w.WriteString(s.LabelIdentifier.Data)
		}
	case javascript.StatementReturn:
		if s.ExpressionStatement == nil {
			w.WriteString("return")
		} else {
			w.WriteString("return ")
			w.WriteExpressionStatement(s.ExpressionStatement)
		}
	case javascript.StatementThrow:
		if s.ExpressionStatement != nil {
			w.WriteString("throw ")
			w.WriteExpressionStatement(s.ExpressionStatement)
		}
	case javascript.StatementDebugger:
		w.WriteString("debugger")
	}
}

func (w *writer) WriteBlockStatement(b *javascript.Block) {
	w.WriteString("{")
	for n := range b.StatementList {
		if n > 0 {
			w.WriteEOS()
		}
		w.WriteStatementListItem(&b.StatementList[n])
	}
	w.WriteString("}")
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
}

func (w *writer) WriteIterationStatementDo(i *javascript.IterationStatementDo) {
}

func (w *writer) WriteIterationStatementWhile(i *javascript.IterationStatementWhile) {
}

func (w *writer) WriteIterationStatementFor(i *javascript.IterationStatementFor) {
}

func (w *writer) WriteSwitchStatement(i *javascript.SwitchStatement) {
}

func (w *writer) WriteWithStatement(i *javascript.WithStatement) {
}

func (w *writer) WriteTryStatement(i *javascript.TryStatement) {
}

func (w *writer) WriteVariableStatement(vd *javascript.VariableStatement) {
}

func (w *writer) WriteDeclaration(d *javascript.Declaration) {
}

func (w *writer) WriteFunctionDeclaration(f *javascript.FunctionDeclaration) {
}

func (w *writer) WriteClassDeclaration(c *javascript.ClassDeclaration) {
}

func (w *writer) WriteAssignmentExpression(ae *javascript.AssignmentExpression) {
}

var ErrInvalidAST = errors.New("invalid AST")
