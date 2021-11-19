package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type GenericThreadEvent struct {
	*GenericEvent
	Thread core.GuildThread
	ThreadID discord.Snowflake
	GuildID discord.Snowflake
	ParentID discord.Snowflake
}

type ThreadCreateEvent struct {
	*GenericThreadEvent
}

type ThreadUpdateEvent struct {
	*GenericThreadEvent
	OldThread core.GuildThread
}

type ThreadDeleteEvent struct {
	*GenericThreadEvent
	ThreadID discord.Snowflake
	ParentID discord.Snowflake
}

type GenericThreadMemberEvent struct {
	*GenericThreadEvent
	ThreadMember *core.ThreadMember
}

type ThreadMemberJoinEvent struct {
	*GenericThreadMemberEvent
}

type ThreadMemberUpdateEvent struct {
	*GenericThreadMemberEvent
	OldThreadMember *core.ThreadMember
}

type ThreadMemberLeaveEvent struct {
	*GenericThreadMemberEvent
}
