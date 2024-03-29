package mqtthelper

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Handler struct {
	pub *mqtt.Client
	sub *mqtt.Client
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Publish(clientID, topic string, qos byte, retained bool, payload interface{}) error {
	if h.pub == nil {
		log.Printf("mqtt: connecting to pub")
		c, err := newclient(fmt.Sprintf("%s_pub", clientID), nil)
		if err != nil {
			return err
		}
		h.pub = &c
	}

	log.Printf("mqtt: publishing to %s: %v", topic, payload)
	if token := (*h.pub).Publish(topic, qos, retained, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (h *Handler) Subscribe(clientID, topic string, qos byte, messageHandler func(client mqtt.Client, msg mqtt.Message)) error {

	s := func(c mqtt.Client) {
		log.Printf("mqtt: subscribing to %s", topic)
		if token := c.Subscribe(topic, qos, messageHandler); token.WaitTimeout(3*time.Second) && token.Error() != nil {
			panic(token.Error())
		}
	}

	if h.sub == nil {
		log.Printf("mqtt: connecting to sub")
		c, err := newclient(fmt.Sprintf("%s_sub", clientID), s)
		if err != nil {
			return err
		}
		h.sub = &c
	}
	/*
		if token := (*h.sub).Subscribe(topic, qos, messageHandler); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	*/
	return nil
}

func (h *Handler) SubscribeWithOnConnect(clientID string, onConnect func(c mqtt.Client)) error {

	if h.sub == nil {
		log.Printf("mqtt: connecting to sub")
		c, err := newclient(fmt.Sprintf("%s_sub", clientID), onConnect)
		if err != nil {
			return err
		}
		h.sub = &c
	}
	return nil
}

func newclient(clientID string, onConnect func(c mqtt.Client)) (mqtt.Client, error) {
	urlStr := getURL()

	uri, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	opts := createClientOptions(clientID, uri)
	opts.OnConnect = onConnect
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		return nil, err
	}
	return client, nil
}

func getURL() string {
	mqttURL, ok := os.LookupEnv("MQTT_URL")
	if !ok {
		panic("missing environment key: MQTT_URL")
	}

	log.Printf("MQTT_URL: %s", mqttURL)
	return mqttURL
}

func Connect(clientId string, urlStr string, onconnect func(c mqtt.Client)) (mqtt.Client, error) {
	uri, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	opts := createClientOptions(clientId, uri)
	opts.OnConnect = onconnect
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		return nil, err
	}
	return client, nil
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

/*func listen(uri *url.URL, topic string) {
	client := Connect("sub", uri)
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
}*/

/*
func main() {
	uri, err := url.Parse(os.Getenv("CLOUDMQTT_URL"))
	if err != nil {
		log.Fatal(err)
	}
	topic := uri.Path[1:len(uri.Path)]
	if topic == "" {
		topic = "test"
	}

	go listen(uri, topic)

	client := connect("pub", uri)
	timer := time.NewTicker(1 * time.Second)
	for t := range timer.C {
		client.Publish(topic, 0, false, t.String())
	}
}
*/
