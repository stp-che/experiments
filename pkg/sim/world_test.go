package sim

import (
	"experiments/pkg/sim/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPossibleDirection(t *testing.T) {
	w := &World{Cols: 10, Rows: 6}
	assert.Equal(t, true, w.PossibleDirection(25, core.UpLeft))
	assert.Equal(t, true, w.PossibleDirection(25, core.Up))
	assert.Equal(t, true, w.PossibleDirection(25, core.UpRight))
	assert.Equal(t, true, w.PossibleDirection(25, core.Right))
	assert.Equal(t, true, w.PossibleDirection(25, core.DownRight))
	assert.Equal(t, true, w.PossibleDirection(25, core.Down))
	assert.Equal(t, true, w.PossibleDirection(25, core.DownLeft))
	assert.Equal(t, true, w.PossibleDirection(25, core.Left))

	assert.Equal(t, false, w.PossibleDirection(1, core.Up))
	assert.Equal(t, false, w.PossibleDirection(1, core.UpRight))
	assert.Equal(t, false, w.PossibleDirection(1, core.UpLeft))
	assert.Equal(t, false, w.PossibleDirection(29, core.Right))
	assert.Equal(t, false, w.PossibleDirection(29, core.UpRight))
	assert.Equal(t, false, w.PossibleDirection(29, core.DownRight))
	assert.Equal(t, false, w.PossibleDirection(20, core.Left))
	assert.Equal(t, false, w.PossibleDirection(20, core.UpLeft))
	assert.Equal(t, false, w.PossibleDirection(20, core.DownLeft))
	assert.Equal(t, false, w.PossibleDirection(55, core.Down))
	assert.Equal(t, false, w.PossibleDirection(55, core.DownRight))
	assert.Equal(t, false, w.PossibleDirection(55, core.DownLeft))
}
