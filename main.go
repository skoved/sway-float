// Copyright skoved
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/goccy/go-yaml"
	"github.com/joshuarubin/go-sway"
)

const configDirSuffix = "sway-float/config.yaml"

type floatConfig struct {
	AppId string `yaml:"app_id"`
	Mark  string `yaml:"con_mark"`
	Title string `yaml:"title"`
}

func (c floatConfig) match(event sway.WindowEvent) bool {
	return c.appIdMatch(event) && c.markMatch(event) && c.titleMatch(event)
}

// should only be called by match
func (c floatConfig) appIdMatch(event sway.WindowEvent) bool {
	if c.AppId == "" {
		return true
	}
	if event.Container.AppID == nil {
		return false
	}
	return c.AppId == *event.Container.AppID
}

func (c floatConfig) markMatch(event sway.WindowEvent) bool {
	if c.Mark == "" {
		return true
	}
	return slices.Contains(event.Container.Marks, c.Mark)
}

func (c floatConfig) titleMatch(event sway.WindowEvent) bool {
	if c.Title == "" {
		return true
	}
	return c.Title == event.Container.Name
}

// Prints err to stderr and calls os.Exit with a non-zero status code
func errorExit(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}

func main() {
	var conf floatConfig
	configDir, confDirErr := os.UserConfigDir()
	if confDirErr != nil {
		errorExit(fmt.Errorf("could not determine the current user's config dir: %w", confDirErr))
	}
	configPath := configDir + "/" + configDirSuffix
	fmt.Fprintln(os.Stderr, "reading config from:", configPath)
	configBytes, fileErr := os.ReadFile(configPath) //gosec:disable G304 -- need to build the path to the config dir
	if fileErr != nil {
		errorExit(fmt.Errorf("could not read file %s: %w", configPath, fileErr))
	}
	yamlErr := yaml.Unmarshal(configBytes, &conf)
	if yamlErr != nil {
		errorExit(fmt.Errorf("could not parse yaml in file %s: %w", configPath, yamlErr))
	}
	fmt.Println(conf.AppId, conf.Mark, conf.Title)

	ctx, cancel := context.WithCancel(context.Background())
	handler, handlerErr := newWindowEventHandler(ctx, conf)
	if handlerErr != nil {
		errorExit(handlerErr)
	}

	errCh := make(chan error)

	go handler.handle(ctx, errCh)

	sigCtx, sigCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	select {
	case subErr := <-errCh:
		sigCancel()
		errorExit(subErr)
	case <-sigCtx.Done():
		cancel()
		fmt.Fprintln(os.Stderr, "Received interrupt. Shutting down")
	}
}
