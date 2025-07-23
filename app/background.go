package main

import (
	"context"
	"log"
	"time"
)

type ButtonEventType string

const (
	ButtonEventTypePress ButtonEventType = "press"
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
			RecordButtonPress(db, evt)
		default:
			log.Printf("unknown event type %s", evt.Event)
		}

		log.Printf("background event %v received for %d, %d, %d", evt.Event, evt.X, evt.Y, evt.ID)

		time.Sleep(time.Millisecond * 100)

		if !open {
			break
		}
	}

	log.Printf("Background event handler stopped")
}

func RecordButtonPress(db ObbDb, evt BackgroundButtonEvent) {
	if err := db.LogButtonEvent(evt.X, evt.Y, evt.ID, evt.Event); err != nil {
		log.Printf("could not log button event: %v", err)
		return
	}

	if err := db.AdjustStat("buttons_pressed", 1); err != nil {
		log.Printf("could not adjust button press stat: %v", err)
		return
	}
}

func BackgroundComputeStatistics(db ObbDb, ctx context.Context) {
	log.Printf("Background statistics worker started")
	done := false
	for !done {
		select {
		case <-ctx.Done():
			log.Printf("Background statistics worker stopping")
			done = true
		case <-time.After(time.Second * 120):
			log.Printf("Refreshing stats")
			if err := db.RefreshStats(); err != nil {
				log.Printf("could not refresh stats: %v", err)
			}
		}
	}

	log.Printf("Background statistics worker stopped")
}
