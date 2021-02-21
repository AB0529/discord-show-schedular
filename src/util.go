package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gookit/color"
	"math/rand"
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
		if e, ok := err.(string); !ok {
			err = errors.New(e)
		}
	}

	if err != nil {
		fmt.Printf("[%s] - %s\n", Yellow.Sprint("WARN"), err)
	}
}

// Die handles errors and exits
func Die(err interface{}) {
	// Handle strings being passed by creating an error type
	if err != nil {
		if e, ok := err.(string); !ok {
			err = errors.New(e)
		}
	}

	if err != nil {
		fmt.Printf("[%s] - %s\n", Red.Sprint("ERR"), err)
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
func (ctx *Context) NewEmbed(content string) {
	_, err := ctx.Session.ChannelMessageSendComplex(ctx.Msg.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Description: content,
			Color:       rand.Intn(10000000),
		},
	})
	Warn(err)
}

// Send sends a message to the same channel as the command
func (ctx *Context) Send(content string) {
	_, err := ctx.Session.ChannelMessageSend(ctx.Msg.ChannelID, content)
	Warn(err)
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