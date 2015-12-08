// Copyright © 2015 Victor Antonovich <victor@antonovich.me>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"

	"k8s.io/kubernetes/pkg/api"
	kubeClient "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
)

const (
	DEFAULT_HOST = "http://localhost:8080"
)

type ClientInterface interface {
	Pods(namespace string, selector string) ([]api.Pod, error)
}

type Client struct {
	kubeClient *kubeClient.Client
}

func newClient(cfg *Config) (*Client, error) {
	host := DEFAULT_HOST
	if cfg.Server != "" {
		host = cfg.Server
	}

	config := kubeClient.Config{
		Host: host,
	}

	c, err := kubeClient.New(&config)
	if err != nil {
		return nil, err
	}

	return &Client{
		kubeClient: c,
	}, nil
}

func (c *Client) Pods(namespace, selector string) ([]api.Pod, error) {
	log.Printf("fetching pods, namespace: %q, selector: %q", namespace, selector)
	s, err := labels.Parse(selector)
	if err != nil {
		return nil, err
	}
	podList, err := c.kubeClient.Pods(namespace).List(s, nil)
	if err != nil {
		return nil, err
	}
	return podList.Items, nil
}

func (c *Client) Services(namespace, selector string) ([]api.Service, error) {
	log.Printf("fetching services, namespace: %q, selector: %q", namespace, selector)
	s, err := labels.Parse(selector)
	if err != nil {
		return nil, err
	}
	svcList, err := c.kubeClient.Services(namespace).List(s, nil)
	if err != nil {
		return nil, err
	}
	return svcList.Items, nil
}
