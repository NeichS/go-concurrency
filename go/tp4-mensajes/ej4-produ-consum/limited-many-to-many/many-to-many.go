package main

import (
	"fmt"
	"math/rand"
	"time"
)


func producer(can_produce, start chan struct{}, can_consume chan int, i int) {

	<-start

	for {
		<-can_produce
		produced_value := rand.Intn(10000)
		fmt.Println("El productor ", i, " produjo el valor ", produced_value)
		can_consume <- produced_value
	}
}

func consumer(can_produce, start chan struct{}, can_consume chan int, i int) {

	<-start

	for {
		consumed_value := <- can_consume
		fmt.Println("El consumidor ", i , " consumio el valor ", consumed_value)
		can_produce <- struct{}{}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	can_consume := make(chan int)
	can_produce := make(chan struct{}, 5) 
	for i:= 0 ; i < 5; i++ {
		can_produce <- struct{}{} 
	}

	start := make(chan struct{})

	for i := 0 ; i < 100 ; i ++ {
		go consumer(can_produce, start, can_consume, i)
	}

	for i := 0 ; i < 20 ; i ++ {
		go producer(can_produce, start, can_consume, i)
	}


	close(start)
	time.Sleep(5 * time.Second) //le doy tiempo de vida al programa para que las gorotuines se ejecuten
}