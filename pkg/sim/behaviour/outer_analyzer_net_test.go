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
		Net    OuterAnalyzerNet
		Signal map[uint8][]uint8
		Result map[uint8]int16
	}{
		{
			Net: testNet,
			Signal: map[uint8][]uint8{
				0: []uint8{2, 1},
				1: []uint8{1, 1},
			},
			Result: map[uint8]int16{5: 15, 8: 250},
		},
		{
			Net:    testNet,
			Signal: map[uint8][]uint8{2: []uint8{2, 1}},
			Result: map[uint8]int16{},
		},
	}

	for i, c := range cases {
		res := c.Net.Activation(c.Signal)
		if !reflect.DeepEqual(res, c.Result) {
			t.Errorf("[%d] expect result to eq %v, got %v", i, c.Result, res)
		}
	}
}
