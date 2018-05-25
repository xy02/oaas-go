package oaas

type Data = interface{}

type ServiceName = string

type Service func(ServiceContext)

type ServiceContext interface {
	Receiver
	Responser
	ServiceClient
	// Config()
}

type ServiceClient interface {
	Watch(ServiceName) Watcher
	Call(ServiceName) (Caller, error)
}

type Caller interface {
	Watcher
	Sender
}

type Sender interface {
	Send(Data) ([]byte, error)
}

type Responser interface {
	Sender
	Error(error) error
	Complete() error
}

type Receiver interface {
	Receive(Data) error
}

type Watcher interface {
	Receiver
	Destroy() error
}

type OaaSProxy interface {
	Register(ServiceName, Service) error
	ServiceClient
}
