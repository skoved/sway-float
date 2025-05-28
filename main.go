// Copyright skoved
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

const configDirSuffix = "/sway-float/config.yaml"

type floatConfig struct {
	appId string `yaml:"app_id"`
	mark  string `yaml:"con_mark"`
	title string `yaml:"title"`
}

// Prints err to stderr and calls os.Exit with statusCode
func errorExit(err error, statusCode int) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(statusCode)
}

func main() {
	var conf floatConfig
	configDir, confDirErr := os.UserConfigDir()
	if confDirErr != nil {
		errorExit(fmt.Errorf("could not determine the current user's config dir: %w", confDirErr), 1)
	}
	configPath := configDir + "/" + configDirSuffix
	configBytes, fileErr := os.ReadFile(configPath) //gosec:disable G304 -- need to build the path to the config dir
	if fileErr != nil {
		errorExit(fmt.Errorf("could not read file %s: %w", configPath, fileErr), 1)
	}
	yamlErr := yaml.Unmarshal(configBytes, &conf)
	if yamlErr != nil {
		errorExit(fmt.Errorf("could not parse yaml in file %s: %w", configPath, yamlErr), 1)
	}

	fmt.Println(conf.appId, conf.mark, conf.title)
}
