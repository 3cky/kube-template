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
	// Template descriptor from configuration
	desc *TemplateDescriptor

	// Template name (base file name)
	name string

	// Template last output (in case of successfully rendered template)
	lastOutput string
}

func newTemplate(d *TemplateDescriptor) *Template {
	// Get last template output, if present
	o, err := ioutil.ReadFile(d.Output)
	if err != nil {
		o = nil
	}
	// Create template
	return &Template{
		desc:       d,
		name:       filepath.Base(d.Path),
		lastOutput: string(o),
	}
}

func (t *Template) Process(c *Client, dryRun bool) (bool, error) {
	if r, err := t.Render(c); err == nil {
		if changed := !(r == t.lastOutput); changed {
			// Template output changed
			if !dryRun {
				// TODO file mode from config
				// TODO atomic write
				if err := ioutil.WriteFile(t.desc.Output, []byte(r), 0644); err != nil {
					return false, err
				}
			}
			t.lastOutput = r
			return true, nil
		}
		// Template output not changed
		return false, nil
	} else {
		// Can't render template
		return false, err
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

// {{pods "selector" "namespace"}}
func pods(c *Client) func(...string) ([]api.Pod, error) {
	return func(s ...string) ([]api.Pod, error) {
		namespace, selector := DEFAULT_NAMESPACE, DEFAULT_SELECTOR
		switch len(s) {
		case 0:
			break
		case 1:
			selector = s[0]
		case 2:
			selector = s[0]
			namespace = s[1]
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

// {{services "selector" "namespace"}}
func services(c *Client) func(...string) ([]api.Service, error) {
	return func(s ...string) ([]api.Service, error) {
		namespace, selector := DEFAULT_NAMESPACE, DEFAULT_SELECTOR
		switch len(s) {
		case 0:
			break
		case 1:
			selector = s[0]
		case 2:
			selector = s[0]
			namespace = s[1]
		default:
			return nil, fmt.Errorf("expected max 2 arguments, got %d", len(s))
		}
		services, err := c.Services(namespace, selector)
		if err != nil {
			return nil, err
		}
		return services, nil
	}
}
