package behaviour

import (
	"experiments/pkg/sim/core"
	"experiments/pkg/test_helpers"
	"reflect"
	"testing"
)

func TestChangeManipulatorValue(t *testing.T) {
	origMSys := func() ManipulationSystem {
		return ManipulationSystem{
			&Manipulator{
				ActionType: AMove,
				DirValues:  [8]int8{10, 0, 0, 120, 0, -100, 0, -5},
			},
		}
	}

	testCases := []struct {
		desc     string
		mutation mChangeManipulatorValue
		newMSys  ManipulationSystem
	}{
		{
			desc:     "increasing",
			mutation: mChangeManipulatorValue{0, core.Left, 20},
			newMSys: ManipulationSystem{
				&Manipulator{
					ActionType: AMove,
					DirValues:  [8]int8{10, 0, 0, 120, 0, -100, 0, 15},
				},
			},
		},
		{
			desc:     "decreasing",
			mutation: mChangeManipulatorValue{0, core.UpLeft, -20},
			newMSys: ManipulationSystem{
				&Manipulator{
					ActionType: AMove,
					DirValues:  [8]int8{-10, 0, 0, 120, 0, -100, 0, -5},
				},
			},
		},
		{
			desc:     "positive overflow",
			mutation: mChangeManipulatorValue{0, core.Right, 20},
			newMSys: ManipulationSystem{
				&Manipulator{
					ActionType: AMove,
					DirValues:  [8]int8{10, 0, 0, 127, 0, -100, 0, -5},
				},
			},
		},
		{
			desc:     "negative overflow",
			mutation: mChangeManipulatorValue{0, core.Down, -50},
			newMSys: ManipulationSystem{
				&Manipulator{
					ActionType: AMove,
					DirValues:  [8]int8{10, 0, 0, 120, 0, -128, 0, -5},
				},
			},
		},
		{
			desc:     "when wrong direction",
			mutation: mChangeManipulatorValue{0, 9, 10},
			newMSys:  origMSys(),
		},
		{
			desc:     "when manipulator does not exist",
			mutation: mChangeManipulatorValue{1, core.Left, 10},
			newMSys:  origMSys(),
		},
	}

	for _, c := range testCases {
		t.Run(c.desc, func(t *testing.T) {
			brain := testBrain()
			brain.ManipulationSystem = origMSys()

			newBrain := c.mutation.apply(brain)

			if !reflect.DeepEqual(newBrain.ManipulationSystem, c.newMSys) {
				t.Errorf("[Expected new ManipulationSystem to equal %s, got %s",
					test_helpers.Inspect(c.newMSys), test_helpers.Inspect(newBrain.ManipulationSystem))
			}
			if !reflect.DeepEqual(brain.ManipulationSystem, origMSys()) {
				t.Errorf("Expected original brain to not change; but ManipulationSystem changed to %v",
					test_helpers.Inspect(brain.ManipulationSystem))
			}
			newBrain.ManipulationSystem = brain.ManipulationSystem
			if !reflect.DeepEqual(brain, newBrain) {
				t.Errorf("Expected the rest of new brain to be a copy of original brain\n\nOriginal brain:\n%v\n\nNew brain:\n%v", brain, newBrain)
			}
		})
	}
}

func TestRandomChangeManipulatorValue(t *testing.T) {
	brain := Brain{ManipulationSystem: ManipulationSystem{}}
	m := randomChangeManipulatorValue(&brain)
	if !reflect.DeepEqual(m, mChangeManipulatorValue{}) {
		t.Errorf("Expected zero value mChangeManipulatorValue, got: %v", m)
	}
	newBrain := m.apply(&brain)
	if !reflect.DeepEqual(&brain, newBrain) {
		t.Errorf("Expected mutation not to be applied because of incorrect link")
	}
}
