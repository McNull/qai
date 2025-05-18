package throbber

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Throbber represents a console-based loading animation
type Throbber struct {
	frames     []string
	interval   time.Duration
	message    string
	ctx        context.Context
	cancel     context.CancelFunc
	isRunning  bool
	lineLength int
	mu         sync.Mutex
}

// NewThrobber creates a new throbber with the given message
func NewThrobber() *Throbber {
	ctx, cancel := context.WithCancel(context.Background())
	return &Throbber{
		frames:     []string{"|", "/", "-", "\\"},
		interval:   100 * time.Millisecond,
		message:    "",
		ctx:        ctx,
		cancel:     cancel,
		isRunning:  false,
		lineLength: 0,
	}
}

// WithMessage sets the message to be displayed with the throbber
func (t *Throbber) WithMessage(message string) *Throbber {
	t.message = message
	return t
}

// WithFrames sets custom animation frames
func (t *Throbber) WithFrames(frames []string) *Throbber {
	t.frames = frames
	return t
}

// WithInterval sets custom animation speed
func (t *Throbber) WithInterval(interval time.Duration) *Throbber {
	t.interval = interval
	return t
}

func (t *Throbber) WithThrob(throb Throb) *Throbber {
	t.frames = throb.Frames
	t.interval = time.Duration(throb.Interval) * time.Millisecond
	return t
}

// Start begins the throbber animation
func (t *Throbber) Start() *Throbber {
	t.mu.Lock()
	if t.isRunning {
		t.mu.Unlock()
		return t
	}
	t.isRunning = true
	t.mu.Unlock()

	go func() {
		i := 0
		for {
			select {
			case <-t.ctx.Done():
				// Instead of calling Stop() which can cause recursion,
				// just clean up and return
				t.mu.Lock()
				if t.isRunning {
					t.isRunning = false
					// Create a blank line with spaces, not null characters
					blanks := make([]byte, t.lineLength)
					for j := range blanks {
						blanks[j] = ' '
					}
					fmt.Printf("\r%s\r", string(blanks))
				}
				t.mu.Unlock()
				return

			default:
				t.mu.Lock()
				line := fmt.Sprintf("%s %s", t.message, t.frames[i])
				if len(line) > t.lineLength {
					t.lineLength = len(line)
				}
				fmt.Printf("\r%s", line)
				t.mu.Unlock()

				i = (i + 1) % len(t.frames)
				time.Sleep(t.interval)
			}
		}
	}()

	return t
}

func (t *Throbber) IsRunning() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.isRunning
}

// Stop halts the throbber animation and clears it
func (t *Throbber) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.isRunning {
		return
	}

	t.cancel()
	t.isRunning = false

	// Create a blank line with spaces, not null characters
	blanks := make([]byte, t.lineLength)
	for j := range blanks {
		blanks[j] = ' '
	}
	fmt.Printf("\r%s\r", string(blanks))
}
