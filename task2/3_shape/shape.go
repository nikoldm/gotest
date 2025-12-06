package main

import (
	"fmt"
	"math"
)

// 定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
// 然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
// 在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
type shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
	height float64
	weight float64
}

type Circle struct {
	radius float64
}

func (r Rectangle) Area() float64 {
	return r.height * r.weight
}

func (c Circle) Area() float64 {
	return c.radius * c.radius * math.Pi
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.height + r.weight)
}

func (c Circle) Perimeter() float64 {
	return 2 * c.radius * math.Pi
}

func main() {
	c := Circle{
		5,
	}
	r := Rectangle{
		2, 5,
	}
	fmt.Println("圆的面积：", c.Area())
	fmt.Println("圆的周长：", c.Perimeter())
	fmt.Println("长方形的面积：", r.Area())
	fmt.Println("长方形的周长：", r.Perimeter())
}
