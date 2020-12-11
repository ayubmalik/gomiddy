package gomiddy

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	mthd = "MThd"
	mtrk = "MTrk"
)

// MidiFile represents a standard midi file.
type MidiFile struct {
	Name   string
	Tracks int
	Tempo  int
}

// MidiFileError represents an error when reading a MIDI file.
type MidiFileError struct {
	err error
	msg string
}

// Load loads a MIDI file.
func Load(file string) (*MidiFile, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, MidiFileError{msg: "File error", err: err}
	}
	defer f.Close()

	br := bufio.NewReader(f)

	h, _ := chunkHeader(br)

	if h.chType != mthd {
		return nil, MidiFileError{msg: "Bad file format"}
	}

	// c, err := chunkData(br, h.len)
	// if err != nil {
	// 	return nil, MidiFileError{msg: "Chunk error", err: err}
	// }

	log.Println(h)

	c, _ := chunkData(br, h.len)

	format := binary.BigEndian.Uint16(c[:2])
	tracks := int(binary.BigEndian.Uint16(c[2:4]))
	pulses := int(binary.BigEndian.Uint16(c[4:]))

	log.Println("format", format)
	log.Println("tracks", tracks)
	log.Println("pulses", pulses)

	midi := &MidiFile{Name: "TODO", Tracks: tracks, Tempo: pulses}
	return midi, nil
}

type header struct {
	chType string
	len    uint32
}

type chunk []byte

func chunkHeader(r io.Reader) (*header, error) {
	buf := make([]byte, 4, 4)

	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	c := header{}
	c.chType = string(buf)

	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	c.len = binary.BigEndian.Uint32(buf)
	return &c, nil
}

func chunkData(r io.Reader, n uint32) (chunk, error) {
	buf := make([]byte, n, n)
	io.ReadFull(r, buf)
	return buf, nil
}

func (e MidiFileError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}
