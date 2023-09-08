package handlers

import (
	"github.com/sabafly/sabafly-disgo/bot"
	"github.com/sabafly/sabafly-disgo/events"
	"github.com/sabafly/sabafly-disgo/gateway"
)

func gatewayHandlerGuildBanAdd(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildBanAdd) {
	client.EventManager().DispatchEvent(&events.GuildBan{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.GuildID,
		User:         event.User,
	})
}

func gatewayHandlerGuildBanRemove(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildBanRemove) {
	client.EventManager().DispatchEvent(&events.GuildUnban{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.GuildID,
		User:         event.User,
	})
}
