package oaas

import (
	"io"
	"time"

	nats "github.com/nats-io/go-nats"
)

type NatsProxyOptions struct {
	ServerAddress string
}

func NewNatsProxy(options NatsProxyOptions) (OaaSProxy, error) {
	url := options.ServerAddress
	if url == "" {
		url = nats.DefaultURL
	}
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	receiveTimeout := 0 * time.Minute
	return NatsProxy{
		NatsServiceClient{
			nc:               nc,
			handshakeTimeout: 3 * time.Second,
			receiveTimeout:   receiveTimeout,
		},
		receiveTimeout,
	}, nil
}

type NatsProxy struct {
	NatsServiceClient
	receiveTimeout time.Duration
}

func (proxy NatsProxy) Register(serviceName ServiceName, service Service) error {
	// 注册服务
	proxy.nc.Subscribe(serviceName, func(m *nats.Msg) {
		//不能阻塞这里
		go reg(m, proxy, service)
	})
	return nil
}

func reg(m *nats.Msg, proxy NatsProxy, service Service) {
	nc := proxy.nc
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
	// log.Println("clientPort", clientPort, servicePort, m.Reply)
	if err != nil {
		return
	}
	//创建上下文
	receiveTimeout := proxy.receiveTimeout
	ctx := NatsServiceContext{
		NatsServiceClient{
			nc:               nc,
			handshakeTimeout: proxy.handshakeTimeout,
			receiveTimeout:   receiveTimeout,
		},
		NatsSubscriber{
			subIn,
			natsPublisher{
				nc:      nc,
				portOut: clientPort,
			},
			receiveTimeout,
		},
	}
	//开始服务, block
	err = service(ctx)
	if err == nil {
		ctx.onComplete()
	} else if err == io.EOF {
		//调用者取消
	} else {
		ctx.onError(err)
	}
}
