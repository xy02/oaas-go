package oaas

import (
	"io"
)

type NatsCaller struct {
	NatsSubscriber
}

func (this NatsCaller) Send(bin RawData) error {
	if !this.subIn.IsValid() {
		return io.EOF
	}
	return this.natsPublisher.onNext(bin)
}
