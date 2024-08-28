package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//quedan inutilizados los canales pb y p1 ya que necesito saber el valor de autos que hay y ya con eso controlo
//los pisos

func enter_procedure(
	entrance,
	pb,
	p1, 
	mutex_p1, 
	mutex_pb, 
	mutex_tot chan struct{}, 
	cars_pb, 
	cars_p1, 
	cars_total *int, 
	i, 
	N int,
) string{

	<-mutex_tot
	if *cars_total < 2*N{
		mutex_tot <- struct{}{}
		
		<-entrance
		fmt.Println("El auto ", i, " esta ingresando")
		entrance <- struct{}{}

		if *cars_pb < N {
			<-mutex_pb
			*cars_pb++
			mutex_pb <- struct{}{}
			<- pb //reclama un "espacio" del canal pb (medio al pedo)
			
			fmt.Println("Ingresa el auto ", i, " al piso pb")
			<-mutex_tot
			*cars_total++
			mutex_tot <- struct{}{}

			return "pb"
		} else {
			<-mutex_p1
			*cars_p1++
			mutex_p1 <- struct{}{}
			<- p1 //reclama un mensaje del canal p1

			fmt.Println("Ingresa el auto ", i, " al piso p1")
			<-mutex_tot
			*cars_total++
			mutex_tot <- struct{}{}

			return "p1"
		}
	} else {
		mutex_tot<- struct{}{}
		fmt.Println("El auto ", i, " no pudo ingresar")
		return "failed"
	}
}

func exit_procedure(
	entrance, 
	pb, 
	p1, 
	mutex_p1, 
	mutex_pb, 
	mutex_tot chan struct{}, 
	cars_pb, 
	cars_p1, 
	cars_total *int, 
	result string, 
	i int,
){
	if result == "pb" {
		<-mutex_pb
		*cars_pb--
		mutex_pb <- struct{}{}
		pb<-struct{}{}
		<-mutex_tot
		*cars_total--
		mutex_tot <- struct{}{}


		<-entrance
		fmt.Println("El auto ", i, " esta saliendo, liberando pb")
		entrance<-struct{}{}

	} else {
		<-mutex_p1
		*cars_p1--
		mutex_p1 <- struct{}{}
		p1<-struct{}{}
		<-mutex_tot
		*cars_total--
		mutex_tot <- struct{}{}

		<-entrance
		fmt.Println("El auto ", i, " esta saliendo, liberando p1")
		entrance<-struct{}{}
	}
}


func cars(entrance, pb, p1, mutex_p1, mutex_pb, mutex_tot chan struct{}, cars_pb, cars_p1, cars_total *int, wg *sync.WaitGroup, i, N int) {
	defer wg.Done()
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	ran_number := rng.Intn(2000)
	//tiempo random antes de llegar al estacionamiento
	time.Sleep(time.Duration(ran_number) * time.Millisecond)

	var result string
	exit := false
	for (exit == false){
		result = enter_procedure(entrance, pb, p1, mutex_p1, mutex_pb, mutex_tot, cars_pb, cars_p1, cars_total, i, N)
		if result == "failed" {
			time.Sleep(1 * time.Second) //intenta entrar cada 1 segundo
		} else {
			exit = true
		}
	}
	
	ran_number2 := rng.Intn(2000)
	//tiempo random en el que se queda estacionado
	time.Sleep(time.Duration(ran_number2) * time.Millisecond)

	exit_procedure(entrance, pb, p1, mutex_p1, mutex_pb, mutex_tot, cars_pb, cars_p1, cars_total, result, i)
}

func main() {
	entrance := make(chan struct{},1)
	entrance <- struct{}{} 

	const N = 20

	pb := make(chan struct{}, N )
	p1 := make(chan struct{}, N )

	for i := 0 ; i < N ; i++ {
		pb <- struct{}{}
		p1 <- struct{}{}
	} 
	
	cars_pb := 0
	mutex_pb := make(chan struct{},1)
	mutex_pb <- struct{}{}
	
	cars_p1 := 0
	mutex_p1 := make(chan struct{},1)
	mutex_p1 <- struct{}{}
	
	cars_total := 0
	mutex_tot := make(chan struct{},1)
	mutex_tot <- struct{}{}
	
	var wg sync.WaitGroup
	for i := 1 ; i < 500; i++ {
		wg.Add(1)
		go cars(entrance, pb, p1, mutex_p1, mutex_pb, mutex_tot, &cars_pb, &cars_p1, &cars_total, &wg, i, N)
	}
	
	wg.Wait()

	fmt.Println("Todos los autos pudieron salir y entrar correctamente")
}