package user

import (
	"github.com/Fogapod/KiwiGO/command"
	"github.com/Fogapod/KiwiGO/context"
	"github.com/Fogapod/KiwiGO/utils/finders"
	"github.com/bwmarrin/discordgo"
)

func Build(base *command.Command) error {
	base.CallFunc = Call
	base.Build()

	return nil
}

func Call(c *command.Command, ctx *context.Context) (string, error) {
	var user *discordgo.User

	if ctx.Argc() < 2 {
		user = ctx.Message.Author
	} else {
		users := finders.FindUser(ctx.Args(1, ctx.Argc()), ctx, &finders.FindUserOptions{})
		if len(users) == 0 {
			return "User not found", nil
		}
		user = users[0].User
	}

	return user.String(), nil
}
