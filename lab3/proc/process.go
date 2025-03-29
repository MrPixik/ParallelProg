package main

import (
	"encoding/binary"
	"os"
	"strconv"
	"sync"
)

func f(x float64) float64 {
	return 4.0 / (1.0 + x*x)
}

func middleRiemannSum(x1, x2 float64) float64 {
	return f((x1 + x2) / 2)
}

type semaphore struct {
	buffer chan struct{}
}

func NewSemaphore(capacity int) semaphore {
	return semaphore{buffer: make(chan struct{}, capacity)}
}

func (s *semaphore) Increment() {
	s.buffer <- struct{}{}
}
func (s *semaphore) Decrement() {
	<-s.buffer
}

func main() {
	threadNum, _ := strconv.Atoi(os.Args[1])
	a, _ := strconv.ParseFloat(os.Args[2], 64)
	b, _ := strconv.ParseFloat(os.Args[3], 64)
	h, _ := strconv.ParseFloat(os.Args[4], 64)

	var wg sync.WaitGroup
	sem := NewSemaphore(1)
	var sum float64

	chunkSize := (b - a) / float64(threadNum)

	for i := 0; i < threadNum; i++ {
		ai := a + chunkSize*float64(i)
		bi := ai + chunkSize

		wg.Add(1)
		go func(a, b float64) {
			defer wg.Done()
			var localSum float64
			for i := a; i < b; i += h {
				localSum += middleRiemannSum(i, i+h)
			}
			sem.Increment()
			sum += localSum
			sem.Decrement()
		}(ai, bi)
	}

	wg.Wait()

	binary.Write(os.Stdout, binary.LittleEndian, sum)
	os.Stdout.Sync()
}
