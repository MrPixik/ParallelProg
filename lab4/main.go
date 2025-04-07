package main

import (
	"ParallelProg/lab4/non_parallel"
	"ParallelProg/lab4/parallel"
)

func main() {
	non_parallel.RunSync()
	parallel.RunAsync()
}
