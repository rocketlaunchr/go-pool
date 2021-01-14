<p align="right">
  <a href="http://godoc.org/github.com/rocketlaunchr/go-pool"><img src="http://godoc.org/github.com/rocketlaunchr/go-pool?status.svg" /></a>
  <a href="https://goreportcard.com/report/github.com/rocketlaunchr/go-pool"><img src="https://goreportcard.com/badge/github.com/rocketlaunchr/go-pool" /></a>
  <a href="https://gocover.io/github.com/rocketlaunchr/go-pool"><img src="http://gocover.io/_badge/github.com/rocketlaunchr/go-pool" /></a>
</p>

<p align="center">
<img src="https://github.com/rocketlaunchr/go-pool/raw/master/assets/logo.png" alt="go-pool" />
</p>

# A Better sync.Pool

This package is a thin wrapper over the `Pool` provided by the `sync` package.

## Extra Features

- Invalidate an item from Pool (so it never gets used again)
- Set a maximum number of items for Pool
- Returns the number of items in the pool (idle and in-use)


Full Documentation will be updated soon


## Example

```go
import "github.com/rocketlaunchr/go-pool"

pool := pool.New(5) // maximum of 5 items in pool
pool.SetFactory(func() interface{} {
	return &X{}
})

item := pool.Borrow()
defer item.Return()

// Use item here or mark as invalid
x := item.Item.(*X) // Use item here
item.MarkAsInvalid()
```

Other useful packages
------------

- [awesome-svelte](https://github.com/rocketlaunchr/awesome-svelte) - Resources for killing react
- [dataframe-go](https://github.com/rocketlaunchr/dataframe-go) - Statistics and data manipulation
- [dbq](https://github.com/rocketlaunchr/dbq) - Zero boilerplate database operations for Go
- [electron-alert](https://github.com/rocketlaunchr/electron-alert) - SweetAlert2 for Electron Applications
- [google-search](https://github.com/rocketlaunchr/google-search) - Scrape google search results
- [igo](https://github.com/rocketlaunchr/igo) - A Go transpiler with cool new syntax such as fordefer (defer for for-loops)
- [mysql-go](https://github.com/rocketlaunchr/mysql-go) - Properly cancel slow MySQL queries
- [react](https://github.com/rocketlaunchr/react) - Build front end applications using Go
- [remember-go](https://github.com/rocketlaunchr/remember-go) - Cache slow database queries
- [testing-go](https://github.com/rocketlaunchr/testing-go) - Testing framework for unit testing