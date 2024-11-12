package dediarize

type Writer interface {
	WriteSegment(*Segment) error
}
