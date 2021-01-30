package gomiddy

import (
	"testing"

	"github.com/deliveroo/assert-go"
)

func TestLoad(t *testing.T) {
	midi, err := Open("testdata/test2.mid")
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, midi)
	assert.Equal(t, midi.TrackCount, 4)
	assert.Equal(t, midi.Tempo, 96)
}
