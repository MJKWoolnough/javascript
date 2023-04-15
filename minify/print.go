package minify

import (
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
	}
	w.WriteString(";")
}

func (w *writer) WriteExportClause(ec *javascript.ExportClause) {
}

func (w *writer) WriteFromClause(fc *javascript.FromClause) {
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
	}
	w.WriteString(";")
}

func (w *writer) WriteImportClause(ic *javascript.ImportClause) {
}

func (w *writer) WriteStatementListItem(si *javascript.StatementListItem) {
}
