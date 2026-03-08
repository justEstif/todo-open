package events_test

import (
	"testing"
	"time"

	"github.com/justEstif/todo-open/internal/events"
)

func TestBrokerFanOut(t *testing.T) {
	t.Parallel()
	b := events.NewBroker()
	ch1, unsub1 := b.Subscribe(4)
	ch2, unsub2 := b.Subscribe(4)
	defer unsub1()
	defer unsub2()

	e := events.Event{Type: events.TypeCreated, At: time.Now()}
	b.Publish(e)

	got1 := <-ch1
	got2 := <-ch2
	if got1.Type != events.TypeCreated {
		t.Errorf("sub1 got type %q, want %q", got1.Type, events.TypeCreated)
	}
	if got2.Type != events.TypeCreated {
		t.Errorf("sub2 got type %q, want %q", got2.Type, events.TypeCreated)
	}
}

func TestBrokerUnsubscribe(t *testing.T) {
	t.Parallel()
	b := events.NewBroker()
	_, unsub := b.Subscribe(4)
	unsub()

	// Publish after unsubscribe should not block or panic.
	b.Publish(events.Event{Type: events.TypeUpdated, At: time.Now()})
}

func TestBrokerSlowSubscriberDrops(t *testing.T) {
	t.Parallel()
	b := events.NewBroker()
	// buf=1 so second publish is dropped
	ch, unsub := b.Subscribe(1)
	defer unsub()

	b.Publish(events.Event{Type: events.TypeCreated, At: time.Now()})
	b.Publish(events.Event{Type: events.TypeUpdated, At: time.Now()}) // dropped

	got := <-ch
	if got.Type != events.TypeCreated {
		t.Errorf("got type %q, want %q", got.Type, events.TypeCreated)
	}
	// channel should have no second event
	select {
	case extra := <-ch:
		t.Errorf("unexpected extra event: %v", extra)
	default:
	}
}
