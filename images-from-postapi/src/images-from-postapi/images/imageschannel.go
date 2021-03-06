package images

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ImagesChannel struct {
	DeviceID  string
	Hostname  string
	TopicName string
	Producer  *kafka.Producer
}

// NewImagesChannel constructs ImagesChannel and sets up connection
func NewImagesChannel(deviceID, hostname, topicName string) (imagesChannel *ImagesChannel) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": hostname})
	if err != nil {
		fmt.Printf("Failed to create '%s' producer: %s\n", topicName, err)
		os.Exit(1)
	}
	return &ImagesChannel{
		DeviceID:  deviceID,
		Hostname:  hostname,
		TopicName: topicName,
		Producer:  producer,
	}
}

// PostImage is an all-in-one helper to post an image to kafka topic
func (imagesChannel *ImagesChannel) PostImage(image *[]byte) {
	// Optional delivery channel, if not specified the Producer object's
	// .Events channel is used.
	deliveryChan := make(chan kafka.Event)

	err := imagesChannel.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &imagesChannel.TopicName,
			Partition: kafka.PartitionAny,
		},
		Value: *image,
	}, deliveryChan)
	if err != nil {
		fmt.Printf("Failed .Produce for'images-topic' message: %s\n", err)
		os.Exit(1)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		// } else {
		// 	fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
		// 		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
}
