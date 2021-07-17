package behaviour

import (
	"reflect"
	"testing"
)

func TestCollectSignal(t *testing.T) {
	receptor := OuterReceptor{
		visionRange: 2,
		cells:       [][]uint8{{1, 2, 2}, {0, 0, 0}},
	}
	sig := []uint8{0, 2, 3, 1, 1, 0}
	res := receptor.CollectSignal(sig)
	expectedRes := CollectedOuterSignal{
		0: {0: 1, 1: 2},
		1: {0: 1},
		2: {2: 1, 3: 1},
	}
	if !reflect.DeepEqual(res, expectedRes) {
		t.Errorf("expect result to eq %v, got %v", expectedRes, res)
	}
}
