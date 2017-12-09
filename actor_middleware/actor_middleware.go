package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor/middleware"
	"github.com/AsynkronIT/goconsole"
	"time"
	"log"
)

type myActor struct{}
type hello struct{ Who string }

func (state *myActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case hello:
		//this actor have been initialized by the receive pipeline
		fmt.Printf("recieved %v\n", msg)
		context.Respond("Hello " + msg.Who + "apang")
	}
}

func tellCorrelated(message interface{}, target *actor.PID, cid string) {
	header := make(map[string]string)
	header["cid"] = cid

	ev := actor.MessageEnvelope{
		Message: message,
		Sender:  nil,
		Header:  header,
	}

	target.Tell(ev)
}

func askCorrelated(message interface{}, target *actor.PID, cid string) *actor.Future {
	header := make(map[string]string)
	header["cid"] = cid

	future := actor.NewFuture(30 * time.Second)

	ev := actor.MessageEnvelope{
		Message: message,
		Sender:  future.PID(),
		Header:  header,
	}

	target.Tell(&ev)
	return future
}

func outLogger(next actor.SenderFunc) actor.SenderFunc {
	fn := func(c actor.Context, target *actor.PID, envelope *actor.MessageEnvelope) {
		message := c.Message()
		log.Printf("%v on recieved %v send %v", c.Self(), message, envelope)
		next(c, target, envelope)
	}

	return fn
}

func outCorrelated(next actor.SenderFunc) actor.SenderFunc {
	fn := func(c actor.Context, target *actor.PID, envelope *actor.MessageEnvelope) {

		header := make(map[string]string)
		header["cid"] = c.MessageHeader().Get("cid")

		ev := actor.MessageEnvelope{
			Message: envelope.Message,
			Sender:  envelope.Sender,
			Header:  header,
		}

		next(c, target, &ev)
	}

	return fn
}

func main() {
	props := actor.
		FromInstance(&myActor{}).
		WithMiddleware(middleware.Logger).
		WithOutboundMiddleware(outCorrelated, outLogger)

	pid := actor.Spawn(props)
	result, _ := askCorrelated(hello{"andi"}, pid, "testcorrelated").Result()
	log.Print(result)
	console.ReadLine()
}
