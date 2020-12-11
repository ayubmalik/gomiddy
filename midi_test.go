package gomiddy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {

	midi, err := Load("testdata/test.mid")
	if err != nil {
		t.Error(err)
	}

	require.NotNil(t, midi)
	assert.Equal(t, 4, midi.Tracks, "wrong track count")
	assert.Equal(t, 120, midi.Tempo, "wrong tempo")

}
