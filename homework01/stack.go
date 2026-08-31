package main

import "fmt"

type Stack struct {
	data []int
}

func (s *Stack) IsEmpty() bool {
	return len(s.data) == 0
}

// 入栈
func (s *Stack) Push(value int) {
	s.data = append(s.data, value)
}

// 出栈
func (s *Stack) Pop() int {
	if s.IsEmpty() {
		panic("stack is empty")
	}

	index := len(s.data) - 1
	value := s.data[index]
	s.data = s.data[:index]
	return value
}

// 查看栈顶
func (s *Stack) Peek() int {
	if s.IsEmpty() {
		panic("stack is empty")
	}

	return s.data[len(s.data)-1]
}

func (s *Stack) Size() int {
	return len(s.data)
}

func (s *Stack) ToString() string {
	return fmt.Sprint(s.data)
}

func (s *Stack) String() string {
	return s.ToString()
}
