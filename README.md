# Art [![GoDoc](https://godoc.org/github.com/purehyperbole/art?status.svg)](https://godoc.org/github.com/purehyperbole/art) [![Go Report Card](https://goreportcard.com/badge/github.com/purehyperbole/art)](https://goreportcard.com/report/github.com/purehyperbole/art) [![Build Status](https://travis-ci.org/purehyperbole/art.svg?branch=master)](https://travis-ci.org/purehyperbole/art)


A thread safe Adaptive Radix Tree implementation in go


# Installation

To start using art, you can run:

`$ go get github.com/purehyperbole/art`

# Usage

To create a new radix tree

```go
package main

import (
    "github.com/purehyperbole/art"
)

func main() {
    // create a new art tree
    r := art.New()
}
```

`Lookup` can be used to retrieve a stored value

```go
value := r.Lookup([]byte("myKey1234"))
```

`Insert` allows a value to be stored for a given key.

A successful insert will return true.

If the operation conflicts with an insert from another thread, it will return false.

```go
if r.Insert([]byte("key"), &Thing{12345}) {
    fmt.Println("success!")
} else {
    fmt.Println("insert failed")
}
```


# Features/Wishlist

- [x] Lock free Insert using CAS (compare & swap)
- [x] Lookup
- [x] Basic key iterator
- [ ] Delete

## Why?

This project was created to explore the performance tradeoffs of a more memory efficient radix tree with my other lock free implementation (github.com/purehyperbole/rad).

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/).

## Copyright and License

Code and documentation copyright since 2019 purehyperbole.

Code released under
[the MIT License](LICENSE).
