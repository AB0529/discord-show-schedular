package main

import (
	"fmt"
	"github.com/darenliang/jikan-go"
	"math"
	"math/rand"
	"strconv"
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

	weekdaysPrev := map[string]string{
		"sun": "sat",
		"mon": "sun",
		"tue": "mon",
		"wed": "tue",
		"thu": "wed",
		"fri": "thu",
		"sat": "fri",
	}

	// Keep track of anime schedules
	db := NewDB()
	l, _ :=  time.LoadLocation("America/New_York")

	for {
		for userID, userShows := range *db {
			for _, show := range userShows {
				weekday := strings.ToLower(time.Now().Weekday().String())[:3]
				anime, _ := jikan.GetAnime(show.MalID)
				// Make sure it's still airing
				// TODO: delete show after it's done airing
				if !anime.Airing {
					continue
				}
				// Get the day of the week, and the time from broadcast time
				bsr := strings.Split(strings.ToLower(anime.Broadcast), " ")
				airWeekday := bsr[0][:len(bsr)-1]

				//fmt.Printf("%s : %s : %s\n", show.Title, weekday, airWeekday)
				// Roll back a day
				if bsr[2][:2] == "00" {
					airWeekday = weekdaysPrev[airWeekday]
					//fmt.Printf("Rolled back %s : %s : %s\n", show.Title, weekday, airWeekday)
				}


				// Check if airing at current weekday
				if weekday == airWeekday && !show.AlreadySent {
					// Get time duration of airing from now
					ahStr, amStr := bsr[2][2:], bsr[2][:2]
					ah, _ := strconv.Atoi(ahStr)
					am, _ := strconv.Atoi(amStr)
					h, m, _ := time.Now().In(l).Clock()

					timeDurH := int(math.Abs(float64(h) - float64(ah)))
					timeDurM := int(math.Abs(float64(m) - float64(am)))

					embedToSend := &discordgo.MessageSend{
						Embed: &discordgo.MessageEmbed{
							Color:       rand.Intn(10000000),
							Description: fmt.Sprintf("🔔 | Airing in **%d hours and %d mins**", timeDurH, timeDurM),
							Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("%s on %s", show.Title, weekday)},
							Image: &discordgo.MessageEmbedImage{URL: show.ImageURL},
						},
					}

					// DM User
					channel, err := s.UserChannelCreate(userID)
					Warn(err)
					_, err = s.ChannelMessageSendComplex(channel.ID, embedToSend)

					if Config.Channel != "" {
						_, err = s.ChannelMessageSendComplex(Config.Channel, embedToSend)
						Warn(err)
					}

					show.AlreadySent = true
					db.Write()
				}

				time.Sleep(time.Second * 5)

			}
		}

		//fmt.Println("------------------------------------------------------")
	}

}

