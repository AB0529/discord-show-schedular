package main

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

// Conf the representation of a config file
type Conf struct {
	// Token the bot token
	Token string
	// Prefix the prefix used to issue commands to the bot
	Prefix string
	// Channel the channel to post all airings
	Channel string
}

// Context the context for the command
type Context struct {
	Session *discordgo.Session
	Msg     *discordgo.MessageCreate
	Command *Command
	LastMessage chan *discordgo.Message
	LastReaction chan *discordgo.MessageReaction
}

// Flag a command flag
type Flag struct {
	Name string
	Value string
	RequiresValue bool
	Exists bool
}

// Command the representation of a bot command
type Command struct {
	Name    string
	Aliases []string
	Example []string
	Desc    []string
	Handler func(*Context)
	Flags []*Flag
}

// MessageCollector waits for user response
type MessageCollector struct {
	MessagesCollected []*discordgo.Message
	Filter func (ctx *Context, m *discordgo.Message) bool
	EndAfterOne bool
	Timeout time.Duration
	UseTimeout bool
	Done chan bool
}

// DBShow the show in the database
type DBShow struct{
	MalID int
	Title string
	ImageURL string
	AlreadySent bool
}
// Database the users database
type Database map[string][]*DBShow