package parallel

import (
	"ParallelProg/lab4/static"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func mergeSortParallel(items []int, maxGoroutines int) []int {
	if len(items) <= 1 {
		return items
	}

	if maxGoroutines <= 1 {
		return mergeSort(items)
	}

	mid := len(items) / 2

	var left, right []int
	var wg sync.WaitGroup

	sem := make(chan struct{}, maxGoroutines)

	wg.Add(1)
	sem <- struct{}{}
	go func() {
		defer wg.Done()
		left = mergeSortParallel(items[:mid], maxGoroutines/2)
		<-sem
	}()

	wg.Add(1)
	sem <- struct{}{}
	go func() {
		defer wg.Done()
		right = mergeSortParallel(items[mid:], maxGoroutines/2)
		<-sem
	}()

	wg.Wait()
	return merge(left, right)
}
func mergeSort(items []int) []int {
	if len(items) < 2 {
		return items
	}
	first := mergeSort(items[:len(items)/2])
	second := mergeSort(items[len(items)/2:])
	return merge(first, second)
}
func merge(a []int, b []int) []int {
	final := []int{}
	i := 0
	j := 0
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			final = append(final, a[i])
			i++
		} else {
			final = append(final, b[j])
			j++
		}
	}
	for ; i < len(a); i++ {
		final = append(final, a[i])
	}
	for ; j < len(b); j++ {
		final = append(final, b[j])
	}
	return final
}

func RunAsync() {
	rand.New(rand.NewSource(static.RandSeed))

	unsorted := make([]int, static.SliceSize)

	for i := 0; i < static.SliceSize; i++ {
		unsorted[i] = rand.Int()
	}

	start := time.Now()
	sorted := mergeSortParallel(unsorted, static.ProcNum)

	fmt.Printf("Parallel sort (proc num %d)\n", static.ProcNum)
	fmt.Printf("Time to sort: %f sec, slice length: %d\n", time.Since(start).Seconds(), len(sorted))
}
