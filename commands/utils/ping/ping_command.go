package ping

import (
	"fmt"
	"strconv"

	"github.com/Fogapod/KiwiGO/command"
	"github.com/Fogapod/KiwiGO/context"
)

func Build(base *command.Command) error {
	base.CallFunc = Call
	base.Build()

	return nil
}

func Call(c *command.Command, ctx *context.Context) (string, error) {
	pingMessage, err := ctx.Send(ctx.Message.ChannelID, "Pinging...")
	if err != nil {
		return "", err
	}

	id1, err := strconv.ParseInt(ctx.Message.ID, 10, 64)
	if err != nil {
		return "", err
	}

	id2, err := strconv.ParseInt(pingMessage.ID, 10, 64)
	if err != nil {
		return "", err
	}

	delta := id2>>22 - id1>>22

	ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, pingMessage.ID, fmt.Sprintf("Pong, it took **%dms** to respond", delta))

	return "", nil
}
