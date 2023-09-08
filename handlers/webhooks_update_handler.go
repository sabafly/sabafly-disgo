package handlers

import (
	"github.com/sabafly/sabafly-disgo/bot"
	"github.com/sabafly/sabafly-disgo/events"
	"github.com/sabafly/sabafly-disgo/gateway"
)

func gatewayHandlerWebhooksUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventWebhooksUpdate) {
	client.EventManager().DispatchEvent(&events.WebhooksUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildId:      event.GuildID,
		ChannelID:    event.ChannelID,
	})
}
