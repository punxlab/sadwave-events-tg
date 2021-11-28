package app

import (
	"context"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
	"github.com/punxlab/sadwave-events-tg/internal/app/api"
	"github.com/punxlab/sadwave-events-tg/internal/app/command"
	"github.com/punxlab/sadwave-events-tg/internal/config"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const defaultCacheExpiration = 24 * time.Hour
const cleanupCacheInterval = 24 * time.Hour

type Runner interface {
	Run(ctx context.Context) error
}

type app struct {
	bot     *tg.BotAPI
	cfg     *config.Config
	handler command.Handler

	fileCache *cache.Cache
}

func NewApp(cfg *config.Config) (Runner, error) {
	bot, err := tg.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return nil, err
	}

	fileCache := cache.New(defaultCacheExpiration, cleanupCacheInterval)
	return &app{
		bot:       bot,
		cfg:       cfg,
		fileCache: fileCache,
		handler: command.NewCommandHandler(
			api.NewSadwaveAPI(cfg.API.Host),
		),
	}, nil
}

func (a *app) Run(ctx context.Context) error {
	cfg := tg.UpdateConfig{
		Offset:  a.cfg.Command.Offset,
		Timeout: a.cfg.Command.Timeout,
	}

	updates, err := a.bot.GetUpdatesChan(cfg)
	if err != nil {
		return err
	}

	for u := range updates {
		if u.Message == nil {
			continue
		}

		messages, err := a.handler.Handle(ctx, u.Message.Text)
		if err != nil {
			log.Print(err)
			continue
		}

		for _, m := range messages {
			msg, err := a.tgMessage(m, u.Message.Chat.ID)
			if err != nil {
				log.Print(err)
				continue
			}

			_, err = a.bot.Send(msg)
			if err != nil {
				log.Print(err)
				continue
			}
		}
	}

	return nil
}

func (a *app) tgMessage(msg *command.Message, chat int64) (tg.Chattable, error) {
	if msg.Photo != "" {
		f, err := a.getFile(msg.Photo)
		if err != nil {
			return nil, err
		}

		return tg.PhotoConfig{
			BaseFile: tg.BaseFile{
				BaseChat: tg.BaseChat{
					ChatID:              chat,
					DisableNotification: true,
				},
				File: f,
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

func (a *app) getFile(url string) (tg.FileBytes, error) {
	file, ok := a.fileCache.Get(url)
	if ok {
		return file.(tg.FileBytes), nil
	}

	r, err := http.Get(url)
	if err != nil {
		return tg.FileBytes{}, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return tg.FileBytes{}, err
	}

	f := tg.FileBytes{
		Name:  url,
		Bytes: body,
	}

	a.fileCache.Set(url, f, defaultCacheExpiration)
	return f, nil
}
