package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const configPath = "config.json"

type Duration time.Duration

type APIConfig struct {
	Host    string   `json:"host"`
	Timeout Duration `json:"timeout"`
}

type CommandConfig struct {
	Timeout int `json:"timeout"`
	Offset  int `json:"offset"`
}

type BotConfig struct {
	Token string `json:"token"`
}

type Config struct {
	API     *APIConfig     `json:"api"`
	Command *CommandConfig `json:"command"`
	Bot     *BotConfig     `json:"bot"`
}

func NewConfig() (*Config, error) {
	cfg := new(Config)

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "read config")
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal config")
	}

	if cfg.API == nil {
		return nil, errors.New("api is nil")
	}

	if strings.TrimSpace(cfg.API.Host) == "" {
		return nil, errors.New("empty api host")
	}

	if cfg.API.Timeout <= 0 {
		return nil, errors.New("incorrect api timeout")
	}

	if cfg.Bot == nil {
		return nil, errors.New("bot is nil")
	}

	if cfg.Bot.Token == "" {
		return nil, errors.New("empty bot token")
	}

	if cfg.Command == nil {
		return nil, errors.New("command is nil")
	}

	if cfg.Command.Timeout <= 0 {
		return nil, errors.New("incorrect command timeout")
	}

	return cfg, nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return d.MarshalJSON()
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	value, ok := v.(string)
	if !ok {
		return errors.New("invalid duration")
	}

	val, err := time.ParseDuration(value)
	if err != nil {
		return err
	}

	*d = Duration(val)

	return nil
}
