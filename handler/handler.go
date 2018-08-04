package handler

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
type Handler struct {
	Bot               *bot.Bot
	CommandMap        map[string]*command.Command
	isFirstReadyEvent bool
}

// New returns new CommandHandler
func New(b *bot.Bot) Handler {
	return Handler{
		Bot:               b,
		CommandMap:        map[string]*command.Command{},
		isFirstReadyEvent: true,
	}
}

func (h *Handler) getPrefix(content string) string {
	// TODO: guild prefix override

	lowerContent := strings.ToLower(content)

	for _, p := range h.Bot.DefaultPrefixes {
		if strings.HasPrefix(lowerContent, p) {
			return p
		}
	}

	return ""
}
