package oaas

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

func RandomID() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
