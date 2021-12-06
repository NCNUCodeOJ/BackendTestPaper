package views

import (
	"math/bits"
	"strconv"

	"net/http"

	"github.com/NCNUCodeOJ/BackendTestPaper/models"
	"github.com/gin-gonic/gin"
	"github.com/vincentinttsh/replace"
	"github.com/vincentinttsh/zero"
)

// CreateTopic 新增
func CreateTopic(c *gin.Context) {
	// 使用者傳過來的檔案格式(名稱、出卷者、對應的課堂、是否亂數出題)
	var topicData struct {
		TestPaperID  *uint    `json:"testPaper_id"`
		Distribution *float64 `json:"distribution"`
		TopicSort    *uint    `json:"topic_sort"`
		QuestionID   *uint    `json:"question_id"`
		QuestionSort *uint    `json:"question_sort"`
		Type         *uint    `json:"type"`
	}
	var topic models.Topic
	if err := c.BindJSON(&topicData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please fill the fields as JSON format.",
		})
		return
	}
	// 如果有空值，則回傳 false
	if zero.IsZero(topicData) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The field cannot be empty.",
		})
		return
	}
	topic.TestPaperID = *topicData.TestPaperID
	topic.Distribution = *topicData.Distribution
	topic.Sort = *topicData.TopicSort
	models.CreateTopic(&topic)
	var questionTopic models.QuestionTopic
	questionTopic.TopicID = topic.ID
	questionTopic.QuestionID = *topicData.QuestionID
	questionTopic.Sort = *topicData.QuestionSort
	questionTopic.Type = *topicData.Type
	models.CreateQuestionTopic(&questionTopic)
	c.JSON(http.StatusOK, gin.H{
		"message": "Create successfully.",
	})
}

// ListTopics 透過 testpaper_id 取得測驗卷
func ListTopics(c *gin.Context) {
	testpaperID, err := strconv.Atoi(c.Params.ByName("tetspaperID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	var allTopicID []uint
	if topics, err := models.ListTopicsByTestpaperID(uint(testpaperID)); err == nil {
		for pos := range topics {
			allTopicID = append(allTopicID, topics[pos].ID)
		}
		c.JSON(http.StatusOK, gin.H{
			"topicsID": allTopicID,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
	}
}

// GetTopicBySort 透過 sort 取得 topic
func GetTopicBySort(c *gin.Context) {
	testpaperID, err := strconv.Atoi(c.Params.ByName("testpaperID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	sort, err := strconv.Atoi(c.Params.ByName("sort"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	topic, err := models.GetTopicBySort(uint(testpaperID), uint(sort))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":           topic.ID,
		"distribution": topic.Distribution,
		"testpaper_id": topic.TestPaperID,
		"sort":         topic.Sort,
	})
}

// UpdateTopic 更新大題
func UpdateTopic(c *gin.Context) {
	testpaperID, err := strconv.Atoi(c.Params.ByName("testpaperID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	sort, err := strconv.Atoi(c.Params.ByName("sort"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	data := struct {
		Distribution *float64 `json:"distribution"`
		TopicSort    *uint    `json:"topic_sort"`
		QuestionID   *uint    `json:"question_id"`
		QuestionSort *uint    `json:"question_sort"`
		Type         *uint    `json:"type"`
	}{}
	c.BindJSON(&data)
	topic, err := models.GetTopicBySort(uint(testpaperID), uint(sort))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
		return
	}
	var newTopic models.Topic
	newTopic.TestPaperID = topic.TestPaperID
	newTopic.Distribution = *data.Distribution
	newTopic.Sort = *data.TopicSort
	replace.Replace(&topic, &newTopic)
	err = models.UpdateTopic(&topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail.",
		})
		return
	}
	questionTopic, err := models.GetQuestionTopic(uint(topic.ID), uint(*data.QuestionID))
	if err != nil {
		var questionTopic models.QuestionTopic
		questionTopic.TopicID = topic.ID
		questionTopic.QuestionID = *data.QuestionID
		questionTopic.Sort = *data.QuestionSort
		questionTopic.Type = *data.Type
		models.CreateQuestionTopic(&questionTopic)
	} else {
		models.DeleteQuestionTopic(questionTopic)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Update successfully.",
	})
}

// DeleteTopic 刪除
func DeleteTopic(c *gin.Context) {
	testpaperID, err := strconv.Atoi(c.Params.ByName("testpaperID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	sort, err := strconv.ParseUint(c.Params.ByName("sort"), 10, bits.UintSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	topic, err := models.GetTopicBySort(uint(testpaperID), uint(sort))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
		return
	}
	err = models.DeleteTopic(topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail.",
		})
		return
	}
	if questionTopics, err := models.ListQuestionTopicsByTopicID(uint(topic.ID)); err == nil {
		for pos := range questionTopics {
			var questionTopicID = questionTopics[pos].ID
			if questionTopic, err := models.GetQuestionTopicByID(uint(questionTopicID)); err == nil {
				models.DeleteQuestionTopic(questionTopic)
			}
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete successfully.",
	})
}
