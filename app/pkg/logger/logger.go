package logger

import (
	"io"
	"os"
	"strings"
	"sync/atomic"
	"time"

	fiberlog "github.com/gofiber/fiber/v3/log"
)

const (
	bufferSize  = 256 << 10
	flushPeriod = 500 * time.Millisecond
	maxFileSize = 128 << 20
	backupCount = 3
)

type Config struct {
	Level       string
	File        string
	Console     bool
	MaxFileSize int64
	Backups     int
}

var (
	level  atomic.Int32
	writer atomic.Pointer[bufferedWriter]
)

func Init(cfg Config) (io.Closer, error) {
	applyLevel(cfg.Level)
	fiberlog.MustSetContextTemplate(fiberlog.ContextConfig{Format: fiberlog.RequestIDFormat})

	closePrevious()

	if cfg.File == "" {
		fiberlog.SetOutput(os.Stdout)
		return nil, nil
	}

	file, err := newRotatingFile(cfg.File, rotationSize(cfg.MaxFileSize), rotationBackups(cfg.Backups))
	if err != nil {
		fiberlog.SetOutput(os.Stdout)
		return nil, err
	}

	var target io.Writer = file
	if cfg.Console {
		target = io.MultiWriter(os.Stdout, file)
	}

	buffered := newBufferedWriter(target, file, bufferSize, flushPeriod)
	writer.Store(buffered)
	fiberlog.SetOutput(buffered)

	return buffered, nil
}

func Enabled(l fiberlog.Level) bool {
	return l >= fiberlog.Level(level.Load())
}

func Flush() {
	if buffered := writer.Load(); buffered != nil {
		_ = buffered.Flush()
	}
}

func ParseLevel(name string) fiberlog.Level {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "trace":
		return fiberlog.LevelTrace
	case "debug":
		return fiberlog.LevelDebug
	case "warn", "warning":
		return fiberlog.LevelWarn
	case "error":
		return fiberlog.LevelError
	case "fatal":
		return fiberlog.LevelFatal
	default:
		return fiberlog.LevelInfo
	}
}

func applyLevel(name string) {
	parsed := ParseLevel(name)
	fiberlog.SetLevel(parsed)
	level.Store(int32(parsed))
}

func rotationSize(size int64) int64 {
	if size <= 0 {
		return maxFileSize
	}
	return size
}

func rotationBackups(count int) int {
	if count <= 0 {
		return backupCount
	}
	return count
}

func closePrevious() {
	if previous := writer.Swap(nil); previous != nil {
		_ = previous.Close()
	}
}
