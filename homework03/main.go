package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"homework01/homework03/model"

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

	user1 := users[0]
	user2 := users[1]

	user2_post := user2.Posts[0]

	for i := 0; i < 3; i++ {
		createComment(user1.ID, user2_post.ID, db)

	}

}

func createComment(UserId uint, postID uint, db *gorm.DB) {
	comment := model.Comment{
		PostID:  postID,
		UserId:  UserId,
		Content: randomChineseString(5, 20),
	}
	if err := db.Create(&comment).Error; err != nil {
		log.Fatal("创建评论失败：", err)
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
