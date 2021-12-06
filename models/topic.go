package models

import (
	"gorm.io/gorm"
)

// 第幾大題(要放在測驗卷的)
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
func CreateTopic(topic *Topic) {
	DB.Create(&topic)
}

// ListTopicsByTestpaperID 取得所有 topic
func ListTopicsByTestpaperID(testpaperID uint) (topics []Topic, err error) {
	err = DB.Where("testpaper_id", testpaperID).Find(&topics).Error
	return
}

// GetTopicBySort 透過 sort 取得 topic
func GetTopicBySort(testpaperID uint, sort uint) (Topic, error) {
	var topic Topic
	if err := DB.Where("testpaper_id = ?", testpaperID).Where("sort = ?", sort).First(&topic).Error; err != nil {
		return Topic{}, err
	}
	return topic, nil
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
