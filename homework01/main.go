package main

import (
	"fmt"
	"sort"
	"strconv"
)

func main() {
	fmt.Println("run main_test.go check result")
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

// 删除有序数组中的重复项
func removeDuplicates(nums []int) int {
	fmt.Println("removeDuplicates nums:", nums)

	if len(nums) == 0 {
		return 0
	}
	result := []int{}
	for i := 0; i < len(nums); i++ {
		nextIndex := i + 1
		if nextIndex >= len(nums) {
			result = append(result, nums[i])
			break
		}
		if nums[i] != nums[nextIndex] {
			result = append(result, nums[i])
		}
	}
	fmt.Println("result", result)
	return len(result)
}

// 区间合并
func Merge(intervals [][]int) [][]int {
	fmt.Println("intervals", intervals)

	if len(intervals) == 0 {
		return intervals
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	fmt.Println(intervals)

	//默认把第一个区间加入结果集
	result := [][]int{intervals[0]}
	fmt.Println("result", result)

	for i := 1; i < len(intervals); i++ {
		current := intervals[i]
		last := result[len(result)-1]

		if last[1] >= current[0] {
			result[len(result)-1][1] = max(last[1], current[1])
		} else {
			result = append(result, current)
		}
	}

	return result
}

// 两数之和，暴力循环算差值
func TwoSum(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		targetNum := target - nums[i]
		for j := i + 1; j < len(nums); j++ {
			if nums[j] == targetNum {
				return []int{i, j}
			}
		}
	}
	return []int{}

}

// 两数之和 解法2
func TwoSum2(nums []int, target int) []int {

	seen := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		need := target - nums[i]
		if j, ok := seen[need]; ok {
			return []int{j, i}
		}
		seen[nums[i]] = i
	}
	return []int{}

}
