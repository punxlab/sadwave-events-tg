package app

import (
	"context"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/punxlab/sadwave-events-tg/internal/app/api"
	"github.com/punxlab/sadwave-events-tg/internal/app/command"
	"github.com/punxlab/sadwave-events-tg/internal/config"
	"io/ioutil"
	"log"
	"net/http"
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

		messages, err := r.handler.Handle(ctx, u.Message.Text)
		if err != nil {
			log.Print(err)
			continue
		}

		for _, m := range messages {
			msg, err := tgMessage(m, u.Message.Chat.ID)
			if err != nil {
				log.Print(err)
				continue
			}

			_, err = r.bot.Send(msg)
			if err != nil {
				log.Print(err)
				continue
			}
		}
	}

	return nil
}

func tgMessage(msg *command.Message, chat int64) (tg.Chattable, error) {
	if msg.Photo != "" {
		r, err := http.Get(msg.Photo)
		if err != nil {
			return nil, err
		}

		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		return tg.PhotoConfig{
			BaseFile: tg.BaseFile{
				BaseChat: tg.BaseChat{
					ChatID:              chat,
					DisableNotification: true,
				},
				File: tg.FileBytes{
					Name:  msg.Photo,
					Bytes: body,
				},
			},
			Caption:   msg.Markup,
			ParseMode: "HTML",
		}, nil
	}

	return tg.MessageConfig{
		BaseChat: tg.BaseChat{
			ChatID:              chat,
			DisableNotification: true,
		},
		Text:                  msg.Markup,
		ParseMode:             tg.ModeHTML,
		DisableWebPagePreview: true,
	}, nil
}
