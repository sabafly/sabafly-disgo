package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

// InteractionResponderFunc is a function that can be used to respond to a discord.Interaction.
type InteractionResponderFunc func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error

// InteractionCreate indicates that a new interaction has been created.
type InteractionCreate struct {
	*GenericEvent
	discord.Interaction
	Respond InteractionResponderFunc
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *InteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches.Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

// ApplicationCommandInteractionCreate is the base struct for all application command interaction create events.
type ApplicationCommandInteractionCreate struct {
	*GenericEvent
	discord.ApplicationCommandInteraction
	RespondFunc  InteractionResponderFunc
	acknowledged bool
}

func (e *ApplicationCommandInteractionCreate) Respond(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error {
	if e.acknowledged {
		return nil
	}
	if err := e.RespondFunc(responseType, data, opts...); err != nil {
		return err
	}
	e.acknowledged = true
	return nil
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *ApplicationCommandInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches.Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

// Acknowledge acknowledges the interaction.
//
// This is used strictly for acknowledging the HTTP interaction request from discord. This responds with 202 Accepted.
//
// When using this, your first http request must be [rest.Interactions.CreateInteractionResponse] or [rest.Interactions.CreateInteractionResponseWithCallback]
//
// This does not produce a visible loading state to the user.
// You are expected to send a new http request within 3 seconds to respond to the interaction.
// This allows you to gracefully handle errors with your sent response & access the resulting message.
//
// If you want to create a visible loading state, use DeferCreateMessage.
//
// Source docs: [Discord Source docs]
//
// [Discord Source docs]: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-callback
func (e *ApplicationCommandInteractionCreate) Acknowledge(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeAcknowledge, nil, opts...)
}

// CreateMessage responds to the interaction with a new message.
func (e *ApplicationCommandInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e *ApplicationCommandInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// Modal responds to the interaction with a new modal.
func (e *ApplicationCommandInteractionCreate) Modal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeModal, modalCreate, opts...)
}

func (e *ApplicationCommandInteractionCreate) RespondMessage(messageBuilder discord.MessageBuilder, opts ...rest.RequestOpt) error {
	if e.acknowledged {
		_, err := e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), messageBuilder.BuildUpdate(), opts...)
		return err
	} else {
		return e.CreateMessage(messageBuilder.BuildCreate())
	}
}

// LaunchActivity responds to the interaction by launching activity associated with the app.
func (e *ApplicationCommandInteractionCreate) LaunchActivity(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeLaunchActivity, nil, opts...)
}

// ComponentInteractionCreate indicates that a new component interaction has been created.
type ComponentInteractionCreate struct {
	*GenericEvent
	discord.ComponentInteraction
	RespondFunc  InteractionResponderFunc
	acknowledged bool
}

func (e *ComponentInteractionCreate) Respond(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error {
	if e.acknowledged {
		return nil
	}
	if err := e.RespondFunc(responseType, data, opts...); err != nil {
		return err
	}
	e.acknowledged = true
	return nil
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *ComponentInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches.Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

// Acknowledge acknowledges the interaction.
//
// This is used strictly for acknowledging the HTTP interaction request from discord. This responds with 202 Accepted.
//
// When using this, your first http request must be [rest.Interactions.CreateInteractionResponse] or [rest.Interactions.CreateInteractionResponseWithCallback]
//
// This does not produce a visible loading state to the user.
// You are expected to send a new http request within 3 seconds to respond to the interaction.
// This allows you to gracefully handle errors with your sent response & access the resulting message.
//
// If you want to create a visible loading state, use DeferCreateMessage.
//
// Source docs: [Discord Source docs]
//
// [Discord Source docs]: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-callback
func (e *ComponentInteractionCreate) Acknowledge(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeAcknowledge, nil, opts...)
}

// CreateMessage responds to the interaction with a new message.
func (e *ComponentInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e *ComponentInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// UpdateMessage responds to the interaction with updating the message the component is from.
func (e *ComponentInteractionCreate) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeUpdateMessage, messageUpdate, opts...)
}

// DeferUpdateMessage responds to the interaction with nothing.
func (e *ComponentInteractionCreate) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil, opts...)
}

// Modal responds to the interaction with a new modal.
func (e *ComponentInteractionCreate) Modal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeModal, modalCreate, opts...)
}

func (e *ComponentInteractionCreate) RespondMessage(messageBuilder discord.MessageBuilder, opts ...rest.RequestOpt) error {
	if e.acknowledged {
		_, err := e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), messageBuilder.BuildUpdate(), opts...)
		return err
	} else {
		return e.CreateMessage(messageBuilder.BuildCreate())
	}
}

// LaunchActivity responds to the interaction by launching activity associated with the app.
func (e *ComponentInteractionCreate) LaunchActivity(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeLaunchActivity, nil, opts...)
}

// AutocompleteInteractionCreate indicates that a new autocomplete interaction has been created.
type AutocompleteInteractionCreate struct {
	*GenericEvent
	discord.AutocompleteInteraction
	Respond InteractionResponderFunc
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *AutocompleteInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches.Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

// Acknowledge acknowledges the interaction.
//
// This is used strictly for acknowledging the HTTP interaction request from discord. This responds with 202 Accepted.
//
// When using this, your first http request must be [rest.Interactions.CreateInteractionResponse] or [rest.Interactions.CreateInteractionResponseWithCallback]
//
// This does not produce a visible loading state to the user.
// You are expected to send a new http request within 3 seconds to respond to the interaction.
// This allows you to gracefully handle errors with your sent response & access the resulting message.
//
// If you want to create a visible loading state, use DeferCreateMessage.
//
// Source docs: [Discord Source docs]
//
// [Discord Source docs]: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-callback
func (e *AutocompleteInteractionCreate) Acknowledge(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeAcknowledge, nil, opts...)
}

// AutocompleteResult responds to the interaction with a slice of choices.
func (e *AutocompleteInteractionCreate) AutocompleteResult(choices []discord.AutocompleteChoice, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeAutocompleteResult, discord.AutocompleteResult{Choices: choices}, opts...)
}

// ModalSubmitInteractionCreate indicates that a new modal submit interaction has been created.
type ModalSubmitInteractionCreate struct {
	*GenericEvent
	discord.ModalSubmitInteraction
	RespondFunc  InteractionResponderFunc
	acknowledged bool
}

func (e *ModalSubmitInteractionCreate) Respond(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error {
	if e.acknowledged {
		return nil
	}
	if err := e.RespondFunc(responseType, data, opts...); err != nil {
		return err
	}
	e.acknowledged = true
	return nil
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *ModalSubmitInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches.Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

// Acknowledge acknowledges the interaction.
//
// This is used strictly for acknowledging the HTTP interaction request from discord. This responds with 202 Accepted.
//
// When using this, your first http request must be [rest.Interactions.CreateInteractionResponse] or [rest.Interactions.CreateInteractionResponseWithCallback]
//
// This does not produce a visible loading state to the user.
// You are expected to send a new http request within 3 seconds to respond to the interaction.
// This allows you to gracefully handle errors with your sent response & access the resulting message.
//
// If you want to create a visible loading state, use DeferCreateMessage.
//
// Source docs: [Discord Source docs]
//
// [Discord Source docs]: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-callback
func (e *ModalSubmitInteractionCreate) Acknowledge(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeAcknowledge, nil, opts...)
}

// CreateMessage responds to the interaction with a new message.
func (e *ModalSubmitInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e *ModalSubmitInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// UpdateMessage responds to the interaction with updating the message the component is from.
func (e *ModalSubmitInteractionCreate) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeUpdateMessage, messageUpdate, opts...)
}

// DeferUpdateMessage responds to the interaction with nothing.
func (e *ModalSubmitInteractionCreate) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil, opts...)
}

func (e *ModalSubmitInteractionCreate) RespondMessage(messageBuilder discord.MessageBuilder, opts ...rest.RequestOpt) error {
	if e.acknowledged {
		_, err := e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), messageBuilder.BuildUpdate(), opts...)
		return err
	} else {
		return e.CreateMessage(messageBuilder.BuildCreate())
	}
}

// LaunchActivity responds to the interaction by launching activity associated with the app.
func (e *ModalSubmitInteractionCreate) LaunchActivity(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeLaunchActivity, nil, opts...)
}
