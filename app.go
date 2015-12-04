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
		case <-time.After(app.config.PollTime):
			app.Run()
		}
	}
}

func (app *App) Run() {
	// Commands to execute are stored in list instead of map to ensure correct execution order
	var commands []string
	// Process templates
	for _, t := range app.templates {
		log.Printf("process template: %s", t.name)
		if updated, err := t.Process(app.client); err == nil {
			if updated {
				log.Printf("template output changed: %s\n %s", t.name, t.lastOutput)
				if cmd := t.desc.Command; len(cmd) > 0 {
					// Check template command is already in list of commands to execute
					if c, err := NormPath(cmd); err == nil {
						if !IsPresent(commands, c) {
							log.Printf("template %s: scheduled command: %q", t.name, c)
							commands = append(commands, c)
						} else {
							log.Printf("template %s: command already scheduled: %q", t.name, c)
						}
					} else {
						log.Printf("template %s: can't schedule command: %v", t.name, err)
					}
				}
			} else {
				log.Printf("template output not changed: %s", t.name)
			}
		} else {
			log.Printf("can't render %v", err)
		}
	}
	// Execute commands for templates
	for _, cmd := range commands {
		log.Printf("executing: %q", cmd)
		if err := Execute(cmd, time.Second); err == nil {
			log.Printf("executed: %q", cmd)
		} else {
			log.Printf("command %q: %v", cmd, err)
		}
	}
}

func (app *App) Stop() {
	log.Println("stopping application...")
	close(app.stopCh)
}
