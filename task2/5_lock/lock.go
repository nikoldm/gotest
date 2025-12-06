package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

//题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
//启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。

//题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
//启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。

func main() {

	//1、使用 sync.Mutex 来保护一个共享的计数器
	var sum int
	var mu sync.Mutex
	var sg sync.WaitGroup

	for i := 0; i < 10; i++ {
		sg.Add(1)
		go func(s int) {
			defer sg.Done()
			for j := 0; j < 1000; j++ {
				// 这里加锁
				mu.Lock()
				sum++
				mu.Unlock()
			}
		}(sum)
	}
	sg.Wait()
	fmt.Println("sum加锁:", sum)

	// 2、原子操作（ sync/atomic 包）实现一个无锁的计数器
	var counter int64
	var sw sync.WaitGroup
	for i := 0; i < 10; i++ {
		sw.Add(1)
		go increateNum(&counter, &sw)
	}
	sw.Wait()
	fmt.Println("counter原子操作:", counter)
}

func increateNum(counter *int64, s *sync.WaitGroup) {
	defer s.Done()
	for i := 0; i < 1000; i++ {
		atomic.AddInt64(counter, 1)
		time.Sleep(time.Millisecond)
	}
}
