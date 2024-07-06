package templates

import (
	templatePkg "html/template"
	loggerPkg "main/pkg/logger"
	mutes "main/pkg/mutes"
	"main/pkg/types"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestDiscordGetTemplateNotExisting(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewDiscordTemplatesManager(logger, timezone)

	template, err := manager.GetTemplate("not-existing")
	require.Error(t, err)
	assert.Nil(t, template)
}

func TestDiscordGetTemplateExistingAndCached(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewDiscordTemplatesManager(logger, timezone)

	template, err := manager.GetTemplate("help")
	require.NoError(t, err)
	assert.NotNil(t, template)

	// this time it should be loaded from cache
	template2, err2 := manager.GetTemplate("help")
	require.NoError(t, err2)
	assert.NotNil(t, template2)
}

func TestDiscordRenderTemplateError(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewDiscordTemplatesManager(logger, timezone)

	template, err := manager.Render("not-existing", nil)
	require.Error(t, err)
	assert.Empty(t, template)
}

func TestDiscordRenderTemplateRenderError(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewDiscordTemplatesManager(logger, timezone)

	value := map[string]interface{}{}

	template, err := manager.Render("voted", value)
	require.Error(t, err)
	assert.Empty(t, template)
}

func TestDiscordRenderTemplateRenderSuccess(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewDiscordTemplatesManager(logger, timezone)

	template, err := manager.Render("mutes", []mutes.Mute{})
	require.NoError(t, err)
	assert.NotEmpty(t, template)
}

func TestDiscordRenderSerializeLink(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewDiscordTemplatesManager(logger, timezone)

	assert.Equal(t, templatePkg.HTML("test"), manager.SerializeLink(types.Link{Name: "test"}))
	assert.Equal(
		t,
		templatePkg.HTML("[test](<href>)"),
		manager.SerializeLink(types.Link{Name: "test", Href: "href"}),
	)
}

func TestDiscordRenderSerializeDate(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewDiscordTemplatesManager(logger, timezone)

	dateStr := manager.SerializeDate(
		time.Date(2000, 1, 1, 0, 0, 0, 0, timezone),
	)
	assert.Equal(t, "Sat, 01 Jan 2000 00:00:00 MSK", dateStr)
}
