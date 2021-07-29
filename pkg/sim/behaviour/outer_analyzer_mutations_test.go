package behaviour

import (
	"experiments/pkg/test_helpers"
	"reflect"
	"testing"
)

func TestChangeOuterAnalyzerLinkPower(t *testing.T) {
	origNet := func() OuterAnalyzerNet {
		return OuterAnalyzerNet{
			&OuterAnalyzerLink{0, 0, 0, 20},
			&OuterAnalyzerLink{1, 0, 1, -5},
		}
	}
	outerAnalyzerCount := 3
	mSys := ManipulationSystem{&Manipulator{}, &Manipulator{}}

	testCases := []struct {
		desc              string
		healthAnalyzerNet HealthAnalyzerNet
		mutation          mChangeOuterAnalyzerLinkPower
		newNet            OuterAnalyzerNet
	}{
		{
			desc:     "increase power",
			mutation: mChangeOuterAnalyzerLinkPower{1, 0, 1, 30},
			newNet: OuterAnalyzerNet{
				&OuterAnalyzerLink{0, 0, 0, 20},
				&OuterAnalyzerLink{1, 0, 1, 25},
			},
		},
		{
			desc:     "decrease power",
			mutation: mChangeOuterAnalyzerLinkPower{0, 0, 0, -50},
			newNet: OuterAnalyzerNet{
				&OuterAnalyzerLink{0, 0, 0, -30},
				&OuterAnalyzerLink{1, 0, 1, -5},
			},
		},
		{
			desc:     "positive overflow",
			mutation: mChangeOuterAnalyzerLinkPower{0, 0, 0, 120},
			newNet: OuterAnalyzerNet{
				&OuterAnalyzerLink{0, 0, 0, 127},
				&OuterAnalyzerLink{1, 0, 1, -5},
			},
		},
		{
			desc:     "negative overflow",
			mutation: mChangeOuterAnalyzerLinkPower{1, 0, 1, -125},
			newNet: OuterAnalyzerNet{
				&OuterAnalyzerLink{0, 0, 0, 20},
				&OuterAnalyzerLink{1, 0, 1, -128},
			},
		},
		{
			desc:     "when power gets zero",
			mutation: mChangeOuterAnalyzerLinkPower{0, 0, 0, -20},
			newNet: OuterAnalyzerNet{
				&OuterAnalyzerLink{1, 0, 1, -5},
			},
		},
		{
			desc: "when power gets zero but there are references from health analyzer",
			healthAnalyzerNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{0, 10, 0, 50},
			},
			mutation: mChangeOuterAnalyzerLinkPower{0, 0, 0, -20},
			newNet: OuterAnalyzerNet{
				&OuterAnalyzerLink{0, 0, 0, 20},
				&OuterAnalyzerLink{1, 0, 1, -5},
			},
		},
		{
			desc:     "when link does not exist (by manipulator)",
			mutation: mChangeOuterAnalyzerLinkPower{0, 0, 1, 50},
			newNet: OuterAnalyzerNet{
				&OuterAnalyzerLink{0, 0, 0, 20},
				&OuterAnalyzerLink{1, 0, 1, -5},
				&OuterAnalyzerLink{0, 0, 1, 50},
			},
		},
		{
			desc:     "when link does not exist (by signal)",
			mutation: mChangeOuterAnalyzerLinkPower{1, 3, 1, -33},
			newNet: OuterAnalyzerNet{
				&OuterAnalyzerLink{0, 0, 0, 20},
				&OuterAnalyzerLink{1, 0, 1, -5},
				&OuterAnalyzerLink{1, 3, 1, -33},
			},
		},
		{
			desc:     "when link does not exist (by analyzer)",
			mutation: mChangeOuterAnalyzerLinkPower{2, 0, 1, 42},
			newNet: OuterAnalyzerNet{
				&OuterAnalyzerLink{0, 0, 0, 20},
				&OuterAnalyzerLink{1, 0, 1, -5},
				&OuterAnalyzerLink{2, 0, 1, 42},
			},
		},
		{
			desc:     "when link does not exist but wrong analyzer",
			mutation: mChangeOuterAnalyzerLinkPower{4, 0, 0, 11},
			newNet: OuterAnalyzerNet{
				&OuterAnalyzerLink{0, 0, 0, 20},
				&OuterAnalyzerLink{1, 0, 1, -5},
			},
		},
		{
			desc:     "when link does not exist but wrong manipulator",
			mutation: mChangeOuterAnalyzerLinkPower{0, 0, 2, 100},
			newNet: OuterAnalyzerNet{
				&OuterAnalyzerLink{0, 0, 0, 20},
				&OuterAnalyzerLink{1, 0, 1, -5},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.desc, func(t *testing.T) {
			brain := testBrain()
			brain.OuterAnalyzerNet = origNet()
			brain.OuterAnalyzersCount = outerAnalyzerCount
			brain.ManipulationSystem = mSys
			brain.HealthAnalyzerNet = c.healthAnalyzerNet

			newBrain := c.mutation.apply(brain)

			if !reflect.DeepEqual(newBrain.OuterAnalyzerNet, c.newNet) {
				t.Errorf("[Expected new OuterAnalyzerNet to equal %s, got %s",
					test_helpers.Inspect(c.newNet), test_helpers.Inspect(newBrain.OuterAnalyzerNet))
			}
			if !reflect.DeepEqual(brain.OuterAnalyzerNet, origNet()) {
				t.Errorf("Expected original brain to not change; but OuterAnalyzerNet changed to %v",
					test_helpers.Inspect(brain.OuterAnalyzerNet))
			}
			newBrain.OuterAnalyzerNet = brain.OuterAnalyzerNet
			if !reflect.DeepEqual(brain, newBrain) {
				t.Errorf("Expected the rest of new brain to be a copy of original brain\n\nOriginal brain:\n%v\n\nNew brain:\n%v", brain, newBrain)
			}
		})
	}
}

func TestRandomChangeOuterAnalyzerLinkPower(t *testing.T) {
	brain := Brain{ManipulationSystem: ManipulationSystem{}}
	m := randomChangeOuterAnalyzerLinkPower(&brain)
	if !reflect.DeepEqual(m, mChangeOuterAnalyzerLinkPower{}) {
		t.Errorf("Expected zero value mChangeOuterAnalyzerLinkPower, got: %v", m)
	}
	newBrain := m.apply(&brain)
	if !reflect.DeepEqual(&brain, newBrain) {
		t.Errorf("Expected mutation not to be applied because of incorrect link")
	}
}
