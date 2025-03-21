package main

import (
	"log"
	"net/http"

	"github.com/olahol/melody"
)

func main() {
	m := melody.New()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		log.Println(string(msg))
		m.Broadcast([]byte("hello " + string(msg)))
	})

	http.ListenAndServe(":8080", nil)
}
