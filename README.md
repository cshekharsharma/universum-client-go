# Go Client for UniversumDB

A light weight go client for [https://github.com/cshekharsharma/universum](UniversumDB).


## Usage 

```go

options := &universum.Options{
	ConnPoolsize:    10,
	ConnWaitTimeout: 10 * time.Second,
	ConnMaxLifetime: 1 * time.Hour,
	HostAddr:        "localhost:11191",
	DialTimeout:     1 * time.Second,
	MaxRetries:      5,
	ReadTimeout:     10 * time.Second,
	WriteTimeout:    10 * time.Second,
  IsReadOnly:      false
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

```

