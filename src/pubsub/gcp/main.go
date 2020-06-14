package pubsub

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
)

// Client wraps *pubsub.Client and allows this struct to satisfy
// the PubberSubber interface
type Client struct {
	PubSubClient *pubsub.Client
}

// Message passes through the native pubsub.Message type
type Message = pubsub.Message

// Publish sends a message to a pubsub topic
// It blocks until the publish is complete, and returns the
// server-generated message id.
func (c *Client) Publish(topicName string, message []byte) (string, error) {
	ctx := context.Background()
	topic := c.PubSubClient.Topic(topicName)
	defer topic.Stop()
	r := topic.Publish(ctx, &pubsub.Message{Data: message})
	id, err := r.Get(ctx)
	if err != nil {
		return "", err
	}
	return id, nil
}

// Subscribe listens for messages on a topic and passes them to a
// custom callback function
func (c *Client) Subscribe(topicName string, cb func([]byte)) error {
	return nil
}

// Subscription wraps pubsub's native Subscription receiver
// Avoid needing to reach into Client#PubSubClient by callers.
func (c *Client) Subscription(name string) *pubsub.Subscription {
	return c.PubSubClient.Subscription(name)
}

// NewClient returns a new pubsub client for a given project
// Relies on GOOGLE_APPLICATION_CREDENTIALS env var being set.
// ADC: https://cloud.google.com/docs/authentication/production
func NewClient(project string) (*Client, error) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS must be set")
	}
	c, err := pubsub.NewClient(context.Background(), project)
	if err != nil {
		return nil, err
	}
	return &Client{c}, nil
}
