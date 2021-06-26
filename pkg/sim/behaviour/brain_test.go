package behaviour

import (
	"experiments/pkg/sim/core"
	"reflect"
	"testing"
)

func testBrain() *Brain {
	return &Brain{
		Structure: BrainStructure{
			VisionRange:            1,
			OuterAnalyzersCount:    2,
			HealthAnalyzerNetSize:  2,
			OuterAnalyzerNetSize:   8,
			ManipulationSystemSize: 6,
		},
		Content: []uint8{
			// OuterReceptor
			0, 1,
			// HealthAnalyzerNet
			0, 10, 4, 255,
			0, 10, 5, 255,
			// OuterAnalyzerNet
			0, 0, 0, 129,
			1, 0, 1, 129,
			0, 1, 0, 28,
			1, 1, 1, 28,
			0, 2, 4, 133,
			1, 2, 5, 133,
			0, 3, 2, 129,
			1, 3, 3, 129,
			// ManipulationSystem
			uint8(AMove), 129, 128, 128, 128, 128, 128, 128, 128,
			uint8(AMove), 128, 129, 128, 128, 128, 128, 128, 128,
			uint8(AMove), 128, 129, 130, 135, 138, 135, 130, 129,
			uint8(AMove), 129, 128, 129, 130, 135, 138, 135, 130,
			uint8(AEat), 129, 128, 128, 128, 128, 128, 128, 128,
			uint8(AEat), 128, 129, 128, 128, 128, 128, 128, 128,
		},
	}
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
			EnergyCost: 46,
		},
		{
			OuterInput: [][]uint8{{2, 0}, {0, 0}, {0, 0}, {0, 0}},
			InnerInput: InnerInput{100},
			Intention:  &Intention{AEat, core.UpLeft},
			EnergyCost: 25,
		},
		{
			OuterInput: [][]uint8{{3, 3}, {0, 0}, {0, 0}, {0, 3}},
			InnerInput: InnerInput{100},
			Intention:  &Intention{AMove, core.DownRight},
			EnergyCost: 25,
		},
		{
			OuterInput: [][]uint8{{3, 0}, {3, 0}, {0, 0}, {0, 0}},
			InnerInput: InnerInput{100},
			Intention:  &Intention{AMove, core.Down},
			EnergyCost: 25,
		},
		{
			OuterInput: [][]uint8{{0, 1}, {1, 0}, {0, 0}, {3, 0}},
			InnerInput: InnerInput{100},
			Intention:  &Intention{AMove, core.Right},
			EnergyCost: 31,
		},
		{
			OuterInput: [][]uint8{{2, 3}, {0, 0}, {0, 0}, {0, 0}},
			InnerInput: InnerInput{100},
			Intention:  &Intention{AMove, core.Down},
			EnergyCost: 25,
		},
		{
			OuterInput: [][]uint8{{2, 3}, {0, 0}, {0, 0}, {0, 0}},
			InnerInput: InnerInput{9},
			Intention:  &Intention{AEat, core.UpLeft},
			EnergyCost: 27,
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

func TestNormalizeContent(t *testing.T) {
	b := &Brain{
		Structure: BrainStructure{
			VisionRange:            1,
			OuterAnalyzersCount:    2,
			HealthAnalyzerNetSize:  2,
			OuterAnalyzerNetSize:   3,
			ManipulationSystemSize: 2,
		},
		Content: []uint8{
			// OuterReceptor
			0, 45,
			// HealthAnalyzerNet
			0, 10, 1, 111,
			0, 10, 44, 222,
			// OuterAnalyzerNet
			200, 0, 0, 129,
			155, 0, 99, 255,
			0, 1, 100, 28,
			// ManipulationSystem
			uint8(AEat), 129, 128, 128, 128, 128, 128, 128, 128,
			uint8(AMove) + ActionTypesCount*5, 128, 129, 128, 128, 128, 255, 255, 255,
		},
	}

	b.NormalizeContent()

	expectedContent := []uint8{
		// OuterReceptor
		0, 1,
		// HealthAnalyzerNet
		0, 10, 1, 111,
		0, 10, 2, 222,
		// OuterAnalyzerNet
		0, 0, 0, 129,
		1, 0, 1, 255,
		0, 1, 0, 28,
		// ManipulationSystem
		uint8(AEat), 129, 128, 128, 128, 128, 128, 128, 128,
		uint8(AMove), 128, 129, 128, 128, 128, 255, 255, 255,
	}

	if !reflect.DeepEqual(b.Content, expectedContent) {
		t.Errorf("Expected Content to eq %v, got %v", expectedContent, b.Content)
	}
}
