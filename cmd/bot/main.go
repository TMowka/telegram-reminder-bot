package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"

	"github.com/tmowka/telegram-reminder-bot/cmd/bot/internal/handlers"
	"github.com/tmowka/telegram-reminder-bot/internal/platform/bot"
	"github.com/tmowka/telegram-reminder-bot/internal/platform/database"
)

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}

func run() error {
	// =========================================================================
	// Logging
	log := log.New(os.Stdout, "BOT : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// =========================================================================
	// Configuration
	var cfg struct {
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:password,noprint"`
			Host       string `conf:"default:0.0.0.0"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
		BOT struct {
			Token    string `conf:""`
			Location string `conf:"default:Europe/Minsk"`
		}
		CHAT struct {
			Id string `conf:""`
		}
	}

	if err := conf.Parse(os.Args[1:], "BOT", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("BOT", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	// =========================================================================
	// Start Database

	log.Println("main : Started : Initializing database support")

	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer func() {
		log.Printf("main : Database Stopping : %s", cfg.DB.Host)
		db.Close()
	}()

	// =========================================================================
	// Start Bot

	log.Println("main : Started : Initializing bot")

	b, err := bot.Create(bot.Config{
		Token:    cfg.BOT.Token,
		Location: cfg.BOT.Location,
	})
	if err != nil {
		return errors.Wrap(err, "creating telebot")
	}

	err = handlers.Telebot(db, b, cfg.CHAT.Id)
	if err != nil {
		return errors.Wrap(err, "registration of telebot handlers")
	}

	log.Println("main : Started : Starting telebot")
	b.Start()

	return nil
}
