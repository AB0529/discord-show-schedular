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
	ImageURL string
	AlreadySent bool
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

// AnimeSchedule currently airing anime schedule
type AnimeSchedule struct {
	RequestHash        string `json:"request_hash"`
	RequestCached      bool   `json:"request_cached"`
	RequestCacheExpiry int    `json:"request_cache_expiry"`
	Monday             []struct {
		MalID       int       `json:"mal_id"`
		URL         string    `json:"url"`
		Title       string    `json:"title"`
		ImageURL    string    `json:"image_url"`
		Synopsis    string    `json:"synopsis"`
		Type        string    `json:"type"`
		AiringStart time.Time `json:"airing_start"`
		Episodes    int       `json:"episodes"`
		Members     int       `json:"members"`
		Genres      []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"genres"`
		Source    string `json:"source"`
		Producers []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"producers"`
		Score     float64  `json:"score"`
		Licensors []string `json:"licensors"`
		R18       bool     `json:"r18"`
		Kids      bool     `json:"kids"`
	} `json:"monday"`
	Tuesday []struct {
		MalID       int       `json:"mal_id"`
		URL         string    `json:"url"`
		Title       string    `json:"title"`
		ImageURL    string    `json:"image_url"`
		Synopsis    string    `json:"synopsis"`
		Type        string    `json:"type"`
		AiringStart time.Time `json:"airing_start"`
		Episodes    int       `json:"episodes"`
		Members     int       `json:"members"`
		Genres      []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"genres"`
		Source    string `json:"source"`
		Producers []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"producers"`
		Score     float64  `json:"score"`
		Licensors []string `json:"licensors"`
		R18       bool     `json:"r18"`
		Kids      bool     `json:"kids"`
	} `json:"tuesday"`
	Wednesday []struct {
		MalID       int       `json:"mal_id"`
		URL         string    `json:"url"`
		Title       string    `json:"title"`
		ImageURL    string    `json:"image_url"`
		Synopsis    string    `json:"synopsis"`
		Type        string    `json:"type"`
		AiringStart time.Time `json:"airing_start"`
		Episodes    int       `json:"episodes"`
		Members     int       `json:"members"`
		Genres      []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"genres"`
		Source    string `json:"source"`
		Producers []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"producers"`
		Score     float64       `json:"score"`
		Licensors []interface{} `json:"licensors"`
		R18       bool          `json:"r18"`
		Kids      bool          `json:"kids"`
	} `json:"wednesday"`
	Thursday []struct {
		MalID       int       `json:"mal_id"`
		URL         string    `json:"url"`
		Title       string    `json:"title"`
		ImageURL    string    `json:"image_url"`
		Synopsis    string    `json:"synopsis"`
		Type        string    `json:"type"`
		AiringStart time.Time `json:"airing_start"`
		Episodes    int       `json:"episodes"`
		Members     int       `json:"members"`
		Genres      []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"genres"`
		Source    string `json:"source"`
		Producers []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"producers"`
		Score     float64  `json:"score"`
		Licensors []string `json:"licensors"`
		R18       bool     `json:"r18"`
		Kids      bool     `json:"kids"`
	} `json:"thursday"`
	Friday []struct {
		MalID       int       `json:"mal_id"`
		URL         string    `json:"url"`
		Title       string    `json:"title"`
		ImageURL    string    `json:"image_url"`
		Synopsis    string    `json:"synopsis"`
		Type        string    `json:"type"`
		AiringStart time.Time `json:"airing_start"`
		Episodes    int       `json:"episodes"`
		Members     int       `json:"members"`
		Genres      []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"genres"`
		Source    string `json:"source"`
		Producers []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"producers"`
		Score     float64  `json:"score"`
		Licensors []string `json:"licensors"`
		R18       bool     `json:"r18"`
		Kids      bool     `json:"kids"`
	} `json:"friday"`
	Saturday []struct {
		MalID       int       `json:"mal_id"`
		URL         string    `json:"url"`
		Title       string    `json:"title"`
		ImageURL    string    `json:"image_url"`
		Synopsis    string    `json:"synopsis"`
		Type        string    `json:"type"`
		AiringStart time.Time `json:"airing_start"`
		Episodes    int       `json:"episodes"`
		Members     int       `json:"members"`
		Genres      []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"genres"`
		Source    string `json:"source"`
		Producers []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"producers"`
		Score     float64       `json:"score"`
		Licensors []interface{} `json:"licensors"`
		R18       bool          `json:"r18"`
		Kids      bool          `json:"kids"`
	} `json:"saturday"`
	Sunday []struct {
		MalID       int         `json:"mal_id"`
		URL         string      `json:"url"`
		Title       string      `json:"title"`
		ImageURL    string      `json:"image_url"`
		Synopsis    string      `json:"synopsis"`
		Type        string      `json:"type"`
		AiringStart time.Time   `json:"airing_start"`
		Episodes    interface{} `json:"episodes"`
		Members     int         `json:"members"`
		Genres      []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"genres"`
		Source    string `json:"source"`
		Producers []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"producers"`
		Score     float64  `json:"score"`
		Licensors []string `json:"licensors"`
		R18       bool     `json:"r18"`
		Kids      bool     `json:"kids"`
	} `json:"sunday"`
	Other []struct {
		MalID       int         `json:"mal_id"`
		URL         string      `json:"url"`
		Title       string      `json:"title"`
		ImageURL    string      `json:"image_url"`
		Synopsis    string      `json:"synopsis"`
		Type        string      `json:"type"`
		AiringStart time.Time   `json:"airing_start"`
		Episodes    interface{} `json:"episodes"`
		Members     int         `json:"members"`
		Genres      []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"genres"`
		Source    string `json:"source"`
		Producers []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"producers"`
		Score     float64       `json:"score"`
		Licensors []interface{} `json:"licensors"`
		R18       bool          `json:"r18"`
		Kids      bool          `json:"kids"`
	} `json:"other"`
	Unknown []struct {
		MalID       int       `json:"mal_id"`
		URL         string    `json:"url"`
		Title       string    `json:"title"`
		ImageURL    string    `json:"image_url"`
		Synopsis    string    `json:"synopsis"`
		Type        string    `json:"type"`
		AiringStart time.Time `json:"airing_start"`
		Episodes    int       `json:"episodes"`
		Members     int       `json:"members"`
		Genres      []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"genres"`
		Source    string `json:"source"`
		Producers []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"producers"`
		Score     float64       `json:"score"`
		Licensors []interface{} `json:"licensors"`
		R18       bool          `json:"r18"`
		Kids      bool          `json:"kids"`
	} `json:"unknown"`
}