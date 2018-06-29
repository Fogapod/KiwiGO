package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	log = getLogger()
)

func main() {
	bot := Bot{}
	handler := CommandHandler{&bot}

	var err error

	bot.Config, err = readConfig()
	if err != nil {
		log.Fatal("Failed to read config file, error:\n%s", err)
		bot.Stop(1, true)
	}

	log.LoggingLevel = bot.Config.LoggingLevel

	log.Trace("Creating session")
	bot.dgSession, err = discordgo.New("Bot " + bot.Config.Token)
	if err != nil {
		log.Fatal("Failed to create Discord session, error:\n%s", err)
		bot.Stop(1, true)
	}

	log.Trace("Openning connection")
	if err = bot.dgSession.Open(); err != nil {
		log.Fatal("Failed to open connection, error:\n%s", err)
		bot.Stop(1, true)
	}

	log.Trace("Registering events")
	bot.dgSession.AddHandler(handler.messageCreate)

	bot.InitPrefixes()

	log.Info("Bot is ready")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Stop(0, false)
}
