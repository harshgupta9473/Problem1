package components

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RunRoutes() {
	route := mux.NewRouter()

	route.HandleFunc("/post", s.PostDataHandler)
	route.HandleFunc("/get", s.GetDataHandler)
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", route))
	
}
