package main

import "fmt"

func main() {
	FindOnlyOne()
}

// FindOnlyOne 找出只出现一次的数字
func FindOnlyOne() {

	arr := [14]int{1, 2, 3, 2, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	count := make(map[int]int)

	for _, num := range arr {
		count[num]++
	}

	for num, n := range count {
		if n == 1 {
			fmt.Println(num)
		}
	}
}
