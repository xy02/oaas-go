package oaas

import (
	"errors"
	"io"
	"time"

	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/go-nats"
)

type NatsSubscriber struct {
	subIn *nats.Subscription
	natsPublisher
	receiveTimeout time.Duration
}

func (this NatsSubscriber) Receive() (RawData, error) {
	subIn := this.subIn
	if this.receiveTimeout == 0 {
		this.receiveTimeout = time.Minute
	}
	// log.Println("receiveTimeout", this.receiveTimeout)
	msg, err := subIn.NextMsg(this.receiveTimeout)
	if err != nil {
		return nil, err
	}
	//序列化
	data := &Data{}
	err = proto.Unmarshal(msg.Data, data)
	if err != nil {
		return nil, err
	}
	switch v := data.Type.(type) {
	case *Data_Raw:
		return v.Raw, nil
	case *Data_Final:
		subIn.Unsubscribe()
		if v.Final == "" {
			return nil, io.EOF
		}
		return nil, errors.New(v.Final)
	}
	subIn.Unsubscribe()
	return nil, io.ErrUnexpectedEOF
}

func (this NatsSubscriber) Unsubscribe() error {
	err := this.subIn.Unsubscribe()
	if err != nil {
		return err
	}
	if this.portOut != "" {
		return this.natsPublisher.onComplete()
	}
	return nil
}
