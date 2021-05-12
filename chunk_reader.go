package gomiddy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	mThd = "MThd"
	mTrk = "MTrk"
)

const (
	msMin = 60000000
)

type chunk struct {
	cType string
	len   uint32
	data  []byte
}

type chunkReader struct {
	reader io.Reader
}

func (cr chunkReader) header() (*Header, error) {
	chunk, err := cr.chunk()
	if err != nil {
		return nil, err
	}

	if chunk.cType != mThd {
		return nil, MidiFileError{msg: "Bad file format, first chunk type not MThd"}
	}
	format := int(binary.BigEndian.Uint16(chunk.data[:2]))
	tracks := int(binary.BigEndian.Uint16(chunk.data[2:4]))
	division := int(binary.BigEndian.Uint16(chunk.data[4:]))

	mask := division
	mask = mask >> 15

	// TODO use SMPTE division
	if mask > 0 {
		fmt.Println("WARNING time is in SMPTE!")
	}
	return &Header{trackCount: tracks, division: division, format: format}, nil
}

func (cr chunkReader) track() (*Track, error) {
	chunk, err := cr.chunk()
	if err != nil {
		return nil, err
	}
	if chunk.cType != mTrk {
		return nil, MidiFileError{msg: "Bad file format, chunk type not MTrk"}
	}

	track := Track{Name: "untitled"}
	events := make([]Event, 0)

	var lastStatus byte
	r := bytes.NewReader(chunk.data)
	for {
		d, err := binary.ReadUvarint(r)
		if err == io.EOF {
			break
		}
		event := Event{delta: d, eventType: unknown}
		statusByte, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if statusByte&0x80 == 0 {
			r.UnreadByte()
		} else {
			lastStatus = statusByte
		}

		event.eventType = EventType(lastStatus)
		event.channel = (lastStatus & 0x0F)

		switch event.eventType {
		case noteOn, noteOff:
			_, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			vel, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			event.velocity = uint8(vel)

		case progChange:
			prg, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			event.program = uint8(prg)
		case polyPressure, controller, pithchBend:
			r.ReadByte()
			r.ReadByte()
		case chanPressure:
			r.ReadByte()
		case sysEx:
			n, _ := r.ReadByte()
			io.CopyN(io.Discard, r, int64(n))
		case meta:
			mtype, _ := r.ReadByte()
			n, _ := binary.ReadUvarint(r)
			buf := make([]byte, n)
			r.Read(buf)

			if mtype == 0x3 {
				track.Name = string(buf)
			}

			if mtype == 0x1 {
				fmt.Println("device", string(buf))
			}

			if mtype == 0x51 {
				ms := int(uint(buf[2]) | uint(buf[1])<<8 | uint(buf[0])<<16)
				bpm := msMin / ms
				track.tempo = bpm
				// TODO set tempo and tick intervals
				//60000 / (BPM * PPQ)
			}

		default:
			return nil, MidiFileError{msg: fmt.Sprintf("Unknown MIDI type %d", lastStatus)}
		}
		events = append(events, event)
	}

	track.Events = events
	return &track, nil
}

func (cr chunkReader) chunk() (*chunk, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(cr.reader, buf)
	if err != nil {
		return nil, err
	}

	chunk := chunk{}
	chunk.cType = string(buf)
	_, err = io.ReadFull(cr.reader, buf)
	if err != nil {
		return nil, err
	}

	chunk.len = binary.BigEndian.Uint32(buf)
	chunk.data = make([]byte, chunk.len)
	_, err = io.ReadFull(cr.reader, chunk.data)
	if err != nil {
		return nil, err
	}
	return &chunk, nil
}
