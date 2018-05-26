package main

import (
	"log"
	"time"

	"math/rand"

	"github.com/syariatifaris/coevent"
)

const (
	topic    = "house"
	pubDelay = time.Millisecond * 50
)

func main() {
	stop := make(chan bool)
	go func() {
		tick := time.NewTicker(pubDelay)
		c := tick.C
		for {
			select {
			case <-c:
				covent.PublishEvent(topic, generateEvent())
			case <-stop:
				tick.Stop()
				return
			default:
				continue
			}
		}
	}()

	go func() {
		time.Sleep(time.Second * 10)
		stop <- true
	}()

	log.Println("starting server..")
	covent.LoadObServer()
}

type EventData struct {
	Message string `json:"message"`
}

func generateEvent() covent.Event {
	var name, msg string
	t := rand.Intn(2)
	switch t {
	case 0:
		name = "door_closed"
		msg = "do something after door is closed"
	case 1:
		name = "door_opened"
		msg = "do something after door is opened"
	}

	return covent.Event{
		Name: name,
		Content: EventData{
			Message: msg,
		},
	}
}
