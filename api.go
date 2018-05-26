package oaas

type RawData []byte

type ServiceName = string

type Subejct = string

type Service func(ServiceContext) error

type (
	OaaSProxy interface {
		Register(ServiceName, Service) error
		ServiceClient
	}

	ServiceContext interface {
		Receive() (RawData, error)
		Send(RawData) error
		ServiceClient
		// Config()
	}

	ServiceClient interface {
		Call(ServiceName) (Caller, error)
		Subscribe(Subejct) (Subscriber, error)
		Publish(Subejct, RawData) error
	}

	Caller interface {
		Send(RawData) error
		Subscriber
	}

	Subscriber interface {
		Receive() (RawData, error)
		Unsubscribe() error
	}
)
