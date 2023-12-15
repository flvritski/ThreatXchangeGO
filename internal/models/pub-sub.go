package models

type Event struct {
	CollectionID int
	Object       Object
}

type EventBus struct {
	subscribers map[int]chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[int]chan Event),
	}
}

func (eb *EventBus) Subscribe(collectionID int) chan Event {
	ch := make(chan Event)
	eb.subscribers[collectionID] = ch
	return ch
}

func (eb *EventBus) Unsubscribe(collectionID int, ch chan Event) {
	close(ch)
	delete(eb.subscribers, collectionID)
}

func (eb *EventBus) Publish(event Event) {
	for _, ch := range eb.subscribers {
		ch <- event
	}
}

var globalEventBus = NewEventBus()
