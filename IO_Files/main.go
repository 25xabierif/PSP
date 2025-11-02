package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// --- Estructuras ---
type Producto struct {
	ID        string
	Nombre    string
	Categoria string
	Precio    float64
	Stock     int
}

func (p Producto) String() string {
	return fmt.Sprintf("ID: %s | %s | Stock actual: %d unidades", p.ID, p.Nombre, p.Stock)
}

type Transaccion struct {
	Tipo       string
	IDProducto string
	Cantidad   int
	Fecha      string
}

func (t Transaccion) String() string {
	return fmt.Sprintf("%s,%s,%d,%s", t.Tipo, t.IDProducto, t.Cantidad, t.Fecha)
}

// --- Funciones requeridas ---

// leerInventario lee el inventario desde un CSV y devuelve slice de Producto
func leerInventario(nombreArchivo string) ([]Producto, error) {
	f, err := os.Open(nombreArchivo)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir inventario: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error leyendo CSV inventario: %w", err)
	}

	var productos []Producto
	for i, rec := range records {
		// saltar cabecera
		if i == 0 {
			continue
		}
		if len(rec) < 5 {
			continue
		}
		precio, err := strconv.ParseFloat(rec[3], 64)
		if err != nil {
			return nil, fmt.Errorf("precio inválido en línea %d: %w", i+1, err)
		}
		stock, err := strconv.Atoi(rec[4])
		if err != nil {
			return nil, fmt.Errorf("stock inválido en línea %d: %w", i+1, err)
		}
		p := Producto{
			ID:        rec[0],
			Nombre:    rec[1],
			Categoria: rec[2],
			Precio:    precio,
			Stock:     stock,
		}
		productos = append(productos, p)
	}
	return productos, nil
}

// leerTransacciones lee las transacciones desde CSV y devuelve slice de Transaccion
func leerTransacciones(nombreArchivo string) ([]Transaccion, error) {
	f, err := os.Open(nombreArchivo)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir transacciones: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error leyendo CSV transacciones: %w", err)
	}

	var trans []Transaccion
	for i, rec := range records {
		// saltar cabecera
		if i == 0 {
			continue
		}
		if len(rec) < 4 {
			continue
		}
		cant, err := strconv.Atoi(rec[2])
		if err != nil {
			return nil, fmt.Errorf("cantidad inválida en línea %d: %w", i+1, err)
		}
		t := Transaccion{
			Tipo:       rec[0],
			IDProducto: rec[1],
			Cantidad:   cant,
			Fecha:      rec[3],
		}
		trans = append(trans, t)
	}
	return trans, nil
}

// procesarTransacciones procesa las transacciones sobre los productos.
// Devuelve una lista de mensajes de error (strings) que no pudieron aplicarse.
func procesarTransacciones(productos []Producto, transacciones []Transaccion) []string {
	var errores []string

	// índice rápido por ID (posición en slice)
	index := make(map[string]int)
	for i, p := range productos {
		index[p.ID] = i
	}

	for _, t := range transacciones {
		pos, existe := index[t.IDProducto]
		if !existe {
			errores = append(errores, fmt.Sprintf("ERROR: Producto %s no encontrado en transacción de tipo %s (fecha: %s)", t.IDProducto, t.Tipo, t.Fecha))
			continue
		}

		switch t.Tipo {
		case "VENTA":
			if productos[pos].Stock < t.Cantidad {
				errores = append(errores, fmt.Sprintf("ERROR: Stock insuficiente para venta. Producto: %s, Stock actual: %d, Cantidad solicitada: %d (fecha: %s)",
					productos[pos].ID, productos[pos].Stock, t.Cantidad, t.Fecha))
				continue
			}
			productos[pos].Stock -= t.Cantidad
		case "COMPRA":
			productos[pos].Stock += t.Cantidad
		case "DEVOLUCION":
			productos[pos].Stock += t.Cantidad
		default:
			errores = append(errores, fmt.Sprintf("ERROR: Tipo de transacción desconocido '%s' para producto %s (fecha: %s)", t.Tipo, t.IDProducto, t.Fecha))
		}
	}

	// Nota: la slice 'productos' se modifica in-place porque es una referencia a backing array,
	// pero si quieres devolverla puedes hacerlo. En este enunciado solo pedían errores.
	return errores
}

// escribirInventario escribe el inventario en formato CSV al archivo indicado.
func escribirInventario(productos []Producto, nombreArchivo string) error {
	f, err := os.Create(nombreArchivo)
	if err != nil {
		return fmt.Errorf("no se pudo crear archivo inventario actualizado: %w", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// cabecera
	if err := writer.Write([]string{"ID", "Nombre", "Categoría", "Precio", "Stock"}); err != nil {
		return fmt.Errorf("error escribiendo cabecera: %w", err)
	}

	for _, p := range productos {
		rec := []string{
			p.ID,
			p.Nombre,
			p.Categoria,
			fmt.Sprintf("%.2f", p.Precio),
			strconv.Itoa(p.Stock),
		}
		if err := writer.Write(rec); err != nil {
			return fmt.Errorf("error escribiendo registro: %w", err)
		}
	}
	return nil
}

// generarReporteBajoStock escribe un archivo con los productos cuyo stock < limite.
func generarReporteBajoStock(productos []Producto, limite int) error {
	nombre := "productos_bajo_stock.txt"
	f, err := os.Create(nombre)
	if err != nil {
		return fmt.Errorf("no se pudo crear reporte bajo stock: %w", err)
	}
	defer f.Close()

	// Header
	fmt.Fprintln(f, "ALERTA: PRODUCTOS CON BAJO STOCK")
	fmt.Fprintln(f, "================================")

	count := 0
	for _, p := range productos {
		if p.Stock < limite {
			fmt.Fprintf(f, "ID: %s | %s | Stock actual: %d unidades\n", p.ID, p.Nombre, p.Stock)
			count++
		}
	}
	fmt.Fprintf(f, "Total de productos con bajo stock: %d\n", count)
	return nil
}

// escribirLog escribe los errores en un archivo de log usando log.SetOutput y también devuelve error si falla
func escribirLog(errores []string, nombreArchivo string) error {
	// asegurarnos de crear directorio si hace falta
	dir := filepath.Dir(nombreArchivo)
	if dir != "." {
		_ = os.MkdirAll(dir, 0755)
	}

	f, err := os.Create(nombreArchivo)
	if err != nil {
		return fmt.Errorf("no se pudo crear archivo de log: %w", err)
	}
	defer f.Close()

	// redirigimos log a ese fichero
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags) // incluir timestamp

	for _, e := range errores {
		// añadimos prefijo ERROR para coherencia con enunciado
		log.Println("[ERROR]:", e)
	}

	return nil
}

// --- main: orquestador ---
func main() {
	invFile := "inventario.txt"
	transFile := "transacciones.txt"

	// 1) leer inventario
	productos, err := leerInventario(invFile)
	if err != nil {
		fmt.Println("Fallo leyendo inventario:", err)
		return
	}

	// 2) leer transacciones
	transacciones, err := leerTransacciones(transFile)
	if err != nil {
		fmt.Println("Fallo leyendo transacciones:", err)
		return
	}

	// 3) procesar transacciones
	errores := procesarTransacciones(productos, transacciones)

	// 4) escribir inventario actualizado
	if err := escribirInventario(productos, "inventario_actualizado.txt"); err != nil {
		fmt.Println("Error escribiendo inventario actualizado:", err)
		return
	}

	// 5) generar reporte bajo stock (limite 10)
	if err := generarReporteBajoStock(productos, 10); err != nil {
		fmt.Println("Error generando reporte bajo stock:", err)
		return
	}

	// 6) escribir log de errores
	if err := escribirLog(errores, "errores.log"); err != nil {
		fmt.Println("Error escribiendo log de errores:", err)
		return
	}

	// resumen por consola
	fmt.Println("Proceso completado.")
	fmt.Println("Inventario actualizado -> inventario_actualizado.txt")
	fmt.Println("Reporte bajo stock -> productos_bajo_stock.txt")
	fmt.Println("Errores (si los hay) -> errores.log")
}
