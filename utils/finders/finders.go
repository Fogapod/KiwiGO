package finders

import (
	"regexp"
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
	CurrentGuild  bool // disable global user lookup using id, search by name only current guild
	GlobalSearch  bool // search by name through all guilds bot is in, can be very expensive
	UseRegex      bool // match names using regex
	CaseSensetive bool // use case-sesnetive comparison of user names
}

func FindUser(pattern string, ctx *context.Context, options *FindUserOptions) (found []*FoundUser, err error) {
	mentionMatch := regexes.UserMentionOrID.FindStringSubmatch(pattern)
	if len(mentionMatch) > 0 {
		var userID string
		if len(mentionMatch[1]) != 0 {
			userID = mentionMatch[1]
		} else {
			userID = mentionMatch[0]
		}

		var guilds []*discordgo.Guild
		if !options.CurrentGuild {
			guilds = ctx.Session.State.Guilds
		} else if ctx.Guild != nil { // not DM
			guilds = append(guilds, ctx.Guild)
		}

		// local cache id lookup
		for _, guild := range guilds {
			member, err := ctx.Session.State.Member(guild.ID, userID)
			if err == nil {
				return append(found, &FoundUser{
					Member: member,
					User:   member.User,
				}), nil
			}
		}

		// global id lookup
		if !options.CurrentGuild {
			user, err := ctx.Session.User(userID)
			if err == nil {
				return append(found, &FoundUser{
					User: user,
				}), nil
			}
		}
	}

	var nameCompareFunc func(n, p string) int

	if options.UseRegex {
		var regex *regexp.Regexp

		if options.CaseSensetive {
			regex, err = regexp.Compile(pattern)
		} else {
			regex, err = regexp.Compile(`(?i)` + pattern)
		}
		if err != nil {
			return
		}

		nameCompareFunc = func(n, p string) int {
			match := regex.FindStringIndex(n)
			if match == nil {
				return -1
			}

			return match[0]
		}
	} else {
		if options.CaseSensetive {
			nameCompareFunc = func(n, p string) int {
				return strings.Index(n, p)
			}
		} else {
			nameCompareFunc = func(n, p string) int {
				return strings.Index(strings.ToLower(n), strings.ToLower(p))
			}
		}
	}

	var guilds []*discordgo.Guild
	if options.GlobalSearch && !options.CurrentGuild {
		guilds = ctx.Session.State.Guilds
	} else if ctx.Guild != nil { // not DM
		guilds = append(guilds, ctx.Guild)
	}

	for _, guild := range guilds {
		for _, member := range guild.Members {
			matchIndex := -1

			if options.GlobalSearch {
				matchIndex = nameCompareFunc(member.Nick, pattern)
			}
			if matchIndex == -1 {
				matchIndex = nameCompareFunc(member.User.String(), pattern)
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

	return
}
