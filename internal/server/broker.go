// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package server

import (
	"encoding/json"
	"fmt"
	"time"
)

// Broker manages Server-Sent Events (SSE) clients and broadcasts messages.
// Architecture Pillar: Fully encapsulated to prevent global state mutations.
type Broker struct {
	Notifier       chan string
	newClients     chan chan string
	closingClients chan chan string
	clients        map[chan string]bool
}

// NewBroker initializes a new SSE broker and starts its listening goroutine.
func NewBroker() *Broker {
	b := &Broker{
		// Zero-Bug Policy: Increased buffer to 10000 to handle extreme concurrency speeds without dropping logs
		Notifier:       make(chan string, 10000),
		newClients:     make(chan chan string),
		closingClients: make(chan chan string),
		clients:        make(map[chan string]bool),
	}
	go b.listen()
	return b
}

// listen continuously monitors channels for new clients, disconnections, and broadcast messages.
func (b *Broker) listen() {
	for {
		select {
		case s := <-b.newClients:
			b.clients[s] = true
		case s := <-b.closingClients:
			// Zero-Bug Policy: Safely delete and close the channel to prevent Goroutine/Memory leaks
			if _, ok := b.clients[s]; ok {
				delete(b.clients, s)
				close(s)
			}
		case event := <-b.Notifier:
			for clientMessageChan := range b.clients {
				// Non-blocking send to prevent a slow client from blocking the entire broker
				select {
				case clientMessageChan <- event:
				default:
				}
			}
		}
	}
}

// AddClient registers a new SSE client channel.
func (b *Broker) AddClient(client chan string) {
	b.newClients <- client
}

// RemoveClient unregisters an existing SSE client channel.
func (b *Broker) RemoveClient(client chan string) {
	b.closingClients <- client
}

// SendLog safely broadcasts a structured JSON log to the UI.
func (b *Broker) SendLog(level, action, path string) {
	event := map[string]string{
		"type":      "log",
		"level":     level,
		"timestamp": time.Now().Format("15:04:05"),
		"action":    action,
		"path":      path,
	}
	data, _ := json.Marshal(event)

	// Non-blocking send to the Notifier
	select {
	case b.Notifier <- string(data):
	default:
	}
}

// SendMetric safely broadcasts live processing metrics to the UI.
func (b *Broker) SendMetric(files, tokens int, elapsed time.Duration) {
	speed := 0
	if elapsed.Seconds() > 0 {
		speed = int(float64(tokens) / elapsed.Seconds())
	}

	m := int(elapsed.Minutes())
	s := int(elapsed.Seconds()) % 60
	ms := int(elapsed.Milliseconds()) % 1000

	event := map[string]interface{}{
		"type":           "metric",
		"filesProcessed": files,
		"tokens":         tokens,
		"elapsedTime":    fmt.Sprintf("%02d:%02d.%03d", m, s, ms),
		"speed":          fmt.Sprintf("%d T/s", speed),
	}
	data, _ := json.Marshal(event)

	// Non-blocking send to the Notifier
	select {
	case b.Notifier <- string(data):
	default:
	}
}
