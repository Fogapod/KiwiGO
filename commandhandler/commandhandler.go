package commandhandler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Fogapod/KiwiGO/bot"
	"github.com/Fogapod/KiwiGO/command"
	"github.com/Fogapod/KiwiGO/context"
	"github.com/Fogapod/KiwiGO/logger"
	"github.com/bwmarrin/discordgo"
)

var (
	log = logger.GetLogger()
)

type CommandHandler struct {
	Bot      *bot.Bot
	commands map[string]*command.Command
}

func NewCommandHandler(b *bot.Bot) CommandHandler {
	return CommandHandler{b, map[string]*command.Command{}}
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

func (h *CommandHandler) Ready(s *discordgo.Session, r *discordgo.Ready) {
	log.Info("Bot is ready to serve %d guilds", len(h.Bot.Session.State.Guilds))
}

func (h *CommandHandler) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	if len(args) == 0 {
		return
	}

	// TODO: actual handler
	// TODO: response registration using redis
	// command := h.commands[strings.ToLower(args[0])]
	command := strings.ToLower(args[0])

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
		h.Bot.Logger.Debug("Failed to create context")
		return
	}

	// switch command.Name {
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