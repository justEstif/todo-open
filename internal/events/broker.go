// Package events provides an in-process publish/subscribe broker for task events.
package events

import (
	"sync"
	"time"

	"github.com/justEstif/todo-open/internal/core"
)

// Type constants for task events.
const (
	TypeCreated       = "task.created"
	TypeUpdated       = "task.updated"
	TypeDeleted       = "task.deleted"
	TypeStatusChanged = "task.status_changed"
)

// Event represents a task domain event.
type Event struct {
	Type      string           `json:"type"`
	Task      *core.Task       `json:"task,omitempty"`
	OldStatus *core.TaskStatus `json:"old_status,omitempty"`
	NewStatus *core.TaskStatus `json:"new_status,omitempty"`
	At        time.Time        `json:"at"`
}

// Broker is an in-process fan-out event broker.
// Slow subscribers drop events rather than blocking the broker.
type Broker struct {
	mu   sync.RWMutex
	subs map[uint64]chan Event
	next uint64
}

// NewBroker returns a ready-to-use Broker.
func NewBroker() *Broker {
	return &Broker{subs: make(map[uint64]chan Event)}
}

// Subscribe registers a subscriber. It returns a channel that receives events
// and an unsubscribe function that must be called when done.
// bufSize controls the subscriber channel buffer; events are dropped when full.
func (b *Broker) Subscribe(bufSize int) (<-chan Event, func()) {
	if bufSize <= 0 {
		bufSize = 64
	}
	ch := make(chan Event, bufSize)
	b.mu.Lock()
	id := b.next
	b.next++
	b.subs[id] = ch
	b.mu.Unlock()

	unsub := func() {
		b.mu.Lock()
		delete(b.subs, id)
		b.mu.Unlock()
		// Non-blocking sends mean the broker never blocks on a removed subscriber.
		// Drain any buffered events so the GC can reclaim the channel.
		for {
			select {
			case <-ch:
			default:
				return
			}
		}
	}
	return ch, unsub
}

// FromMutation converts a core.MutationEvent to a broker Event.
// Use this at wiring boundaries so the field mapping lives in one place.
func FromMutation(m core.MutationEvent) Event {
	return Event{
		Type:      m.Type,
		Task:      m.Task,
		OldStatus: m.OldStatus,
		NewStatus: m.NewStatus,
		At:        m.At,
	}
}

// Publish fans out an event to all current subscribers.
// Subscribers that are full drop the event (non-blocking send).
func (b *Broker) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subs {
		select {
		case ch <- e:
		default:
			// subscriber too slow; drop event
		}
	}
}
