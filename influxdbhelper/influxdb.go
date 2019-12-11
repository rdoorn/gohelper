package influxdbhelper

import (
	"os"
	"time"

	_ "github.com/influxdata/influxdb1-client"
	client "github.com/influxdata/influxdb1-client/v2"
)

type Handler struct {
	client client.Client
	db     string
}

func New(db string) *Handler {
	influxAddr, ok := os.LookupEnv("INFLUXDB_URL")
	if !ok {
		panic("Missing environment variable INFLUXDB_URL")
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: influxAddr,
	})

	if err != nil {
		panic(err)
	}

	return &Handler{
		client: c,
		db:     db,
	}
}

func (h *Handler) Insert(bucket string, tags map[string]string, fields map[string]interface{}) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  h.db,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	pt, err := client.NewPoint(bucket, tags, fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := h.client.Write(bp); err != nil {
		return err
	}

	return nil
}

func (h *Handler) Close() {
	h.client.Close()
}
