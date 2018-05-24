package oaas

import (
	"fmt"

	nats "github.com/nats-io/go-nats"
)

type NatsProxy struct {
	nc *nats.Conn
}

func (proxy NatsProxy) Register(serviceName ServiceName, service Service) error {
	nc := proxy.nc
	// 注册服务
	nc.Subscribe(serviceName, func(m *nats.Msg) {
		// 期望得到客户的接收数据地址
		clientAddress := string(m.Data)
		fmt.Printf("Received a message: %s\n", clientAddress)
		//注册服务端的接收地址
		serverAddress := ""
		//创建上下文
		ctx, err := NewNatsSerivceContext(proxy, serverAddress, clientAddress)
		if err != nil {
			return
		}
		//RPC返回服务接收地址
		nc.Publish(m.Reply, []byte(serverAddress))
		//开始服务
		go service(ctx)
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

func NewNatsProxy(options NatsProxyOptions) (OaaSProxy, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}
	return NatsProxy{
		nc: nc,
	}, nil
}
