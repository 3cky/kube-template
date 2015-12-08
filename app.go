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
	"os"
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

	// Initial templates processing run
	app.Run()

	if app.config.RunOnce {
		log.Println("run once requested, exiting...")
		close(app.doneCh)
		return
	}

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
		if updated, err := t.Process(app.client, app.config.DryRun); err == nil {
			if updated {
				if !app.config.DryRun {
					log.Printf("template output updated: %s", t.name)
				} else {
					log.Printf("(dry-run) %s:\n%s", t.name, t.lastOutput)
				}
				if cmd := t.desc.Command; len(cmd) > 0 {
					// Normalize command path, if applicable
					if _, err := os.Stat(cmd); err == nil {
						if c, err := NormPath(cmd); err == nil {
							cmd = c
						}
					}
					// Check template command is already in list of commands to execute
					if !IsPresent(commands, cmd) {
						log.Printf("template %s: scheduled command: %q", t.name, cmd)
						commands = append(commands, cmd)
					} else {
						log.Printf("template %s: command already scheduled: %q", t.name, cmd)
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
		if !app.config.DryRun {
			log.Printf("executing: %q", cmd)
			if err := Execute(cmd, time.Second); err == nil {
				log.Printf("executed: %q", cmd)
			} else {
				log.Printf("command %q: %v", cmd, err)
			}
		} else {
			log.Printf("(dry-run) executing: %q", cmd)
		}
	}
}

func (app *App) Stop() {
	log.Println("stopping application...")
	close(app.stopCh)
}
