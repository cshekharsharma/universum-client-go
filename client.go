package universum

import (
	"context"
	"encoding/base64"
	"strconv"
	"sync"
	"time"
)

var ncmu sync.Mutex

type Client struct {
	id   string
	pool *connPool
	opts *Options
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	return true, nil
}

func (c *Client) Get(ctx context.Context, key string) (*CommandResult, error) {
	rawResult, err := sendCommand(ctx, c, commandGet, key)
	if err != nil {
		return nil, err
	}

	cmdResult, err := toCommandResult(rawResult)
	if err != nil {
		return nil, err
	}

	return cmdResult, nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl int64) (*CommandResult, error) {
	rawResult, err := sendCommand(ctx, c, commandGet, key, value, ttl)

	if err != nil {
		return nil, err
	}

	cmdResult, err := toCommandResult(rawResult)
	if err != nil {
		return nil, err
	}

	return cmdResult, nil
}

func NewClient(opts *Options) (*Client, error) {
	ncmu.Lock()
	defer ncmu.Unlock()

	opts.init()
	connPool, err := newConnPool(opts)

	if err != nil {
		return nil, err
	}

	currTime := time.Now().UnixNano()
	uniqueId := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(currTime))))

	return &Client{
		id:   uniqueId,
		opts: opts,
		pool: connPool,
	}, nil
}
