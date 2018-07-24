package help

import (
	"github.com/Fogapod/KiwiGO/command"
	"github.com/Fogapod/KiwiGO/context"
)

// Build command
func Build(base *command.Command) error {
	base.CallFunc = Call
	base.Build()

	return nil
}

// Call command
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
	return "Commands: **" + commands + "**\nDefault prefix: **" + ctx.Bot.DefaultPrefixes[0] + "**", nil
}
