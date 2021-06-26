package behaviour

import (
	"reflect"
	"testing"
)

// func inspect(val interface{}) string {
// 	switch v := reflect.ValueOf(val); v.Kind() {
// 	case reflect.Array, reflect.Slice:
// 		s := "["
// 		for i := 0; i < v.Len(); i++ {
// 			if i > 0 {
// 				s += ", "
// 			}
// 			s += fmt.Sprintf("%v", v.Index(i))
// 		}
// 		s += "]"
// 		return s
// 	default:
// 		return fmt.Sprintf("%v", val)
// 	}
// }

func TestSimpleMutation(t *testing.T) {
	brain := &Brain{
		Structure: BrainStructure{
			VisionRange:            1,
			OuterAnalyzersCount:    2,
			HealthAnalyzerNetSize:  1,
			OuterAnalyzerNetSize:   2,
			ManipulationSystemSize: 1,
		},
		Content: []uint8{
			// OuterReceptor
			0, 0,
			// HealthAnalyzerNet
			0, 10, 1, 111,
			// OuterAnalyzerNet
			1, 0, 0, 255,
			0, 1, 0, 28,
			// ManipulationSystem
			1, 0, 0, 0, 0, 0, 0, 0, 0,
		},
	}

	cases := []struct {
		mutation   simpleMutation
		newContent []uint8
	}{
		{
			simpleMutation{0: 1, 4: 100, 9: 33, 21: 88},
			[]uint8{1, 0, 0, 10, 0, 111, 1, 0, 0, 33, 0, 1, 0, 28, 1, 0, 0, 0, 0, 0, 0, 88, 0},
		},
		{
			simpleMutation{1: 20, 99: 100},
			[]uint8{0, 0, 0, 10, 1, 111, 1, 0, 0, 255, 0, 1, 0, 28, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for i, c := range cases {
		newBrain := c.mutation.Apply(brain)
		if !reflect.DeepEqual(newBrain.Structure, brain.Structure) {
			t.Errorf("[%d] Expected new brain Structure to be the same, but changed to %v", i, newBrain.Structure)
		}
		if !reflect.DeepEqual(newBrain.Content, c.newContent) {
			t.Errorf("[%d] Expected new brain Content to eq %v, got %v", i, c.newContent, newBrain.Content)
		}
	}
}
