package core

import "github.com/DisgoOrg/disgo/discord"

// GetMemberPermissions returns all Permissions from the provided Member
func GetMemberPermissions(member *Member) discord.Permissions {
	if member.IsOwner() {
		return discord.PermissionsAll
	}

	var permissions discord.Permissions
	if publicRole := member.Bot.Caches.RoleCache().Get(member.GuildID, member.GuildID); publicRole != nil {
		permissions = publicRole.Permissions
	}

	for _, role := range member.Roles() {
		permissions = permissions.Add(role.Permissions)
		if permissions.Has(discord.PermissionAdministrator) {
			return discord.PermissionsAll
		}
	}
	if member.TimedOutUntil != nil {
		permissions &= discord.PermissionViewChannel | discord.PermissionReadMessageHistory
	}
	return permissions
}

func GetMemberPermissionsInChannel(channel *Channel, member *Member) discord.Permissions {
	if !channel.IsGuildChannel() {
		unsupportedChannelType(channel)
	}
	if *channel.GuildID != member.GuildID {
		panic("channel and member need to be part of the same guild")
	}

	if member.IsOwner() {
		return discord.PermissionsAll
	}
	permissions := GetMemberPermissions(member)
	if permissions.Has(discord.PermissionAdministrator) {
		return discord.PermissionsAll
	}

	var (
		allowRaw discord.Permissions
		denyRaw  discord.Permissions
	)
	if overwrite := channel.PermissionOverwrite(discord.PermissionOverwriteTypeRole, *channel.GuildID); overwrite != nil {
		allowRaw = overwrite.Allow
		denyRaw = overwrite.Deny
	}

	var (
		allowRole discord.Permissions
		denyRole  discord.Permissions
	)
	for _, roleID := range member.RoleIDs {
		if roleID == *channel.GuildID {
			continue
		}

		overwrite := channel.PermissionOverwrite(discord.PermissionOverwriteTypeRole, roleID)
		if overwrite == nil {
			break
		}
		allowRole = allowRole.Add(overwrite.Allow)
		denyRole = denyRole.Add(overwrite.Deny)
	}

	allowRaw = (allowRaw & (denyRole - 1)) | allowRole
	denyRaw = (denyRaw & (allowRole - 1)) | denyRole

	if overwrite := channel.PermissionOverwrite(discord.PermissionOverwriteTypeMember, member.ID); overwrite != nil {
		allowRaw = (allowRaw & (overwrite.Deny - 1)) | overwrite.Allow
		denyRaw = (denyRaw & (overwrite.Allow - 1)) | overwrite.Deny
	}

	permissions &= denyRaw - 1
	permissions |= allowRaw

	return permissions
}
