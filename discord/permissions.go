package discord

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/disgoorg/json/v2"

	"github.com/disgoorg/disgo/internal/flags"
)

// Permissions extends the Bit structure, and is used within roles and channels (https://discord.com/developers/docs/topics/permissions#permissions)
type Permissions int64

const (
	PermissionCreateInstantInvite Permissions = 1 << iota
	PermissionKickMembers
	PermissionBanMembers
	PermissionAdministrator
	PermissionManageChannels
	PermissionManageGuild
	PermissionAddReactions
	PermissionViewAuditLog
	PermissionPrioritySpeaker
	PermissionStream
	PermissionViewChannel
	PermissionSendMessages
	PermissionSendTTSMessages
	PermissionManageMessages
	PermissionEmbedLinks
	PermissionAttachFiles
	PermissionReadMessageHistory
	PermissionMentionEveryone
	PermissionUseExternalEmojis
	PermissionViewGuildInsights
	PermissionConnect
	PermissionSpeak
	PermissionMuteMembers
	PermissionDeafenMembers
	PermissionMoveMembers
	PermissionUseVAD
	PermissionChangeNickname
	PermissionManageNicknames
	PermissionManageRoles
	PermissionManageWebhooks
	PermissionManageGuildExpressions
	PermissionUseApplicationCommands
	PermissionRequestToSpeak
	PermissionManageEvents
	PermissionManageThreads
	PermissionCreatePublicThreads
	PermissionCreatePrivateThreads
	PermissionUseExternalStickers
	PermissionSendMessagesInThreads
	PermissionUseEmbeddedActivities
	PermissionModerateMembers
	PermissionViewCreatorMonetizationAnalytics
	PermissionUseSoundboard
	PermissionCreateGuildExpressions
	PermissionCreateEvents
	PermissionUseExternalSounds
	PermissionSendVoiceMessages
	_
	_
	PermissionSendPolls
	PermissionUseExternalApps

	PermissionsAllText = PermissionCreateInstantInvite |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionViewChannel |
		PermissionSendMessages |
		PermissionSendTTSMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone |
		PermissionUseExternalEmojis |
		PermissionManageRoles |
		PermissionManageWebhooks |
		PermissionUseApplicationCommands |
		PermissionManageThreads |
		PermissionCreatePublicThreads |
		PermissionCreatePrivateThreads |
		PermissionUseExternalStickers |
		PermissionSendMessagesInThreads |
		PermissionUseEmbeddedActivities |
		PermissionSendVoiceMessages |
		PermissionSendPolls |
		PermissionUseExternalApps

	PermissionsAllVoice = PermissionCreateInstantInvite |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionPrioritySpeaker |
		PermissionStream |
		PermissionViewChannel |
		PermissionSendMessages |
		PermissionSendTTSMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone |
		PermissionUseExternalEmojis |
		PermissionConnect |
		PermissionSpeak |
		PermissionStream |
		PermissionMuteMembers |
		PermissionDeafenMembers |
		PermissionMoveMembers |
		PermissionUseVAD |
		PermissionManageRoles |
		PermissionManageWebhooks |
		PermissionUseApplicationCommands |
		PermissionManageEvents |
		PermissionUseExternalStickers |
		PermissionUseEmbeddedActivities |
		PermissionUseSoundboard |
		PermissionUseExternalSounds |
		PermissionSendVoiceMessages |
		PermissionRequestToSpeak |
		PermissionUseEmbeddedActivities |
		PermissionCreateGuildExpressions |
		PermissionCreateEvents |
		PermissionSendPolls

	PermissionsAllStage = PermissionCreateInstantInvite |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionStream |
		PermissionViewChannel |
		PermissionSendMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone |
		PermissionUseExternalEmojis |
		PermissionConnect |
		PermissionMuteMembers |
		PermissionMoveMembers |
		PermissionManageRoles |
		PermissionManageWebhooks |
		PermissionUseApplicationCommands |
		PermissionRequestToSpeak |
		PermissionManageEvents |
		PermissionUseExternalStickers |
		PermissionSendVoiceMessages |
		PermissionSendPolls

	PermissionsAllForum = PermissionCreateInstantInvite |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionViewChannel |
		PermissionSendMessages |
		PermissionSendTTSMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone |
		PermissionUseExternalEmojis |
		PermissionManageRoles |
		PermissionManageWebhooks |
		PermissionUseApplicationCommands |
		PermissionManageThreads |
		PermissionUseExternalStickers |
		PermissionUseEmbeddedActivities |
		PermissionSendVoiceMessages |
		PermissionSendPolls

	PermissionsAllChannel = PermissionsAllText |
		PermissionsAllVoice |
		PermissionsAllStage |
		PermissionsAllForum

	PermissionsAllGuild = PermissionKickMembers |
		PermissionBanMembers |
		PermissionManageGuild |
		PermissionAdministrator |
		PermissionManageWebhooks |
		PermissionManageGuildExpressions |
		PermissionViewCreatorMonetizationAnalytics |
		PermissionViewGuildInsights |
		PermissionViewAuditLog |
		PermissionManageRoles |
		PermissionChangeNickname |
		PermissionManageNicknames |
		PermissionModerateMembers

	PermissionsAll = PermissionsAllChannel |
		PermissionsAllGuild

	PermissionsNone Permissions = 0
)

var permissions = map[Permissions]string{
	PermissionCreateInstantInvite:              "CREATE_INSTANT_INVITE",
	PermissionKickMembers:                      "KICK_MEMBERS",
	PermissionBanMembers:                       "BAN_MEMBERS",
	PermissionAdministrator:                    "ADMINISTRATOR",
	PermissionManageChannels:                   "MANAGE_CHANNELS",
	PermissionManageGuild:                      "MANAGE_GUILD",
	PermissionAddReactions:                     "ADD_REACTIONS",
	PermissionViewAuditLog:                     "VIEW_AUDIT_LOG",
	PermissionViewChannel:                      "VIEW_CHANNEL",
	PermissionSendMessages:                     "SEND_MESSAGES",
	PermissionSendTTSMessages:                  "SEND_TTS_MESSAGES",
	PermissionManageMessages:                   "MANAGE_MESSAGES",
	PermissionEmbedLinks:                       "EMBED_LINKS",
	PermissionAttachFiles:                      "ATTACH_FILES",
	PermissionReadMessageHistory:               "READ_MESSAGE_HISTORY",
	PermissionMentionEveryone:                  "MENTION_EVERYONE",
	PermissionUseExternalEmojis:                "USE_EXTERNAL_EMOJIS",
	PermissionConnect:                          "CONNECT",
	PermissionSpeak:                            "SPEAK",
	PermissionMuteMembers:                      "MUTE_MEMBERS",
	PermissionDeafenMembers:                    "DEAFEN_MEMBERS",
	PermissionMoveMembers:                      "MOVE_MEMBERS",
	PermissionUseVAD:                           "USE_VAD",
	PermissionPrioritySpeaker:                  "PRIORITY_SPEAKER",
	PermissionChangeNickname:                   "CHANGE_NICKNAME",
	PermissionManageNicknames:                  "MANAGE_NICKNAMES",
	PermissionManageRoles:                      "MANAGE_ROLES",
	PermissionManageWebhooks:                   "MANAGE_WEBHOOKS",
	PermissionManageGuildExpressions:           "MANAGE_GUILD_EXPRESSIONS",
	PermissionUseApplicationCommands:           "USE_APPLICATION_COMMANDS",
	PermissionRequestToSpeak:                   "REQUEST_TO_SPEAK",
	PermissionManageEvents:                     "MANAGE_EVENTS",
	PermissionManageThreads:                    "MANAGE_THREADS",
	PermissionCreatePublicThreads:              "CREATE_PUBLIC_THREADS",
	PermissionCreatePrivateThreads:             "CREATE_PRIVATE_THREADS",
	PermissionUseExternalStickers:              "USE_EXTERNAL_STICKERS",
	PermissionSendMessagesInThreads:            "SEND_MESSAGES_IN_THREADS",
	PermissionUseEmbeddedActivities:            "USE_EMBEDDED_ACTIVITIES",
	PermissionModerateMembers:                  "MODERATE_MEMBERS",
	PermissionViewCreatorMonetizationAnalytics: "VIEW_CREATOR_MONETIZATION_ANALYTICS",
	PermissionUseSoundboard:                    "USE_SOUNDBOARD",
	PermissionUseExternalSounds:                "USE_EXTERNAL_SOUNDS",
	PermissionStream:                           "STREAM",
	PermissionViewGuildInsights:                "VIEW_GUILD_INSIGHTS",
	PermissionSendVoiceMessages:                "SEND_VOICE_MESSAGES",
	PermissionSendPolls:                        "SEND_POLLS",
	PermissionUseExternalApps:                  "USE_EXTERNAL_APPS",
}

func (p Permissions) String() string {
	if p == PermissionsNone {
		return "None"
	}
	perms := new(strings.Builder)
	for permission, name := range permissions {
		if p.Has(permission) {
			perms.WriteString(name)
			perms.WriteString(", ")
		}
	}
	return perms.String()[:perms.Len()-2] // remove trailing comma and space
}

// MarshalJSON marshals permissions into a string
func (p Permissions) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(p), 10))
}

// UnmarshalJSON unmarshalls permissions into an int64
func (p *Permissions) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("")) || bytes.Equal(data, []byte("null")) {
		return nil
	}

	str, _ := strconv.Unquote(string(data))
	perms, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}

	*p = Permissions(perms)
	return nil
}

// Add allows you to add multiple bits together, producing a new bit
func (p Permissions) Add(bits ...Permissions) Permissions {
	return flags.Add(p, bits...)
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (p Permissions) Remove(bits ...Permissions) Permissions {
	return flags.Remove(p, bits...)
}

// Has will ensure that the bit includes all the bits entered
func (p Permissions) Has(bits ...Permissions) bool {
	return flags.Has(p, bits...)
}

// Missing will check whether the bit is missing any one of the bits
func (p Permissions) Missing(bits ...Permissions) bool {
	return flags.Missing(p, bits...)
}
