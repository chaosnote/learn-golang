package main

import (
	"log"
	"net/url"
	"os"

	"github.com/chaosnote/melody"
)

func main() {

	m := &melody.Monkey{
		Melody: *melody.New(),
	}

	m.HandleConnect(func(s *melody.Session) {
		s.Write([]byte("chris"))
	})
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		log.Println(string(msg))
	})
	m.HandleDisconnect(func(s *melody.Session) {
		os.Exit(0)
	})

	// 測試對外連線功能
	m.Dial(
		url.URL{
			Scheme: "ws",
			Host:   "localhost:8080",
			Path:   "/ws",
		},
		map[string]any{},
	)

}
