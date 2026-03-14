package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type DailyErrorHook struct {
	logDir   string
	file     *os.File
	currDate string
	mu       sync.Mutex
	formatter logrus.Formatter
}

func NewDailyErrorHook(logDir string) *DailyErrorHook {
	os.MkdirAll(logDir, 0755)

	hook := &DailyErrorHook{
		logDir:    logDir,
		formatter: &logrus.JSONFormatter{},
	}

	hook.rotateLog()

	go func() {
		for {
			hook.cleanupOldLogs()
			time.Sleep(24 * time.Hour)
		}
	}()

	return hook
}

func (h *DailyErrorHook) rotateLog() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now()
	dateStr := now.Format("2006-01-02")

	if h.currDate == dateStr && h.file != nil {
		return nil
	}

	if h.file != nil {
		h.file.Close()
	}

	filename := filepath.Join(h.logDir, fmt.Sprintf("error-%s.log", dateStr))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	h.file = file
	h.currDate = dateStr
	return nil
}

func (h *DailyErrorHook) cleanupOldLogs() {
	files, err := os.ReadDir(h.logDir)
	if err != nil {
		return
	}

	cutoff := time.Now().AddDate(0, 0, -7)
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), "error-") && strings.HasSuffix(f.Name(), ".log") {
			dateStr := strings.TrimSuffix(strings.TrimPrefix(f.Name(), "error-"), ".log")
			t, err := time.Parse("2006-01-02", dateStr)
			if err == nil && t.Before(cutoff) {
				os.Remove(filepath.Join(h.logDir, f.Name()))
			}
		}
	}
}

func (h *DailyErrorHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (h *DailyErrorHook) Fire(entry *logrus.Entry) error {
	now := time.Now()
	dateStr := now.Format("2006-01-02")

	h.mu.Lock()
	if h.currDate != dateStr {
		h.mu.Unlock()
		if err := h.rotateLog(); err != nil {
			return err
		}
		h.mu.Lock()
	}

	file := h.file
	h.mu.Unlock()

	if file == nil {
		return nil
	}

	b, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	h.mu.Lock()
	_, err = file.Write(b)
	h.mu.Unlock()

	return err
}
