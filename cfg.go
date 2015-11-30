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

	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var cfgFile string

type Config struct {
	Server string

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

func readConfig(cmd *cobra.Command) {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("kube-template") // name of config file (without extension)
	viper.AddConfigPath(".")             // adding home directory as first search path
	viper.AutomaticEnv()                 // read in environment variables that match

	viper.BindPFlag("server", cmd.Flags().Lookup("server"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("using config file: %s", viper.ConfigFileUsed())
	}
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
	// Read config from file, if present
	readConfig(cmd)
	config.Server = viper.GetString("server")
	// Add template descriptors specified by command line
	cmdTemplates, err := cmd.Flags().GetStringSlice("template")
	if err != nil {
		return nil, err
	}
	config.templatePaths = make(map[string]bool)
	config.TemplateDescriptors = make([]*TemplateDescriptor, 0)
	for _, template := range cmdTemplates {
		d, err := parseTemplateDescriptor(template)
		if err != nil {
			log.Printf("can't parse '%s': %v", template, err)
		} else {
			log.Printf("adding template from command line: %s", d.Path)
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
				log.Printf("skipped non-complete template descriptor: %#v", cfgTemplate)
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
			log.Printf("adding template from config file: %s", d.Path)
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
		log.Printf("template already added: %s", d.Path)
	}
}