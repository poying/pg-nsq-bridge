package bridge

import (
	"errors"
	"strings"

	nsq "github.com/nsqio/go-nsq"
)

type Options struct {
	PostgresAddress string
	NSQAddress      string
	Channel         string
	OnError         func(err error)
	OnProcess       func(postgresChannel, nsqChannel string, payload Payload)
}

type Bridge struct {
	Options
	channels    map[string]string
	nsqProducer *nsq.Producer
}

func New(options Options) (*Bridge, error) {
	nsqConfig := nsq.NewConfig()
	nsqProducer, err := nsq.NewProducer(options.NSQAddress, nsqConfig)
	if err != nil {
		return nil, err
	}
	channels, err := parseChannel(options.Channel)
	if err != nil {
		return nil, err
	}
	return &Bridge{
		options,
		channels,
		nsqProducer,
	}, nil
}

func parseChannel(input string) (map[string]string, error) {
	channels := make(map[string]string)
	parts := strings.Split(input, ",")
	for _, part := range parts {
		ch := strings.Split(part, ":")
		if len(ch) != 2 {
			return channels, errors.New("Invalid channel format")
		}
		channels[strings.TrimSpace(ch[0])] = strings.TrimSpace(ch[1])
	}
	return channels, nil
}
