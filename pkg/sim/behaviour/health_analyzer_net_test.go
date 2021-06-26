package behaviour

import (
	"reflect"
	"testing"
)

func TestCorrection(t *testing.T) {
	testNet := HealthAnalyzerNet{
		// MinHealth, MaxHealth, OuterAnalyzerLink, Correction
		0, 50, 0, 50,
		30, 70, 0, 200,
		40, 150, 1, 100,
		200, 255, 2, 1,
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
