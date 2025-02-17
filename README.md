# steps
[![GoDoc](https://godoc.org/github.com/JensRantil/steps?status.svg)](https://pkg.go.dev/github.com/JensRantil/steps)

"Steps" is a simple discrete event simulator in Go. It's useful for simulations of systems that are driven by events, such as queues, workflows, etc.

See [the documentation](https://pkg.go.dev/github.com/JensRantil/steps) for API and examples.

## Example

```go
sim := NewSimulation()

sim.Schedule(Event{When: sim.Now.Add(10 * time.Second), Action: func(s *Simulation) {
  fmt.Println("Actor 1:", sim.Now)
}})
sim.Schedule(Event{When: sim.Now.Add(time.Second), Action: func(s *Simulation) {
  fmt.Println("Actor 2:", sim.Now)
}})

sim.RunUntilDone()

// Output:
// Actor 2: 0001-01-01 00:00:01 +0000 UTC
// Actor 1: 0001-01-01 00:00:10 +0000 UTC
```

See [here](https://pkg.go.dev/github.com/JensRantil/steps#pkg-examples) for more examples.

## Why write yet another discrete event simulator?

 * [simgo](https://github.com/fschuetz04/simgo) was too slow for my needs. I needed to run simulations in an inner loop.
 * [godes](https://github.com/agoussia/godes) was too heavyweight and complex for my needs. I just needed a simple performant scheduler (without any goroutines).

## Contact & support

This library was coded up by [Jens Rantil](https://jensrantil.github.io). Do not hesitate to contact [Sweet Potato Tech](https://jensrantil.github.io/pages/services/) for support.