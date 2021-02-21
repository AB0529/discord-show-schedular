package main

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

// Conf the representation of a config file
type Conf struct {
	// Token the bot token
	Token string `json:"token"`
	// Prefix the prefix used to issue commands to the bot
	Prefix string `json:"prefix"`
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

// ReactionCollector waits for user response
type ReactionCollector struct {
	ReactionsCollected []*discordgo.MessageReaction
	Filter func (ctx *Context, m *discordgo.MessageReaction) bool
	EndAfterOne bool
	Timeout time.Duration
	UseTimeout bool
	Done chan bool
}


// DBShow the show in the database
type DBShow struct{
	MalID int
	Title string
}
// Database the users database
type Database map[string][]*DBShow

// Show representation of a show query response
type Show struct {
	RequestHash        string `json:"request_hash"`
	RequestCached      bool   `json:"request_cached"`
	RequestCacheExpiry int    `json:"request_cache_expiry"`
	Results            []struct {
		MalID     int       `json:"mal_id"`
		URL       string    `json:"url"`
		ImageURL  string    `json:"image_url"`
		Title     string    `json:"title"`
		Airing    bool      `json:"airing"`
		Synopsis  string    `json:"synopsis"`
		Type      string    `json:"type"`
		Episodes  int       `json:"episodes"`
		Score     float64   `json:"score"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		Members   int       `json:"members"`
		Rated     string    `json:"rated"`
	} `json:"results"`
	LastPage int `json:"last_page"`
}