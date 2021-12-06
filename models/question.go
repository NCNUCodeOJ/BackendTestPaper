package models

import (
	"gorm.io/gorm"
)

// Question 題目 data model
type Question struct {
	gorm.Model
	Question   string `gorm:"type:text;"`
	Author     uint   `gorm:"NOT NULL;"`
	Layer      uint   `gorm:"NOT NULL;"`
	Source     string `gorm:"NOT NULL;"`
	Difficulty uint   `gorm:"NOT NULL;"`
	Type       uint   `gorm:"NOT NULL;"`
	// 題目
	// 出題者
	// 層級(校內、區域、全國)
	// 題目出處(學校 id、單位 id)
	// 難易度
	// 類型(多選、單選、填充)
	// 選項/答案
}

// Option 選項 data model
type Option struct {
	gorm.Model
	Content    string `gorm:"type:text;"`
	Answer     bool   `gorm:"NOT NULL;"`
	QuestionID uint   `gorm:"NOT NULL;"`
	Sort       uint   `gorm:"NOT NULL;"`
	// 內容
	// 是否為正確答案
	// 對應的題目
	// 這是第幾個選項(若為填充題則填 -1)
}

// CreateQuestion 新增題目
func CreateQuestion(question *Question) (err error) {
	err = DB.Create(&question).Error
	return
}

// CreateOption 新增選項/答案
func CreateOption(option *Option) (err error) {
	err = DB.Create(&option).Error
	return
}

// ListQuestions 取得所有 Question
func ListQuestions() (questions []Question, err error) {
	err = DB.Find(&questions).Error
	return
}

// GetQuestion  透過 id 取得 question
func GetQuestion(id uint) (Question, error) {
	var question Question
	if err := DB.Where("id = ?", id).First(&question).Error; err != nil {
		return Question{}, err
	}
	return question, nil
}

// ListOptionsByQuestionID 透過 questionID 取得該題目下的所有 option
func ListOptionsByQuestionID(questionID uint) (option Option, err error) {
	err = DB.First(&option, questionID).Error
	return
}

// GetAnswerByQuestionID  透過 id 取得 question
func GetAnswerByQuestionID(questionID uint) (Option, error) {
	var option Option
	if err := DB.Where("question_id = ?", questionID).First(&option).Error; err != nil {
		return Option{}, err
	}
	return option, nil
}
