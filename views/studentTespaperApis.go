package views

import (
	"strconv"

	"github.com/NCNUCodeOJ/BackendTestPaper/models"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vincentinttsh/zero"
)

// CreateStudentTestPaper 新增已評分測驗卷
func CreateStudentTestPaper(c *gin.Context) {
	var studentTestpaper models.StudentTestPaper
	userID := c.MustGet("userID").(uint)
	var studentData struct {
		StudentID   *uint    `json:"student_id"`
		TestpaperID *uint    `json:"testpaper_id"`
		Score       *float64 `json:"score"`
	}
	if err := c.BindJSON(&studentData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please fill the fields as JSON format.",
		})
		return
	}
	// 如果有空值，則回傳 false
	if zero.IsZero(studentData) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The field cannot be empty.",
		})
		return
	}
	// 算總分 score
	// 取得 distribution: questionTopicData.distribution
	// studentTestPaperAnswer []*uint
	// topicData, err := models.GetTopicBySort(uint(*studentData.TestpaperID), uint())
	// if err != nil {
	// 	fmt.Println("topicID not found.")
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"message": "Not found.",
	// 	})
	// 	return
	// }
	studentTestpaper.StudentID = userID
	studentTestpaper.TestPaperID = *studentData.TestpaperID
	// studentTestpaper.Score = score
	models.CreateStudentTestPaper(&studentTestpaper)
	c.JSON(http.StatusOK, gin.H{
		"message": "Create successfully.",
	})
}

// ListStudentTestPapers 取得這個測驗卷的所有已評分的學生考卷
func ListStudentTestPapers(c *gin.Context) {
	var studentTestPapersID []uint
	if studentTestPapers, err := models.ListStudentTestPapers(); err == nil {
		for pos := range studentTestPapers {
			studentTestPapersID = append(studentTestPapersID, studentTestPapers[pos].ID)
		}
		c.JSON(http.StatusOK, gin.H{
			"studentTestPapersID": studentTestPapersID,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
	}
}

// GetStudentTestPaperTestPaperIDandStudentID 透過 StudentID 取得該生的已評分測驗卷
func GetStudentTestPaperTestPaperIDandStudentID(c *gin.Context) {
	studentID, err := strconv.Atoi(c.Params.ByName("studentID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	testpaperID, err := strconv.Atoi(c.Params.ByName("testpaperID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	studentTestpaper, err := models.GetStudentTestPaperByTestPaperIDandStudentID(uint(studentID), uint(testpaperID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         studentTestpaper.ID,
		"student_id": studentTestpaper.StudentID,
		"score":      studentTestpaper.Score,
	})
}
