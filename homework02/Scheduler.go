package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Task 表示一个可以被调度执行的任务。
type Task func(context.Context) (int, error)

// Scheduler 管理任务队列和 worker。
type Scheduler struct {
	Tasks  chan Task          // 任务队列
	ctx    context.Context    // 调度器的上下文
	cancel context.CancelFunc // 用于取消调度器的上下文

	workers sync.WaitGroup // 等待 worker 退出
	tasksWG sync.WaitGroup // 等待已提交任务完成

	closeOnce sync.Once // 确保只关闭一次
}

// NewScheduler 创建任务调度器。
func NewScheduler(workerCount, queueSize int) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	scheduler := &Scheduler{
		Tasks:  make(chan Task, queueSize),
		ctx:    ctx,
		cancel: cancel,
	}

	scheduler.workers.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go scheduler.worker(i)
	}

	return scheduler
}

func (scheduler *Scheduler) worker(workerID int) {
	defer scheduler.workers.Done()
	for {
		select {
		case <-scheduler.ctx.Done():
			return
		case task, ok := <-scheduler.Tasks:
			if !ok {
				fmt.Println("worker", workerID, "exit")
				return
			}
			fmt.Println("worker", workerID, "execute")
			scheduler.execute(workerID, task)
			scheduler.tasksWG.Done()
		}
	}
}

func (scheduler *Scheduler) execute(workerID int, task Task) {
	defer func() {
		if value := recover(); value != nil {
			fmt.Printf("worker %d: 任务发生 panic: %v\n", workerID, value)
		}
	}()

	taskID, err := task(scheduler.ctx)
	if err != nil {
		fmt.Printf("workerId %d: 任务执行失败: %v\n", workerID, err)
		return
	}

	fmt.Printf("workerId: %d: 任务执行完成: taskId:%d\n", workerID, taskID)
}

// Submit 提交任务。
func (scheduler *Scheduler) Submit(task Task) error {
	if task == nil {
		return errors.New("任务不能为空")
	}

	scheduler.tasksWG.Add(1)
	select {
	case <-scheduler.ctx.Done():
		return errors.New("调度器已关闭")
	case scheduler.Tasks <- task:
		return nil
	}
}

// Wait 等待所有已提交任务完成。
func (scheduler *Scheduler) Wait() {
	scheduler.tasksWG.Wait()
}

// Shutdown 等待所有已提交任务完成，然后关闭 workers。
func (scheduler *Scheduler) Shutdown() {
	scheduler.closeOnce.Do(func() {
		scheduler.tasksWG.Wait()
		close(scheduler.Tasks)
		scheduler.workers.Wait()
		scheduler.cancel()
	})
}
