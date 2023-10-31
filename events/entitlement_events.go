package events

import "github.com/sabafly/sabafly-disgo/discord"

type GenericEntitlementEvent struct {
	*GenericEvent
	discord.Entitlement
}

type EntitlementCreate struct {
	*GenericEntitlementEvent
}

type EntitlementUpdate struct {
	*GenericEntitlementEvent
}

type EntitlementDelete struct {
	*GenericEntitlementEvent
}
