package formatters

import "strings"

// ReplaceMassMentions replaces @everyone and @here mentions
func ReplaceMassMentions(s *string) *string {
	if s == nil {
		return s
	}

	var formatted string
	formatted = strings.Replace(strings.Replace(*s, "@everyone", "@\u200beveryone", -1), "@here", "@\u200bhere", -1)

	return &formatted
}

// TODO: ReplaceUserMentions
// TODO: TrimText
