package main

import (
	logger "log"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"

	"github.com/tmowka/telegram-reminder-bot/cmd/bot/internal/handlers"
	"github.com/tmowka/telegram-reminder-bot/internal/platform/bot"
	"github.com/tmowka/telegram-reminder-bot/internal/platform/database"
	"github.com/tmowka/telegram-reminder-bot/internal/schema"
)

type config struct {
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

func main() {
	cfg, err := configure()
	if err != nil {
		logger.Println("error :", err)
		os.Exit(1)
	}

	if err := migrate(cfg); err != nil {
		logger.Println("error :", err)
		os.Exit(1)
	}

	if err := run(cfg); err != nil {
		logger.Println("error :", err)
		os.Exit(1)
	}
}

func run(cfg *config) error {
	// =========================================================================
	// Logging
	log := logger.New(os.Stdout, "BOT : ",
		logger.LstdFlags|logger.Lmicroseconds|logger.Lshortfile)

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
	// Migrate Database

	log.Println("main : Started : Migrating database")

	if err := schema.Migrate(db); err != nil {
		return errors.Wrap(err, "migrating db")
	}

	// =========================================================================
	// Start Bot

	log.Println("main : Started : Initializing bot")

	b, err := bot.Create(bot.Config{
		Token: cfg.BOT.Token,
	})
	if err != nil {
		return errors.Wrap(err, "creating telebot")
	}

	err = handlers.Telebot(db, b, cfg.CHAT.Id, cfg.BOT.Location)
	if err != nil {
		return errors.Wrap(err, "registration of telebot handlers")
	}

	log.Println("main : Started : Starting telebot")
	b.Start()

	return nil
}

func migrate(cfg *config) error {
	// =========================================================================
	// Logging
	log := logger.New(os.Stdout, "BOT : ",
		logger.LstdFlags|logger.Lmicroseconds|logger.Lshortfile)

	// =========================================================================
	// Migrate Database

	log.Println("migrate : Started : Migrating database")

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
	defer db.Close()

	if err := schema.Migrate(db); err != nil {
		return errors.Wrap(err, "migrating database")
	}

	log.Println("migrate : Completed : Migrating database")
	return nil
}

func configure() (*config, error) {
	// =========================================================================
	// Logging
	log := logger.New(os.Stdout, "BOT : ",
		logger.LstdFlags|logger.Lmicroseconds|logger.Lshortfile)

	log.Println("configure : Started : Initializing config")

	var cfg config

	if err := conf.Parse(os.Args[1:], "BOT", &cfg); err != nil {
		if err != conf.ErrHelpWanted {
			return nil, errors.Wrap(err, "parsing config")
		}

		_, err := conf.Usage("BOT", &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "generating config usage")
		}
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "generating config for output")
	}
	log.Printf("configure : Config :\n%v\n", out)

	return &cfg, nil
}
