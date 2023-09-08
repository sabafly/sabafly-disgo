package handlers

import (
	"github.com/sabafly/sabafly-disgo/bot"
	"github.com/sabafly/sabafly-disgo/events"
	"github.com/sabafly/sabafly-disgo/gateway"
)

func gatewayHandlerUserUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventUserUpdate) {
	oldUser, _ := client.Caches().SelfUser()
	client.Caches().SetSelfUser(event.OAuth2User)

	client.EventManager().DispatchEvent(&events.SelfUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		SelfUser:     event.OAuth2User,
		OldSelfUser:  oldUser,
	})

}
