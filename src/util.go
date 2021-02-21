package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gookit/color"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	// Red color red
	Red = color.Red
	// Purple color purple
	Purple = color.Magenta
	// Green color green
	Green = color.LightGreen
	// Yellow color yellow
	Yellow = color.Yellow
)

// StoredShows the shows stored
type StoredShows struct {
	Title string `json:"title"`
	MalID int    `json:"malID"`
}

// Show repersentation of a show query response
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

// Episode response type of episode request
type Episode struct {
	RequestHash        string `json:"request_hash"`
	RequestCached      bool   `json:"request_cached"`
	RequestCacheExpiry int    `json:"request_cache_expiry"`
	EpisodesLastPage   int    `json:"episodes_last_page"`
	Episodes           []struct {
		EpisodeID     int       `json:"episode_id"`
		Title         string    `json:"title"`
		TitleJapanese string    `json:"title_japanese"`
		TitleRomanji  string    `json:"title_romanji"`
		Aired         time.Time `json:"aired"`
		Filler        bool      `json:"filler"`
		Recap         bool      `json:"recap"`
		VideoURL      string    `json:"video_url"`
		ForumURL      string    `json:"forum_url"`
	} `json:"episodes"`
}

// Warn logs warning to stdout
func Warn(err interface{}) {
	// Handle strings being passed by creating an error type
	if err != nil {
		if e, ok := err.(string); ok {
			err = errors.New(e)
		}
		fmt.Printf("[%s] - %v\n", Yellow.Sprint("WARN"), err)
	}
}

// Die handles errors and exits
func Die(err interface{}) {
	// Handle strings being passed by creating an error type
	if err != nil {
		if e, ok := err.(string); ok {
			err = errors.New(e)
		}
		fmt.Printf("[%s] - %v\n", Red.Sprint("ERR"), err)
		os.Exit(1)
	}
}

// RegisterCommands register provided commands
func RegisterCommands(cmds []*Command) {
	// Loop through each command and add it to Commands slice
	for _, c := range cmds {
		Commands[c.Name] = c
		for _, a := range c.Aliases {
			Commands[a] = c
		}
		fmt.Printf("[%s] - Loaded %s command\n", Purple.Sprint("CMD"), Yellow.Sprint(c.Name))
	}
}

// NewEmbed creates a simple embed with description only and a random color
func (ctx *Context) NewEmbed(content string) *discordgo.Message {
	m, err := ctx.Session.ChannelMessageSendComplex(ctx.Msg.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color:       rand.Intn(10000000),
			Description: content,
		},
	})
	Warn(err)
	return m
}

// Send sends a message to the same channel as the command
func (ctx *Context) Send(content string) *discordgo.Message {
	m, err := ctx.Session.ChannelMessageSend(ctx.Msg.ChannelID, content)
	Warn(err)
	return m
}

// SendErr sends a error message to the same channel as the command
func (ctx *Context) SendErr(content interface{}) {
	// Handle strings being passed by creating an error type
	if content != nil {
		if e, ok := content.(string); ok {
			content = errors.New(e)
		}
		ctx.NewEmbed(fmt.Sprintf(":x: | Uh oh, something **went wrong**!\n```css\n%s\n```", content))
	}
}

// Edit edits a message with a new content
func (ctx *Context) Edit(m *discordgo.Message, content string) *discordgo.Message {
	m, err := ctx.Session.ChannelMessageEdit(m.ChannelID, m.ID, content)
	Warn(err)
	return m
}

// EditEmbed edits a embed with a new content
func (ctx *Context) EditEmbed(m *discordgo.Message, content string) *discordgo.Message {
	m, err := ctx.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Embed: &discordgo.MessageEmbed{
			Color:       rand.Intn(10000000),
			Description: content,
		},
		ID: m.ID,
		Channel: m.ChannelID,
	})
	Warn(err)
	return m
}

// FindCommandFlag finds a flag and the value in a command string
func (ctx *Context) FindCommandFlag() map[string]string {
	// Find the flags
	flags := strings.Split(strings.ToLower(ctx.Msg.Message.Content), " ")[1:]
	// Make sure flags is not empty
	if len(flags) <= 0 {
		return nil
	}

	return map[string]string{flags[0]: strings.Join(flags[1:], " ")}
}

// SendCommandHelp properly formats and shows the help page of a command
func (ctx *Context) SendCommandHelp() {
	_, err := ctx.Session.ChannelMessageSendComplex(ctx.Msg.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Description: fmt.Sprintf("`%s%s` Command Help", Config.Prefix, ctx.Command.Name),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name: "📜 | Description",
					Value: fmt.Sprintf("```css\n%s\n```", strings.Join(ctx.Command.Desc, "\n")),
				},
				{
					Name: "🤖 | Example",
					Value: fmt.Sprintf("```css\n%s\n```", strings.Join(ctx.Command.Example, "\n")),
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Aliases: %s", strings.Join(ctx.Command.Aliases, " | ")),
				IconURL: ctx.Msg.Message.Author.AvatarURL("512x512"),
			},
			Color:       rand.Intn(10000000),
		},
	})
	Warn(err)
}

// FindShows find anime from mal with the provided input
func FindShows(query string) (*Show, error) {
	// Form http request
	resp, err := http.Get("https://api.jikan.moe/v3/search/anime?q=" + strings.Join(strings.Split(query, " "), "%20"))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	show := &Show{}
	_ = json.Unmarshal(body, &show)

	// Return the first show if only one is present
	if len(show.Results) >= 1 {
		return show, nil
	}
	if len(show.Results) <= 0 {
		return nil, errors.New("could not find show")
	}

	return nil, errors.New("could not find show")
}