package handlers

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/arvinkulagin/beeper/pubsub"
	"github.com/arvinkulagin/beeper/log"
	"github.com/gorilla/mux"
)

type List struct {
	Broker *pubsub.Broker
	Logger log.Logger
}

func (l List) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list := l.Broker.Topics()
	out, err := json.Marshal(list)
	if err != nil {
		l.Logger.Err.Println(err)
		w.WriteHeader(500)
		return
	}
	fmt.Fprint(w, string(out))
}

type Add struct {
	Broker *pubsub.Broker
	Logger log.Logger
}

func (a Add) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := ioutil.ReadAll(r.Body) // Ограничить количество байтов из ReadAll с помощью io.LimitReader
	defer r.Body.Close()
	if err != nil {
		a.Logger.Err.Println(err)
		w.WriteHeader(403)
		return
	}
	err = a.Broker.AddTopic(string(id))
	if err != nil {
		a.Logger.Err.Println(err)
		w.WriteHeader(403)
		return
	}
	a.Logger.Out.Printf("Add %s\n", id)
}

type Del struct {
	Broker *pubsub.Broker
	Logger log.Logger
}

func (d Del) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := d.Broker.DelTopic(id)
	if err != nil {
		d.Logger.Err.Println(err)
		w.WriteHeader(404)
	}
	d.Logger.Out.Printf("Delete %s\n", id)
}

type Pub struct {
	Broker *pubsub.Broker
	Logger log.Logger
}

func (p Pub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.Logger.Err.Println(err)
		w.WriteHeader(403)
		return
	}
	err = p.Broker.Publish(id, data)
	if err != nil {
		p.Logger.Err.Println(err)
		w.WriteHeader(404)
	}
	p.Logger.Out.Printf("Publish %s: %s\n", id, string(data))
}

type Ping struct {
	Logger log.Logger
}

func (p Ping) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Logger.Out.Printf("Pong %s\n", r.RemoteAddr)
	return
}