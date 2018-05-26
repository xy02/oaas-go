package main

import (
	"log"
	"time"

	"github.com/xy02/oaas-go"
)

func main() {
	proxy, err := oaas.NewNatsProxy(oaas.NatsProxyOptions{})
	if err != nil {
		log.Fatal(err)
	}
	proxy.Register("aips.taxi.data", Hello)
	// go work(1, proxy)
	work(2, proxy)
}

func work(id int, proxy oaas.OaaSProxy) {
	caller, err := proxy.Call("aips.taxi.data")
	if err != nil {
		log.Fatal(err)
	}
	errCh := make(chan error)
	go func() {
		for {
			data, err := caller.Receive()
			log.Println("onData", id, string(data))
			if err != nil {
				errCh <- err
				return
			}
		}
	}()
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case t := <-ticker.C:
			err := caller.Send([]byte(t.String()))
			if err != nil {
				log.Fatal("failedSend", id, err)
			}
			log.Println("sent", t.String())
		case err := <-errCh:
			log.Println("onError", id, err)
		}
	}
}

func Hello(ctx oaas.ServiceContext) error {
	for count := 0; count < 5; count++ {
		data, err := ctx.Receive()
		// log.Println("service onData", string(data), err)
		if err != nil {
			return err
		}
		ctx.Send(data)
		ctx.Send(data)
	}
	return nil
}
