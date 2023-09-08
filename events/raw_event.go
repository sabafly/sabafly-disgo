package events

import "github.com/sabafly/sabafly-disgo/gateway"

type Raw struct {
	*GenericEvent
	gateway.EventRaw
}
