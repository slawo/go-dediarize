package dediarize_test

import (
	"strings"
	"testing"

	"github.com/slawo/go-dediarize/dediarize"
	"github.com/stretchr/testify/assert"
)

// func TestParserDoubleDip(t *testing.T) {
// 	p := dediarize.NewParser()
// 	_, err := p.Parse(strings.NewReader("{}"))
// 	assert.Nil(t, err)
// 	_, err = p.Parse(strings.NewReader("{}"))
// 	assert.EqualError(t, err, "diary already parsed")
// }

func TestParse(t *testing.T) {
	//t.Parallel()
	for _, tt := range []struct {
		name    string
		json    string
		wantErr string
		want    dediarize.Diary
	}{
		{
			name: "empty",
			json: `{}`,
		},
		{
			name: "valid",
			json: `{
			"segments": [
				{ "text": "Knock knock.", "speaker": "Ed", "start": "0.0", "end": "1.0" },
				{ "text": "Who's there?", "speaker": "Sam", "start": "1.0", "end": "2.0", "words": [{"word": "Who's"}, {"word": "there"}] },
				{ "text": "Go fmt.", "speaker": "Ed", "start": "2.0", "end": "3.0" },
				{ "text": "Go fmt who?", "speaker": "Sam", "start": "3.0", "end": "4.0" },
				{ "speaker": "Ed", "text": "Go fmt yourself!" }
			], "something-else": "not a segment", "language": "en"
		}
		`,
			want: dediarize.Diary{
				Segments: []dediarize.Segment{
					{Text: "Knock knock.", Speaker: "Ed", Start: "0.0", End: "1.0"},
					{Text: "Who's there?", Speaker: "Sam", Start: "1.0", End: "2.0"},
					{Text: "Go fmt.", Speaker: "Ed", Start: "2.0", End: "3.0"},
					{Text: "Go fmt who?", Speaker: "Sam", Start: "3.0", End: "4.0"},
					{Text: "Go fmt yourself!", Speaker: "Ed"},
				},
			},
		},
		{
			name: "parser ignores unknown fields",
			json: `{"ignored": "", "other": {"a":{}, "b": [[], {}], "c": "other", "d": 10, "e": false}, "another": [[], {}, "other", 10, false]}`,
		},
		{
			name: "segments ignores unknown fields",
			json: `{"segments": [ {"a":{}, "b": [[], {}], "c": "other", "d": 10, "e": false}, {"another": [[], {}, "other", 10, false]}]}`,
			want: dediarize.Diary{Segments: []dediarize.Segment{{}, {}}},
		},
		{
			name:    "first token array error",
			json:    `[]`,
			wantErr: "first token expected to be object, got [",
		},
		{
			name:    "first token string error",
			json:    `"s"`,
			wantErr: "first token expected to be object, got s",
		},
		{
			name:    "empty error",
			json:    "",
			wantErr: "failed to decode first token: EOF",
		},
		{
			name:    "first token error",
			json:    "]",
			wantErr: "failed to decode first token: invalid character ']' looking for beginning of value",
		},
		{
			name:    "wrong segments type",
			json:    `{"segments": "s"}`,
			wantErr: "failed to decode segments: segments expected to be an array, got string, s",
		},
		{
			name:    "invalid segments first token",
			json:    `{"segments": }`,
			wantErr: "failed to decode segments: invalid character '}' looking for beginning of value",
		},
		{
			name:    "invalid segments first token",
			json:    `{"segments": [} }`,
			wantErr: "failed to decode segments: segments failed to parse item token: invalid character '}' looking for beginning of value",
		},
		{
			name:    "invalid segment not object",
			json:    `{"segments": ["t"] }`,
			wantErr: "failed to decode segments: segments item expected to be an object, got string, t",
		},
		{
			name:    "segment malformed",
			json:    `{"segments": [{[]}] }`,
			wantErr: "failed to decode segments: segments failed to parse item, invalid character '['",
		},
		{
			name:    "segment malformed",
			json:    `{"segments": [{[]}] }`,
			wantErr: "failed to decode segments: segments failed to parse item, invalid character '['",
		},
		{
			name:    "segment field malformed",
			json:    `{"segments": [{"text":]}] }`,
			wantErr: "failed to decode segments: segments failed to parse item, invalid character ']' looking for beginning of value",
		},
		{
			name:    "segment field wrong type",
			json:    `{"segments": [{"text":[]}] }`,
			wantErr: "failed to decode segments: segments failed to parse item, expected string, got json.Delim, [",
		},
		{
			name:    "ignored invalid error",
			json:    `{"ignored": }}`,
			wantErr: "failed to decode token to skip: invalid character '}' looking for beginning of value",
		},
		{
			name:    "ignored invalid array malformed error",
			json:    `{"ignored": ,}`,
			wantErr: "failed to decode token to skip: invalid character ',' looking for beginning of value",
		},
		{
			name:    "ignored invalid array error",
			json:    `{"ignored": [}]}`,
			wantErr: "failed to skip token: failed to skip array: invalid character '}' looking for beginning of value",
		},
		{
			name:    "ignored invalid array malformed error",
			json:    `{"ignored": [,]}`,
			wantErr: "failed to skip token: failed to skip array: invalid character ',' looking for beginning of value",
		},
		{
			name:    "ignored object invalid label error",
			json:    `{"ignored": {"a": [{]}]}}`,
			wantErr: "failed to skip token: failed to skip object: failed to skip array: failed to skip object: invalid character ']'",
		},
		{
			name:    "ignored object invalid value error",
			json:    `{"ignored": {"a": [{"b":]}]}}`,
			wantErr: "failed to skip token: failed to skip object: failed to skip array: invalid character ']' looking for beginning of value",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := dediarize.Diary{}
			pw, err := dediarize.NewDiaryWriter(&d)
			assert.NoError(t, err)
			p, err := dediarize.NewParser(pw)
			assert.NoError(t, err)
			err = p.Parse(strings.NewReader(tt.json))
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualErrorf(t, err, tt.wantErr, "expected error %s", tt.wantErr)
			}
			assert.Equal(t, tt.want, d)
		})
	}
}
