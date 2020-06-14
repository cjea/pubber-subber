package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cjea/pubber-subber/src/config"
	pubsub "github.com/cjea/pubber-subber/src/pubsub/gcp"

	"github.com/valyala/fasthttp"
)

func deliverAndAck(ctx context.Context, m *pubsub.Message) {
	endpoint := os.Getenv("RECEIVE_ENDPOINT")
	if endpoint == "" {
		panic("Must set RECEIVE_ENDPOINT to subscribe to topics.")
	}
	err := _post(endpoint, m.Data)
	if err != nil {
		fmt.Println("Something bad happened :(")
		return
	}
	m.Ack()
	fmt.Printf("Acked a message %q :)\n", m.ID)
}

func publishHandler(ctx *fasthttp.RequestCtx, clientPools config.ClientPools) {
	path := string(ctx.Path())
	project, topic, err := parsePublishPath(path)
	if err != nil {
		ctx.Error(badPublishPath(path), fasthttp.StatusNotFound)
		return
	}
	client, ok := clientPools[project]
	if !ok {
		ctx.Error("Project not found: "+project, fasthttp.StatusNotFound)
		return
	}
	id, err := client.Publish(topic, ctx.PostBody())
	if err != nil {
		ctx.Error("Error: "+err.Error(), 400)
		return
	}
	fmt.Fprintf(ctx, "Success! Message ID: %q", id)
}

func router(clientPools config.ClientPools) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		if !strings.HasPrefix(path, "/publish") {
			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
			return
		}
		publishHandler(ctx, clientPools)
	}
}

// parsePublishPath requires a path of /publish/<project>/<topic>
// It returns an error if the path is malformed.
func parsePublishPath(path string) (project, topic string, err error) {
	parts := strings.Split(path, "/")
	if len(parts) != 4 || parts[1] != "publish" { // ["", "publish", <project>, <topic>]
		return "", "", errors.New(badPublishPath(path))
	}
	project = parts[2]
	topic = parts[3]
	return project, topic, nil
}

func badPublishPath(path string) string {
	return fmt.Sprintf("Path must be /publish/<project>/<topic>; got %s", path)
}

// StartPolling is a blocking call that uses a subscription's
// project name to find the appropriate client on which to receive messages.
// Concurrent-safe.
func StartPolling(clientPools config.ClientPools, s *config.Subscription) {
	client, ok := clientPools[s.Project]
	if !ok {
		panic(fmt.Errorf("Project %s not found", s.Project))
	}
	err := client.Subscription(s.Name).
		Receive(context.Background(), deliverAndAck)
	if err != nil {
		panic(err)
	}
}

// _post returns an error if it does not get a 200 response
func _post(host string, body []byte) error {
	req := fasthttp.AcquireRequest()
	req.SetBody(body)
	req.Header.SetMethod("POST")
	req.SetRequestURI(host)
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		return err
	}
	fasthttp.ReleaseRequest(req)

	if code := res.StatusCode(); code != 200 {
		return errors.New("non-200 status code: " + string(code))
	}
	fasthttp.ReleaseResponse(res)
	return nil
}

func main() {
	cfg := config.FromFile(os.Getenv("CONFIG_PATH"))
	clientPools := cfg.GetClientPools()
	for _, subscription := range cfg.GetSubscriptions() {
		go StartPolling(clientPools, subscription)
	}
	fasthttp.ListenAndServe(":30980", router(clientPools))
}
