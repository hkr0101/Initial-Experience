package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// Question 结构体
type Question struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Content  string `json:"content"`
	Status   string `json:"status"`
	Name     string `json:"name"`
	Tags     string `json:"tags"`
	UserName string `json:"user_name"`
}

// User 结构体
type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var DB *gorm.DB

var Logintate bool = false

func main() {
	r := gin.Default()
	Connect()
	Migrate()
	//用户注册
	r.POST("/register", registerHandler)
	// 用户登录
	r.POST("/login", loginHandler)
	// 用户登出
	r.POST("/logout", logoutHandler)

	// 问题管理
	auth := r.Group("/questions")
	auth.Use(AuthMiddleware()) // 使用身份验证中间件
	{
		auth.POST("", createQuestion)
		auth.DELETE("/:id", deleteQuestion)
		auth.PUT("/:id", updateQuestion)
		auth.GET("", getQuestions)
		auth.GET("/:id", getQuestionByID)
	}

	r.Run(":8080")
}

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
	err := DB.AutoMigrate(&User{}, &Question{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
}

// 注册
func registerHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 检查用户名是否已存在
	var existingUser User
	if err := DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Username already exists"})
		return
	}

	// 创建新用户
	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "id": user.Username})
}

// 登录处理
func loginHandler(c *gin.Context) {
	if Logintate {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already logged"})
		return
	}
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundUser User
	if err := DB.Where("username = ?", user.Username).First(&foundUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	if foundUser.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	// 登录成功，可以设置用户身份信息到上下文中
	//fmt.Print(foundUser.Username)
	//c.Set("user", foundUser.Username)
	Logintate = true
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// 身份验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//user11, exists := c.Get("user")
		//fmt.Println("User:", user11)
		if !Logintate {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 登出处理
func logoutHandler(c *gin.Context) {
	Logintate = false // 重置登录状态
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// 创建问题
func createQuestion(c *gin.Context) {
	var question Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Create(&question)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "question": question})
}

// 删除问题
func deleteQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := DB.Delete(&Question{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// 修改问题
func updateQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var question Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	question.ID = uint(id)
	if err := DB.Save(&question).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "question": question})
}

// 获取所有问题
func getQuestions(c *gin.Context) {
	var questions []Question
	DB.Find(&questions)
	c.JSON(http.StatusOK, questions)
}

// 根据 ID 查询问题
func getQuestionByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var question Question
	if err := DB.First(&question, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	c.JSON(http.StatusOK, question)
}
