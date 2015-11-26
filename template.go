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
	"path/filepath"
	"text/template"
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

func (t *Template) Render() (string, error) {
	// Read template data
	data, err := ioutil.ReadFile(t.desc.Path)
	if err != nil {
		return "", err
	}
	s := string(data)
	// Create template from read data
	template, err := template.New(t.name).Parse(s)
	if err != nil {
		return "", fmt.Errorf("%s: %v", t.name, err)
	}
	// Render template to buffer
	buf := new(bytes.Buffer)
	if err := template.Execute(buf, nil); err != nil {
		return "", fmt.Errorf("%s: %v", t.name, err)
	}

	return buf.String(), nil
}
