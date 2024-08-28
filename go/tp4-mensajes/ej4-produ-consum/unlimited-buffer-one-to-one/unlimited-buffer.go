package main

import (
	"fmt"
	"time"
	"math/rand"
)

func producer(buffer chan int) {

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for {
		produced_number := rng.Intn(100)
		fmt.Println("El productor produjo el numero", produced_number)
		buffer <- produced_number
	}
	
}

func consumer(buffer chan int) {
	for {
		consumed_number := <- buffer
		fmt.Println("Se consumio el numero", consumed_number)
	}
}

func main() {

	buffer := make(chan int)
	
	go consumer(buffer)
	go producer(buffer)

	time.Sleep(2 * time.Second)
}