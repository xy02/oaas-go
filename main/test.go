package main

import "github.com/xy02/oaas-go"

func main() {
	node := oaas.NewNatsProxy(oaas.NatsProxyOptions{
		ServerAddress: "nats address",
	})
	node.Register("aips.taxi.data", Hello)
	// caller := node.Call("aips.taxi.data")
	// watcher := node.Watch("aips.taxi.data")
}

func Hello(ctx oaas.ServiceContext) {

}
