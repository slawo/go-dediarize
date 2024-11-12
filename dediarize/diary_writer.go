package dediarize

import "fmt"

type DiaryWriter struct {
	Diary *Diary `json:"diary"`
}

func NewDiaryWriter(d *Diary) (*DiaryWriter, error) {
	if d == nil {
		return nil, fmt.Errorf("diary is nil")
	}
	return &DiaryWriter{Diary: d}, nil
}

func (w *DiaryWriter) WriteSegment(s *Segment) error {

	w.Diary.Segments = append(w.Diary.Segments, Segment{
		Text:    s.Text,
		Speaker: s.Speaker,
		Start:   s.Start,
		End:     s.End,
	})

	return nil
}
