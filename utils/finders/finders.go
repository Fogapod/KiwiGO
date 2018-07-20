package finders

import (
	"sort"
	"strings"
	"time"

	"github.com/Fogapod/KiwiGO/context"
	"github.com/Fogapod/KiwiGO/regexes"
	"github.com/bwmarrin/discordgo"
)

type FoundUser struct {
	Member     *discordgo.Member
	User       *discordgo.User
	matchIndex int
}

type FindUserOptions struct {
	StrictGuild  bool
	GlobalSearch bool
}

func FindUser(pattern string, ctx *context.Context, options *FindUserOptions) []*FoundUser {
	var found []*FoundUser

	mentionMatch := regexes.UserMentionOrID.FindStringSubmatch(pattern)
	if len(mentionMatch) > 0 {
		var userID string
		if len(mentionMatch[1]) != 0 {
			userID = mentionMatch[1]
		} else {
			userID = mentionMatch[0]
		}

		var guilds []*discordgo.Guild
		if options.StrictGuild {
			guilds = append(guilds, ctx.Guild)
		} else {
			guilds = ctx.Session.State.Guilds
		}

		for _, guild := range guilds {
			member, err := ctx.Session.State.Member(guild.ID, userID)
			if err == nil {
				return append(found, &FoundUser{
					Member: member,
					User:   member.User,
				})
			}
		}

		if !options.StrictGuild {
			user, err := ctx.Session.User(userID)
			if err == nil {
				return append(found, &FoundUser{
					User: user,
				})
			}
		}
	}

	lowerPattern := strings.ToLower(pattern)

	var guilds []*discordgo.Guild
	if options.GlobalSearch {
		guilds = ctx.Session.State.Guilds
	} else {
		guilds = append(guilds, ctx.Guild)
	}

	for _, guild := range guilds {
		for _, member := range guild.Members {
			matchIndex := -1

			if options.GlobalSearch {
				matchIndex = strings.Index(strings.ToLower(member.Nick), lowerPattern)
			}
			if matchIndex == -1 {
				matchIndex = strings.Index(strings.ToLower(member.User.String()), lowerPattern)
			}
			if matchIndex == -1 {
				continue
			}

			found = append(found, &FoundUser{
				Member:     member,
				User:       member.User,
				matchIndex: matchIndex,
			})

			sort.Slice(found, func(i, j int) bool {
				lastMessageTimestampDelta := ctx.Bot.GetLatsUserMessageTimestamp(ctx.Channel.ID, found[j].User.ID).Sub(ctx.Bot.GetLatsUserMessageTimestamp(ctx.Channel.ID, found[i].User.ID))
				if lastMessageTimestampDelta < 0 {
					return true
				} else if lastMessageTimestampDelta > 0 {
					return false
				}

				matchIndexDelta := found[i].matchIndex - found[j].matchIndex
				if matchIndexDelta < 0 {
					return true
				} else if matchIndexDelta > 0 {
					return false
				}

				joinedAt1, err := time.Parse(time.RFC3339, found[i].Member.JoinedAt)
				if err != nil {
					return false
				}
				joinedAt2, err := time.Parse(time.RFC3339, found[j].Member.JoinedAt)
				if err != nil {
					return false
				}

				return joinedAt1.Sub(joinedAt2) > 0
			})
		}
	}

	return found
}
