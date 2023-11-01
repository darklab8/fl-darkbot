package consoler

import "strings"

// Capturer of Cobra output
type Writer struct {
	sb *strings.Builder
}

func NewWriter() Writer {
	w := Writer{}
	w.sb = &(strings.Builder{})
	return w
}

func (w Writer) Write(p []byte) (n int, err error) {
	w.sb.Write(p)
	return len(p), nil
}

func (w Writer) String() string {
	return w.sb.String()
}
