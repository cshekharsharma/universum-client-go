package universum

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func RedisTest() {
	c := redis.NewClient(&redis.Options{
		Network: "localhost:6379",
	})

	a := c.Get(context.TODO(), "key")
	a.Bool()
}
