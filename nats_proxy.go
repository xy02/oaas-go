package oaas

import (
	"time"

	nats "github.com/nats-io/go-nats"
)

type NatsProxy struct {
	NatsServiceClient
}

func (proxy NatsProxy) Register(serviceName ServiceName, service Service) error {
	nc := proxy.nc
	// 注册服务
	nc.Subscribe(serviceName, func(m *nats.Msg) {
		clientPort := string(m.Data)
		//注册服务端的接收地址
		servicePort := "service." + RandomID()
		//注册接收请求数据的通道
		subIn, err := nc.SubscribeSync(servicePort)
		if err != nil {
			return
		}
		defer subIn.Unsubscribe()
		//RPC返回服务接收地址
		err = nc.Publish(m.Reply, []byte(servicePort))
		if err != nil {
			return
		}
		//创建上下文
		ctx := NatsServiceContext{
			NatsReceiver{
				subIn: subIn,
			},
			NatsResponser{
				NatsSender: NatsSender{
					portOut: clientPort,
					publish: nc.Publish,
				},
				broadcastPort: serviceName + ".bc",
			},
			NatsServiceClient{
				nc:               nc,
				handshakeTimeout: proxy.handshakeTimeout,
			},
		}
		//开始服务
		service(ctx)
	})
	return nil
}

type NatsProxyOptions struct {
	ServerPort string
}

func NewNatsProxy(options NatsProxyOptions) (OaaSProxy, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}
	return NatsProxy{
		NatsServiceClient{
			nc:               nc,
			handshakeTimeout: 3 * time.Second,
		},
	}, nil
}
