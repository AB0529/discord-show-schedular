package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Commands all the commands for the bot
var Commands = make(map[string]*Command)

func main() {
	// Load config
	Config = NewConfig("../config.yml")
	Config = NewConfig("../config.yml")
	// Setup Discord
	dg, _ := discordgo.New("Bot " + Config.Token)
	// Register events
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages|discordgo.IntentsGuildMessageReactions|discordgo.IntentsGuildEmojis)
	dg.AddHandler(MessageCreate)
	dg.AddHandler(Ready)
	dg.AddHandler(MessageReactionAdd)

	// Register commands
	RegisterCommands([]*Command{
		{
			Name:    "ping",
			Aliases: []string{"pong"},
			Example: []string{Config.Prefix + "ping"},
			Desc:    []string{"Generic Ping-Pong command"},
			Handler: Ping,
		},
		{
			Name:    "schedule",
			Aliases: []string{"sched", "sc"},
			Example: []string{Config.Prefix + "schedule <flag> <value>", Config.Prefix + "schedule add Bleach"},
			Desc:    []string{"Controls the user's schedule", "Flag 'add': adds an anime to your schedule", "Flag 'list': lists your schedule", "Flag 'remove': removes an anime to your schedule"},
			Handler: Schedule,
			Flags: []*Flag{ { Name: "add", RequiresValue: true }, { Name: "list" }, { Name: "li" }, { Name: "remove" }, { Name: "rm" } },
		},
		{
			Name:    "test",
			Aliases: []string{},
			Example: []string{Config.Prefix + "test <flag> <value>"},
			Desc:    []string{"Command used for testing"},
			Handler: Test,
			Flags: []*Flag{ { Name: "add", RequiresValue: true }, { Name: "list" } },
		},
	})

	// Open a websocket connection to Discord and begin listening.
	err := dg.Open()
	if err != nil {
		Die("could not creating session")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	err = dg.Close()
	Die(err)
}
