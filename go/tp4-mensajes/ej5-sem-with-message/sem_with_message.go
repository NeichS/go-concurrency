package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	num int
	queue chan struct{}
	mutex chan struct{}
}

//inicializo el semaforo con la capacidad requerida
func initializer(capacity int, mutex chan struct{}) chan struct{} {
	for i := 0 ; i < capacity ; i++ { 
		mutex <- struct{}{}
	}
	return mutex
}

func NewSemaphore(capacity int) *Semaphore {
	
	mutex := make(chan struct{}, capacity)
	mutex = initializer(capacity, mutex)
	return &Semaphore{
		num : capacity,
		queue: make(chan struct{}),
		mutex: mutex,
	}
}

func (s *Semaphore) wait() {
	<-s.mutex
	s.num--
	if s.num < 0 {
		s.mutex <- struct{}{}
		<-s.queue
	} else {
		s.mutex <- struct{}{}
	}
}

func (s *Semaphore) signal() {
	<-s.mutex
	s.num++
	if s.num <= 0{
		s.queue <- struct{}{}
		s.mutex <- struct{}{}
	} else {
		s.mutex <- struct{}{}
	}
}


func proceso (id int, sem *Semaphore, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Proceso %d esperando...\n", id)
	sem.wait() // Esperar en el semáforo

	fmt.Printf("Proceso %d accediendo a la sección crítica...\n", id)
	time.Sleep(1 * time.Second) // Simular trabajo en la sección crítica

	fmt.Printf("Proceso %d saliendo de la sección crítica...\n", id)
	sem.signal() // Señalizar la salida de la sección crítica
}

func main() {
	my_semaphore := NewSemaphore(1)
	var wg sync.WaitGroup
	
	for i := 1; i <= 15; i++ {
		wg.Add(1)
		go proceso(i, my_semaphore, &wg)
	}

	wg.Wait() // Esperar a que todos los procesos terminen
	fmt.Println("Todos los procesos han terminado.")
}