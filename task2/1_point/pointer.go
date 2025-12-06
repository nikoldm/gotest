package main

import (
	"fmt"
)

// 编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10
func pointerAddTen(p *int) {
	*p += 10
}

// 实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
func pointerSlice(slice *[]int) {
	for k, v := range *slice {
		(*slice)[k] = v * 2
	}
}

func main() {

	// 指针
	var p *int
	var a = 3
	p = &a
	pointerAddTen(p)
	fmt.Println("1、指针应用+10：", a)

	ps := []int{3, 4, 5, 2}
	p1 := &ps
	pointerSlice(p1)
	fmt.Println("2、整数切片的指针：", *p1)

}
