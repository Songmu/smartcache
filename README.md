smartcache
=======

[![Test Status](https://github.com/Songmu/smartcache/workflows/test/badge.svg?branch=master)][actions]
[![Coverage Status](https://coveralls.io/repos/Songmu/smartcache/badge.svg?branch=master)][coveralls]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/smartcache?status.svg)][godoc]

[actions]: https://github.com/Songmu/smartcache/actions?workflow=test
[coveralls]: https://coveralls.io/r/Songmu/smartcache?branch=master
[license]: https://github.com/Songmu/smartcache/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/Songmu/smartcache

The smartcache realizes smart in memory cache generation to minimize process blocks by using soft expire limit

## Synopsis

```go
var (
    expire     = 5*time.Minute
    softExpire = 1*time.Minute
)
ca := smartcache.New(expire, softExpire, func(ctx context.Context) (interface{}, error) {
    val, err := genCache(ctx)
    return val, err
})

val, err := ca.Get(context.Background())
```

## Description

In in-memory cache generation, there are problems, for example, thundering herd problems, blocking processing when regenerating or etc., but smartcache avoids them by setting a soft expire limit.

## Installation

```console
% go get github.com/Songmu/smartcache
```

## Author

[Songmu](https://github.com/Songmu)
