package dediarize

import (
	"fmt"
	"io"
	"os"
)

func LoadJsonFile(file string) (*Diary, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file '%s': %w", file, err)
	}
	defer jsonFile.Close()
	return ParseJson(jsonFile)
}

func ParseJson(s io.Reader) (*Diary, error) {
	d := Diary{}
	pw, err := NewDiaryWriter(&d)
	if err != nil {
		return nil, err
	}

	p, err := NewParser(pw)
	if err != nil {
		return nil, err
	}

	if err := p.Parse(s); err != nil {
		return nil, err
	}

	return &d, nil

	// dec := json.NewDecoder(s)

	// for t, err := dec.Token(); ; t, err = dec.Token() {
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	fmt.Printf("%T: %v\n", t, t)
	// }

	// // while the array contains values
	// for dec.More() {
	// 	var m Message
	// 	// decode an array value (Message)
	// 	err := dec.Decode(&m)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	fmt.Printf("%v: %v\n", m.Name, m.Text)
	// }

	// // read closing bracket
	// t, err = dec.Token()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%T: %v\n", t, t)

	return &d, nil
}
