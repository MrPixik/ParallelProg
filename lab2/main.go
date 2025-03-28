package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

const (
	processFileName = "lab2/proc/process.exe"
	a               = 0.0
	b               = 1.0
	n               = 1000000000
	h               = (b - a) / n
)

func runProc(outCh chan<- float64, wg *sync.WaitGroup, goroutinesNum, ai, bi, h string) {
	defer wg.Done()

	process := exec.Command(processFileName, goroutinesNum, ai, bi, h)
	var out bytes.Buffer
	process.Stdout = &out
	err := process.Run()
	if err != nil {
		fmt.Println("Run failed", err)
		return
	}

	var readSum float64
	err = binary.Read(&out, binary.LittleEndian, &readSum) // Должен быть тот же порядок, что и при записи!
	if err != nil {
		fmt.Println("binary.Read failed:", err)
		return
	}
	outCh <- readSum
}

func main() {
	processNum, _ := strconv.Atoi(os.Args[1])
	goroutinesNum := os.Args[2]

	resCh := make(chan float64)
	var wg sync.WaitGroup

	tStart := time.Now()
	for i := 0; i < processNum; i++ {
		ai := a + (b-a)/float64(processNum)*float64(i)
		bi := a + (b-a)/float64(processNum)*float64(i+1)

		wg.Add(1)
		go runProc(resCh, &wg, goroutinesNum,
			strconv.FormatFloat(ai, 'f', -1, 64),
			strconv.FormatFloat(bi, 'f', -1, 64),
			strconv.FormatFloat(h, 'f', -1, 64))
	}
	go func() {
		wg.Wait()
		close(resCh)
	}()

	var integral float64
	for sum := range resCh {
		integral += sum
	}

	integral *= h
	fmt.Printf("time=%f sum=%22.15e\n", time.Since(tStart).Seconds(), integral)
}
