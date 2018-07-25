package commandhandler

import (
	"strings"

	"github.com/Fogapod/KiwiGO/bot"
	"github.com/Fogapod/KiwiGO/command"
	"github.com/Fogapod/KiwiGO/logger"
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
