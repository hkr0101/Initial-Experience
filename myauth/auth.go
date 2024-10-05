package myauth

import (
	"Initial_Experience/db"
	"Initial_Experience/myModels"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

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

	user.Password, _ = hashPassword(user.Password)

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "id": user.Username})
}

// 登录处理

func LoginHandler(c *gin.Context) {
	var onlineUser mymodels.OnlineUser
	var user mymodels.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var foundUser mymodels.User
	if err := db.DB.Where("username = ?", user.Username).First(&foundUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username"})
		return
	}

	check := checkPassword(foundUser.Password, user.Password)

	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
		return
	}
	// 匹配成功，可以设置用户身份信息到上下文中
	// 重复登录
	//fmt.Println(user.UserID)
	if err := db.DB.Where("user_id = ?", foundUser.UserID).First(&onlineUser).Error; err != nil {
		onlineUser.UserID = foundUser.UserID
		fmt.Println(onlineUser)
		db.DB.Create(&onlineUser)
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"message": "User already logged in"})

}

// 身份验证中间件

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("my_id"))
		var curOnlineUser mymodels.OnlineUser
		if err := db.DB.Where("user_id = ?", id).First(&curOnlineUser).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 登出处理

func LogoutHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("my_id"))
	if err := db.DB.Delete(&mymodels.OnlineUser{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// AI需求注册
func RegisterAndChangeAI(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("my_id"))
	var user_AI mymodels.AIRequest
	if err := c.ShouldBindJSON(&user_AI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user_AI.UserID = id
	if err := db.DB.Save(&user_AI).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "save AI"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User_AI registered successfully"})
}
