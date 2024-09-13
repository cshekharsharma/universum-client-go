# Go Client for UniversumDB

A light weight go client for [https://github.com/cshekharsharma/universum](https://github.com/cshekharsharma/universum)


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
| `SNAPSHOT`    | Trigger a snapshot of the current database state.     |
| `INFO`        | Retrieve server and database information.             |
| `HELP`        | Display available commands and their usage.           |


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

	"github.com/cshekharsharma/universum"
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

## Running Tests
```bash
make test
```

## Contributing

To contribute:
1. Fork the repository
2. Submit a pull request
3. Open issues for bugs or feature requests

## License

This project is licensed under the Apache 2.0

----
