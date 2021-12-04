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
	userID := c.MustGet("userID").(uint)
	// 使用者傳過來的檔案格式(測驗卷名稱、出卷者、對應的課堂、是否亂數出題)
	var data struct {
		TestPaperName *string `json:"testpaper_name"`
		Author        *uint   `json:"author"`
		ClassID       *uint   `json:"class_id"`
		Random        *bool   `json:"random"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫或未使用 json",
		})
		return
	}
	// 如果有空值，則回傳 false
	if zero.IsZero(data) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未填寫完成",
		})
		return
	}
	testpaper.TestPaperName = *data.TestPaperName
	testpaper.Author = userID
	testpaper.ClassID = *data.ClassID
	testpaper.Random = *data.Random
	models.CreateTestPaper(&testpaper)
	c.JSON(http.StatusOK, gin.H{
		"message": "新增成功",
	})
}

// ListTestPapers 取得全部測驗卷的 id
func ListTestPapers(c *gin.Context) {
	var testpapersID []uint
	if testpapers, err := models.ListTestPapers(); err == nil {
		for pos := range testpapers {
			testpapersID = append(testpapersID, testpapers[pos].ID)
		}
		c.JSON(http.StatusOK, gin.H{
			"testpapersID": testpapersID,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "查無資料",
		})
	}
}

// GetTestPaperByID 透過 id 取得測驗卷
func GetTestPaperByID(c *gin.Context) {
	// check redis
	// ParseUint convert strings to values
	id, err := strconv.Atoi(c.Params.ByName("testpaperID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	testpaper, err := models.GetTestPaperByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "查無資料",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":             testpaper.ID,
		"testpaper_name": testpaper.TestPaperName,
		"author":         testpaper.Author,
		"class_id":       testpaper.ClassID,
		"random":         testpaper.Random,
	})
}

// UpdateTestPaper 更新測驗卷
func UpdateTestPaper(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("testpaperID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
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
			"message": "查無資料",
		})
		return
	}
	replace.Replace(&testpaper, &data)
	err = models.UpdateTestPaper(&testpaper)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未填寫完成",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
	})
}

// DeleteTestPaper 刪除測驗卷
func DeleteTestPaper(c *gin.Context) {
	id, err := strconv.ParseUint(c.Params.ByName("testpaperID"), 10, bits.UintSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	testpaper, err := models.GetTestPaperByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "查無資料",
		})
		return
	}
	err = models.DeleteTestPaper(testpaper)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "刪除失敗",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "刪除成功",
	})
}
