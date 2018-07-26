package ping

import (
	"fmt"
	"os/exec"

	"github.com/Fogapod/KiwiGO/command"
	"github.com/Fogapod/KiwiGO/context"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/idna"
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
func Call(c *command.Command, ctx *context.Context) (string, error) {
	pingMessage, err := ctx.Send("Pinging...")
	if err != nil {
		return "", err
	}

	var pingTarget string

	if ctx.Argc() > 1 {
		pingTarget = ctx.Args(1, ctx.Argc())
		encoded, err := idna.ToASCII(pingTarget)
		if err == nil {
			cmd := exec.Command("ping", encoded, "-c", "3")
			ret, err := cmd.CombinedOutput()
			if err == nil {
				ctx.Edit(pingMessage.ID, "```\n"+string(ret)+"```")
				return "", nil
			}
		}
	}

	var userMessageTimestamp discordgo.Timestamp
	if ctx.Message.EditedTimestamp != "" {
		userMessageTimestamp = ctx.Message.EditedTimestamp
	} else {
		userMessageTimestamp = ctx.Message.Timestamp
	}

	ts1, err := userMessageTimestamp.Parse()
	if err != nil {
		return "", err
	}

	ts2, err := pingMessage.Timestamp.Parse()
	if err != nil {
		return "", err
	}

	response := "Pong, it took **%dms** to "
	if pingTarget == "" {
		response += "respond"
	} else {
		response += "ping **" + pingTarget + "**"
	}

	delta := int(ts2.Sub(ts1).Seconds() * 1000)
	ctx.Edit(pingMessage.ID, fmt.Sprintf(response, delta))

	return "", nil
}
