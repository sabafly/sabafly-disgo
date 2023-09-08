package discord

import (
	"fmt"
	"io"

	"github.com/disgoorg/snowflake/v2"
	"github.com/sabafly/sabafly-disgo/internal/builtin"
)

type MessageBuilder interface {
	AddActionRow(components ...InteractiveComponent) MessageBuilder
	AddContainerComponents(containers ...ContainerComponent) MessageBuilder
	AddEmbeds(embeds ...Embed) MessageBuilder
	AddFile(name string, description string, reader io.Reader, flags ...FileFlags) MessageBuilder
	AddFiles(files ...*File) MessageBuilder
	AddFlags(flags ...MessageFlags) MessageBuilder
	AddStickers(stickerIds ...snowflake.ID) MessageBuilder
	Create() MessageCreate
	Update() MessageUpdate
	ClearAllowedMentions() MessageBuilder
	ClearContainerComponents() MessageBuilder
	ClearEmbeds() MessageBuilder
	ClearFiles() MessageBuilder
	ClearFlags() MessageBuilder
	ClearStickers() MessageBuilder
	RemoveContainerComponent(i int) MessageBuilder
	RemoveEmbed(i int) MessageBuilder
	RemoveFile(i int) MessageBuilder
	RemoveFlags(flags ...MessageFlags) MessageBuilder
	SetAllowedMentions(allowedMentions *AllowedMentions) MessageBuilder
	SetContainerComponent(i int, container ContainerComponent) MessageBuilder
	SetContainerComponents(containerComponents ...ContainerComponent) MessageBuilder
	SetContent(content string) MessageBuilder
	SetContentf(content string, a ...any) MessageBuilder
	SetEmbed(i int, embed Embed) MessageBuilder
	SetEmbeds(embeds ...Embed) MessageBuilder
	SetEphemeral(ephemeral bool) MessageBuilder
	SetFile(i int, file *File) MessageBuilder
	SetFiles(files ...*File) MessageBuilder
	SetFlags(flags MessageFlags) MessageBuilder
	SetMessageReference(messageReference *MessageReference) MessageBuilder
	SetMessageReferenceByID(messageID snowflake.ID) MessageBuilder
	SetStickers(stickerIds ...snowflake.ID) MessageBuilder
	SetSuppressEmbeds(suppressEmbeds bool) MessageBuilder
}

func NewMessageBuilder() MessageBuilder {
	return &messageBuilderImpl{}
}

func NewMessageBuilderFromMessageCreate(m MessageCreate) MessageBuilder {
	return &messageBuilderImpl{
		MessageCreate: m,
	}
}

var _ MessageBuilder = (*messageBuilderImpl)(nil)

type messageBuilderImpl struct {
	MessageCreate
	attachments *[]AttachmentUpdate
}

// SetContent sets the content of the Message
func (b *messageBuilderImpl) SetContent(content string) MessageBuilder {
	b.Content = content
	return b
}

// SetContentf sets the content of the Message but with format
func (b *messageBuilderImpl) SetContentf(content string, a ...any) MessageBuilder {
	return b.SetContent(fmt.Sprintf(content, a...))
}

// SetTTS sets whether the Message should be text to speech
func (b *messageBuilderImpl) SetTTS(tts bool) MessageBuilder {
	b.TTS = tts
	return b
}

// SetEmbeds sets the Embed(s) of the Message
func (b *messageBuilderImpl) SetEmbeds(embeds ...Embed) MessageBuilder {
	b.Embeds = embeds
	return b
}

// SetEmbed sets the provided Embed at the index of the Message
func (b *messageBuilderImpl) SetEmbed(i int, embed Embed) MessageBuilder {
	if len(b.Embeds) > i {
		b.Embeds[i] = embed
	}
	return b
}

// AddEmbeds adds multiple embeds to the Message
func (b *messageBuilderImpl) AddEmbeds(embeds ...Embed) MessageBuilder {
	b.Embeds = append(b.Embeds, embeds...)
	return b
}

// ClearEmbeds removes all the embeds from the Message
func (b *messageBuilderImpl) ClearEmbeds() MessageBuilder {
	b.Embeds = []Embed{}
	return b
}

// RemoveEmbed removes an embed from the Message
func (b *messageBuilderImpl) RemoveEmbed(i int) MessageBuilder {
	if len(b.Embeds) > i {
		b.Embeds = append(b.Embeds[:i], b.Embeds[i+1:]...)
	}
	return b
}

// SetContainerComponents sets the discord.ContainerComponent(s) of the Message
func (b *messageBuilderImpl) SetContainerComponents(containerComponents ...ContainerComponent) MessageBuilder {
	b.Components = containerComponents
	return b
}

// SetContainerComponent sets the provided discord.InteractiveComponent at the index of discord.InteractiveComponent(s)
func (b *messageBuilderImpl) SetContainerComponent(i int, container ContainerComponent) MessageBuilder {
	if len(b.Components) > i {
		b.Components[i] = container
	}
	return b
}

// AddActionRow adds a new discord.ActionRowComponent with the provided discord.InteractiveComponent(s) to the Message
func (b *messageBuilderImpl) AddActionRow(components ...InteractiveComponent) MessageBuilder {
	b.Components = append(b.Components, ActionRowComponent(components))
	return b
}

// AddContainerComponents adds the discord.ContainerComponent(s) to the Message
func (b *messageBuilderImpl) AddContainerComponents(containers ...ContainerComponent) MessageBuilder {
	b.Components = append(b.Components, containers...)
	return b
}

// RemoveContainerComponent removes a discord.ActionRowComponent from the Message
func (b *messageBuilderImpl) RemoveContainerComponent(i int) MessageBuilder {
	if len(b.Components) > i {
		b.Components = append(b.Components[:i], b.Components[i+1:]...)
	}
	return b
}

// ClearContainerComponents removes all the discord.ContainerComponent(s) of the Message
func (b *messageBuilderImpl) ClearContainerComponents() MessageBuilder {
	b.Components = []ContainerComponent{}
	return b
}

// AddStickers adds provided stickers to the Message
func (b *messageBuilderImpl) AddStickers(stickerIds ...snowflake.ID) MessageBuilder {
	b.StickerIDs = append(b.StickerIDs, stickerIds...)
	return b
}

// SetStickers sets the stickers of the Message
func (b *messageBuilderImpl) SetStickers(stickerIds ...snowflake.ID) MessageBuilder {
	b.StickerIDs = stickerIds
	return b
}

// ClearStickers removes all Sticker(s) from the Message
func (b *messageBuilderImpl) ClearStickers() MessageBuilder {
	b.StickerIDs = []snowflake.ID{}
	return b
}

// SetFiles sets the File(s) for this MessageCreate
func (b *messageBuilderImpl) SetFiles(files ...*File) MessageBuilder {
	b.Files = files
	return b
}

// SetFile sets the discord.File at the index for this discord.MessageCreate
func (b *messageBuilderImpl) SetFile(i int, file *File) MessageBuilder {
	if len(b.Files) > i {
		b.Files[i] = file
	}
	return b
}

// AddFiles adds the discord.File(s) to the discord.MessageCreate
func (b *messageBuilderImpl) AddFiles(files ...*File) MessageBuilder {
	b.Files = append(b.Files, files...)
	return b
}

// AddFile adds a discord.File to the discord.MessageCreate
func (b *messageBuilderImpl) AddFile(name string, description string, reader io.Reader, flags ...FileFlags) MessageBuilder {
	b.Files = append(b.Files, NewFile(name, description, reader, flags...))
	return b
}

// ClearFiles removes all discord.File(s) of this discord.MessageCreate
func (b *messageBuilderImpl) ClearFiles() MessageBuilder {
	b.Files = []*File{}
	return b
}

// RemoveFile removes the discord.File at this index
func (b *messageBuilderImpl) RemoveFile(i int) MessageBuilder {
	if len(b.Files) > i {
		b.Files = append(b.Files[:i], b.Files[i+1:]...)
	}
	return b
}

// SetAllowedMentions sets the AllowedMentions of the Message
func (b *messageBuilderImpl) SetAllowedMentions(allowedMentions *AllowedMentions) MessageBuilder {
	b.AllowedMentions = allowedMentions
	return b
}

// ClearAllowedMentions clears the discord.AllowedMentions of the Message
func (b *messageBuilderImpl) ClearAllowedMentions() MessageBuilder {
	return b.SetAllowedMentions(nil)
}

// SetMessageReference allows you to specify a MessageReference to reply to
func (b *messageBuilderImpl) SetMessageReference(messageReference *MessageReference) MessageBuilder {
	b.MessageReference = messageReference
	return b
}

// SetMessageReferenceByID allows you to specify a Message CommandID to reply to
func (b *messageBuilderImpl) SetMessageReferenceByID(messageID snowflake.ID) MessageBuilder {
	if b.MessageReference == nil {
		b.MessageReference = &MessageReference{}
	}
	b.MessageReference.MessageID = &messageID
	return b
}

// SetFlags sets the message flags of the Message
func (b *messageBuilderImpl) SetFlags(flags MessageFlags) MessageBuilder {
	b.Flags = flags
	return b
}

// AddFlags adds the MessageFlags of the Message
func (b *messageBuilderImpl) AddFlags(flags ...MessageFlags) MessageBuilder {
	b.Flags = b.Flags.Add(flags...)
	return b
}

// RemoveFlags removes the MessageFlags of the Message
func (b *messageBuilderImpl) RemoveFlags(flags ...MessageFlags) MessageBuilder {
	b.Flags = b.Flags.Remove(flags...)
	return b
}

// ClearFlags clears the discord.MessageFlags of the Message
func (b *messageBuilderImpl) ClearFlags() MessageBuilder {
	return b.SetFlags(MessageFlagsNone)
}

// SetEphemeral adds/removes discord.MessageFlagEphemeral to the Message flags
func (b *messageBuilderImpl) SetEphemeral(ephemeral bool) MessageBuilder {
	if ephemeral {
		b.Flags = b.Flags.Add(MessageFlagEphemeral)
	} else {
		b.Flags = b.Flags.Remove(MessageFlagEphemeral)
	}
	return b
}

// SetSuppressEmbeds adds/removes discord.MessageFlagSuppressEmbeds to the Message flags
func (b *messageBuilderImpl) SetSuppressEmbeds(suppressEmbeds bool) MessageBuilder {
	if suppressEmbeds {
		b.Flags = b.Flags.Add(MessageFlagSuppressEmbeds)
	} else {
		b.Flags = b.Flags.Remove(MessageFlagSuppressEmbeds)
	}
	return b
}

// RetainAttachments removes all Attachment(s) from this Message except the ones provided
func (b *messageBuilderImpl) RetainAttachments(attachments ...Attachment) MessageBuilder {
	if b.attachments == nil {
		b.attachments = new([]AttachmentUpdate)
	}
	for _, attachment := range attachments {
		*b.attachments = append(*b.attachments, AttachmentKeep{ID: attachment.ID})
	}
	return b
}

// RetainAttachmentsByID removes all Attachment(s) from this Message except the ones provided
func (b *messageBuilderImpl) RetainAttachmentsByID(attachmentIDs ...snowflake.ID) MessageBuilder {
	if b.attachments == nil {
		b.attachments = new([]AttachmentUpdate)
	}
	for _, attachmentID := range attachmentIDs {
		*b.attachments = append(*b.attachments, AttachmentKeep{ID: attachmentID})
	}
	return b
}

// Create builds the MessageBuilder to a MessageCreate struct
func (b *messageBuilderImpl) Create() MessageCreate {
	return b.MessageCreate
}

// Update builds the MessageBuilder to a MessageUpdate struct
func (b *messageBuilderImpl) Update() MessageUpdate {
	return MessageUpdate{
		Content:         builtin.Or(b.Content != "", &b.Content, nil),
		Embeds:          builtin.Or(len(b.Embeds) > 0, &b.Embeds, nil),
		Components:      builtin.Or(len(b.Components) > 0, &b.Components, nil),
		Attachments:     b.attachments,
		Files:           b.Files,
		AllowedMentions: b.AllowedMentions,
		Flags:           builtin.Or(b.Flags != MessageFlagsNone, &b.Flags, nil),
	}
}
