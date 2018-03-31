package bridge

import (
	"time"

	"github.com/lib/pq"
	nsq "github.com/nsqio/go-nsq"
)

type Options struct {
	PostgresAddress string
	NSQAddress      string
	PostgresChannel string
	NSQChannel      string
	OnError         func(err error)
	OnProcess       func(payload Payload)
}

type Bridge struct {
	Options
	stopCh           chan int
	postgresListener *pq.Listener
	nsqProducer      *nsq.Producer
}

func New(options Options) (*Bridge, error) {
	nsqConfig := nsq.NewConfig()
	nsqProducer, err := nsq.NewProducer(options.NSQAddress, nsqConfig)
	if err != nil {
		return nil, err
	}
	minReconn := 10 * time.Second
	maxReconn := time.Minute
	postgresListener := pq.NewListener(options.PostgresAddress, minReconn, maxReconn, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			options.OnError(err)
		}
	})
	return &Bridge{
		options,
		make(chan int, 1),
		postgresListener,
		nsqProducer,
	}, nil
}

func (b *Bridge) Stop() {
	b.stopCh <- 1
}
