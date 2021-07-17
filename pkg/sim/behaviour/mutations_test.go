package behaviour

import (
	"reflect"
	"testing"
)

func TestIncreaseVisionRange(t *testing.T) {
	testCases := []struct {
		origReceptor OuterReceptor
		newReceptor  OuterReceptor
	}{
		{
			origReceptor: OuterReceptor{},
			newReceptor:  OuterReceptor{1, [][]uint8{{0, 0}}},
		},
		{
			origReceptor: OuterReceptor{1, [][]uint8{{0, 0}}},
			newReceptor:  OuterReceptor{2, [][]uint8{{0, 0, 0}, {0, 0, 0}}},
		},
		{
			origReceptor: OuterReceptor{
				5,
				[][]uint8{
					{0, 1, 1, 1, 2, 2},
					{3, 3, 1, 1, 2, 2},
					{3, 3, 3, 1, 5, 5},
					{4, 4, 3, 5, 5, 5},
					{4, 4, 4, 5, 5, 5},
				},
			},
			newReceptor: OuterReceptor{
				6,
				[][]uint8{
					{0, 0, 1, 1, 1, 2, 2},
					{0, 0, 1, 1, 1, 2, 2},
					{3, 3, 3, 1, 1, 2, 2},
					{3, 3, 3, 3, 1, 5, 5},
					{4, 4, 4, 3, 5, 5, 5},
					{4, 4, 4, 4, 5, 5, 5},
				},
			},
		},
		{
			// When vision range has maximal value it is not increased
			origReceptor: NewOuterReceptor(maxVisionRange),
			newReceptor:  NewOuterReceptor(maxVisionRange),
		},
	}

	for i, c := range testCases {
		brain := testBrain()
		brain.OuterReceptor = c.origReceptor
		newBrain := mIncreaseVisionRange{}.Apply(brain)

		if !reflect.DeepEqual(newBrain.OuterReceptor, c.newReceptor) {
			t.Errorf("[%d] Expected new new outer receptor to equal %v, got %v", i, c.newReceptor, newBrain.OuterReceptor)
		}
	}

	brain := testBrain()
	brain.OuterReceptor = OuterReceptor{1, [][]uint8{{0, 1}}}
	newBrain := mIncreaseVisionRange{}.Apply(brain)

	if !reflect.DeepEqual(brain.OuterReceptor, OuterReceptor{1, [][]uint8{{0, 1}}}) {
		t.Errorf("Expected original brain to not change; but OuterReceptor changed to %v", brain.OuterReceptor)
	}

	newBrain.OuterReceptor = brain.OuterReceptor
	if !reflect.DeepEqual(brain, newBrain) {
		t.Errorf("Expected the rest of new brain to be a copy of original brain\n\nOriginal brain:\n%v\n\nNew brain:\n%v", brain, newBrain)
	}
}

func TestDencreaseVisionRange(t *testing.T) {
	testCases := []struct {
		origReceptor OuterReceptor
		newReceptor  OuterReceptor
	}{
		{
			origReceptor: OuterReceptor{},
			newReceptor:  OuterReceptor{},
		},
		{
			origReceptor: OuterReceptor{1, [][]uint8{{0, 0}}},
			newReceptor:  OuterReceptor{},
		},
		{
			origReceptor: OuterReceptor{
				4,
				[][]uint8{
					{3, 1, 1, 2, 2},
					{3, 3, 1, 5, 5},
					{4, 3, 5, 5, 5},
					{4, 4, 5, 5, 5},
				},
			},
			newReceptor: OuterReceptor{
				3,
				[][]uint8{
					{3, 1, 5, 5},
					{3, 5, 5, 5},
					{4, 5, 5, 5},
				},
			},
		},
	}

	for i, c := range testCases {
		brain := testBrain()
		brain.OuterReceptor = c.origReceptor
		newBrain := mDecreaseVisionRange{}.Apply(brain)

		if !reflect.DeepEqual(newBrain.OuterReceptor, c.newReceptor) {
			t.Errorf("[%d] Expected new new outer receptor to equal %#v, got %#v", i, c.newReceptor, newBrain.OuterReceptor)
		}
	}

	brain := testBrain()
	brain.OuterReceptor = OuterReceptor{1, [][]uint8{{0, 1}}}
	newBrain := mIncreaseVisionRange{}.Apply(brain)

	if !reflect.DeepEqual(brain.OuterReceptor, OuterReceptor{1, [][]uint8{{0, 1}}}) {
		t.Errorf("Expected original brain to not change; but OuterReceptor changed to %v", brain.OuterReceptor)
	}

	newBrain.OuterReceptor = brain.OuterReceptor
	if !reflect.DeepEqual(brain, newBrain) {
		t.Errorf("Expected the rest of new brain to be a copy of original brain\n\nOriginal brain:\n%v\n\nNew brain:\n%v", brain, newBrain)
	}
}
