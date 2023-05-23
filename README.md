# safe-go

[![CircleCI](https://img.shields.io/circleci/build/github/deliveroo/safe-go?token=25175417b032a53cc407a1b1616d5a87813b9f8f)](https://circleci.com/gh/deliveroo/safe-go/tree/master)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg)](http://godoc.deliveroo.net/github.com/deliveroo/safe-go)

safe-go provides helpers for writing concurrent code safely.

## Package `safe`

Package `safe` provides helpers for gracefully handling panics in background
goroutines.

See the [GoDocs](http://godoc.deliveroo.net/github.com/deliveroo/safe-go) for
more information.

## Package `concurrent`

Package `concurrent` provides tools for making concurrent, inter-dependent calls that should yield data.
A good example is a microservice that needs to fetch data from a number of other services to perform its business logic.
This package helps to parallelize these calls in a way that's readable and safe.

See the [GoDocs](http://godoc.deliveroo.net/github.com/deliveroo/safe-go/concurrent) for
more information.
