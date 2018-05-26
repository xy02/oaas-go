package oaas

import (
	"time"

	nats "github.com/nats-io/go-nats"
)

type NatsServiceClient struct {
	nc               *nats.Conn
	handshakeTimeout time.Duration
	receiveTimeout   time.Duration
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
	// log.Println("servicePort", servicePort)
	receiveTimeout := sclient.receiveTimeout
	return NatsCaller{
		NatsSubscriber{
			subIn,
			natsPublisher{
				nc:      nc,
				portOut: servicePort,
			},
			receiveTimeout,
		},
	}, nil
}

func (sclient NatsServiceClient) Subscribe(subject Subejct) (Subscriber, error) {
	subIn, err := sclient.nc.SubscribeSync("bc." + subject)
	if err != nil {
		return nil, err
	}
	receiveTimeout := sclient.receiveTimeout
	return NatsSubscriber{
		subIn,
		natsPublisher{
			nc: sclient.nc,
		},
		receiveTimeout,
	}, nil
}

func (sclient NatsServiceClient) Publish(subject Subejct, bin RawData) error {
	portOut := "bc." + subject
	publisher := natsPublisher{
		nc:      sclient.nc,
		portOut: portOut,
	}
	return publisher.onNext(bin)
}
