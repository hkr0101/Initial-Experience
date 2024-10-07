package myauth

import (
	"Initial_Experience/db"
	"Initial_Experience/myModels"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET")) // JWT 密钥，可以通过环境变量配置

// Hash 密码
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// 检查密码

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// 创建 JWT

// 创建 JWT，包含用户ID和用户名

func GenerateJWT(user mymodels.User) (string, error) {
	// 设置过期时间为24小时
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.MapClaims{
		"userID":   user.UserID,   // 用户ID
		"username": user.Username, // 用户名
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 注册处理

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

	if !checkPassword(foundUser.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
		return
	}

	// 生成 JWT Token

	tokenString, err := GenerateJWT(foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	// 返回 JWT Token 给客户端
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": tokenString})
}

// 验证JWT Token

func ValidateJWT(tokenStr string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// JWT身份验证中间件

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中提取 JWT
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization token required"})
			c.Abort()
			return
		}

		// 验证并解析 Token
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		// 从 claims 中提取用户信息
		userID := (*claims)["userID"].(float64) // 注意：JWT中的数字会被解析为 float64 类型
		username := (*claims)["username"].(string)

		// 设置用户信息到上下文中，供后续处理使用
		c.Set("userID", int(userID)) // 转为int类型
		c.Set("username", username)

		c.Next()
	}
}

// 登出处理 (JWT无状态不需要特殊登出逻辑，只需在客户端删除Token)

func LogoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// AI需求注册

func RegisterAndChangeAI(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := userID.(int)
	var userAI mymodels.AIRequest
	if err := c.ShouldBindJSON(&userAI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userAI.UserID = id
	if err := db.DB.Save(&userAI).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "save AI"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User_AI registered successfully"})
}
