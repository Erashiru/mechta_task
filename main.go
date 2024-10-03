package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
)

type Object struct {
	A int `json:"a"`
	B int `json:"b"`
}

func calculateSum(object []Object, resultChan chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0
	for _, obj := range object {
		sum += obj.A + obj.B
	}
	resultChan <- sum
}

func main() {
	if len(os.Args[1:]) < 1 {
		log.Fatalf("No Args for goroutines number")
	}

	numOfGoroutines, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to parse number of goroutines: %v", err)
	}

	file, err := os.Open("data.json")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var objects []Object
	err = json.Unmarshal(byteValue, &objects)
	if err != nil {
		log.Fatalf("Failed to unmarshal data: %v", err)
	}

	resultChan := make(chan int, numOfGoroutines)
	var wg sync.WaitGroup

	objectsPerGoroutine := len(objects) / numOfGoroutines
	for i := 0; i < numOfGoroutines; i++ {
		start := i * objectsPerGoroutine
		end := start + objectsPerGoroutine

		if i == numOfGoroutines-1 {
			end = len(objects)
		}

		wg.Add(1)
		go calculateSum(objects[start:end], resultChan, &wg)
	}

	wg.Wait()
	close(resultChan)

	totalSum := 0
	for result := range resultChan {
		totalSum += result
	}

	fmt.Printf("Total sum: %d\n", totalSum)
}
