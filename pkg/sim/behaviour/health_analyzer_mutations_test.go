package behaviour

import (
	"experiments/pkg/test_helpers"
	"reflect"
	"testing"
)

func TestAddHealthAnalyzerLink(t *testing.T) {
	brain := testBrain()
	brain.OuterAnalyzerNet = randomOuterAnalyzerNet(3, 2)
	link := randomHealthAnalyzerLink(3)
	brain.HealthAnalyzerNet = HealthAnalyzerNet{link}

	newLink := HealthAnalyzerLink{OuterAnalyzerLink: 1}
	newBrain := mAddHealthAnalyzerLink{newLink}.apply(brain)

	expectedNet := HealthAnalyzerNet{link, &newLink}
	if !reflect.DeepEqual(newBrain.HealthAnalyzerNet, expectedNet) {
		t.Errorf("Expected new HealthAnalyzerNet to equal %v, got %v", newBrain.HealthAnalyzerNet, expectedNet)
	}

	if !reflect.DeepEqual(brain.HealthAnalyzerNet, HealthAnalyzerNet{link}) {
		t.Errorf("Expected original brain to not change; but HealthAnalyzerNet changed to %v", brain.HealthAnalyzerNet)
	}

	newBrain.HealthAnalyzerNet = brain.HealthAnalyzerNet
	if !brainsEqual(brain, newBrain) {
		t.Errorf("Expected the rest of new brain to be a copy of original brain\n\nOriginal brain:\n%v\n\nNew brain:\n%v", brain, newBrain)
	}

	newLink = HealthAnalyzerLink{OuterAnalyzerLink: 3}
	newBrain = mAddHealthAnalyzerLink{newLink}.apply(brain)
	if !reflect.DeepEqual(newBrain.HealthAnalyzerNet, HealthAnalyzerNet{link}) {
		t.Errorf("Expected mutation not to be applied because of incorrect link")
	}
}

func TestRandomAddHealthAnalyzerLink(t *testing.T) {
	brain := Brain{OuterAnalyzerNet: OuterAnalyzerNet{}}
	m := randomAddHealthAnalyzerLink(&brain)
	if !reflect.DeepEqual(m, mAddHealthAnalyzerLink{}) {
		t.Errorf("Expected zero value mChangeOuterReceptor, got: %v", m)
	}
	newBrain := m.apply(&brain)
	if !brainsEqual(&brain, newBrain) {
		t.Errorf("Expected mutation not to be applied because of incorrect link")
	}
}

func TestChangeHealthAnalyzerCorrection(t *testing.T) {
	origNet := func() HealthAnalyzerNet {
		return HealthAnalyzerNet{
			&HealthAnalyzerLink{Correction: 5},
			&HealthAnalyzerLink{Correction: 240},
		}
	}
	testCases := []struct {
		mutation mChangeHealthAnalyzerCorrection
		newNet   HealthAnalyzerNet
	}{
		{
			mutation: mChangeHealthAnalyzerCorrection{0, 10},
			newNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{Correction: 15},
				&HealthAnalyzerLink{Correction: 240},
			},
		},
		{
			mutation: mChangeHealthAnalyzerCorrection{1, -50},
			newNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{Correction: 5},
				&HealthAnalyzerLink{Correction: 190},
			},
		},
		{
			mutation: mChangeHealthAnalyzerCorrection{1, 50},
			newNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{Correction: 5},
				&HealthAnalyzerLink{Correction: 255},
			},
		},
		{
			mutation: mChangeHealthAnalyzerCorrection{2, -50},
			newNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{Correction: 5},
				&HealthAnalyzerLink{Correction: 240},
			},
		},
		{
			mutation: mChangeHealthAnalyzerCorrection{0, -10},
			newNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{Correction: 0},
				&HealthAnalyzerLink{Correction: 240},
			},
		},
		{
			mutation: mChangeHealthAnalyzerCorrection{0, 11},
			newNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{Correction: 240},
			},
		},
	}

	for i, c := range testCases {
		brain := testBrain()
		brain.HealthAnalyzerNet = origNet()
		newBrain := c.mutation.apply(brain)

		if !reflect.DeepEqual(newBrain.HealthAnalyzerNet, c.newNet) {
			t.Errorf("[%d] Expected new HealthAnalyzerNet to equal %s, got %s",
				i, test_helpers.Inspect(c.newNet), test_helpers.Inspect(newBrain.HealthAnalyzerNet))
		}
	}

	brain := testBrain()
	brain.HealthAnalyzerNet = origNet()
	newBrain := mChangeHealthAnalyzerCorrection{0, 10}.apply(brain)

	if !reflect.DeepEqual(brain.HealthAnalyzerNet, origNet()) {
		t.Errorf("Expected original brain to not change; but HealthAnalyzerNet changed to %v",
			test_helpers.Inspect(brain.HealthAnalyzerNet))
	}

	newBrain.HealthAnalyzerNet = brain.HealthAnalyzerNet
	if !brainsEqual(brain, newBrain) {
		t.Errorf("Expected the rest of new brain to be a copy of original brain\n\nOriginal brain:\n%v\n\nNew brain:\n%v", brain, newBrain)
	}
}

func TestRandomChangeHealthAnalyzerCorrection(t *testing.T) {
	brain := Brain{HealthAnalyzerNet: HealthAnalyzerNet{}}
	m := randomChangeHealthAnalyzerCorrection(&brain)
	if !reflect.DeepEqual(m, mChangeHealthAnalyzerCorrection{}) {
		t.Errorf("Expected zero value mChangeHealthAnalyzerCorrection, got: %v", m)
	}
	newBrain := m.apply(&brain)
	if !brainsEqual(&brain, newBrain) {
		t.Errorf("Expected mutation not to be applied because of incorrect link")
	}
}

func TestChangeHealthAnalyzerMinMax(t *testing.T) {
	origNet := func() HealthAnalyzerNet {
		return HealthAnalyzerNet{
			&HealthAnalyzerLink{MinHealth: 5, MaxHealth: 245},
			&HealthAnalyzerLink{MinHealth: 150, MaxHealth: 170},
		}
	}
	testCases := []struct {
		mutation mChangeHealthAnalyzerMinMax
		newNet   HealthAnalyzerNet
	}{
		{
			mutation: mChangeHealthAnalyzerMinMax{0, 10, -30},
			newNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{MinHealth: 15, MaxHealth: 215},
				&HealthAnalyzerLink{MinHealth: 150, MaxHealth: 170},
			},
		},
		{
			mutation: mChangeHealthAnalyzerMinMax{0, -10, 20},
			newNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{MinHealth: 0, MaxHealth: 255},
				&HealthAnalyzerLink{MinHealth: 150, MaxHealth: 170},
			},
		},
		{
			mutation: mChangeHealthAnalyzerMinMax{1, 10, -10},
			newNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{MinHealth: 5, MaxHealth: 245},
			},
		},
		{
			mutation: mChangeHealthAnalyzerMinMax{1, 50, -50},
			newNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{MinHealth: 5, MaxHealth: 245},
			},
		},
		{
			mutation: mChangeHealthAnalyzerMinMax{0, -50, -250},
			newNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{MinHealth: 150, MaxHealth: 170},
			},
		},
		{
			mutation: mChangeHealthAnalyzerMinMax{0, 253, 50},
			newNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{MinHealth: 150, MaxHealth: 170},
			},
		},
		{
			mutation: mChangeHealthAnalyzerMinMax{2, 50, 0},
			newNet: HealthAnalyzerNet{
				&HealthAnalyzerLink{MinHealth: 5, MaxHealth: 245},
				&HealthAnalyzerLink{MinHealth: 150, MaxHealth: 170},
			},
		},
	}

	for i, c := range testCases {
		brain := testBrain()
		brain.HealthAnalyzerNet = origNet()
		newBrain := c.mutation.apply(brain)

		if !reflect.DeepEqual(newBrain.HealthAnalyzerNet, c.newNet) {
			t.Errorf("[%d] Expected new HealthAnalyzerNet to equal %s, got %s",
				i, test_helpers.Inspect(c.newNet), test_helpers.Inspect(newBrain.HealthAnalyzerNet))
		}
	}

	brain := testBrain()
	brain.HealthAnalyzerNet = origNet()
	newBrain := mChangeHealthAnalyzerCorrection{0, 10}.apply(brain)

	if !reflect.DeepEqual(brain.HealthAnalyzerNet, origNet()) {
		t.Errorf("Expected original brain to not change; but HealthAnalyzerNet changed to %v",
			test_helpers.Inspect(brain.HealthAnalyzerNet))
	}

	newBrain.HealthAnalyzerNet = brain.HealthAnalyzerNet
	if !brainsEqual(brain, newBrain) {
		t.Errorf("Expected the rest of new brain to be a copy of original brain\n\nOriginal brain:\n%v\n\nNew brain:\n%v", brain, newBrain)
	}
}

func TestRandomChangeHealthAnalyzerMinMax(t *testing.T) {
	brain := Brain{HealthAnalyzerNet: HealthAnalyzerNet{}}
	m := randomChangeHealthAnalyzerMinMax(&brain)
	if !reflect.DeepEqual(m, mChangeHealthAnalyzerMinMax{}) {
		t.Errorf("Expected zero value mChangeHealthAnalyzerCorrection, got: %v", m)
	}
	newBrain := m.apply(&brain)
	if !brainsEqual(&brain, newBrain) {
		t.Errorf("Expected mutation not to be applied because of incorrect link")
	}
}
