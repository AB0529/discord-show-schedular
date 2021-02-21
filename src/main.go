package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Context the context for the command
type Context struct {
	Session *discordgo.Session
	Msg     *discordgo.MessageCreate
	Command *Command
}

// Commands all the commands for the bot
var Commands = make(map[string]*Command)

func main() {
	// Load config
	Config = NewConfig("../config.yml")
	Config = NewConfig("../config.yml")
	// Setup Discord
	dg, _ := discordgo.New("Bot " + Config.Token)
	// Register events
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	dg.AddHandler(MessageCreate)
	dg.AddHandler(Ready)

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
			Example: []string{Config.Prefix + "schedule <flag> <value>"},
			Desc:    []string{"Controls the user's schedule"},
			Handler: Schedule,
			Flags:   []string{"add", "remove", "list"},
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
