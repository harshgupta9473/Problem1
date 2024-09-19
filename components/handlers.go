package components

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	Producer sarama.SyncProducer
	Consumer sarama.Consumer
	Client   *redis.Client
	wg     *sync.WaitGroup
}

func CreateNewServer(producer sarama.SyncProducer, consumer sarama.Consumer, client *redis.Client) *Server {
	return &Server{
		Producer: producer,
		Consumer: consumer,
		Client:   client,
		wg   :  &sync.WaitGroup{},
}
	}


 func (s *Server)  SendProducerMsg(topic string, message []byte) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := s.Producer.SendMessage(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (s *Server) ConsumeMessage(topic string, key string) {

	partitionConsumer, _ := s.Consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	for msg := range partitionConsumer.Messages() {
		err := s.Client.Set(context.Background(), key, string(msg.Value), 0).Err()
		if err != nil {
			log.Println(err)
			return 
		}
	}
	
}

// func (s *Server)FetchFromRedis(key string){
// 	data, err := s.Client.Get(context.Background(),key).Result()
// 	if err != nil {
// 		log.Printf("Error fetching data from Redis: %v", err)
// 		return
// 	}
	
// }



func (s *Server) PostDataHandler(w http.ResponseWriter, r *http.Request) {
	var msg PostRequest
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	bytesMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.wg.Add(1)
	 go func(){
		defer s.wg.Done()
		s.SendProducerMsg("hospital", bytesMsg)
	 }()

	
	
	 
	 
	s.wg.Wait()

	response:=map[string]interface{}{
		"success":true,
		"msg":"Details of patient "+msg.PatientName+"sent successfully!",
	}

	w.Header().Set("Content-Type","application/json")
	err=json.NewEncoder(w).Encode(response)
	if err!=nil{
		http.Error(w,"Error sending details",http.StatusInternalServerError)
		return
	}
	go func(){
		
		s.ConsumeMessage("hospital",msg.HospitalID)
	 }()

}

func (s *Server)GetDataHandler(w http.ResponseWriter,r *http.Request){
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing 'key' query parameter", http.StatusBadRequest)
		return
	}
	value,err:=s.Client.Get(context.Background(),key).Result()

		if err == redis.Nil {
			http.Error(w, "No data found for key", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to get data from Redis", http.StatusInternalServerError)
			return
		}
		s.wg.Add(1)
		go func(){
			defer s.wg.Done()
			bytesMsg, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	     }
			s.SendProducerMsg("fromHospital",bytesMsg)

		}()

	s.wg.Wait()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("Data sent to Kafka: %s", value)))
	

}
