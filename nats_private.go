package oaas

import (
	"log"

	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/go-nats"
)

var finalCompelete []byte

func init() {
	out, err := proto.Marshal(&Data{
		Type: &Data_Final{
			Final: "",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	finalCompelete = out
}

type natsPublisher struct {
	nc      *nats.Conn
	portOut string
}

func (this natsPublisher) send(data Data) error {
	out, err := proto.Marshal(&data)
	if err != nil {
		return err
	}
	return this.nc.Publish(this.portOut, out)
}

func (this natsPublisher) onNext(bin RawData) error {
	return this.send(Data{
		Type: &Data_Raw{
			Raw: bin,
		},
	})
}

func (this natsPublisher) onError(err error) error {
	return this.send(Data{
		Type: &Data_Final{
			Final: err.Error(),
		},
	})
}

func (this natsPublisher) onComplete() error {
	return this.nc.Publish(this.portOut, finalCompelete)
}
