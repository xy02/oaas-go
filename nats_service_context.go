package oaas

type NatsServiceContext struct {
	NatsServiceClient
	NatsSubscriber
}

func (this NatsServiceContext) Send(bin RawData) error {
	return this.natsPublisher.onNext(bin)
}
