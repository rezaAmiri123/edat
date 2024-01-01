package edatkafka

import (
	"github.com/rezaAmiri123/edat/msg"
	"github.com/segmentio/kafka-go"
)

type Serializer interface {
	Serialize(message msg.Message) (kafka.Message, error)
	Deserialize(message kafka.Message) (msg.Message, error)
}

var DefaultSerializer = KafkaSerializer{}

type KafkaSerializer struct{}

var _ Serializer = (*KafkaSerializer)(nil)

func (KafkaSerializer) Serialize(message msg.Message) (kafka.Message, error) {
	headers := make([]kafka.Header, 0, len(message.Headers()))

	for key, value := range message.Headers() {
		headers = append(headers, kafka.Header{
			Key:   key,
			Value: []byte(value),
		})
	}

	kafkaMsg := kafka.Message{
		Value:   message.Payload(),
		Headers: headers,
	}

	return kafkaMsg, nil
}

func (KafkaSerializer) Deserialize(message kafka.Message) (msg.Message, error) {
	var id string

	headers := make(msg.Headers, len(message.Headers))

	for _, header := range message.Headers {
		if header.Key == msg.MessageID {
			id = string(header.Value)
		} else {
			headers.Set(header.Key, string(header.Key))
		}
	}
	msg := msg.NewMessage(message.Value,
		msg.WithMessageID(id),
		msg.WithHeaders(headers),
	)

	return msg, nil
}
