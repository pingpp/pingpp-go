package utils

import "testing"

func Test_IsRepeated(t *testing.T) {
	array := []string{"sdhj", "ewrqyui", "hdjsabnmxsa", "sdhj", "sdhjsahdjk"}
	var pa *[]int
	testCases := []struct {
		Array        interface{}
		IsRepeated   bool
		RepeatedElem interface{}
	}{
		{
			[]int{1, 23, 4, 565, 4},
			true,
			4,
		}, {
			[]string{"hee", "sdhjs", "qwe", "asd", "qwe"},
			true,
			"qwe",
		}, {
			[]int{1, 2, 3, 4, 5},
			false,
			nil,
		}, {
			&array,
			true,
			"sdhj",
		}, {
			pa,
			false,
			nil,
		},
	}
	for _, tc := range testCases {
		isRepeated, elem := IsRepeated(tc.Array)
		if isRepeated != tc.IsRepeated {
			t.Errorf("should be %v, but result is %v", tc.IsRepeated, isRepeated)
		}
		if elem != tc.RepeatedElem {
			t.Errorf("should be %v, but result is %v", tc.RepeatedElem, elem)
		}
	}
}
