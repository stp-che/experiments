package sim

import (
	"experiments/pkg/sim/behaviour"
	"experiments/pkg/sim/core"
	"reflect"
	"testing"
)

// B . W . .
// . F W . .
// . . B . .
// . . W F .
// . . . . .
var cfg testSimConfig = testSimConfig{
	H:       5,
	W:       5,
	Walls:   []int{2, 7, 17},
	Food:    []int{6, 18},
	BotsPos: []int{0, 12},
}

func TestLookAround(t *testing.T) {
	sim, brains := prepare(cfg)
	brains[0].SetVisionRange(1)
	brains[1].SetVisionRange(2)

	expectedRes := behaviour.OuterInput{
		Direction: core.Up,
		Signal:    []uint8{0, 0},
	}
	res := sim.Bots[0].LookAround()
	if !reflect.DeepEqual(res, expectedRes) {
		t.Errorf("expect result to eq %v, got %v", expectedRes, res)
	}

	expectedRes = behaviour.OuterInput{
		Direction: core.Up,
		Signal:    []uint8{uint8(RCBot), uint8(RCNone), uint8(RCWall), uint8(RCNone), uint8(RCFood), uint8(RCWall)},
	}
	res = sim.Bots[1].LookAround()
	if !reflect.DeepEqual(res, expectedRes) {
		t.Errorf("expect result to eq %v, got %v", expectedRes, res)
	}
}
