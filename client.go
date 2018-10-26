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
	"sync"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/exec"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	_ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	DEFAULT_MASTER_HOST = "http://127.0.0.1:8080/"
)

type Client struct {
	sync.RWMutex
	kubeClient      kubernetes.Interface
	stopCh          chan struct{}
	useInformers    bool
	informerFactory informers.SharedInformerFactory
	listers         map[string]interface{}
}

func newClientForConfig(cfg *Config, stopCh chan struct{}) (*Client, error) {
	var config *rest.Config
	if cfg.GuessKubeAPISettings {
		var err error
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else if cfg.KubeConfig != "" {
		var err error
		config, err = clientcmd.BuildConfigFromFlags("", cfg.KubeConfig)
		if err != nil {
			return nil, err
		}
	} else {
		host := DEFAULT_MASTER_HOST
		if cfg.Master != "" {
			host = cfg.Master
		}
		config = &rest.Config{
			Host: host,
		}
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return newClient(c, stopCh, cfg.PollingEnabled())
}

func newClient(c kubernetes.Interface, stopCh chan struct{}, useInformers bool) (*Client, error) {
	return &Client{
		kubeClient:      c,
		stopCh:          stopCh,
		useInformers:    useInformers,
		informerFactory: informers.NewSharedInformerFactory(c, 0),
		listers:         make(map[string]interface{}),
	}, nil
}
