package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	opensearch "github.com/opensearch-project/opensearch-go"
)

type DeviceData struct {
	DeviceID    string    `json:"device_id"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Timestamp   time.Time `json:"timestamp"`
}

func main() {
	cfg := opensearch.Config{
		Addresses: []string{"https://localhost:9200"},
		Username:  "",
		Password:  "!",
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	client, err := opensearch.NewClient(cfg)
	if err != nil {
		panic(fmt.Sprintf("Erro ao criar cliente: %s", err))
	}

	indexName := "device_data"

	for i := 0; i < 10; i++ {
		data := DeviceData{
			DeviceID:    fmt.Sprintf("sensor-%d", rand.Intn(5)+1),
			Temperature: 20.0 + rand.Float64()*10.0,
			Humidity:    40.0 + rand.Float64()*20.0,
			Timestamp:   time.Now(),
		}

		body, _ := json.Marshal(data)

		res, err := client.Index(
			indexName,
			bytes.NewReader(body),
			client.Index.WithContext(context.Background()),
			client.Index.WithRefresh("true"),
		)
		if err != nil {
			fmt.Println("Erro ao indexar:", err)
		} else {
			fmt.Printf("Enviado: %+v\n", data)
			res.Body.Close()
		}

		time.Sleep(1 * time.Second)
	}
}
