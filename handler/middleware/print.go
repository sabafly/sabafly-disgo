package middleware

import (
	"github.com/sabafly/sabafly-disgo/events"
	"github.com/sabafly/sabafly-disgo/handler"
)

func Print(content string) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(event *events.InteractionCreate) error {
			println(content)
			return next(event)
		}
	}
}
