package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
func sum(filePath string, barrier *sync.WaitGroup, mutex *sync.Mutex, totalSum *int64, sums *map[int][]string) {
	defer barrier.Done()

	data, _ := readFile(filePath)

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}

	mutex.Lock()
	*totalSum += int64(_sum)
	(*sums)[_sum] = append((*sums)[_sum], filePath)
	mutex.Unlock()

}

// print the totalSum for all files and the files with equal sum
func main() {
	now := time.Now().UnixNano()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	var barrier sync.WaitGroup
	var mutex sync.Mutex

	var totalSum int64
	sums := make(map[int][]string)
	for _, path := range os.Args[1:] {
		barrier.Add(1)
		go sum(path, &barrier, &mutex, &totalSum, &sums)
	}

	barrier.Wait()

	fmt.Println(totalSum)

	for sum, files := range sums {
		if len(files) > 1 {
			fmt.Printf("Sum %d: %v\n", sum, files)
		}
	}
	fmt.Printf("executado em: %d\n", time.Now().UnixNano()-now)
}
