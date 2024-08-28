package main

import (
        "fmt"
        "sync"
)

const (
        clientesMax = 5 // Máximo de clientes en la barbería
        barberos   = 2  // Número de barberos
        sillasEspera = clientesMax - 2 // Sillas de espera
)

type Barberia struct {
        barberos     []*Barbero
        clientes     []*Cliente
        semillas     chan struct{}
        semCobrar    *sync.Mutex
        semSillas    *sync.Mutex
        cond         *sync.Cond
}

type Barbero struct {
        silla       int
        cobrando   bool
        sillaEspera int
}

type Cliente struct {
        silla  int
        cortado bool
        pagado bool
}

func (b *Barberia) run() {
        for {
                // Esperar a que haya un cliente
                <-b.semillas
                cliente := b.clientes[0]
                cliente.cortado = true

                // Cortar el pelo al cliente
                fmt.Println("Cortando el pelo al cliente", cliente.silla)

                // Liberar la silla del cliente
                b.semSillas.Lock()
                b.clientes[0] = nil
                b.semSillas.Unlock()

                // Cobrar al cliente
                b.semCobrar.Lock()
                cliente.pagado = true
                b.semCobrar.Unlock()

                // Despertar al cliente
                b.cond.Signal()
        }
}

func (b *Barberia) main() {
        b.barberos = make([]*Barbero, barberos)
        b.clientes = make([]*Cliente, clientesMax)
        b.semillas = make(chan struct{}, clientesMax)
        b.semCobrar = &sync.Mutex{}
        b.semSillas = &sync.Mutex{}
        b.cond = sync.NewCond(b.semSillas)

        // Inicializar los barberos
        for i := 0; i < barberos; i++ {
                barbero := &Barbero{i, false, i}
                b.barberos[i] = barbero
                go b.barberos[i].run()
        }

        // Inicializar los clientes
        for i := 0; i < clientesMax; i++ {
                cliente := &Cliente{i, false, false}
                b.clientes[i] = cliente
        }

        // Esperar a que los barberos despierten
        for i := 0; i < barberos; i++ {
                b.cond.Wait()
        }

        // Probar la barbería
        for i := 0; i < 10; i++ {
                cliente := &Cliente{i, false, false}
                b.clientes[i] = cliente
                b.semillas <- struct{}{}
        }

        // Esperar a que los clientes salgan
        for i := 0; i < clientesMax; i++ {
                b.cond.Wait()
        }

        fmt.Println("La barbería ha cerrado")
}

func (b *Barbero) run() {
        for {
                // Esperar a que haya un cliente
                b.cond.Wait()

                // Aceptar un cliente
                b.sillaEspera++
                b.cobrando = true

                // Esperar a que el cliente escoja su silla
                b.semSillas.Lock()
                for i := 0; i < clientesMax; i++ {
                        if b.clientes[i] != nil && !b.clientes[i].cortado {
                                b.clientes[i].silla = b.sillaEspera
                                b.clientes[i].cortado = false
                                break
                        }
                }
                b.semSillas.Unlock()

                // Cortar el pelo al cliente
                fmt.Println("Cortando el pelo al cliente", b.sillaEspera)

                // Liberar la silla del cliente
                b.semSillas.Lock()
                for i := 0; i < clientesMax; i++ {
                        if b.clientes[i] != nil && b.clientes[i].silla == b.sillaEspera {
                                b.clientes[i] = nil
                                break
                        }
                }
                b.semSillas.Unlock()

                // Cobrar al cliente
                b.semCobrar.Lock()
                if b.cobrando {
                        b.cobrando = false
                        // Cobrar al cliente
                        fmt.Println("Cobrando al cliente", b.sillaEspera)
                }
                b.semCobrar.Unlock()

                // Esperar a que el cliente pague
                b.cond.Wait()

                // Liberar la silla del cliente
                b.semSillas
        }
}