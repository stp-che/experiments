package sim

import "testing"

func TestGenomeComputeAction(t *testing.T) {
	controllers := []*ControllerGene{
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
		Direction  Direction
	}{
		{[]int16{1, 1, 1}, AMove, Up},
		{[]int16{1, 6, 6}, AMove, UpLeft},
		{[]int16{1, 2, 3}, AMove, Up},
		{[]int16{1, 1, 4}, AEat, DownLeft},
		{[]int16{-5, 1, 2}, AMove, Right},
	}

	for i, c := range cases {
		a := genomeComputeAction(controllers, c.Activity)
		if a.Type != c.ActionType || a.Direction != c.Direction {
			t.Errorf("[%d] Expected action Action{%v, %v}, got Action{%v, %v}", i, c.ActionType, c.Direction, a.Type, a.Direction)
		}
	}

	a := genomeComputeAction(controllers, []int16{0, -1, -1})
	if a != nil {
		t.Errorf("Expected action nil, got Action{%v, %v}", a.Type, a.Direction)
	}
}
