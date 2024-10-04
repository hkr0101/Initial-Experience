package main

import (
	"Initial_Experience/AI_answer"
	"Initial_Experience/db"
	"Initial_Experience/myauth"
	"Initial_Experience/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db.Connect()
	db.Migrate()
	//用户注册
	r.POST("/register", myauth.RegisterHandler)
	// 用户登录
	r.POST("/login", myauth.LoginHandler)
	// 用户登出
	r.POST("/logout", myauth.LogoutHandler)
	// 问题管理

	//查看所有问题
	r.GET("/questions", routes.GetQuestions)
	//查看某个问题
	r.GET("/questions/:question_id", routes.GetQuestionByID)
	//查看某个问题的所有答案
	r.GET("/questions/:question_id/answer", routes.GetAnswerListByQuestion)
	//查看某个问题的某个答案
	r.GET("/questions/answer/:answer_id", routes.GetAnswerByID)

	auth := r.Group("/my")
	auth.Use(myauth.AuthMiddleware()) // 使用身份验证中间件
	{
		//创建问题
		auth.POST("/questions", routes.CreateQuestion)
		//创建答案
		auth.POST("/questions/:question_id/answer", routes.CreateAnswer)
		//删除问题
		auth.DELETE("/questions/:question_id", routes.DeleteQuestion)
		//删除答案
		auth.DELETE("/questions/answer/:answer_id", routes.DeleteAnswer)
		//更新问题
		auth.PUT("/questions/:question_id", routes.UpdateQuestion)
		//更新答案
		auth.PUT("/:answer_id", routes.UpdateAnswer)
		//给出当前用户的所有答案
		auth.GET("/answer", routes.GetAnswerListByUser)
		//给出当前用户的所有问题
		auth.GET("/questions", routes.GetQuestionByUser)
		//登出
		auth.POST("/logout", myauth.LogoutHandler)
		//调用ai，未完成
		auth.POST("/chat", aianswer.ChatGPTHandler)
	}
	r.Run(":8080")
}
