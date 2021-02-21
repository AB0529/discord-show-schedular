package main

// Command the representation of a bot command
type Command struct {
	Name    string
	Aliases []string
	Example []string
	Desc    []string
	Handler func(*Context)
	Flags   []string
}

// Ping command which returns a message
func Ping(ctx *Context) {
	ctx.NewEmbed("Pong nigga")
}

// Schedule hub command which handles actions for the user's schedule
func Schedule(ctx *Context) {
	args := ctx.FindCommandFlag()
	if args == nil {
		ctx.SendCommandHelp()
	}
}
