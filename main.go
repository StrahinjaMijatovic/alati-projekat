package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Config struct {
	Entries map[string]string `json:"entries"`
}

type Service struct {
	data map[string][]*Config `json:"data"`
}

func NewService() *Service {
	return &Service{
		data: make(map[string][]*Config),
	}
}

func main() {
	quit := make(chan os.Signal)
	service := NewService()

	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/config/{id}", service.AddConfig).Methods("POST")
	router.HandleFunc("/config-group/{id}", service.AddConfigGroup).Methods("POST")
	router.HandleFunc("/config/{id}", service.GetConfig).Methods("GET")
	router.HandleFunc("/config-group/{id}", service.GetConfigGroup).Methods("GET")
	router.HandleFunc("/config/{id}", service.DeleteConfig).Methods("DELETE")
	router.HandleFunc("/config-group/{id}", service.DeleteConfigGroup).Methods("DELETE")
	router.HandleFunc("/config-group/{id}", service.AddConfigToGroup).Methods("PUT")

	//log.Fatal(http.ListenAndServe(":8080", router))

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")

}

//{
//"entries": {
//"key1": "value1",
//"key2": "value2"
//}
//}
