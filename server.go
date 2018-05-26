package covent

import (
	"fmt"
	"log"
	"net"
	"time"

	"container/list"
	"encoding/json"

	"github.com/dustin/go-coap"
	"github.com/tokopedia/cartapp/errors"
)

const (
	delay      = time.Millisecond * 100
	maxBuffLen = 100
)

var motes map[*net.UDPAddr]*mote

type mote struct {
	buff  *list.List
	topic string
	evts  []string
	max   int
}

func LoadObServer() {
	motes = make(map[*net.UDPAddr]*mote)

	log.Fatal(coap.ListenAndServe("udp", ":5683",
		coap.FuncHandler(func(conn *net.UDPConn, addr *net.UDPAddr, message *coap.Message) *coap.Message {
			if message.Code == coap.GET && message.Option(coap.Observe) != nil {
				if message.Path() != nil {
					log.Println(message.Path())
					if message.Path()[0] == "reg" {
						go sendMoteEvent(conn, addr, message)
					}
				}
			}
			return nil
		})))
}

func PublishEvent(topic string, evt Event) {
	for _, clt := range motes {
		if clt.topic == topic && eventRegistered(clt, evt.Name) {
			if clt.buff.Len() <= clt.max {
				clt.buff.PushBack(evt)
			}
		}
	}
}

func eventRegistered(clt *mote, evt string) bool {
	for _, e := range clt.evts {
		if e == evt {
			return true
		}
	}

	return false
}

func sendMoteEvent(conn *net.UDPConn, addr *net.UDPAddr, message *coap.Message) error {
	if _, ok := motes[addr]; ok {
		return errors.New("mote has been registered")
	}

	var regData RegEvtData
	err := json.Unmarshal(message.Payload, &regData)
	if err != nil {
		return err
	}

	motes[addr] = &mote{
		topic: regData.Topic,
		evts:  regData.Events,
		buff:  list.New(),
		max:   maxBuffLen,
	}

	log.Println("mote", addr.String(), "registered..")

	for {
		clt := motes[addr]
		for e := clt.buff.Front(); e != nil; e = e.Next() {
			if evt, ok := e.Value.(Event); ok {
				payload, err := json.Marshal(evt)
				if err != nil {
					return err
				}

				msg := coap.Message{
					Type:      coap.Acknowledgement,
					Code:      coap.Content,
					MessageID: message.MessageID,
					Payload:   payload,
				}

				msg.SetOption(coap.ContentFormat, coap.TextPlain)
				msg.SetOption(coap.LocationPath, message.Path())

				err = coap.Transmit(conn, addr, msg)
				if err != nil {
					log.Println("error transmitting:", err.Error())
					return err
				}

				if e.Prev() != nil {
					clt.buff.Remove(e.Prev())
				}

				time.Sleep(delay)
			}
		}
	}
}

func periodicTransmitter(conn *net.UDPConn, addr *net.UDPAddr, message *coap.Message) {
	now := time.Now()

	for {
		msg := coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: message.MessageID,
			Payload:   []byte(fmt.Sprintf("Been running for %v", time.Since(now))),
		}

		msg.SetOption(coap.ContentFormat, coap.TextPlain)
		msg.SetOption(coap.LocationPath, message.Path())

		log.Printf("Transmitting %v", msg)
		err := coap.Transmit(conn, addr, msg)
		if err != nil {
			log.Printf("Error on transmitter, stopping: %v", err)
			return
		}

		time.Sleep(time.Millisecond * 10)
	}
}
