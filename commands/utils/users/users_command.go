package users

import (
	"github.com/Fogapod/KiwiGO/command"
	"github.com/Fogapod/KiwiGO/context"
	"github.com/Fogapod/KiwiGO/utils/finders"
)

// Build command
func Build(base *command.Command) error {
	base.CallFunc = Call
	base.Build()

	return nil
}

// Call command
func Call(c *command.Command, ctx *context.Context) (string, error) {
	if ctx.Argc() < 2 {
		return "Not enough arguments", nil
	} else {
		users, err := finders.FindUser(ctx.Args(1, ctx.Argc()), ctx, &finders.FindUserOptions{
			UseRegex: true, // temporary. TODO: command flags
		})
		if err != nil {
			return "", err
		}

		if len(users) == 0 {
			return "Users not found", nil
		}

		response := ""
		for i, u := range users {
			if i >= 20 {
				break // limited to 20 until paginator implementation
			}
			response += u.User.String() + "\n"
		}

		return "```" + response + "```", nil
	}
}
