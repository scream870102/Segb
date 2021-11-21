package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/scream870102/segb/commands"
	"github.com/scream870102/segb/database"
	"github.com/scream870102/segb/events"
	"github.com/scream870102/segb/misc"
)

// Variables used for command line parameters
var Token string

func main() {
	configText := os.Getenv("CONFIG")
	var cfg *misc.Config
	if configText != "" {
		json.Unmarshal([]byte(configText), &cfg)
	} else {
		const fileName string = "./config.json"
		tmpCfg, err := misc.ParseConfigFromJSONFile(fileName)
		if err != nil {
			panic(err)
		}
		cfg = tmpCfg
	}

	database.DatabaseInstance().Init()
	database.SpreadSheetInstance().Init()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	registerEvents(dg)
	registerCommands(dg, cfg)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	botKeeper := misc.Keeper{
		Guild:   cfg.Guild,
		Channel: cfg.Channel,
		Offset:  cfg.Delay,
		Session: dg,
	}

	delay := time.Duration(botKeeper.Offset) * time.Minute

	time.AfterFunc(delay, botKeeper.AwakeBOT)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func registerEvents(dg *discordgo.Session) {
	dg.AddHandler(events.NewReadyHandler().Handler)
}

func registerCommands(dg *discordgo.Session, cfg *misc.Config) {
	cmdHandler := commands.NewCommandHandler(cfg.Prefix)

	cmdHandler.RegisterCommand(&commands.CmdShow{})
	cmdHandler.RegisterCommand(&commands.CmdAdd{})
	cmdHandler.RegisterCommand(&commands.CmdUpdate{})
	cmdHandler.RegisterCommand(&commands.CmdList{})
	cmdHandler.RegisterCommand(&commands.CmdHelp{})
	cmdHandler.RegisterMiddleware(&commands.MwPermissions{})
	dg.AddHandler(cmdHandler.HandleMessage)
}
