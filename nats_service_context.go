package oaas

import (
	"io"

	nats "github.com/nats-io/go-nats"
)

type NatsServiceContext struct {
	nc         *nats.Conn
	subIn      *nats.Subscription
	chIn       chan *nats.Msg
	addressOut string
}

func NewNatsSerivceContext(proxy NatsProxy, addressIn string, addressOut string) (ServiceContext, error) {
	nc := proxy.nc
	//注册接收请求数据的通道
	chIn := make(chan *nats.Msg, 64)
	subIn, err := nc.ChanSubscribe(addressIn, chIn)
	if err != nil {
		return nil, err
	}
	return NatsServiceContext{
		nc:         nc,
		subIn:      subIn,
		chIn:       chIn,
		addressOut: addressOut,
	}, nil
}

func (nsc NatsServiceContext) Receive(Data) error {
	msg, ok := <-nsc.chIn
	if !ok {
		return io.EOF
	}
	//序列化
}
func (nsc NatsServiceContext) Send(Data) error {
	// Simple Publisher
	return nsc.nc.Publish(nsc.addressOut, []byte("Hello World"))
}
func (nsc NatsServiceContext) Error(error) {

}
func (nsc NatsServiceContext) Complete() {

}
func (nsc NatsServiceContext) Call(ServiceName) Caller {

}
func (nsc NatsServiceContext) Watch(ServiceName) Watcher {

}
