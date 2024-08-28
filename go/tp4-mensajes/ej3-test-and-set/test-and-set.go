package main

import (
	"fmt" 
	"time"
	"sync"
)

func test_and_set(lock *bool, mutex chan struct{}) bool {
	<- mutex //espera a que haya un mensaje
	initial := *lock
	*lock = true
	mutex <- struct{}{} //envia el mensaje para liberar
	return initial
}


func generic_func(mutex chan struct{}, lock *bool, i int, wg *sync.WaitGroup, start chan struct{} ) {
	
	defer wg.Done() // decrementa el contador del WaitGroup al final
	<-start // espera la señal de inicio

	for (test_and_set(lock, mutex)){
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Println("accediendo a la seccion critica, proceso ", i )

	*lock = false
}

func main() {

	mutex := make(chan struct{}, 1) //struct{}{} es una estructura vacía y se usa comúnmente para señalización.
	mutex <- struct{}{} //meto un mensaje vacio

	lock := false
	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 1 ; i < 250 ; i++ {
		wg.Add(1)
		go generic_func(mutex, &lock, i, &wg, start)
	}

	time.Sleep(2 * time.Second)
	
	close(start) //es como un signal all 

	wg.Wait()
}