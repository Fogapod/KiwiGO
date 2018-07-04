package bot

import (
	"os"

	"github.com/Fogapod/KiwiGO/logger"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Logger          *logger.Logger
	Session         *discordgo.Session
	DefaultPrefixes []string
	GuildPrefixes   map[string]string // id:prefix
}

func (bot *Bot) InitPrefixes(defaultPrefix string) {
	bot.DefaultPrefixes = append(bot.DefaultPrefixes, defaultPrefix, "<@"+bot.Session.State.User.ID+">", "<@!"+bot.Session.State.User.ID+">")
}

func (b *Bot) GetAllChannels() []*discordgo.Channel {
	var channels []*discordgo.Channel

	for _, g := range b.Session.State.Guilds {
		channels = append(channels, g.Channels...)
	}

	// TODO: DM channels

	return channels
}

func (b *Bot) GetAllEmojis() []*discordgo.Emoji {
	var emojis []*discordgo.Emoji

	for _, g := range b.Session.State.Guilds {
		emojis = append(emojis, g.Emojis...)
	}

	return emojis
}

func (b *Bot) GetAllUsers() []*discordgo.User {
	// use map to avoid duplicates
	userMap := make(map[string]*discordgo.User)

	for _, g := range b.Session.State.Guilds {
		for _, m := range g.Members {
			userMap[m.User.ID] = m.User
		}
	}

	users := make([]*discordgo.User, len(userMap))

	var i int
	for _, u := range userMap {
		users[i] = u
	}

	return users
}

func (bot *Bot) Stop(exitCode int, forceStop bool) {
	bot.Logger.Info("Closing connection and exiting with code %d", exitCode)

	if !forceStop {
		bot.Session.Close()
	}

	os.Exit(exitCode)
}
