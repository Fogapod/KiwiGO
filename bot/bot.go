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

func (bot *Bot) Stop(exitCode int, forceStop bool) {
	bot.Logger.Info("Closing connection and exiting with code %d", exitCode)

	if !forceStop {
		bot.Session.Close()
	}
	os.Exit(exitCode)
}
