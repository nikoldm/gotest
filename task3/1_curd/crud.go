package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Student 假设有一个名为 students 的表，对表做curd
type Student struct {
	ID    uint   // Standard field for the primary key
	Name  string // A regular string field
	Age   uint8  // An unsigned 8-bit integer
	Grade string
}

func main() {
	dsn := "root:asdfasdf@tcp(localhost:3306)/dev?charset=utf8mb4&parseTime=True&loc=Local"
	// 连接数据库
	db, errOpen := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if errOpen != nil {
		log.Fatal("连接失败:", errOpen)
	}
	fmt.Println("MySQL连接成功!")

	// 自动迁移（创建表）
	if err := db.AutoMigrate(&Student{}); err != nil {
		log.Println("A账户创建失败:", err.Error())
		return
	}

	//题目1：基本CRUD操作
	//创建记录 ---后面被删除
	student := Student{Name: "张得分", Age: 13, Grade: "三年级"}

	// 1、插入数据 ,批量插入：db.CreateInBatches(users, 100)

	if result := db.Create(&student); result.Error != nil {
		log.Fatal("创建失败:", result.Error)
	}
	fmt.Printf("创建成功，ID: %d\n", student.ID)

	// 2、查询年龄大于18岁的记录
	var users []Student
	if err := db.Where("age >= ?", 18).Find(&users); err.Error != nil {
		log.Println("查询失败:", err.Error)
		return
	}
	for _, u := range users {
		fmt.Printf("ID: %d, Name: %s, Age:%d, grade:%s\n", u.ID, u.Name, u.Age, u.Grade)
	}

	//3、编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	if err := db.Model(&Student{}).Where("Name=?", "张三").Update("Grade", "四年级"); err.Error != nil {
		log.Println("更新失败:", err.Error)
	}

	//4、编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。

	if err := db.Where("Age<?", 15).Delete(&Student{}); err.Error != nil {
		log.Println("删除失败:", err.Error)
	}

	//==== 题目2 ========================================

	//题目2：事务语句:编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作
	//准备新建两个表，初始化数据A的id=1，B的id=2

	if err := db.AutoMigrate(&Account{}); err != nil {
		log.Println("A账户创建失败:", err.Error())
	}
	if err := db.AutoMigrate(&Transaction{}); err != nil {
		log.Println("B账户创建失败:", err.Error())
	}
	//if result := db.Create(&Account{1, 1000}); result.Error != nil {
	//	log.Fatal("A初始化失败:", result.Error)
	//}
	//if result := db.Create(&Account{2, 1000}); result.Error != nil {
	//	log.Fatal("B初始化失败:", result.Error)
	//}

	accTran := 100.0
	resErr := db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		//1、先检查账户 A 的余额是否足够，
		var accA Account
		if err := db.First(&accA, 1); err.Error != nil || err.RowsAffected == 0 {
			// 返回任何错误都会回滚事务
			return fmt.Errorf("A账户查询失败：%w", err.Error)
		}
		if accA.Balance < accTran {
			return fmt.Errorf("用户余额不足")
		}
		//2、如果足够则从账户 A 扣除 100 元，
		if err := tx.Model(&accA).Update("Balance", accA.Balance-accTran); err.Error != nil {
			return fmt.Errorf("A账户扣除100元失败：%v", err.Error)
		}

		//3、向账户 B 增加 100 元
		var accB Account
		if err := db.First(&accB, 2); err.Error != nil {
			return fmt.Errorf("B账户查询失败：%s", err.Error)
		}
		if err := tx.Model(&accB).Update("Balance", accB.Balance+accTran); err.Error != nil {
			return fmt.Errorf("B账户更新失败：%s", err.Error)
		}

		// 4、向交易表记录一笔
		trans := Transaction{}
		trans.FromAccountId = accA.ID
		trans.ToAccountId = accB.ID
		trans.Amount = accTran
		// 返回任何错误都会回滚事务
		if err := tx.Create(&trans).Error; err != nil {
			return fmt.Errorf("交易记录失败：%s", err.Error)
		}
		// 返回 nil 提交事务
		return nil
	})

	if resErr != nil {
		fmt.Printf("A 向账户 B 转账 100 元失败：%s", resErr)
		return
	}
	fmt.Println("A 向账户 B转账100元成功")
}

type Account struct {
	ID      uint
	Balance float64
}

type Transaction struct {
	ID            uint
	FromAccountId uint
	ToAccountId   uint
	Amount        float64
}
