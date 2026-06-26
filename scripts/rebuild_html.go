package main

import (
	"fmt"
	"log"

	"github.com/smileyan/backend/config"
	"github.com/smileyan/backend/models"
	"github.com/smileyan/backend/utils"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库
	db := config.InitDatabase()

	fmt.Println("开始重新渲染文章 HTML...")

	// 处理 Post
	var posts []models.Post
	if err := db.Find(&posts).Error; err != nil {
		log.Fatalf("查询文章失败: %v", err)
	}

	updatedPosts := 0
	for i := range posts {
		newHTML := string(utils.MarkdownToHTML([]byte(posts[i].Content)))
		if newHTML != posts[i].HTMLContent {
			posts[i].HTMLContent = newHTML
			if err := db.Save(&posts[i]).Error; err != nil {
				log.Printf("更新文章 [%d] %s 失败: %v", posts[i].ID, posts[i].Title, err)
			} else {
				updatedPosts++
			}
		}
	}
	fmt.Printf("文章: 共 %d 篇，更新 %d 篇\n", len(posts), updatedPosts)

	// 处理 Page
	var pages []models.Page
	if err := db.Find(&pages).Error; err != nil {
		log.Fatalf("查询页面失败: %v", err)
	}

	updatedPages := 0
	for i := range pages {
		newHTML := string(utils.MarkdownToHTML([]byte(pages[i].Content)))
		if newHTML != pages[i].HTMLContent {
			pages[i].HTMLContent = newHTML
			if err := db.Save(&pages[i]).Error; err != nil {
				log.Printf("更新页面 [%d] %s 失败: %v", pages[i].ID, pages[i].Title, err)
			} else {
				updatedPages++
			}
		}
	}
	fmt.Printf("页面: 共 %d 个，更新 %d 个\n", len(pages), updatedPages)

	fmt.Println("处理完成！")
}
