package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zane868/golang_study/homework04/util"
)

func TestRequestLoggingAndRecovery(t *testing.T) {
	var logs bytes.Buffer
	original := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(original) })
	router := gin.New()
	router.Use(RequestLog(), Recover())
	router.POST("/ok", func(c *gin.Context) { c.Set("userID", uint(7)); c.Status(200) })
	router.POST("/bad", func(c *gin.Context) { util.Error(c, 401, "Unauthorized") })
	router.POST("/error", func(c *gin.Context) { util.HandleError(c, errors.New("database unavailable")) })
	router.POST("/panic", func(c *gin.Context) { panic("private-panic-content") })
	ids := make(map[string]bool)
	for _, tc := range []struct {
		path   string
		status int
		level  string
	}{
		{"/ok", 200, "INFO"}, {"/bad", 401, "WARN"}, {"/error", 500, "ERROR"}, {"/panic", 500, "ERROR"}, {"/missing", 404, "WARN"},
	} {
		logs.Reset()
		req := httptest.NewRequest("POST", tc.path+"?token=private-query", strings.NewReader(`{"password":"private-password"}`))
		req.Header.Set("Authorization", "private-token")
		req.Header.Set("Cookie", "session=private-cookie")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		id := w.Header().Get("X-Request-ID")
		if w.Code != tc.status || id == "" || ids[id] {
			t.Fatalf("bad response: %s %d %q", tc.path, w.Code, id)
		}
		ids[id] = true
		if strings.Contains(logs.String(), "private-") {
			t.Fatalf("sensitive data logged: %s", logs.String())
		}
		lines := strings.Split(strings.TrimSpace(logs.String()), "\n")
		var entry map[string]any
		if err := json.Unmarshal([]byte(lines[len(lines)-1]), &entry); err != nil {
			t.Fatal(err)
		}
		if entry["level"] != tc.level || entry["request_id"] != id || entry["status"] != float64(tc.status) || entry["duration_ms"] == nil {
			t.Fatalf("bad access log: %v", entry)
		}
		if tc.path == "/panic" && (!strings.Contains(logs.String(), "HTTP handler panic") || !strings.Contains(logs.String(), "stack")) {
			t.Fatal("missing panic log")
		}
		if tc.path == "/error" {
			if !strings.Contains(logs.String(), "database unavailable") || strings.Contains(w.Body.String(), "database unavailable") {
				t.Fatal("error must be logged but not exposed")
			}
		}
	}
}
