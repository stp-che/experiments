package behaviour

import (
	"reflect"
	"testing"
)

func TestCollectSignal(t *testing.T) {
	receptor := &OuterReceptor{1, 2, 2, 0, 0, 0}
	sig := []uint8{0, 2, 3, 1, 1, 0}
	res := receptor.CollectSignal(sig)
	expectedRes := CollectedOuterSignal{
		0: {1, 2, 0, 0, 0},
		1: {1, 0, 0, 0, 0},
		2: {0, 0, 1, 1, 0},
	}
	if !reflect.DeepEqual(res, expectedRes) {
		t.Errorf("expect result to eq %v, got %v", expectedRes, res)
	}
}
