package covent

import (
	"log"

	"encoding/json"

	"github.com/dustin/go-coap"
)

type EvtFunc func(event string, value interface{}) error

type Processor struct {
	Callback EvtFunc
	Events   []string
}

func LoadClient(topic string, listenEvts []string, procs []Processor) {
	regPayload, _ := json.Marshal(RegEvtData{
		Topic:  topic,
		Events: listenEvts,
	})

	req := coap.Message{
		Type:      coap.NonConfirmable,
		Code:      coap.GET,
		MessageID: 12345,
		Payload:   regPayload,
	}

	req.AddOption(coap.Observe, 1)
	req.SetPathString("/reg")

	c, err := coap.Dial("udp", "localhost:5683")
	if err != nil {
		log.Fatalln("error dialing:", err.Error())
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalln("error sending request:", err.Error())
	}

	for err == nil {
		if rv != nil {
			if err != nil {
				log.Fatalln("error receiving:", err.Error())
			}

			var evt Event
			err := json.Unmarshal(rv.Payload, &evt)
			if err != nil {
				log.Fatalln("unable to marshall event data", err.Error())
			}

			go func(event Event) {
				for _, proc := range procs {
					if eventListed(proc, event.Name) {
						proc.Callback(event.Name, event.Content)
					}
				}
			}(evt)
		}
		rv, err = c.Receive()

	}
	log.Println("mote stopped..", err)
}

func eventListed(proc Processor, evtName string) bool {
	for _, event := range proc.Events {
		if evtName == event {
			return true
		}
	}

	return false
}
