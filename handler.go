// Copyright skoved
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joshuarubin/go-sway"
)

const floatCmd = "floating enable"

var _ sway.EventHandler = (*windowEventHandler)(nil)

type windowEventHandler struct {
	client sway.Client
	sway.EventHandler
	confs []floatConfig
}

func (h windowEventHandler) Window(ctx context.Context, event sway.WindowEvent) {
	// only want window events where a title or mark is changed
	if event.Change != sway.WindowTitle && event.Change != sway.WindowMark {
		return
	}
	if event.Container.Type == sway.NodeFloatingCon {
		return
	}
	for _, conf := range h.confs {
		if conf.match(event) {
			replies, cmdErr := h.client.RunCommand(ctx, floatCmd)
			if cmdErr != nil {
				fmt.Fprintf(os.Stderr, "Could not run command '%s': %s\n", floatCmd, cmdErr.Error())
			}
			for _, reply := range replies {
				if reply.Error != "" {
					fmt.Printf("Command '%s' failed: %s\n", floatCmd, reply.Error)
				}
			}
		}
	}
}

func (h windowEventHandler) handle(ctx context.Context, quit chan<- error) {
	subErr := sway.Subscribe(ctx, h, sway.EventTypeWindow)
	if subErr != nil {
		quit <- fmt.Errorf("could not subscribe to window events: %w", subErr)
	}
	close(quit)
}

// Creates a new WindowEventHandler which only responds to sway.WindowEvent
func newWindowEventHandler(ctx context.Context, c []floatConfig) (windowEventHandler, error) {
	client, clientErr := sway.New(ctx)
	if clientErr != nil {
		return windowEventHandler{}, fmt.Errorf("could not create sway client: %w", clientErr)
	}
	return windowEventHandler{
		client:       client,
		confs:        c,
		EventHandler: sway.NoOpEventHandler(),
	}, nil
}
