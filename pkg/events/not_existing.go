package events

type NotExistingEvent struct{}

func (e NotExistingEvent) Name() string {
	return "not_existing"
}

func (e NotExistingEvent) IsAlert() bool {
	return false
}
