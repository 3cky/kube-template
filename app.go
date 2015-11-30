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
	"time"
)

type App struct {
	// Stop channel
	stopCh chan interface{}
	// Done channel
	doneCh chan interface{}

	// Application config
	config *Config

	// Kubernetes client
	client *Client

	// Templates to process
	templates []*Template
}

func newApp(cfg *Config) (*App, error) {
	app := &App{
		config: cfg,
	}

	app.stopCh = make(chan interface{})
	app.doneCh = make(chan interface{})

	// Add all configured templates
	app.templates = make([]*Template, 0, len(cfg.TemplateDescriptors))
	for _, d := range cfg.TemplateDescriptors {
		app.templates = append(app.templates, newTemplate(d))
	}

	// Create Kubernetes client
	c, err := newClient(cfg)
	if err != nil {
		return nil, err
	}
	app.client = c

	return app, nil
}

func (app *App) Start() {
	log.Println("starting application...")

	defer log.Println("application stopped")

	app.Run()

MainLoop:
	for {
		select {
		case <-app.stopCh:
			close(app.doneCh)
			break MainLoop
		case <-time.After(5 * time.Second):
			app.Run()
		}
	}
}

func (app *App) Run() {
	for _, t := range app.templates {
		log.Printf("render: %s", t.name)
		if rendered, err := t.Render(app.client); err == nil {
			log.Printf("rendered %s:\n %s", t.name, rendered)
		} else {
			log.Printf("can't render %v", err)
		}
	}
}

func (app *App) Stop() {
	log.Println("stopping application...")
	close(app.stopCh)
}
