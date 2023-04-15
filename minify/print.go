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

func (w *writer) Write(p []byte) {
	if w.err == nil {
		var n int
		n, w.err = w.Writer.Write(p)
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
}

func (w *writer) WriteImportDeclaration(ed *javascript.ImportDeclaration) {
}

func (w *writer) WriteStatementListItem(si *javascript.StatementListItem) {
}
