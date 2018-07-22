package command

import (
	"github.com/Fogapod/KiwiGO/bot"
	"github.com/Fogapod/KiwiGO/context"
)

// TODO: interraction with handler somehow

type AllowedChatType int

const (
	AllowedChatTypeAll AllowedChatType = iota
	AllowedChatTypeGuildOnly
	AllowedChatTypeDMOnly
)

type Command struct {
	UsageDoc string
	ShortDoc string
	LongDoc  string

	Name     string
	Aliases  []string
	Category string // *Category ?
	// bot_perms  []*Permission
	// user_perms []*Permission
	MinArgs int
	// flags      []*FlagOptions
	AllowedChatType AllowedChatType
	Nsfw            bool
	Hidden          bool
	Disabled        bool
	// events     ???
	IsSubCommand bool
	parent       *Command
	children     map[string]*Command

	CallFunc   func(*Command, *context.Context) (string, error)
	HelpFunc   func(*Command, *context.Context) (string, error)
	UnloadFUnc func(*Command) error

	Bot        *bot.Bot
	CommandMap *map[string]*Command // TODO: pointer to habdler
}

func New(b *bot.Bot, name string, commands *map[string]*Command) *Command {
	return &Command{
		Name:            name,
		AllowedChatType: AllowedChatTypeAll,
		children:        map[string]*Command{},
		Bot:             b,
		CommandMap:      commands,
	}
}

func NewSubcommand(b *bot.Bot, name string, commands *map[string]*Command, parent *Command) *Command {
	cmd := &Command{
		Name:            name,
		AllowedChatType: AllowedChatTypeAll,
		IsSubCommand:    true,
		parent:          parent,
		children:        map[string]*Command{},
		Bot:             b,
		CommandMap:      commands,
	}

	return cmd
}

func (c *Command) Build() {
	if c.ShortDoc == "" {
		c.UsageDoc = "{prefix}" + c.Name
	}

	c.Aliases = append(c.Aliases, c.Name)

	if c.Category == "" {
		c.Category = "No category"
	}

	if c.CallFunc == nil {
		c.CallFunc = DefaultCall
	}

	if c.HelpFunc == nil {
		c.HelpFunc = DefaultHelp
	}

	if c.UnloadFUnc == nil {
		c.UnloadFUnc = DefaultUnload
	}

	// TODO: build subcommands
}

// Wrappers

func (c *Command) Call(ctx *context.Context) (string, error) {
	return c.CallFunc(c, ctx)
}

func (c *Command) Help(ctx *context.Context) (string, error) {
	return c.HelpFunc(c, ctx)
}

func (c *Command) Unload() error {
	return c.UnloadFUnc(c)
}

// Default funcs

func DefaultCall(c *Command, ctx *context.Context) (string, error) {
	return "", nil
}

func DefaultHelp(c *Command, ctx *context.Context) (string, error) {
	// TODO: complete formatting, add flag/aliases/local prefix support
	// TODO: use embed
	return ctx.Prefix + c.Name, nil
}

func DefaultUnload(c *Command) error {
	return nil
}

// TODO: check funcs
// TODO: errors
