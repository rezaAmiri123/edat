package msg

import "github.com/rezaAmiri123/edat/core"

// WithReply starts a reply builder allowing custom headers to be injected
func WithReply(reply core.Reply) *ReplyBuilder {
	return &ReplyBuilder{
		reply:   reply,
		headers: make(Headers),
	}
}

// ReplyBuilder is used to build custom replies
type ReplyBuilder struct {
	reply   core.Reply
	headers Headers
}

// Reply replaces the reply to be wrapped
func (b *ReplyBuilder) Reply(reply core.Reply) *ReplyBuilder {
	b.reply = reply
	return b
}

// Headers adds headers to include with the reply
func (b *ReplyBuilder) Headers(headers Headers) *ReplyBuilder {
	for key, value := range headers {
		b.headers[key] = value
	}
	return b
}

// Success wraps the reply with the custom headers as a Success reply
func (b *ReplyBuilder) Success() Reply {
	if b.reply == nil {
		b.reply = Success{}
	}

	b.headers[MessageReplyOutcome] = ReplyOutcomeSuccess
	b.headers[MessageReplyName] = b.reply.ReplyName()

	return &replyMessage{
		reply:   b.reply,
		headers: b.headers,
	}
}

// Failure wraps the reply with the custom headers as a Failure reply
func (b *ReplyBuilder) Failure() Reply {
	if b.reply == nil {
		b.reply = Failure{}
	}

	b.headers[MessageReplyOutcome] = ReplyOutcomesFailure
	b.headers[MessageReplyName] = b.reply.ReplyName()

	return &replyMessage{
		reply:   b.reply,
		headers: b.headers,
	}
}
