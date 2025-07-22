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
	X     uint32
	Y     uint32
	ID    int64
	Event ButtonEventType
}

func BackgroundEventHandler(c <-chan BackgroundButtonEvent) {
	for {
		evt, open := <-c

		_ = evt

		log.Printf("background event received")

		time.Sleep(time.Millisecond * 100)

		if !open {
			break
		}
	}
}
