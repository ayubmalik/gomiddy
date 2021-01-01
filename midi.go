package gomiddy

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// MIDIFile represents a standard midi file.
type MIDIFile struct {
	Name      string
	NumTracks int
	Tempo     int
}

// MTrack represents a MIDI track
type MTrack struct {
	name   string
	events []MEvent
}

// MHeader represents MIDI header
type MHeader struct {
	numTracks int
	division  int
	format    int
}

// MEvent represents a MIDI event
type MEvent struct {
	delta     uint64
	eventType string
}

// Load loads specified file as a MIDIFile.
func Load(file string) (*MIDIFile, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, MidiFileError{msg: "File error", err: err}
	}
	defer f.Close()

	cr := chunkReader{reader: bufio.NewReader(f)}

	hdr, err := cr.header()
	if err != nil {
		return nil, MidiFileError{msg: "Bad Header", err: err}
	}

	log.Println("format", hdr.format)
	log.Println("numTracks", hdr.numTracks)
	log.Println("division", hdr.division)

	for {
		track, err := cr.track()
		if err != nil {
			fmt.Println("got error", err)
			if err == io.EOF {
				break
			}
			return nil, err
		}

		fmt.Println("TRACK", track)
	}
	midi := &MIDIFile{Name: "TODO", NumTracks: hdr.numTracks, Tempo: 0}
	return midi, nil
}

// MidiFileError represents an error when reading a MIDI file.
type MidiFileError struct {
	err error
	msg string
}

func (e MidiFileError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}
