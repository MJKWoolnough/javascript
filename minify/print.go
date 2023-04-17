package minify

import (
	"errors"
	"io"

	"vimagination.zapto.org/javascript"
)

type writer struct {
	io.Writer
	count int64
	err   error
}

func (w *writer) WriteString(str string) {
	if w.err == nil {
		var n int
		n, w.err = io.WriteString(w.Writer, str)
		w.count += int64(n)
	}
}

func Print(w io.Writer, m *javascript.Module) (int64, error) {
	wr := writer{Writer: w}
	for _, mi := range m.ModuleListItems {
		wr.WriteModuleListItem(mi)
	}
	return wr.count, wr.err
}

func (w *writer) WriteModuleListItem(mi javascript.ModuleItem) {
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
		w.WriteString(";")
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
	w.WriteString(";")
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
		w.WriteString(";")
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
}

func (w *writer) WriteStatementListItem(si *javascript.StatementListItem) {
}

var ErrInvalidAST = errors.New("invalid AST")
