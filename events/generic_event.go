package events

import (
	"context"

	"github.com/disgoorg/disgo/bot"
)

// NewGenericEvent constructs a new GenericEvent with the provided Client instance
func NewGenericEvent(client *bot.Client, sequenceNumber int, shardID int) *GenericEvent {
	ctx, cancelCause := context.WithCancelCause(context.Background())
	ctx, cancelFunc := context.WithCancel(ctx)
	return &GenericEvent{client: client, sequenceNumber: sequenceNumber, shardID: shardID, Context: ctx, cancelFunc: cancelFunc, cancelCauseFunc: cancelCause}
}

// GenericEvent the base event structure
type GenericEvent struct {
	client          *bot.Client
	sequenceNumber  int
	shardID         int
	cancelFunc      context.CancelFunc
	cancelCauseFunc context.CancelCauseFunc
	context.Context
}

// Client returns the bot.Client instance that dispatched the event
func (e *GenericEvent) Client() *bot.Client {
	return e.client
}

// SequenceNumber returns the sequence number of the gateway event
func (e *GenericEvent) SequenceNumber() int {
	return e.sequenceNumber
}

// ShardID returns the shard ID the event was dispatched from
func (e *GenericEvent) ShardID() int {
	return e.shardID
}

func (e *GenericEvent) IsCanceled() bool {
	return e.Context.Err() != nil
}

func (e *GenericEvent) Cancel() {
	if e.cancelFunc == nil {
		panic("cancelFunc is nil")
	}
	e.cancelFunc()
}

func (e *GenericEvent) CancelCause(err error) {
	if e.cancelCauseFunc == nil {
		panic("cancelFunc is nil")
	}
	e.cancelCauseFunc(err)
}
