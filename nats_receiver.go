package oaas

import (
	"encoding/json"

	nats "github.com/nats-io/go-nats"
)

type NatsReceiver struct {
	subIn *nats.Subscription
}

func (receiver NatsReceiver) Receive(data Data) error {
	msg, err := receiver.subIn.NextMsg(0)
	if err != nil {
		return err
	}
	//序列化
	return json.Unmarshal(msg.Data, data)
}

type NatsWatcher struct {
	NatsReceiver
}

func (watcher NatsWatcher) Destroy() error {
	return watcher.subIn.Unsubscribe()
}
