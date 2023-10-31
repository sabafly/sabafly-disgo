package handlers

import (
	"github.com/sabafly/sabafly-disgo/bot"
	"github.com/sabafly/sabafly-disgo/events"
	"github.com/sabafly/sabafly-disgo/gateway"
)

func gatewayHandlerEntitlementCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventEntitlementCreate) {
	client.EventManager().DispatchEvent(&events.EntitlementCreate{
		GenericEntitlementEvent: &events.GenericEntitlementEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			Entitlement:  event.Entitlement,
		},
	})
}

func gatewayHandlerEntitlementUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventEntitlementUpdate) {
	client.EventManager().DispatchEvent(&events.EntitlementUpdate{
		GenericEntitlementEvent: &events.GenericEntitlementEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			Entitlement:  event.Entitlement,
		},
	})
}

func gatewayHandlerEntitlementDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventEntitlementDelete) {
	client.EventManager().DispatchEvent(&events.EntitlementDelete{
		GenericEntitlementEvent: &events.GenericEntitlementEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			Entitlement:  event.Entitlement,
		},
	})
}
