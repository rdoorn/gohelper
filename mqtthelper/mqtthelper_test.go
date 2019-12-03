package mqtthelper

import (
	"fmt"
	"log"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type TestHandler struct {
	mqtt *Handler
}

func (h *TestHandler) mqttNestOut(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
}

func TestMQTT(t *testing.T) {

	h := TestHandler{
		mqtt: New(),
	}

	// Setup MQTT Sub
	err := h.mqtt.Subscribe("testnestout", "nest/out", 0, h.mqttNestOut)
	if err != nil {
		panic(err)
	}

	timer := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timer.C:
			log.Printf("Program timeout")
			return
		}
	}
}
