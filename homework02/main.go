package main

import (
	"fmt"
	"sync"
)

func main() {
	Goroutine01()
	fmt.Println("run main_test.go check result")
}

// 协程简单的使用
func Goroutine01() {
	var wg sync.WaitGroup
	//wg := sync.WaitGroup{}
	wg.Add(2)
	go PrintlNoddNumber(&wg)
	go PrintlEvenNumber(&wg)
	wg.Wait()
}

// 打印1-10奇数
func PrintlNoddNumber(wg *sync.WaitGroup) {
	for i := 1; i <= 10; i++ {
		if i%2 != 0 {
			fmt.Println("NoddNumber", i)
		}
	}
	defer wg.Done()
}

// 打印2-10偶数
func PrintlEvenNumber(wg *sync.WaitGroup) {
	for i := 2; i <= 10; i++ {
		if i%2 == 0 {
			fmt.Println("EvenNumber", i)
		}
	}
	defer wg.Done()
}

func add(x *int) {
	*x = *x + 10
}

func double(x []*int) {
	for i := 0; i < len(x); i++ {
		*x[i] = *x[i] * 2
	}
}
