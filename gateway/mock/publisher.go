package mock

import (
	"context"
	"errors"
)

type publisher struct {
	Channels map[string][]interface{}
}

func NewPublisher() publisher {
	return publisher{
		Channels: make(map[string][]interface{}),
	}
}

func (p *publisher) Publish(ctx context.Context, channel string, message interface{}) error {
	if channel == "" {
		return errors.New("no channel specified")
	}

	if _, ok := p.Channels[channel]; !ok {
		p.Channels[channel] = make([]interface{}, 0)
	}

	p.Channels[channel] = append(p.Channels[channel], message)
	return nil
}

func (p *publisher) GetMessages(channel string) ([]interface{}, error) {
	if channel == "" {
		return nil, errors.New("no channel specified")
	}

	messages, ok := p.Channels[channel]
	if !ok {
		return make([]interface{}, 0), nil
	}

	return messages, nil
}
