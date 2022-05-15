package bot

import (
	"context"

	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var _ Client = (*clientImpl)(nil)

type Client interface {
	Logger() log.Logger

	// Close will clean up all disgo internals and close the discord gracefully.
	Close(ctx context.Context)

	// Token returns the configured bot token.
	Token() string

	// ApplicationID returns the application id.
	ApplicationID() snowflake.ID

	// ID returns the bot id.
	ID() snowflake.ID

	// Caches returns the cache.Caches used by the Client.
	Caches() cache.Caches

	// Rest returns the rest.Rest used by the Client.
	Rest() rest.Rest

	AddEventListeners(listeners ...EventListener)
	RemoveEventListeners(listeners ...EventListener)

	// EventManager returns the EventManager used by the Client.
	EventManager() EventManager

	// ConnectGateway connects to the configured gateway.Gateway.
	ConnectGateway(ctx context.Context) error

	// Gateway returns the gateway.Gateway used by the Client.
	Gateway() gateway.Gateway

	// HasGateway returns whether the Client has a configured gateway.Gateway.
	HasGateway() bool

	// ConnectShardManager connects to the configured sharding.ShardManager.
	ConnectShardManager(ctx context.Context) error

	// ShardManager returns the sharding.ShardManager used by the Client.
	ShardManager() sharding.ShardManager

	// HasShardManager returns whether the Client has a configured sharding.ShardManager.
	HasShardManager() bool

	// Shard returns the gateway.Gateway the specific guildID runs on.
	Shard(guildID snowflake.ID) (gateway.Gateway, error)

	// Connect sends a discord.GatewayMessageDataVoiceStateUpdate to the specific gateway.Gateway and connects the bot to the specified channel.
	Connect(ctx context.Context, guildID snowflake.ID, channelID snowflake.ID) error

	// Disconnect sends a discord.GatewayMessageDataVoiceStateUpdate to the specific gateway.Gateway and disconnects the bot from this guild.
	Disconnect(ctx context.Context, guildID snowflake.ID) error

	// RequestMembers sends a discord.GatewayMessageDataRequestGuildMembers to the specific gateway.Gateway and requests the Member(s) of the specified guild.
	//  guildID  : is the snowflake of the guild to request the members of.
	//  presence : Weather or not to include discord.Presence data.
	//  nonce	 : The nonce to return to the discord.GatewayEventGuildMembersChunk.
	//  userIDs  : The snowflakes of the users to request the members of.
	RequestMembers(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, userIDs ...snowflake.ID) error

	// RequestMembersWithQuery sends a discord.GatewayMessageDataRequestGuildMembers to the specific gateway.Gateway and requests the Member(s) of the specified guild.
	//  guildID  : is the snowflake of the guild to request the members of.
	//  presence : Weather or not to include discord.Presence data.
	//  nonce    : The nonce to return to the discord.GatewayEventGuildMembersChunk.
	//  query    : The query to use for the request.
	//  limit    : The number of discord.Member(s) to return.
	RequestMembersWithQuery(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, query string, limit int) error

	// SetPresence sends a discord.GatewayMessageDataPresenceUpdate to the gateway.Gateway.
	SetPresence(ctx context.Context, presenceUpdate discord.GatewayMessageDataPresenceUpdate) error

	// SetPresenceForShard sends a discord.GatewayMessageDataPresenceUpdate to the specific gateway.Gateway.
	SetPresenceForShard(ctx context.Context, shardId int, presenceUpdate discord.GatewayMessageDataPresenceUpdate) error

	// MemberChunkingManager returns the MemberChunkingManager used by the Client.
	MemberChunkingManager() MemberChunkingManager

	// StartHTTPServer starts the configured HTTPServer used for interactions over webhooks.
	StartHTTPServer() error

	// HTTPServer returns the configured HTTPServer used for interactions over webhooks.
	HTTPServer() httpserver.Server

	// HasHTTPServer returns whether the Client has a configured HTTPServer.
	HasHTTPServer() bool
}

type clientImpl struct {
	token         string
	applicationID snowflake.ID

	logger log.Logger

	restServices rest.Rest

	eventManager EventManager

	shardManager sharding.ShardManager
	gateway      gateway.Gateway

	httpServer httpserver.Server

	caches cache.Caches

	memberChunkingManager MemberChunkingManager
}

func (c *clientImpl) Logger() log.Logger {
	return c.logger
}

func (c *clientImpl) Close(ctx context.Context) {
	if c.restServices != nil {
		c.restServices.Close(ctx)
	}
	if c.gateway != nil {
		c.gateway.Close(ctx)
	}
	if c.shardManager != nil {
		c.shardManager.Close(ctx)
	}
	if c.httpServer != nil {
		c.httpServer.Close(ctx)
	}
}

func (c *clientImpl) Token() string {
	return c.token
}

func (c *clientImpl) ApplicationID() snowflake.ID {
	return c.applicationID
}

func (c *clientImpl) ID() snowflake.ID {
	if selfUser, ok := c.Caches().GetSelfUser(); ok {
		return selfUser.ID
	}
	return 0
}

func (c *clientImpl) Caches() cache.Caches {
	return c.caches
}

func (c *clientImpl) Rest() rest.Rest {
	return c.restServices
}

// AddEventListeners adds one or more EventListener(s) to the EventManager
func (c *clientImpl) AddEventListeners(listeners ...EventListener) {
	c.eventManager.AddEventListeners(listeners...)
}

// RemoveEventListeners removes one or more EventListener(s) from the EventManager
func (c *clientImpl) RemoveEventListeners(listeners ...EventListener) {
	c.eventManager.RemoveEventListeners(listeners...)
}

func (c *clientImpl) EventManager() EventManager {
	return c.eventManager
}

// ConnectGateway opens the gateway connection to discord
func (c *clientImpl) ConnectGateway(ctx context.Context) error {
	if c.gateway == nil {
		return discord.ErrNoGateway
	}
	return c.gateway.Open(ctx)
}

func (c *clientImpl) Gateway() gateway.Gateway {
	return c.gateway
}

// HasGateway returns whether this Client has an active gateway.Gateway connection
func (c *clientImpl) HasGateway() bool {
	return c.gateway != nil
}

// ConnectShardManager opens the gateway connection to discord
func (c *clientImpl) ConnectShardManager(ctx context.Context) error {
	if c.shardManager == nil {
		return discord.ErrNoShardManager
	}
	c.shardManager.Open(ctx)
	return nil
}

func (c *clientImpl) ShardManager() sharding.ShardManager {
	return c.shardManager
}

// HasShardManager returns whether this Client is sharded
func (c *clientImpl) HasShardManager() bool {
	return c.shardManager != nil
}

func (c *clientImpl) Shard(guildID snowflake.ID) (gateway.Gateway, error) {
	if c.HasGateway() {
		return c.gateway, nil
	} else if c.HasShardManager() {
		if shard := c.shardManager.GetGuildShard(guildID); shard != nil {
			return shard, nil
		}
		return nil, discord.ErrShardNotFound
	}
	return nil, discord.ErrNoGatewayOrShardManager
}

func (c *clientImpl) Connect(ctx context.Context, guildID snowflake.ID, channelID snowflake.ID) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.GatewayOpcodeVoiceStateUpdate, discord.GatewayMessageDataVoiceStateUpdate{
		GuildID:   guildID,
		ChannelID: &channelID,
	})
}

func (c *clientImpl) Disconnect(ctx context.Context, guildID snowflake.ID) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.GatewayOpcodeVoiceStateUpdate, discord.GatewayMessageDataVoiceStateUpdate{
		GuildID:   guildID,
		ChannelID: nil,
	})
}

func (c *clientImpl) RequestMembers(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, userIDs ...snowflake.ID) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.GatewayOpcodeRequestGuildMembers, discord.GatewayMessageDataRequestGuildMembers{
		GuildID:   guildID,
		Presences: presence,
		UserIDs:   userIDs,
		Nonce:     nonce,
	})
}
func (c *clientImpl) RequestMembersWithQuery(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, query string, limit int) error {
	shard, err := c.Shard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.GatewayOpcodeRequestGuildMembers, discord.GatewayMessageDataRequestGuildMembers{
		GuildID:   guildID,
		Query:     &query,
		Limit:     &limit,
		Presences: presence,
		Nonce:     nonce,
	})
}

func (c *clientImpl) SetPresence(ctx context.Context, presenceUpdate discord.GatewayMessageDataPresenceUpdate) error {
	if !c.HasGateway() {
		return discord.ErrNoGateway
	}
	return c.gateway.Send(ctx, discord.GatewayOpcodePresenceUpdate, presenceUpdate)
}

// SetPresenceForShard sets the Presence of this Client for the provided shard
func (c *clientImpl) SetPresenceForShard(ctx context.Context, shardId int, presenceUpdate discord.GatewayMessageDataPresenceUpdate) error {
	if !c.HasShardManager() {
		return discord.ErrNoShardManager
	}
	shard := c.shardManager.Shard(shardId)
	if shard == nil {
		return discord.ErrShardNotFound
	}
	return shard.Send(ctx, discord.GatewayOpcodePresenceUpdate, presenceUpdate)
}

func (c *clientImpl) MemberChunkingManager() MemberChunkingManager {
	return c.memberChunkingManager
}

// StartHTTPServer starts the interaction webhook server
func (c *clientImpl) StartHTTPServer() error {
	if c.httpServer == nil {
		return discord.ErrNoHTTPServer
	}
	c.httpServer.Start()
	return nil
}

func (c *clientImpl) HTTPServer() httpserver.Server {
	return c.httpServer
}

// HasHTTPServer returns whether Client has an active httpserver.Server
func (c *clientImpl) HasHTTPServer() bool {
	return c.httpServer != nil
}