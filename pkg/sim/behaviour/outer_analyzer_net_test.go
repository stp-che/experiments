package behaviour

import (
	"reflect"
	"testing"
)

func TestActivation(t *testing.T) {
	testNet := OuterAnalyzerNet{
		{Analyzer: 0, Signal: 0, Manipulator: 5, Power: 10},
		{Analyzer: 0, Signal: 0, Manipulator: 8, Power: 100},
		{Analyzer: 0, Signal: 1, Manipulator: 5, Power: -5},
		{Analyzer: 1, Signal: 0, Manipulator: 8, Power: 50},
	}

	cases := []struct {
		Net        OuterAnalyzerNet
		Signal     CollectedOuterSignal
		Correction OuterAnalyzerNetCorrection
		Result     map[uint8]int16
	}{
		{
			Net: testNet,
			Signal: CollectedOuterSignal{
				0: {2, 1},
				1: {1, 1},
			},
			Correction: OuterAnalyzerNetCorrection{},
			Result:     map[uint8]int16{5: 15, 8: 250},
		},
		{
			Net: testNet,
			Signal: CollectedOuterSignal{
				0: {2, 1},
				1: {1, 1},
			},
			Correction: OuterAnalyzerNetCorrection{1: 0.11},
			Result:     map[uint8]int16{5: 15, 8: 72},
		},
		{
			Net:        testNet,
			Signal:     CollectedOuterSignal{2: {2, 1}},
			Correction: OuterAnalyzerNetCorrection{},
			Result:     map[uint8]int16{},
		},
	}

	for i, c := range cases {
		res := c.Net.Activation(c.Signal, c.Correction)
		if !reflect.DeepEqual(res, c.Result) {
			t.Errorf("[%d] expect result to eq %v, got %v", i, c.Result, res)
		}
	}
}
