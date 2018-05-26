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
	for i := 0; i < 5; i++ {
		go work(i+1, proxy)
	}
	forever := make(chan bool)
	<-forever
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
			if err != nil {
				errCh <- err
				return
			}
			log.Println(id, "onData", string(data))
		}
	}()
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			data := oaas.RandomID()
			err := caller.Send([]byte(data))
			if err != nil {
				// log.Fatal("failedSend", id, err)
				return
			}
			log.Println(id, "sent", data)
		case err := <-errCh:
			log.Println(id, "onError", err)
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
