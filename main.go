package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	prefixes []string
)

type Config struct {
	Token  string `json: "token"`
	Prefix string `json: "prefix"`
}

// TODO: move to separate file
func readConfig() (Config, error) {
	var c Config

	raw, err := ioutil.ReadFile("./config.json")
	if err == nil {
		err = json.Unmarshal(raw, &c)
		if err == nil {
			return c, err
		}
	}
	return c, err
}

func main() {
	config, err := readConfig()
	if err != nil {
		fmt.Println("Fatal: failed to read config file. Exiting")
		os.Exit(1)
	}

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Fatal: creating Discord session,", err)
		return
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("Fatal: opening connection,", err)
		return
	}

	prefixes = append(prefixes, config.Prefix, "<@"+dg.State.User.ID+">", "<@!"+dg.State.User.ID+">")

	dg.AddHandler(messageCreate)

	fmt.Println("Info: Bot is ready") // TODO: logger

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
			fmt.Println("Channel with id %d not found", m.ChannelID)
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
	switch strings.ToLower(strings.ToLower(args[0])) {
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

		delta := id2>>22 - id1>>22 // TODO: use timestamps

		s.ChannelMessageEdit(m.ChannelID, pingMessage.ID, fmt.Sprintf("Pong, it took **%dms** to respond", delta))
	case "uptime":
		s.ChannelMessageSend(m.ChannelID, "Todo")
	case "user":
		// TODO: getUser
		// TODO: other get functions (utilize?)
		s.ChannelMessageSend(m.ChannelID, "Todo")
	}
}
