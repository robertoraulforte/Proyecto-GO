package cmd

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var consultarDbCmd = &cobra.Command{
	Use:   "consultar-db",
	Short: "Muestra los registros válidos almacenados en la base de datos",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", "demo.db")
		if err != nil {
			log.Fatalf("Error al conectar con la base de datos: %v", err)
		}
		defer db.Close()

		rows, err := db.Query("SELECT nombre, apellido, correo, genero, ip FROM usuarios_limpios")
		if err != nil {
			log.Fatalf("Error al ejecutar la consulta: %v", err)
		}
		defer rows.Close()

		fmt.Println("Registros válidos almacenados:")

		var nombre, apellido, correo, genero, ip string
		for rows.Next() {
			err := rows.Scan(&nombre, &apellido, &correo, &genero, &ip)
			if err != nil {
				log.Println("Error al leer fila:", err)
				continue
			}
			fmt.Printf("- %s %s | %s | %s | IP: %s\n", nombre, apellido, correo, genero, ip)
		}

		err = rows.Err()
		if err != nil {
			log.Fatalf("Error al recorrer los resultados: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(consultarDbCmd)
}
