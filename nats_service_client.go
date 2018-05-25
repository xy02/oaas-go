package oaas

import (
	"time"

	nats "github.com/nats-io/go-nats"
)

type NatsServiceClient struct {
	nc               *nats.Conn
	handshakeTimeout time.Duration
}

func (sclient NatsServiceClient) Call(serviceName ServiceName) (Caller, error) {
	nc := sclient.nc
	clientPort := "client." + RandomID()
	// Simple Sync Subscriber
	subIn, err := nc.SubscribeSync(clientPort)
	if err != nil {
		return nil, err
	}
	msg, err := nc.Request(serviceName, []byte(clientPort), sclient.handshakeTimeout)
	if err != nil {
		return nil, err
	}
	servicePort := string(msg.Data)
	return NatsCaller{
		NatsWatcher{
			NatsReceiver{
				subIn: subIn,
			},
		},
		NatsSender{
			publish: nc.Publish,
			portOut: servicePort,
		},
	}, nil
}

func (sclient NatsServiceClient) Watch(ServiceName) Watcher {
	return nil
}

type NatsServiceContext struct {
	NatsReceiver
	NatsResponser
	NatsServiceClient
}

type NatsCaller struct {
	NatsWatcher
	NatsSender
}
