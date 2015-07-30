package handlers

import (
	"github.com/arvinkulagin/beeper/config"
	"github.com/arvinkulagin/beeper/log"
	"github.com/arvinkulagin/beeper/pubsub"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

type WSHandler struct {
	Broker *pubsub.Broker
	Config config.Config
	Logger log.Logger
}

func (h WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if h.Config.Origin == "" {
				return true
			}
			if r.Header.Get("Origin") == h.Config.Origin {
				return true
			}
			h.Logger.Out.Printf("Deny %s: wrong origin %s\n", r.RemoteAddr, r.Header.Get("Origin"))
			return false
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Logger.Err.Println(err)
		return
	}
	defer conn.Close()
	receiver, err := h.Broker.Subscribe(id)
	if err != nil {
		h.Logger.Err.Println(err)
		return
	}
	h.Logger.Out.Printf("Subscribe %s: %s\n", conn.RemoteAddr().String(), id)
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
		receiver.Close()
		h.Logger.Out.Printf("Unsubscribe %s: %s\n", conn.RemoteAddr().String(), id)
	}()
	for {
		message, err := receiver.Receive()
		if err != nil {
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}











