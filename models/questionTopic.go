package models

import "gorm.io/gorm"

// QuestionTopic 選擇題在哪個 Topic
type QuestionTopic struct {
	gorm.Model
	TopicID      uint    `gorm:"NOT NULL;"`
	Distribution float64 `gorm:"NOT NULL;"`
	QuestionID   uint    `gorm:"NOT NULL;"`
	Sort         uint    `gorm:"NOT NULL;"`
	Random       bool    `gorm:"NOT NULL;"`
	Type         uint    `gorm:"NOT NULL;"`
	// 對應的第幾大題
	// 配分
	// 對應的題目 ID
	// 排序(第幾大題下的第幾題)
	// 選項是否隨機呈現(作答時)
	// 題型(選擇或填充)
}

// AddQuestionTopic 新增填充題們
func AddQuestionTopic(QuestionTopic *QuestionTopic) {
	DB.Create(&QuestionTopic)
}

// GetQuestionTopic 透過 ID 取得 QuestionTopic
// 從 Topic 那邊 Get 了，所以不用 Get
func GetQuestionTopic(id uint) (QuestionTopic, error) {
	var QuestionTopic QuestionTopic
	if err := DB.Where("id = ?", id).First(&QuestionTopic).Error; err != nil {
		return QuestionTopic, err
	}
	return QuestionTopic, nil
}

// ChangeQuestionTopic 修改
func ChangeQuestionTopic(QuestionTopic QuestionTopic) {
	DB.Where("id = ?", QuestionTopic.ID).Save(&QuestionTopic)
}

// DeleteQuestionTopic 刪除
func DeleteQuestionTopic(QuestionTopic QuestionTopic) {
	DB.Where("id = ?", QuestionTopic.ID).Delete(&QuestionTopic)
}
