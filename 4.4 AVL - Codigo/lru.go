package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Pelicula struct {
	ID     int
	Titulo string
	Prom   float64
}

type NodoAVL struct {
	clave    float64
	datos    []Pelicula
	alt      int
	izq, der *NodoAVL
}

var rotaciones int

func Altura(n *NodoAVL) int {
	if n == nil {
		return 0
	}
	return n.alt
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func rotarDer(y *NodoAVL) *NodoAVL {
	rotaciones++
	x, t2 := y.izq, y.izq.der
	x.der = y
	y.izq = t2
	y.alt = max(Altura(y.izq), Altura(y.der)) + 1
	x.alt = max(Altura(x.izq), Altura(x.der)) + 1
	return x
}

func rotarIzq(x *NodoAVL) *NodoAVL {
	rotaciones++
	y, t2 := x.der, x.der.izq
	y.izq = x
	x.der = t2
	x.alt = max(Altura(x.izq), Altura(x.der)) + 1
	y.alt = max(Altura(y.izq), Altura(y.der)) + 1
	return y
}

func Insertar(raiz *NodoAVL, clave float64, dato Pelicula) *NodoAVL {
	if raiz == nil {
		return &NodoAVL{clave: clave, datos: []Pelicula{dato}, alt: 1}
	}
	if clave < raiz.clave {
		raiz.izq = Insertar(raiz.izq, clave, dato)
	} else if clave > raiz.clave {
		raiz.der = Insertar(raiz.der, clave, dato)
	} else {
		raiz.datos = append(raiz.datos, dato)
		return raiz
	}

	raiz.alt = max(Altura(raiz.izq), Altura(raiz.der)) + 1
	bal := Altura(raiz.izq) - Altura(raiz.der)

	if bal > 1 && clave < raiz.izq.clave {
		return rotarDer(raiz)
	}
	if bal < -1 && clave > raiz.der.clave {
		return rotarIzq(raiz)
	}
	if bal > 1 && clave > raiz.izq.clave {
		raiz.izq = rotarIzq(raiz.izq)
		return rotarDer(raiz)
	}
	if bal < -1 && clave < raiz.der.clave {
		raiz.der = rotarDer(raiz.der)
		return rotarIzq(raiz)
	}
	return raiz
}

func ConsultaRango(raiz *NodoAVL, a, b float64) []Pelicula {
	if raiz == nil {
		return nil
	}
	var res []Pelicula
	if a < raiz.clave {
		res = append(res, ConsultaRango(raiz.izq, a, b)...)
	}
	if a <= raiz.clave && raiz.clave <= b {
		res = append(res, raiz.datos...)
	}
	if raiz.clave < b {
		res = append(res, ConsultaRango(raiz.der, a, b)...)
	}
	return res
}

func BalanceOK(n *NodoAVL) bool {
	if n == nil {
		return true
	}
	b := Altura(n.izq) - Altura(n.der)
	return b >= -1 && b <= 1 && BalanceOK(n.izq) && BalanceOK(n.der)
}

func Contar(n *NodoAVL) int {
	if n == nil {
		return 0
	}
	return 1 + Contar(n.izq) + Contar(n.der)
}

func leerPeliculas(ruta string) (map[int]string, error) {
	f, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.FieldsPerRecord = -1

	if _, err := r.Read(); err != nil {
		return nil, err
	}

	m := map[int]string{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(row) < 2 {
			continue
		}
		id, _ := strconv.Atoi(strings.TrimSpace(row[0]))
		m[id] = row[1]
	}
	return m, nil
}

func leerPromedios(ruta string) (map[int]float64, map[int]int, error) {
	f, err := os.Open(ruta)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.FieldsPerRecord = -1

	if _, err := r.Read(); err != nil {
		return nil, nil, err
	}

	suma := map[int]float64{}
	cant := map[int]int{}

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(row) < 3 {
			continue
		}
		id, err1 := strconv.Atoi(strings.TrimSpace(row[1]))
		rat, err2 := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err1 != nil || err2 != nil {
			continue
		}
		suma[id] += rat
		cant[id]++
	}
	return suma, cant, nil
}

func main() {
	movies := flag.String("movies", "movies.csv", "ruta de movies.csv")
	ratings := flag.String("ratings", "ratings.csv", "ruta de ratings.csv")
	a := flag.Float64("a", 3.5, "inicio del rango")
	b := flag.Float64("b", 4.0, "fin del rango")
	flag.Parse()

	titulos, err := leerPeliculas(*movies)
	if err != nil {
		fmt.Println("error movies:", err)
		return
	}

	suma, cant, err := leerPromedios(*ratings)
	if err != nil {
		fmt.Println("error ratings:", err)
		return
	}

	var raiz *NodoAVL
	for id, s := range suma {
		avg := s / float64(cant[id])
		raiz = Insertar(raiz, avg, Pelicula{ID: id, Titulo: titulos[id], Prom: avg})
	}

	res := ConsultaRango(raiz, *a, *b)

	fmt.Printf("Resultados en [%.2f, %.2f]\n", *a, *b)
	for _, p := range res {
		fmt.Printf("%.2f | %d | %s\n", p.Prom, p.ID, p.Titulo)
	}
	fmt.Printf("\nAltura: %d\nNodos: %d\nRotaciones: %d\nBalanceado: %v\n",
		Altura(raiz), Contar(raiz), rotaciones, BalanceOK(raiz))
}
