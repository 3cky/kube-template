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

	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	FLAG_SERVER    = "server"
	FLAG_CONFIG    = "config"
	FLAG_POLL_TIME = "poll-time"
	FLAG_TEMPLATE  = "template"
)

func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kube-template",
		Short: "kube-template",
		Long: `Watches a series of templates on the file system, writing new changes when
  Kubernetes is updated.`,
		Run: runCmd,
	}
	initCmd(cmd)
	return cmd
}

func initCmd(cmd *cobra.Command) {
	cmd.Flags().StringP(FLAG_SERVER, "s", "", "The address and port of the Kubernetes API server")
	cmd.Flags().DurationP(FLAG_POLL_TIME, "p", 15*time.Second, "Kubernetes API server poll time")
	cmd.Flags().StringVarP(&cfgFile, FLAG_CONFIG, "c", "", fmt.Sprintf("config file (default is ./%s.(yaml|json))", CFG_FILE))
	cmd.Flags().StringSliceP(FLAG_TEMPLATE, "t", nil, `Adds a new template to watch on disk in the format
		'templatePath:outputPath[:command]'. This option is additive
		and may be specified multiple times for multiple templates.`)
}

func runCmd(cmd *cobra.Command, args []string) {
	config, err := newConfig(cmd)

	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}
	if len(config.TemplateDescriptors) == 0 {
		log.Fatalf("no templates to process found, exiting...")
	}

	// Start application
	app, err := newApp(config)
	if err != nil {
		log.Fatalf("can't create application: %v", err)
	}

	go app.Start()

	// Listen for signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// Event loop
EventLoop:
	for {
		select {
		case signal := <-signalCh:
			log.Printf("received %v signal, stopping", signal)
			app.Stop()
		case <-app.doneCh:
			break EventLoop
		}
	}
}
