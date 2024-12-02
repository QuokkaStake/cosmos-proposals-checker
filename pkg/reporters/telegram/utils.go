package telegram

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) ReplyRender(
	c tele.Context,
	templateName string,
	renderStruct any,
) error {
	template, err := reporter.TemplatesManager.Render(templateName, renderStruct)
	if err != nil {
		reporter.Logger.Error().Str("template", templateName).Err(err).Msg("Error rendering template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	return reporter.BotReply(c, template)
}

func (reporter *Reporter) EditRender(
	c tele.Context,
	message *tele.Message,
	templateName string,
	renderStruct interface{},
	opts ...interface{},
) error {
	opts = append(opts, tele.ModeHTML, tele.NoPreview)

	template, renderErr := reporter.TemplatesManager.Render(templateName, renderStruct)

	if renderErr != nil {
		reporter.Logger.Error().Str("template", templateName).Err(renderErr).Msg("Error rendering template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", renderErr))
	}

	if _, editErr := reporter.TelegramBot.Edit(message, strings.TrimSpace(template), opts...); editErr != nil {
		reporter.Logger.Error().Err(editErr).Msg("Error editing message")
		return editErr
	}

	return nil
}
