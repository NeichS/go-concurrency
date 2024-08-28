package main

import (
	"sync"
	"time"
)

func esquiadores(
	medio_elevacion,
	llega [N]chan struct{},
	mutex_siguiente_silla,
	start chan struct{},
	siguiente_silla *int,
	wg *sync.WaitGroup,
) {

	defer wg.Done()
	<-start

	for {
		<-mutex_siguiente_silla
		silla_actual := *siguiente_silla
		medio_elevacion[*siguiente_silla] <- struct{}{} // se sube a la silla
		mutex_siguiente_silla <- struct{}{}  // Se libera el mutex aquí
		<-llega[silla_actual]
	}
}

func silla(
	silla_espacio,
	llega,
	mutex_esquiadores,
	mutex_siguiente_silla,
	start chan struct{},
	cantidad_esquiadores *int,
	siguiente_silla *int,
	wg *sync.WaitGroup,
) {

	defer wg.Done()

	<-start 

	for {
		<-silla_espacio // receive
		<-silla_espacio // receive
		<-silla_espacio // receive
		<-silla_espacio // receive
		
		<-mutex_esquiadores
		*cantidad_esquiadores = *cantidad_esquiadores + 4 
		mutex_esquiadores <- struct{}{}

		<-mutex_siguiente_silla
		*siguiente_silla = (*siguiente_silla + 1) % N
		mutex_siguiente_silla <- struct{}{}

		// la silla esta viajando sleep
		time.Sleep(100 * time.Millisecond) // Simula el viaje de la silla

		// la silla llega a destino
		llega <- struct{}{}
		llega <- struct{}{}
		llega <- struct{}{}
		llega <- struct{}{}
		
		<-mutex_esquiadores
		*cantidad_esquiadores = *cantidad_esquiadores - 4 
		mutex_esquiadores <- struct{}{}
	}
}

const N int = 150

func main() {
	cantidad_esquiadores := 0
	var medio_elevacion [N]chan struct{}
	mutex_esquiadores := make(chan struct{}, 1)
	mutex_esquiadores <- struct{}{}

	siguiente_silla := 0
	mutex_siguiente_silla := make(chan struct{}, 1)
	mutex_siguiente_silla <- struct{}{}

	start := make(chan struct{})

	var wg sync.WaitGroup

	var llega [N]chan struct{}
	for i := 0; i < N; i++ {
		medio_elevacion[i] = make(chan struct{}, 4)
		llega[i] = make(chan struct{}, 4)
		wg.Add(1)
		go silla(
			medio_elevacion[i], 
			llega[i], 
			mutex_esquiadores, 
			mutex_siguiente_silla,
			start, 
			&cantidad_esquiadores, 
			&siguiente_silla,
			&wg,
		)
	}

	for i := 0; i < 300; i++ {
		wg.Add(1)
		go esquiadores(
			medio_elevacion,
			llega, 
			mutex_siguiente_silla,
			start, 
			&siguiente_silla,
			&wg,
		)
	}

	close(start)
	wg.Wait()
}
