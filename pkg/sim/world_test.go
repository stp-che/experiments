package sim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPossibleDirection(t *testing.T) {
	w := &World{Cols: 20, Rows: 15}
	assert.Equal(t, true, w.PossibleDirection(Pos{1, 1}, Up))
	assert.Equal(t, false, w.PossibleDirection(Pos{1, 0}, Up))
	assert.Equal(t, true, w.PossibleDirection(Pos{5, 13}, Down))
	assert.Equal(t, false, w.PossibleDirection(Pos{5, 14}, Down))
}
