package gomiddy

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	mThd = "MThd"
	mTrk = "MTrk"
)

type chunk struct {
	cType string
	len   uint32
	data  []byte
}

type chunkReader struct {
	reader io.Reader
}

func (cr chunkReader) header() (*MHeader, error) {
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
	if mask > 0 {
		fmt.Println("WARNING time is in SMPTE!")
	}
	return &MHeader{numTracks: tracks, division: division, format: format}, nil
}

func (cr chunkReader) track() (*MTrack, error) {
	chunk, err := cr.chunk()
	if err != nil {
		return nil, err
	}
	if chunk.cType != mTrk {
		return nil, MidiFileError{msg: "Bad file format, chunk type not MTrk"}
	}

	track := MTrack{name: "untitled"}
	events := make([]MEvent, 0)
	var lastStatus byte

	r := bytes.NewReader(chunk.data)
	for {
		d, err := binary.ReadUvarint(r)
		if err == io.EOF {
			break
		}
		event := MEvent{delta: d, eventType: "midi"}
		statusByte, err := r.ReadByte()

		if statusByte&0x80 == 0 {
			r.UnreadByte()
		} else {
			lastStatus = statusByte
		}
		msg := (lastStatus & 0xF0) >> 4

		switch msg {
		case 0x2, 0x3, 0x4, 0x5, 0x6:
			event.eventType = "unknown"
			r.ReadByte()
		case 0x8, 0x9, 0xA, 0xB, 0xE:
			r.ReadByte()
			r.ReadByte()
		case 0xC, 0xD:
			r.ReadByte()
		case 0xF:
			event.eventType = "meta"
			b, _ := r.ReadByte()
			meta := int(b)
			n, _ := binary.ReadUvarint(r)
			buf := make([]byte, n, n)
			r.Read(buf)

			if meta == 3 {
				//name := string(buf)
				track.name = string(buf)
				fmt.Println("META, N, BUF =", meta, n, buf)
			}

		default:
			return nil, errors.New("unknown msg type " + string(msg))
		}
		events = append(events, event)
	}

	track.events = events
	return &track, nil
}

func (cr chunkReader) chunk() (*chunk, error) {
	buf := make([]byte, 4, 4)
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
	chunk.data = make([]byte, chunk.len, chunk.len)
	_, err = io.ReadFull(cr.reader, chunk.data)
	if err != nil {
		return nil, err
	}
	return &chunk, nil
}

func decodeVarInt(buf []byte) (x uint32, n int) {
	if len(buf) < 1 {
		return 0, 0
	}

	if buf[0] <= 0x80 {
		return uint32(buf[0]), 1
	}

	var b byte
	for _, b = range buf {
		x = x << 7
		x |= uint32(b) & 0x7F
		n++
		if b&0x80 == 0 {
			return x, n
		}
	}
	return x, n
}
