package main

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Config          *Config
	dgSession       *discordgo.Session
	DefaultPrefixes []string
	GuildPrefixes   map[string]string // id: prefix
}

func (bot *Bot) InitPrefixes() {
	bot.DefaultPrefixes = append(bot.DefaultPrefixes, bot.Config.Prefix, "<@"+bot.dgSession.State.User.ID+">", "<@!"+bot.dgSession.State.User.ID+">")
}

func (bot *Bot) Stop(exitCode int, forceStop bool) {
	log.Info("Closing connection and exiting with code %d", exitCode)

	if !forceStop {
		bot.dgSession.Close()
	}
	os.Exit(exitCode)
}
