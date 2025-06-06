// Copyright 2022-2023 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}
	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "TEST_STREAM",
		Subjects: []string{"FOO.*"},
	})
	if err != nil {
		log.Fatal(err)
	}

	cons, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "TestConsumerConsume",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatal(err)
	}
	go endlessPublish(ctx, nc, js)

	cc, err := cons.Consume(
		func(msg jetstream.Msg) {
			fmt.Println(string(msg.Data()))
			msg.Ack()
		},
		jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
			fmt.Println(err)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

}

func endlessPublish(ctx context.Context, nc *nats.Conn, js jetstream.JetStream) {
	var i int
	for {
		time.Sleep(500 * time.Millisecond)
		if nc.Status() != nats.CONNECTED {
			continue
		}
		if _, err := js.Publish(ctx, "FOO.TEST1", []byte(fmt.Sprintf("msg %d", i))); err != nil {
			fmt.Println("pub error: ", err)
		}
		i++
	}
}
