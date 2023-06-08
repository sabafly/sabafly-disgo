package handlers

import (
	"github.com/sabafly/disgo/bot"
	"github.com/sabafly/disgo/events"
	"github.com/sabafly/disgo/gateway"
)

func gatewayHandlerVoiceStateUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventVoiceStateUpdate) {
	member := event.Member

	oldVoiceState, oldOk := client.Caches().VoiceState(event.GuildID, event.UserID)
	if event.ChannelID == nil {
		client.Caches().RemoveVoiceState(event.GuildID, event.UserID)
	} else {
		client.Caches().AddVoiceState(event.VoiceState)
	}
	client.Caches().AddMember(member)

	if event.UserID == client.ID() && client.VoiceManager() != nil {
		client.VoiceManager().HandleVoiceStateUpdate(event)
	}

	genericGuildVoiceEvent := &events.GenericGuildVoiceState{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		VoiceState:   event.VoiceState,
		Member:       member,
	}

	client.EventManager().DispatchEvent(&events.GuildVoiceStateUpdate{
		GenericGuildVoiceState: genericGuildVoiceEvent,
		OldVoiceState:          oldVoiceState,
	})

	if oldOk && oldVoiceState.ChannelID != nil && event.ChannelID != nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceMove{
			GenericGuildVoiceState: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else if (oldOk || oldVoiceState.ChannelID == nil) && event.ChannelID != nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceJoin{
			GenericGuildVoiceState: genericGuildVoiceEvent,
		})
	} else if event.ChannelID == nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceLeave{
			GenericGuildVoiceState: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else {
		client.Logger().Warnf("could not decide which GuildVoice to fire")
	}
}

func gatewayHandlerVoiceServerUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventVoiceServerUpdate) {
	if client.VoiceManager() != nil {
		client.VoiceManager().HandleVoiceServerUpdate(event)
	}

	client.EventManager().DispatchEvent(&events.VoiceServerUpdate{
		GenericEvent:           events.NewGenericEvent(client, sequenceNumber, shardID),
		EventVoiceServerUpdate: event,
	})
}
