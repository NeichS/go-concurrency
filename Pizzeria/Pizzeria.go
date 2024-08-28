package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

// Definimos un semáforo utilizando un canal
type Semaphore struct{
	sem chan struct{}
}

func NewSemaphore(size int) *Semaphore {
	return &Semaphore{
		sem: make(chan struct{}, size),
	}
}

func (s *Semaphore) Wait() {
	s.sem <- struct{}{}
}

func (s *Semaphore) Signal() {
	<-s.sem
}

// Definimos un tipo Set (conjunto) que es un mapa donde la clave es un string y el valor es un struct vacío.
type Set map[string]struct{}

// Función para agregar un elemento al conjunto.
func (s Set) Add(elemento string) {
	s[elemento] = struct{}{}
}

// Función para verificar si un elemento está presente en el conjunto.
func (s Set) Contains(elemento string) bool {
	_, existe := s[elemento]
	return existe
}

// Función para eliminar un elemento del conjunto.
func (s Set) Remove(elemento string) {
	delete(s, elemento)
}

func (s Set) RandomElement() string {
	elementos := make([]string, 0, len(s))
	for key := range s {
		elementos = append(elementos, key)
	}

	if len(elementos) == 0 {
		return ""
	}

	indice := rand.Intn(len(elementos))
	return elementos[indice]
}

func Difference(a, b Set) Set {
	resultado := make(Set)
	for key := range a {
		if !b.Contains(key) {
			resultado.Add(key)
		}
	}
	return resultado
}

var pizza = make(Set)
var mutex = &sync.Mutex{}

var pizza_completa = NewSemaphore(0)
var falta_queso = NewSemaphore(0)
var falta_salsa = NewSemaphore(0)
var falta_tomate = NewSemaphore(0)

func chef() {
	ingredientes := make(Set)
	ingredientes.Add("salsa")
	ingredientes.Add("tomate")
	ingredientes.Add("queso")

	contador_pizzas := 0

	for {
		mutex.Lock()
		pizza.Add(ingredientes.RandomElement())
		ingredientes_diferencia := Difference(ingredientes, pizza)
		if len(ingredientes_diferencia) > 0 {
			pizza.Add(ingredientes_diferencia.RandomElement())
		}

		if !pizza.Contains("salsa") {
			mutex.Unlock()
			fmt.Println("Faltaba salsa")
			falta_salsa.Signal()
		} else if !pizza.Contains("queso") {
			mutex.Unlock()
			fmt.Println("Faltaba queso")
			falta_queso.Signal()
		} else if !pizza.Contains("tomate") {
			mutex.Unlock()
			fmt.Println("Faltaba tomate")
			falta_tomate.Signal()
		} else {
			mutex.Unlock()
			fmt.Println("Pizza completa")
			pizza_completa.Signal()
		}

		pizza_completa.Wait()
		mutex.Lock()
		// Resetear la pizza para la siguiente iteración
		pizza = make(Set)
		mutex.Unlock()

		contador_pizzas++
		fmt.Printf("Pizzas hechas %s \n", strconv.Itoa(contador_pizzas))
	}
}

func ayudanteQueso() {
	for {
		falta_queso.Wait()
		mutex.Lock()
		pizza.Add("queso")
		mutex.Unlock()
		fmt.Println("El ayudante agregó queso")
		pizza_completa.Signal()
	}
}

func ayudanteSalsa() {
	for {
		falta_salsa.Wait()
		mutex.Lock()
		pizza.Add("salsa")
		mutex.Unlock()
		fmt.Println("El ayudante agregó salsa")
		pizza_completa.Signal()
	}
}

func ayudanteTomate() {
	for {
		falta_tomate.Wait()
		mutex.Lock()
		pizza.Add("tomate")
		mutex.Unlock()
		fmt.Println("El ayudante agregó tomate")
		pizza_completa.Signal()
	}
}

func main() {
	go chef()
	go ayudanteTomate()
	go ayudanteQueso()
	go ayudanteSalsa()

	select {}
}
