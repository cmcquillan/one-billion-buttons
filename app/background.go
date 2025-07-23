package main

import (
	"log"
	"time"
)

type ButtonEventType int

const (
	ButtonEventTypePress ButtonEventType = iota
)

type BackgroundButtonEvent struct {
	X     uint64
	Y     uint64
	ID    int64
	Event ButtonEventType
}

func BackgroundEventHandler(db ObbDb, c <-chan BackgroundButtonEvent) {
	log.Printf("Background event handler started")

	for {
		evt, open := <-c

		switch evt.Event {
		case ButtonEventTypePress:
			RecordButtonPressStatistics(db, evt)
		default:
			log.Printf("unknown event type %d", evt.Event)
		}

		log.Printf("background event %v received for %d, %d, %d", evt.Event, evt.X, evt.Y, evt.ID)

		time.Sleep(time.Millisecond * 100)

		if !open {
			break
		}
	}

	log.Printf("Background event handler stopped")
}

func RecordButtonPressStatistics(db ObbDb, evt BackgroundButtonEvent) {
	if err := db.LogButtonEvent(evt.X, evt.Y, evt.ID, evt.Event); err != nil {
		log.Printf("could not log button event: %v", err)
		return
	}

	if err := db.AdjustStat("buttons_pressed", 1); err != nil {
		log.Printf("could not adjust button press stat: %v", err)
		return
	}
}
