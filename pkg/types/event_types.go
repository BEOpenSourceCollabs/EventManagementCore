package types

// Represents various types of events
type EventType string

func (event_type EventType) IsValid() bool {
	switch event_type {
	case OfflineEventType, OnlineEventType, BothEventType:
		return true
	default:
		return false
	}
}

const (
	OfflineEventType EventType = "offline"
	OnlineEventType  EventType = "online"
	BothEventType    EventType = "both"
)
