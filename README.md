<p align="right">
  ⭐ &nbsp;&nbsp;<strong>the project to show your appreciation.</strong> :arrow_upper_right:
</p>

<p align="center">
<img src="https://github.com/rocketlaunchr/go-pool/raw/master/assets/logo.png" alt="go-pool" />
</p>

<p align="right">
  <a href="http://godoc.org/github.com/rocketlaunchr/go-pool"><img src="http://godoc.org/github.com/rocketlaunchr/go-pool?status.svg" /></a>
  <a href="https://goreportcard.com/report/github.com/rocketlaunchr/go-pool"><img src="https://goreportcard.com/badge/github.com/rocketlaunchr/go-pool" /></a>
  <a href="https://gocover.io/github.com/rocketlaunchr/go-pool"><img src="http://gocover.io/_badge/github.com/rocketlaunchr/go-pool" /></a>
</p>

# A Better sync.Pool  ![](https://img.shields.io/static/v1?label=%E2%9A%A0%EF%B8%8F&message=BETA&labelColor=fbcf11&color=red) 

This package is a **thin** wrapper over the `Pool` provided by the `sync` package. The `Pool` is an essential package to obtain maximum performance by reducing the number of memory allocations.

## Extra Features

- Invalidate an item from the Pool (so it never gets used again)
- Set a maximum number of items for the Pool
- Returns the number of items in the pool (idle and in-use)

## When should I use a pool?

If you frequently allocate many objects of the same type and you want to save some memory allocation and garbage allocation overhead — @jrv

[How did I improve latency by 700% using sync.Pool](https://www.akshaydeo.com/blog/2017/12/23/How-did-I-improve-latency-by-700-percent-using-syncPool)

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