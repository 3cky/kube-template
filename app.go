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

//go:generate go run cmd/genclient/main.go
//go:generate go run cmd/gendeps/main.go
//go:generate go run cmd/gentemplate/main.go

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
)

type App struct {
	// Stop channel
	stopCh chan struct{}
	// Done channel
	doneCh chan struct{}

	// Do not write template output flag
	dryRun bool

	// Template output update period
	updatePeriod time.Duration

	// Dependency manager
	dm *DependencyManager

	// Templates to process
	templates []*Template
}

func newApp(cfg *Config) (*App, error) {

	// Create Kubernetes client
	stopCh := make(chan struct{})
	client, err := newClientForConfig(cfg, stopCh)
	if err != nil {
		close(stopCh)
		return nil, err
	}

	// Create dependency manager
	dm := newDependencyManager(client)

	// Add all configured templates
	templates, err := newTemplatesFromConfig(cfg, dm)
	if err != nil {
		return nil, err
	}

	doneCh := make(chan struct{})

	return &App{
		stopCh:       stopCh,
		doneCh:       doneCh,
		dm:           dm,
		templates:    templates,
		dryRun:       cfg.DryRun,
		updatePeriod: cfg.PollPeriod,
	}, nil
}

func (app *App) Start() {
	glog.V(1).Infoln("starting templates processing...")

	defer close(app.doneCh)

	defer glog.V(1).Infoln("templates processing stopped")

	// Initial templates processing run
	app.Run()

	if app.updatePeriod.Nanoseconds() <= 0 {
		<-app.stopCh
		return
	}

	for {
		select {
		case <-app.stopCh:
			return
		case <-time.After(app.updatePeriod):
			app.Run()
		}
	}
}

func (app *App) RunOnce() {
	glog.V(1).Infoln("run once templates processing...")
	app.Run()
	glog.V(1).Infoln("templates processed")
}

func (app *App) Run() {
	// Commands to execute are stored in list instead of map to ensure correct execution order
	var commands []string
	commandTimeouts := make(map[string]time.Duration)
	// Flush cached dependencies
	app.dm.flushCachedDependencies()
	// Process templates
	for _, t := range app.templates {
		glog.V(2).Infof("processing template: %s", t.name)
		if updated, err := t.Process(app.dryRun); err == nil {
			if updated {
				if !app.dryRun {
					glog.V(2).Infof("template output updated: %s", t.name)
				} else {
					fmt.Printf("(dry-run) %s:\n%s", t.name, t.lastOutput)
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
						glog.V(4).Infof("template %s: scheduled command: %q", t.name, cmd)
						commands = append(commands, cmd)
						commandTimeouts[cmd] = t.desc.CommandTimeout
					} else {
						glog.V(4).Infof("template %s: command already scheduled: %q", t.name, cmd)
					}
				}
			} else {
				glog.V(2).Infof("template output not changed: %s", t.name)
			}
		} else {
			glog.Errorf("can't render %v", err)
		}
	}
	// Execute commands for templates
	for _, cmd := range commands {
		if !app.dryRun {
			glog.V(4).Infof("executing: %q", cmd)
			if err := Execute(cmd, commandTimeouts[cmd]); err == nil {
				glog.V(4).Infof("executed: %q", cmd)
			} else {
				glog.Errorf("command %q: %v", cmd, err)
			}
		} else {
			fmt.Printf("(dry-run) executing: %q\n", cmd)
		}
	}
}

func (app *App) Stop() {
	glog.V(1).Infoln("stopping templates processing...")
	close(app.stopCh)
}
