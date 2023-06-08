package events

import "github.com/sabafly/disgo/gateway"

type Raw struct {
	*GenericEvent
	gateway.EventRaw
}
