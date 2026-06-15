package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Nodo struct {
	clave int
	valor int

	prev *Nodo
	next *Nodo
}

type LRU struct {
	capacidad int

	mapa map[int]*Nodo

	head *Nodo
	tail *Nodo
}

func nuevaLRU(num int) *LRU {
	return &LRU{
		capacidad: num,
		mapa:      make(map[int]*Nodo),
	}
}

func (l *LRU) agregarAlFrente(n *Nodo) {
	n.prev = nil
	n.next = l.head

	if l.head != nil {
		l.head.prev = n
	}

	l.head = n

	if l.tail == nil {
		l.tail = n
	}

}

func (l *LRU) eliminarNodo(n *Nodo) {

	if n.prev != nil {
		n.prev.next = n.next
	} else { // por si es el head
		l.head = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else { // por si es el tail
		l.tail = n.prev
	}
}

func (l *LRU) moverAlFrente(n *Nodo) {
	l.eliminarNodo(n)
	l.agregarAlFrente(n)
}

func (l *LRU) eliminarTail() *Nodo {

	if l.tail == nil {
		return nil
	}

	eliminado := l.tail //Lo guardo para eliminarlo del mapa despues

	l.eliminarNodo(eliminado)

	return eliminado
}

func (l *LRU) Put(clave int, valor int) {

	if nodo, existe := l.mapa[clave]; existe {
		nodo.valor = valor
		l.moverAlFrente(nodo)

		return
	}
	if len(l.mapa) >= l.capacidad {
		eliminado := l.eliminarTail()
		delete(l.mapa, eliminado.clave)
	}

	nuevo := &Nodo{
		clave: clave,
		valor: valor,
	}

	l.agregarAlFrente(nuevo)
	l.mapa[clave] = nuevo
}

func (l *LRU) Get(clave int) (int, bool) {
	nodo, existe := l.mapa[clave]

	if !existe {
		return 0, false
	}

	l.moverAlFrente(nodo)

	return nodo.valor, true
}

func CargarSecuencia(ruta string) []int {

	archivo, err := os.Open(ruta)

	if err != nil {
		fmt.Println("Error al abrir archivo:", err)
		return nil
	}

	defer archivo.Close()

	lector := csv.NewReader(archivo)

	registros, err := lector.ReadAll()

	if err != nil {
		fmt.Println("Error al leer CSV:", err)
		return nil
	}

	var secuencia []int

	for i := 1; i < len(registros); i++ {

		movieId, err := strconv.Atoi(registros[i][1])

		if err != nil {
			continue
		}

		secuencia = append(secuencia, movieId)
	}

	return secuencia
}
func main() {

	secuencia := CargarSecuencia("ratings.csv")
	capacidades := []int{50, 100, 500, 1000}

	for _, cap := range capacidades {

		cache := nuevaLRU(cap)

		hits := 0
		accesos := 0

		for _, movieId := range secuencia {

			_, existe := cache.Get(movieId)

			if existe {
				hits++
			} else {
				cache.Put(movieId, movieId)
			}

			accesos++
		}

		hitRatio := float64(hits) / float64(accesos)

		fmt.Printf("Capacidad: %d\n", cap)
		fmt.Printf("Hits: %d\n", hits)
		fmt.Printf("Hit Ratio: %.2f%%\n\n", hitRatio*100)
	}
}
