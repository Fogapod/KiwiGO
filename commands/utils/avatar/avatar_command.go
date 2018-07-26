package avatar

import (
	"fmt"
	"strings"

	"github.com/Fogapod/KiwiGO/command"
	"github.com/Fogapod/KiwiGO/constants"
	"github.com/Fogapod/KiwiGO/context"
	"github.com/Fogapod/KiwiGO/utils/finders"
	"github.com/bwmarrin/discordgo"
)

// Build builds command
func Build(base *command.Command) error {
	base.UsageDoc = "{prefix}{aliases} [user]"
	base.ShortDoc = "Get user's avatar. Defaults to you"
	base.CallFunc = Call
	base.Aliases = append(base.Aliases, "pfp")
	base.Build()

	return nil
}

// Call calls command
func Call(c *command.Command, ctx *context.Context) (string, error) {
	var user *discordgo.User

	if ctx.Argc() < 2 {
		user = ctx.Message.Author
	} else {
		users, err := finders.FindUser(ctx.Args(1, ctx.Argc()), ctx, &finders.FindUserOptions{})
		if err != nil {
			return "", err
		}

		if len(users) == 0 {
			return "User not found", nil
		}
		user = users[0].User
	}

	avatarURLTemplate := "https://cdn.discordapp.com/avatars/" + user.ID + "/" + user.Avatar + ".%s?size=2048"
	extensions := []string{"png", "webp", "jpg"}

	if strings.HasPrefix(user.Avatar, "a_") { // animated avatar
		extensions = append([]string{"gif"}, extensions...)
	}

	for i, extensionName := range extensions {
		extensions[i] = "[" + extensionName + "](" + fmt.Sprintf(avatarURLTemplate, extensionName) + ")"
	}

	avatarURL := user.AvatarURL("2048")

	embed := discordgo.MessageEmbed{
		Title:       user.String() + " 's avatar",
		Description: strings.Join(extensions, " | "),
		Image:       &discordgo.MessageEmbedImage{URL: avatarURL},
		Color:       constants.GopherCOlor,
	}

	ctx.SendComplex(&discordgo.MessageSend{Embed: &embed})

	return "", nil
}
