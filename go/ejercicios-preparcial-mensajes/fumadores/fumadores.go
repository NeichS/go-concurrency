package main

import (
	"fmt"
	"math/rand"
	"time"
)

func fumador_uno(tabaco, fin_fumar chan struct{}){
	
	for {
		//espera a que falte tabaco
		<-tabaco

		fmt.Println("El fumador uno esta fumando")
		time.Sleep(120 * time.Millisecond)

		//avisa que termina de fumar

		fin_fumar <- struct{}{}
	}
	
}	

func fumador_dos(papel, fin_fumar chan struct{}) {
	for {
		//espera a que falte papel
		<-papel

		fmt.Println("El fumador dos esta fumando")
		time.Sleep(120 * time.Millisecond)
		//avisa que termina de fumar
		fin_fumar <- struct{}{}
	}
		
}

func fumador_tres(fosforo, fin_fumar chan struct{}) {
	for {
		//espera a que falte fosforos
		<-fosforo

		fmt.Println("El fumador tres esta fumando")
		time.Sleep(120 * time.Millisecond)
		//avisa que termina de fumar
		fin_fumar <- struct{}{}
	}
		
}

func agente(ingredientes [3]chan struct{}, fin_fumar chan struct{}) {
	

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for {
		//coloca 2 ingredientes
		ingrediente_faltante := rng.Intn(3)
		if ingrediente_faltante == 0 {
			ingredientes[0] <- struct{}{}
			fmt.Println("El agente espera a que termine tabaco")
			<-fin_fumar
		} else if ingrediente_faltante == 1{
			ingredientes[1] <- struct{}{}
			fmt.Println("El agente espera a que termine papel")
			<-fin_fumar
		} else {
			ingredientes[2] <- struct{}{}
			fmt.Println("El agente espera a que termine fosforo")
			<-fin_fumar
		}
	}
	
	
}

func main(){

	var ingredientes [3]chan struct{}
	for i:= 0; i < len(ingredientes); i++ {
		ingredientes[i] = make(chan struct{},1)
	}

	fin_fumar := make(chan struct{},1)


	go fumador_uno(ingredientes[0], fin_fumar)
	go fumador_dos(ingredientes[1], fin_fumar)
	go fumador_tres(ingredientes[2], fin_fumar)
	go agente(ingredientes, fin_fumar)

	time.Sleep(3 * time.Second)
}