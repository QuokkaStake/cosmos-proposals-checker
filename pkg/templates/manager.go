package templates

type Manager interface {
	Render(templateName string, data interface{}) (string, error)
}
