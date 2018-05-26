package main

import (
	"log"

	"github.com/syariatifaris/coevent"
)

func main() {
	covent.LoadClient("house",
		[]string{"door_opened", "door_closed"},
		[]covent.Processor{
			{
				Events:   []string{"door_opened", "microwave_on"},
				Callback: eventA,
			},
			{
				Events:   []string{"door_closed", "microwave_off"},
				Callback: eventB,
			},
		},
	)
}

func eventA(event string, value interface{}) error {
	log.Println("handling event:", event, "with value:", value)
	return nil
}

func eventB(event string, value interface{}) error {
	log.Println("handling event:", event, "with value:", value)
	return nil
}
