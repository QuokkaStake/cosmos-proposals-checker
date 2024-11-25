package telegram

import (
	"fmt"
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
