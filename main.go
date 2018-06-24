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
	"time"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Token string `json: "token"`
}

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
		fmt.Println("Error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	fmt.Println("Bot is ready")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	switch strings.ToLower(m.Content) {
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

		delta := (id2 - id1) >> 22

		s.ChannelMessageEdit(m.ChannelID, pingMessage.ID, fmt.Sprintf("Pong, it took **%d**ms to respond", delta))
	case "sleep":
		time.Sleep(5 * time.Second)
		s.ChannelMessageSend(m.ChannelID, "Done sleeping")
	}
}
