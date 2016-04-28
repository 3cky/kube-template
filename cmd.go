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
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"bytes"
	"k8s.io/kubernetes/pkg/util"
)

const (
	FLAG_RUN_ONCE  = "once"
	FLAG_DRY_RUN   = "dry-run"
	FLAG_MASTER    = "master"
	FLAG_CONFIG    = "config"
	FLAG_POLL_TIME = "poll-time"
	FLAG_TEMPLATE  = "template"
	FLAG_HELP_MD   = "help-md"
)

func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "kube-template",
		Long: "Watches Kubernetes for updates, writing output of a series of templates to files.",
		Run:  runCmd,
	}
	initCmd(cmd)
	return cmd
}

func initCmd(cmd *cobra.Command) {
	// Command-related flags set
	f := cmd.Flags()
	f.Bool(FLAG_DRY_RUN, false, "don't write template output, dump result to stdout")
	f.Bool(FLAG_RUN_ONCE, false, "run template processing once and exit")
	f.String(FLAG_MASTER, "", fmt.Sprintf("Kubernetes API server address (default is %s)", DEFAULT_MASTER_HOST))
	f.DurationP(FLAG_POLL_TIME, "p", 15*time.Second, "Kubernetes API server poll time")
	f.StringVarP(&cfgFile, FLAG_CONFIG, "c", "", fmt.Sprintf("config file (default is ./%s.(yaml|json))", CFG_FILE))
	f.StringSliceP(FLAG_TEMPLATE, "t", nil, `adds a new template to watch on disk in the format
		'templatePath:outputPath[:command]'. This option is additive
		and may be specified multiple times for multiple templates`)
	f.Bool(FLAG_HELP_MD, false, "get help in Markdown format")
	// Merge flags
	pflag.CommandLine.SetNormalizeFunc(func(_ *pflag.FlagSet, name string) pflag.NormalizedName {
		if strings.Contains(name, "_") {
			return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
		}
		return pflag.NormalizedName(name)
	})
	pflag.CommandLine.AddFlagSet(f)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	// Init logs
	// FIXME probably we shouldn't use k8s utils there
	util.InitLogs()
	defer util.FlushLogs()
}

func runCmd(cmd *cobra.Command, _ []string) {
	if f, _ := cmd.Flags().GetBool(FLAG_HELP_MD); f {
		out := new(bytes.Buffer)
		cobra.GenMarkdown(cmd, out)
		fmt.Println(out)
		return
	}

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
