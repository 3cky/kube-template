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
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/kubernetes/pkg/util"
)

const (
	FLAG_V             = "v"
	FLAG_STDERR_THRESH = "stderrthreshold"
	FLAG_RUN_ONCE      = "once"
	FLAG_DRY_RUN       = "dry-run"
	FLAG_SERVER        = "server"
	FLAG_CONFIG        = "config"
	FLAG_POLL_TIME     = "poll-time"
	FLAG_TEMPLATE      = "template"
)

func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kube-template",
		Short: "kube-template",
		Long:  "Watches Kubernetes for updates, writing output of a series of templates to files.",
		Run:   runCmd,
	}
	initCmd(cmd)
	return cmd
}

func initCmd(cmd *cobra.Command) {
	// Command-related flags set
	f := cmd.Flags()
	f.Bool(FLAG_DRY_RUN, false, "don't write template output, dump result to stdout")
	f.Bool(FLAG_RUN_ONCE, false, "run template processing once and exit")
	f.StringP(FLAG_SERVER, "s", "", "the address and port of the Kubernetes API server")
	f.DurationP(FLAG_POLL_TIME, "p", 15*time.Second, "Kubernetes API server poll time")
	f.StringVarP(&cfgFile, FLAG_CONFIG, "c", "", fmt.Sprintf("config file (default is ./%s.(yaml|json))", CFG_FILE))
	f.StringSliceP(FLAG_TEMPLATE, "t", nil, `adds a new template to watch on disk in the format
		'templatePath:outputPath[:command]'. This option is additive
		and may be specified multiple times for multiple templates.`)
	// Merge glog-related flags
	// FIXME probably we shouldn't use k8s utils there
	pflag.CommandLine.AddFlagSet(f)
	util.InitFlags()
	util.InitLogs()
	defer util.FlushLogs()
}

func runCmd(cmd *cobra.Command, args []string) {
	config, err := newConfig(cmd)

	if err != nil {
		glog.Fatalf("configuration error: %v", err)
	}
	if len(config.TemplateDescriptors) == 0 {
		glog.Fatalf("no templates to process (use --help to get configuration options), exiting...")
	}

	// Start application
	app, err := newApp(config)
	if err != nil {
		glog.Fatalf("can't create application: %v", err)
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
			glog.V(2).Infof("received %v signal, stopping", signal)
			app.Stop()
		case <-app.doneCh:
			break EventLoop
		}
	}
}
