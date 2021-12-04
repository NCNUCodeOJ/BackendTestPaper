package models

import (
	"gorm.io/gorm"
)

//TestPaper Database - database
type TestPaper struct {
	gorm.Model
	TestPaperName string `gorm:"type:text;"`
	Author        uint   `gorm:"NOT NULL;"`
	ClassID       uint   `gorm:"NOT NULL;"`
	Random        bool   `gorm:"NOT NULL;"`
	// 測驗卷名稱
	// 出卷者
	// 對應的課堂
	// 是否隨機出題
}

// CreateTestPaper 新增測驗卷
func CreateTestPaper(testpaper *TestPaper) {
	DB.Create(&testpaper)
}

// ListTestPapers 取得所有 testpaper
func ListTestPapers() (testpapers []TestPaper, err error) {
	err = DB.Find(&testpapers).Error
	return
}

// GetTestPaperByID 透過 id 取得
func GetTestPaperByID(testpaperID uint) (testpaper TestPaper, err error) {
	err = DB.First(&testpaper, testpaperID).Error
	return
}

// UpdateTestPaper 更新
func UpdateTestPaper(testpaper *TestPaper) (err error) {
	err = DB.Where("id = ?", testpaper.ID).Save(&testpaper).Error
	return
}

// DeleteTestPaper 刪除
func DeleteTestPaper(testpaper TestPaper) (err error) {
	err = DB.Where("id = ?", testpaper.ID).Delete(&testpaper).Error
	return
}
