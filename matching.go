// Copyright skoved
// SPDX-License-Identifier: GPL-3.0-or-later
//go:generate go-enum --marshal -t match_enum.tmpl

package main

import "strings"

// ENUM(equal, prefix, suffix)
//
//nolint:recvcheck
type matcherEnum int

type MatcherOption func(floatConfig) floatConfig

type MatcherFunc func(string) bool

// Returns a MatcherOption that always returns true
//
//nolint:unused
func trueMatcher(_ string) bool {
	return true
}

// Returns a MatcherOption that checks if the title of a window is the title from floatingConfig
func WithEqualMatcher() MatcherOption {
	return func(fc floatConfig) floatConfig {
		fc.titleMatch = equalMatcher(fc.Title)
		return fc
	}
}

func equalMatcher(confString string) MatcherFunc {
	if confString == "" {
		return trueMatcher
	}
	return func(eventString string) bool { return confString == eventString }
}

// Returns a MatcherOption that checks if the title of a window starts with the title from floatingConfig
func WithPrefixMatcher() MatcherOption {
	return func(fc floatConfig) floatConfig {
		fc.titleMatch = prefixMatcher(fc.Title)
		return fc
	}
}

func prefixMatcher(confString string) MatcherFunc {
	if confString == "" {
		return trueMatcher
	}
	return func(eventString string) bool { return strings.HasPrefix(eventString, confString) }
}

// Returns a MatcherOption that checks if the title of a window ends with the title from floatingConfig
func WithSuffixMatcher() MatcherOption {
	return func(fc floatConfig) floatConfig {
		fc.titleMatch = suffixMatcher(fc.Title)
		return fc
	}
}

func suffixMatcher(confString string) MatcherFunc {
	if confString == "" {
		return trueMatcher
	}
	return func(eventString string) bool { return strings.HasSuffix(eventString, confString) }
}
