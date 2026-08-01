package logger

import (
	"bufio"
	"bytes"
	"io"
	"sync"
	"time"
)

var urgentTags = [][]byte{[]byte("[Error]"), []byte("[Fatal]"), []byte("[Panic]")}

type bufferedWriter struct {
	mu     sync.Mutex
	buf    *bufio.Writer
	closer io.Closer
	stop   chan struct{}
	once   sync.Once
	wg     sync.WaitGroup
}

func newBufferedWriter(target io.Writer, closer io.Closer, size int, period time.Duration) *bufferedWriter {
	w := &bufferedWriter{
		buf:    bufio.NewWriterSize(target, size),
		closer: closer,
		stop:   make(chan struct{}),
	}

	w.wg.Add(1)
	go w.flushLoop(period)

	return w
}

func (w *bufferedWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	n, err := w.buf.Write(p)
	if err != nil || !isUrgent(p) {
		return n, err
	}

	return n, w.buf.Flush()
}

func (w *bufferedWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.buf.Flush()
}

func (w *bufferedWriter) Close() error {
	w.once.Do(func() { close(w.stop) })
	w.wg.Wait()

	w.mu.Lock()
	defer w.mu.Unlock()

	err := w.buf.Flush()
	if w.closer == nil {
		return err
	}

	if closeErr := w.closer.Close(); err == nil {
		err = closeErr
	}

	return err
}

func (w *bufferedWriter) flushLoop(period time.Duration) {
	defer w.wg.Done()

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = w.Flush()
		case <-w.stop:
			return
		}
	}
}

func isUrgent(record []byte) bool {
	for _, tag := range urgentTags {
		if bytes.Contains(record, tag) {
			return true
		}
	}
	return false
}
