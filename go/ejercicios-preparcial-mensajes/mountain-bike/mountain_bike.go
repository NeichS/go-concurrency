package main

import (
	"fmt"
	"time"
)

func op_one(rueda chan struct{}) {
	//prerara rueda (una a la vez)
	for {
		rueda <- struct{}{}
	}
}

func op_two(cuadro chan struct{}) {
	//prerara cuadro
	for {
		cuadro <- struct{}{}
	}
}

func op_three(manillar chan struct{}) {

	//prepara manillar
	for {
		manillar <- struct{}{}
	}
}

func mounter(rueda, cuadro, manillar chan struct{}) {

	bicis_finalizadas := 0

	for {
		//toma 2 ruedas
		for i := 0; i < 2; i++ {
			<-rueda
		}

		//toma cuadro
		<-cuadro

		//manillar
		<-manillar

		//se arma la bici
		bicis_finalizadas++
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Bicis manufacturadas ", bicis_finalizadas)
	}
	
}

func main() {

	rueda := make(chan struct{},1)
	cuadro := make(chan struct{},1)
	manillar := make(chan struct{},1)

	go op_one(rueda)
	go op_two(cuadro)
	go op_three(manillar)
	go mounter(rueda, cuadro, manillar)

	time.Sleep(3 * time.Second)
}