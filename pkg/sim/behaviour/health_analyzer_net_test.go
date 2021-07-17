package behaviour

import (
	"reflect"
	"testing"
)

func TestCorrection(t *testing.T) {
	testNet := HealthAnalyzerNet{
		{MinHealth: 0, MaxHealth: 50, OuterAnalyzerLink: 0, Correction: 50},
		{MinHealth: 30, MaxHealth: 70, OuterAnalyzerLink: 0, Correction: 200},
		{MinHealth: 40, MaxHealth: 150, OuterAnalyzerLink: 1, Correction: 100},
		{MinHealth: 200, MaxHealth: 255, OuterAnalyzerLink: 2, Correction: 1},
	}

	cases := []struct {
		healthIndicator int
		Result          OuterAnalyzerNetCorrection
	}{
		{1, OuterAnalyzerNetCorrection{0: 3.125}},
		{40, OuterAnalyzerNetCorrection{0: 39.0625, 1: 6.25}},
		{199, OuterAnalyzerNetCorrection{}},
		{1000, OuterAnalyzerNetCorrection{2: 0.0625}},
	}

	for i, c := range cases {
		res := testNet.Correction(c.healthIndicator)
		if !reflect.DeepEqual(res, c.Result) {
			t.Errorf("[%d] Expected correction to eq %v, got %v", i, c.Result, res)
		}
	}
}
