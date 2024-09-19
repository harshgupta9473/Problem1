package components

import "github.com/IBM/sarama"

func InitProducerKafka() (sarama.SyncProducer,error){
	return sarama.NewSyncProducer([]string{"localhost:9092"},nil)
}

func InitConsumerKafka()(sarama.Consumer,error){
	return sarama.NewConsumer([]string{"localhost:9092"},nil)
}