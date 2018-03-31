package bridge

import (
	"encoding/json"

	"github.com/lib/pq"
)

func (b *Bridge) Start() {
	err := b.postgresListener.Listen(b.PostgresChannel)
	if err != nil {
		b.OnError(err)
		return
	}
	for {
		select {
		case c := <-b.postgresListener.Notify:
			go b.process(c)
		case <-b.stopCh:
			err = b.postgresListener.Unlisten(b.PostgresChannel)
			if err != nil {
				b.OnError(err)
			}
			err = b.postgresListener.Close()
			if err != nil {
				b.OnError(err)
			}
		}
	}
}

func (b *Bridge) process(c *pq.Notification) {
	payload := Payload{
		Pid:             c.BePid,
		PostgresChannel: c.Channel,
		Data:            c.Extra,
	}
	b.OnProcess(payload)
	payloadBuf, err := json.Marshal(payload)
	if err != nil {
		b.OnError(err)
		return
	}
	err = b.nsqProducer.Publish(b.NSQChannel, payloadBuf)
	if err != nil {
		b.OnError(err)
	}
}
