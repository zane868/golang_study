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

func TestPalindromicNumber(t *testing.T) {
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
			if got := PalindromicNumber(tt.num); got != tt.want {
				t.Fatalf("PalindromicNumber(%d) = %v, want %v", tt.num, got, tt.want)
			}
		})
	}
}

func TestValidParentheses(t *testing.T) {
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
			if got := ValidParentheses(tt.in); got != tt.want {
				t.Fatalf("ValidParentheses(%q) = %v, want %v", tt.in, got, tt.want)
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
			if got := plusOne(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("plusOne(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
