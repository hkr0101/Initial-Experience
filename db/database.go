package db

import (
	"Initial_Experience/myModels"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// 数据库的连接
func Connect() {
	var err error
	dsn := "root:123456@tcp(localhost:3306)/initial_experience"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	} else {
		log.Println("Successfully connected to the database")
	}
}

// 自动迁移模型
func Migrate() {
	err := DB.AutoMigrate(&mymodels.User{}, &mymodels.Question{}, &mymodels.Answer{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
}
