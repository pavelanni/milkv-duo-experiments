package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	//rand.Seed(time.Now().UTC().UnixNano())

	const points = 1_000_000
	fmt.Printf("Estimating pi with %d point(s).\n\n", points)

	var sucess int
	start := time.Now()
	for i := 0; i < points; i++ {
		x, y := genRandomPoint(r)

		// Check if point lies within the circular region:
		if x*x+y*y < 1 {
			sucess++
		}
	}

	piApprox := 4.0 * (float64(sucess) / float64(points))
	errorPct := 100.0 * math.Abs(piApprox-math.Pi) / math.Pi
	duration := time.Since(start)

	time.Sleep(2 * time.Second)
	println("Estimated pi: ", piApprox)
	println("pi: ", math.Pi)
	println("Error: ", errorPct)
	println("Duration: ", duration)
}

// generates a random point p = (x, y)
func genRandomPoint(r *rand.Rand) (x, y float64) {
	x = 2.0*r.Float64() - 1.0
	y = 2.0*r.Float64() - 1.0
	return x, y
}
