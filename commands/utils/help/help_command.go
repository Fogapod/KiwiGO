package help

import (
	"github.com/Fogapod/KiwiGO/command"
	"github.com/Fogapod/KiwiGO/context"
)

// Build builds command
func Build(base *command.Command) error {
	base.UsageDoc = "{prefix}{aliases} [command]"
	base.ShortDoc = "Get info about given command or list of all commands"
	base.CallFunc = Call
	base.Build()

	return nil
}

// Call calls command
func Call(c *command.Command, ctx *context.Context) (string, error) {
	// TODO: c.commandMap.Unique() ?
	uniqueCOmmands := map[*command.Command]string{}
	for alias, command := range *c.CommandMap {
		uniqueCOmmands[command] = alias
	}

	commands := ""
	for command := range uniqueCOmmands {
		commands += command.Name + " " // trailing space ...
	}

	// very temporary solution, TODO
	if ctx.Argc() == 2 {
		cmd, ok := (*c.CommandMap)[ctx.Arg(1)]
		if !ok {
			return "Command not found", nil
		}

		return cmd.Help(ctx)
	}
	return "Commands: **" + commands + "**\nDefault prefix: **" + ctx.Bot.DefaultPrefixes[0] + "**", nil
}
