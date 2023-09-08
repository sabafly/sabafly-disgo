package handlers

import (
	"github.com/sabafly/sabafly-disgo/bot"
	"github.com/sabafly/sabafly-disgo/events"
	"github.com/sabafly/sabafly-disgo/gateway"
)

func gatewayHandlerIntegrationCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventIntegrationCreate) {
	client.EventManager().DispatchEvent(&events.IntegrationCreate{
		GenericIntegration: &events.GenericIntegration{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.GuildID,
			Integration:  event.Integration,
		},
	})
}

func gatewayHandlerIntegrationUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventIntegrationUpdate) {
	client.EventManager().DispatchEvent(&events.IntegrationUpdate{
		GenericIntegration: &events.GenericIntegration{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.GuildID,
			Integration:  event.Integration,
		},
	})
}

func gatewayHandlerIntegrationDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventIntegrationDelete) {
	client.EventManager().DispatchEvent(&events.IntegrationDelete{
		GenericEvent:  events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:       event.GuildID,
		ID:            event.ID,
		ApplicationID: event.ApplicationID,
	})
}
