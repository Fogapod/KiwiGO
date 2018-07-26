package ping

import (
	"fmt"

	"github.com/Fogapod/KiwiGO/command"
	"github.com/Fogapod/KiwiGO/context"
	"github.com/bwmarrin/discordgo"
)

// Build builds command
func Build(base *command.Command) error {
	base.UsageDoc = "{prefix}{aliases} [target]"
	base.ShortDoc = "Ping ip / check bot's response delay"
	base.CallFunc = Call
	base.Build()

	return nil
}

// Call calls command
func Call(c *command.Command, ctx *context.Context) (response string, err error) {
	pingMessage, err := ctx.Send("Pinging...")
	if err != nil {
		return
	}

	// TODO: ping ip

	var userMessageTimestamp discordgo.Timestamp
	if ctx.Message.EditedTimestamp != "" {
		userMessageTimestamp = ctx.Message.EditedTimestamp
	} else {
		userMessageTimestamp = ctx.Message.Timestamp
	}

	ts1, err := userMessageTimestamp.Parse()
	if err != nil {
		return
	}

	ts2, err := pingMessage.Timestamp.Parse()
	if err != nil {
		return
	}

	delta := int(ts2.Sub(ts1).Seconds() * 1000)
	ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, pingMessage.ID, fmt.Sprintf("Pong, it took **%dms** to respond", delta))

	return
}
