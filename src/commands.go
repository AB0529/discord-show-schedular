package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
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
	flags, err := ctx.FindCommandFlag()
	if err != nil {
		ctx.SendErr(err)
		return
	}

	for _, flag := range flags {
		ctx.Send(fmt.Sprintf("Name: %s\nValue: %s", flag.Name, flag.Value))
	}

}

// Schedule hub command which handles actions for the user's schedule
// TODO: make a better flag system
func Schedule(ctx *Context) {
	flags, err := ctx.FindCommandFlag()
	if err != nil {
		ctx.SendCommandHelp()
		return
	}
	if flags[0].RequiresValue && flags[0].Value == "" {
		ctx.SendErr("no value provided for the flag " + flags[0].Name)
		return
	}

	db := *NewDB()

	// Get the first flag
	switch flags[0].Name {
	case "add":
		show := flags[0].Value

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
				j, err := strconv.Atoi(m.Content)

				if m.Author.ID != ctx.Msg.Author.ID || err != nil {
					return false
				}
				if j > l {
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
			// Show hasn't aired yet
			if y <= 1 {
				ctx.SendErr("show has not yet aired!")
				return
			}

			ctx.SendErr(fmt.Sprintf("%s stopped airing on %d/%d/%d", shows.Results[res-1].Title, mo, d, y))
			// Delete message
			err := ctx.Session.ChannelMessageDelete(m.ChannelID, m.ID)
			ctx.SendErr(err)
			return
		}

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
		ctx.SendErr(err)

		AddShowToDatabase(&DBShow{
			MalID: shows.Results[res-1].MalID,
			Title: shows.Results[res-1].Title,
			ImageURL: shows.Results[res-1].ImageURL,
			AlreadySent: false,
		}, ctx.Msg.Author.ID)
		break
	case "li":
		fallthrough
	case "list":
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
		break
	case "rm":
		fallthrough
	case "remove":
		shows := db[ctx.Msg.Author.ID]
		if len(shows) < 1 {
			ctx.SendErr(fmt.Sprintf("You have no shows, add them with %s%s add ShowName", Config.Prefix, ctx.Command.Name))
			return
		}
		msg := ""
		for i, show := range shows {
			msg += fmt.Sprintf("%d) %s\n", i+1, show.Title)
		}
		m := ctx.NewEmbed(fmt.Sprintf("📚 | **%s**, you have `%d` show(s)!\n```css\n%s\nc. Cancel\n```", ctx.Msg.Author.Username, len(shows), msg))
		collector := &MessageCollector{
			MessagesCollected: []*discordgo.Message{},
			Filter: func(ctx *Context, m *discordgo.Message) bool {
				i, err := strconv.Atoi(m.Content)

				if m.Author.ID != ctx.Msg.Author.ID || err != nil {
					return false
				}
				if i > len(db[ctx.Msg.Author.ID]) {
					return false
				}

				return true
			},
			EndAfterOne:       true,
			Timeout:           time.Second * 10,
			UseTimeout:        true,
		}
		err := collector.New(ctx)
		if err != nil {
			// Delete message
			err := ctx.Session.ChannelMessageDelete(m.ChannelID, m.ID)
			ctx.SendErr(err)
			return
		}
		resp, _ := strconv.Atoi(collector.MessagesCollected[0].Content)

		// Remove from DB
		ctx.NewEmbed(fmt.Sprintf("🚫 | **Removed** `%s` from your schedule!", db[ctx.Msg.Author.ID][resp-1].Title))

		db[ctx.Msg.Author.ID] = append(db[ctx.Msg.Author.ID][:resp-1], db[ctx.Msg.Author.ID][resp:]...)
		db.Write()

		break
	}
}
