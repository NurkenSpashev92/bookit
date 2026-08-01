package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	logFileMode = 0o644
	logDirMode  = 0o755
)

type rotatingFile struct {
	path    string
	maxSize int64
	backups int
	file    *os.File
	size    int64
}

func newRotatingFile(path string, maxSize int64, backups int) (*rotatingFile, error) {
	file, size, err := openLogFile(path)
	if err != nil {
		return nil, err
	}

	return &rotatingFile{
		path:    path,
		maxSize: maxSize,
		backups: backups,
		file:    file,
		size:    size,
	}, nil
}

func (f *rotatingFile) Write(p []byte) (int, error) {
	if f.size+int64(len(p)) > f.maxSize {
		if err := f.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := f.file.Write(p)
	f.size += int64(n)

	return n, err
}

func (f *rotatingFile) Close() error {
	return f.file.Close()
}

func (f *rotatingFile) rotate() error {
	if err := f.file.Close(); err != nil {
		return err
	}

	_ = os.Remove(backupPath(f.path, f.backups))
	for i := f.backups - 1; i >= 1; i-- {
		_ = os.Rename(backupPath(f.path, i), backupPath(f.path, i+1))
	}
	_ = os.Rename(f.path, backupPath(f.path, 1))

	file, size, err := openLogFile(f.path)
	if err != nil {
		return err
	}

	f.file = file
	f.size = size

	return nil
}

func backupPath(path string, index int) string {
	return path + "." + strconv.Itoa(index)
}

func openLogFile(path string) (*os.File, int64, error) {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, logDirMode); err != nil {
			return nil, 0, fmt.Errorf("create log directory %q: %w", dir, err)
		}
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, logFileMode)
	if err != nil {
		return nil, 0, fmt.Errorf("open log file %q: %w", path, err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, 0, fmt.Errorf("stat log file %q: %w", path, err)
	}

	return file, info.Size(), nil
}
