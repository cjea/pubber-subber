package config

import (
	"encoding/json"
	"io/ioutil"

	pubsub "github.com/cjea/pubber-subber/src/pubsub/gcp"
)

// Project has a 1:1 relationship with connection pool clients
type Project struct {
	Name          string   `json:"name"`
	Subscriptions []string `json:"subscriptions"`
	Topics        []string `json:"topics"`
}

// Config is a list of projects, each containing topics and subscriptions
type Config struct {
	Projects []*Project `json:"projects"`
}

// ClientPools maps a project name to a client
type ClientPools = map[string]*pubsub.Client

// GetClientPools maps a project name to a client
func (cfg *Config) GetClientPools() ClientPools {
	ret := make(map[string]*pubsub.Client)
	for _, project := range cfg.Projects {
		p, err := pubsub.NewClient(project.Name)
		if err != nil {
			panic(err)
		}
		ret[project.Name] = p
	}
	return ret
}

// Subscription represents a subscription to a topic
type Subscription struct {
	Name    string
	Project string
}

// GetSubscriptions maps a config to a list of Subscriptions
func (cfg *Config) GetSubscriptions() []*Subscription {
	ret := []*Subscription{}
	for _, project := range cfg.Projects {
		for _, sub := range project.Subscriptions {
			ret = append(ret, &Subscription{Name: sub, Project: project.Name})
		}
	}
	return ret
}

// FromFile takes a path to a JSON config and returns a structured *Config
func FromFile(path string) *Config {
	cfgBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic("CONFIG_PATH must point to a valid JSON config:\n" + err.Error())
	}
	var cfg Config
	err = json.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		panic("CONFIG_PATH is not valid json:\n" + err.Error())
	}
	return &cfg
}
