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

func (t *Template) Process(dm *DependencyManager, dryRun bool) (bool, error) {
	if r, err := t.Render(dm); err == nil {
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

func (t *Template) Render(dm *DependencyManager) (string, error) {
	// Read template data
	data, err := ioutil.ReadFile(t.desc.Path)
	if err != nil {
		return "", err
	}
	s := string(data)
	// Create template from read data
	template, err := template.New(t.name).Funcs(funcMap(dm)).Parse(s)
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

func funcMap(dm *DependencyManager) template.FuncMap {
	return template.FuncMap{
		// Kubernetes objects
		"pods":                   pods(dm),
		"services":               services(dm),
		"replicationcontrollers": replicationcontrollers(dm),
		"events":                 events(dm),
		"endpoints":              endpoints(dm),
		"nodes":                  nodes(dm),
		"namespaces":             namespaces(dm),
		// Utils
		"add": add,
		"sub": sub,
	}
}

// Parse template tag with max 1 argument - selector
func parseSelector(s ...string) (string, error) {
	selector := DEFAULT_SELECTOR
	switch len(s) {
	case 0:
		break
	case 1:
		selector = s[0]
	default:
		return "", fmt.Errorf("expected max 1 argument, got %d", len(s))
	}
	return selector, nil
}

// Parse template tag with max 2 arguments - selector and namespace (in given order)
func parseNamespaceSelector(s ...string) (string, string, error) {
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
		return "", "", fmt.Errorf("expected max 2 arguments, got %d", len(s))
	}
	return namespace, selector, nil
}

// {{pods "selector" "namespace"}}
func pods(dm *DependencyManager) func(...string) ([]api.Pod, error) {
	return func(s ...string) ([]api.Pod, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.Pods(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{services "selector" "namespace"}}
func services(dm *DependencyManager) func(...string) ([]api.Service, error) {
	return func(s ...string) ([]api.Service, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.Services(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{replicationcontrollers "selector" "namespace"}}
func replicationcontrollers(dm *DependencyManager) func(...string) ([]api.ReplicationController, error) {
	return func(s ...string) ([]api.ReplicationController, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.ReplicationControllers(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{events "selector" "namespace"}}
func events(dm *DependencyManager) func(...string) ([]api.Event, error) {
	return func(s ...string) ([]api.Event, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.Events(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{endpoints "selector" "namespace"}}
func endpoints(dm *DependencyManager) func(...string) ([]api.Endpoints, error) {
	return func(s ...string) ([]api.Endpoints, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.Endpoints(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{nodes "selector"}}
func nodes(dm *DependencyManager) func(...string) ([]api.Node, error) {
	return func(s ...string) ([]api.Node, error) {
		if selector, err := parseSelector(s...); err == nil {
			return dm.Nodes(selector)
		} else {
			return nil, err
		}
	}
}

// {{namespaces "selector"}}
func namespaces(dm *DependencyManager) func(...string) ([]api.Namespace, error) {
	return func(s ...string) ([]api.Namespace, error) {
		if selector, err := parseSelector(s...); err == nil {
			return dm.Namespaces(selector)
		} else {
			return nil, err
		}
	}
}

// {{add a b}}
func add(a, b int) int {
	return a + b
}

// {{sub a b}}
func sub(a, b int) int {
	return a - b
}
