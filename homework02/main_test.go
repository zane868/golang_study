package main
package main

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{name: "positive", in: 5, want: 15},
		{name: "zero", in: 0, want: 10},
		{name: "negative", in: -3, want: 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in
			add(&got)
			if got != tt.want {
				t.Fatalf("add(%d) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestDouble(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{name: "multiple values", in: []int{2, 7, 10}, want: []int{4, 14, 20}},
		{name: "zero and negative", in: []int{0, -3, 1}, want: []int{0, -6, 2}},
		{name: "empty", in: []int{}, want: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := make([]*int, len(tt.in))
			for i := range tt.in {
				v := tt.in[i]
				values[i] = &v
			}

			double(values)

			for i := range values {
				if *values[i] != tt.want[i] {
					t.Fatalf("double(%v) = %v, want %v", tt.in, extractValues(values), tt.want)
				}
			}
		})
	}
}

func extractValues(values []*int) []int {
	result := make([]int, len(values))
	for i, v := range values {
		result[i] = *v
	}
	return result
}
