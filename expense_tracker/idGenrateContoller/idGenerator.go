package idgenerator

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

const (
	alphabet36        = "abcdefghijklmnopqrstuvwxyz0123456789"
	base36            = 36
	idWidth           = 6
	maxIDs            = 10_000_000 // cap if you want to stop at 10M
	a          uint64 = 13007
	b          uint64 = 987653
	M          uint64 = 2176782336
)

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
}

func toBase36Fixed(n uint64) string {
	var buf [idWidth]byte
	for i := idWidth - 1; i >= 0; i-- {
		buf[i] = alphabet36[n%base36]
		n /= base36
	}
	return string(buf[:])
}

func GenerateUserID() (string, error) {
	seq, err := rdb.Incr(ctx, "user_seq").Result() // 1,2,3,...
	if err != nil {
		return "", err
	}
	idx := uint64(seq - 1) // 0-based
	if idx >= 10_000_000 {
		return "", fmt.Errorf("id space exhausted")
	}
	n := (a*idx + b) % M // scramble for non-sequential look
	return toBase36Fixed(n), nil
}
