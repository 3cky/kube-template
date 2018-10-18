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
	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	kubeClient kubernetes.Interface
}

func newClient(cfg *Config) (*Client, error) {
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

	return &Client{
		kubeClient: c,
	}, nil
}

func (c *Client) Pods(namespace, selector string) ([]corev1.Pod, error) {
	glog.V(4).Infof("fetching pods, namespace: %q, selector: %q", namespace, selector)
	options := metav1.ListOptions{LabelSelector: selector}
	podList, err := c.kubeClient.CoreV1().Pods(namespace).List(options)
	if err != nil {
		return nil, err
	}
	return podList.Items, nil
}

func (c *Client) Services(namespace, selector string) ([]corev1.Service, error) {
	glog.V(4).Infof("fetching services, namespace: %q, selector: %q", namespace, selector)
	options := metav1.ListOptions{LabelSelector: selector}
	svcList, err := c.kubeClient.CoreV1().Services(namespace).List(options)
	if err != nil {
		return nil, err
	}
	return svcList.Items, nil
}

func (c *Client) ReplicationControllers(namespace, selector string) ([]corev1.ReplicationController, error) {
	glog.V(4).Infof("fetching replication controllers, namespace: %q, selector: %q", namespace, selector)
	options := metav1.ListOptions{LabelSelector: selector}
	rcList, err := c.kubeClient.CoreV1().ReplicationControllers(namespace).List(options)
	if err != nil {
		return nil, err
	}
	return rcList.Items, nil
}

func (c *Client) Events(namespace, selector string) ([]corev1.Event, error) {
	glog.V(4).Infof("fetching events, namespace: %q, selector: %q", namespace, selector)
	options := metav1.ListOptions{LabelSelector: selector}
	evList, err := c.kubeClient.CoreV1().Events(namespace).List(options)
	if err != nil {
		return nil, err
	}
	return evList.Items, nil
}

func (c *Client) Endpoints(namespace, selector string) ([]corev1.Endpoints, error) {
	glog.V(4).Infof("fetching endpoints, namespace: %q, selector: %q", namespace, selector)
	options := metav1.ListOptions{LabelSelector: selector}
	epList, err := c.kubeClient.CoreV1().Endpoints(namespace).List(options)
	if err != nil {
		return nil, err
	}
	return epList.Items, nil
}

func (c *Client) Nodes(selector string) ([]corev1.Node, error) {
	glog.V(4).Infof("fetching nodes, selector: %q", selector)
	options := metav1.ListOptions{LabelSelector: selector}
	nodeList, err := c.kubeClient.CoreV1().Nodes().List(options)
	if err != nil {
		return nil, err
	}
	return nodeList.Items, nil
}

func (c *Client) Namespaces(selector string) ([]corev1.Namespace, error) {
	glog.V(4).Infof("fetching namespaces, selector: %q", selector)
	options := metav1.ListOptions{LabelSelector: selector}
	nsList, err := c.kubeClient.CoreV1().Namespaces().List(options)
	if err != nil {
		return nil, err
	}
	return nsList.Items, nil
}
