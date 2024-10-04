package mymodels

// Question 结构体
type Question struct {
	QuestionID int    `json:"question_id" gorm:"primaryKey;autoIncrement"`
	Content    string `json:"content" binding:"required"`
	Status     string `json:"status" binding:"required"`
	Name       string `json:"name" binding:"required"`
	UserID     int    `json:"user_id"` // 外键，关联用户
}

// User 结构体
type User struct {
	UserID   int    `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Answer struct {
	AnswerID   int    `json:"answer_id" gorm:"primaryKey;autoIncrement"`
	Content    string `json:"content" binding:"required"`
	QuestionID int    `json:"question_id"` // 外键，关联问题
	UserID     int    `json:"user_id"`     // 外键，关联用户
}
