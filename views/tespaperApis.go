package views

import (
	"math/bits"
	"strconv"

	"github.com/NCNUCodeOJ/BackendTestPaper/models"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vincentinttsh/replace"
	"github.com/vincentinttsh/zero"
)

// CreateTestPaper 新增測驗卷
func CreateTestPaper(c *gin.Context) {
	var testpaper models.TestPaper
	var err error
	userID := c.MustGet("userID").(uint)

	// 使用者傳過來的檔案格式(測驗卷名稱、出卷者、對應的課堂、是否亂數出題)
	var data struct {
		TestPaperName *string `json:"testpaper_name"`
		ClassID       *uint   `json:"class_id"`
	}

	if err = c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please fill the field according to the form.",
		})
		return
	}

	// 如果有空值，則回傳 false
	if zero.IsZero(data) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The field cannot be empty.",
		})
		return
	}

	testpaper.TestPaperName = *data.TestPaperName
	testpaper.Author = userID
	testpaper.ClassID = *data.ClassID

	if err = models.CreateTestPaper(&testpaper); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Create successfully.",
		"testpaper_id": testpaper.ID,
	})
}

// GetTestPaperByID 透過 id 取得測驗卷
func GetTestPaperByID(c *gin.Context) {
	// check redis
	// ParseUint convert strings to values
	id, err := strconv.Atoi(c.Params.ByName("testpaperID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	testpaper, err := models.GetTestPaperByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":             testpaper.ID,
		"testpaper_name": testpaper.TestPaperName,
		"author":         testpaper.Author,
		"class_id":       testpaper.ClassID,
	})
}

// UpdateTestPaper 更新測驗卷
func UpdateTestPaper(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("testpaperID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	var data struct {
		TestPaperName *string `json:"testpaper_name"`
		Author        *uint   `json:"author"`
		ClassID       *uint   `json:"class_id"`
		Random        *bool   `json:"random"`
	}
	c.BindJSON(&data)
	testpaper, err := models.GetTestPaperByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
		return
	}
	replace.Replace(&testpaper, &data)
	err = models.UpdateTestPaper(&testpaper)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The field cannot be empty.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Update successfully.",
	})
}

// DeleteTestPaper 刪除測驗卷
func DeleteTestPaper(c *gin.Context) {
	id, err := strconv.ParseUint(c.Params.ByName("testpaperID"), 10, bits.UintSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System error.",
		})
		return
	}
	testpaper, err := models.GetTestPaperByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
		return
	}
	err = models.DeleteTestPaper(testpaper)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Fail.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete successfully.",
	})
}
