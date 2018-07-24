package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Fogapod/KiwiGO/bot"
	"github.com/Fogapod/KiwiGO/commandhandler"
	"github.com/Fogapod/KiwiGO/config"
	"github.com/Fogapod/KiwiGO/logger"
	"github.com/bwmarrin/discordgo"
)

//go:generate go run gencommands/main.go

var (
	log = logger.GetLogger()
)

func main() {
	bot := bot.New()
	bot.Logger = log

	handler := commandhandler.New(bot)

	var err error

	bot.Config, err = config.ReadConfig()
	if err != nil {
		log.Fatal("Failed to read config file, error:\n%s", err)
		bot.Stop(1, true)
	}

	err = log.SetLoggingLevel(bot.Config.LoggingLevel)
	if err != nil {
		log.Fatal("Failed to set logging level, error:\n%s", err)
		bot.Stop(1, true)
	}

	log.Info("Initializing bot")
	log.Trace("Creating session")
	bot.Session, err = discordgo.New("Bot " + bot.Config.Token)
	if err != nil {
		log.Fatal("Failed to create Discord session, error:\n%s", err)
		bot.Stop(1, true)
	}

	log.Trace("Registering events")
	bot.Session.AddHandlerOnce(handler.HandleReady)
	bot.Session.AddHandler(handler.HandleMessage)
	bot.Session.AddHandler(handler.HandleMessageUpdate)
	bot.Session.AddHandler(handler.HandleMessageDelete)

	log.Trace("Openning connection")
	if err = bot.Session.Open(); err != nil {
		log.Fatal("Failed to open connection, error:\n%s", err)
		bot.Stop(1, true)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Trace("Exit code received")

	bot.Stop(0, false)
}
