package cmd

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

var formato string
var salida string

// Tipo de dato global accesible desde todas las funciones
type Registro struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Correo   string `json:"correo"`
	IP       string `json:"ip"`
	Genero   string `json:"genero"`
}

var exportarCmd = &cobra.Command{
	Use:   "exportar",
	Short: "Exporta los datos limpios de la base de datos a CSV, JSON o Excel",
	Run: func(cmd *cobra.Command, args []string) {
		if formato == "" || salida == "" {
			log.Fatal("Debe indicar el formato (--formato=csv|json|xlsx) y el archivo de salida con --salida")
		}

		db, err := sql.Open("sqlite3", "demo.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := db.Query("SELECT nombre, apellido, correo, ip, genero FROM usuarios_limpios")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var datos []Registro
		for rows.Next() {
			var r Registro
			if err := rows.Scan(&r.Nombre, &r.Apellido, &r.Correo, &r.IP, &r.Genero); err != nil {
				log.Println("Error al leer fila:", err)
				continue
			}
			datos = append(datos, r)
		}

		switch strings.ToLower(formato) {
		case "csv":
			exportarCSV(datos, salida)
		case "json":
			exportarJSON(datos, salida)
		case "xlsx":
			exportarExcel(datos, salida)
		default:
			log.Fatalf("Formato no soportado: %s", formato)
		}

		fmt.Println("Datos exportados correctamente a", salida)
	},
}

func exportarCSV(data []Registro, archivo string) {
	f, err := os.Create(archivo)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	writer.Write([]string{"Nombre", "Apellido", "Correo", "IP", "Genero"})
	for _, r := range data {
		writer.Write([]string{r.Nombre, r.Apellido, r.Correo, r.IP, r.Genero})
	}
}

func exportarJSON(data []Registro, archivo string) {
	f, err := os.Create(archivo)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Fatal(err)
	}
}

func exportarExcel(data []Registro, archivo string) {
	f := excelize.NewFile()
	sheet := "Datos"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Nombre", "Apellido", "Correo", "IP", "Genero"}
	for i, h := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, h)
	}

	for i, r := range data {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", i+2), r.Nombre)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", i+2), r.Apellido)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", i+2), r.Correo)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", i+2), r.IP)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", i+2), r.Genero)
	}

	if err := f.SaveAs(archivo); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(exportarCmd)
	exportarCmd.Flags().StringVarP(&formato, "formato", "f", "", "Formato de exportación: csv, json, xlsx")
	exportarCmd.Flags().StringVarP(&salida, "salida", "o", "", "Ruta del archivo de salida")
}
