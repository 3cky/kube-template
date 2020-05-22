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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	gotemplate "text/template"

	"strings"

	"github.com/Masterminds/sprig"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DefaultNamespace = metav1.NamespaceDefault
	DefaultSelector  = ""
)

type Template struct {
	// Template descriptor from configuration
	desc *TemplateDescriptor

	// Template name (base file name)
	name string

	// Go template to render
	template *gotemplate.Template

	// Template last output (in case of successfully rendered template)
	lastOutput string
}

func newTemplate(cfg *Config, dm *DependencyManager, d *TemplateDescriptor) (*Template, error) {
	// Template name
	name := filepath.Base(d.Path)
	// Get last template output, if present
	o, err := ioutil.ReadFile(d.Output)
	if err != nil {
		o = nil
	}
	// Read template data
	data, err := ioutil.ReadFile(d.Path)
	if err != nil {
		return nil, err
	}
	s := string(data)
	// Create Go template from read data
	template, err := gotemplate.New(name).Delims(cfg.LeftDelimiter, cfg.RightDelimiter).Funcs(funcMap(dm)).Parse(s)
	if err != nil {
		return nil, err
	}
	// Create template
	return &Template{
		desc:       d,
		name:       name,
		template:   template,
		lastOutput: string(o),
	}, nil
}

func newTemplatesFromConfig(cfg *Config, dm *DependencyManager) ([]*Template, error) {
	templates := make([]*Template, 0, len(cfg.TemplateDescriptors))
	for _, d := range cfg.TemplateDescriptors {
		t, err := newTemplate(cfg, dm, d)
		if err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (t *Template) Process(dryRun bool) (bool, error) {
	if r, err := t.Render(); err == nil {
		if changed := !(r == t.lastOutput); changed {
			// Template output changed
			if !dryRun {
				if err := t.Write([]byte(r)); err != nil {
					// Can't write template output
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

func (t *Template) Write(content []byte) error {
	dir := filepath.Dir(t.desc.Output)
	if _, err := os.Stat(t.desc.Output); os.IsNotExist(err) {
		// Output file does not exist, create intermediate dirs and write directly
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		// TODO file mode from config
		if err := ioutil.WriteFile(t.desc.Output, content, 0644); err != nil {
			return err
		}
	} else {
		// Output file exist, update atomically using temp file
		var f *os.File
		if f, err = ioutil.TempFile(dir, t.name); err != nil {
			return err
		}
		defer UnlinkQuietly(f.Name())
		// Write template output to temp file
		if _, err := f.Write(content); err != nil {
			return err
		}
		if err := f.Sync(); err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
		// TODO file mode from config
		if err := os.Chmod(f.Name(), 0644); err != nil {
			return err
		}
		// Rename temp file to output file
		if err := os.Rename(f.Name(), t.desc.Output); err != nil {
			return err
		}
	}
	return nil
}

func (t *Template) Render() (string, error) {
	// Render template to buffer
	buf := new(bytes.Buffer)
	if err := t.template.Execute(buf, nil); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func funcMap(dm *DependencyManager) gotemplate.FuncMap {
	f := gotemplate.FuncMap{
		// Legacy helper functions
		"toLower":   strings.ToLower,
		"toUpper":   strings.ToUpper,
		"toTitle":   strings.Title,
		"trimSpace": strings.TrimSpace,
	}

	// Kubernetes objects functions
	for k, v := range kubeObjectsFuncMap(dm) {
		f[k] = v
	}

	// Sprig helper functions
	for k, v := range sprig.FuncMap() {
		f[k] = v
	}

	return f
}

// Parse template tag with max 1 argument - selector
func parseSelector(s ...string) (string, error) {
	selector := DefaultSelector
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
	namespace, selector := DefaultNamespace, DefaultSelector
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
