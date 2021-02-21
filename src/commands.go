package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"strings"
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
	if _, ok := args["add"]; ok {
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
		// Make sure show is still airing before adding it
		if !shows.Results[res-1].Airing {
			y, mo, d := shows.Results[res-1].EndDate.Date()
			ctx.SendErr(fmt.Sprintf("%s stopped airing on %d/%d/%d", shows.Results[res-1].Title, mo, d, y))
			// Delete message
			err := ctx.Session.ChannelMessageDelete(m.ChannelID, m.ID)
			ctx.SendErr(err)
			return
		}
		db := *NewDB()
		// Check for duplicates
		for _, show := range db[ctx.Msg.Author.ID] {
			if show.MalID == shows.Results[res-1].MalID {
				ctx.SendErr("Show already in user schedule, you silly goose")
				return
			}
		}

		_, err = ctx.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Embed: &discordgo.MessageEmbed{
				Color:       rand.Intn(10000000),
				Description: fmt.Sprintf("📺 | Got it, I'll add `%s` to your schedule!", shows.Results[res-1].Title),
				Image: &discordgo.MessageEmbedImage{URL: shows.Results[res-1].ImageURL},
				Footer: &discordgo.MessageEmbedFooter{IconURL: ctx.Msg.Author.AvatarURL(""), Text: ctx.Msg.Author.Username},
			},
			ID: m.ID,
			Channel: m.ChannelID,
		})
		// Delete message
		err = ctx.Session.ChannelMessageDelete(m.ChannelID, m.ID)
		ctx.SendErr(err)

		AddShowToDatabase(&DBShow{
			MalID: shows.Results[res-1].MalID,
			Title: shows.Results[res-1].Title,
		}, ctx.Msg.Author.ID)
	}
	// List
	if strings.Split(ctx.Msg.Content, " ")[1] == "list" {
		db := *NewDB()
		// Loop through users shows
		shows := db[ctx.Msg.Author.ID]
		if len(shows) < 1 {
			ctx.SendErr(fmt.Sprintf("You have no shows, add them with %s%s add ShowName", Config.Prefix, ctx.Command.Name))
			return
		}
		msg := ""
		for i, show := range shows {
			msg += fmt.Sprintf("```css\n%d) %s\n```", i+1, show.Title)
		}
		ctx.NewEmbed(fmt.Sprintf("📚 | **%s**, you have `%d` show(s)!\n%s", ctx.Msg.Author.Username, len(shows), msg))
	}
}
