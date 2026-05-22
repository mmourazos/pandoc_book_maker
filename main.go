// Este es el fichero principal (probablemente el único) de la aplicación.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const pandocCommand = "pandoc.exe"

func wrap(s string, w string) string {
	sb := strings.Builder{}

	sb.WriteString(w)
	sb.WriteString(s)
	sb.WriteString(w)

	return sb.String()
}

func buildMarkdownFiles(files []string) string {
	sb := strings.Builder{}

	for _, file := range files {
		sb.WriteString(wrap(file, "'"))
		sb.WriteString(" ")
	}

	return sb.String()
}

func getMarkdownFiles(dir string, recursive bool, files *[]string) (*[]string, error) {
	fmt.Printf("Procesando directorio: %s\n", dir)
	if files == nil {
		slice := make([]string, 0, 128)
		files = &slice
	}

	// DirEntry es una interfaz que representa una entrada en un directorio, que puede ser un archivo o un subdirectorio.
	// Tiene los métodos Name(), IsDir(), Type()  y FileInfo() (que devuelve un FileInfo con información sobre el archivo o directorio).
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Error al leer el directorio %s: %v", dir, err)
		return nil, err
	}

	for _, entry := range entries {
		if recursive && entry.IsDir() {
			fmt.Printf("Procesando entrada de directorio: %s\n", entry.Name())
			path := filepath.Join(dir, entry.Name())
			fmt.Printf("Ruta generada para el directorio: %s\n", path)
			getMarkdownFiles(path, recursive, files)
		}

		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			path := filepath.Join(dir, entry.Name())
			fmt.Printf("Archivo Markdown encontrado: %s\n", path)

			*files = append(*files, path)
		}
	}
	return files, nil
}

func main() {
	// Usamos flag para obtener los argumentos de la línea de comandos.
	rootDir := flag.String("b", ".", "El directorio base del proyecto")
	recursive := flag.Bool("r", false, "Si se deben buscar archivos de forma recursiva")
	outputFile := flag.String("o", "output.pdf", "El nombre del archivo de salida")
	defaultsFile := flag.String("d", ".\\defaults.yaml", "El archivo de configuración de pandoc")
	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error al obtener el directorio de trabajo actual: %v", err)
	}

	fullRootDir := filepath.Join(pwd, *rootDir)
	fmt.Printf("Ruta completa del directorio base: %s\n", fullRootDir)

	fullDefaultsFile := filepath.Join(pwd, *defaultsFile)

	fmt.Printf("Directorio base: %s\n", fullRootDir)
	// Obtener los subdirectorios del directorio base.

	mdFiles, err := getMarkdownFiles(fullRootDir, *recursive, nil)
	if err != nil {
		log.Fatalf("Error al obtener los archivos Markdown: %v", err)
	}

	for _, mdFile := range *mdFiles {
		fmt.Printf("Archivo Markdown encontrado: %s\n", mdFile)
	}

	commandOpts := []string{fmt.Sprintf("-d %s", fullDefaultsFile), fmt.Sprintf("-o %s", *outputFile)}
	commandOpts = append(commandOpts, buildMarkdownFiles(*mdFiles))

	cmd := exec.Command(pandocCommand, commandOpts...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Comando a ejecutar: %s\n", cmd.String())
	if err := cmd.Run(); err != nil {
		log.Printf("Error al ejecutar el comando pandoc: %v", err)
	}
}
