package oaas

import (
	"encoding/json"
)

type NatsSender struct {
	publish func(subj string, data []byte) error
	portOut string
}

func (sender NatsSender) Send(data Data) ([]byte, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	// Simple Publisher
	err = sender.publish(sender.portOut, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

type NatsResponser struct {
	NatsSender
	broadcastPort string
}

func (res NatsResponser) Send(data Data) ([]byte, error) {
	buf, err := res.NatsSender.Send(data)
	if err != nil {
		return nil, err
	}
	//broadcast
	err = res.publish(res.broadcastPort, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (res NatsResponser) Error(err error) error {
	buf := []byte(err.Error())
	err = res.publish(res.portOut, buf)
	if err != nil {
		return err
	}
	//broadcast
	return res.publish(res.broadcastPort, buf)
}

func (res NatsResponser) Complete() error {
	err := res.publish(res.portOut, nil)
	if err != nil {
		return err
	}
	//broadcast
	return res.publish(res.broadcastPort, nil)
}
