// Config Service API
//
//	   Title: Config Service API
//
//	   Schemes: http
//		  Version: 0.0.1
//		  BasePath: /
//
//		  Produces:
//			- application/json
//
// swagger:meta
package main

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type dbServerConfig struct {
	data      map[string]*DBConfig
	dataGroup map[string]*DBGroup
}

func main() {

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	server := dbServerConfig{
		data:      map[string]*DBConfig{},
		dataGroup: map[string]*DBGroup{},
	}

	router.HandleFunc("/config/", server.createConfigHandler).Methods("POST")
	router.HandleFunc("/configs/", server.getAllConfigHandler).Methods("GET")
	router.HandleFunc("/configs/{id}/", server.getConfigVersionsHandler).Methods("GET")
	router.HandleFunc("/configs/{id}/", server.delConfigVersionsHandler).Methods("DELETE")
	router.HandleFunc("/config/{id}/{version}/", server.getConfigHandler).Methods("GET")
	router.HandleFunc("/config/{id}/{version}/", server.delConfigHandler).Methods("DELETE")

	router.HandleFunc("/group/", server.createGroupHandler).Methods("POST")
	router.HandleFunc("/groups/", server.getAllGroupHandler).Methods("GET")
	router.HandleFunc("/groups/{id}/", server.getGroupVersionsHandler).Methods("GET")
	router.HandleFunc("/groups/{id}/", server.delGroupVersionsHandler).Methods("DELETE")
	router.HandleFunc("/group/{id}/{version}/", server.getGroupHandler).Methods("GET")
	router.HandleFunc("/group/{id}/{version}/", server.delGroupHandler).Methods("DELETE")
	router.HandleFunc("/group/{id}/{version}/{new}/", server.appendGroupHandler).Methods("POST")

	router.HandleFunc("/group/{id}/{version}/{labels}/", server.getConfigsByLabel).Methods("GET")
	router.HandleFunc("/group/{id}/{version}/{new}/{labels}/", server.delConfigsByLabel).Methods("DELETE")

	router.HandleFunc("/swagger.yaml", server.swaggerHandler).Methods("GET")
	router.Path("/metrics").Handler(metricsHandler())

	// SwaggerUI
	optionsDevelopers := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	developerDocumentationHandler := middleware.SwaggerUI(optionsDevelopers, nil)
	router.Handle("/docs", developerDocumentationHandler)

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("...server starting...")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("...server shutting down...")

	//graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("...server stopped...")

}
