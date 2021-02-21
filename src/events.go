package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// MessageCreate the function which handles message events
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bots
	if m.Author.Bot {
		return
	}

	// Create new context on each message
	c := strings.Split(strings.ToLower(m.Message.Content)[1:], " ")[0]
	// Find the command with the matching name alias and run it
	cmd, ok := Commands[c]
	if !ok {
		return
	}
	ctx := &Context{
		Session: s,
		Msg:     m,
		Command: cmd,
	}
	cmd.Handler(ctx)
}

// Ready the function which handles when the bot is ready
func Ready(_ *discordgo.Session, e *discordgo.Ready) {
	fmt.Printf("[%s] - in as %s%s with prefix: \"%s\"\n", Purple.Sprint("BOT"), Yellow.Sprint(e.User.Username), Yellow.Sprint("#"+e.User.Discriminator), Green.Sprint(Config.Prefix))
}
