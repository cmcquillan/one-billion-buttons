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

	ticker := time.NewTicker(time.Second * 2)
	presses := make([]BackgroundButtonEvent, 0, 1000)

	for {
		closed := false

		select {
		case <-ticker.C:
			if len(presses) > 0 {
				RecordButtonPress(db, presses)
				log.Printf("processing %d button press events", len(presses))
				presses = make([]BackgroundButtonEvent, 0, 1000)
			}
		case evt, open := <-c:
			switch evt.Event {
			case ButtonEventTypePress:
				log.Printf("background event %v received for %d, %d, %d", evt.Event, evt.X, evt.Y, evt.ID)
				presses = append(presses, evt)
			default:
				log.Printf("unknown event type %s", evt.Event)
			}

			if !open {
				closed = true
			}
		}

		if closed {
			break
		}
	}

	RecordButtonPress(db, presses)

	time.Sleep(time.Millisecond * 100)

	log.Printf("Background event handler stopped")
}

func RecordButtonPress(db ObbDb, events []BackgroundButtonEvent) {
	err := db.LogButtonEvents(events)

	if err != nil {
		log.Printf("could not save button press events %v", err)
	}
}

func BackgroundComputeStatistics(db ObbDb, ctx context.Context) {
	log.Printf("Background statistics worker started")
	ticker := time.NewTicker(time.Second * 120)
	done := false
	for !done {
		select {
		case <-ctx.Done():
			log.Printf("Background statistics worker stopping")
			done = true
		case <-ticker.C:
			log.Printf("Refreshing stats")
			if err := db.RefreshStats(); err != nil {
				log.Printf("could not refresh stats: %v", err)
			}
		}
	}

	log.Printf("Background statistics worker stopped")
}
