/*
Package gomiddy defines types and functions to read standard MIDI files.
The full spefication is defined at https://www.midi.org/specifications-old/item/standard-midi-files-smf.
*/
package gomiddy

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	noteOff      EventType = 0x80
	noteOn       EventType = 0x90
	polyPressure EventType = 0xA0
	controller   EventType = 0xB0
	progChange   EventType = 0xC0
	chanPressure EventType = 0xD0
	pithchBend   EventType = 0xE0
	sysEx        EventType = 0xF0
	meta         EventType = 0xFF
	unknown      EventType = 0x00
)

// MIDIFile represents a standard midi file.
type MIDIFile struct {
	Name       string
	TrackCount int
	Tempo      int
}

// Track represents a MIDI track
type Track struct {
	Name   string
	Events []Event
}

func (t Track) String() string {
	return fmt.Sprintf("Track name: %s, events: %d", t.Name, len(t.Events))
}

// Header represents MIDI header
type Header struct {
	trackCount int
	division   int
	format     int
}

// EventType represents MIDI event types.
type EventType uint8

func (e EventType) String() string {
	switch e {
	case noteOff:
		return "NOTE_OFF"
	case noteOn:
		return "NOTE_ON"
	case polyPressure:
		return "POLYPHONIC_PRESSURE"
	case controller:
		return "CONTROLLER"
	case progChange:
		return "PROGRAM_CHANGE"
	case chanPressure:
		return "CHANNEL_PRESSURE"
	case pithchBend:
		return "PITCH_BEND"
	case sysEx:
		return "SYSTEM_EXCLUSIVE"
	case meta:
		return "META"
	default:
		return "UNKNOWN"
	}
}

// Event represents a MIDI event.
type Event struct {
	delta     uint64
	eventType EventType
	channel   uint8
	velocity  uint8
	program   uint8
}

// MidiFileError represents an error when reading a MIDI file.
type MidiFileError struct {
	err error
	msg string
}

func (e MidiFileError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

// Open loads specified file as a MIDIFile.
func Open(file string) (*MIDIFile, error) {
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
	log.Println("numTracks", hdr.trackCount)
	log.Println("division", hdr.division)

	for {
		track, err := cr.track()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		fmt.Println("TRACK", track)
	}
	midi := &MIDIFile{Name: "TODO", TrackCount: hdr.trackCount, Tempo: 0}
	return midi, nil
}
