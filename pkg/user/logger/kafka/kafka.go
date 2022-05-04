package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/kafkapubsub"
)

type topicName string

const LOG topicName = "log"

type Kafka struct {
	topic map[topicName]*pubsub.Topic
}

func New(kafkaBrokers []string) (*Kafka, error) {
	kafka := &Kafka{
		topic: make(map[topicName]*pubsub.Topic),
	}
	for _, name := range []topicName{
		// список топиков в брокере сообщений
		LOG,
	} {
		// ждем пока кафка прогрузится 5 сек
		<-time.After(5 * time.Second)
		log.Default().Printf("CreateKafkaConnect: brokers(%s) topic(%s) ", kafkaBrokers, name)
		t, err := kafkapubsub.OpenTopic(
			kafkaBrokers,
			kafkapubsub.MinimalConfig(),
			string(name),
			nil)
		if err != nil {
			return nil, err
		}
		kafka.topic[name] = t
	}
	return kafka, nil
}

func (x *Kafka) Shutdown(ctx context.Context) {
	for _, t := range x.topic {
		t.Shutdown(ctx)
	}
}

func (x *Kafka) LogNewUser(ctx context.Context) error {
	m := make(map[string]interface{})
	m["timestamp"] = time.Now().Unix()
	m["message_type"] = "INFO"
	m["message"] = "New user added"

	body, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return x.topic[LOG].Send(ctx, &pubsub.Message{
		Body: body,
	})
}
