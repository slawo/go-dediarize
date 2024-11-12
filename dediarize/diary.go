package dediarize

type Diary struct {
	Segments []Segment `json:"segments"`
}

type Segment struct {
	Text    string `json:"text"`    // Text is the text of the segment
	Speaker string `json:"speaker"` // Speaker is the speaker of the segment
	Start   string `json:"start"`   // Start is the start time of the segment in seconds since the beginning of the audio
	End     string `json:"end"`     // End is the end time of the segment in seconds since the beginning of the audio
}
