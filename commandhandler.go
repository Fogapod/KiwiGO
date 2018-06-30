package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Fogapod/KiwiGO/bot"
	"github.com/Fogapod/KiwiGO/context"
	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	Bot *bot.Bot
}

func (h *CommandHandler) getPrefix(content string) string {
	// TODO: guild prefix override

	lowerContent := strings.ToLower(content)

	for _, p := range h.Bot.DefaultPrefixes {
		if strings.HasPrefix(lowerContent, p) {
			return p
		}
	}

	return ""
}

func (h *CommandHandler) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	prefix := h.getPrefix(m.Content)

	isDM := m.GuildID == ""

	if prefix == "" {
		if !isDM { // empty prefix is allowed for direct messages
			return
		}
	}

	// TODO: blacklist guild/user/(channel?)

	// TODO: argument parser, flag parser
	args := strings.Fields(m.Content[len(prefix):])

	// TODO: actual handler
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

	ctx, err := context.MakeContext(h.Bot, m.Message, prefix)
	if err != nil {
		log.Debug("Failed to create context")
		return
	}

	switch command {
	case "help":
		ctx.Send(m.ChannelID, "Commands: help, ping, uptime, user\nPrefix: **"+h.Bot.DefaultPrefixes[0]+"**")
	case "ping":
		pingMessage, err := ctx.Send(m.ChannelID, "Pinging...")
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
		ctx.Send(m.ChannelID, "Todo")
	case "user":
		// TODO: getUser
		// TODO: other get functions (utilize?)
		ctx.Send(m.ChannelID, "Todo")
	}
}
