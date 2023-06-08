package middleware

import (
	"github.com/sabafly/disgo/events"
	"github.com/sabafly/disgo/handler"
)

var Logger handler.Middleware = func(next handler.Handler) handler.Handler {
	return func(e *events.InteractionCreate) error {
		e.Client().Logger().Infof("handling interaction: %s\n", e.Interaction.ID())
		return next(e)
	}
}
