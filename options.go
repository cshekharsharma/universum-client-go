package universum

import (
	"crypto/tls"
	"time"
)

const MaxDialTimeout = 5 * time.Second
const MaxReadTimeout = 1 * time.Second
const MaxWriteTimeout = 3 * time.Second
const AllowedMaxRetries = 1 << 4
const MaxConnPoolsize = 1 << 16

type Options struct {
	HostAddr   string
	ClientName string

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	MaxRetries   int64
	RetryBackoff time.Duration

	ConnPoolsize    int64
	ConnMaxLifetime time.Duration
	IsReadonly      bool
	TLSConfig       *tls.Config
}

func (o *Options) init() {
	if o.HostAddr == "" {
		o.HostAddr = "localhost:11191"
	}

	if o.ClientName == "" {
		o.ClientName = "GoClient"
	}

	if o.DialTimeout > MaxDialTimeout {
		o.DialTimeout = MaxDialTimeout
	}

	if o.ReadTimeout > MaxReadTimeout {
		o.ReadTimeout = MaxReadTimeout
	}

	if o.WriteTimeout > MaxWriteTimeout {
		o.WriteTimeout = MaxWriteTimeout
	}

	if o.MaxRetries > AllowedMaxRetries {
		o.MaxRetries = AllowedMaxRetries
	}

	if o.ConnPoolsize > MaxConnPoolsize {
		o.ConnPoolsize = MaxConnPoolsize
	}
}
