package gomiddy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	midi, err := Open("testdata/test2.mid")
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, midi)
	assert.Equal(t, midi.TrackCount, 4)
	assert.Equal(t, midi.Tempo, 124)
}

func TestNoteString(t *testing.T) {
	x := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	fmt.Println("x", x)
}
