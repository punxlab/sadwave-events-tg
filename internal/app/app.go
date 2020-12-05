package app

import (
	"context"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"

	"github.com/punxlab/sadwave-events-tg/internal/app/api"
	"github.com/punxlab/sadwave-events-tg/internal/app/command"
	"github.com/punxlab/sadwave-events-tg/internal/config"
)

type Runner interface {
	Run(ctx context.Context) error
}

type app struct {
	bot     *tg.BotAPI
	cfg     *config.Config
	handler command.Handler
}

func NewApp(cfg *config.Config) (Runner, error) {
	bot, err := tg.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return nil, err
	}

	return &app{
		bot: bot,
		cfg: cfg,
		handler: command.NewCommandHandler(
			api.NewSadwaveAPI(cfg.API.Host),
		),
	}, nil
}

func (r *app) Run(ctx context.Context) error {
	cfg := tg.UpdateConfig{
		Offset:  r.cfg.Command.Offset,
		Timeout: r.cfg.Command.Timeout,
	}

	updates, err := r.bot.GetUpdatesChan(cfg)
	if err != nil {
		return err
	}

	for u := range updates {
		if u.Message == nil {
			continue
		}

		res, err := r.handler.Handle(ctx, u.Message.Text)
		if err != nil {
			log.Print(err)
			continue
		}

		msg := tg.MessageConfig{
			BaseChat: tg.BaseChat{
				ChatID: u.Message.Chat.ID,
			},
			Text:                  res,
			ParseMode:             tg.ModeHTML,
			DisableWebPagePreview: true,
		}

		_, err = r.bot.Send(msg)
		if err != nil {
			log.Print(err)
			continue
		}
	}

	return nil
}
