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
	config    Config
	log       = getLogger()
	prefixes  []string
	dgSession *discordgo.Session
)

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatal("Failed to read config file, error:\n%s", err)
		stopBot(1, true)
	}

	log.LoggingLevel = config.LoggingLevel

	log.Trace("Creating session")
	dgSession, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("Failed to create Discord session, error:\n%s", err)
		stopBot(1, true)
	}

	log.Trace("Openning connection")
	if err = dgSession.Open(); err != nil {
		log.Fatal("Failed to open connection, error:\n%s", err)
		stopBot(1, true)
	}

	prefixes = append(prefixes, config.Prefix, "<@"+dgSession.State.User.ID+">", "<@!"+dgSession.State.User.ID+">")

	log.Trace("Registering events")
	dgSession.AddHandler(messageCreate)

	log.Info("Bot is ready")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	stopBot(0, false)
}

func stopBot(exitCode int, forceStop bool) {
	log.Info("Closing connection and exiting with code %d", exitCode)

	if !forceStop {
		dgSession.Close()
	}
	os.Exit(exitCode)
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

	isDM := m.GuildID == ""

	if prefix == "" {
		if !isDM { // empty prefix is  allowed for direct messages
			return
		}
	}

	// TODO: blacklist guild/user/(channel?)

	// TODO: argument parser, flag parser
	args := strings.Fields(m.Content[len(prefix):])

	// TODO: command handler
	// TODO: response registration using redis
	command := strings.ToLower(strings.ToLower(args[0]))

	var commandUseLocation string

	if isDM {
		commandUseLocation = "DM"
	} else {
		commandUseLocation = m.GuildID
	}

	// temporary, command handler needed, dispatching commandUse event after checks
	log.Trace("%s in %s -> %s", m.Author.ID, commandUseLocation, command)

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
