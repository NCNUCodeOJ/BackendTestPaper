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
	var data struct {
		Description *string `json:"description"`
		TestPaperID *uint   `json:"testPaper_id"`
		Sort        *uint   `json:"sort"`
	}
	var topic models.Topic
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	// 如果有空值，則回傳 false
	if zero.IsZero(data) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "所有欄位不可為空值",
		})
		return
	}

	topic.Description = *data.Description
	topic.TestPaperID = *data.TestPaperID
	topic.Sort = *data.Sort
	models.CreateTopic(&topic)
	c.JSON(http.StatusOK, gin.H{
		"message": "新增成功",
	})
}

// ListTopics 透過 testpaper_id 取得測驗卷
func ListTopics(c *gin.Context) {
	var allTopicID []uint
	if topics, err := models.ListTopics(); err == nil {
		for pos := range topics {
			allTopicID = append(allTopicID, topics[pos].ID)
		}
		c.JSON(http.StatusOK, gin.H{
			"topicsID": allTopicID,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "尚無內容",
		})
	}
}

// GetTopicBySort 透過 sort 取得 topic
func GetTopicBySort(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("testpaperID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	sort, err := strconv.Atoi(c.Params.ByName("sort"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	topic, err := models.GetTopicBySort(uint(id), uint(sort))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "查無資料",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":           topic.ID,
		"description":  topic.Description,
		"testpaper_id": topic.TestPaperID,
		"sort":         topic.Sort,
	})
}

// UpdateTopic 更新大題
func UpdateTopic(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("testpaperID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	sort, err := strconv.Atoi(c.Params.ByName("sort"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	data := struct {
		Description *string `json:"description"`
		TestPaperID *uint   `json:"testPaper_id"`
		Sort        *uint   `json:"sort"`
	}{}
	c.BindJSON(&data)
	topic, err := models.GetTopicBySort(uint(id), uint(sort))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "查無資料",
		})
		return
	}
	replace.Replace(&topic, &data)
	err = models.UpdateTopic(&topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新失敗",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
	})
}

// DeleteTopic 刪除
func DeleteTopic(c *gin.Context) {
	id, err := strconv.ParseUint(c.Params.ByName("testpaperID"), 10, bits.UintSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	sort, err := strconv.ParseUint(c.Params.ByName("sort"), 10, bits.UintSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	topic, err := models.GetTopicBySort(uint(id), uint(sort))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "查無資料",
		})
		return
	}
	err = models.DeleteTopic(topic)
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
