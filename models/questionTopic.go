package models

import "gorm.io/gorm"

// QuestionTopic 選擇題在哪個 Topic
type QuestionTopic struct {
	gorm.Model
	TopicID    uint `gorm:"NOT NULL;"`
	QuestionID uint `gorm:"NOT NULL;"`
	Sort       uint `gorm:"NOT NULL;"`
	Random     bool `gorm:"NOT NULL;"`
	Type       uint `gorm:"NOT NULL;"`
	// 對應的第幾大題
	// 配分
	// 對應的題目 ID
	// 排序(第幾大題下的第幾題)
	// 選項是否隨機呈現(作答時)
	// 題型(選擇或填充)
}

// CreateQuestionTopic 新增填充題們
func CreateQuestionTopic(QuestionTopic *QuestionTopic) {
	DB.Create(&QuestionTopic)
}

// ListQuestionTopicsByTopicID 用 topic_id 取得所有 QuestionTopic
func ListQuestionTopicsByTopicID(topicID uint) (questionTopics []QuestionTopic, err error) {
	err = DB.Table("question_topics").Where("topic_id = ?", topicID).Find(&questionTopics).Error
	return
}

// GetQuestionTopicByID 透過 ID 取得 QuestionTopic
func GetQuestionTopicByID(id uint) (QuestionTopic, error) {
	var QuestionTopic QuestionTopic
	if err := DB.Where("id = ?", id).First(&QuestionTopic).Error; err != nil {
		return QuestionTopic, err
	}
	return QuestionTopic, nil
}

// GetQuestionTopic 透過 ID 取得 QuestionTopic
func GetQuestionTopic(topicID uint, questionID uint) (QuestionTopic, error) {
	var QuestionTopic QuestionTopic
	if err := DB.Where("topic_id = ?", topicID).Where("question_id = ?", questionID).First(&QuestionTopic).Error; err != nil {
		return QuestionTopic, err
	}
	return QuestionTopic, nil
}

// DeleteQuestionTopic 刪除
func DeleteQuestionTopic(QuestionTopic QuestionTopic) {
	DB.Where("id = ?", QuestionTopic.ID).Delete(&QuestionTopic)
}

// GetQuestionTopicByTopicIDandSort 透過 TopicID 和 Sort 取得 QuestionTopic
func GetQuestionTopicByTopicIDandSort(topicID uint, sort uint) (QuestionTopic, error) {
	var QuestionTopic QuestionTopic
	if err := DB.Where("topic_id = ?", topicID).Where("sort = ?", sort).First(&QuestionTopic).Error; err != nil {
		return QuestionTopic, err
	}
	return QuestionTopic, nil
}
