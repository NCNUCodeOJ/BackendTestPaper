package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Topic 第幾大題(要放在測驗卷的)
type Topic struct {
	gorm.Model
	TestPaperID  uint    `gorm:"NOT NULL;"`
	Distribution float64 `gorm:"NOT NULL;"`
	Sort         uint    `gorm:"NOT NULL;"`
	// 對應的測驗卷
	// 配分
	// 排序(這是第幾大題)
}

// CreateTopic 新增大題
func CreateTopic(topic *Topic, questions []*uint) (err error) {
	err = DB.Create(&topic).Error
	if err != nil {
		return
	}
	for i, questionID := range questions {
		var question2Topic QuestionTopic

		question2Topic.QuestionID = *questionID
		question2Topic.TopicID = topic.ID
		question2Topic.Sort = uint(i + 1)
		err = DB.Create(&question2Topic).Error
		if err != nil {
			return
		}
	}
	return
}

// GetTestpaperTopicCount 取得測驗卷的大題數
func GetTestpaperTopicCount(testpaperID uint) (count int64, err error) {
	err = DB.Model(&Topic{}).Where(&Topic{TestPaperID: testpaperID}).Count(&count).Error
	return
}

// ListTopicsByTestpaperID 取得所有 topic
func ListTopicsByTestpaperID(testpaperID uint) (topics []Topic, err error) {
	err = DB.Model(&Topic{}).Where(&Topic{TestPaperID: testpaperID}).Find(&topics).Error
	return
}

// GetTopicBySort 透過 sort 取得 topic
func GetTopicBySort(testpaperID uint, sort uint) (Topic, error) {
	var topic Topic
	if err := DB.Model(&Topic{}).Where(&Topic{TestPaperID: testpaperID, Sort: sort}).First(&topic).Error; err != nil {
		return Topic{}, err
	}
	return topic, nil
}

// GetTopicDataBySort 取得測驗卷的大題資料
func GetTopicDataBySort(testpaperID uint, sort uint) (gin.H, error) {
	var topic Topic
	var question2Topic []QuestionTopic
	var data = gin.H{}

	if err := DB.Model(&Topic{}).Where(&Topic{TestPaperID: testpaperID, Sort: sort}).First(&topic).Error; err != nil {
		return data, err
	}

	data = gin.H{
		"topic_id":     topic.ID,
		"distribution": topic.Distribution,
		"sort":         topic.Sort,
		"questions":    []gin.H{},
	}

	if err := DB.Model(&QuestionTopic{}).
		Where(&QuestionTopic{TopicID: topic.ID}).Find(&question2Topic).Error; err != nil {
		return data, err
	}
	for _, question := range question2Topic {
		questionData := gin.H{
			"question_id": question.QuestionID,
			"sort":        question.Sort,
		}
		data["questions"] = append(data["questions"].([]gin.H), questionData)
	}

	return data, nil
}

// UpdateTopic 更新
func UpdateTopic(topic *Topic) (err error) {
	err = DB.Where("sort = ?", topic.Sort).Save(&topic).Error
	return
}

// DeleteTopic 刪除
func DeleteTopic(topic Topic) (err error) {
	DB.Where("id = ?", topic.TestPaperID).Delete(&topic)
	return
}
