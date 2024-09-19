package main

import (
	"kafkaReddis/components"
	"log"
)

func main() {
	producer, err := components.InitProducerKafka()
	if err!=nil{
		log.Fatal("problem in initialising producer")
		return
	}
	consumer,err:=components.InitConsumerKafka()
	if err!=nil{
		log.Fatal("problem in initialising consumer")
		return
	}
	client,err:=components.InitRedis()
	if err!=nil{
		log.Fatal("problem in initialising Reddis")
		return
	}
	server:=components.CreateNewServer(producer,consumer,client)
	server.RunRoutes()
}