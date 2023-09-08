package handlers

import (
	"github.com/sabafly/sabafly-disgo/bot"
	"github.com/sabafly/sabafly-disgo/events"
	"github.com/sabafly/sabafly-disgo/gateway"
)

func gatewayHandlerGuildIntegrationsUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildIntegrationsUpdate) {
	client.EventManager().DispatchEvent(&events.GuildIntegrationsUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.GuildID,
	})
}
