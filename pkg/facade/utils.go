package facade

import (
	"io"
	"math/rand"
	"time"
)

// SlowWriter wraps an io.Writer and writes data in random-sized chunks,
// adding a random delay between chunks to simulate slow or unstable connections.
type SlowWriter struct {
	W            io.Writer
	MinDelay     time.Duration
	MaxDelay     time.Duration
	MinChunkSize int
	MaxChunkSize int
	rng          *rand.Rand
}

// NewSlowWriter constructs a SlowWriter with randomized delay and chunk size behavior.
func NewSlowWriter(w io.Writer, minDelay, maxDelay time.Duration, minChunkSize, maxChunkSize int) *SlowWriter {
	if minChunkSize <= 0 {
		minChunkSize = 1
	}
	if maxChunkSize < minChunkSize {
		maxChunkSize = minChunkSize
	}
	return &SlowWriter{
		W:            w,
		MinDelay:     minDelay,
		MaxDelay:     maxDelay,
		MinChunkSize: minChunkSize,
		MaxChunkSize: maxChunkSize,
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Write splits data into random-sized chunks and sleeps a random amount between writes.
func (s *SlowWriter) Write(p []byte) (int, error) {
	totalWritten := 0
	for totalWritten < len(p) {
		// Pick random chunk size
		chunkSize := s.MinChunkSize
		if s.MaxChunkSize > s.MinChunkSize {
			chunkSize += s.rng.Intn(s.MaxChunkSize - s.MinChunkSize + 1)
		}
		end := totalWritten + chunkSize
		if end > len(p) {
			end = len(p)
		}
		n, err := s.W.Write(p[totalWritten:end])
		totalWritten += n
		if err != nil {
			return totalWritten, err
		}
		// Pick random delay
		delay := s.MinDelay
		if s.MaxDelay > s.MinDelay {
			delay += time.Duration(s.rng.Int63n(int64(s.MaxDelay - s.MinDelay)))
		}
		time.Sleep(delay)
	}

	return totalWritten, nil
}
