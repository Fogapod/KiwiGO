package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Fogapod/KiwiGO/bot"
	"github.com/Fogapod/KiwiGO/commandhandler"
	"github.com/Fogapod/KiwiGO/logger"
	"github.com/bwmarrin/discordgo"
)

var (
	log = logger.GetLogger()
)

func main() {
	bot := bot.Bot{}
	handler := commandhandler.NewCommandHandler(&bot)

	bot.Logger = log

	var err error

	config, err := readConfig()
	if err != nil {
		log.Fatal("Failed to read config file, error:\n%s", err)
		bot.Stop(1, true)
	}

	err = log.SetLoggingLevel(config.LoggingLevel)
	if err != nil {
		log.Fatal("Failed to set logging level, error:\n%s", err)
		bot.Stop(1, true)
	}

	log.Info("Initializing bot")
	log.Trace("Creating session")
	bot.Session, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("Failed to create Discord session, error:\n%s", err)
		bot.Stop(1, true)
	}

	bot.Session.AddHandlerOnce(handler.Ready)

	log.Trace("Openning connection")
	if err = bot.Session.Open(); err != nil {
		log.Fatal("Failed to open connection, error:\n%s", err)
		bot.Stop(1, true)
	}

	log.Trace("Registering events")
	bot.Session.AddHandler(handler.HandleMessage)

	bot.InitPrefixes(config.Prefix)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Stop(0, false)
}
