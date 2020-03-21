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

The smartcache is an in-memory cache library with avoiding the following problems.

- thundering herd
- block processing when regenerating
- etc.

To avoid the above problems, you can set a soft expire limit to Cache. The soft expired cached value is internally pre-warmed by a single goroutine and the value is replaced seamlessly.

## Installation

```console
% go get github.com/Songmu/smartcache
```

## Author

[Songmu](https://github.com/Songmu)
