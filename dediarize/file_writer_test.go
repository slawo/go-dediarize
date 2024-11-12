package dediarize

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileWriterFailsWithMissingWriter(t *testing.T) {
	w, err := NewFileWriter(nil)
	assert.Nil(t, w)
	assert.EqualError(t, err, "missing writer")
}

func TestNewFileWriter(t *testing.T) {
	buf := bytes.NewBufferString("")
	w, err := NewFileWriter(buf)
	assert.NoError(t, err)
	assert.NotNil(t, w)
}

func TestWriteSegment(t *testing.T) {
	buf := bytes.NewBufferString("")
	w, _ := NewFileWriter(buf)

	s := &Segment{Text: "Knock knock.", Speaker: "Ed", Start: "0.0", End: "1.0"}

	err := w.WriteSegment(s)
	assert.NoError(t, err)
	assert.Equal(t, "[Ed]: Knock knock.\n", buf.String())
}
