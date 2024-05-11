package universum

import (
	"encoding/base64"
	"strconv"
	"time"
)

type Client struct {
	id   string
	pool *connPool
	opts *Options
}

func (c *Client) Exists(key string) (bool, error) {
	return true, nil
}

func (c *Client) Get(key string) (interface{}, error) {
	return 0, nil
}

func (c *Client) Set(key string, value interface{}, ttl int64) (bool, error) {
	return true, nil
}

func NewClient(opts *Options) *Client {
	opts.init()
	currTime := time.Now().UnixNano()
	return &Client{
		id:   base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(currTime)))),
		opts: opts,
		pool: newConnPool(),
	}
}
