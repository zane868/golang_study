package main

import (
	"reflect"
	"testing"
)

func TestFindOnlyOne(t *testing.T) {
	input := []int{1, 2, 3, 2, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	got := FindOnlyOne(input)
	want := []int{4, 5, 6, 7, 8, 9, 10}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FindOnlyOne(%v) = %v, want %v", input, got, want)
	}
}

func TestIsPalindromicNumber(t *testing.T) {
	tests := []struct {
		name string
		num  int
		want bool
	}{
		{name: "palindrome", num: 11211, want: true},
		{name: "not palindrome", num: 12345, want: false},
		{name: "single digit", num: 7, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPalindromicNumber(tt.num); got != tt.want {
				t.Fatalf("IsPalindromicNumber(%d) = %v, want %v", tt.num, got, tt.want)
			}
		})
	}
}

func TestIsValidParentheses(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{name: "empty", in: "", want: true},
		{name: "single pair", in: "()", want: true},
		{name: "multiple pairs", in: "()[]{}", want: true},
		{name: "nested", in: "([{}])", want: true},
		{name: "mismatch", in: "([)]", want: false},
		{name: "extra close", in: "]", want: false},
		{name: "missing close", in: "([{}", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidParentheses(tt.in); got != tt.want {
				t.Fatalf("IsValidParentheses(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestLongestCommonPrefix(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want string
	}{
		{name: "common prefix", in: []string{"flower", "flow", "flight"}, want: "fl"},
		{name: "no common prefix", in: []string{"dog", "racecar", "car"}, want: ""},
		{name: "single string", in: []string{"hello"}, want: "hello"},
		{name: "empty array", in: []string{}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LongestCommonPrefix(tt.in); got != tt.want {
				t.Fatalf("LongestCommonPrefix(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestPlusOne(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{name: "normal increment", in: []int{1, 2, 3}, want: []int{1, 2, 4}},
		{name: "carry", in: []int{9}, want: []int{1, 0}},
		{name: "multi carry", in: []int{9, 9, 9}, want: []int{1, 0, 0, 0}},
		{name: "no carry", in: []int{1, 9, 9}, want: []int{2, 0, 0}},
		{name: "empty", in: []int{}, want: []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PlusOne(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("PlusOne(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{name: "empty", in: []int{}, want: 0},
		{name: "single number", in: []int{7}, want: 1},
		{name: "no duplicates", in: []int{1, 2, 3}, want: 3},
		{name: "with duplicates", in: []int{1, 1, 2}, want: 2},
		{name: "all same", in: []int{1, 1, 1, 1}, want: 1},
		{name: "multiple duplicates", in: []int{0, 0, 1, 1, 2, 2, 3}, want: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicates(tt.in); got != tt.want {
				t.Fatalf("RemoveDuplicates(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		want [][]int
	}{
		{name: "overlap", in: [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}, want: [][]int{{1, 6}, {8, 10}, {15, 18}}},
		{name: "adjacent", in: [][]int{{1, 2}, {2, 3}, {4, 5}}, want: [][]int{{1, 3}, {4, 5}}},
		{name: "disjoint", in: [][]int{{1, 2}, {4, 5}, {7, 8}}, want: [][]int{{1, 2}, {4, 5}, {7, 8}}},
		{name: "empty", in: [][]int{}, want: [][]int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Merge(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Merge(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestTwoSum2(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		target int
		want   []int
	}{
		{name: "example", input: []int{2, 7, 11, 15}, target: 9, want: []int{0, 1}},
		{name: "duplicate values", input: []int{3, 3}, target: 6, want: []int{0, 1}},
		{name: "no solution", input: []int{1, 2, 3}, target: 10, want: []int{}},
		{name: "single element", input: []int{5}, target: 10, want: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TwoSum2(tt.input, tt.target); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("TwoSum2(%v, %d) = %v, want %v", tt.input, tt.target, got, tt.want)
			}
		})
	}
}

func TestTwoSum(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		target int
		want   []int
	}{
		{name: "example", input: []int{2, 7, 11, 15}, target: 9, want: []int{0, 1}},
		{name: "duplicate values", input: []int{3, 3}, target: 6, want: []int{0, 1}},
		{name: "no solution", input: []int{1, 2, 3}, target: 10, want: []int{}},
		{name: "single element", input: []int{5}, target: 10, want: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TwoSum(tt.input, tt.target); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("TwoSum(%v, %d) = %v, want %v", tt.input, tt.target, got, tt.want)
			}
		})
	}
}
