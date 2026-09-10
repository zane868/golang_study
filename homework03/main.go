package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strconv"

	"github.com/zane868/golang_study/homework03/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	chineseStart = 0x4E00
	chineseEnd   = 0x9FFF
)

func main() {
	//创建DB
	db := createDb()

	//清理历史数据
	clean(db)

	//创建用户和文章
	createUserAndPost(db)

	//获取所有用户
	users := findAll(db)

	//进行评论
	comment(users, db)

	//统计最大评论数
	countCommentMax(db)

	//清理评论
	deleComment(users, db)

}

func countCommentMax(db *gorm.DB) []model.CommentGroup {

	var results []model.CommentGroup

	//在实际业务场景中 要根据具体的数据大小来确定是否直接用group
	err := db.Model(&model.Comment{}).
		Select("post_id, COUNT(post_id) AS total").
		Group("post_id").Scan(&results).Error
	if err != nil {
		return nil
	}

	//找出最大评论数
	sort.Slice(results, func(i, j int) bool {
		return results[i].Total > results[j].Total
	})

	if len(results) == 0 {
		return results
	}

	//返回所有评论数等于最大值的文章，主要处理两篇文章并行top1
	maxTotal := results[0].Total
	maxCount := 1
	for maxCount < len(results) && results[maxCount].Total == maxTotal {
		maxCount++
	}
	for _, group := range results[:maxCount] {
		fmt.Println("文章 ID:", group.PostID, "最大评论数是:", group.Total)
	}

	return results[:maxCount]
}

func comment(users []model.User, db *gorm.DB) {

	//进行评论
	user1 := users[0]
	user2 := users[1]
	user3 := users[2]

	user2_post := user2.Posts[0]
	user3_post := user3.Posts[0]

	for i := 0; i < 3; i++ {
		createComment(user1.ID, user2_post.ID, db)
	}
	createComment(user1.ID, user3_post.ID, db)

	//重新加载reload
	u2 := get(user2.ID, db)
	fmt.Println(u2.Posts[0].Comments[0].CreatedAt)

}

func createComment(userId uint, postID uint, db *gorm.DB) {
	comment := model.Comment{
		PostID:  postID,
		UserId:  userId,
		Content: randomChineseString(5, 20),
	}
	if err := db.Create(&comment).Error; err != nil {
		log.Fatal("创建评论失败：", err)
	}
}

func deleComment(users []model.User, db *gorm.DB) {

	user1 := users[0]
	user2 := users[1]

	user2_post := user2.Posts[0]

	if err := db.Delete(&model.Comment{UserId: user1.ID, PostID: user2_post.ID}, "user_id = ? AND post_id = ?", user1.ID, user2_post.ID).Error; err != nil {
		log.Fatal("删除评论失败", err)
	}
}

func createUserAndPost(db *gorm.DB) {
	// 创建用户和文章
	users := make([]model.User, 3)
	for i := 0; i < 3; i++ {
		user := createUser(strconv.Itoa(i) + "号用户")
		post := createPost(user.Name)
		user.Posts = append(user.Posts, *post)
		users[i] = *user
	}

	if err := db.CreateInBatches(&users, 3).Error; err != nil {
		log.Fatal("批量创建用户失败：", err)
	}

}

func createPost(title string) *model.Post {
	return &model.Post{
		Title:   title + "文章测试",
		Content: randomChineseString(),
	}
}

func createUser(name string) *model.User {
	return &model.User{
		Name:  name,
		Email: name + "@example.com",
		Age:   20,
	}
}

func findAll(db *gorm.DB) []model.User {
	users := make([]model.User, 0, 10)
	db.Preload("Posts.Comments").Where("id > 0").Find(&users)
	return users
}

func get(id uint, db *gorm.DB) model.User {
	users := model.User{}
	db.Preload("Posts.Comments").Where("id =?", id).Find(&users)
	return users

}

func clean(db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Comment{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Post{}).Error; err != nil {
			return err
		}
		return tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.User{}).Error
	})
	if err != nil {
		log.Fatal("清理数据失败：", err)
	}
	fmt.Println("清理数据")
}

func createDb() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("data/demo.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("连接数据库失败: %w", err))

	}

	// 根据 User 结构自动创建或更新表
	if err := db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{}); err != nil {
		panic(fmt.Errorf("创建数据表失败：: %w", err))

	}

	return db
}

func randomChineseString(length ...int) string {
	stringLength := randomInt()
	if len(length) > 0 {
		stringLength = length[0]
	}

	if stringLength <= 0 {
		return ""
	}

	result := make([]rune, stringLength)

	for i := range result {
		value, err := rand.Int(
			rand.Reader,
			big.NewInt(chineseEnd-chineseStart+1),
		)
		if err != nil {
			log.Fatal(err)
		}

		result[i] = rune(chineseStart + value.Int64())
	}

	return string(result)
}

func randomInt() int {
	value, err := rand.Int(rand.Reader, big.NewInt(2971))
	if err != nil {
		log.Fatal("生成随机整数失败：", err)
	}

	return int(value.Int64()) + 30
}
