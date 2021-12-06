package models

import (
	"gorm.io/gorm"
)

type StudentTestPaperAnswer struct {
	gorm.Model
	StudentTestPaperID uint   `gorm:"NOT NULL;"`
	TopicSort          uint   `gorm:"NOT NULL;"`
	QuestionSort       uint   `gorm:"NOT NULL;"`
	Resort             uint   `gorm:"NOT NULL;"`
	Content            string `gorm:"NOT NULL;"`
	Correct            bool   `gorm:"NOT NULL;"`
	// 答題者(學生)
	// 對應的大題
	// 對應的題目
	// 重新排序
	// 回答內容
	// 是否正確
}

// GetStudentTestPaperAnswersByID 透過 id 取得已被批改的題目
func GetStudentTestPaperAnswerByID(id uint) (studentTestpaperAnswers StudentTestPaperAnswer, err error) {
	err = DB.Where("id = ?", id).First(&studentTestpaperAnswers).Error
	return
}

// // GetAnswerByQuestionID 用 question_id 找 option 裡的 answer
// func GetStudentTestPaperAnswerByQuestionID(questionID uint) (option Option, err error) {
// 	err = DB.Table("options").Where("topic_id = ?", questionID).Where("answer = ?", true).First(&option).Error
// 	return
// }

// CreateStudentTestPaperAnswer 新增該學生已被批改的題目
func CreateStudentTestPaperAnswer(studentTestpaperAnswer *StudentTestPaperAnswer) {
	// testpaper.GET("/:testpaperID/graded", view.GetTestPaperByID)
	DB.Create(&studentTestpaperAnswer)
}

// ListStudentTestPaperAnswersByStudentTestPaperID 透過 studentTestPaperID 取得該學生所有已批改的題目
func ListStudentTestPaperAnswersByStudentTestPaperID(studentTestPaperID uint) (studentTestpaperAnswers []StudentTestPaperAnswer, err error) {
	err = DB.Where("student_testpaper_id = ?", studentTestPaperID).Find(&studentTestpaperAnswers).Error
	return
}
