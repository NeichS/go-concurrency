package main

import (
	"fmt"
	"sync"
	"time"
)


func reader(readers_active *int, mutex_ra, start,readWrite chan struct{}, i int, wg *sync.WaitGroup) {

	defer wg.Done()

	<-start

	<-mutex_ra
	*readers_active++
	if *readers_active == 1 {
		mutex_ra <- struct{}{}
		<-readWrite
	} else {
		mutex_ra <- struct{}{}
	} 


	//lee 
	fmt.Println("El lector ", i, " lee")
	time.Sleep(80 * time.Millisecond)
	
	<-mutex_ra
	*readers_active--
	if *readers_active == 0 {
		mutex_ra <- struct{}{}
		readWrite <- struct{}{}
	} else {
		mutex_ra <- struct{}{}
	}
}

func writer(readWrite, start chan struct{}, i int, wg *sync.WaitGroup) {

	defer wg.Done()

	<-start

	<-readWrite

	//escribo
	fmt.Println("El escritor ", i, " escribe")

	readWrite <- struct{}{}

}

func main() {

	reader_active := 0 
	mutex_ra := make(chan struct{},1)
	mutex_ra <- struct{}{}

	readWrite := make(chan struct{},1)
	readWrite <- struct{}{}

	start := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0 ; i < 20 ; i++ {
		wg.Add(2)
		go writer(readWrite, start, i, &wg)
		go reader(&reader_active, mutex_ra, start, readWrite, i, &wg)
	}

	close(start)
	wg.Wait()

}