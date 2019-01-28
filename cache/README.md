# ThinkGo-Cache

`ThinkGo-Cache` is a cache library for Golang，it currently supports redis, memory, and can customize the store adapter.

## Installation

```
go get github.com/thinkoner/thinkgo/cache
```

## Usage

```go
import (
	"github.com/thinkoner/thinkgo/cache"
	"time"
)

func main() {
	var foo string 
	
	// Create a cache with memory store
	c, _ := cache.NewCache(cache.NewMemoryStore("thinkgo"))
	
	// Set the value
    c.Put("foo", "thinkgo", 10 * time.Minute)
	
	// Get the string associated with the key "foo" from the cache
	c.Get("foo", &foo)
}
```

## License

This project is licensed under the `Apache 2.0 license`.