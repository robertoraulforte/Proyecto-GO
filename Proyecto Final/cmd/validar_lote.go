package cmd

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var archivo string

var validarLoteCmd = &cobra.Command{
	Use:   "validar-lote",
	Short: "Valida registros de un archivo CSV e inserta los válidos en una base de datos",
	Run: func(cmd *cobra.Command, args []string) {
		if archivo == "" {
			log.Fatal("Debe proporcionar un archivo CSV con --archivo")
		}

		// Configuración de logs estructurados
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logFile, err := os.OpenFile("pipeline.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logrus.SetOutput(logFile)
		} else {
			log.Fatal("No se pudo abrir pipeline.log para escritura")
		}

		// Abrir y leer el CSV
		f, err := os.Open(archivo)
		if err != nil {
			logrus.Fatal(err)
		}
		defer f.Close()

		reader := csv.NewReader(f)
		registros, err := reader.ReadAll()
		if err != nil {
			logrus.Fatal(err)
		}

		if len(registros) < 1 {
			logrus.Fatal("El archivo CSV no contiene registros")
		}

		headers := registros[0]
		_ = headers // por si luego se quiere usar
		data := registros[1:]

		// Conexión a SQLite
		db, err := sql.Open("sqlite3", "demo.db")
		if err != nil {
			logrus.Fatal(err)
		}
		defer db.Close()

		// Crear tabla si no existe
		crearTabla := `CREATE TABLE IF NOT EXISTS usuarios_limpios (
            nombre TEXT,
            apellido TEXT,
            correo TEXT,
            genero TEXT,
            ip TEXT
        )`
		_, err = db.Exec(crearTabla)
		if err != nil {
			logrus.Fatal(err)
		}

		total := 0
		validados := 0
		for _, fila := range data {
			total++
			if len(fila) != 6 {
				logrus.WithField("fila", fila).Warn("Fila incompleta")
				continue
			}

			if err := validarFila(fila); err != nil {
				logrus.WithField("fila", fila).WithError(err).Warn("Datos inválidos")
				continue
			}

			// Insertar en la base de datos
			_, err = db.Exec(`INSERT INTO usuarios_limpios (nombre, apellido, correo, genero, ip) VALUES (?, ?, ?, ?, ?)`,
				fila[1], fila[2], fila[3], fila[4], fila[5])
			if err != nil {
				logrus.WithField("fila", fila).WithError(err).Warn("Error insertando en DB")
				continue
			}

			validados++
		}

		logrus.WithFields(logrus.Fields{
			"archivo":    archivo,
			"procesados": total,
			"insertados": validados,
		}).Info("Proceso finalizado")

		fmt.Println("Proceso completado. Ver 'pipeline.log' para detalles.")
	},
}

// validarFila valida que ciertos campos no estén vacíos y que el correo sea válido
func validarFila(fila []string) error {
	if fila[1] == "" || fila[2] == "" || fila[3] == "" {
		return errors.New("campos vacíos obligatorios (nombre, apellido, correo)")
	}

	// Validación simple de correo
	re := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	if !re.MatchString(fila[3]) {
		return errors.New("correo inválido")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(validarLoteCmd)
	validarLoteCmd.Flags().StringVarP(&archivo, "archivo", "a", "", "Ruta del archivo CSV a procesar")
}
