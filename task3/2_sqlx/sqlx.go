package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

/*
题目1：使用SQL扩展库进行查询

	假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
	要求 ：
	编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
	编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

题目2：实现类型安全映射

	假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
	要求 ：
	定义一个 Book 结构体，包含与 books 表对应的字段。
	编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
*/
type Employee struct {
	ID         int
	Name       string `db:"name"`
	Department string
	Salary     int
}
type Book struct {
	ID     int
	Title  string
	Author string
	Price  float64
}

func main() {

	dsn := "root:asdfasdf@tcp(localhost:3306)/dev?charset=utf8mb4&parseTime=True&loc=Local"
	// 使用sqlx.Connect（包含Open和Ping）
	db, errConn := sqlx.Connect("mysql", dsn)
	if errConn != nil {
		log.Fatalf("连接数据库失败: %v", errConn)
	}
	// 设置连接池参数（重要！）
	db.SetMaxOpenConns(25)                 // 最大打开连接数
	db.SetMaxIdleConns(10)                 // 最大空闲连接数
	db.SetConnMaxLifetime(5 * time.Minute) // 连接最大生命周期

	// 测试查询
	var version string

	if err := db.Get(&version, "SELECT VERSION()"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("数据库版本:", version)

	defer db.Close()

	fmt.Println("数据库连接成功!")
	var empDept []Employee
	if err := db.Select(&empDept, "SELECT * FROM employees WHERE department=?", "技术部"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", empDept)

	var emp []Employee
	if err := db.Select(&emp, "SELECT * FROM `employees` where salary in (select max(salary) from  `employees`)"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", emp)

	//题目2：实现类型安全映射：查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全

	var book []Book
	if err := db.Select(&book, "SELECT * FROM Books where Price > 50 "); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("书：%#v\n", book)

}
