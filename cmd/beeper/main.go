package main

import (
	"flag"
	"github.com/arvinkulagin/beeper/config"
	"github.com/arvinkulagin/beeper/handlers"
	"github.com/arvinkulagin/beeper/pubsub"
	"github.com/gorilla/mux"
	"github.com/arvinkulagin/beeper/log"
	"net/http"
)

func main() {
	logger := log.NewLogger()
	cfg := config.Default()
	flag.StringVar(&cfg.WSAddress, "ws", cfg.WSAddress, "Websocket address")
	flag.StringVar(&cfg.RESTAddress, "rest", cfg.RESTAddress, "REST address")
	flag.StringVar(&cfg.Origin, "o", cfg.Origin, "Origin URL")
	file := flag.String("config", "", "Config file")
	flag.Parse()
	if *file != "" {
		err := config.FromFile(*file, &cfg)
		if err != nil {
			logger.Err.Fatal(err)
		}
		flag.Parse()
	}
	broker := pubsub.NewBroker()
	go func() {
		wsRouter := mux.NewRouter()
		wsRouter.Handle("/{id}", handlers.WSHandler{Broker: broker, Config: cfg, Logger: logger})
		logger.Out.Println("Listening websockets on", cfg.WSAddress)
		logger.Err.Fatal(http.ListenAndServe(cfg.WSAddress, wsRouter))
	}()
	go func() {
		restRouter := mux.NewRouter()
		restRouter.Handle("/topic", handlers.List{Broker: broker, Logger: logger}).Methods("GET")
		restRouter.Handle("/topic", handlers.Add{Broker: broker, Logger: logger}).Methods("POST")
		restRouter.Handle("/topic/{id}", handlers.Pub{Broker: broker, Logger: logger}).Methods("POST")
		restRouter.Handle("/topic/{id}", handlers.Del{Broker: broker, Logger: logger}).Methods("DELETE")
		restRouter.Handle("/ping", handlers.Ping{Logger: logger}).Methods("GET")
		logger.Out.Println("Listening REST on", cfg.RESTAddress)
		logger.Err.Fatal(http.ListenAndServe(cfg.RESTAddress, restRouter))
	}()
	wait := make(chan struct{})
	<-wait
}