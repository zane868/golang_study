package main

import (
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zane868/golang_study/homework04/config"
	"github.com/zane868/golang_study/homework04/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestBlogPageAndAPIWorkflow(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "test.db")), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { sqlDB.Close() })
	if err := db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{}); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{}
	cfg.Jwt.Secret = "integration-test-secret"
	router := regRouter(NewContext(db, cfg), cfg)
	page := httptest.NewRecorder()
	router.ServeHTTP(page, httptest.NewRequest("GET", "/index", nil))
	if page.Code != 200 || !strings.Contains(page.Body.String(), "博客实验室") || !strings.Contains(page.Header().Get("Content-Type"), "text/html") {
		t.Fatal("index page missing")
	}
	token := ""
	call := func(method, path, body string, status int) map[string]json.RawMessage {
		t.Helper()
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", token)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != status {
			t.Fatalf("%s %s: got %d, want %d: %s", method, path, w.Code, status, w.Body.String())
		}
		var result map[string]json.RawMessage
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		return result
	}
	call("GET", "/api/v1/posts", "", 401)
	call("GET", "/api/v1/comments?postId=1", "", 401)
	credentials := `{"username":"alice","email":"alice@example.com","password":"password123"}`
	call("POST", "/api/v1/users/register", credentials, 200)
	call("POST", "/api/v1/users/register", credentials, 409)
	call("POST", "/api/v1/users/login", `{"username":"alice","password":"wrong"}`, 401)
	login := call("POST", "/api/v1/users/login", credentials, 200)
	var session struct {
		Authorization string `json:"authorization"`
	}
	if err := json.Unmarshal(login["data"], &session); err != nil {
		t.Fatal(err)
	}
	token = session.Authorization
	if token == "" {
		t.Fatal("missing token")
	}
	if string(call("GET", "/api/v1/posts", "", 200)["data"]) != "[]" {
		t.Fatal("expected empty array")
	}
	call("POST", "/api/v1/posts", `{"title":"First","content":"你好","username":"spoofed"}`, 200)
	var post model.Post
	if err := db.First(&post, 1).Error; err != nil || post.UserID != 1 {
		t.Fatalf("post identity: %v %+v", err, post)
	}
	call("GET", "/api/v1/posts/1", "", 200)
	for _, query := range []string{"", "?postId=0", "?postId=-1", "?postId=abc"} {
		call("GET", "/api/v1/comments"+query, "", 400)
	}
	call("GET", "/api/v1/comments?postId=999", "", 404)
	if string(call("GET", "/api/v1/comments?postId=1", "", 200)["data"]) != "[]" {
		t.Fatal("expected empty comments")
	}
	call("GET", "/api/v1/posts/0", "", 400)
	call("GET", "/api/v1/posts/999", "", 404)
	call("PUT", "/api/v1/posts/1", `{"title":"Updated","content":"新内容","postId":"999"}`, 200)
	if err := db.First(&post, 1).Error; err != nil || post.Title != "Updated" || post.Count != 3 {
		t.Fatalf("update failed: %v %+v", err, post)
	}
	call("POST", "/api/v1/comments", `{"postId":999,"content":"Missing"}`, 404)
	call("POST", "/api/v1/comments", `{"postId":1,"content":""}`, 422)
	comment := call("POST", "/api/v1/comments", `{"postId":1,"content":"Nice","UserId":999}`, 200)
	var created struct {
		ID     uint `json:"id"`
		UserID uint `json:"user_id"`
	}
	if err := json.Unmarshal(comment["data"], &created); err != nil || created.ID != 1 || created.UserID != 1 {
		t.Fatalf("comment identity: %v %+v", err, created)
	}
	aliceToken := token
	assertComments := func(want int) {
		t.Helper()
		data := call("GET", "/api/v1/comments?postId=1", "", 200)["data"]
		var comments []model.CommentResponse
		if err := json.Unmarshal(data, &comments); err != nil || len(comments) != want {
			t.Fatalf("comment list: %s, %v", data, err)
		}
		if strings.Contains(string(data), "password") || strings.Contains(string(data), "Password") {
			t.Fatal("comment list exposes user model")
		}
		if want > 0 && (comments[0].Username == "" || comments[0].PostID != 1 || len(comments[0].CreatedAt) != 23) {
			t.Fatalf("missing comment metadata: %+v", comments[0])
		}
	}
	assertComments(1)
	other := `{"username":"bob","email":"bob@example.com","password":"password123"}`
	call("POST", "/api/v1/users/register", other, 200)
	if err := json.Unmarshal(call("POST", "/api/v1/users/login", other, 200)["data"], &session); err != nil {
		t.Fatal(err)
	}
	token = session.Authorization
	call("POST", "/api/v1/comments", `{"postId":1,"content":"Bob comment"}`, 200)
	assertComments(2)
	call("GET", "/api/v1/posts/1", "", 403)
	call("PUT", "/api/v1/posts/1", `{"title":"No"}`, 403)
	call("DELETE", "/api/v1/posts/1", "", 403)
	call("DELETE", "/api/v1/comments/1", "", 403)
	token = aliceToken
	var articleList []model.PostResponse
	if err := json.Unmarshal(call("GET", "/api/v1/posts", "", 200)["data"], &articleList); err != nil || len(articleList) != 1 || len(articleList[0].Comments) != 2 {
		t.Fatalf("article comments missing: %+v, %v", articleList, err)
	}
	call("POST", "/api/v1/posts", `{"title":"Second","content":"No comments"}`, 200)
	if string(call("GET", "/api/v1/comments?postId=2", "", 200)["data"]) != "[]" {
		t.Fatal("comments leaked between posts")
	}
	call("DELETE", "/api/v1/posts/2", "", 200)
	call("DELETE", "/api/v1/comments/2", "", 403)
	call("PUT", "/api/v1/posts/1", `{"title":"Still here","content":"Edited"}`, 200)
	assertComments(2)
	call("DELETE", "/api/v1/comments/1", "", 200)
	assertComments(1)
	token = session.Authorization
	call("DELETE", "/api/v1/comments/2", "", 200)
	token = aliceToken
	assertComments(0)
	call("DELETE", "/api/v1/comments/1", "", 404)
	call("DELETE", "/api/v1/posts/1", "", 200)
	call("GET", "/api/v1/posts/1", "", 404)
	call("GET", "/api/v1/comments?postId=1", "", 404)
	if string(call("GET", "/api/v1/posts", "", 200)["data"]) != "[]" {
		t.Fatal("expected empty list after deletion")
	}
}
