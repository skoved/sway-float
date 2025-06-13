// Copyright skoved
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// set by the compiler
var version string

// Prints err to stderr and calls os.Exit with a non-zero status code
func errorExit(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}

// flag stuff
var (
	helpFlag    bool
	versionFlag bool
	appIdFlag   string
	markFlag    string
	titleFlag   string
	matchFlag   string
)

// usage strings
const (
	helpUsage    = "print usage info"
	versionUsage = "print the version"
	appIdUsage   = "specify the appId of the window you want to float"
	markUsage    = "specify the con_mark of the window you want to float"
	titleUsage   = "specify the title of the window you want to match"
	matchUsage   = "specify how you want to match the title of windows. Valid options are: equal, prefix, suffix"
)

func setFlags() {
	flag.BoolVar(&helpFlag, "help", false, helpUsage)
	flag.BoolVar(&helpFlag, "h", false, helpUsage+" (shorthand)")
	flag.BoolVar(&versionFlag, "version", false, versionUsage)
	flag.BoolVar(&versionFlag, "v", false, versionUsage+" (shorthand)")
	flag.StringVar(&appIdFlag, "app_id", "", appIdUsage)
	flag.StringVar(&appIdFlag, "a", "", appIdUsage+" (shorthand)")
	flag.StringVar(&markFlag, "con_mark", "", markUsage)
	flag.StringVar(&markFlag, "c", "", markUsage+" (shorthand)")
	flag.StringVar(&titleFlag, "title", "", titleUsage)
	flag.StringVar(&titleFlag, "t", "", titleUsage+" (shorthand)")
	flag.StringVar(&matchFlag, "match", "equal", matchUsage)
	flag.StringVar(&matchFlag, "m", "equal", matchUsage+"(shorthand)")

	flag.Parse()
}

func main() {
	setFlags()
	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}
	if versionFlag {
		fmt.Fprintln(flag.CommandLine.Output(), os.Args[0], "version", version)
		fmt.Fprintln(flag.CommandLine.Output(), "Copyright (C) 2025 Sam Koved <https://github.com/skoved>")
		fmt.Fprintln(flag.CommandLine.Output())
		fmt.Fprintln(
			flag.CommandLine.Output(), "License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>",
		)
		fmt.Fprintln(
			flag.CommandLine.Output(),
			"This program comes with ABSOLUTELY NO WARRANTY;\nThis is free software, and you are welcome to redistribute it",
		)
		os.Exit(0)
	}

	var confs []floatConfig
	configDir, confDirErr := os.UserConfigDir()
	if confDirErr != nil {
		errorExit(fmt.Errorf("could not determine the current user's config dir: %w", confDirErr))
	}
	configPath := configDir + "/" + configDirSuffix
	fmt.Fprintln(os.Stderr, "reading config from:", configPath)
	_, statErr := os.Stat(configPath)
	if errors.Is(statErr, os.ErrNotExist) {
		if appIdFlag == "" && markFlag == "" && titleFlag == "" {
			fmt.Fprintf(
				flag.CommandLine.Output(),
				"Could not find config file %s. Please create a config file or use the %s, %s, or %s flags\n",
				configPath,
				appIdFlag,
				markFlag,
				titleFlag,
			)
			os.Exit(1)
		}
		me, meErr := ParsematcherEnum(matchFlag)
		if meErr != nil {
			errorExit(fmt.Errorf("invalid matcher %w", meErr))
		}
		opt, valid := me.toMatcher()
		if !valid {
			errorExit(errors.New("this should not happen"))
		}
		confs = []floatConfig{
			newFloatingConfig(appIdFlag, markFlag, titleFlag, opt),
		}
	} else if statErr != nil {
		errorExit(fmt.Errorf("encountered an issue checking if config file %s exists: %w", configPath, statErr))
	} else {
		configBytes, fileErr := os.ReadFile(configPath) //gosec:disable G304 -- need to build the path to the config dir
		if fileErr != nil {
			errorExit(fmt.Errorf("could not read file %s: %w", configPath, fileErr))
		}
		var yamlErr error
		confs, yamlErr = floatingConfigFromYaml(configBytes)
		if yamlErr != nil {
			errorExit(fmt.Errorf("could not parse yaml in file %s: %w", configPath, yamlErr))
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	handler, handlerErr := newWindowEventHandler(ctx, confs)
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
