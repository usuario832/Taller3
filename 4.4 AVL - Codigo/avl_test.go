package main

import "testing"

func TestArbolVacio(t *testing.T) {
	if Altura(nil) != 0 {
		t.Fatal("la altura de un árbol vacío debe ser 0")
	}

	if !BalanceOK(nil) {
		t.Fatal("un árbol vacío debe estar balanceado")
	}

	res := ConsultaRango(nil, 3.0, 4.0)
	if len(res) != 0 {
		t.Fatal("la consulta en árbol vacío debe devolver 0 resultados")
	}
}

func TestInsertarUnElemento(t *testing.T) {
	var raiz *NodoAVL

	p := Pelicula{ID: 1, Titulo: "Pelicula prueba", Prom: 4.0}
	raiz = Insertar(raiz, p.Prom, p)

	if raiz == nil {
		t.Fatal("la raíz no debe ser nil")
	}

	if Altura(raiz) != 1 {
		t.Fatal("la altura debe ser 1")
	}

	if !BalanceOK(raiz) {
		t.Fatal("el árbol debe estar balanceado")
	}
}

func TestInsertarVariosBalanceado(t *testing.T) {
	var raiz *NodoAVL

	for i := 1; i <= 10; i++ {
		p := Pelicula{ID: i, Titulo: "Pelicula", Prom: float64(i)}
		raiz = Insertar(raiz, p.Prom, p)
	}

	if !BalanceOK(raiz) {
		t.Fatal("el árbol debe mantenerse balanceado")
	}

	if Contar(raiz) != 10 {
		t.Fatal("el árbol debe tener 10 nodos")
	}
}

func TestConsultaRango(t *testing.T) {
	var raiz *NodoAVL

	for i := 1; i <= 5; i++ {
		p := Pelicula{ID: i, Titulo: "Pelicula", Prom: float64(i)}
		raiz = Insertar(raiz, p.Prom, p)
	}

	res := ConsultaRango(raiz, 2.0, 4.0)

	if len(res) != 3 {
		t.Fatalf("se esperaban 3 resultados, se obtuvo %d", len(res))
	}

	if res[0].Prom != 2.0 || res[1].Prom != 3.0 || res[2].Prom != 4.0 {
		t.Fatal("los resultados no están en orden ascendente")
	}
}

func TestConsultaSinResultados(t *testing.T) {
	var raiz *NodoAVL

	for i := 1; i <= 5; i++ {
		p := Pelicula{ID: i, Titulo: "Pelicula", Prom: float64(i)}
		raiz = Insertar(raiz, p.Prom, p)
	}

	res := ConsultaRango(raiz, 10.0, 20.0)

	if len(res) != 0 {
		t.Fatal("la consulta no debería devolver resultados")
	}
}

func TestMismaClave(t *testing.T) {
	var raiz *NodoAVL

	p1 := Pelicula{ID: 1, Titulo: "Pelicula A", Prom: 4.0}
	p2 := Pelicula{ID: 2, Titulo: "Pelicula B", Prom: 4.0}

	raiz = Insertar(raiz, p1.Prom, p1)
	raiz = Insertar(raiz, p2.Prom, p2)

	res := ConsultaRango(raiz, 4.0, 4.0)

	if len(res) != 2 {
		t.Fatalf("se esperaban 2 películas con la misma clave, se obtuvo %d", len(res))
	}
}

func BenchmarkInsertarAVL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var raiz *NodoAVL

		for j := 0; j < 1000; j++ {
			p := Pelicula{ID: j, Titulo: "Pelicula", Prom: float64(j)}
			raiz = Insertar(raiz, p.Prom, p)
		}
	}
}

func BenchmarkConsultaRangoAVL(b *testing.B) {
	var raiz *NodoAVL

	for i := 0; i < 10000; i++ {
		p := Pelicula{ID: i, Titulo: "Pelicula", Prom: float64(i)}
		raiz = Insertar(raiz, p.Prom, p)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ConsultaRango(raiz, 3000.0, 4000.0)
	}
}
