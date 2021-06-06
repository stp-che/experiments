package behaviour

import (
	"experiments/pkg/sim/core"
	"testing"
)

func TestComputeIntention(t *testing.T) {
	manipulators := []*manipulator{
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
		Activity   []int16
		ActionType ActionType
		Direction  core.Direction
	}{
		{[]int16{1, 1, 1}, AMove, core.Up},
		{[]int16{1, 6, 6}, AMove, core.UpLeft},
		{[]int16{1, 2, 3}, AMove, core.Up},
		{[]int16{1, 1, 4}, AEat, core.DownLeft},
		{[]int16{-5, 1, 2}, AMove, core.Right},
	}

	for i, c := range cases {
		a := computeIntention(manipulators, c.Activity)
		if a.ActionType != c.ActionType || a.Direction != c.Direction {
			t.Errorf("[%d] Expected action Action{%v, %v}, got Action{%v, %v}", i, c.ActionType, c.Direction, a.ActionType, a.Direction)
		}
	}

	a := computeIntention(manipulators, []int16{0, -1, -1})
	if a != nil {
		t.Errorf("Expected action nil, got Action{%v, %v}", a.ActionType, a.Direction)
	}
}
