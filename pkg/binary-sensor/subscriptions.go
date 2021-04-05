package binarysensor

import (
	"context"
	"fmt"

	"github.com/spritkopf/esb-bridge/pkg/esbbridge"
)

const (
	NotificationIndexChannel = 0
	NotificationIndexValue   = 1
)

type subscription struct {
	channel uint8
	rxChan  chan bool
	ctx     context.Context
	cancel  context.CancelFunc
}

// Subscribe subscribes to the selected channel of a binary sensor. Incoming sensor states will be sent to
// the channel rxChan in the returned subscription structure.
func (b *BinarySensor) Subscribe(channel uint8) (*subscription, error) {

	ctx, cancel := context.WithCancel(context.Background())
	s := subscription{
		channel: channel,
		rxChan:  make(chan bool, 1),
		ctx:     ctx,
		cancel:  cancel}

	esbMsgChan, err := b.esbClient.Listen(ctx, b.esbAddress, cmdStateNotification)

	if err != nil {
		return nil, fmt.Errorf("error starting listening: %v", err)

	}

	go func(s subscription, esbChan <-chan esbbridge.EsbMessage) {
		select {
		case <-s.ctx.Done():
			return
		case e := <-esbChan:
			if e.Payload[NotificationIndexChannel] == s.channel {
				if e.Payload[NotificationIndexValue] == 0 {
					s.rxChan <- false
				} else {
					s.rxChan <- true
				}
			}
		}
	}(s, esbMsgChan)

	return &s, nil
}

// Unubscribe stops the subscription to the selected channel of a binary sensor.
func (b *BinarySensor) Unubscribe(sub *subscription) error {

	sub.cancel()
	return nil
}
