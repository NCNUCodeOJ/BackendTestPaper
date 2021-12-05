package router

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"time"

	"github.com/NCNUCodeOJ/BackendTestPaper/views"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func getUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(jwt.ExtractClaims(c)["id"].(string))
		if err != nil {
			c.Abort()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "系統錯誤",
				"error":   err.Error(),
			})
		} else {
			c.Set("userID", uint(id))
			c.Set("teacher", jwt.ExtractClaims(c)["teacher"].(bool))
			c.Set("admin", jwt.ExtractClaims(c)["admin"].(bool))
			c.Next()
		}
	}
}

// SetupRouter index
func SetupRouter() *gin.Engine {
	if gin.Mode() == "test" {
		err := godotenv.Load(".env.test")
		if err != nil {
			log.Println("Error loading .env.test file")
		}
	} else if gin.Mode() == "debug" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "NCNUOJ",
		SigningAlgorithm: "HS512",
		Key:              []byte(os.Getenv("SECRET_KEY")),
		MaxRefresh:       time.Hour,
		TimeFunc:         time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	baseURL := "api/v1"
	privateURL := "api/private/v1"
	r := gin.Default()
	// testpaper 測驗卷
	// Group 可以讓網址延伸，不用重複寫
	testpaper := r.Group(privateURL + "/testpaper")
	testpaper.Use(authMiddleware.MiddlewareFunc())
	testpaper.Use(getUserInfo())
	{
		testpaper.POST("", views.CreateTestPaper)
		testpaper.GET("", views.ListTestPapers)
		testpaper.GET("/:testpaperID", views.GetTestPaperByID)
		testpaper.PATCH("/:testpaperID", views.UpdateTestPaper)
		testpaper.DELETE("/:testpaperID", views.DeleteTestPaper)
	}
	// topic 大題
	topic := r.Group(privateURL + "/testpaper/:testpaperID/topic")
	{
		topic.POST("", views.CreateTopic)
		topic.GET("", views.ListTopics)
		topic.GET("/:sort", views.GetTopicBySort)
		topic.PATCH("/:sort", views.UpdateTopic)
		topic.DELETE("/:sort", views.DeleteTopic)
	}
	// Question 題目 (含選項/答案)
	question := r.Group(baseURL + "/question")
	question.Use(authMiddleware.MiddlewareFunc())
	question.Use(getUserInfo())
	{
		question.POST("", views.CreateQuestion)
		question.GET("", views.ListQuestions)
		question.GET("/:questionID", views.GetQuestion)
		// question.PATCH("/:questionID", views.UpdateQuestion)
		// question.DELETE("/:questionID", views.DeleteQuestion)
		// 對使用者而言，一個問題就是一個物件
	}
	privatequestion := r.Group(privateURL + "/question")
	question.Use(authMiddleware.MiddlewareFunc())
	question.Use(getUserInfo())
	{
		privatequestion.GET("/:questionID", views.GetQuestionPrivate)
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Page not found"})
	})
	return r
}
