package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"

	"github.com/NCNUCodeOJ/BackendTestPaper/models"
	"github.com/NCNUCodeOJ/BackendTestPaper/router"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

// cspell:disable-next-line
var token = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjo0NzYwNjk2NDkyLCJpZCI6IjEiLCJvcmlnX2lhdCI6MTYzODYzMjQ5MiwidGVhY2hlciI6dHJ1ZSwidXNlcm5hbWUiOiJ2aW5jZW50In0.SUnwDQX_wkWlZdTMyCjhqIX4TIIzYrrY7lTiR_E2K8tvQBU1pyUgja60K0xcF1_x0m-egvRJQmhix5l6wdoR6g"
var sigs = make(chan os.Signal, 1)

var questionID int

func init() {
	gin.SetMode(gin.TestMode)
	models.Setup()
}

// POST 測驗卷
func TestCreateTestpaper(t *testing.T) {
	var testpaperData = []byte(`{
		"testpaper_name": "testpaper1",
		"author": 1,
		"class_id": 1,
		"random": false
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/private/v1/testpaper", bytes.NewBuffer(testpaperData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// POST 大題
func TestCreateTopic(t *testing.T) {
	var topicData = []byte(`{
		"distribution": 1,
		"testpaper_id": 1,
		"sort": 1,
		"question_id": 1,
		"type": 1
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/private/v1/testpaper/1/topic", bytes.NewBuffer(topicData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// POST 題目
func TestCreateQuestion(t *testing.T) {
	var questionData = []byte(`{
		"question": "question1",
		"author": 1,
		"layer": 1,
		"source": "source",
		"difficulty": 1,
		"type": 1,
		"option": [
			{
				"content": "content1",
				"answer": true,
				"question_id": 1,
				"sort": 1
			},
			{
				"content": "content2",
				"answer": true,
				"question_id": 1,
				"sort": 2
			}
		  ]
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/v1/question", bytes.NewBuffer(questionData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		QuestionID int `json:"question_id"`
	}{}
	json.Unmarshal(body, &s)
	questionID = s.QuestionID
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListQuestions(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/question", bytes.NewBuffer(make([]byte, 1000)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// GET 所有測驗卷
func TestListTestPapers(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/private/v1/testpaper", bytes.NewBuffer(make([]byte, 1000)))
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		TestpapersID []uint `json:"testpapers_id"`
	}{}
	json.Unmarshal(body, &s)
	r.ServeHTTP(w, req)
}

// GET 測驗卷
func TestGetTestpaperByID(t *testing.T) {
	r := router.SetupRouter()
	// 取得 ResponseRecorder 物件，用來記錄 response 狀態
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/private/v1/testpaper/1", bytes.NewBuffer(make([]byte, 1000)))
	req.Header.Set("Authorization", token)
	// gin.Engine.ServerHttp 實作 http.Handler 介面，用來處理 HTTP 請求及回應
	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		TestPaperName string `json:"testpaper_name"`
		Author        uint   `json:"author"`
		ClassID       uint   `json:"class_id"`
		Random        bool   `json:"random"`
	}{}
	json.Unmarshal(body, &s)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	assert.Equal(t, http.StatusOK, w.Code)
}

// GET 所有大題
func TestListTopics(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/private/v1/testpaper/1/topic", bytes.NewBuffer(make([]byte, 1000)))
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		TopicID []uint `json:"topics_id"`
	}{}
	json.Unmarshal(body, &s)
	r.ServeHTTP(w, req)
}

// Get 大題
func TestGetTopic(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/private/v1/testpaper/1/topic/1", bytes.NewBuffer(make([]byte, 1000)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// // GET 所有題目
// func TestListQuestions(t *testing.T) {
// 	r := router.SetupRouter()
// 	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
// 	req, _ := http.NewRequest("GET", "/api/v1/question", bytes.NewBuffer(make([]byte, 1000)))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", token)
// 	r.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)
// }

// GET 題目
func TestGetQuestion(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/question/"+strconv.Itoa(questionID), bytes.NewBuffer(make([]byte, 1000)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// PATCH 測驗卷
func TestUpdateTestpaper(t *testing.T) {
	var testpaperPatchData = []byte(`{
		"testpaper_name": "testpaper1patchtest",
		"author": 1,
		"class_id": 1,
		"random": false
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("PATCH", "/api/private/v1/testpaper/1", bytes.NewBuffer(testpaperPatchData))
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		Message string `json:"message"`
	}{}
	json.Unmarshal(body, &s)
	assert.Equal(t, http.StatusOK, w.Code)
}

// PATCH 大題
func TestUpdateTopic(t *testing.T) {
	var topicData = []byte(`{
		"distribution": 2,
		"testpaper_id": 1,
		"sort": 1
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("PATCH", "/api/private/v1/testpaper/1/topic/1", bytes.NewBuffer(topicData))
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		Message string `json:"message"`
	}{}
	json.Unmarshal(body, &s)
}
func TestCleanup(t *testing.T) {
	e := os.Remove("test.db")
	if e != nil {
		t.Fail()
	}
}

// POST 該測驗卷下已評分學生考卷
func TestCreateStudentTestpaper(t *testing.T) {
	var StudentTestpaperData = []byte(`{
		"student_id": 1,
		"testpaper_id": 1,
		"score": 60
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/private/v1/testpaper/1/graded", bytes.NewBuffer(StudentTestpaperData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// GET 該測驗卷下所有已評分的學生考卷
func TestListStudentTestPaper(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/private/v1/testpaper/1/graded", bytes.NewBuffer(make([]byte, 1000)))
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		StudentTestPaper []uint `json:"student_testpaper_id"`
	}{}
	json.Unmarshal(body, &s)
	r.ServeHTTP(w, req)
}

// Get 該測驗卷下該學生已作答考卷
func TestGetStudentTestPaperByStudentID(t *testing.T) {
	r := router.SetupRouter()
	// 取得 ResponseRecorder 物件，用來記錄 response 狀態
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/private/v1/testpaper/1/graded/1", bytes.NewBuffer(make([]byte, 1000)))
	req.Header.Set("Authorization", token)
	// gin.Engine.ServerHttp 實作 http.Handler 介面，用來處理 HTTP 請求及回應
	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		TestpaperID uint    `json:"testpaper_id"`
		StudentID   uint    `json:"student_id"`
		Score       float64 `json:"score"`
	}{}
	json.Unmarshal(body, &s)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	assert.Equal(t, http.StatusOK, w.Code)
}

// POST 該測驗卷下已批改的學生題目
func TestCreateStudentTestpaperAnswer(t *testing.T) {
	var StudentTestpaperAnswerData = []byte(`{
		"student_testpaper_id": 1,
		"topic_sort": 1,
		"question_sort": 1,
		"resort": 1,
		"content": "content",
		"correct": true
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/private/v1/testpaper/1/answer/1", bytes.NewBuffer(StudentTestpaperAnswerData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// GET 該測驗卷下該學生所有已被批改的題目
func TestListStudentTestPaperAnswer(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/private/v1/testpaper/1/answer/1", bytes.NewBuffer(make([]byte, 1000)))
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		StudentTestPaperAnswer []uint `json:"student_testpaper_answer_id"`
	}{}
	json.Unmarshal(body, &s)
	r.ServeHTTP(w, req)
}

// Delete 測驗卷
func TestDeleteTestpaper(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("DELETE", "/api/private/v1/testpaper/1", bytes.NewBuffer(make([]byte, 1000)))
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		Message string `json:"message"`
	}{}
	json.Unmarshal(body, &s)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Delete 大題
func TestDeleteTopic(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("DELETE", "/api/private/v1/testpaper/1/topic/1", bytes.NewBuffer(make([]byte, 1000)))
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		Message string `json:"message"`
	}{}
	json.Unmarshal(body, &s)
	assert.Equal(t, http.StatusOK, w.Code)
}
