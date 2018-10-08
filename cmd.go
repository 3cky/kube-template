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
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
)

const (
	FLAG_VERSION                 = "version"
	FLAG_RUN_ONCE                = "once"
	FLAG_DRY_RUN                 = "dry-run"
	FLAG_MASTER                  = "master"
	FLAG_CONFIG                  = "config"
	FLAG_POLL_TIME               = "poll-time"
	FLAG_TEMPLATE                = "template"
	FLAG_HELP_MD                 = "help-md"
	FLAG_GUESS_KUBE_API_SETTINGS = "guess-kube-api-settings"
	FLAG_KUBE_CONFIG             = "kube-config"
	FLAG_LEFT_DELIM              = "left-delimiter"
	FLAG_RIGHT_DELIM             = "right-delimiter"
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
	f.Bool(FLAG_VERSION, false, "display the version number and build timestamp")
	f.Bool(FLAG_DRY_RUN, false, "don't write template output, dump result to stdout")
	f.Bool(FLAG_RUN_ONCE, false, "run template processing once and exit")
	f.Bool(FLAG_GUESS_KUBE_API_SETTINGS, false, "guess Kubernetes API settings from POD environment")
	f.String(FLAG_MASTER, "", fmt.Sprintf("Kubernetes API server address (default is %s)", DEFAULT_MASTER_HOST))
	f.DurationP(FLAG_POLL_TIME, "p", 15*time.Second, "Kubernetes API server poll time (0 disables server polling)")
	f.StringP(FLAG_KUBE_CONFIG, "k", "", "Kubernetes config file to use")
	f.StringP(FLAG_LEFT_DELIM, "l", "{{", "templating left delimiter (default is {{)")
	f.StringP(FLAG_RIGHT_DELIM, "r", "}}", "templating right delimiter (default is }})")
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
	// Init logging
	initLogs()
	defer flushLogs()
}

func runCmd(cmd *cobra.Command, _ []string) {
	if f, _ := cmd.Flags().GetBool(FLAG_VERSION); f {
		fmt.Printf("Build version: %s\n", BuildVersion)
		fmt.Printf("Build timestamp: %s\n", BuildTimestamp)
		return
	}

	if f, _ := cmd.Flags().GetBool(FLAG_HELP_MD); f {
		out := new(bytes.Buffer)
		doc.GenMarkdown(cmd, out)
		fmt.Println(out)
		return
	}

	getConfig := func() (*Config, error) {
		config, err := newConfig(cmd)
		if err != nil {
			return nil, err
		}
		if len(config.TemplateDescriptors) == 0 {
			return nil, errors.New("no templates to process")
		}
		return config, nil
	}

	config, err := getConfig()
	if err != nil {
		glog.Fatalf("config error: %v, exiting...", err)
	}

	app, err := newApp(config)
	if err != nil {
		glog.Fatalf("config couldn't be used: %v", err)
	}

	if config.RunOnce {
		app.RunOnce()
		return
	}

	// Start templates processing
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
		case sig := <-signalCh:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				glog.V(2).Infof("received %v signal...", sig)
				// Stop templates processing and exit
				app.Stop()
				select {
				case <-app.doneCh:
					break EventLoop
				}
			case syscall.SIGHUP:
				glog.V(2).Infof("received %v signal, reloading config", sig)
				config, err := getConfig()
				if err != nil {
					glog.Errorf("config reloading error: %v", err)
					continue
				}
				newApp, err := newApp(config)
				if err != nil {
					glog.Errorf("reloaded config couldn't be used: %v", err)
					continue
				}
				// Stop templates processing using current config
				app.Stop()
				select {
				case <-app.doneCh:
					// Start templates processing using new config
					app = newApp
					go app.Start()

				}
			}
		}
	}
}
