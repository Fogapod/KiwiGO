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

// Handler of commands
type CommandHandler struct {
	Bot        *bot.Bot
	CommandMap map[string]*command.Command
}

// Returns new CommandHandler
func New(b *bot.Bot) CommandHandler {
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

// Handles ready event
func (h *CommandHandler) HandleReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Trace("Loading commands")
	h.LoadCommands(true)

	h.Bot.InitPrefixes(h.Bot.Config.Prefix)

	log.Info("%s is ready to serve %d guilds", s.State.User, len(s.State.Guilds))

	if len(h.Bot.DefaultPrefixes) == 0 {
		log.Warn("Bot is ready, but prefix list is empty")
	} else {
		log.Info("Default prefix: %s", h.Bot.DefaultPrefixes[0])
	}
}

// Handles new message event
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
		ctx.Send(m.ChannelID, "Error occurred running command **"+cmd.Name+"**. Developer was notified") // developer wasn't notified, TODO
	}

	if response != "" {
		ctx.Send(m.ChannelID, response)
	}
}

// Handles message update event
func (h *CommandHandler) HandleMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.Author.Bot {
		return
	}

	/* Not implemented yet: https://github.com/bwmarrin/discordgo/pull/545
	if m.Content == m.BeforeUpdate.Content {
		return
	}
	*/

	// TODO: delete bot's response
	h.HandleMessage(s, &discordgo.MessageCreate{m.Message})
}

// Handles message delete event
func (h *CommandHandler) HandleMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	// TODO: delete bot's response
}
