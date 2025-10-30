package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "Herramienta CLI para procesamiento de datos",
	Long:  `Esta herramienta permite procesar, limpiar y analizar datos en CSV.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
