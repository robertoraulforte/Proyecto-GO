package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	_ "github.com/mattn/go-sqlite3"
)

var analizarCmd = &cobra.Command{
	Use:   "analizar-db",
	Short: "Realiza un análisis exploratorio de los datos almacenados",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Análisis Exploratorio de Datos")
		fmt.Println("========================================")

		db, err := sql.Open("sqlite3", "demo.db")
		if err != nil {
			log.Fatal("Error al abrir la base de datos:", err)
		}
		defer db.Close()

		var total int
		err = db.QueryRow("SELECT COUNT(*) FROM usuarios_limpios").Scan(&total)
		if err != nil {
			log.Fatal("Error al contar registros:", err)
		}
		fmt.Printf("Total de registros: %d\n\n", total)

		fmt.Println("Distribución por género:")
		rowsGenero, err := db.Query("SELECT genero, COUNT(*) FROM usuarios_limpios GROUP BY genero ORDER BY COUNT(*) DESC")
		if err != nil {
			log.Fatalf("Error al obtener géneros: %v", err)
		}
		defer rowsGenero.Close()
		for rowsGenero.Next() {
			var genero string
			var count int
			rowsGenero.Scan(&genero, &count)
			fmt.Printf("  - %s: %d\n", genero, count)
		}

		fmt.Println("Nombres más comunes:")
		rowsNombres, err := db.Query("SELECT nombre, COUNT(*) FROM usuarios_limpios GROUP BY nombre ORDER BY COUNT(*) DESC LIMIT 10")
		if err != nil {
			log.Fatalf("Error al obtener nombres: %v", err)
		}
		defer rowsNombres.Close()
		for rowsNombres.Next() {
			var nombre string
			var count int
			rowsNombres.Scan(&nombre, &count)
			fmt.Printf("  - %s: %d\n", nombre, count)
		}

		fmt.Println("Dominios de correo más comunes:")
		rowsCorreo, err := db.Query("SELECT correo FROM usuarios_limpios")
		if err != nil {
			log.Fatalf("Error al obtener correos: %v", err)
		}
		defer rowsCorreo.Close()

		contadores := make(map[string]int)
		for rowsCorreo.Next() {
			var correo string
			rowsCorreo.Scan(&correo)
			parts := strings.Split(correo, "@")
			if len(parts) == 2 {
				dominio := parts[1]
				contadores[dominio]++
			}
		}

		for dominio, count := range contadores {
			if count > 5 {
				fmt.Printf("  - %s: %d\n", dominio, count)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(analizarCmd)
}
