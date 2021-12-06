package views

import (
	"strconv"
	"strings"

	"github.com/NCNUCodeOJ/BackendTestPaper/models"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vincentinttsh/zero"
)

// CreateStudentTestPaperAnswer 新增已批改的題目
func CreateStudentTestPaperAnswer(c *gin.Context) {
	var studentTestpaperAnswer models.StudentTestPaperAnswer
	var studentData struct {
		StudentTestpaperID *uint   `json:"student_testpaper_id"`
		TopicSort          *uint   `json:"topic_sort"`
		QuestionSort       *uint   `json:"question_sort"`
		Resort             *uint   `json:"resort"`
		Content            *string `json:"content"`
		Correct            *bool   `json:"correct"`
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
	// 取得 testpaper_id: studentTestpaperData.TestPaperID
	studentTestpaperData, err := models.GetStudentTestPaper(uint(*studentData.StudentTestpaperID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
		return
	}
	// TopicSort
	topicData, err := models.GetTopicBySort(uint(studentTestpaperData.TestPaperID), uint(*studentData.TopicSort))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Topic not found.",
		})
		return
	}
	// QuestionSort
	questionTopicData, err := models.GetQuestionTopicByTopicIDandSort(uint(topicData.ID), uint(*studentData.QuestionSort))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "questionTopic not found.",
		})
		return
	}
	// 取得 Answer(正確答案) optionData.Content
	optionData, err := models.GetAnswerByQuestionID(questionTopicData.QuestionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Option not found.",
		})
		return
	}
	// 學生答對
	studentTestpaperAnswer.StudentTestPaperID = *studentData.StudentTestpaperID
	studentTestpaperAnswer.TopicSort = topicData.Sort
	studentTestpaperAnswer.QuestionSort = *studentData.QuestionSort
	studentTestpaperAnswer.Resort = *studentData.Resort
	studentTestpaperAnswer.Content = *studentData.Content
	if strings.Contains(studentTestpaperAnswer.Content, optionData.Content) {
		studentTestpaperAnswer.Correct = false
	} else {
		studentTestpaperAnswer.Correct = true
	}
	models.CreateStudentTestPaperAnswer(&studentTestpaperAnswer)
	c.JSON(http.StatusOK, gin.H{
		"message": "Create successfully.",
	})
}

// ListStudentTestPaperAnswersByStudentTestPaperID 透過 studentTestPaperID 取得測驗卷中該學生所有已被批改的題目
func ListStudentTestPaperAnswersByStudentTestPaperID(c *gin.Context) {
	testpaperID, err := strconv.Atoi(c.Params.ByName("testpaperID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	var studentTestpaperAnswersID []uint
	if studentTestpaperAnswers, err := models.ListStudentTestPaperAnswersByStudentTestPaperID(uint(testpaperID)); err == nil {
		for pos := range studentTestpaperAnswers {
			studentTestpaperAnswersID = append(studentTestpaperAnswersID, studentTestpaperAnswers[pos].ID)
		}
		c.JSON(http.StatusOK, gin.H{
			"studentTestpaperAnswersID": studentTestpaperAnswersID,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
	}
}

// ListStudentTestPaperAnswersByStudentID 透過 studentID 取得該生的已批改的題目
func ListStudentTestPaperAnswersByStudentID(c *gin.Context) {
	testpaperID, err := strconv.Atoi(c.Params.ByName("testpaperID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	studentID, err := strconv.Atoi(c.Params.ByName("student_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	studentTestpaperData, err := models.GetStudentTestPaperByTestPaperIDandStudentID(uint(testpaperID), uint(studentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
		return
	}
	var allstudentTestpaperAnswerID []uint
	if studentTestpaperAnswerData, err := models.ListStudentTestPaperAnswersByStudentTestPaperID(uint(studentTestpaperData.ID)); err == nil {
		for pos := range studentTestpaperAnswerData {
			allstudentTestpaperAnswerID = append(allstudentTestpaperAnswerID, studentTestpaperAnswerData[pos].ID)
		}
		c.JSON(http.StatusOK, gin.H{
			"studentTestpaperAnswersID": allstudentTestpaperAnswerID,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
	}
}
