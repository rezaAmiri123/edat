package outbox

import (
	"context"
	"sync"
	"time"

	"github.com/rezaAmiri123/edat/log"
	"github.com/rezaAmiri123/edat/msg"
	"github.com/rezaAmiri123/edat/retry"
	"golang.org/x/sync/errgroup"
)

// PollingProcessor implements MessageProcessor
type PollingProcessor struct {
	in               MessageStore
	out              msg.MessagePublisher
	messagePerPollig int
	pollingInterval  time.Duration
	purgeOlderThan   time.Duration
	purgeIterval     time.Duration
	retryer          retry.Retryer
	logger           edatlog.Logger
	stopping         chan struct{}
	close            sync.Once
}

var _ MessageProcessor = (*PollingProcessor)(nil)

// NewPollingProcessor constructs a new PollingProcessor
func NewPollingProcessor(in MessageStore, out msg.MessagePublisher, options ...PollingProcessorOption) *PollingProcessor {
	p := &PollingProcessor{
		in:               in,
		out:              out,
		messagePerPollig: DefaultMessagesPerPolling,
		pollingInterval:  DefaultPollingInterval,
		purgeOlderThan:   DefaultPurgeOlderThan,
		purgeIterval:     DefaultPurgeInterval,
		retryer:          DefaultRetryer,
		logger:           edatlog.DefaultLogger,
		stopping:         make(chan struct{}),
	}

	for _, option := range options {
		option(p)
	}

	p.logger.Trace("outbox.PollingProcessor constructed")

	return p
}

// Start implements MessageProcessor.Start
func (p *PollingProcessor) Start(ctx context.Context) error {
	cCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	group, gCtx := errgroup.WithContext(cCtx)

	group.Go(func() error {
		if <-p.stopping; true {
			cancel()
		}
		return nil
	})

	group.Go(func() error {
		return p.processMessages(gCtx)
	})

	group.Go(func() error {
		return p.purgePublished(gCtx)
	})

	p.logger.Trace("processor started")

	return group.Wait()
}

// Stop implements MessageProcessor.Stop
func (p *PollingProcessor) Stop(ctx context.Context) (err error) {
	p.close.Do(func() {
		close(p.stopping)

		done := make(chan struct{})
		go func() {
			// anything to wait for?
			close(done)
		}()

		select {
		case <-done:
			p.logger.Trace("done with internal cleanup")
		case <-ctx.Done():
			p.logger.Warn("timed out waiting for internal cleanup to complete")
		}
	})
	return
}

func (p *PollingProcessor) processMessages(ctx context.Context) error {
	pollingTimer := time.NewTimer(0)

	for {
		var err error
		var messages []Message

		err = p.retryer.Retry(ctx, func() error {
			messages, err = p.in.Fetch(ctx, p.messagePerPollig)
			return err
		})

		if err != nil {
			p.logger.Error("error fetching message", edatlog.Error(err))
			return err
		}

		if len(messages) > 0 {
			p.logger.Trace("processong message", edatlog.Int("MessageCount", len(messages)))
			ids := make([]string, 0, len(messages))
			for _, message := range messages {
				err := p.processMessage(ctx, message)
				if err != nil {
					return err
				}

				ids = append(ids, message.MessageID)
			}

			err = p.retryer.Retry(ctx, func() error {
				return p.in.MarkPublished(ctx, ids)
			})
			if err != nil {
				return err
			}

			continue
		}

		if !pollingTimer.Stop() {
			select {
			case <-pollingTimer.C:
			default:
			}
		}

		pollingTimer.Reset(p.pollingInterval)

		select {
		case <-ctx.Done():
			return nil
		case <-pollingTimer.C:
		}
	}
}

func (p *PollingProcessor) processMessage(ctx context.Context, message Message) error {
	var err error
	var outgoingMsg msg.Message

	logger := p.logger.Sub(
		edatlog.String("MessageID", message.MessageID),
		edatlog.String("DestinationChannel", message.Destination),
	)

	outgoingMsg, err = message.ToMessage()
	if err != nil {
		logger.Error("error with transforming stored message", edatlog.Error(err))
		// TODO this has potential to halt processing; systems need to be in place to fix or address
		return err
	}
	err = p.out.Publish(ctx, outgoingMsg)
	if err != nil {
		logger.Error("error publishing message", edatlog.Error(err))
		// TODO this has potential to halt processing; systems need to be in place to fix or address
		return err
	}

	return nil
}

func (p *PollingProcessor) purgePublished(ctx context.Context) error {
	purgeTimer := time.NewTimer(0)

	for {
		err := p.retryer.Retry(ctx, func() error {
			return p.in.PurgePublished(ctx, p.purgeOlderThan)
		})

		if err != nil {
			p.logger.Error("error purging pubished message", edatlog.Error(err))
			return err
		}

		if !purgeTimer.Stop() {
			select {
			case <-purgeTimer.C:
			default:
			}
		}

		purgeTimer.Reset(p.purgeIterval)

		select {
		case <-ctx.Done():
			return nil
		case <-purgeTimer.C:
		}
	}
}
