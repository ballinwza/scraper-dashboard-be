package helper

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func GenerateUniqeId() string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	return id
}
