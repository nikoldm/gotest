package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

//题目1：模型定义
//	假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
//	要求 ：
//		使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
//		编写Go代码，使用Gorm创建这些模型对应的数据库表。
//题目2：关联查询
//	基于上述博客系统的模型定义。
//	要求 ：
//		编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
//		编写Go代码，使用Gorm查询评论数量最多的文章信息。
//题目3：钩子函数
//	继续使用博客系统的模型。
//	要求 ：
//		为Post模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
//		为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。

type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(32);not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Age       uint8  `gorm:"check:age>0"`
	PostCount uint   `gorm:"default:0"`
	Posts     []Post `gorm:"foreignKey:UserID"` // 一对多：一个用户多篇文章
}

type Post struct {
	gorm.Model
	Title         string    `gorm:"type:varchar(64);not null"`
	PostContent   string    `gorm:"type:text;not null"`
	CommentCount  uint      `gorm:"default:0"`
	Status        string    `gorm:"type:varchar(32);default:'published';not null'"`
	CommentStatus string    `gorm:"type:varchar(32);default:'无评论';not null'"`
	Comments      []Comment `gorm:"foreignKey:PostID"` // 一对多：一篇文章多个评论
	User          User      `gorm:"foreignKey:UserID"` // 属于关系
	UserID        uint      `gorm:"not null;index"`    // 外键指向User
}

type Comment struct {
	gorm.Model
	Content     string    `gorm:"type:text;not null"`
	CommentTime time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;not null"`
	UserID      uint      `gorm:"not null;index"` //评论人的id
	User        User      `gorm:"foreignKey:UserID"`
	PostID      uint      `gorm:"not null;index"`    // 外键，指向Post
	Post        Post      `gorm:"foreignKey:PostID"` //属于关系
}

// 全局连接；
var db *gorm.DB

// 初始化数据库连接
func initDB() error {
	dsn := "root:asdfasdf@tcp(localhost:3306)/dev?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印sql
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute * 10)

	fmt.Println("数据库连接成功")
	return nil
}

func createSampleData() error {
	// 创建user数据，无法向 ‘create’ 传递结构体，因此你应该传入数据的指针.
	users := []User{
		{Username: "张三", Email: "zhangsan@google.com", Password: "123456", Age: 18},
		{Username: "李四", Email: "lisi@google.com", Password: "123456677", Age: 25},
		{Username: "王五", Email: "wangwu@google.com", Password: "888888", Age: 36},
		{Username: "狄仁杰", Email: "direnjie@google.com", Password: "9997888", Age: 55},
	}
	if err := db.Create(&users).Error; err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	// 创建文章
	posts := []*Post{
		{Title: "Go语言入门指南", PostContent: "Go语言的基础语法介绍...", UserID: 1},
		{Title: "GORM使用技巧", PostContent: "深入讲解GORM的高级特性...", UserID: 1},
		{Title: "Web开发实践", PostContent: "使用Go构建RESTful API...", UserID: 2},
		{Title: "数据库设计", PostContent: "关系型数据库设计原则...", UserID: 3},
	}
	if err := db.Create(posts).Error; err != nil {
		return fmt.Errorf("创建文章失败: %w", err)
	}

	// 创建评论
	comments := []Comment{
		{Content: "好文章，学习了！", UserID: 2, PostID: 1},
		{Content: "期待更多实战例子", UserID: 3, PostID: 1},
		{Content: "GORM真的很强大", UserID: 1, PostID: 2},
		{Content: "对初学者很有帮助", UserID: 2, PostID: 2},
		{Content: "大神的理解十分顶级啊", UserID: 2, PostID: 2},
		{Content: "API设计部分讲得很好", UserID: 3, PostID: 3},
		{Content: "感谢分享！", UserID: 1, PostID: 4},
	}
	if err := db.CreateInBatches(comments, 3).Error; err != nil {
		return fmt.Errorf("创建评论失败: %w", err)
	}

	return nil
}

// 查询用户的所有文章及评论, 使用Joins和Select
func getUserPostAndComment1(userID uint) (posts []Post, err error) {

	err = db.Table("posts").
		Select(`
            posts.*,
            users.username as author_name,
            users.email as author_email,
            COUNT(comments.id) as comment_count
        `).
		Joins("LEFT JOIN users ON users.id = posts.user_id").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Where("posts.user_id = ?", userID).
		Group("posts.id").
		Order("posts.created_at DESC").
		Scan(&posts).Error

	return posts, err
}
func getUserPostAndComment(userID uint) (posts []Post, err error) {

	var user User
	err = db.Preload("Posts.Comments.User").Preload("Posts.User").Preload(clause.Associations).First(&user, userID).Error
	return user.Posts, err

	// 另一种写法
	//if err = db.Where("user_id = ?", userID).
	//	Preload("User", func(db *gorm.DB) *gorm.DB {
	//		return db.Select("id, username, email")
	//	}).
	//	Preload("Comments", func(db *gorm.DB) *gorm.DB {
	//		return db.Preload("User", func(db *gorm.DB) *gorm.DB {
	//			return db.Select("id, username")
	//		})
	//	}).
	//	Find(&posts).Error; err != nil {
	//	return nil, fmt.Errorf("查询失败: %w", err)
	//}
	//return posts, nil

}

// 打印查询结果
func printUserPostAndComment(userID uint) {
	if posts, err := getUserPostAndComment(userID); err == nil {
		for _, post := range posts {
			fmt.Printf("\n文章：%s (ID:%d)\n", post.Title, post.ID)
			fmt.Printf("作者：%s\n", post.User.Username)
			fmt.Printf("评论数：%d\n", len(post.Comments))

			for i, comment := range post.Comments {
				fmt.Printf("  %d. %s - 评论者：%s\n",
					i+1, comment.Content[:min(30, len(comment.Content))],
					comment.User.Username)
			}
		}
	} else {
		fmt.Errorf("查询用户发表的文章和评论失败：%v", err.Error())
	}

}

// 查询评论最多的文章 使用原生SQL（性能更好）
func getPostMostComment() (post Post, err error) {

	err = db.Raw("SELECT p.id,p.Title, u.username, u.email,COUNT(c.id) as comment_count " +
		"FROM posts p LEFT JOIN comments c ON p.id = c.post_id LEFT JOIN users u ON p.user_id = u.id " +
		"GROUP BY p.id ORDER BY comment_count DESC LIMIT 1").Scan(&post).Error

	db.Where("post_id = ?", post.ID).Find(&post.Comments)

	return post, err

}
func getPostMostComment1() (post Post, err error) {
	err = db.Select("posts.ID,posts.Title,count(c.id) as CommentCount ").Preload("User").Preload("Comments").Joins("LEFT JOIN comments c ON c.post_id = posts.id").Group("posts.id").Order("count(c.id) desc").First(&post).Error
	return post, err
}

//钩子函数============================================================================================

// 开始事务
// BeforeSave
// BeforeCreate
// BeforeDelete
// 关联前的 save
// 插入记录至 db
// 关联后的 save
// AfterCreate
// AfterSave
// 提交或回滚事务

// AfterCreate Post的钩子函数：创建后更新 用户的 文章数量
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	// 后更新 用户的 文章数量
	if err = tx.Model(&User{}).Where("ID=?", p.UserID).UpdateColumn("post_count", gorm.Expr("post_count+1")).Error; err != nil {
		return fmt.Errorf("更新用户文章数量失败: %w", err.Error())
	}

	fmt.Printf("文章创建成功，已更新用户%d的文章数量\n", p.UserID)
	return nil
}

// BeforeDelete Post的钩子函数：删除前减少用户文章数量
func (p *Post) BeforeDelete(tx *gorm.DB) (err error) {
	if err = tx.First(&p, p.ID).Error; err != nil {
		return fmt.Errorf("文章不存在：%w", err.Error())
	}
	if err = tx.Model(&User{}).Where("ID=?", p.UserID).UpdateColumn("post_count", gorm.Expr("post_count-1")).Error; err != nil {
		return fmt.Errorf("更新用户文章数量失败: %w", err.Error())
	}
	fmt.Printf("文章删除前，已更新用户%d的文章数量\n", p.UserID)
	return nil
}

// 更新Comment模型定义，添加钩子函数
// Comment的钩子函数：创建后更新文章评论数量
func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	if err = tx.Model(&Post{}).Where("id=?", c.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count+1")).UpdateColumn("comment_status", "有评论").Error; err != nil {
		return fmt.Errorf("更新文章评论数量失败: %w", err.Error)
	}
	fmt.Printf("评论创建成功，已更新文章%d的评论数量\n", c.PostID)
	return nil
}

// Comment的钩子函数：删除时检查文章评论数量
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	// 删除后，获取评论数量
	var commentCount int64
	commentStatus := "有评论"
	tx.Model(&Comment{}).Where("Post_id=?", c.PostID).Count(&commentCount)
	if commentCount <= 0 {
		// 更新文章的评论数量
		commentStatus = "无评论"
	}
	tx.Model(&Post{}).Where("ID=?", c.PostID).UpdateColumn("comment_count", commentCount).UpdateColumn("comment_status", commentStatus)
	fmt.Printf("评论删除成功，文章%d还有%d条评论\n", c.PostID, commentCount)
	return nil
}

// 批量操作时钩子不会触发，需要手动处理:  1、单个创建，触发钩子。2、手动调用钩子函数  ?

func main() {
	// 1. 初始化数据库连接
	if err := initDB(); err != nil {
		panic(err)
	}
	// 关闭连接
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	//2. 创建表
	if err := db.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
		panic("自动创建表失败")
	}
	fmt.Println("自动创建表成功")

	// 3. 创建示例数据
	if err := createSampleData(); err != nil {
		fmt.Println("表数据已存在，继续执行")
		//panic(err)
	}

	// 查询用户1的所有文章及评论
	printUserPostAndComment(1)

	// 查询评论最多的文章
	post, err := getPostMostComment()
	if err != nil {
		panic(err)
	}
	fmt.Printf("评论最多的文章：%v,评论数是：%v, 分别是:\n", post.Title, post.CommentCount)
	for i, comment := range post.Comments {
		fmt.Printf("%v、%v\n", i+1, comment.Content)
	}

	// 5. 测试钩子函数
	//if err := testHooks(); err != nil {
	//	fmt.Printf("钩子测试失败: %v\n", err)
	//}

	fmt.Println("\n 所有操作完成！")

}

// 测试钩子函数
func testHooks() error {
	fmt.Println("\n测试钩子函数...")

	// 测试1：创建文章，触发AfterCreate钩子
	newPost := Post{
		Title:       "钩子函数测试",
		PostContent: "测试GORM钩子函数的功能",
		UserID:      1,
	}

	if err := db.Create(&newPost).Error; err != nil {
		return fmt.Errorf("创建文章失败: %w", err)
	}

	// 验证用户文章数量
	var user User
	db.First(&user, 1)
	fmt.Printf("用户%s的文章数量：%d\n", user.Username, user.PostCount)

	// 测试2：创建评论，触发Comment的AfterCreate
	newComment := Comment{
		Content: "测试评论钩子",
		UserID:  2,
		PostID:  newPost.ID,
	}

	if err := db.Create(&newComment).Error; err != nil {
		return fmt.Errorf("创建评论失败: %w", err)
	}

	// 验证文章评论数量
	var post Post
	db.First(&post, newPost.ID)
	fmt.Printf("文章%s的评论数量：%d，状态：%s\n",
		post.Title, post.CommentCount, post.CommentStatus)

	// 测试3：删除评论，触发AfterDelete
	if err := db.Delete(&newComment).Error; err != nil {
		return fmt.Errorf("删除评论失败: %w", err)
	}

	// 再次验证文章状态
	db.First(&post, newPost.ID)
	fmt.Printf("删除评论后，文章状态：%s，评论数量：%d\n",
		post.CommentStatus, post.CommentCount)

	// 测试4：删除文章，触发BeforeDelete
	if err := db.Delete(&newPost).Error; err != nil {
		return fmt.Errorf("删除文章失败: %w", err)
	}

	// 验证用户文章数量
	db.First(&user, 1)
	fmt.Printf("删除文章后，用户%s的文章数量：%d\n", user.Username, user.PostCount)

	return nil
}
