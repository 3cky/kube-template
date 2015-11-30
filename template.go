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
	"bytes"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"fmt"
	"k8s.io/kubernetes/pkg/api"
)

const (
	DEFAULT_NAMESPACE = api.NamespaceDefault
	DEFAULT_SELECTOR  = ""
)

type Template struct {
	desc *TemplateDescriptor
	name string
}

func newTemplate(d *TemplateDescriptor) *Template {
	return &Template{
		desc: d,
		name: filepath.Base(d.Path),
	}
}

func (t *Template) Render(c *Client) (string, error) {
	// Read template data
	data, err := ioutil.ReadFile(t.desc.Path)
	if err != nil {
		return "", err
	}
	s := string(data)
	// Create template from read data
	template, err := template.New(t.name).Funcs(funcMap(c)).Parse(s)
	if err != nil {
		return "", err
	}
	// Render template to buffer
	buf := new(bytes.Buffer)
	if err := template.Execute(buf, nil); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func funcMap(c *Client) template.FuncMap {
	return template.FuncMap{
		"pods":     pods(c),
		"services": services(c),
	}
}

// {{pods "namespace" "selector"}}
func pods(c *Client) func(...string) ([]api.Pod, error) {
	return func(s ...string) ([]api.Pod, error) {
		namespace := DEFAULT_NAMESPACE
		selector := DEFAULT_SELECTOR
		switch len(s) {
		case 0:
			break
		case 1:
			namespace = s[0]
		case 2:
			namespace = s[0]
			selector = s[1]
		default:
			return nil, fmt.Errorf("expected max 2 arguments, got %d", len(s))
		}
		pods, err := c.Pods(namespace, selector)
		if err != nil {
			return nil, err
		}
		return pods, nil
	}
}

// {{services "namespace" "selector"}}
func services(c *Client) func(...string) ([]api.Service, error) {
	return func(s ...string) ([]api.Service, error) {
		namespace := DEFAULT_NAMESPACE
		selector := DEFAULT_SELECTOR
		switch len(s) {
		case 0:
			break
		case 1:
			namespace = s[0]
		case 2:
			namespace = s[0]
			selector = s[1]
		default:
			return nil, fmt.Errorf("expected max 2 arguments, got %d", len(s))
		}
		svcs, err := c.Services(namespace, selector)
		if err != nil {
			return nil, err
		}
		return svcs, nil
	}
}
