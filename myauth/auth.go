package myauth

import (
	"Initial_Experience/db"
	"Initial_Experience/mymodels"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Loginstate bool = false
var Curuser mymodels.User = mymodels.User{}

// 注册
func RegisterHandler(c *gin.Context) {
	var user mymodels.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 检查用户名是否已存在
	var existingUser mymodels.User
	if err := db.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Username already exists"})
		return
	}

	// 创建新用户
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "id": user.Username})
}

// 登录处理
func LoginHandler(c *gin.Context) {
	if Loginstate {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already logged"})
		return
	}
	var user mymodels.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundUser mymodels.User
	if err := db.DB.Where("username = ?", user.Username).First(&foundUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	if foundUser.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}
	// 登录成功，可以设置用户身份信息到上下文中
	Loginstate = true
	Curuser = foundUser
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// 身份验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//user11, exists := c.Get("user")
		//fmt.Println("User:", user11)
		if !Loginstate {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 登出处理
func LogoutHandler(c *gin.Context) {
	Loginstate = false // 重置登录状态
	Curuser = mymodels.User{}
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
