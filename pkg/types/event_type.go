package types

// EventType representing a type of event within the system.
type EventType string

func (etype EventType) IsValid() bool {
	switch etype {
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
