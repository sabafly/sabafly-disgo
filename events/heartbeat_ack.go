package events

import "github.com/sabafly/disgo/gateway"

type HeartbeatAck struct {
	*GenericEvent
	gateway.EventHeartbeatAck
}
