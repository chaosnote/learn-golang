package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type player struct {
	UID string
}

func enter(w http.ResponseWriter, r *http.Request) {
	log.Println("handle_enter")
	timeout := 20 * time.Second
	st := time.Now()

	// 訂單流程

	var p = player{}

	ctx := r.Context()
	ctx_value := context.WithValue(ctx, p, "value")
	ctx_first, cancel_first := context.WithCancel(ctx_value)
	defer cancel_first()
	ctx_sec, cancel_sec := context.WithTimeout(ctx_first, timeout)
	defer cancel_sec()

	/*
		go func() {
			time.Sleep(10 * time.Second)
			log.Println("值傳遞", ctx_sec.Value(p))

			// 取 DB 值 (過程略)
			// http 請求資訊 (過程略)
			// ∟ 成功(更新 DB 值)
			// ∟ 失敗(中止並還原 DB 值)

			select {
			case <-ctx.Done():
				log.Println("[2]玩家連線斷開", ctx.Err())
			case <-ctx_sec.Done():
				log.Println("[2]超時", ctx.Err())
			default:
				complete <- true
			}

			log.Println("end")
		}()
	*/

	// cls&curl http://localhost:8080

	time.Sleep(3 * time.Second)
	log.Println("等待 N 秒結束")

	select {
	case <-ctx.Done():
		log.Println("[0]玩家連線斷開", ctx.Err())
	case <-ctx_first.Done():
		log.Println("[1]介於 玩家連線斷開 至 超時之間", ctx_first.Err())
	case <-ctx_sec.Done():
		if time.Since(st) >= timeout {
			log.Println("[2-0]超時", ctx_sec.Err())
		} else {
			log.Println("[2-1]玩家連線斷開", ctx.Err())
		}
	default:
		log.Println("[3]結束", ctx.Err())
	}

	log.Println("ctx.Err():", ctx.Err())
	log.Println("ctx_first.Err():", ctx_first.Err())
	log.Println("ctx_sec.Err():", ctx_sec.Err())
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc(`/`, enter).Methods(http.MethodGet)

	log.Println("http://localhost:8080")
	e := http.ListenAndServe(":8080", router)
	if e != nil && e != http.ErrServerClosed {
		panic(e)
	}
}
