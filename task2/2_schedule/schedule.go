package main

import (
	"fmt"
	"sync"
	"time"
)

// Task 设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
// 定义任务：
type Task func()

// TaskSchedule 定义调度器
func TaskSchedule(tasks []Task) {

	var sg sync.WaitGroup
	// 为每个任务增加一个计数器
	sg.Add(len(tasks))

	//执行任务
	for _, task := range tasks {
		go func(t Task) {
			start := time.Now()
			defer sg.Done()
			t()
			gap := time.Since(start)
			fmt.Println("任务耗时：", gap)
		}(task)
	}
	// 等待所有任务完成
	sg.Wait()

}

func main() {
	tasks := []Task{
		func() {
			fmt.Println("Task 1 started,时间：", time.Now())
			time.Sleep(1 * time.Second)
			fmt.Println("Task 1 finished")
		},
		func() {
			fmt.Println("Task 2 started,时间：", time.Now())
			time.Sleep(2 * time.Second)
			fmt.Println("Task 2 finished")
		},
		func() {
			fmt.Println("Task 3 started,时间：", time.Now())
			time.Sleep(3 * time.Second)
			fmt.Println("Task 3 finished")
		},
	}
	TaskSchedule(tasks)
	fmt.Println("所有任务完成。", time.Now())
}
