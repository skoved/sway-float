// Copyright skoved
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"slices"

	"github.com/goccy/go-yaml"
	"github.com/joshuarubin/go-sway"
)

const configDirSuffix = "sway-float/config.yaml"

type floatConfig struct {
	AppId      string      `yaml:"app_id"`
	Mark       string      `yaml:"con_mark"`
	Title      string      `yaml:"title"`
	Matcher    matcherEnum `yaml:"matcher"`
	titleMatch MatcherFunc
}

func newFloatingConfig(appId, mark, title string, opt MatcherOption) floatConfig {
	c := floatConfig{
		AppId:      appId,
		Mark:       mark,
		Title:      title,
		Matcher:    MatcherEnumEqual,
		titleMatch: nil,
	}
	return opt(c)
}

func floatingConfigFromYaml(y []byte) ([]floatConfig, error) {
	var confs []floatConfig
	yErr := yaml.Unmarshal(y, &confs)
	if yErr != nil {
		return nil, fmt.Errorf("could not parse config file %w", yErr)
	}
	for i, conf := range confs {
		opt, valid := conf.Matcher.toMatcher()
		if !valid {
			return nil, fmt.Errorf("invalid matcher %s", conf.Matcher.String())
		}
		confs[i] = opt(conf)
	}
	return confs, nil
}

func (c floatConfig) match(event sway.WindowEvent) bool {
	if c.AppId == "" && c.Mark == "" && c.Title == "" {
		return false
	}
	return c.appIdMatch(event) && c.markMatch(event) && c.titleMatch(event.Container.Name)
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
