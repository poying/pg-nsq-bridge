package bridge

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/lib/pq"
)

func (b *Bridge) Start(ctx context.Context) {
	var wg sync.WaitGroup
	for postgresChannel, nsqChannel := range b.channels {
		wg.Add(1)
		go func(postgresChannel, nsqChannel string) {
			b.listen(ctx, postgresChannel, nsqChannel)
			wg.Done()
		}(postgresChannel, nsqChannel)
	}
	wg.Wait()
}

func (b *Bridge) listen(ctx context.Context, postgresChannel, nsqChannel string) {
	minReconn := 10 * time.Second
	maxReconn := time.Minute
	postgresListener := pq.NewListener(b.PostgresAddress, minReconn, maxReconn, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			b.OnError(err)
		}
	})
	err := postgresListener.Listen(postgresChannel)
	if err != nil {
		b.OnError(err)
		return
	}
Loop:
	for {
		select {
		case c := <-postgresListener.Notify:
			if c != nil {
				go b.process(c, postgresChannel, nsqChannel)
			}
		case <-ctx.Done():
			err = postgresListener.Unlisten(postgresChannel)
			if err != nil {
				b.OnError(err)
			}
			err = postgresListener.Close()
			if err != nil {
				b.OnError(err)
			}
			break Loop
		}
	}
}

func (b *Bridge) process(c *pq.Notification, postgresChannel, nsqChannel string) {
	payload := Payload{
		Pid:             c.BePid,
		PostgresChannel: c.Channel,
		Data:            c.Extra,
	}
	b.OnProcess(postgresChannel, nsqChannel, payload)
	payloadBuf, err := json.Marshal(payload)
	if err != nil {
		b.OnError(err)
		return
	}
	err = b.nsqProducer.Publish(nsqChannel, payloadBuf)
	if err != nil {
		b.OnError(err)
	}
}
