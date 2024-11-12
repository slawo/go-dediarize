package dediarize_test

import (
	"testing"

	"github.com/slawo/go-dediarize/dediarize"
	"github.com/stretchr/testify/assert"
)

func TestNewDiaryWriterFailsWithMissingDiary(t *testing.T) {
	pw, err := dediarize.NewDiaryWriter(nil)
	assert.Nil(t, pw)
	assert.EqualError(t, err, "diary is nil")
}

func TestNewDiaryWriter(t *testing.T) {

	d := dediarize.Diary{}
	pw, err := dediarize.NewDiaryWriter(&d)
	assert.NoError(t, err)
	assert.NotNil(t, pw)
}

func TestDiaryWriterWriteSegment(t *testing.T) {
	d := dediarize.Diary{}
	pw, _ := dediarize.NewDiaryWriter(&d)
	err := pw.WriteSegment(&dediarize.Segment{Text: "Knock knock.", Speaker: "Ed", Start: "0.0", End: "1.0"})
	assert.NoError(t, err)
	assert.Equal(t, dediarize.Diary{Segments: []dediarize.Segment{{Text: "Knock knock.", Speaker: "Ed", Start: "0.0", End: "1.0"}}}, d)
}
