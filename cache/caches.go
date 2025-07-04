package cache

import (
	"iter"
	"slices"
	"sync"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type SelfUserCache interface {
	SelfUser() (discord.OAuth2User, bool)
	SetSelfUser(selfUser discord.OAuth2User)
}

func NewSelfUserCache() SelfUserCache {
	return &selfUserCacheImpl{}
}

type selfUserCacheImpl struct {
	selfUserMu sync.Mutex
	selfUser   *discord.OAuth2User
}

func (c *selfUserCacheImpl) SelfUser() (discord.OAuth2User, bool) {
	c.selfUserMu.Lock()
	defer c.selfUserMu.Unlock()

	if c.selfUser == nil {
		return discord.OAuth2User{}, false
	}
	return *c.selfUser, true
}

func (c *selfUserCacheImpl) SetSelfUser(user discord.OAuth2User) {
	c.selfUserMu.Lock()
	defer c.selfUserMu.Unlock()

	c.selfUser = &user
}

type GuildCache interface {
	GuildCache() Cache[discord.Guild]

	IsGuildUnready(guildID snowflake.ID) bool
	SetGuildUnready(guildID snowflake.ID, unready bool)
	UnreadyGuildIDs() []snowflake.ID

	IsGuildUnavailable(guildID snowflake.ID) bool
	SetGuildUnavailable(guildID snowflake.ID, unavailable bool)
	UnavailableGuildIDs() []snowflake.ID

	Guild(guildID snowflake.ID) (discord.Guild, bool)
	Guilds() iter.Seq[discord.Guild]
	GuildsLen() int
	AddGuild(guild discord.Guild)
	RemoveGuild(guildID snowflake.ID) (discord.Guild, bool)
}

func NewGuildCache(cache Cache[discord.Guild], unreadyGuilds Set[snowflake.ID], unavailableGuilds Set[snowflake.ID]) GuildCache {
	return &guildCacheImpl{
		cache:             cache,
		unreadyGuilds:     unreadyGuilds,
		unavailableGuilds: unavailableGuilds,
	}
}

type guildCacheImpl struct {
	cache             Cache[discord.Guild]
	unreadyGuilds     Set[snowflake.ID]
	unavailableGuilds Set[snowflake.ID]
}

func (c *guildCacheImpl) GuildCache() Cache[discord.Guild] {
	return c.cache
}

func (c *guildCacheImpl) IsGuildUnready(guildID snowflake.ID) bool {
	return c.unreadyGuilds.Has(guildID)
}

func (c *guildCacheImpl) SetGuildUnready(guildID snowflake.ID, unready bool) {
	if c.unreadyGuilds.Has(guildID) && !unready {
		c.unreadyGuilds.Remove(guildID)
	} else if !c.unreadyGuilds.Has(guildID) && unready {
		c.unreadyGuilds.Add(guildID)
	}
}

func (c *guildCacheImpl) UnreadyGuildIDs() []snowflake.ID {
	var guilds []snowflake.ID
	for guildID := range c.unreadyGuilds.All() {
		guilds = append(guilds, guildID)
	}
	return guilds
}

func (c *guildCacheImpl) IsGuildUnavailable(guildID snowflake.ID) bool {
	return c.unavailableGuilds.Has(guildID)
}

func (c *guildCacheImpl) SetGuildUnavailable(guildID snowflake.ID, unavailable bool) {
	if c.unavailableGuilds.Has(guildID) && unavailable {
		c.unavailableGuilds.Remove(guildID)
	} else if !c.unavailableGuilds.Has(guildID) && !unavailable {
		c.unavailableGuilds.Add(guildID)
	}
}

func (c *guildCacheImpl) UnavailableGuildIDs() []snowflake.ID {
	var guilds []snowflake.ID
	for guildId := range c.unavailableGuilds.All() {
		guilds = append(guilds, guildId)
	}
	return guilds
}

func (c *guildCacheImpl) Guild(guildID snowflake.ID) (discord.Guild, bool) {
	return c.cache.Get(guildID)
}

func (c *guildCacheImpl) Guilds() iter.Seq[discord.Guild] {
	return c.cache.All()
}

func (c *guildCacheImpl) GuildsLen() int {
	return c.cache.Len()
}

func (c *guildCacheImpl) AddGuild(guild discord.Guild) {
	c.cache.Put(guild.ID, guild)
}

func (c *guildCacheImpl) RemoveGuild(guildID snowflake.ID) (discord.Guild, bool) {
	return c.cache.Remove(guildID)
}

type ChannelCache interface {
	ChannelCache() Cache[discord.GuildChannel]

	Channel(channelID snowflake.ID) (discord.GuildChannel, bool)
	Channels() iter.Seq[discord.GuildChannel]
	ChannelsLen() int
	AddChannel(channel discord.GuildChannel)
	RemoveChannel(channelID snowflake.ID) (discord.GuildChannel, bool)
	RemoveChannelsByGuildID(guildID snowflake.ID)
}

func NewChannelCache(cache Cache[discord.GuildChannel]) ChannelCache {
	return &channelCacheImpl{
		cache: cache,
	}
}

type channelCacheImpl struct {
	cache Cache[discord.GuildChannel]
}

func (c *channelCacheImpl) ChannelCache() Cache[discord.GuildChannel] {
	return c.cache
}

func (c *channelCacheImpl) Channel(channelID snowflake.ID) (discord.GuildChannel, bool) {
	return c.cache.Get(channelID)
}

func (c *channelCacheImpl) Channels() iter.Seq[discord.GuildChannel] {
	return c.cache.All()
}

func (c *channelCacheImpl) ChannelsLen() int {
	return c.cache.Len()
}

func (c *channelCacheImpl) AddChannel(channel discord.GuildChannel) {
	c.cache.Put(channel.ID(), channel)
}

func (c *channelCacheImpl) RemoveChannel(channelID snowflake.ID) (discord.GuildChannel, bool) {
	return c.cache.Remove(channelID)
}

func (c *channelCacheImpl) RemoveChannelsByGuildID(guildID snowflake.ID) {
	c.cache.RemoveIf(func(channel discord.GuildChannel) bool {
		return channel.GuildID() == guildID
	})
}

type StageInstanceCache interface {
	StageInstanceCache() GroupedCache[discord.StageInstance]

	StageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID) (discord.StageInstance, bool)
	StageInstances(guildID snowflake.ID) iter.Seq[discord.StageInstance]
	StageInstancesAllLen() int
	StageInstancesLen(guildID snowflake.ID) int
	AddStageInstance(stageInstance discord.StageInstance)
	RemoveStageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID) (discord.StageInstance, bool)
	RemoveStageInstancesByGuildID(guildID snowflake.ID)
}

func NewStageInstanceCache(cache GroupedCache[discord.StageInstance]) StageInstanceCache {
	return &stageInstanceCacheImpl{
		cache: cache,
	}
}

type stageInstanceCacheImpl struct {
	cache GroupedCache[discord.StageInstance]
}

func (c *stageInstanceCacheImpl) StageInstanceCache() GroupedCache[discord.StageInstance] {
	return c.cache
}

func (c *stageInstanceCacheImpl) StageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID) (discord.StageInstance, bool) {
	return c.cache.Get(guildID, stageInstanceID)
}

func (c *stageInstanceCacheImpl) StageInstances(guildID snowflake.ID) iter.Seq[discord.StageInstance] {
	return c.cache.GroupAll(guildID)
}

func (c *stageInstanceCacheImpl) StageInstancesAllLen() int {
	return c.cache.Len()
}

func (c *stageInstanceCacheImpl) StageInstancesLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *stageInstanceCacheImpl) AddStageInstance(stageInstance discord.StageInstance) {
	c.cache.Put(stageInstance.GuildID, stageInstance.ID, stageInstance)
}

func (c *stageInstanceCacheImpl) RemoveStageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID) (discord.StageInstance, bool) {
	return c.cache.Remove(guildID, stageInstanceID)
}

func (c *stageInstanceCacheImpl) RemoveStageInstancesByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type GuildScheduledEventCache interface {
	GuildScheduledEventCache() GroupedCache[discord.GuildScheduledEvent]

	GuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID) (discord.GuildScheduledEvent, bool)
	GuildScheduledEvents(guildID snowflake.ID) iter.Seq[discord.GuildScheduledEvent]
	GuildScheduledEventsAllLen() int
	GuildScheduledEventsLen(guildID snowflake.ID) int
	AddGuildScheduledEvent(guildScheduledEvent discord.GuildScheduledEvent)
	RemoveGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID) (discord.GuildScheduledEvent, bool)
	RemoveGuildScheduledEventsByGuildID(guildID snowflake.ID)
}

func NewGuildScheduledEventCache(cache GroupedCache[discord.GuildScheduledEvent]) GuildScheduledEventCache {
	return &guildScheduledEventCacheImpl{
		cache: cache,
	}
}

type guildScheduledEventCacheImpl struct {
	cache GroupedCache[discord.GuildScheduledEvent]
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEventCache() GroupedCache[discord.GuildScheduledEvent] {
	return c.cache
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID) (discord.GuildScheduledEvent, bool) {
	return c.cache.Get(guildID, guildScheduledEventID)
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEvents(guildID snowflake.ID) iter.Seq[discord.GuildScheduledEvent] {
	return c.cache.GroupAll(guildID)
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEventsAllLen() int {
	return c.cache.Len()
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEventsLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *guildScheduledEventCacheImpl) AddGuildScheduledEvent(guildScheduledEvent discord.GuildScheduledEvent) {
	c.cache.Put(guildScheduledEvent.GuildID, guildScheduledEvent.ID, guildScheduledEvent)
}

func (c *guildScheduledEventCacheImpl) RemoveGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID) (discord.GuildScheduledEvent, bool) {
	return c.cache.Remove(guildID, guildScheduledEventID)
}

func (c *guildScheduledEventCacheImpl) RemoveGuildScheduledEventsByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type GuildSoundboardSoundCache interface {
	GuildSoundboardSoundCache() GroupedCache[discord.SoundboardSound]
	GuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID) (discord.SoundboardSound, bool)
	GuildSoundboardSounds(guildID snowflake.ID) iter.Seq[discord.SoundboardSound]
	GuildSoundboardSoundsAllLen() int
	GuildSoundboardSoundsLen(guildID snowflake.ID) int
	AddGuildSoundboardSound(sound discord.SoundboardSound)
	RemoveGuildSoundboardSound(guildID snowflake.ID, sound snowflake.ID) (discord.SoundboardSound, bool)
	RemoveGuildSoundboardSoundsByGuildID(guildID snowflake.ID)
}

func NewGuildSoundboardSoundCache(cache GroupedCache[discord.SoundboardSound]) GuildSoundboardSoundCache {
	return &guildSoundboardSoundCacheImpl{
		cache: cache,
	}
}

type guildSoundboardSoundCacheImpl struct {
	cache GroupedCache[discord.SoundboardSound]
}

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSoundCache() GroupedCache[discord.SoundboardSound] {
	return c.cache
}

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID) (discord.SoundboardSound, bool) {
	return c.cache.Get(guildID, soundID)
}

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSounds(guildID snowflake.ID) iter.Seq[discord.SoundboardSound] {
	return c.cache.GroupAll(guildID)
}

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSoundsAllLen() int {
	return c.cache.Len()
}

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSoundsLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *guildSoundboardSoundCacheImpl) AddGuildSoundboardSound(sound discord.SoundboardSound) {
	c.cache.Put(*sound.GuildID, sound.SoundID, sound)
}

func (c *guildSoundboardSoundCacheImpl) RemoveGuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID) (discord.SoundboardSound, bool) {
	return c.cache.Remove(guildID, soundID)
}

func (c *guildSoundboardSoundCacheImpl) RemoveGuildSoundboardSoundsByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type RoleCache interface {
	RoleCache() GroupedCache[discord.Role]

	Role(guildID snowflake.ID, roleID snowflake.ID) (discord.Role, bool)
	Roles(guildID snowflake.ID) iter.Seq[discord.Role]
	RolesAllLen() int
	RolesLen(guildID snowflake.ID) int
	AddRole(role discord.Role)
	RemoveRole(guildID snowflake.ID, roleID snowflake.ID) (discord.Role, bool)
	RemoveRolesByGuildID(guildID snowflake.ID)
}

func NewRoleCache(cache GroupedCache[discord.Role]) RoleCache {
	return &roleCacheImpl{
		cache: cache,
	}
}

type roleCacheImpl struct {
	cache GroupedCache[discord.Role]
}

func (c *roleCacheImpl) RoleCache() GroupedCache[discord.Role] {
	return c.cache
}

func (c *roleCacheImpl) Role(guildID snowflake.ID, roleID snowflake.ID) (discord.Role, bool) {
	return c.cache.Get(guildID, roleID)
}

func (c *roleCacheImpl) Roles(guildID snowflake.ID) iter.Seq[discord.Role] {
	return c.cache.GroupAll(guildID)
}

func (c *roleCacheImpl) RolesAllLen() int {
	return c.cache.Len()
}

func (c *roleCacheImpl) RolesLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *roleCacheImpl) AddRole(role discord.Role) {
	c.cache.Put(role.GuildID, role.ID, role)
}

func (c *roleCacheImpl) RemoveRole(guildID snowflake.ID, roleID snowflake.ID) (discord.Role, bool) {
	return c.cache.Remove(guildID, roleID)
}

func (c *roleCacheImpl) RemoveRolesByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type MemberCache interface {
	MemberCache() GroupedCache[discord.Member]

	Member(guildID snowflake.ID, userID snowflake.ID) (discord.Member, bool)
	Members(guildID snowflake.ID) iter.Seq[discord.Member]
	MembersAllLen() int
	MembersLen(guildID snowflake.ID) int
	AddMember(member discord.Member)
	RemoveMember(guildID snowflake.ID, userID snowflake.ID) (discord.Member, bool)
	RemoveMembersByGuildID(guildID snowflake.ID)
}

func NewMemberCache(cache GroupedCache[discord.Member]) MemberCache {
	return &memberCacheImpl{
		cache: cache,
	}
}

type memberCacheImpl struct {
	cache GroupedCache[discord.Member]
}

func (c *memberCacheImpl) MemberCache() GroupedCache[discord.Member] {
	return c.cache
}

func (c *memberCacheImpl) Member(guildID snowflake.ID, userID snowflake.ID) (discord.Member, bool) {
	return c.cache.Get(guildID, userID)
}

func (c *memberCacheImpl) Members(guildID snowflake.ID) iter.Seq[discord.Member] {
	return c.cache.GroupAll(guildID)
}

func (c *memberCacheImpl) MembersAllLen() int {
	return c.cache.Len()
}

func (c *memberCacheImpl) MembersLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *memberCacheImpl) AddMember(member discord.Member) {
	c.cache.Put(member.GuildID, member.User.ID, member)
}

func (c *memberCacheImpl) RemoveMember(guildID snowflake.ID, userID snowflake.ID) (discord.Member, bool) {
	return c.cache.Remove(guildID, userID)
}

func (c *memberCacheImpl) RemoveMembersByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type ThreadMemberCache interface {
	ThreadMemberCache() GroupedCache[discord.ThreadMember]

	ThreadMember(threadID snowflake.ID, userID snowflake.ID) (discord.ThreadMember, bool)
	ThreadMembers(threadID snowflake.ID) iter.Seq[discord.ThreadMember]
	ThreadMembersAllLen() int
	ThreadMembersLen(guildID snowflake.ID) int
	AddThreadMember(threadMember discord.ThreadMember)
	RemoveThreadMember(threadID snowflake.ID, userID snowflake.ID) (discord.ThreadMember, bool)
	RemoveThreadMembersByThreadID(threadID snowflake.ID)
}

func NewThreadMemberCache(cache GroupedCache[discord.ThreadMember]) ThreadMemberCache {
	return &threadMemberCacheImpl{
		cache: cache,
	}
}

type threadMemberCacheImpl struct {
	cache GroupedCache[discord.ThreadMember]
}

func (c *threadMemberCacheImpl) ThreadMemberCache() GroupedCache[discord.ThreadMember] {
	return c.cache
}

func (c *threadMemberCacheImpl) ThreadMember(threadID snowflake.ID, userID snowflake.ID) (discord.ThreadMember, bool) {
	return c.cache.Get(threadID, userID)
}

func (c *threadMemberCacheImpl) ThreadMembers(threadID snowflake.ID) iter.Seq[discord.ThreadMember] {
	return c.cache.GroupAll(threadID)
}

func (c *threadMemberCacheImpl) ThreadMembersAllLen() int {
	return c.cache.Len()
}

func (c *threadMemberCacheImpl) ThreadMembersLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *threadMemberCacheImpl) AddThreadMember(threadMember discord.ThreadMember) {
	c.cache.Put(threadMember.ThreadID, threadMember.UserID, threadMember)
}

func (c *threadMemberCacheImpl) RemoveThreadMember(threadID snowflake.ID, userID snowflake.ID) (discord.ThreadMember, bool) {
	return c.cache.Remove(threadID, userID)
}

func (c *threadMemberCacheImpl) RemoveThreadMembersByThreadID(threadID snowflake.ID) {
	c.cache.GroupRemove(threadID)
}

type PresenceCache interface {
	PresenceCache() GroupedCache[discord.Presence]

	Presence(guildID snowflake.ID, userID snowflake.ID) (discord.Presence, bool)
	Presences(guildID snowflake.ID) iter.Seq[discord.Presence]
	PresencesAllLen() int
	PresencesLen(guildID snowflake.ID) int
	AddPresence(presence discord.Presence)
	RemovePresence(guildID snowflake.ID, userID snowflake.ID) (discord.Presence, bool)
	RemovePresencesByGuildID(guildID snowflake.ID)
}

func NewPresenceCache(cache GroupedCache[discord.Presence]) PresenceCache {
	return &presenceCacheImpl{
		cache: cache,
	}
}

type presenceCacheImpl struct {
	cache GroupedCache[discord.Presence]
}

func (c *presenceCacheImpl) PresenceCache() GroupedCache[discord.Presence] {
	return c.cache
}

func (c *presenceCacheImpl) Presence(guildID snowflake.ID, userID snowflake.ID) (discord.Presence, bool) {
	return c.cache.Get(guildID, userID)
}

func (c *presenceCacheImpl) Presences(guildID snowflake.ID) iter.Seq[discord.Presence] {
	return c.cache.GroupAll(guildID)
}

func (c *presenceCacheImpl) PresencesAllLen() int {
	return c.cache.Len()
}

func (c *presenceCacheImpl) PresencesLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *presenceCacheImpl) AddPresence(presence discord.Presence) {
	c.cache.Put(presence.GuildID, presence.PresenceUser.ID, presence)
}

func (c *presenceCacheImpl) RemovePresence(guildID snowflake.ID, userID snowflake.ID) (discord.Presence, bool) {
	return c.cache.Remove(guildID, userID)
}

func (c *presenceCacheImpl) RemovePresencesByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type VoiceStateCache interface {
	VoiceStateCache() GroupedCache[discord.VoiceState]

	VoiceState(guildID snowflake.ID, userID snowflake.ID) (discord.VoiceState, bool)
	VoiceStates(guildID snowflake.ID) iter.Seq[discord.VoiceState]
	VoiceStatesAllLen() int
	VoiceStatesLen(guildID snowflake.ID) int
	AddVoiceState(voiceState discord.VoiceState)
	RemoveVoiceState(guildID snowflake.ID, userID snowflake.ID) (discord.VoiceState, bool)
	RemoveVoiceStatesByGuildID(guildID snowflake.ID)
}

func NewVoiceStateCache(cache GroupedCache[discord.VoiceState]) VoiceStateCache {
	return &voiceStateCacheImpl{
		cache: cache,
	}
}

type voiceStateCacheImpl struct {
	cache GroupedCache[discord.VoiceState]
}

func (c *voiceStateCacheImpl) VoiceStateCache() GroupedCache[discord.VoiceState] {
	return c.cache
}

func (c *voiceStateCacheImpl) VoiceState(guildID snowflake.ID, userID snowflake.ID) (discord.VoiceState, bool) {
	return c.cache.Get(guildID, userID)
}

func (c *voiceStateCacheImpl) VoiceStates(guildID snowflake.ID) iter.Seq[discord.VoiceState] {
	return c.cache.GroupAll(guildID)
}

func (c *voiceStateCacheImpl) VoiceStatesAllLen() int {
	return c.cache.Len()
}

func (c *voiceStateCacheImpl) VoiceStatesLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *voiceStateCacheImpl) AddVoiceState(voiceState discord.VoiceState) {
	c.cache.Put(voiceState.GuildID, voiceState.UserID, voiceState)
}

func (c *voiceStateCacheImpl) RemoveVoiceState(guildID snowflake.ID, userID snowflake.ID) (discord.VoiceState, bool) {
	return c.cache.Remove(guildID, userID)
}

func (c *voiceStateCacheImpl) RemoveVoiceStatesByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type MessageCache interface {
	MessageCache() GroupedCache[discord.Message]

	Message(channelID snowflake.ID, messageID snowflake.ID) (discord.Message, bool)
	Messages(channelID snowflake.ID) iter.Seq[discord.Message]
	MessagesAllLen() int
	MessagesLen(guildID snowflake.ID) int
	AddMessage(message discord.Message)
	RemoveMessage(channelID snowflake.ID, messageID snowflake.ID) (discord.Message, bool)
	RemoveMessagesByChannelID(channelID snowflake.ID)
	RemoveMessagesByGuildID(guildID snowflake.ID)
}

func NewMessageCache(cache GroupedCache[discord.Message]) MessageCache {
	return &messageCacheImpl{
		cache: cache,
	}
}

type messageCacheImpl struct {
	cache GroupedCache[discord.Message]
}

func (c *messageCacheImpl) MessageCache() GroupedCache[discord.Message] {
	return c.cache
}

func (c *messageCacheImpl) Message(channelID snowflake.ID, messageID snowflake.ID) (discord.Message, bool) {
	return c.cache.Get(channelID, messageID)
}

func (c *messageCacheImpl) Messages(channelID snowflake.ID) iter.Seq[discord.Message] {
	return c.cache.GroupAll(channelID)
}

func (c *messageCacheImpl) MessagesAllLen() int {
	return c.cache.Len()
}

func (c *messageCacheImpl) MessagesLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *messageCacheImpl) AddMessage(message discord.Message) {
	c.cache.Put(message.ChannelID, message.ID, message)
}

func (c *messageCacheImpl) RemoveMessage(channelID snowflake.ID, messageID snowflake.ID) (discord.Message, bool) {
	return c.cache.Remove(channelID, messageID)
}

func (c *messageCacheImpl) RemoveMessagesByChannelID(channelID snowflake.ID) {
	c.cache.GroupRemove(channelID)
}

func (c *messageCacheImpl) RemoveMessagesByGuildID(guildID snowflake.ID) {
	c.cache.RemoveIf(func(_ snowflake.ID, message discord.Message) bool {
		return message.GuildID != nil && *message.GuildID == guildID
	})
}

type EmojiCache interface {
	EmojiCache() GroupedCache[discord.Emoji]

	Emoji(guildID snowflake.ID, emojiID snowflake.ID) (discord.Emoji, bool)
	Emojis(guildID snowflake.ID) iter.Seq[discord.Emoji]
	EmojisAllLen() int
	EmojisLen(guildID snowflake.ID) int
	AddEmoji(emoji discord.Emoji)
	RemoveEmoji(guildID snowflake.ID, emojiID snowflake.ID) (discord.Emoji, bool)
	RemoveEmojisByGuildID(guildID snowflake.ID)
}

func NewEmojiCache(cache GroupedCache[discord.Emoji]) EmojiCache {
	return &emojiCacheImpl{
		cache: cache,
	}
}

type emojiCacheImpl struct {
	cache GroupedCache[discord.Emoji]
}

func (c *emojiCacheImpl) EmojiCache() GroupedCache[discord.Emoji] {
	return c.cache
}

func (c *emojiCacheImpl) Emoji(guildID snowflake.ID, emojiID snowflake.ID) (discord.Emoji, bool) {
	return c.cache.Get(guildID, emojiID)
}

func (c *emojiCacheImpl) Emojis(guildID snowflake.ID) iter.Seq[discord.Emoji] {
	return c.cache.GroupAll(guildID)
}

func (c *emojiCacheImpl) EmojisAllLen() int {
	return c.cache.Len()
}

func (c *emojiCacheImpl) EmojisLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *emojiCacheImpl) AddEmoji(emoji discord.Emoji) {
	c.cache.Put(emoji.GuildID, emoji.ID, emoji)
}

func (c *emojiCacheImpl) RemoveEmoji(guildID snowflake.ID, emojiID snowflake.ID) (discord.Emoji, bool) {
	return c.cache.Remove(guildID, emojiID)
}

func (c *emojiCacheImpl) RemoveEmojisByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type StickerCache interface {
	StickerCache() GroupedCache[discord.Sticker]

	Sticker(guildID snowflake.ID, stickerID snowflake.ID) (discord.Sticker, bool)
	Stickers(guildID snowflake.ID) iter.Seq[discord.Sticker]
	StickersAllLen() int
	StickersLen(guildID snowflake.ID) int
	AddSticker(sticker discord.Sticker)
	RemoveSticker(guildID snowflake.ID, stickerID snowflake.ID) (discord.Sticker, bool)
	RemoveStickersByGuildID(guildID snowflake.ID)
}

func NewStickerCache(cache GroupedCache[discord.Sticker]) StickerCache {
	return &stickerCacheImpl{
		cache: cache,
	}
}

type stickerCacheImpl struct {
	cache GroupedCache[discord.Sticker]
}

func (c *stickerCacheImpl) StickerCache() GroupedCache[discord.Sticker] {
	return c.cache
}

func (c *stickerCacheImpl) Sticker(guildID snowflake.ID, stickerID snowflake.ID) (discord.Sticker, bool) {
	return c.cache.Get(guildID, stickerID)
}

func (c *stickerCacheImpl) Stickers(guildID snowflake.ID) iter.Seq[discord.Sticker] {
	return c.cache.GroupAll(guildID)
}

func (c *stickerCacheImpl) StickersAllLen() int {
	return c.cache.Len()
}

func (c *stickerCacheImpl) StickersLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *stickerCacheImpl) AddSticker(sticker discord.Sticker) {
	if sticker.GuildID == nil {
		return
	}
	c.cache.Put(*sticker.GuildID, sticker.ID, sticker)
}

func (c *stickerCacheImpl) RemoveSticker(guildID snowflake.ID, stickerID snowflake.ID) (discord.Sticker, bool) {
	return c.cache.Remove(guildID, stickerID)
}

func (c *stickerCacheImpl) RemoveStickersByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

// Caches combines all different entity caches into one with some utility methods.
type Caches interface {
	SelfUserCache
	GuildCache
	ChannelCache
	StageInstanceCache
	GuildScheduledEventCache
	GuildSoundboardSoundCache
	RoleCache
	MemberCache
	ThreadMemberCache
	PresenceCache
	VoiceStateCache
	MessageCache
	EmojiCache
	StickerCache

	// CacheFlags returns the current configured FLags of the caches.
	CacheFlags() Flags

	// MemberPermissions returns the calculated permissions of the given member.
	// This requires the FlagRoles to be set.
	MemberPermissions(member discord.Member) discord.Permissions

	// MemberPermissionsInChannel returns the calculated permissions of the given member in the given channel.
	// This requires the FlagRoles and FlagChannels to be set.
	MemberPermissionsInChannel(channel discord.GuildChannel, member discord.Member) discord.Permissions

	// MemberRoles returns all roles of the given member.
	// This requires the FlagRoles to be set.
	MemberRoles(member discord.Member) []discord.Role

	// AudioChannelMembers returns all members which are in the given audio channel.
	// This requires the FlagVoiceStates to be set.
	AudioChannelMembers(channel discord.GuildAudioChannel) []discord.Member

	// SelfMember returns the current bot member from the given guildID.
	// This is only available after we received the gateway.EventTypeGuildCreate event for the given guildID.
	SelfMember(guildID snowflake.ID) (discord.Member, bool)

	// GuildThreadsInChannel returns all discord.GuildThread from the ChannelCache and a bool indicating if it exists.
	GuildThreadsInChannel(channelID snowflake.ID) []discord.GuildThread

	// GuildMessageChannel returns a discord.GuildMessageChannel from the ChannelCache and a bool indicating if it exists.
	GuildMessageChannel(channelID snowflake.ID) (discord.GuildMessageChannel, bool)

	// GuildThread returns a discord.GuildThread from the ChannelCache and a bool indicating if it exists.
	GuildThread(channelID snowflake.ID) (discord.GuildThread, bool)

	// GuildAudioChannel returns a discord.GetGuildAudioChannel from the ChannelCache and a bool indicating if it exists.
	GuildAudioChannel(channelID snowflake.ID) (discord.GuildAudioChannel, bool)

	// GuildTextChannel returns a discord.GuildTextChannel from the ChannelCache and a bool indicating if it exists.
	GuildTextChannel(channelID snowflake.ID) (discord.GuildTextChannel, bool)

	// GuildVoiceChannel returns a discord.GuildVoiceChannel from the ChannelCache and a bool indicating if it exists.
	GuildVoiceChannel(channelID snowflake.ID) (discord.GuildVoiceChannel, bool)

	// GuildCategoryChannel returns a discord.GuildCategoryChannel from the ChannelCache and a bool indicating if it exists.
	GuildCategoryChannel(channelID snowflake.ID) (discord.GuildCategoryChannel, bool)

	// GuildNewsChannel returns a discord.GuildNewsChannel from the ChannelCache and a bool indicating if it exists.
	GuildNewsChannel(channelID snowflake.ID) (discord.GuildNewsChannel, bool)

	// GuildNewsThread returns a discord.GuildThread from the ChannelCache and a bool indicating if it exists.
	GuildNewsThread(channelID snowflake.ID) (discord.GuildThread, bool)

	// GuildPublicThread returns a discord.GuildThread from the ChannelCache and a bool indicating if it exists.
	GuildPublicThread(channelID snowflake.ID) (discord.GuildThread, bool)

	// GuildPrivateThread returns a discord.GuildThread from the ChannelCache and a bool indicating if it exists.
	GuildPrivateThread(channelID snowflake.ID) (discord.GuildThread, bool)

	// GuildStageVoiceChannel returns a discord.GuildStageVoiceChannel from the ChannelCache and a bool indicating if it exists.
	GuildStageVoiceChannel(channelID snowflake.ID) (discord.GuildStageVoiceChannel, bool)

	// GuildForumChannel returns a discord.GuildForumChannel from the ChannelCache and a bool indicating if it exists.
	GuildForumChannel(channelID snowflake.ID) (discord.GuildForumChannel, bool)

	// GuildMediaChannel returns a discord.GuildMediaChannel from the ChannelCache and a bool indicating if it exists.
	GuildMediaChannel(channelID snowflake.ID) (discord.GuildMediaChannel, bool)
}

// New returns a new default Caches instance with the given ConfigOpt(s) applied.
func New(opts ...ConfigOpt) Caches {
	cfg := defaultConfig()
	cfg.apply(opts)

	return &cachesImpl{
		config:                    cfg,
		selfUserCache:             cfg.SelfUserCache,
		guildCache:                cfg.GuildCache,
		channelCache:              cfg.ChannelCache,
		stageInstanceCache:        cfg.StageInstanceCache,
		guildScheduledEventCache:  cfg.GuildScheduledEventCache,
		guildSoundboardSoundCache: cfg.GuildSoundboardSoundCache,
		roleCache:                 cfg.RoleCache,
		memberCache:               cfg.MemberCache,
		threadMemberCache:         cfg.ThreadMemberCache,
		presenceCache:             cfg.PresenceCache,
		voiceStateCache:           cfg.VoiceStateCache,
		messageCache:              cfg.MessageCache,
		emojiCache:                cfg.EmojiCache,
		stickerCache:              cfg.StickerCache,
	}
}

// these type aliases are needed to allow having the GuildCache, ChannelCache, etc. as methods on the cachesImpl struct
type (
	guildCache                = GuildCache
	channelCache              = ChannelCache
	stageInstanceCache        = StageInstanceCache
	guildScheduledEventCache  = GuildScheduledEventCache
	guildSoundboardSoundCache = GuildSoundboardSoundCache
	roleCache                 = RoleCache
	memberCache               = MemberCache
	threadMemberCache         = ThreadMemberCache
	presenceCache             = PresenceCache
	voiceStateCache           = VoiceStateCache
	messageCache              = MessageCache
	emojiCache                = EmojiCache
	stickerCache              = StickerCache
	selfUserCache             = SelfUserCache
)

type cachesImpl struct {
	config config

	guildCache
	channelCache
	stageInstanceCache
	guildScheduledEventCache
	guildSoundboardSoundCache
	roleCache
	memberCache
	threadMemberCache
	presenceCache
	voiceStateCache
	messageCache
	emojiCache
	stickerCache
	selfUserCache
}

func (c *cachesImpl) CacheFlags() Flags {
	return c.config.CacheFlags
}

func (c *cachesImpl) MemberPermissions(member discord.Member) discord.Permissions {
	if guild, ok := c.Guild(member.GuildID); ok && guild.OwnerID == member.User.ID {
		return discord.PermissionsAll
	}

	var permissions discord.Permissions
	if publicRole, ok := c.Role(member.GuildID, member.GuildID); ok {
		permissions = publicRole.Permissions
	}

	for _, role := range c.MemberRoles(member) {
		permissions = permissions.Add(role.Permissions)
		if permissions.Has(discord.PermissionAdministrator) {
			return discord.PermissionsAll
		}
	}
	if member.CommunicationDisabledUntil != nil && member.CommunicationDisabledUntil.After(time.Now()) {
		permissions &= discord.PermissionViewChannel | discord.PermissionReadMessageHistory
	}
	return permissions
}

func (c *cachesImpl) MemberPermissionsInChannel(channel discord.GuildChannel, member discord.Member) discord.Permissions {
	permissions := c.MemberPermissions(member)
	if permissions.Has(discord.PermissionAdministrator) {
		return discord.PermissionsAll
	}

	var (
		allow discord.Permissions
		deny  discord.Permissions
	)

	if overwrite, ok := channel.PermissionOverwrites().Role(channel.GuildID()); ok {
		permissions |= overwrite.Allow
		permissions &= ^overwrite.Deny
	}

	for _, roleID := range member.RoleIDs {
		if roleID == channel.GuildID() {
			continue
		}

		if overwrite, ok := channel.PermissionOverwrites().Role(roleID); ok {
			allow |= overwrite.Allow
			deny |= overwrite.Deny
		}
	}

	if overwrite, ok := channel.PermissionOverwrites().Member(member.User.ID); ok {
		allow |= overwrite.Allow
		deny |= overwrite.Deny
	}

	permissions &= ^deny
	permissions |= allow

	if member.CommunicationDisabledUntil != nil && member.CommunicationDisabledUntil.After(time.Now()) {
		permissions &= discord.PermissionViewChannel | discord.PermissionReadMessageHistory
	}

	return permissions
}

func (c *cachesImpl) MemberRoles(member discord.Member) []discord.Role {
	var roles []discord.Role

	for role := range c.Roles(member.GuildID) {
		if slices.Contains(member.RoleIDs, role.ID) {
			roles = append(roles, role)
		}
	}
	return roles
}

func (c *cachesImpl) AudioChannelMembers(channel discord.GuildAudioChannel) []discord.Member {
	var members []discord.Member
	for state := range c.VoiceStates(channel.GuildID()) {
		if member, ok := c.Member(channel.GuildID(), state.UserID); ok && state.ChannelID != nil && *state.ChannelID == channel.ID() {
			members = append(members, member)
		}
	}
	return members
}

func (c *cachesImpl) SelfMember(guildID snowflake.ID) (discord.Member, bool) {
	selfUser, ok := c.SelfUser()
	if !ok {
		return discord.Member{}, false
	}
	return c.Member(guildID, selfUser.ID)
}

func (c *cachesImpl) GuildThreadsInChannel(channelID snowflake.ID) []discord.GuildThread {
	var threads []discord.GuildThread
	for channel := range c.Channels() {
		if thread, ok := channel.(discord.GuildThread); ok && *thread.ParentID() == channelID {
			threads = append(threads, thread)
		}
	}
	return threads
}

func (c *cachesImpl) MessageChannel(channelID snowflake.ID) (discord.MessageChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(discord.MessageChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *cachesImpl) GuildMessageChannel(channelID snowflake.ID) (discord.GuildMessageChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if chM, ok := ch.(discord.GuildMessageChannel); ok {
			return chM, true
		}
	}
	return nil, false
}

func (c *cachesImpl) GuildThread(channelID snowflake.ID) (discord.GuildThread, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(discord.GuildThread); ok {
			return cCh, true
		}
	}
	return discord.GuildThread{}, false
}

func (c *cachesImpl) GuildAudioChannel(channelID snowflake.ID) (discord.GuildAudioChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(discord.GuildAudioChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *cachesImpl) GuildTextChannel(channelID snowflake.ID) (discord.GuildTextChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(discord.GuildTextChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildTextChannel{}, false
}

func (c *cachesImpl) GuildVoiceChannel(channelID snowflake.ID) (discord.GuildVoiceChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(discord.GuildVoiceChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildVoiceChannel{}, false
}

func (c *cachesImpl) GuildCategoryChannel(channelID snowflake.ID) (discord.GuildCategoryChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(discord.GuildCategoryChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildCategoryChannel{}, false
}

func (c *cachesImpl) GuildNewsChannel(channelID snowflake.ID) (discord.GuildNewsChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(discord.GuildNewsChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildNewsChannel{}, false
}

func (c *cachesImpl) GuildNewsThread(channelID snowflake.ID) (discord.GuildThread, bool) {
	if ch, ok := c.GuildThread(channelID); ok && ch.Type() == discord.ChannelTypeGuildNewsThread {
		return ch, true
	}
	return discord.GuildThread{}, false
}

func (c *cachesImpl) GuildPublicThread(channelID snowflake.ID) (discord.GuildThread, bool) {
	if ch, ok := c.GuildThread(channelID); ok && ch.Type() == discord.ChannelTypeGuildPublicThread {
		return ch, true
	}
	return discord.GuildThread{}, false
}

func (c *cachesImpl) GuildPrivateThread(channelID snowflake.ID) (discord.GuildThread, bool) {
	if ch, ok := c.GuildThread(channelID); ok && ch.Type() == discord.ChannelTypeGuildPrivateThread {
		return ch, true
	}
	return discord.GuildThread{}, false
}

func (c *cachesImpl) GuildStageVoiceChannel(channelID snowflake.ID) (discord.GuildStageVoiceChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(discord.GuildStageVoiceChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildStageVoiceChannel{}, false
}

func (c *cachesImpl) GuildForumChannel(channelID snowflake.ID) (discord.GuildForumChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(discord.GuildForumChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildForumChannel{}, false
}

func (c *cachesImpl) GuildMediaChannel(channelID snowflake.ID) (discord.GuildMediaChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(discord.GuildMediaChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildMediaChannel{}, false
}
