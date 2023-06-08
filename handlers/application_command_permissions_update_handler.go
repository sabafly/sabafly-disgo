package handlers

import (
	"github.com/sabafly/disgo/bot"
	"github.com/sabafly/disgo/events"
	"github.com/sabafly/disgo/gateway"
)

func gatewayHandlerApplicationCommandPermissionsUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventApplicationCommandPermissionsUpdate) {
	client.EventManager().DispatchEvent(&events.GuildApplicationCommandPermissionsUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		Permissions:  event.ApplicationCommandPermissions,
	})
}
