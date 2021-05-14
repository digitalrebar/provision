package cli

import (
	"reflect"
	"testing"
)

func TestColorProcessor(t *testing.T) {
	orig := colorPatterns

	colorPatterns = nil
	colorString = "0=1;1=2,3;2=3;3=5;4=5,6,7;5=1;6=5;7=9;8=13;9=33"
	processColorPatterns()

	expected := [][]int{
		[]int{1},
		[]int{2, 3},
		[]int{3},
		[]int{5},
		[]int{5, 6, 7},
		[]int{1},
		[]int{5},
		[]int{9},
		[]int{13},
	}
	if !reflect.DeepEqual(colorPatterns, expected) {
		t.Errorf("Failed to set colorPatterns correctly: expected: %v actual: %v", expected, colorPatterns)
	}
	colorPatterns = orig
}
