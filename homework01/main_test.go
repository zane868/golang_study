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
