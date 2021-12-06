package models

import (
	"gorm.io/gorm"
)

type StudentTestPaper struct {
	gorm.Model
	StudentID   uint    `gorm:"NOT NULL;"`
	TestPaperID uint    `gorm:"NOT NULL;"`
	Score       float64 `gorm:"NOT NULL;"`
	// 答題者
	// 對應的課堂
	// 學生
}

// CreateStudentTestPaper 新增該測驗卷下已評分學生考卷
func CreateStudentTestPaper(studentTestpaper *StudentTestPaper) {
	DB.Create(&studentTestpaper)
}

// ListStudentTestPapers 取得該測驗卷下所有學生的已評分測驗卷
func ListStudentTestPapers() (studentTestpapers []StudentTestPaper, err error) {
	err = DB.Find(&studentTestpapers).Error
	return
}

// GetStudentTestPaperByStudentByID 透過 studentID 取得該學生的已評分測驗卷
func GetStudentTestPaperByStudentByID(studentID uint) (studentTestpaper StudentTestPaper, err error) {
	err = DB.Where("student_id = ?", studentID).First(&studentTestpaper).Error
	return
}

// GetStudentTestPaper 透過 id 取得已評分學生考卷
func GetStudentTestPaper(id uint) (studentTestpaper StudentTestPaper, err error) {
	err = DB.Where("id = ?", id).First(&studentTestpaper).Error
	return
}

// GetStudentTestPaperByTestPaperIDandStudentID 透過 TestPaperID 和 StudentID 取得已評分學生考卷
func GetStudentTestPaperByTestPaperIDandStudentID(testpaperID uint, studentID uint) (studentTestpaper StudentTestPaper, err error) {
	err = DB.Where("testpaper_id = ?", testpaperID).Where("student_id = ?", studentID).First(&studentTestpaper).Error
	return
}
