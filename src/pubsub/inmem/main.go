package inmem_pubsub

import (
	"fmt"
	"time"
)

type Publisher struct{}

func (p *Publisher) Publish(topic string, msg []byte) error {
	fmt.Printf("You will publish: %q", msg)
	return nil
}

type Subscriber struct{}

func (s *Subscriber) Subscribe(topic string, cb func([]byte)) error {
	for {
		fmt.Printf("You will listen to: %q", topic)
		cb([]byte("fart face :)"))
		time.Sleep(2 * time.Second)
	}
}
