package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	fiberlog "github.com/gofiber/fiber/v3/log"

	"github.com/nurkenspashev92/bookit/pkg/logger"
)

func TestInitWritesToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "logs", "app.log")

	closer, err := logger.Init(logger.Config{Level: "info", File: path})
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	t.Cleanup(func() {
		if closer != nil {
			closer.Close()
		}
		fiberlog.SetOutput(os.Stdout)
	})

	if closer == nil {
		t.Fatal("closer must be returned when a log file is configured")
	}

	fiberlog.Infow("file logging works", "key", "value")
	logger.Flush()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}

	line := string(data)
	for _, want := range []string{"[Info]", "file logging works", "key=value"} {
		if !strings.Contains(line, want) {
			t.Errorf("log file %q must contain %q", line, want)
		}
	}
}

func TestInitWithoutFileReturnsNoCloser(t *testing.T) {
	closer, err := logger.Init(logger.Config{Level: "info"})
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	if closer != nil {
		t.Error("closer must be nil when no log file is configured")
	}
}

func TestInitRotatesFileBySize(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")

	closer, err := logger.Init(logger.Config{Level: "info", File: path, MaxFileSize: 512, Backups: 2})
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	t.Cleanup(func() {
		closer.Close()
		fiberlog.SetOutput(os.Stdout)
	})

	for i := 0; i < 20; i++ {
		fiberlog.Infow("rotation", "payload", strings.Repeat("x", 128))
		logger.Flush()
	}

	if _, err := os.Stat(path + ".1"); err != nil {
		t.Fatalf("first backup must exist: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("active log must exist: %v", err)
	}
	if info.Size() > 512 {
		t.Errorf("active log grew to %d bytes, limit is 512", info.Size())
	}
}

func TestLogRecordSurvivesWithoutFlushWhenUrgent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")

	closer, err := logger.Init(logger.Config{Level: "info", File: path})
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	t.Cleanup(func() {
		closer.Close()
		fiberlog.SetOutput(os.Stdout)
	})

	fiberlog.Errorw("database unreachable", "error", "timeout")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}

	if !strings.Contains(string(data), "database unreachable") {
		t.Errorf("error records must be flushed immediately, got %q", data)
	}
}

func TestInitReportsUnwritableFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "blocker"), nil, 0o644); err != nil {
		t.Fatalf("prepare: %v", err)
	}

	closer, err := logger.Init(logger.Config{Level: "info", File: filepath.Join(dir, "blocker", "app.log")})
	if err == nil {
		t.Fatal("init must fail when the log path is not writable")
	}
	if closer != nil {
		t.Error("closer must be nil when the log file could not be opened")
	}
}
