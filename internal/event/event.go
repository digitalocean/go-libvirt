package event

// Event represents an internal Event.
type Event interface {
	GetCallbackID() uint32
}
