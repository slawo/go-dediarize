package dediarize

import (
	"encoding/json"
	"fmt"
	"io"
)

type Parser struct {
	w Writer
}

func NewParser(w Writer) (*Parser, error) {
	if w == nil {
		return nil, fmt.Errorf("writer is nil")
	}
	return &Parser{w: w}, nil
}

func (p *Parser) Parse(s io.Reader) error {
	dec := json.NewDecoder(s)
	t, err := dec.Token()
	if err != nil {
		return fmt.Errorf("failed to decode first token: %w", err)
	}
	if t != json.Delim('{') {
		return fmt.Errorf("first token expected to be object, got %v", t)
	}
	for t, err = dec.Token(); err == nil; t, err = dec.Token() {
		switch t {
		case "segments":
			err = p.parseSegments(dec)
			if err != nil {
				return fmt.Errorf("failed to decode segments: %w", err)
			}
		case json.Delim('}'):
			return nil
		default:
			err = p.skipNext(dec)
		}
		if err != nil {
			return err
		}
	}
	return err
}

func (p *Parser) parseSegments(dec *json.Decoder) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if t != json.Delim('[') {
		return fmt.Errorf("segments expected to be an array, got %T, %v", t, t)
	}
	for t, err = dec.Token(); err == nil; t, err = dec.Token() {
		switch t {
		case json.Delim(']'):
			return nil
		case json.Delim('{'):
			err = p.parseSegment(dec)
		default:
			return fmt.Errorf("segments item expected to be an object, got %T, %v", t, t)
		}
		if err != nil {
			return fmt.Errorf("segments failed to parse item, %w", err)
		}
	}
	return fmt.Errorf("segments failed to parse item token: %w", err)
}

func (p *Parser) parseSegment(dec *json.Decoder) error {
	seg := Segment{}
	for t, err := dec.Token(); ; t, err = dec.Token() {
		if err != nil {
			return err
		}
		switch t {
		case "text":
			seg.Text, err = p.getString(dec)
		case "speaker":
			seg.Speaker, err = p.getString(dec)
		case "start":
			seg.Start, err = p.getString(dec)
		case "end":
			seg.End, err = p.getString(dec)
		case json.Delim('}'):
			return p.w.WriteSegment(&seg)
		default:
			err = p.skipNext(dec)
		}
		if err != nil {
			return err
		}
	}
}

func (p *Parser) getString(dec *json.Decoder) (string, error) {
	t, err := dec.Token()
	if err != nil {
		return "", err
	}
	switch s := t.(type) {
	case string:
		return s, err
	case float64:
		return fmt.Sprintf("%f", s), err
	}
	return "", fmt.Errorf("expected string, got %T, %v", t, t)
}

func (p *Parser) skipNext(dec *json.Decoder) error {
	t, err := dec.Token()
	if err != nil {
		return fmt.Errorf("failed to decode token to skip: %w", err)
	}
	switch t {
	case json.Delim('{'):
		err = p.skipObject(dec)
	case json.Delim('['):
		err = p.skipArray(dec)
	default:
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to skip token: %w", err)
	}
	return nil
}

func (p *Parser) skipObject(dec *json.Decoder) error {
	for t, err := dec.Token(); ; t, err = dec.Token() {
		if nil != err {
			return fmt.Errorf("failed to skip object: %w", err)
		}
		switch t {
		case json.Delim('}'):
			return nil
		}
		t, err := dec.Token()
		if err != nil {
			return err
		}
		switch t {
		case json.Delim('{'):
			err = p.skipObject(dec)
		case json.Delim('['):
			err = p.skipArray(dec)
		default:
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to skip object: %w", err)
		}
	}
}

func (p *Parser) skipArray(dec *json.Decoder) error {
	for t, err := dec.Token(); ; t, err = dec.Token() {
		if nil != err {
			return fmt.Errorf("failed to skip array: %w", err)
		}
		switch t {
		case json.Delim('{'):
			err = p.skipObject(dec)
		case json.Delim('['):
			err = p.skipArray(dec)
		case json.Delim(']'):
			return nil
		default:
		}
		if err != nil {
			return fmt.Errorf("failed to skip array: %w", err)
		}
	}
}
