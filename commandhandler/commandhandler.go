package commandhandler

import (
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
	Bot        *bot.Bot
	CommandMap map[string]*command.Command
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

func (h *CommandHandler) HandleReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Trace("Loading commands")
	h.LoadCommands(true)

	log.Info("%s is ready to serve %d guilds", s.State.User, len(s.State.Guilds))

	if len(h.Bot.DefaultPrefixes) == 0 {
		log.Warn("Bot is ready, but prefix list is empty")
	} else {
		log.Info("Default prefix: %s", h.Bot.DefaultPrefixes[0])
	}
}

func (h *CommandHandler) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	h.Bot.RegisterMessageTimestamp(m.Message)

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

	ctx, err := context.New(h.Bot, s, m.Message, prefix)
	if err != nil {
		log.Warn("Failed to create context:\n%s", err)
		return
	}

	// TODO: argument parser, flag parser
	err = ctx.ParseContent(m.Content)
	if err != nil {
		log.Debug("Failed to parse arguments:\n%s", err)
		return
	}
	args := strings.Fields(m.Content[len(prefix):])

	if len(args) == 0 {
		return
	}

	// TODO: actual handler
	// TODO: response registration using redis
	cmd, found := h.CommandMap[strings.ToLower(ctx.Arg(0))]
	if !found {
		return
	}

	var commandUseLocation string

	if isDM {
		commandUseLocation = "DM"
	} else {
		commandUseLocation = m.GuildID
	}

	log.Debug("%s in %s -> %s", m.Author.ID, commandUseLocation, cmd.Name)

	response, err := cmd.Call(ctx)
	if err != nil { // TODO: errors, error handler
		log.Warn("Error running command %s:\n%s", cmd.Name, err)
		ctx.Send(m.ChannelID, "Error occured running command **"+cmd.Name+"**. Developer was notified") // developer wasn't notified, TODO
	}

	if response != "" {
		ctx.Send(m.ChannelID, response)
	}
}
