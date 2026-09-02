package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	s := Rectangle{
		length: 5,
		width:  3,
	}
	fmt.Printf("Area: %d\n", s.Area())
	fmt.Printf("Perimeter: %d\n", s.Perimeter())
	fmt.Println("run main_test.go check result")
}

// 任务调度器
func ExecuteScheduler() {
	scheduler := NewScheduler(3, 10)

	// 提交任务
	for i := 0; i < 1; i++ {
		taskID := i
		scheduler.Submit(func(ctx context.Context) (int, error) {
			fmt.Printf("%s Task %d start\n", nowTime(), taskID)
			select {
			case <-time.After(2 * time.Second):
				fmt.Printf("%s Task %d done\n", nowTime(), taskID)
				return taskID, nil

			case <-ctx.Done():
				return taskID, ctx.Err()
			}

		})
	}

	//等待一段时间后关闭调度器
	scheduler.Shutdown()
}

// 获取当前时间
func nowTime() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
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
