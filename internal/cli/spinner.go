// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package cli

import (
	"fmt"
	"sync"
	"time"
)

// Spinner handles the animated loading state and live metrics in the terminal.
type Spinner struct {
	mu     sync.Mutex
	active bool
	files  int
	tokens int
	start  time.Time
	stop   chan struct{}
	done   chan struct{}
}

// NewSpinner initializes a new thread-safe CLI spinner instance.
func NewSpinner() *Spinner {
	return &Spinner{}
}

// Start begins the spinner animation in a separate goroutine.
func (s *Spinner) Start(action string) {
	s.mu.Lock()
	s.active = true
	s.files = 0
	s.tokens = 0
	s.start = time.Now()
	s.stop = make(chan struct{})
	s.done = make(chan struct{})
	s.mu.Unlock()

	go func() {
		defer close(s.done)
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.stop:
				s.mu.Lock()
				f := s.files
				t := s.tokens
				elapsed := time.Since(s.start)
				s.mu.Unlock()

				m := int(elapsed.Minutes())
				sec := int(elapsed.Seconds()) % 60
				ms := int(elapsed.Milliseconds()) % 1000

				// Zero-Bug Policy: Print final state with a checkmark and newline to preserve metrics
				fmt.Printf("\r\033[K%s✓%s %s... | Files: %s%d%s | Tokens: %s%d%s | Time: %02d:%02d.%03d\n",
					ColorGreen, ColorReset, action,
					ColorGreen, f, ColorReset,
					ColorYellow, t, ColorReset,
					m, sec, ms)
				return
			default:
				s.mu.Lock()
				f := s.files
				t := s.tokens
				elapsed := time.Since(s.start)
				s.mu.Unlock()

				m := int(elapsed.Minutes())
				sec := int(elapsed.Seconds()) % 60
				ms := int(elapsed.Milliseconds()) % 1000

				fmt.Printf("\r\033[K%s%s%s %s... | Files: %s%d%s | Tokens: %s%d%s | Time: %02d:%02d.%03d",
					ColorCyan, frames[i%len(frames)], ColorReset,
					action,
					ColorGreen, f, ColorReset,
					ColorYellow, t, ColorReset,
					m, sec, ms)
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// Update safely modifies the live metrics displayed by the spinner.
func (s *Spinner) Update(files, tokens int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		s.files = files
		s.tokens = tokens
	}
}

// Stop halts the spinner animation and prints the final state.
func (s *Spinner) Stop() {
	s.mu.Lock()
	if s.active {
		s.active = false
		close(s.stop)
		s.mu.Unlock()
		<-s.done // Wait for the goroutine to finish printing the final state
		return
	}
	s.mu.Unlock()
}
