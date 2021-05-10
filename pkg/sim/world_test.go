package sim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPossibleDirection(t *testing.T) {
	w := &World{Cols: 10, Rows: 6}
	assert.Equal(t, true, w.PossibleDirection(25, UpLeft))
	assert.Equal(t, true, w.PossibleDirection(25, Up))
	assert.Equal(t, true, w.PossibleDirection(25, UpRight))
	assert.Equal(t, true, w.PossibleDirection(25, Right))
	assert.Equal(t, true, w.PossibleDirection(25, DownRight))
	assert.Equal(t, true, w.PossibleDirection(25, Down))
	assert.Equal(t, true, w.PossibleDirection(25, DownLeft))
	assert.Equal(t, true, w.PossibleDirection(25, Left))

	assert.Equal(t, false, w.PossibleDirection(1, Up))
	assert.Equal(t, false, w.PossibleDirection(1, UpRight))
	assert.Equal(t, false, w.PossibleDirection(1, UpLeft))
	assert.Equal(t, false, w.PossibleDirection(29, Right))
	assert.Equal(t, false, w.PossibleDirection(29, UpRight))
	assert.Equal(t, false, w.PossibleDirection(29, DownRight))
	assert.Equal(t, false, w.PossibleDirection(20, Left))
	assert.Equal(t, false, w.PossibleDirection(20, UpLeft))
	assert.Equal(t, false, w.PossibleDirection(20, DownLeft))
	assert.Equal(t, false, w.PossibleDirection(55, Down))
	assert.Equal(t, false, w.PossibleDirection(55, DownRight))
	assert.Equal(t, false, w.PossibleDirection(55, DownLeft))
}
