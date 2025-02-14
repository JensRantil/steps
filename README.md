# discrete-events
[![GoDoc](https://godoc.org/github.com/JensRantil/discrete-events?status.svg)](https://pkg.go.dev/github.com/JensRantil/discrete-events)

A simple discrete event simulator/scheduler in Go.

See [simulation_test.go](simulation_test.go) for examples. You can find the documentation [here](https://pkg.go.dev/github.com/JensRantil/discrete-events).

## Why write yet another discrete event simulator?

 * https://github.com/fschuetz04/simgo was too slow for my needs. I needed to run simulations in an inner loop.
 * https://github.com/agoussia/godes was too heavyweight and complex for my needs. I just needed a simple performant scheduler (without any goroutines).