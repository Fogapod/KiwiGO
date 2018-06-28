package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	log      = getLogger()
	prefixes []string
)

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatal("Failed to read config file, error: \n%s", err)
		os.Exit(1)
	}

	log.LoggingLevel = config.LoggingLevel

	log.Trace("Creating session")
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("Failed to create Discord session, error: \n%s", err)
		return
	}

	log.Trace("Openning connection")
	if err = dg.Open(); err != nil {
		log.Fatal("Failed to open connection, error: \n%s", err)
		return
	}

	prefixes = append(prefixes, config.Prefix, "<@"+dg.State.User.ID+">", "<@!"+dg.State.User.ID+">")

	log.Trace("Registering events")
	dg.AddHandler(messageCreate)

	log.Info("Info: Bot is ready")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func checkPrefixes(content string) string {
	// TODO: guild prefix override

	lowerContent := strings.ToLower(content)

	for _, p := range prefixes {
		if strings.HasPrefix(lowerContent, p) {
			return p
		}
	}

	return ""
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	prefix := checkPrefixes(m.Content)

	if prefix == "" {
		channel, err := s.State.Channel(m.ChannelID)
		if err != nil {
			log.Debug("Channel with id %d not found", m.ChannelID)
			return
		}
		if channel.Type != discordgo.ChannelTypeDM { // empty prefix is  allowed for direct messages
			return
		}
	}

	// TODO: blacklist guild/user/(channel?)

	// TODO: argument parser, flag parser
	args := strings.Fields(m.Content[len(prefix):])

	// TODO: command handler
	// TODO: response registration using redis
	command := strings.ToLower(strings.ToLower(args[0]))
	log.Trace("called command %s", command)

	switch command {
	case "help":
		s.ChannelMessageSend(m.ChannelID, "help, ping, uptime, user")
	case "ping":
		pingMessage, err := s.ChannelMessageSend(m.ChannelID, "Pinging...")
		if err != nil {
			return
		}

		id1, err := strconv.ParseInt(m.ID, 10, 64)
		if err != nil {
			return
		}

		id2, err := strconv.ParseInt(pingMessage.ID, 10, 64)
		if err != nil {
			return
		}

		delta := id2>>22 - id1>>22

		s.ChannelMessageEdit(m.ChannelID, pingMessage.ID, fmt.Sprintf("Pong, it took **%dms** to respond", delta))
	case "uptime":
		s.ChannelMessageSend(m.ChannelID, "Todo")
	case "user":
		// TODO: getUser
		// TODO: other get functions (utilize?)
		s.ChannelMessageSend(m.ChannelID, "Todo")
	}
}
