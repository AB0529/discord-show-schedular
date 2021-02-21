package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var LastMessage = make(chan *discordgo.Message)
var LastReaction = make(chan *discordgo.MessageReaction)

// MessageCreate the function which handles message events
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bots
	if m.Author.Bot {
		return
	}
	// Ignore no content messages
	if m.Content == "" {
		return
	}

	// Create new context on each message1
	msg := strings.Split(strings.ToLower(m.Message.Content)[1:], " ")

	if len(msg) <= 0 {
		// Send message to LastMessage channel
		go func() { LastMessage <- m.Message }()
		return
	}

	c := msg[0]
	// Find the command with the matching name alias and run it
	cmd, ok := Commands[c]
	if !ok {
		go func() { LastMessage <- m.Message }()
		return
	}
	// Make sure message starts with prefix
	if string(m.Message.Content[0]) != Config.Prefix {
		return
	}
	ctx := &Context{
		Session: s,
		Msg:     m,
		Command: cmd,
		LastMessage: LastMessage,
		LastReaction: LastReaction,
	}
	cmd.Handler(ctx)
}

// MessageReactionAdd the function which handles message reaction events
func MessageReactionAdd(_ *discordgo.Session, m *discordgo.MessageReactionAdd) {
	go func() { LastReaction <- m.MessageReaction }()
}

// Ready the function which handles when the bot is ready
func Ready(s *discordgo.Session, e *discordgo.Ready) {
	fmt.Printf("[%s] - in as %s%s with prefix: \"%s\"\n", Purple.Sprint("BOT"), Yellow.Sprint(e.User.Username), Yellow.Sprint("#"+e.User.Discriminator), Green.Sprint(Config.Prefix))
	err := s.UpdateGameStatus(0, "with yo momma")
	Die(err)

	// Keep track of anime schedules
	for {
		db := NewDB()
		weekday := strings.ToLower(time.Now().Weekday().String())
		for userID, userShows := range *db {
			for _, show := range userShows {
				// Form http request
				resp, err := http.Get("https://api.jikan.moe/v3/schedule/" + weekday)
				if err != nil {
					panic(err)
				}
				body, _ := ioutil.ReadAll(resp.Body)
				resp.Body.Close()
				schedule := &AnimeSchedule{}
				_ = json.Unmarshal(body, &schedule)

				SendMsg := func() {
					// DM User
					channel, err := s.UserChannelCreate(userID)
					Warn(err)
					_, err = s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
						Embed: &discordgo.MessageEmbed{
							Color:       rand.Intn(10000000),
							Description: fmt.Sprintf("🔔 | Ding, ding, `%s` **new episode** aired!", show.Title),
							Image: &discordgo.MessageEmbedImage{URL: show.ImageURL},
						},
					})
					Warn(err)
					show.AlreadySent = true
					db.Write()
				}


				// Find show in schedule
				switch weekday {
				case "monday":
					for _, ss := range schedule.Monday {
						if ss.MalID == show.MalID && !show.AlreadySent {
							SendMsg()
						}
					}
				case "tuesday":
					for _, ss := range schedule.Tuesday {
						if ss.MalID == show.MalID && !show.AlreadySent {
							SendMsg()
						}
					}
				case "wednesday":
					for _, ss := range schedule.Wednesday {
						if ss.MalID == show.MalID && !show.AlreadySent {
							SendMsg()
						}
					}
				case "thursday":
					for _, ss := range schedule.Thursday {
						if ss.MalID == show.MalID && !show.AlreadySent {
							SendMsg()
						}
					}
				case "friday":
					for _, ss := range schedule.Friday {
						if ss.MalID == show.MalID && !show.AlreadySent {
							SendMsg()
						}
					}
				case "saturday":
					for _, ss := range schedule.Saturday {
						if ss.MalID == show.MalID && !show.AlreadySent {
							SendMsg()
						}
					}
				case "sunday":
					for _, ss := range schedule.Sunday {
						if ss.MalID == show.MalID && !show.AlreadySent {
							SendMsg()
						}
					}
				default:
					show.AlreadySent = false
					db.Write()
				}

			}
		}

		// Check every 10 seconds
		time.Sleep(time.Second * 10)
	}

}

