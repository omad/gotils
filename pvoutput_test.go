package main

import (
	"fmt"
	"reflect"
	"testing"
)

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func TestIntMinBasic(t *testing.T) {
	ans := IntMin(2, -2)
	if ans != -2 {
		t.Errorf("IntMin(2, -2) = %d; want -2", ans)
	}
}
func TestIntMinTableDriven(t *testing.T) {
	var tests = []struct {
		a, b int
		want int
	}{
		{0, 1, 0},
		{1, 0, 0},
		{2, -2, -2},
		{0, -1, -1},
		{-1, 0, -1},
	}
	for _, tt := range tests {

		testname := fmt.Sprintf("%d,%d", tt.a, tt.b)
		t.Run(testname, func(t *testing.T) {
			ans := IntMin(tt.a, tt.b)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

func TestParseVals(t *testing.T) {
	s := `20200104    25.973 kWh
	20200103    25.386 kWh
	20200102    25.332 kWh
	20200101    26.817 kWh
	20191231    26.513 kWh`

	output := parseAuroraOutputs(s)
	expected := map[string]int{
		"20191231": 26513,
		"20200101": 26816,
		"20200102": 25332,
		"20200103": 25385,
		"20200104": 25972,
	}
	if output == nil {
		t.Error()
	}
	if !reflect.DeepEqual(output, expected) {
		t.Error()

	}
}

// func ReverseRunes(in string) string {
// 	return sort.Sort(sort.Reverse(sort.StringSlice(in)))
// }
// func TestReverseRunes(t *testing.T) {
// 	cases := []struct {
// 		in, want string
// 	}{
// 		{"Hello, world", "dlrow ,olleH"},
// 		{"Hello, 世界", "界世 ,olleH"},
// 		{"", ""},
// 	}
// 	for _, c := range cases {
// 		got := ReverseRunes(c.in)
// 		if got != c.want {
// 			t.Errorf("ReverseRunes(%q) == %q, want %q", c.in, got, c.want)
// 		}
// 	}
// }
