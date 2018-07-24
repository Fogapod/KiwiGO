package bot

import (
	"os"
	"time"

	"github.com/Fogapod/KiwiGO/config"
	"github.com/Fogapod/KiwiGO/logger"
	"github.com/bwmarrin/discordgo"
)

// Main bot struct
type Bot struct {
	Logger            *logger.Logger
	Config            *config.Config
	Session           *discordgo.Session
	DefaultPrefixes   []string
	GuildPrefixes     map[string]string                         // id:prefix
	MessageTimestamps map[string]map[string]discordgo.Timestamp // channelID: userID: timestmp
}

// Returns new Bot instance
func New() *Bot {
	return &Bot{
		GuildPrefixes:     make(map[string]string),
		MessageTimestamps: make(map[string]map[string]discordgo.Timestamp),
	}
}

//
func (bot *Bot) InitPrefixes(defaultPrefix string) {
	bot.DefaultPrefixes = append(bot.DefaultPrefixes, defaultPrefix, "<@"+bot.Session.State.User.ID+">", "<@!"+bot.Session.State.User.ID+">")
}

// Registers timestamp of last user's message
// TODO: fix data race?
func (bot *Bot) RegisterMessageTimestamp(m *discordgo.Message) {
	if _, ok := bot.MessageTimestamps[m.ChannelID]; !ok {
		bot.MessageTimestamps[m.ChannelID] = map[string]discordgo.Timestamp{m.Author.ID: m.Timestamp}
	} else {
		bot.MessageTimestamps[m.ChannelID][m.Author.ID] = m.Timestamp
	}
}

// Gets timestamp of last user's message
// TODO: fix data race?
func (bot *Bot) GetLatsUserMessageTimestamp(channelID, userID string) time.Time {
	channel, ok := bot.MessageTimestamps[channelID]
	if ok {
		timestamp, ok := channel[userID]
		if ok {
			ts, _ := timestamp.Parse()
			return ts
		}
	}

	return time.Unix(0, 0)
}

// Returns all cached channels
func (b *Bot) GetAllChannels() (channels []*discordgo.Channel) {
	for _, g := range b.Session.State.Guilds {
		channels = append(channels, g.Channels...)
	}

	// TODO: DM channels

	return channels
}

// Returns all cached emojis
func (b *Bot) GetAllEmojis() (emojis []*discordgo.Emoji) {
	for _, g := range b.Session.State.Guilds {
		emojis = append(emojis, g.Emojis...)
	}

	return emojis
}

// Returns all cached users
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
		i++
	}

	return users
}

// Stops bot
func (bot *Bot) Stop(exitCode int, forceStop bool) {
	bot.Logger.Info("Closing connection and exiting with code %d", exitCode)

	if !forceStop {
		bot.Session.Close()
	}

	os.Exit(exitCode)
}
