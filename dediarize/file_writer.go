package dediarize

import (
	"fmt"
	"io"
)

type FileWriter struct {
	w io.Writer
}

func NewFileWriter(w io.Writer) (*FileWriter, error) {
	if w == nil {
		return nil, fmt.Errorf("missing writer")
	}
	return &FileWriter{w: w}, nil
}

func (w *FileWriter) WriteSegment(s *Segment) error {
	_, err := fmt.Fprintf(w.w, "[%s]: %s\n", s.Speaker, s.Text)
	return err
}
