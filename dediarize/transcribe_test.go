package dediarize_test

import (
	"os"
	"path"
	"testing"

	"github.com/slawo/go-dediarize/dediarize"
	"github.com/stretchr/testify/assert"
)

func TestTranscribeJsonFile(t *testing.T) {
	d := t.TempDir()
	for _, tt := range []struct {
		name   string
		input  string
		output string
		err    string
	}{
		{name: "empty", input: "{}", output: ""},
		{name: "empty", input: `{
			"segments": [
				{ "text": "Knock knock.", "speaker": "Ed", "start": "0.0", "end": "1.0" },
				{ "text": "Who's there?", "speaker": "Sam", "start": "1.0", "end": "2.0", "words": [{"word": "Who's"}, {"word": "there"}] },
				{ "text": "Go fmt.", "speaker": "Ed", "start": "2.0", "end": "3.0" },
				{ "text": "Go fmt who?", "speaker": "Sam", "start": "3.0", "end": "4.0" },
				{ "speaker": "Ed", "text": "Go fmt yourself!" }
			], "something-else": "not a segment", "language": "en"
		}
		`, output: `[Ed]: Knock knock.
[Sam]: Who's there?
[Ed]: Go fmt.
[Sam]: Go fmt who?
[Ed]: Go fmt yourself!
`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			inf := path.Join(d, tt.name+".json")
			ouf := path.Join(d, tt.name+".txt")

			err := os.WriteFile(inf, []byte(tt.input), 0644)
			assert.NoError(t, err)

			err = dediarize.TranscribeJsonFile(inf, ouf)
			assert.NoError(t, err)

			b, err := os.ReadFile(ouf)
			assert.NoError(t, err)
			assert.Equal(t, tt.output, string(b))
		})
	}
}

func TestTranscribeJsonFileErrorOnMissingInputFile(t *testing.T) {
	d := t.TempDir()
	inf := path.Join(d, "missingInput.json")
	ouf := path.Join(d, "missingInput.txt")

	err := dediarize.TranscribeJsonFile(inf, ouf)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestTranscribeJsonFileErrorOnMissingOutputFolder(t *testing.T) {
	d := t.TempDir()
	inf := path.Join(d, "missingOutput.json")
	ouf := path.Join(d, "dne", "missingOutput.txt")

	err := os.WriteFile(inf, []byte(`{ "segments": [ { "text": "Knock knock.", "speaker": "Ed", "start": "0.0", "end": "1.0" } ] }`), 0644)
	assert.NoError(t, err)

	err = dediarize.TranscribeJsonFile(inf, ouf)
	assert.ErrorIs(t, err, os.ErrNotExist)
}
