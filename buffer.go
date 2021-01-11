package meta

import (
	"bytes"
)

type snippetBuffer struct {
	max int
	buf bytes.Buffer
}

func newSnippetBuffer(max int) *snippetBuffer {
	return &snippetBuffer{
		max: max,
	}
}

func (s *snippetBuffer) Write(p []byte) (n int, err error) {
	if len(p) == 0 || s.IsFull() {
		return 0, nil
	}
	if !s.IsEmpty() {
		s.buf.WriteByte(' ')
	}
	if len(p) > s.Room() {
		return s.buf.Write(p[:s.Room()])
	}
	if p[len(p)-1] == '\n' { // Skip the occasional newline that slips in
		return s.buf.Write(p[:len(p)-1])
	}
	return s.buf.Write(p)
}

func (s *snippetBuffer) Room() int {
	return s.max - s.buf.Len()
}

func (s *snippetBuffer) IsEmpty() bool {
	return s.buf.Len() == 0
}

func (s *snippetBuffer) IsFull() bool {
	return s.max-s.buf.Len() <= 0
}

func (s *snippetBuffer) String() string {
	return s.buf.String()
}
