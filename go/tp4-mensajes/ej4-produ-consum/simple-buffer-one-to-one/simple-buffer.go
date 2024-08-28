package main

import (
	"fmt"
	"time"
	"math/rand"
)

func producer(buffer *int, can_consume chan struct{}, can_produce chan struct{}) {

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for {
		<- can_produce
		*buffer =  rng.Intn(100)
		fmt.Println("El productor produjo el numero", *buffer)
		can_consume <- struct{}{} //almaceno un mensaje
	}
	
}

func consumer(buffer *int, can_consume chan struct{}, can_produce chan struct{}) {
	for {
		<- can_consume
		fmt.Println("Se consumio el numero", *buffer)
		can_produce <- struct{}{}
	}
}

func main() {
	
	buffer_simple := 0 //buffer de mensajes vacios

	can_consume := make(chan struct{},1)
	can_produce := make(chan struct{}, 1)
	can_produce <- struct{}{}

	
	go consumer(&buffer_simple, can_consume, can_produce)
	go producer(&buffer_simple, can_consume, can_produce)

	time.Sleep(4 * time.Second)
}