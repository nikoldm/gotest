package main

import (
	"fmt"
	"sync"
)

// Goroutine:一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
func printOdd(done func()) {
	defer done()
	for i := 1; i <= 100; i++ {
		if i%2 == 1 {
			fmt.Printf("%d ", i)
		}
	}
	fmt.Println(" ")
}

func event(done func()) {
	defer done()
	for i := 1; i <= 100; i++ {
		if i%2 == 0 {
			fmt.Printf("(%d) ", i)
		}
	}
	fmt.Println(" == ")
}

// 闭包
func addAdapt() func(int) int {
	n := 10
	return func(x int) int {
		n = n + x
		return n
	}
}

func main() {

	//1、通过 计数器 等待所有协程完成
	var wg sync.WaitGroup
	wg.Add(2)
	fmt.Println("3、打印从0到10的奇偶数==========")
	go printOdd(wg.Done) //也可以传递指针，使用同一个地址
	go event(wg.Done)
	wg.Wait()

	// 2、通过通道通知主协程协程已完成
	workDone := make(chan bool)
	workDone1 := make(chan bool)
	go func() {
		for i := 0; i < 100; i++ {
			if i%2 == 0 {
				fmt.Printf("=%d= ", i)
			}
		}
		workDone <- true
	}()
	go func() {
		for i := 0; i < 100; i++ {
			if i%2 == 1 {
				fmt.Printf("*%d* ", i)
			}
		}
		workDone1 <- true
	}()
	<-workDone
	<-workDone1
	fmt.Println("协程已完成。")

}
