package handlers

import (
	"github.com/sabafly/disgo/bot"
	"github.com/sabafly/disgo/discord"
	"github.com/sabafly/disgo/events"
	"github.com/sabafly/disgo/gateway"
	"github.com/sabafly/disgo/httpserver"
	"github.com/sabafly/disgo/rest"
)

func gatewayHandlerInteractionCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventInteractionCreate) {
	handleInteraction(client, sequenceNumber, shardID, nil, event.Interaction)
}

func respond(client bot.Client, respondFunc httpserver.RespondFunc, interaction discord.Interaction) events.InteractionResponderFunc {
	return func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error {
		response := discord.InteractionResponse{
			Type: responseType,
			Data: data,
		}
		if respondFunc != nil {
			return respondFunc(response)
		}
		return client.Rest().CreateInteractionResponse(interaction.ID(), interaction.Token(), response, opts...)
	}
}

func handleInteraction(client bot.Client, sequenceNumber int, shardID int, respondFunc httpserver.RespondFunc, interaction discord.Interaction) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	client.EventManager().DispatchEvent(&events.InteractionCreate{
		GenericEvent: genericEvent,
		Interaction:  interaction,
		Respond:      respond(client, respondFunc, interaction),
	})

	switch i := interaction.(type) {
	case discord.ApplicationCommandInteraction:
		client.EventManager().DispatchEvent(&events.ApplicationCommandInteractionCreate{
			GenericEvent:                  genericEvent,
			ApplicationCommandInteraction: i,
			RespondFunc:                   respond(client, respondFunc, interaction),
		})

	case discord.ComponentInteraction:
		client.EventManager().DispatchEvent(&events.ComponentInteractionCreate{
			GenericEvent:         genericEvent,
			ComponentInteraction: i,
			RespondFunc:          respond(client, respondFunc, interaction),
		})

	case discord.AutocompleteInteraction:
		client.EventManager().DispatchEvent(&events.AutocompleteInteractionCreate{
			GenericEvent:            genericEvent,
			AutocompleteInteraction: i,
			Respond:                 respond(client, respondFunc, interaction),
		})

	case discord.ModalSubmitInteraction:
		client.EventManager().DispatchEvent(&events.ModalSubmitInteractionCreate{
			GenericEvent:           genericEvent,
			ModalSubmitInteraction: i,
			RespondFunc:            respond(client, respondFunc, interaction),
		})

	default:
		client.Logger().Errorf("unknown interaction with type %T received", interaction)
	}
}
