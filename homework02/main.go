package main

import "fmt"

func main() {
	fmt.Println("run main_test.go check result")
}

func add(x *int) {
	*x = *x + 10
}

func double(x []*int) {
	for i := 0; i < len(x); i++ {
		*x[i] = *x[i] * 2
	}
}
