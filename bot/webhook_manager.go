package bot

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
)

var (
	WebhookDefaultName                 = "disgo"
	WebhookDefaultAvatar *discord.Icon = nil
)

var _ WebhookManager = (*webhookManagerImpl)(nil)

func NewWebhookManager(client *Client, logger *slog.Logger) WebhookManager {
	return &webhookManagerImpl{
		client: client,
		logger: logger,
	}
}

type WebhookManager interface {
	GetMessenger(channel discord.Channel) (discord.WebhookMessenger, error)
	GetWebhook(channel discord.Channel) (discord.Webhook, error)
}

type webhookManagerImpl struct {
	client *Client
	logger *slog.Logger
}

func (w *webhookManagerImpl) GetMessenger(channel discord.Channel) (discord.WebhookMessenger, error) {
	if channel.Type() == discord.ChannelTypeDM || channel.Type() == discord.ChannelTypeGroupDM {
		return nil, discord.ErrUnsupportedType
	}

	if channel.Type().IsThread() {
		var c discord.GuildThread
	interaction:
		switch ch := channel.(type) {
		case discord.GuildThread:
			c = ch
		case discord.MessageThread:
			c = ch.GuildThread
		case discord.InteractionChannel:
			channel = ch.MessageChannel
			goto interaction
		default:
			return nil, discord.ErrUnsupportedType
		}
		return NewThreadWebhookMessenger(w.client, *c.ParentID(), c.ID())
	}

	return NewChannelWebhookMessenger(w.client, channel.ID())
}

func (w *webhookManagerImpl) GetWebhook(channel discord.Channel) (discord.Webhook, error) {
	m, err := w.GetMessenger(channel)
	if err != nil {
		return nil, err
	}
	return m.Webhook(), nil
}

var (
	_ discord.WebhookMessenger = (*channelWebhookMessenger)(nil)
	_ discord.WebhookMessenger = (*threadWebhookMessenger)(nil)
)

func NewChannelWebhookMessenger(client *Client, channelID snowflake.ID) (discord.WebhookMessenger, error) {
	webhooks, err := client.Rest.GetWebhooks(channelID)
	if err != nil {
		return nil, err
	}

	for _, webhook := range webhooks {
		if webhook.Type() == discord.WebhookTypeIncoming && webhook.(discord.IncomingWebhook).User.ID == client.ApplicationID {
			return channelWebhookMessenger{channelID: channelID, webhook: webhook.(discord.IncomingWebhook), client: client}, nil
		}
	}

	webhook, err := client.Rest.CreateWebhook(channelID, discord.WebhookCreate{
		Name:   WebhookDefaultName,
		Avatar: WebhookDefaultAvatar,
	})
	if err != nil {
		return nil, err
	}
	return channelWebhookMessenger{channelID: channelID, webhook: *webhook, client: client}, nil
}

type channelWebhookMessenger struct {
	channelID snowflake.ID
	webhook   discord.IncomingWebhook
	client    *Client
}

func (c channelWebhookMessenger) SendWebhook(message discord.MessageBuilder, username, avatarURL, threadName string) (*discord.Message, error) {
	create := message.BuildWebhookCreate(username, avatarURL, threadName)
	return c.client.Rest.CreateWebhookMessage(
		c.webhook.ID(),
		c.webhook.Token,
		create,
		rest.CreateWebhookMessageParams{
			Wait:           true,
			WithComponents: create.Flags.Has(discord.MessageFlagIsComponentsV2),
		},
	)
}

func (c channelWebhookMessenger) Webhook() discord.Webhook {
	return c.webhook
}

func (c channelWebhookMessenger) Send(message discord.MessageBuilder) (*discord.Message, error) {
	return c.SendWebhook(message, "", "", "")
}

func (c channelWebhookMessenger) Update(target snowflake.ID, message discord.MessageBuilder) (*discord.Message, error) {
	update := message.BuildWebhookUpdate()
	return c.client.Rest.UpdateWebhookMessage(
		c.webhook.ID(),
		c.webhook.Token,
		target,
		update,
		rest.UpdateWebhookMessageParams{
			WithComponents: update.Flags.Has(discord.MessageFlagIsComponentsV2),
		},
	)
}

func (c channelWebhookMessenger) Delete(message snowflake.ID) error {
	return c.client.Rest.DeleteWebhookMessage(c.webhook.ID(), c.webhook.Token, message, 0)
}

func NewThreadWebhookMessenger(client *Client, parentID, threadID snowflake.ID) (discord.WebhookMessenger, error) {
	webhooks, err := client.Rest.GetWebhooks(parentID)
	if err != nil {
		return nil, err
	}

	for _, webhook := range webhooks {
		if webhook.Type() == discord.WebhookTypeIncoming && webhook.(discord.IncomingWebhook).User.ID == client.ApplicationID {
			return threadWebhookMessenger{channelID: parentID, threadID: threadID, webhook: webhook.(discord.IncomingWebhook), client: client}, nil
		}
	}

	webhook, err := client.Rest.CreateWebhook(parentID, discord.WebhookCreate{
		Name:   WebhookDefaultName,
		Avatar: WebhookDefaultAvatar,
	})
	if err != nil {
		return nil, err
	}
	return threadWebhookMessenger{channelID: parentID, threadID: threadID, webhook: *webhook, client: client}, nil
}

type threadWebhookMessenger struct {
	channelID snowflake.ID
	threadID  snowflake.ID
	webhook   discord.IncomingWebhook
	client    *Client
}

func (t threadWebhookMessenger) SendWebhook(message discord.MessageBuilder, username, avatarURL, threadName string) (*discord.Message, error) {
	create := message.BuildWebhookCreate(username, avatarURL, threadName)
	return t.client.Rest.CreateWebhookMessage(
		t.webhook.ID(),
		t.webhook.Token,
		create,
		rest.CreateWebhookMessageParams{
			Wait:           true,
			ThreadID:       t.threadID,
			WithComponents: create.Flags.Has(discord.MessageFlagIsComponentsV2),
		})
}

func (t threadWebhookMessenger) Webhook() discord.Webhook {
	return t.webhook
}

func (t threadWebhookMessenger) Send(message discord.MessageBuilder) (*discord.Message, error) {
	return t.SendWebhook(message, "", "", "")
}

func (t threadWebhookMessenger) Update(target snowflake.ID, message discord.MessageBuilder) (*discord.Message, error) {
	update := message.BuildWebhookUpdate()
	return t.client.Rest.UpdateWebhookMessage(
		t.webhook.ID(),
		t.webhook.Token,
		target,
		update,
		rest.UpdateWebhookMessageParams{
			ThreadID:       t.threadID,
			WithComponents: update.Flags.Has(discord.MessageFlagIsComponentsV2),
		})
}

func (t threadWebhookMessenger) Delete(message snowflake.ID) error {
	return t.client.Rest.DeleteWebhookMessage(t.webhook.ID(), t.webhook.Token, message, t.threadID)
}
