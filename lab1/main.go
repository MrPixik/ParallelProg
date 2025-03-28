package main

import (
	"fmt"
	"math"
	"time"
)

const numOps = 1000000000

func f(x float64) float64 {
	return 4.0 / (1.0 + x*x)
}
func addition(x float64) float64 {
	return (4.0 + x) / (1.0 + x*x)
}
func subtraction(x float64) float64 {
	return (4.0 - x) / (1.0 + x*x)
}
func multiplication(x float64) float64 {
	return (4.0 * x) / (1.0 + x*x)
}
func division(x float64) float64 {
	return (4.0 / x) / (1.0 + x*x)
}
func pow(x float64) float64 {
	return math.Pow(x, 2.0) / (1.0 + x*x)
}
func exp(x float64) float64 {
	return math.Exp(x) / (1.0 + x*x)
}
func log(x float64) float64 {
	return math.Log(x) / (1.0 + x*x)
}
func sin(x float64) float64 {
	return math.Sin(x) / (1.0 + x*x)
}
func benchmark(ms float64, operation func(float64) float64, name string) {
	sum := 0.0
	start := time.Now()
	for i := 1; i <= numOps; i++ {
		sum += operation(float64(i))
	}
	elapsed := time.Since(start).Seconds()
	res := elapsed - ms
	gflops := float64(numOps) / (res * 1e9)
	fmt.Printf("Time: %.6f sec '%s' perf.: %.6e GFlops\n", res, name, gflops)
}

func measurmentStandart(operation func(float64) float64) float64 {
	sum := 0.0
	start := time.Now()
	for i := 1; i <= numOps; i++ {
		sum += operation(float64(i))
	}
	elapsed := time.Since(start).Seconds()
	return elapsed
}

func main() {

	ms := measurmentStandart(f)

	benchmark(ms, addition, "+")
	benchmark(ms, subtraction, "-")
	benchmark(ms, multiplication, "*")
	benchmark(ms, division, "/")
	benchmark(ms, pow, "pow")
	benchmark(ms, exp, "exp")
	benchmark(ms, log, "log")
	benchmark(ms, sin, "sin")
}
