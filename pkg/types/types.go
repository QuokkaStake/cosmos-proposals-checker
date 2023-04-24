package types

import (
	"fmt"
)

type Link struct {
	Name string
	Href string
}

func (l Link) Serialize() string {
	if l.Href == "" {
		return l.Name
	}

	return fmt.Sprintf("<a href='%s'>%s</a>", l.Href, l.Name)
}
