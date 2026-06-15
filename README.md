## 4.3 Caché LRU mediante Lista Doblemente Enlazada

### Descripción

Se implementó una caché LRU (Least Recently Used) utilizando una lista doblemente enlazada y una tabla hash (`map`). Esta combinación permite realizar operaciones de búsqueda, inserción y actualización de elementos en tiempo O(1).

La lista doblemente enlazada mantiene el orden de uso de los elementos, mientras que el mapa permite acceder rápidamente a cada nodo mediante su clave.

### Estructuras utilizadas

#### Nodo

```go
type Nodo struct {
    clave int
    valor int
    prev *Nodo
    next *Nodo
}
```

Cada nodo almacena una clave, un valor y referencias al nodo anterior y siguiente.

#### LRU

```go
type LRU struct {
    capacidad int
    mapa map[int]*Nodo
    head *Nodo
    tail *Nodo
}
```

- `head`: elemento más recientemente utilizado.
- `tail`: elemento menos recientemente utilizado.
- `mapa`: acceso rápido a los nodos.
- `capacidad`: tamaño máximo de la caché.

### Operaciones principales

#### Get(clave)

Busca una clave dentro de la caché.

- Si existe, el nodo se mueve al frente de la lista.
- Si no existe, retorna que la clave no fue encontrada.

#### Put(clave, valor)

Inserta o actualiza un elemento.

- Si la clave ya existe, actualiza su valor y mueve el nodo al frente.
- Si la clave no existe y la caché está llena, elimina el nodo menos recientemente utilizado (`tail`).
- Finalmente agrega el nuevo nodo al frente de la lista.

### Complejidad

| Operación | Complejidad |
|------------|------------|
| Get | O(1) |
| Put | O(1) |
| Búsqueda en mapa | O(1) |
| Eliminación de tail | O(1) |
| Movimiento al frente | O(1) |

La combinación de una tabla hash y una lista doblemente enlazada permite mantener todas las operaciones principales en tiempo constante O(1).

### Lectura del dataset

Se utilizó el archivo `ratings.csv`. De cada registro se extrajo el campo `movieId`, generando una secuencia de accesos que posteriormente fue utilizada para evaluar el comportamiento de la caché.

### Resultados

| Capacidad | Hits | Hit Ratio |
|------------|--------:|-----------:|
| 50 | 784 | 0.78% |
| 100 | 2481 | 2.46% |
| 500 | 19629 | 19.47% |
| 1000 | 38399 | 38.08% |

### Conclusión

Los resultados muestran que al incrementar la capacidad de la caché aumenta el hit ratio. Esto se debe a que una mayor cantidad de elementos puede mantenerse almacenada, reduciendo la probabilidad de eliminar elementos que serán utilizados nuevamente en accesos futuros.



# Taller 3 — Índice AVL con consultas por rango

## 1. Descripción del proyecto

Este proyecto corresponde al problema **4.4 Árboles binarios — Índice AVL con consultas por rango**.

El objetivo es implementar desde cero un árbol AVL en Go para indexar películas del dataset MovieLens. La clave utilizada para el árbol es el **rating promedio** de cada película, calculado a partir del archivo `ratings.csv`.

Luego de construir el árbol, el programa permite realizar consultas por rango, por ejemplo:

```text
[3.5, 4.0]
```

Esto devuelve todas las películas cuyo rating promedio se encuentra dentro de ese intervalo.

El árbol AVL mantiene su balance automáticamente mediante rotaciones, lo que permite conservar una altura O(log n).

---

## 2. Estructura del proyecto

```text
curso vsc/
│
├── cmd/
├── internal/
│   └── lru/
│       ├── lru.go
│       └── avl_test.go
│
├── movies.csv
├── ratings.csv
├── go.mod
├── README.md
└── performance.md
```

---

## 3. Dataset utilizado

Dataset:

```text
MovieLens Latest Small
```

Fuente:

```text
https://grouplens.org/datasets/movielens/
```

Archivos usados:

```text
movies.csv
ratings.csv
```

El archivo `movies.csv` contiene los identificadores y títulos de las películas.

El archivo `ratings.csv` contiene las calificaciones de los usuarios. A partir de estas calificaciones se calcula el promedio de rating para cada película.

---

## 4. Requisitos

* Go 1.26.4 o superior
* Windows amd64
* Archivos `movies.csv` y `ratings.csv` ubicados en la raíz del proyecto

Verificar versión de Go:

```bash
go version
```

---

## 5. Compilación

Para comprobar que el proyecto compila correctamente:

```bash
go build ./...
```

Este comando debe ejecutarse sin errores.

---

## 6. Ejecución del programa

Para ejecutar el programa principal:

```bash
go run .\internal\lru\lru.go
```

También se puede indicar un rango personalizado:

```bash
go run .\internal\lru\lru.go -a 3.5 -b 4.0
```

Ejemplo de salida esperada:

```text
Resultados en [3.50, 4.00]
4.00 | 3414 | Love Is a Many-Splendored Thing (1955)
4.00 | 107953 | Dragon Ball Z: Battle of Gods (2013)
...

Altura: 12
Nodos: 1286
Rotaciones: 888
Balanceado: true
```

---

## 7. Pruebas unitarias

Para ejecutar las pruebas unitarias:

```bash
go test ./...
```

Resultado obtenido:

```text
ok      taller3/internal/lru    0.429s
```

Las pruebas incluyen:

* Árbol vacío
* Inserción de un solo elemento
* Inserción de varios elementos
* Verificación de balance
* Consulta por rango
* Consulta sin resultados
* Inserción de varias películas con la misma clave

---

## 8. Benchmarks

Para ejecutar los benchmarks:

```bash
go test ./internal/lru -bench=Benchmark -benchmem -count=1 -v
```

Resultado obtenido:

```text
BenchmarkInsertarAVL-4              5582            223770 ns/op           96001 B/op       2000 allocs/op
BenchmarkConsultaRangoAVL-4         1743            730986 ns/op          700194 B/op       1954 allocs/op
```

Estos resultados se encuentran analizados en el archivo:

```text
performance.md
```

---

## 9. Complejidad

| Operación               | Complejidad  |
| ----------------------- | ------------ |
| Inserción AVL           | O(log n)     |
| Rotación izquierda      | O(1)         |
| Rotación derecha        | O(1)         |
| Consulta por rango      | O(log n + k) |
| Verificación de balance | O(n)         |
| Conteo de nodos         | O(n)         |

Donde `k` representa la cantidad de resultados devueltos en la consulta por rango.

---

## 10. Funciones principales

### `Insertar`

```go
func Insertar(raiz *NodoAVL, clave float64, dato Pelicula) *NodoAVL
```

Inserta una película en el árbol AVL usando como clave su rating promedio. Luego actualiza la altura del nodo y aplica rotaciones si el árbol se desbalancea.

Complejidad: O(log n)

### `rotarIzq`

```go
func rotarIzq(x *NodoAVL) *NodoAVL
```

Realiza una rotación hacia la izquierda para corregir desbalances del árbol.

Complejidad: O(1)

### `rotarDer`

```go
func rotarDer(y *NodoAVL) *NodoAVL
```

Realiza una rotación hacia la derecha para corregir desbalances del árbol.

Complejidad: O(1)

### `ConsultaRango`

```go
func ConsultaRango(raiz *NodoAVL, a, b float64) []Pelicula
```

Devuelve todas las películas cuya clave se encuentra entre `a` y `b`, en orden ascendente.

Complejidad: O(log n + k)

### `BalanceOK`

```go
func BalanceOK(n *NodoAVL) bool
```

Verifica que todos los nodos del árbol cumplan con la condición AVL:

```text
|factor de balance| <= 1
```

---

## 11. Diagramas

Los diagramas de funciones fueron elaborados para las funciones principales:

* `Insertar`
* `ConsultaRango`
* `rotarIzq`
* `rotarDer`

Estos diagramas muestran la firma de cada función, los parámetros, los retornos, la lógica interna y las llamadas entre funciones.

---

## 12. Video explicativo

Enlace del video de YouTube:

```text
PEGAR AQUÍ EL ENLACE DEL VIDEO NO LISTADO
```

El video incluye:

* Presentación del problema
* Explicación de la estructura AVL
* Explicación de inserción y rotaciones
* Explicación de consulta por rango
* Ejecución del código sobre el dataset MovieLens
* Comentario de resultados de performance

---

## 13. Integrantes

* Rodrigo Villavicencio
* Fabrizio Rodriguez
* Fernando

---

## 14. Conclusión

El proyecto implementa correctamente un árbol AVL desde cero en Go. El programa procesa datos reales del dataset MovieLens, calcula el rating promedio por película e indexa los registros en un árbol AVL.

La ejecución muestra que el árbol permanece balanceado:

```text
Balanceado: true
```

Además, los benchmarks muestran que la inserción y la consulta por rango se comportan de acuerdo con la complejidad teórica esperada.
