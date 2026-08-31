package main

import (
	"fmt"
	"sort"
	"strconv"
)

func main() {
	fmt.Println("test run main_test.go check result")
}

// FindOnlyOne 找出只出现一次的数字
func FindOnlyOne(arr []int) []int {
	count := make(map[int]int)

	for _, num := range arr {
		count[num]++
	}

	result := make([]int, 0)
	for num, n := range count {
		if n == 1 {
			result = append(result, num)
		}
	}

	sort.Ints(result)
	return result
}

// 判断是否回文数
func PalindromicNumber(num int) bool {
	numStr := strconv.Itoa(num)
	reverseResult := Reverse(numStr)

	return reverseResult == numStr
}

// 有效的括号
func ValidParentheses(str string) bool {
	r := []rune(str)
	stack := &Stack{}
	pairs := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, ch := range r {
		switch ch {
		case '(', '[', '{':
			stack.Push(int(ch))
		case ')', ']', '}':
			if stack.IsEmpty() {
				return false
			}
			top := rune(stack.Pop())
			if pairs[ch] != top {
				return false
			}
		}
	}
	return stack.IsEmpty()
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
