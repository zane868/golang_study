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
	reverseResult := reverse(numStr)

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

// 最长前缀
func LongestCommonPrefix(strs []string) string {

	if len(strs) == 0 {
		return ""
	}
	first := strs[0]
	for i := 1; i < len(strs); i++ {
		for j := 0; j < len(first) && j < len(strs[i]); j++ {
			if first[j] != strs[i][j] {
				first = first[:j]
				break
			}
		}
	}
	return first
}

// 加一
func plusOne(digits []int) []int {

	for i := len(digits) - 1; i >= 0; i-- {

		if digits[i] < 9 {
			digits[i]++
			return digits
		}

		digits[i] = 0
	}

	return append([]int{1}, digits...)
}

// 反转字符串
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
