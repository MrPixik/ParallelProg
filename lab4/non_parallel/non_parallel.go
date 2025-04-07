package non_parallel

import (
	"ParallelProg/lab4/static"
	"fmt"
	"math/rand"
	"time"
)

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

func RunSync() {

	rand.New(rand.NewSource(static.RandSeed))

	unsorted := make([]int, static.SliceSize)

	for i := 0; i < static.SliceSize; i++ {
		unsorted[i] = rand.Int()
	}

	start := time.Now()
	sorted := mergeSort(unsorted)

	fmt.Println("Non-Parallel sort:")
	fmt.Printf("Time to sort: %f sec, slice length: %d\n", time.Since(start).Seconds(), len(sorted))
}
