package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"ludobot/internal/arguments"
	"ludobot/internal/di"
	"ludobot/internal/handlers"
	"os"
	"os/signal"

	"github.com/akamensky/argparse"
	"github.com/go-telegram/bot"
	"gopkg.in/yaml.v3"
)

func getToken(dependencies *di.Dependencies) string {
	fileData, err := os.ReadFile(dependencies.Args.PathToBotConfig)
	if err != nil {
		slog.Error(fmt.Sprintf("error opening file: %s", err.Error()))
		os.Exit(1)
	}

	yamlData := make(map[string]interface{})
	yaml.Unmarshal(fileData, &yamlData)

	return yamlData["token"].(string)
}

func parseArgs() arguments.Arguments {
	parser := argparse.NewParser("ludobot", "Ludobot, an entertainment gambling bot for telegram")

	pathToLocalizationConfig := parser.String("", "localization-config", &argparse.Options{
		Required: true,
		Help:     "Path to localization config",
	})
	pathToBotConfig := parser.String("", "bot-config", &argparse.Options{
		Required: true,
		Help:     "Path to bot config",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return arguments.Arguments{
		PathToLocalizationConfig: *pathToLocalizationConfig,
		PathToBotConfig:          *pathToBotConfig,
	}
}

func createDependencies() *di.Dependencies {
	return &di.Dependencies{
		Args:    parseArgs(),
		Context: context.Background(),
	}
}

func registerHandlers(b *bot.Bot, dependencies *di.Dependencies) {
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		handlers.NewDefaultHandler(dependencies),
	)
}

func run() error {
	dependencies := createDependencies()
	token := getToken(dependencies)

	if len(token) == 0 {
		return errors.New("bot token is empty")
	}

	b, err := bot.New(token)

	if err != nil {
		return fmt.Errorf("bot creation error: %s", err.Error())
	}

	registerHandlers(b, dependencies)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b.Start(ctx)

	return nil
}

func main() {
	err := run()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
