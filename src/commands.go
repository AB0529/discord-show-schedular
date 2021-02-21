package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

// Ping command which returns a message
func Ping(ctx *Context) {
	m := ctx.NewEmbed("Pinging....")
	ts, _ := m.Timestamp.Parse()
	now := time.Now()
	ctx.EditEmbed(m, fmt.Sprintf("🏓 | **Pong my ping**\n\n💗 | **Heartbeat**: `%1.fms`\n ⏱️| **Message Delay**: `%1.fms`",
		float64(ctx.Session.HeartbeatLatency().Milliseconds()),
		float64(now.Sub(ts).Milliseconds())))
}

// Test command used for testing
func Test(ctx *Context) {
	ctx.Send("Say something nigga")
	collector := &MessageCollector{
		MessagesCollected: []*discordgo.Message{},
		Filter:            func(ctx *Context, m *discordgo.Message) bool {
			return m.Timestamp >= ctx.Msg.Timestamp
		},
		EndAfterOne:       false,
		Timeout:           time.Second * 5,
		UseTimeout:        true,
	}
	err := collector.New(ctx)
	if err != nil {
		ctx.SendErr(err)
		return
	}

	ctx.Send(fmt.Sprintf("You said %d many things", len(collector.MessagesCollected)))


}

// Schedule hub command which handles actions for the user's schedule
func Schedule(ctx *Context) {
	args := ctx.FindCommandFlag()
	if args == nil {
		ctx.SendCommandHelp()
		return
	}
	// Add
	if args["add"] != "" {
		show := args["add"]
		shows, err := FindShows(show)
		if err != nil {
			ctx.SendErr(err)
			return
		}
		m := ctx.NewEmbed("Loading...")
		// List 1-5 shows
		l := 1
		msg := ""
		if len(shows.Results) > 1 {
			l = 5
		}
		// Format the results
		for i := 0; i < l; i++ {
			msg += fmt.Sprintf("%d) %s\n", i+1, shows.Results[i].Title)
		}
		// Get the user input
		m = ctx.EditEmbed(m, fmt.Sprintf("Results for `%s`:\n```css\n%s\nc) Cancel\n```\n*Type the number of the show you want; you have 10 seconds*", show, msg))
		collector := &MessageCollector{
			MessagesCollected: []*discordgo.Message{},
			Filter: func(ctx *Context, m *discordgo.Message) bool {
				if _, err := strconv.Atoi(m.Content); m.Author.ID != ctx.Msg.Author.ID || err != nil {
					return false
				}

				return true
			},
			EndAfterOne:       true,
			Timeout:           time.Second * 10,
			UseTimeout:        true,
		}
		err = collector.New(ctx)
		if err != nil {
			// Delete message
			err := ctx.Session.ChannelMessageDelete(m.ChannelID, m.ID)
			ctx.SendErr(err)
			return
		}
		resMsg := collector.MessagesCollected[0]

		// Return the selected output
		res, err := strconv.Atoi(resMsg.Content)
		if err != nil {
			ctx.SendErr(err)
			return
		}
		ctx.EditEmbed(m, "You chose: " + shows.Results[res-1].Title)


	} else {
		ctx.NewEmbed(":x: | What **anime** do you want to add?")
	}
}
