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
	"errors"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	CFG_FILE        = "kube-template"
	CFG_MASTER      = FLAG_MASTER
	CFG_POLL_TIME   = FLAG_POLL_TIME
	CFG_POLL_PERIOD = FLAG_POLL_PERIOD
)

var cfgFile string

type Config struct {
	// Do not write template output
	DryRun bool
	// Run template processing once and exit
	RunOnce bool
	// Guess Kubernetes API settings from POD environment
	GuessKubeAPISettings bool
	// Kubernetes config file
	KubeConfig string
	// Kubernetes API server address
	Master string
	// Kubernetes API server poll period
	PollPeriod time.Duration

	// Template delimiters
	LeftDelimiter  string
	RightDelimiter string

	// Template paths
	templatePaths map[string]bool

	// Template descriptors
	TemplateDescriptors []*TemplateDescriptor
}

type TemplateDescriptor struct {
	// Template file path
	Path string
	// Template output path
	Output string
	// Optional command to execute after template output updating
	Command string
}

func readConfig(cmd *cobra.Command) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile) // specify config file set by flag
	} else {
		viper.SetConfigName(CFG_FILE) // default name of config file (without extension)
		viper.AddConfigPath(".")      // adding home directory as first search path
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.BindPFlag(CFG_MASTER, cmd.Flags().Lookup(FLAG_MASTER)); err != nil {
		return err
	}

	if err := viper.BindPFlag(CFG_POLL_PERIOD, cmd.Flags().Lookup(FLAG_POLL_PERIOD)); err != nil {
		return err
	}

	err := viper.ReadInConfig()

	if err == nil {
		glog.V(1).Infof("using config file: %s", viper.ConfigFileUsed())
		return nil
	}

	// raise error if config file used but can't be read
	if viper.ConfigFileUsed() != "" {
		return err
	}

	return nil
}

// Parses a string in format 'templatePath:outputPath[:command]' into a TemplateDescriptor struct
func parseTemplateDescriptor(s string) (*TemplateDescriptor, error) {
	if len(strings.TrimSpace(s)) == 0 {
		return nil, errors.New("empty template descriptor string")
	}

	var path, output, command string
	parts := strings.SplitN(s, ":", 3)

	switch len(parts) {
	case 2:
		path, output = parts[0], parts[1]
	case 3:
		path, output, command = parts[0], parts[1], parts[2]
	default:
		return nil, errors.New("invalid template descriptor, should be 'templatePath:outputPath[:command]'")
	}

	return &TemplateDescriptor{
		Path:    path,
		Output:  output,
		Command: command,
	}, nil
}

func newConfig(cmd *cobra.Command) (*Config, error) {
	// Create empty config
	config := new(Config)
	// Get command-line only options
	dryRun, err := cmd.Flags().GetBool(FLAG_DRY_RUN)
	if err != nil {
		return nil, err
	}
	config.DryRun = dryRun
	runOnce, err := cmd.Flags().GetBool(FLAG_RUN_ONCE)
	if err != nil {
		return nil, err
	}
	config.RunOnce = runOnce
	guessKubeAPISettings, err := cmd.Flags().GetBool(FLAG_GUESS_KUBE_API_SETTINGS)
	if err != nil {
		return nil, err
	}
	config.GuessKubeAPISettings = guessKubeAPISettings
	kubeConfig, err := cmd.Flags().GetString(FLAG_KUBE_CONFIG)
	if err != nil {
		return nil, err
	}
	config.KubeConfig = kubeConfig

	leftDelimiter, err := cmd.Flags().GetString(FLAG_LEFT_DELIM)
	if err != nil {
		return nil, err
	}
	config.LeftDelimiter = leftDelimiter

	rightDelimiter, err := cmd.Flags().GetString(FLAG_RIGHT_DELIM)
	if err != nil {
		return nil, err
	}
	config.RightDelimiter = rightDelimiter

	// Read config from file, if present
	err = readConfig(cmd)
	if err != nil {
		return nil, err
	}
	// Get command line / config options
	config.Master = viper.GetString(CFG_MASTER)
	if viper.IsSet(CFG_POLL_TIME) {
		config.PollPeriod = viper.GetDuration(CFG_POLL_TIME)
		glog.Warningf("'%s' parameter is deprecated, use '%s' instead", CFG_POLL_TIME, CFG_POLL_PERIOD)
	} else {
		config.PollPeriod = viper.GetDuration(CFG_POLL_PERIOD)
	}
	glog.V(2).Infof("poll period set to %v", config.PollPeriod)
	// Add template descriptors specified by command line
	cmdTemplates, err := cmd.Flags().GetStringSlice(FLAG_TEMPLATE)
	if err != nil {
		return nil, err
	}
	config.templatePaths = make(map[string]bool)
	config.TemplateDescriptors = make([]*TemplateDescriptor, 0)
	for _, template := range cmdTemplates {
		d, err := parseTemplateDescriptor(template)
		if err != nil {
			glog.Errorf("can't parse '%s': %v", template, err)
		} else {
			glog.V(2).Infof("adding template from command line: %s", d.Path)
			config.appendTemplateDescriptor(d)
		}
	}
	// Merge template descriptors from config file
	if iCfgTemplates := viper.Get("templates"); iCfgTemplates != nil {
		for _, iCfgTemplate := range iCfgTemplates.([]interface{}) {
			cfgTemplate := iCfgTemplate.(map[interface{}](interface{}))
			// Check template path and output path are present
			iPath, pathPresent := cfgTemplate["path"]
			iOutput, outputPresent := cfgTemplate["output"]
			if !pathPresent || !outputPresent {
				glog.Warningf("skipped non-complete template descriptor: %#v", cfgTemplate)
				continue
			}
			path, output := iPath.(string), iOutput.(string)
			// Command is optional
			var cmd string
			if iCmd, cmdPresent := cfgTemplate["command"]; cmdPresent {
				cmd = iCmd.(string)
			}
			// Add template descriptor
			d := &TemplateDescriptor{
				Path:    path,
				Output:  output,
				Command: cmd,
			}
			glog.V(2).Infof("adding template from config file: %s", d.Path)
			config.appendTemplateDescriptor(d)
		}
	}

	return config, nil
}

func (cfg *Config) appendTemplateDescriptor(d *TemplateDescriptor) {
	if _, added := cfg.templatePaths[d.Path]; !added {
		cfg.templatePaths[d.Path] = true
		cfg.TemplateDescriptors = append(cfg.TemplateDescriptors, d)
	} else {
		glog.Warningf("template already added: %s", d.Path)
	}
}
