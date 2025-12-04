package handlers

import (
	"context"
	"ludobot/internal/di"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func NewDefaultHandler(dependencies *di.Dependencies) func(context.Context, *bot.Bot, *models.Update) {
	return func(ctx context.Context, b *bot.Bot, m *models.Update) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: m.Message.Chat.ID,
			Text:   m.Message.Text,
		})
	}
}
