package main

import (
	"log"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/router"
	"time"
	"strconv"
	console "github.com/AsynkronIT/goconsole"
)

type myMessage struct{ i int }

func (m *myMessage) Hash() string {
	return strconv.Itoa(m.i)
}

func main() {

	act := func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *myMessage:
			log.Printf("%v got message %d", context.Self(), msg.i)
		}
	}

	log.Println("RoundRobin routing:")
	pid := actor.Spawn(router.NewRoundRobinPool(5).WithFunc(act))
	for i := 0; i < 10; i++ {
		pid.Tell(&myMessage{i})
	}
	time.Sleep(1 * time.Second)
	log.Println("")

	log.Println("Random routing:")
	pid = actor.Spawn(router.NewRandomPool(5).WithFunc(act))
	for i := 0; i < 10; i++ {
		pid.Tell(&myMessage{i})
	}
	time.Sleep(1 * time.Second)
	log.Println("")


	log.Println("ConsistentHash routing:")
	pid = actor.Spawn(router.NewConsistentHashPool(5).WithFunc(act))
	for i := 0; i < 10; i++ {
		pid.Tell(&myMessage{i})
	}
	time.Sleep(1 * time.Second)
	log.Println("")

	log.Println("BroadcastPool routing:")
	pid = actor.Spawn(router.NewBroadcastPool(5).WithFunc(act))
	for i := 0; i < 10; i++ {
		pid.Tell(&myMessage{i})
	}
	console.ReadLine()
}
