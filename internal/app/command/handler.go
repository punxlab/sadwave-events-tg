package command

import (
	"bytes"
	"context"
	"github.com/goodsign/monday"
	"html/template"
	"strings"
	"sync"
	"time"

	"github.com/punxlab/sadwave-events-tg/internal/app/api/model"
)

const (
	commandStart = "/start"
	commandHelp  = "/help"
)

type Handler interface {
	Handle(ctx context.Context, cmd string) (string, error)
}

type handler struct {
	api model.SadwaveAPI
	cmd map[string]*model.City
	mu  sync.RWMutex
}

func NewCommandHandler(api model.SadwaveAPI) Handler {
	return &handler{
		api: api,
	}
}

func (h *handler) Handle(ctx context.Context, cmd string) (string, error) {
	err := h.fillCitiesCommands(ctx)
	if err != nil {
		return "", err
	}

	if cmd == commandStart {
		return h.startResponse()
	}

	if cmd == commandHelp {
		return h.helpResponse()
	}

	if h.isCityCommand(cmd) {
		events, err := h.api.Events(ctx, cmd)
		if err != nil {
			return "", err
		}
		return eventsResponse(events)
	}

	return h.helpResponse()
}

func (h *handler) fillCitiesCommands(ctx context.Context) error {
	if len(h.cmd) > 0 {
		return nil
	}

	cities, err := h.api.Cities(ctx)
	if err != nil {
		return err
	}

	h.setCommands(citiesToCommands(cities))

	return nil
}

func citiesToCommands(cities []*model.City) map[string]*model.City {
	res := make(map[string]*model.City, len(cities))
	for _, c := range cities {
		res["/"+c.Alias] = c
	}

	return res
}

func (h *handler) isCityCommand(cmd string) bool {
	_, ok := h.commands()[cmd]
	return ok
}

func eventsResponse(events []*model.Event) (string, error) {
	if len(events) == 0 {
		return "Гигов нет", nil
	}

	t, err := template.
		New("events").
		Funcs(template.FuncMap{
			"formatDate": formatDate,
			"formatTime": formatTime,
		}).
		Parse(`{{range .}}
<a href="{{.Url}}"><b>{{formatDate .Date.Date}} в {{formatTime .Date.Time}}</b></a>
{{.Overview}}
{{end}}`)
	if err != nil {
		return "", err
	}

	return renderTemplate(t, events)
}

func (h *handler) commands() map[string]*model.City {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.cmd
}

func (h *handler) setCommands(commands map[string]*model.City) {
	h.mu.Lock()
	h.cmd = commands
	h.mu.Unlock()
}

func (h *handler) helpResponse() (string, error) {
	t, err := template.
		New("help").
		Parse(`{{range .}}/{{.Alias}} - {{.Name}}
{{end}}`)
	if err != nil {
		return "", err
	}

	return renderTemplate(t, h.commands())
}

func (h *handler) startResponse() (string, error) {
	t, err := template.
		New("start").
		Parse(`Привет. Здесь ты можешь найти афишу гигов под редакцией <a href="https://sadwave.com/">sadwave</a>.
Вот список команд, которые ты можешь использовать:
{{.}}/help - Напомнит команды выше`)
	if err != nil {
		return "", err
	}

	helpCommand, err := h.helpResponse()
	if err != nil {
		return "", err
	}

	return renderTemplate(t, helpCommand)
}

func formatDate(d string) string {
	res, err := time.Parse("2006-01-02T15:00:00", d)
	if err != nil {
		return d
	}

	return strings.ToLower(monday.Format(res, "02 January, Monday", monday.LocaleRuRU))
}

func formatTime(t string) string {
	res, err := time.Parse("15:00:00", t)
	if err != nil {
		return t
	}

	return res.Format("15:00")
}

func renderTemplate(t *template.Template, data interface{}) (string, error) {
	buf := bytes.NewBufferString("")
	err := t.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
