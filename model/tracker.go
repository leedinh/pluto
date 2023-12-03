package model

type TrackerUpdate struct {
	Ch chan TrackerEvent
}

type TrackerEvent struct {
	EventType string
	Data      interface{}
}

func NewTrackerUpdate() *TrackerUpdate {
	return &TrackerUpdate{
		Ch: make(chan TrackerEvent),
	}
}

func (t *TrackerUpdate) GetUpdates() {
	for {
		select {
		case <-t.Ch:

		}
	}
}
