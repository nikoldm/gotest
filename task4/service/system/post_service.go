package system

import (
	"fmt"
	"task4/config"
	"task4/model/system"
)

type PostService struct{}

func (P *PostService) CreatePost(post *system.Post) (*system.Post, error) {
	if err := config.DB.Create(&post).Error; err != nil {
		return nil, err
	}
	// 预加载用户信息
	config.DB.Preload("User").First(&post, post.ID)
	return post, nil
}

// UpdatePost 更新文章
func (P *PostService) UpdatePost(post *system.Post) error {
	var queryPost system.Post
	if err := config.DB.First(&queryPost, post.ID).Error; err != nil {
		return err
	}

	// 检查是否是文章作者
	if post.UserID != queryPost.UserID {
		return fmt.Errorf("修改的文章不属于本人")
	}

	queryPost.Title = post.Title
	queryPost.PostContent = post.PostContent
	if err := config.DB.Save(&queryPost).Error; err != nil {
		return err
	}

	// 预加载用户信息
	config.DB.Preload("User").First(&post, post.ID)
	return nil
}

func (P *PostService) DeletePost(postID uint64, userID any) error {
	var post system.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		return fmt.Errorf("文章不存在：" + err.Error())
	}

	// 检查是否是文章作者
	if post.UserID != userID.(uint) {
		return fmt.Errorf("文章不属于作者")
	}

	// 删除文章（软删除）
	if err := config.DB.Delete(&post).Error; err != nil {
		return fmt.Errorf("文章删除失败：" + err.Error())
	}
	return nil
}

func (P *PostService) GetPost(postID uint64) (system.Post, error) {
	var post system.Post
	if err := config.DB.Preload("User").Preload("Comments.User").First(&post, postID).Error; err != nil {
		return post, fmt.Errorf("文章查询失败：" + err.Error())
	}
	return post, nil
}

func (P *PostService) GetPosts(page int, size int, total *int64) (posts []system.Post, err error) {
	offset := (page - 1) * size
	if err = config.DB.Preload("User").Order("created_at DESC").Limit(size).Offset(offset).Find(&posts).Error; err != nil {
		return posts, err
	}

	config.DB.Model(&system.Post{}).Count(total)
	return posts, err
}
