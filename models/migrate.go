package models

import (
	"github.com/smileyan/backend/config"
)

func AutoMigrate() {
	db := config.DB
	db.AutoMigrate(
		&User{},
		&Category{},
		&Tag{},
		&Post{},
		&Page{},
		&Comment{},
		&Subscription{},
	)

	// 创建默认管理员用户
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		defaultAdmin := User{
			Email:    "root@smileyan.cn",
			Nickname: "管理员",
			Role:     RoleAdmin,
			Status:   StatusActive,
		}
		db.Create(&defaultAdmin)
	}
}
