package main

import (
	"github.com/bwmarrin/discordgo"
)

type Context struct {
	Bot               *Bot
	Message           *discordgo.Message
	Channel           *discordgo.Channel
	Guild             *discordgo.Guild
	Author            *discordgo.User
	Prefix            string
	RegisterResponses bool
	Commands          []*Command
}

// experimental, arguments may change
// TODO: ChannelMessageSendComplex implementation
func (ctx *Context) Send(channelID, content string) (*discordgo.Message, error) {
	if ctx.RegisterResponses {
		// register with redis
	}
	return ctx.Bot.dgSession.ChannelMessageSend(channelID, content)
}

// experimental, arguments may change
func (ctx *Context) React(emoji string) (*discordgo.MessageReaction, error) {
	if ctx.RegisterResponses {
		// register with redis
	}
	// add reaction
	return nil, nil
}

func MakeContext(bot *Bot, msg *discordgo.Message, prefix string) (*Context, error) {
	channel, err := bot.dgSession.State.Channel(msg.ChannelID)
	if err != nil {
		log.Debug("Channel with id %d not found", msg.ChannelID)
		return nil, err
	}

	var guild *discordgo.Guild

	if msg.GuildID != "" {
		guild, err = bot.dgSession.State.Guild(msg.GuildID)
		if err != nil {
			log.Debug("Guild with id %d not found", msg.ChannelID)
			return nil, err
		}
	}

	return &Context{
		bot,
		msg,
		channel,
		guild,
		msg.Author,
		prefix,
		true,
		[]*Command{},
	}, nil
}
