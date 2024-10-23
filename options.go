package universum

import (
	"time"
)

const DefaultHostAddr = "localhost:11191"
const DefaultClientName = "GoClient"

const DefaultDialTimeout = 1 * time.Second
const MaxDialTimeout = 5 * time.Second

const DefaultReadTimeout = 1 * time.Second
const MaxReadTimeout = 3 * time.Second

const DefaultWriteTimeout = 1 * time.Second
const MaxWriteTimeout = 3 * time.Second

const DefaultMaxRetries = 1 << 2 // 4
const AllowedMaxRetries = 1 << 4 // 16

const DefaultRetryBackoff = 50 * time.Millisecond
const MaxRetryBackoff = 500 * time.Millisecond

const DefaultConnPoolsize = 1 << 4 // 16
const MaxConnPoolsize = 1 << 16    // 65636

const DefaultConnMaxLifetime = 10 * time.Minute
const MaxConnMaxLifetime = 30 * time.Minute

type Options struct {
	HostAddr   string
	ClientName string

	DialTimeout     time.Duration
	ConnWaitTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration

	MaxRetries   int64
	RetryBackoff time.Duration

	ConnPoolsize    int64
	ConnMaxLifetime time.Duration
	IsReadonly      bool

	EnableTLS          bool
	TLSCertFile        string
	TLSKeyFile         string
	CAFile             string
	InsecureSkipVerify bool
}

func (opts *Options) Init() {
	// HostAddr validation
	if opts.HostAddr == "" {
		opts.HostAddr = DefaultHostAddr
	}

	// ClientName validation
	if opts.ClientName == "" {
		opts.ClientName = DefaultClientName
	}

	// DialTimeout validation
	if opts.DialTimeout <= 0 {
		opts.DialTimeout = DefaultDialTimeout
	} else if opts.DialTimeout > MaxDialTimeout {
		opts.DialTimeout = MaxDialTimeout
	}

	// ReadTimeout validation
	if opts.ReadTimeout <= 0 {
		opts.ReadTimeout = DefaultReadTimeout
	} else if opts.ReadTimeout > MaxReadTimeout {
		opts.ReadTimeout = MaxReadTimeout
	}

	// WriteTimeout validation
	if opts.WriteTimeout <= 0 {
		opts.WriteTimeout = DefaultWriteTimeout
	} else if opts.WriteTimeout > MaxWriteTimeout {
		opts.WriteTimeout = MaxWriteTimeout
	}

	// MaxRetries validation
	if opts.MaxRetries <= 0 {
		opts.MaxRetries = DefaultMaxRetries
	} else if opts.MaxRetries > AllowedMaxRetries {
		opts.MaxRetries = AllowedMaxRetries
	}

	// RetryBackoff validation
	if opts.RetryBackoff <= 0 {
		opts.RetryBackoff = DefaultRetryBackoff
	} else if opts.RetryBackoff > MaxRetryBackoff {
		opts.RetryBackoff = MaxRetryBackoff
	}

	// ConnPoolsize validation
	if opts.ConnPoolsize <= 0 {
		opts.ConnPoolsize = DefaultConnPoolsize
	} else if opts.ConnPoolsize > MaxConnPoolsize {
		opts.ConnPoolsize = MaxConnPoolsize
	}

	// ConnMaxLifetime validation
	if opts.ConnMaxLifetime <= 0 {
		opts.ConnMaxLifetime = DefaultConnMaxLifetime
	} else if opts.ConnMaxLifetime > MaxConnMaxLifetime {
		opts.ConnMaxLifetime = MaxConnMaxLifetime
	}
}
