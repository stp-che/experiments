package behaviour

import (
	"experiments/pkg/sim/core"
	"reflect"
	"testing"
)

func testBrain() *Brain {
	return &Brain{
		OuterReceptor:       OuterReceptor{1, [][]uint8{{0, 1}}},
		OuterAnalyzersCount: 2,
		HealthAnalyzerNet: HealthAnalyzerNet{
			{0, 10, 4, 255},
			{0, 10, 5, 255},
		},
		OuterAnalyzerNet: OuterAnalyzerNet{
			{0, 0, 0, 1},
			{1, 0, 1, 1},
			{0, 1, 0, -100},
			{1, 1, 1, -100},
			{0, 2, 4, 5},
			{1, 2, 5, 5},
			{0, 3, 2, 1},
			{1, 3, 3, 1},
		},
		ManipulationSystem: []*Manipulator{
			{
				ActionType: AMove,
				DirValues:  [8]int8{1, 0, 0, 0, 0, 0, 0, 0},
			},
			{
				ActionType: AMove,
				DirValues:  [8]int8{0, 1, 0, 0, 0, 0, 0, 0},
			},
			{
				ActionType: AMove,
				DirValues:  [8]int8{0, 1, 2, 7, 10, 7, 2, 1},
			},
			{
				ActionType: AMove,
				DirValues:  [8]int8{1, 0, 1, 2, 7, 10, 7, 2},
			},
			{
				ActionType: AEat,
				DirValues:  [8]int8{1, 0, 0, 0, 0, 0, 0, 0},
			},
			{
				ActionType: AEat,
				DirValues:  [8]int8{0, 1, 0, 0, 0, 0, 0, 0},
			},
		},
	}
}

func brainsEqual(b1, b2 *Brain) bool {
	u1, u2 := b1.uuid, b1.origUuid
	defer (func() {
		b1.uuid, b1.origUuid = u1, u2
	})()

	b1.uuid, b1.origUuid = b2.uuid, b2.origUuid
	return reflect.DeepEqual(b1, b2)
}

type processTestCase struct {
	OuterInput [][]uint8
	InnerInput InnerInput
	Result     ProcessingResult
	Intention  *Intention
	EnergyCost int
}

func TestProcess(t *testing.T) {
	cases := []processTestCase{
		{
			OuterInput: [][]uint8{{1, 1}, {1, 1}, {0, 1}, {1, 1}},
			InnerInput: InnerInput{100},
			Intention:  &Intention{AMove, core.DownRight},
			EnergyCost: 43,
		},
		{
			OuterInput: [][]uint8{{2, 0}, {0, 0}, {0, 0}, {0, 0}},
			InnerInput: InnerInput{100},
			Intention:  &Intention{AEat, core.UpLeft},
			EnergyCost: 22,
		},
		{
			OuterInput: [][]uint8{{3, 3}, {0, 0}, {0, 0}, {0, 3}},
			InnerInput: InnerInput{100},
			Intention:  &Intention{AMove, core.DownRight},
			EnergyCost: 22,
		},
		{
			OuterInput: [][]uint8{{3, 0}, {3, 0}, {0, 0}, {0, 0}},
			InnerInput: InnerInput{100},
			Intention:  &Intention{AMove, core.Down},
			EnergyCost: 22,
		},
		{
			OuterInput: [][]uint8{{0, 1}, {1, 0}, {0, 0}, {3, 0}},
			InnerInput: InnerInput{100},
			Intention:  &Intention{AMove, core.Right},
			EnergyCost: 28,
		},
		{
			OuterInput: [][]uint8{{2, 3}, {0, 0}, {0, 0}, {0, 0}},
			InnerInput: InnerInput{100},
			Intention:  &Intention{AMove, core.Down},
			EnergyCost: 22,
		},
		{
			OuterInput: [][]uint8{{2, 3}, {0, 0}, {0, 0}, {0, 0}},
			InnerInput: InnerInput{9},
			Intention:  &Intention{AEat, core.UpLeft},
			EnergyCost: 24,
		},
	}

	brain := testBrain()

	for i, c := range cases {
		res := brain.Process(
			[]OuterInput{
				{core.Up, c.OuterInput[0]},
				{core.Right, c.OuterInput[1]},
				{core.Down, c.OuterInput[2]},
				{core.Left, c.OuterInput[3]},
			},
			c.InnerInput,
		)

		if !reflect.DeepEqual(res.Decision, c.Intention) {
			t.Errorf("[%d] Expected Decision to eq %v, got %v", i, c.Intention, res.Decision)
		}
		if res.EnergyCost != c.EnergyCost {
			t.Errorf("[%d] Expected EnergyCost to eq %d, got %d", i, c.EnergyCost, res.EnergyCost)
		}
	}
}
