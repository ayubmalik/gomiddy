package gomiddy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	midi, err := Open("testdata/test2.mid")
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, midi)
	assert.Equal(t, midi.TrackCount, 2)
	assert.Equal(t, midi.Tempo, 124)
}
