package oaas

type Data = interface{}

type ServiceName = string

type Service func(ServiceContext)

type ServiceContext interface {
	Receiver
	Sender
	ServiceClient
	// Config()
}

type ServiceClient interface {
	Call(ServiceName) Caller
	Watch(ServiceName) Watcher
}

type Caller interface {
	Watcher
	Send(Data) error
}

type Sender interface {
	Send(Data) error
	Error(error)
	Complete()
}

type Receiver interface {
	Receive(Data) error
}

type Watcher interface {
	Receiver
	Destroy()
}

type OaaSProxy interface {
	Register(ServiceName, Service) error
	ServiceClient
}
