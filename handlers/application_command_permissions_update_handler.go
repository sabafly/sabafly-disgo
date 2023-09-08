package handlers

import (
	"github.com/sabafly/sabafly-disgo/bot"
	"github.com/sabafly/sabafly-disgo/events"
	"github.com/sabafly/sabafly-disgo/gateway"
)

func gatewayHandlerApplicationCommandPermissionsUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventApplicationCommandPermissionsUpdate) {
	client.EventManager().DispatchEvent(&events.GuildApplicationCommandPermissionsUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		Permissions:  event.ApplicationCommandPermissions,
	})
}
