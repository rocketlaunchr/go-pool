<p align="right">
  <a href="http://godoc.org/github.com/rocketlaunchr/go-pool"><img src="http://godoc.org/github.com/rocketlaunchr/go-pool?status.svg" /></a>
  <a href="https://goreportcard.com/report/github.com/rocketlaunchr/go-pool"><img src="https://goreportcard.com/badge/github.com/rocketlaunchr/go-pool" /></a>
  <a href="https://gocover.io/github.com/rocketlaunchr/go-pool"><img src="http://gocover.io/_badge/github.com/rocketlaunchr/go-pool" /></a>
</p>

<p align="center">
<img src="https://github.com/rocketlaunchr/go-pool/raw/master/assets/logo.png" alt="go-pool" />
</p>

# sync.Pool wrapper

This package is a thin wrapper over the Pool provided by the `sync` package.

## Extra Features

- Invalidate an item from Pool (so it never gets used again)
- Set a maximum number of items for Pool
- Returns the number if items in the pool (idle and in-use)


Full Documentation will be updated soon


## Example

```go
import "github.com/rocketlaunchr/go-pool"

pool := pool.New(5) // maximum of 5 items in pool
pool.SetFactory(func() interface{} {
	x := &X{}
	return x
})

item := pool.GetItem()
defer item.Close()

// Use item here or mark as invalid
item.MarkAsInvalid()
``