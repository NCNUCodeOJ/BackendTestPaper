package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/NCNUCodeOJ/BackendTestPaper/models"
	"github.com/NCNUCodeOJ/BackendTestPaper/router"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

// cspell:disable-next-line
var token = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjo0NzYwNjk2NDkyLCJpZCI6IjEiLCJvcmlnX2lhdCI6MTYzODYzMjQ5MiwidGVhY2hlciI6dHJ1ZSwidXNlcm5hbWUiOiJ2aW5jZW50In0.SUnwDQX_wkWlZdTMyCjhqIX4TIIzYrrY7lTiR_E2K8tvQBU1pyUgja60K0xcF1_x0m-egvRJQmhix5l6wdoR6g"

var questionID int

func init() {
	gin.SetMode(gin.TestMode)
	models.Setup()
}

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

func TestGetQuestion(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/question/"+strconv.Itoa(questionID), bytes.NewBuffer(make([]byte, 1000)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateQuestion(t *testing.T) {
	var questionData = []byte(`{
		"question": "question1patchtest",
		"author": 1,
		"layer": 1,
		"source": "1",
		"difficulty": 1,
		"type": 1
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("PATCH", "/api/v1/question/"+strconv.Itoa(questionID), bytes.NewBuffer(questionData))
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)
	fmt.Println(string(body))
	s := struct {
		Message string `json:"message"`
	}{}
	json.Unmarshal(body, &s)
}

func TestDeleteQuestion(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("DELETE", "/api/v1/question/"+strconv.Itoa(questionID), bytes.NewBuffer(make([]byte, 1000)))
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		Message string `json:"message"`
	}{}
	json.Unmarshal(body, &s)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCleanup(t *testing.T) {
	e := os.Remove("test.db")
	if e != nil {
		t.Fail()
	}
}
