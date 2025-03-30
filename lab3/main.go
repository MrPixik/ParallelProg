package main

import (
	"ParallelProg/lab3/async"
	"ParallelProg/lab3/syncr"
	"fmt"
)

func main() {
	async.Run()
	fmt.Println("")
	syncr.Run()
}
