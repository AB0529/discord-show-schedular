package main

import (
	"fmt"
	"time"
)

// Command the representation of a bot command
type Command struct {
	Name    string
	Aliases []string
	Example []string
	Desc    []string
	Handler func(*Context)
	Flags   []string
}

// Ping command which returns a message
func Ping(ctx *Context) {
	m := ctx.NewEmbed("Pinging....")
	ts, _ := m.Timestamp.Parse()
	now := time.Now()
	ctx.EditEmbed(m, fmt.Sprintf("🏓 | **Pong my ping**\n\n💗 | **Heartbeat**: `%1.fms`\n ⏱️| **Message Delay**: `%1.fms`",
		float64(ctx.Session.HeartbeatLatency().Milliseconds()),
		float64(now.Sub(ts).Milliseconds())))
}

// Schedule hub command which handles actions for the user's schedule
func Schedule(ctx *Context) {
	args := ctx.FindCommandFlag()
	if args == nil {
		ctx.SendCommandHelp()
	}
	// Add
	if args["add"] != "" {
		show := args["add"]
		shows, err := FindShows(show)
		if err != nil {
			ctx.SendErr(err)
			return
		}
		m := ctx.Send("Loading...")
		ctx.Edit(m, fmt.Sprintf("Title: %s\nDesc: %s", shows.Results[0].Title, shows.Results[0].Synopsis))
	} else {
		ctx.NewEmbed(":x: | What **anime** do you want to add?")
	}
}
