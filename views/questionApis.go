package views

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NCNUCodeOJ/BackendTestPaper/models"
	"github.com/gin-gonic/gin"
	"github.com/vincentinttsh/zero"
)

// CreateQuestion 新增題目
func CreateQuestion(c *gin.Context) {
	// 使用者傳過來的檔案格式 (題目、出題者、範圍、出處)
	var question models.Question
	userID := c.MustGet("userID").(uint)
	teacher := c.MustGet("teacher").(bool)
	if !teacher {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denial.",
		})
		return
	}
	var questionData struct {
		Question   *string          `json:"question"`
		Author     *uint            `json:"author"`
		Layer      *uint            `json:"layer"`
		Source     *string          `json:"source"`
		Difficulty *uint            `json:"difficulty"`
		Type       *uint            `json:"type"`
		OptionList []*models.Option `json:"options"`
	}
	if err := c.BindJSON(&questionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please fill the field according to the form.",
		})
		return
	}
	// 如果有空值，則回傳 false
	if zero.IsZero(&questionData) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The field cannot be empty.",
		})
		return
	}
	question.Question = *questionData.Question
	question.Author = userID
	question.Layer = *questionData.Layer
	question.Source = *questionData.Source
	question.Difficulty = *questionData.Difficulty
	question.Type = *questionData.Type
	models.CreateQuestion(&question)

	for i, optionData := range questionData.OptionList {
		var option models.Option
		option.Content = optionData.Content
		option.Answer = optionData.Answer
		option.QuestionID = question.ID
		option.Sort = uint(i + 1)
		if err := models.CreateOption(&option); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Fail.",
			})
			fmt.Println(option)
			return
		}
		fmt.Println(option)
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "Create successfully.",
		"question_id": question.ID,
	})
}

// ListQuestions 取得全部題目的 id
func ListQuestions(c *gin.Context) {
	var questionsID []uint
	teacher := c.MustGet("teacher").(bool)
	if !teacher {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denial.",
		})
		return
	}
	if questions, err := models.ListQuestions(); err == nil {
		for pos := range questions {
			questionsID = append(questionsID, questions[pos].ID)
		}
		c.JSON(http.StatusOK, gin.H{
			"questions_id": questionsID,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
	}
}

// GetQuestion 查詢題目 (也會列出該題目所有選項/答案)
func GetQuestionPrivate(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("questionID"))
	var optionList []models.Option
	var option = make([]gin.H, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	question, err := models.GetQuestion(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
		return
	}
	for _, opt := range optionList {
		option = append(option, gin.H{
			"option_content":     opt.Content,
			"option_answer":      opt.Answer,
			"option_question_id": opt.QuestionID,
			"option_sort":        opt.Sort,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         question.ID,
		"question":   question.Question,
		"author":     question.Author,
		"layer":      question.Layer,
		"source":     question.Source,
		"difficulty": question.Difficulty,
		"type":       question.Type,
		"option":     optionList,
	})
}

// GetQuestion (老師)查詢題目 (也會列出該題目所有選項/答案)
func GetQuestion(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("questionID"))
	var optionList []models.Option
	var option = make([]gin.H, 0)
	teacher := c.MustGet("teacher").(bool)
	if !teacher {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denial.",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	question, err := models.GetQuestion(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
		return
	}
	for _, opt := range optionList {
		option = append(option, gin.H{
			"option_content":     opt.Content,
			"option_answer":      opt.Answer,
			"option_question_id": opt.QuestionID,
			"option_sort":        opt.Sort,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         question.ID,
		"question":   question.Question,
		"author":     question.Author,
		"layer":      question.Layer,
		"source":     question.Source,
		"difficulty": question.Difficulty,
		"type":       question.Type,
		"option":     optionList,
	})
}
