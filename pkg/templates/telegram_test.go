package templates

import (
	templatePkg "html/template"
	loggerPkg "main/pkg/logger"
	"main/pkg/types"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestTelegramGetTemplateNotExisting(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewTelegramTemplatesManager(logger, timezone)

	template, err := manager.GetTemplate("not-existing")
	require.Error(t, err)
	assert.Nil(t, template)
}

func TestTelegramGetTemplateExistingAndCached(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewTelegramTemplatesManager(logger, timezone)

	template, err := manager.GetTemplate("help")
	require.NoError(t, err)
	assert.NotNil(t, template)

	// this time it should be loaded from cache
	template2, err2 := manager.GetTemplate("help")
	require.NoError(t, err2)
	assert.NotNil(t, template2)
}

func TestTelegramRenderTemplateError(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewTelegramTemplatesManager(logger, timezone)

	template, err := manager.Render("not-existing", nil)
	require.Error(t, err)
	assert.Empty(t, template)
}

func TestTelegramRenderTemplateRenderError(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewTelegramTemplatesManager(logger, timezone)

	value := map[string]interface{}{}

	template, err := manager.Render("voted", value)
	require.Error(t, err)
	assert.Empty(t, template)
}

func TestTelegramRenderTemplateRenderSuccess(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewTelegramTemplatesManager(logger, timezone)

	template, err := manager.Render("help", "1.0.0")
	require.NoError(t, err)
	assert.NotEmpty(t, template)
}

func TestTelegramRenderSerializeLink(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewTelegramTemplatesManager(logger, timezone)

	assert.Equal(t, templatePkg.HTML("test"), manager.SerializeLink(types.Link{Name: "test"}))
	assert.Equal(
		t,
		templatePkg.HTML("<a href='href'>test</a>"),
		manager.SerializeLink(types.Link{Name: "test", Href: "href"}),
	)
}

func TestTelegramRenderSerializeDate(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	timezone, _ := time.LoadLocation("Europe/Moscow")
	manager := NewTelegramTemplatesManager(logger, timezone)

	dateStr := manager.SerializeDate(
		time.Date(2000, 1, 1, 0, 0, 0, 0, timezone),
	)
	assert.Equal(t, "Sat, 01 Jan 2000 00:00:00 MSK", dateStr)
}
