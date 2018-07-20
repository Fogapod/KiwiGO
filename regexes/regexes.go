package regexes

import (
	"regexp"
)

var (
	IDExpr             = `\d{17,19}`
	UserMentionExpr    = `<@!?(` + IDExpr + `)>`
	RoleMentionExpr    = `<@&(` + IDExpr + `)>`
	ChannelMentionExpr = `<#(` + IDExpr + `)>`
	EmojiExpr          = `^<(?P<animated>a?):(?P<name>[_a-zA-Z]{2,32}):(?P<id>` + IDExpr + `)>$`

	ID             = regexp.MustCompile(`^` + IDExpr + `$`)
	UserMention    = regexp.MustCompile(`^` + UserMentionExpr + `$`)
	RoleMention    = regexp.MustCompile(`^` + RoleMentionExpr + `$`)
	ChannelMention = regexp.MustCompile(`^` + ChannelMentionExpr + `$`)
	Emoji          = regexp.MustCompile(`^` + EmojiExpr + `$`)

	UserMentionOrID    = regexp.MustCompile(`^(?:` + UserMentionExpr + `)|^` + IDExpr + `$`)
	RoleMentionOrID    = regexp.MustCompile(`^(?:` + RoleMentionExpr + `)|^` + IDExpr + `$`)
	ChannelMentionOrID = regexp.MustCompile(`^(?:` + ChannelMentionExpr + `)|^` + IDExpr + `$`)
)
