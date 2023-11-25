package slice

import (
	"reflect"
	"testing"
)

func TestSliceDelete(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []int
		index    int
		expected []int
	}{
		{"DeleteMiddle", []int{1, 2, 3, 4, 5}, 2, []int{1, 2, 4, 5}},
		{"DeleteFrist", []int{1, 2, 3, 4, 5}, 0, []int{2, 3, 4, 5}},
		{"DeleteLast", []int{1, 2, 3, 4, 5}, 4, []int{1, 2, 3, 4}},
		{"DeleteInvalid", []int{1, 2, 3, 4, 5}, 5, []int{1, 2, 3, 4, 5}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SliceDelete(tc.index, tc.slice)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Test case %s failed, Expected %v, but got %v", tc.name, tc.expected, result)
			}
		})
	}
}
