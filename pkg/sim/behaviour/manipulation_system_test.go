package behaviour

import (
	"experiments/pkg/sim/core"
	"testing"
)

func TestComputeIntention(t *testing.T) {
	ms := ManipulationSystem{
		{
			ActionType: AMove,
			DirValues:  [8]int8{1, 5, 0, -1, 0, 0, 0, 0},
		},
		{
			ActionType: AMove,
			DirValues:  [8]int8{2, 1, 0, 0, 0, 0, 0, 0},
		},
		{
			ActionType: AEat,
			DirValues:  [8]int8{1, 1, 1, 1, 1, 1, 2, 1},
		},
	}

	cases := []struct {
		Activation map[uint8]int16
		ActionType ActionType
		Direction  core.Direction
	}{
		{map[uint8]int16{0: 1, 1: 1, 2: 1}, AMove, core.Up},
		{map[uint8]int16{0: 1, 1: 6, 2: 6}, AMove, core.UpLeft},
		{map[uint8]int16{0: 1, 1: 2, 2: 3}, AMove, core.Up},
		{map[uint8]int16{0: 1, 1: 1, 2: 4}, AEat, core.DownLeft},
		{map[uint8]int16{0: -5, 1: 1, 2: 2}, AMove, core.Right},
	}

	for i, c := range cases {
		a := ms.ComputeIntention(c.Activation)
		if a.ActionType != c.ActionType || a.Direction != c.Direction {
			t.Errorf("[%d] Expected action Action{%v, %v}, got Action{%v, %v}", i, c.ActionType, c.Direction, a.ActionType, a.Direction)
		}
	}

	a := ms.ComputeIntention(map[uint8]int16{1: -1, 2: -1})
	if a != nil {
		t.Errorf("Expected action nil, got Action{%v, %v}", a.ActionType, a.Direction)
	}
}
