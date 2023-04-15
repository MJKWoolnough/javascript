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

	return wr.count, wr.err
}
