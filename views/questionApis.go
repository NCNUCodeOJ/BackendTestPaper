package views

import (
	"fmt"
	"math/bits"
	"net/http"
	"strconv"

	"github.com/NCNUCodeOJ/BackendTestPaper/models"
	"github.com/gin-gonic/gin"
	"github.com/vincentinttsh/replace"
	"github.com/vincentinttsh/zero"
)

// CreateQuestion 新增題目
func CreateQuestion(c *gin.Context) {
	// 使用者傳過來的檔案格式 (題目、出題者、範圍、出處)
	var question models.Question
	userID := c.MustGet("userID").(uint)
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
			"message": "未按照格式填寫",
		})
		return
	}
	// 如果有空值，則回傳 false
	if zero.IsZero(&questionData) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "所有欄位不可為空值",
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
				"message": "題目新增失敗",
			})
			fmt.Println(option)
			return
		}
		fmt.Println(option)
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "新增成功",
		"question_id": question.ID,
	})
}

// ListQuestions 取得全部題目的 id
func ListQuestions(c *gin.Context) {
	var questionsID []uint
	if questions, err := models.ListQuestions(); err == nil {
		for pos := range questions {
			questionsID = append(questionsID, questions[pos].ID)
		}
		c.JSON(http.StatusOK, gin.H{
			"questions_id": questionsID,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "查無資料",
		})
	}
}

// GetQuestion 查詢題目 (也會列出該題目所有選項/答案)
func GetQuestion(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("questionID"))
	var optionList []models.Option
	var option = make([]gin.H, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	question, err := models.GetQuestion(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "查無資料",
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

// UpdateQuestion 更新題目 (包含選項/答案)
func UpdateQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Params.ByName("questionID"), 10, bits.UintSize)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
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
			"message": "json error",
			"error":   err.Error(),
		})
		return
	}

	question, err := models.GetQuestion(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "question not found",
		})
		return
	}

	if err := models.UpdateOption(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
			"error":   err.Error(),
		})
		return
	}
	replace.Replace(&question, &questionData)
	err = models.UpdateQuestion(&question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "update error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "update success",
	})
}

// DeleteQuestion 刪除題目 (包含選項/答案)
func DeleteQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Params.ByName("questionID"), 10, bits.UintSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	question, err := models.GetQuestion(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "查無資料",
		})
		return
	}
	err = models.DeleteQuestion(question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "刪除失敗",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "刪除成功",
	})
}
