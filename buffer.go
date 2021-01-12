package meta

import (
	"bytes"
	"strings"
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

func (s *snippetBuffer) Write(p []byte) (int, error) {
	if len(p) == 0 || s.IsFull() {
		return 0, nil
	}
	room := s.max - s.buf.Len()
	if len(p) > room {
		return s.buf.Write(p[:room])
	}
	return s.buf.Write(p)
}

func (s *snippetBuffer) WriteByte(b byte) error {
	if s.IsFull() {
		return nil
	}
	return s.buf.WriteByte(b)
}

func (s *snippetBuffer) IsEmpty() bool {
	return s.buf.Len() == 0
}

func (s *snippetBuffer) IsFull() bool {
	return s.max-s.buf.Len() <= 0
}

func (s *snippetBuffer) Reset() {
	s.buf.Reset()
}

func (s *snippetBuffer) String() string {
	return strings.TrimSuffix(s.buf.String(), " ")
}
