package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gookit/color"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
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
			Footer: &discordgo.MessageEmbedFooter{IconURL: ctx.Msg.Author.AvatarURL("")},
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
			Footer: &discordgo.MessageEmbedFooter{IconURL: ctx.Msg.Author.AvatarURL("512x512")},
		},
		ID: m.ID,
		Channel: ctx.Msg.ChannelID,
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

// New collects user messages
// TODO: implement done channel
// TODO: implement reaction collector
func (collector *MessageCollector) New(ctx *Context) error {
	// Use timeout instead of channel
	if collector.UseTimeout {
		// Create timeout context
		c, cancel := context.WithTimeout(context.Background(), collector.Timeout)
		defer cancel()

		sel:
		select {
			case msg := <-ctx.LastMessage:
				// Cancel
				if msg.Timestamp >= ctx.Msg.Timestamp && msg.Author.ID == ctx.Msg.Author.ID && strings.ToLower(msg.Content) == "c" {
					return errors.New("collector canceled")
				}

				if msg.Timestamp >= ctx.Msg.Timestamp && collector.Filter(ctx, msg) {
					if collector.EndAfterOne {
						collector.MessagesCollected = append(collector.MessagesCollected, msg)
						return nil
					}
					collector.MessagesCollected = append(collector.MessagesCollected, msg)
					goto sel
				} else {
					goto sel
				}
			case <-c.Done():
				if !collector.EndAfterOne {
					return nil
				}
				return c.Err()
			}
	}

	return nil
}

// NewDB opens the database
func NewDB() *Database {
	// Make sure file exists, if not create it
	if _, err := os.Stat("../users.yml"); err != nil {
		ioutil.WriteFile("../users.yml", []byte{}, 666)
	}
	// Open the database for reading
	db := &Database{}
	d, _ := ioutil.ReadFile("../users.yml")
	err := yaml.Unmarshal(d, &db)
	Die(err)

	return db
}
// Write writes to the database
func (db *Database) Write() {
	d, err := yaml.Marshal(&db)
	Die(err)
	err = ioutil.WriteFile("../users.yml", d, 666)
	Die(err)
}
// AddShowToDatabase adds a show to the database
func AddShowToDatabase(show *DBShow, userID string) {
	db := *NewDB()

	// Attempt to find user, if not found add them
	if _, ok := db[userID]; !ok {
		db[userID] = []*DBShow{show}
		db.Write()
		return
	}

	// Add to user shows
	db[userID] = append(db[userID], show)
	db.Write()
	return
}