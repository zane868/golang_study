package logging

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConsoleAndFileAppend(t *testing.T) {
	path := filepath.Join(t.TempDir(), "logs", "app.log")
	var console bytes.Buffer
	for _, message := range []string{"first start", "second start"} {
		logger, file, err := New(path, &console)
		if err != nil {
			t.Fatal(err)
		}
		logger.Info(message, "request_id", "test-id")
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, console.Bytes()) {
		t.Fatal("file and console differ")
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Fatal("previous logs were overwritten")
	}
	for _, line := range lines {
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatal(err)
		}
		if entry["request_id"] != "test-id" {
			t.Fatal("missing structured field")
		}
	}
}

func TestInvalidLogPath(t *testing.T) {
	path := t.TempDir()
	if _, file, err := New(path, &bytes.Buffer{}); err == nil || file != nil {
		t.Fatal("expected failure when log path is directory")
	}
}
