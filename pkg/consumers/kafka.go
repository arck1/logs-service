package consumers

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"os"
	"sync"
)

type KafkaConsumer struct {
	Config *kafka.ConfigMap
	// config =
	//	&kafka.ConfigMap{
	//	"bootstrap.servers":  kc.hosts,
	//	"group.id":           kc.groupId,
	//	"auto.offset.reset":  kc.autoOffsetReset,
	//	"enable.auto.commit": kc.autoCommit,
	//}
	Topics []string

	WaitGroup sync.WaitGroup
	CloseSignal     chan os.Signal
}

func (kc *KafkaConsumer) Stop() {
	close(kc.CloseSignal)
	kc.WaitGroup.Wait()
}

func (kc *KafkaConsumer) Listen(messages chan ConsumerMessage) error {

	kc.WaitGroup.Add(1)
	defer kc.WaitGroup.Done()

	consumer, err := kafka.NewConsumer(kc.Config)

	if err != nil {
		panic(err)
	}
	err = consumer.SubscribeTopics(kc.Topics, nil)

	if err != nil {
		consumer.Close()
		panic(err)
	}

	for {
		msg, err := consumer.ReadMessage(-1)

		if err != nil {
			log.Fatalf("Error on message consume: %v (%v)\n", err, msg)
		}

		select {
		case <-kc.CloseSignal:
			log.Println("Disconnecting...")
			break
		default:
			messages <- ConsumerMessage{msg.Key, msg.Value}
		}
	}
	consumer.Close()
	return nil
}
