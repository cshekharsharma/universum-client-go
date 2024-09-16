# Go Client for UniversumDB

 
A light weight go client for [UniversumDB](https://github.com/cshekharsharma/universum), an in-memory key-value store with Redis-like features. The client offers high-performance communication with the Universum database server using RESP3 serialisation protocol.

[![Go Report Card](https://goreportcard.com/badge/github.com/cshekharsharma/universum-client-go)](https://goreportcard.com/badge/github.com/cshekharsharma/universum-client-go)


## Features

- **RESP3 Protocol Support**: Built-in support for encoding and decoding commands and responses using the RESP3 protocol.
- **Connection Pooling**: Efficient management of multiple connections for high concurrency and load management.
- **Timeout Management**: Configurable read, write, and request execution timeouts.
- **Error Handling**: Graceful error handling and connection recovery strategies to ensure high availability.
- **Client Authentication**: (Coming Soon) Support for username/password authentication and TLS encryption.

## Supported Commands

| Command       | Description                                           |
|---------------|-------------------------------------------------------|
| `PING`        | Check if the server is responsive.                    |
| `EXISTS`      | Determine if a key exists in the database.            |
| `GET`         | Retrieve the value associated with a key.             |
| `SET`         | Set a value for a key with optional expiration time.  |
| `DELETE`      | Remove a key and its associated value from the database. |
| `INCR`        | Increment the integer value of a key.                 |
| `DECR`        | Decrement the integer value of a key.                 |
| `APPEND`      | Append a value to a string key.                       |
| `MGET`        | Retrieve values for multiple keys at once.            |
| `MSET`        | Set multiple key-value pairs at once.                 |
| `MDELETE`     | Delete multiple keys at once.                         |
| `TTL`         | Get the remaining time-to-live (TTL) of a key.        |
| `EXPIRE`      | Set a timeout on a key, after which it will be deleted. |
| `INFO`        | Retrieve server and database information.             |


## Installation

```
 go get github.com/cshekharsharma/universum-client-go@latest
```

## Usage 

```go

package main

import (
	"context"
	"time"

	"github.com/cshekharsharma/universum-client-go"
)

func main() {

	options := &universum.Options{
		ConnPoolsize:    10,
		ConnWaitTimeout: 10 * time.Second,
		ConnMaxLifetime: 1 * time.Hour,
		HostAddr:        "localhost:11191",
		DialTimeout:     1 * time.Second,
		MaxRetries:      5,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		IsReadOnly:      false,
	}

	client := universum.NewClient(options)
	result, err := client.Get(context.Background(), "key")

	if err != nil {
		// handle the error
	}

	// record found
	if result.Code == universum.RespRecordFound {
		value := result.Value
	}
}


```

## Configuration Options

The client can be configured via the Options struct. Here are some of the configurable fields:

| Setting         | Description                                           |
|-----------------|-------------------------------------------------------|
| HostAddr        | Address of the Universum DB server. (ip:port)         |
| DialTimeout     | Timeout duration (in seconds) for establishing connections. |
| MaxRetries      | Number of retry attempts for connecting to the server. |
| ConnPoolsize    | Number of connections in the connection pool. |
| ConnWaitTimeout | Duration (in seconds) to wait for an available connection from the pool. |
| ConnMaxLifetime | Maximum lifetime of a connection, after which it will be dropped. |
| ReadTimeout     | Timeout duration (in seconds) for reading from the network |
| WriteTimeout    | Timeout duration (in seconds) for writing to the network |
| IsReadOnly      | Mark connection as readonly (disallowing write commands) |


## Running Tests

```bash
make test
```

## Contributing

To contribute:
1. Fork the repository
2. Submit a pull request
3. Open issues for bugs or feature requests
4. 

## License

This project is licensed under the Apache 2.0

----
