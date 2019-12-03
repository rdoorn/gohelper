package statsdhelper

import (
	"log"
	"os"

	"github.com/peterbourgon/g2s"
)

type Handler struct {
	*g2s.Statsd
}

func New() *Handler {
	statsdURL, ok := os.LookupEnv("STATSD_URL")
	if !ok {
		panic("missing environment key: STATSD_URL")
	}

	log.Printf("STATSD_URL: %s", statsdURL)

	// "http://localhost:8125"
	statdsClient, err := g2s.Dial("udp", statsdURL)
	if err != nil {
		panic("Couldn't connect to statsd!")
	}
	//statdsClient.Timing(1.0, "open_website", 8*time.Minute)

	return &Handler{
		statdsClient,
	}
}
