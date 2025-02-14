# steps
[![GoDoc](https://godoc.org/github.com/JensRantil/steps?status.svg)](https://pkg.go.dev/github.com/JensRantil/steps)

"Steps" is a simple discrete event simulator in Go. It's useful for simulations of systems that are driven by events, such as queues, workflows, etc.

See [the documentation](https://pkg.go.dev/github.com/JensRantil/steps) for API and examples.

## Why write yet another discrete event simulator?

 * [simgo](https://github.com/fschuetz04/simgo) was too slow for my needs. I needed to run simulations in an inner loop.
 * [godes](https://github.com/agoussia/godes) was too heavyweight and complex for my needs. I just needed a simple performant scheduler (without any goroutines).
