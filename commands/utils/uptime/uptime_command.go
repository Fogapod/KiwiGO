package uptime

import (
	"strconv"
	"time"

	"github.com/Fogapod/KiwiGO/command"
	"github.com/Fogapod/KiwiGO/context"
)

var startTime = time.Now()

func Build(base *command.Command) error {
	base.CallFunc = Call
	base.Aliases = append(base.Aliases, "up")
	base.Build()

	return nil
}

func Call(c *command.Command, ctx *context.Context) (string, error) {
	deltaSEconds := int(time.Now().Sub(startTime).Seconds())
	divmod := func(numerator, denominator int) (quotient, remainder int) {
		return numerator / denominator, numerator % denominator
	}

	niceTime := ""

	minutes, seconds := divmod(deltaSEconds, 60)
	hours, minutes := divmod(minutes, 60)
	days, hours := divmod(hours, 24)
	months, days := divmod(days, 30)
	years, months := divmod(months, 12)

	if years > 0 {
		niceTime += strconv.Itoa(years) + "y "
	}
	if months > 0 {
		niceTime += strconv.Itoa(months) + "mon "
	}
	if days > 0 {
		niceTime += strconv.Itoa(days) + "d "
	}
	if hours > 0 {
		niceTime += strconv.Itoa(hours) + "h "
	}
	if minutes > 0 {
		niceTime += strconv.Itoa(minutes) + "m "
	}
	if seconds > 0 {
		niceTime += strconv.Itoa(seconds) + "s"
	}

	return "Bot is running for **" + niceTime + "**", nil
}
