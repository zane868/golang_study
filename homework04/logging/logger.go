package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// New 自动创建日志目录，追加写入文件，同时保留控制台输出。
// 调用方负责关闭返回的文件。
func New(path string, console io.Writer) (*slog.Logger, *os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, nil, fmt.Errorf("创建日志目录: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, nil, fmt.Errorf("打开日志文件: %w", err)
	}
	writer := io.MultiWriter(file, console)
	logger := slog.New(slog.NewJSONHandler(writer, nil))
	return logger, file, nil
}
