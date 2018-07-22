package context

import (
	"strings"

	"github.com/Fogapod/KiwiGO/bot"
	"github.com/bwmarrin/discordgo"
)

type Context struct {
	Bot               *bot.Bot
	Session           *discordgo.Session
	Message           *discordgo.Message
	Channel           *discordgo.Channel
	Guild             *discordgo.Guild
	Author            *discordgo.User
	Prefix            string
	args              []string // TODO: separate arg type
	separators        []string // TODO: separate arg type
	RegisterResponses bool     // delete?
	// Flags          map[string]*Flag
}

// experimental, arguments may change
// TODO: ChannelMessageSendComplex implementation
func (ctx *Context) Send(channelID, content string) (*discordgo.Message, error) {
	if ctx.RegisterResponses {
		// register with redis
	}
	return ctx.Bot.Session.ChannelMessageSend(channelID, content)
}

// experimental, arguments may change
func (ctx *Context) React(emoji string) (*discordgo.MessageReaction, error) {
	if ctx.RegisterResponses {
		// register with redis
	}
	// add reaction
	return nil, nil
}

func New(b *bot.Bot, s *discordgo.Session, msg *discordgo.Message, prefix string) (*Context, error) {
	channel, err := b.Session.State.Channel(msg.ChannelID)
	if err != nil {
		b.Logger.Debug("Channel with id %d not found", msg.ChannelID)
		return nil, err
	}

	var guild *discordgo.Guild

	if msg.GuildID != "" {
		guild, err = b.Session.State.Guild(msg.GuildID)
		if err != nil {
			b.Logger.Debug("Guild with id %d not found", msg.ChannelID)
			return nil, err
		}
	}

	return &Context{
		Bot:               b,
		Session:           s,
		Message:           msg,
		Channel:           channel,
		Guild:             guild,
		Author:            msg.Author,
		Prefix:            prefix,
		RegisterResponses: true,
	}, nil
}

// TODO: arg parser
// TODO: flag parser
func (ctx *Context) ParseContent(content string) error {
	ctx.args = strings.Fields(content[len(ctx.Prefix):])

	return nil
}

func (ctx *Context) Arg(index int) string {
	// TODO: prevent segfault
	return ctx.args[index]
}

func (ctx *Context) Args(begin, end int) string {
	// TODO: prevent segfault
	return strings.Join(ctx.args[begin:end], " ")
}

func (ctx *Context) Argc() int {
	return len(ctx.args)
}

func (ctx *Context) ArgArray() []string {
	args := make([]string, len(ctx.args))

	copy(args, ctx.args)

	return args
}
