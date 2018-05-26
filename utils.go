package oaas

import "github.com/xy02/oaas-go/objectid"

func RandomID() string {
	return objectid.New().Hex()
}
