package broadcast

import (
	"log"
	"sync"
)

type Event struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SessionBroadcaster struct {
	mu          sync.RWMutex
	subscribers map[chan Event]struct{}
}

func NewSessionBroadcaster() *SessionBroadcaster {
	return &SessionBroadcaster{
		subscribers: make(map[chan Event]struct{}),
	}
}

func (b *SessionBroadcaster) Subscribe() chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan Event, 512)
	b.subscribers[ch] = struct{}{}
	return ch
}

func (b *SessionBroadcaster) Unsubscribe(ch chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.subscribers[ch]; ok {
		delete(b.subscribers, ch)
		close(ch)
	}
}

func (b *SessionBroadcaster) ListenerCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subscribers)
}

func (b *SessionBroadcaster) Broadcast(eventType, text string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	evt := Event{Type: eventType, Text: text}
	var dead []chan Event

	for ch := range b.subscribers {
		select {
		case ch <- evt:
		default:
			dead = append(dead, ch)
		}
	}

	if len(dead) > 0 {
		b.mu.RUnlock()
		b.mu.Lock()
		for _, ch := range dead {
			if _, ok := b.subscribers[ch]; ok {
				delete(b.subscribers, ch)
				close(ch)
			}
		}
		b.mu.Unlock()
		b.mu.RLock()
	}
}

var (
	managerMu    sync.RWMutex
	broadcasters = make(map[string]*SessionBroadcaster)
)

func GetOrCreateBroadcaster(sessionID string) *SessionBroadcaster {
	managerMu.Lock()
	defer managerMu.Unlock()
	if b, ok := broadcasters[sessionID]; ok {
		return b
	}
	b := NewSessionBroadcaster()
	broadcasters[sessionID] = b
	log.Printf("Created broadcaster for session '%s'", sessionID)
	return b
}

func CleanupSession(sessionID string) {
	managerMu.Lock()
	defer managerMu.Unlock()
	delete(broadcasters, sessionID)
	log.Printf("Cleaned up session '%s'", sessionID)
}

func GetSessionStats() map[string]map[string]int {
	managerMu.RLock()
	defer managerMu.RUnlock()
	stats := make(map[string]map[string]int)
	for id, b := range broadcasters {
		stats[id] = map[string]int{"listeners": b.ListenerCount()}
	}
	return stats
}
