package routes

import (
	"Initial_Experience/AI_answer"
	"Initial_Experience/db"
	"Initial_Experience/myModels"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//创建答案

func CreateAnswer(c *gin.Context) {
	var answer mymodels.Answer
	Curuser := getCurUser(c)
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	qID, _ := strconv.Atoi(c.Param("question_id"))
	answer.QuestionID = qID
	answer.UserID = Curuser.UserID
	db.DB.Create(&answer)
	c.JSON(http.StatusOK, gin.H{"data": answer})
}

//删除答案

func DeleteAnswer(c *gin.Context) {
	Curuser := getCurUser(c)
	id, _ := strconv.Atoi(c.Param("answer_id"))
	var answer = mymodels.Answer{}
	if err := db.DB.Where("answer_id = ?", id).First(&answer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if Curuser.UserID != answer.UserID && Curuser.Username != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	if err := db.DB.Delete(&mymodels.Answer{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

//根据id查找答案

func GetAnswerByID(c *gin.Context) {
	var answer = mymodels.Answer{}
	id, _ := strconv.Atoi(c.Param("answer_id"))
	if err := db.DB.Where("answer_id = ?", id).First(&answer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "answer not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": answer})
}

//当前用户的所有答案

func GetAnswerListByUser(c *gin.Context) {
	Curuser := getCurUser(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	user := Curuser
	var answerList []mymodels.Answer
	if err := db.DB.Where("user_id = ?", user.UserID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&answerList).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "answer not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": answerList})
}

//当前问题的所有答案

func GetAnswerListByQuestion(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	id, _ := strconv.Atoi(c.Param("answer_id"))
	var answerList []mymodels.Answer
	if err := db.DB.Where("question_id = ?", id).Offset((page - 1) * pageSize).Limit(pageSize).Find(&answerList).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "answer not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": answerList})
}

//更新答案

func UpdateAnswer(c *gin.Context) {
	Curuser := getCurUser(c)
	id, _ := strconv.Atoi(c.Param("answer_id"))
	var answer = mymodels.Answer{}
	if err := db.DB.Where("answer_id = ?", id).First(&answer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "answer not found"})
		return
	}

	if Curuser.UserID != answer.UserID && Curuser.Username != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var newAnswer = mymodels.Answer{}
	if err := c.ShouldBind(&newAnswer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newAnswer.AnswerID = answer.AnswerID
	newAnswer.QuestionID = answer.QuestionID
	newAnswer.UserID = answer.UserID
	if err := db.DB.Save(&newAnswer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": newAnswer})
}

// 创建一个ai回答

func CreateAnswerByAI(c *gin.Context) {
	Curuser := getCurUser(c)
	questionId, _ := strconv.Atoi(c.Param("question_id"))
	var AI = mymodels.AIRequest{}
	if err := db.DB.Where("user_id = ?", Curuser.UserID).First(&AI).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "AI not found"})
		return
	}
	var question = mymodels.Question{}
	if err := db.DB.Where("question_id = ?", questionId).First(&question).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "question not found"})
		return
	}
	answer_content, err := AI_answer.CallAI(question.Content, AI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var answer mymodels.Answer
	answer.QuestionID = questionId
	answer.UserID = question.UserID
	answer.Content = answer_content
	db.DB.Create(&answer)
	c.JSON(http.StatusOK, gin.H{"data": answer})
}
