package oaas

import (
	"fmt"

	nats "github.com/nats-io/go-nats"
)

type NatsProxy struct {
}

func (c NatsProxy) Register(serviceName ServiceName, service Service) error {
	nc, _ := nats.Connect(nats.DefaultURL)

	// Simple Publisher
	nc.Publish("foo", []byte("Hello World"))

	// Simple Async Subscriber
	nc.Subscribe(serviceName, func(m *nats.Msg) {
		// 期望得到客户的接收数据地址
		clientAddress := string(m.Data)
		fmt.Printf("Received a message: %s\n", clientAddress)
		//注册服务端的接收地址
		serverAddress := ""
		nc.Subscribe(serverAddress, func(m *nats.Msg) {

		})
		//RPC返回服务接收地址
		nc.Publish(m.Reply, []byte(serverAddress))
	})
	return nil
}

func (c NatsProxy) Call(ServiceName) Caller {
	return nil
}

func (c NatsProxy) Watch(ServiceName) Watcher {
	return nil
}

type NatsProxyOptions struct {
	ServerAddress string
}

func NewNatsProxy(options NatsProxyOptions) OaaSProxy {
	return NatsProxy{}
}
