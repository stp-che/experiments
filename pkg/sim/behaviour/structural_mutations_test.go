package behaviour

import (
	"experiments/pkg/test_helpers"
	"reflect"
	"testing"
)

func TestReductiveMutation(t *testing.T) {
	testCases := []struct {
		desc      string
		origBrain func() *Brain
		newBrain  *Brain
	}{
		{
			desc: "removes the first orphan manipulator",
			origBrain: func() *Brain {
				return &Brain{
					OuterAnalyzersCount: 1,
					OuterAnalyzerNet: OuterAnalyzerNet{
						&OuterAnalyzerLink{0, 0, 0, 10},
						&OuterAnalyzerLink{0, 1, 3, 10},
					},
					ManipulationSystem: ManipulationSystem{
						&Manipulator{AMove, [8]int8{1}},
						&Manipulator{AMove, [8]int8{2}},
						&Manipulator{AMove, [8]int8{3}},
						&Manipulator{AMove, [8]int8{4}},
					},
				}
			},
			newBrain: &Brain{
				OuterAnalyzersCount: 1,
				OuterAnalyzerNet: OuterAnalyzerNet{
					&OuterAnalyzerLink{0, 0, 0, 10},
					&OuterAnalyzerLink{0, 1, 2, 10},
				},
				ManipulationSystem: ManipulationSystem{
					&Manipulator{AMove, [8]int8{1}},
					&Manipulator{AMove, [8]int8{3}},
					&Manipulator{AMove, [8]int8{4}},
				},
			},
		},
		{
			desc: "removes the first orphan outer analyzer when there is no ophan manipulator",
			origBrain: func() *Brain {
				return &Brain{
					OuterAnalyzersCount: 3,
					OuterReceptor:       OuterReceptor{2, [][]uint8{{0, 1, 1}, {1, 2, 2}}},
					OuterAnalyzerNet: OuterAnalyzerNet{
						&OuterAnalyzerLink{0, 0, 0, 0},
						&OuterAnalyzerLink{2, 1, 1, 10},
					},
					ManipulationSystem: ManipulationSystem{
						&Manipulator{AMove, [8]int8{1}},
						&Manipulator{AMove, [8]int8{2}},
					},
				}
			},
			newBrain: &Brain{
				OuterAnalyzersCount: 2,
				OuterReceptor:       OuterReceptor{2, [][]uint8{{0, 1, 1}, {1, 1, 1}}},
				OuterAnalyzerNet: OuterAnalyzerNet{
					&OuterAnalyzerLink{0, 0, 0, 0},
					&OuterAnalyzerLink{2, 1, 1, 10},
				},
				ManipulationSystem: ManipulationSystem{
					&Manipulator{AMove, [8]int8{1}},
					&Manipulator{AMove, [8]int8{2}},
				},
			},
		},
		{
			desc: "removes the first effectless outer analyzer link when there is no ophan manipulator or outer analyzer",
			origBrain: func() *Brain {
				return &Brain{
					OuterAnalyzersCount: 1,
					OuterAnalyzerNet: OuterAnalyzerNet{
						&OuterAnalyzerLink{0, 0, 0, 0},
						&OuterAnalyzerLink{0, 1, 1, 10},
						&OuterAnalyzerLink{0, 1, 2, 0},
						&OuterAnalyzerLink{0, 2, 1, 10},
					},
					HealthAnalyzerNet: HealthAnalyzerNet{
						&HealthAnalyzerLink{0, 200, 0, 50},
						&HealthAnalyzerLink{100, 150, 3, 100},
					},
					ManipulationSystem: ManipulationSystem{
						&Manipulator{AMove, [8]int8{1}},
						&Manipulator{AMove, [8]int8{2}},
						&Manipulator{AMove, [8]int8{3}},
					},
				}
			},
			newBrain: &Brain{
				OuterAnalyzersCount: 1,
				OuterAnalyzerNet: OuterAnalyzerNet{
					&OuterAnalyzerLink{0, 0, 0, 0},
					&OuterAnalyzerLink{0, 1, 1, 10},
					&OuterAnalyzerLink{0, 2, 1, 10},
				},
				HealthAnalyzerNet: HealthAnalyzerNet{
					&HealthAnalyzerLink{0, 200, 0, 50},
					&HealthAnalyzerLink{100, 150, 2, 100},
				},
				ManipulationSystem: ManipulationSystem{
					&Manipulator{AMove, [8]int8{1}},
					&Manipulator{AMove, [8]int8{2}},
					&Manipulator{AMove, [8]int8{3}},
				},
			},
		},
		{
			desc: "removes the first effectless health analyzer link in other cases",
			origBrain: func() *Brain {
				return &Brain{
					OuterAnalyzersCount: 1,
					OuterAnalyzerNet: OuterAnalyzerNet{
						&OuterAnalyzerLink{0, 0, 0, 0},
						&OuterAnalyzerLink{0, 1, 1, 10},
						&OuterAnalyzerLink{0, 2, 1, 10},
					},
					HealthAnalyzerNet: HealthAnalyzerNet{
						&HealthAnalyzerLink{0, 200, 0, 50},
						&HealthAnalyzerLink{100, 150, 2, 100},
					},
					ManipulationSystem: ManipulationSystem{
						&Manipulator{AMove, [8]int8{1}},
						&Manipulator{AMove, [8]int8{2}},
					},
				}
			},
			newBrain: &Brain{
				OuterAnalyzersCount: 1,
				OuterAnalyzerNet: OuterAnalyzerNet{
					&OuterAnalyzerLink{0, 0, 0, 0},
					&OuterAnalyzerLink{0, 1, 1, 10},
					&OuterAnalyzerLink{0, 2, 1, 10},
				},
				HealthAnalyzerNet: HealthAnalyzerNet{
					&HealthAnalyzerLink{100, 150, 2, 100},
				},
				ManipulationSystem: ManipulationSystem{
					&Manipulator{AMove, [8]int8{1}},
					&Manipulator{AMove, [8]int8{2}},
				},
			},
		},
		{
			desc: "does not make change if nothing can be reduced",
			origBrain: func() *Brain {
				return &Brain{
					OuterAnalyzersCount: 1,
					OuterAnalyzerNet: OuterAnalyzerNet{
						&OuterAnalyzerLink{0, 0, 0, 10},
						&OuterAnalyzerLink{0, 1, 1, 10},
					},
					HealthAnalyzerNet: HealthAnalyzerNet{
						&HealthAnalyzerLink{100, 150, 1, 100},
					},
					ManipulationSystem: ManipulationSystem{
						&Manipulator{AMove, [8]int8{1}},
						&Manipulator{AMove, [8]int8{2}},
					},
				}
			},
			newBrain: &Brain{
				OuterAnalyzersCount: 1,
				OuterAnalyzerNet: OuterAnalyzerNet{
					&OuterAnalyzerLink{0, 0, 0, 10},
					&OuterAnalyzerLink{0, 1, 1, 10},
				},
				HealthAnalyzerNet: HealthAnalyzerNet{
					&HealthAnalyzerLink{100, 150, 1, 100},
				},
				ManipulationSystem: ManipulationSystem{
					&Manipulator{AMove, [8]int8{1}},
					&Manipulator{AMove, [8]int8{2}},
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.desc, func(t *testing.T) {
			brain := c.origBrain()
			newBrain := new(reductiveMutation).apply(brain)

			if !reflect.DeepEqual(newBrain, c.newBrain) {
				t.Errorf("Expected new brain to equal\n\n%s\n\ngot\n\n%s",
					test_helpers.Inspect(c.newBrain), test_helpers.Inspect(newBrain))
			}

			if !reflect.DeepEqual(brain, c.origBrain()) {
				t.Errorf("Expected original brain to not change but changed to:\n%s", test_helpers.Inspect(brain))
			}
		})
	}
}

func TestAddOuterAnalyzerMutation(t *testing.T) {
	testCases := []struct {
		desc                    string
		origOuterAnalyzersCount int
		newOuterAnalyzersCount  int
	}{
		{"increases OuterAnalyzerCount", 45, 46},
		{"does not increase OuterAnalyzersCount when maximum reached", 255, 255},
	}

	for _, c := range testCases {
		t.Run(c.desc, func(t *testing.T) {
			brain := testBrain()
			brain.OuterAnalyzersCount = c.origOuterAnalyzersCount
			newBrain := (mAddOuterAnalyzer{}).apply(brain)

			if newBrain.OuterAnalyzersCount != c.newOuterAnalyzersCount {
				t.Errorf("Expect new brain OuterAnalyzerCount to eq %d, got %d", c.newOuterAnalyzersCount, newBrain.OuterAnalyzersCount)
			}

			origBrain := testBrain()
			origBrain.OuterAnalyzersCount = c.origOuterAnalyzersCount
			if !reflect.DeepEqual(brain, origBrain) {
				t.Errorf("Expected original brain to not change but changed to:\n%s", test_helpers.Inspect(brain))
			}
		})
	}
}

func TestAddManipulator(t *testing.T) {
	testCases := []struct {
		desc                   string
		origManipulationSystem ManipulationSystem
		mutation               mAddManipulator
		newManipulationSystem  ManipulationSystem
	}{
		{
			desc: "adds manipulator",
			origManipulationSystem: ManipulationSystem{
				&Manipulator{AMove, [8]int8{1}},
				&Manipulator{AEat, [8]int8{2}},
			},
			mutation: mAddManipulator{AMove, [8]int8{3}},
			newManipulationSystem: ManipulationSystem{
				&Manipulator{AMove, [8]int8{1}},
				&Manipulator{AEat, [8]int8{2}},
				&Manipulator{AMove, [8]int8{3}},
			},
		},
		{
			desc:                   "does not add manipulator if maximum reached",
			origManipulationSystem: make(ManipulationSystem, 255),
			mutation:               mAddManipulator{},
			newManipulationSystem:  make(ManipulationSystem, 255),
		},
	}

	for _, c := range testCases {
		t.Run(c.desc, func(t *testing.T) {
			brain := testBrain()
			brain.ManipulationSystem = c.origManipulationSystem
			newBrain := c.mutation.apply(brain)

			if !reflect.DeepEqual(newBrain.ManipulationSystem, c.newManipulationSystem) {
				t.Errorf("Expected new brain ManipulationSystem to eq\n\n%s\n\ngot\n\n%s",
					test_helpers.Inspect(c.newManipulationSystem), test_helpers.Inspect(newBrain.ManipulationSystem))
			}

			origBrain := testBrain()
			origBrain.ManipulationSystem = c.origManipulationSystem
			if !reflect.DeepEqual(brain, origBrain) {
				t.Errorf("Expected original brain to not change but changed to:\n%s", test_helpers.Inspect(brain))
			}

			newBrain.ManipulationSystem = brain.ManipulationSystem
			if !reflect.DeepEqual(brain, newBrain) {
				t.Errorf("Expected only ManipulationSystem to be changed but changed something else as well:\n\norig brain\n\n%s\n\nnew brain\n\n%s",
					test_helpers.Inspect(brain), test_helpers.Inspect(newBrain))
			}
		})
	}
}
