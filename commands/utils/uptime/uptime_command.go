package uptime

import (
	"github.com/Fogapod/KiwiGO/command"
	"github.com/Fogapod/KiwiGO/context"
)

func Build(base *command.Command) error {
	base.CallFunc = Call
	base.Build()

	return nil
}

func Call(c *command.Command, ctx *context.Context) (string, error) {
	return "TODO", nil
}
