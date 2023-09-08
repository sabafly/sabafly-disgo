package events

import "github.com/sabafly/sabafly-disgo/gateway"

type HeartbeatAck struct {
	*GenericEvent
	gateway.EventHeartbeatAck
}
