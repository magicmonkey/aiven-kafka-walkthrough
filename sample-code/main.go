/*
A demo application which periodically publishes mock messages onto
a kafka topic, which is hosted by Aiven.

The pretend payload data represents IoT sensors inside a wood drying
kiln, with readings of temperature, humidity, and weight from a load-cell.

Run this using `go run .` so it picks up the kafka.go file as well
*/

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
)

// Configure some things
var (
	serviceCertPath string = "./service.cert"
	serviceKeyPath  string = "./service.key"
	projectCaPath   string = "./ca.pem"
	serviceURI      string = "my-test-kafka-magicmonkey-f69b.aivencloud.com:17855"
	topicName       string = "kiln"
)

// Define the payload to send to Kafka
type payloadId struct {
	Id uuid.UUID `json:"id"`
}

type payloadEvent struct {
	Timestamp string  `json:"timestamp"`
	T1        float32 `json:"temperature1"`
	T2        float32 `json:"temperature2"`
	Humidity  float32 `json:"humidity"`
	Weight    int     `json:"weight"`
}

// Let's go!
func main() {

	// Initialise the RNG so we get a good selection of mock payload events
	rand.Seed(time.Now().UnixNano())

	// Initialise Kafka client (see kafka.go)
	client := getKafkaClient()

	// Make a producer object for publishing data to kafka
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	// Loop over producing a mock event every second
	for {
		payloadEvent, payloadId := getPayloadAsJson()
		fmt.Println(payloadId, payloadEvent)

		msg := &sarama.ProducerMessage{
			Topic: topicName,
			Key:   sarama.StringEncoder(payloadId),
			Value: sarama.StringEncoder(payloadEvent),
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			panic(err)
		}
		fmt.Println("+++ Written to partition", partition, "and offset", offset)
		fmt.Println("")

		time.Sleep(1 * time.Second)
	}

}

// getPayloadAsJson generates a mock payload object
func getPayloadAsJson() (evSerialized string, idSerialized string) {

	id := payloadId{
		Id: uuid.New(),
	}

	ev := payloadEvent{
		Timestamp: time.Now().Format(time.RFC3339),
		T1:        rand.Float32() * 40,
		T2:        rand.Float32() * 40,
		Humidity:  rand.Float32() * 100,
		Weight:    rand.Intn(10000),
	}

	evAsJson, err := json.Marshal(ev)
	if err != nil {
		panic(err)
	}

	idAsJson, err := json.Marshal(id)
	if err != nil {
		panic(err)
	}

	// json.Marshal outputs []byte, so turn it into strings for sarama's StringEncoder
	evSerialized = string(evAsJson)
	idSerialized = string(idAsJson)

	return
}
