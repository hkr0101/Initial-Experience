package routes

import (
	"Initial_Experience/db"
	"Initial_Experience/myauth"
	"Initial_Experience/mymodels"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var user = myauth.Curuser

// 创建问题

func CreateQuestion(c *gin.Context) {
	var question mymodels.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	question.UserID = user.UserID
	//fmt.Print(question)
	db.DB.Create(&question)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "question": question})
}

// 删除问题

func DeleteQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("question_id"))
	var question mymodels.Question
	//根据当前问题的id找到问题的具体内容
	if err := db.DB.Where("question_id = ?", id).First(&question).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	//如果不是你的问题
	if myauth.Curuser.UserID != question.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	//删除问题
	if err := db.DB.Delete(&mymodels.Question{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// 修改问题

func UpdateQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("question_id"))
	var question mymodels.Question
	//根据当前问题的id找到问题的具体内容
	if err := db.DB.Where("question_id = ?", id).First(&question).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	//如果不是你的问题
	if myauth.Curuser.UserID != question.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	//构建一个新的问题
	var newQuestion mymodels.Question
	if err := c.ShouldBindJSON(&newQuestion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newQuestion.UserID = question.UserID
	newQuestion.QuestionID = question.QuestionID

	if err := db.DB.Save(&question).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "question": question})
}

// 获取所有问题

func GetQuestions(c *gin.Context) {
	var questions []mymodels.Question
	db.DB.Find(&questions)
	c.JSON(http.StatusOK, questions)
}

// 根据 ID 查询问题

func GetQuestionByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("question_id"))
	var question mymodels.Question
	if err := db.DB.First(&question, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	c.JSON(http.StatusOK, question)
}

//获取当前用户的所有问题

func GetQuestionByUser(c *gin.Context) {
	var question []mymodels.Question
	myID := myauth.Curuser.UserID
	if err := db.DB.Where("user_id = ?", myID).Find(&question).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	c.JSON(http.StatusOK, question)
}
