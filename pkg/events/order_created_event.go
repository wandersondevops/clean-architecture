package events

import "time"

type OrderCreatedEvent struct {
	Payload   interface{}
	Timestamp time.Time
}

// Constructor for OrderCreatedEvent
func NewOrderCreatedEvent() *OrderCreatedEvent {
	return &OrderCreatedEvent{
		Timestamp: time.Now(), // Set the creation time
	}
}

// GetName returns the name of the event
func (e *OrderCreatedEvent) GetName() string {
	return "OrderCreated"
}

// GetPayload returns the event payload
func (e *OrderCreatedEvent) GetPayload() interface{} {
	return e.Payload
}

// SetPayload sets the event payload
func (e *OrderCreatedEvent) SetPayload(payload interface{}) {
	e.Payload = payload
}

// GetDateTime returns the time when the event was created
func (e *OrderCreatedEvent) GetDateTime() time.Time {
	return e.Timestamp
}
