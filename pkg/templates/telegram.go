package templates

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"main/pkg/types"
	"main/pkg/utils"
	"main/templates"
	"time"

	"github.com/rs/zerolog"
)

type TelegramTemplatesManager struct {
	Templates map[string]*template.Template
	Logger    zerolog.Logger
	Timezone  *time.Location
}

func NewTelegramTemplatesManager(
	logger *zerolog.Logger,
	timezone *time.Location,
) *TelegramTemplatesManager {
	return &TelegramTemplatesManager{
		Templates: map[string]*template.Template{},
		Logger:    logger.With().Str("component", "telegram_templates_manager").Logger(),
		Timezone:  timezone,
	}
}

func (m *TelegramTemplatesManager) Render(templateName string, data interface{}) (string, error) {
	templateToRender, err := m.GetTemplate(templateName)
	if err != nil {
		m.Logger.Error().
			Err(err).
			Str("name", templateName).
			Msg("Error getting template")
		return "", err
	}

	var buffer bytes.Buffer
	if err := templateToRender.Execute(&buffer, data); err != nil {
		m.Logger.Error().
			Err(err).
			Str("name", templateName).
			Msg("Error rendering template")
		return "", err
	}

	return buffer.String(), nil
}

func (m *TelegramTemplatesManager) GetTemplate(templateName string) (*template.Template, error) {
	if cachedTemplate, ok := m.Templates[templateName]; ok {
		m.Logger.Trace().Str("type", templateName).Msg("Using cached template")
		return cachedTemplate, nil
	}

	m.Logger.Trace().Str("type", templateName).Msg("Loading template")

	filename := templateName + ".html"

	t, err := template.New(filename).Funcs(template.FuncMap{
		"SerializeLink":  m.SerializeLink,
		"SerializeDate":  m.SerializeDate,
		"FormatDuration": utils.FormatDuration,
	}).ParseFS(templates.TemplatesFs, "telegram/"+filename)
	if err != nil {
		return nil, err
	}

	m.Templates[templateName] = t

	return t, nil
}

func (m *TelegramTemplatesManager) SerializeLink(link types.Link) template.HTML {
	if link.Href != "" {
		return template.HTML(fmt.Sprintf("<a href='%s'>%s</a>", link.Href, html.EscapeString(link.Name)))
	}

	return template.HTML(link.Name)
}

func (m *TelegramTemplatesManager) SerializeDate(date time.Time) string {
	return date.In(m.Timezone).Format(time.RFC1123)
}
