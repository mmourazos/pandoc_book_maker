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

func buildMarkdownFiles(files []string) []string {
	output := make([]string, 0, len(files))
	sb := strings.Builder{}

	for _, file := range files {
		sb.WriteString(wrap(file, "'"))
		sb.WriteString(" ")
	}

	output = append(output, sb.String())

	return output
}

func getMarkdownFiles(dir string, recursive bool, files *[]string) error {
	fmt.Printf("Procesando directorio: %s\n", dir)

	// DirEntry es una interfaz que representa una entrada en un directorio, que puede ser un archivo o un subdirectorio.
	// Tiene los métodos Name(), IsDir(), Type()  y FileInfo() (que devuelve un FileInfo con información sobre el archivo o directorio).
	entries, err := os.ReadDir(dir)
	if err != nil {
		err = fmt.Errorf("error al leer el directorio %s: %v", dir, err)
		return err
	}

	for _, entry := range entries {
		if recursive && entry.IsDir() {
			fmt.Printf("Procesando entrada de directorio: %s\n", entry.Name())
			path := filepath.Join(dir, entry.Name())
			fmt.Printf("Ruta generada para el directorio: %s\n", path)
			if err := getMarkdownFiles(path, recursive, files); err != nil {
				return err
			}
		}

		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			path := filepath.Join(dir, entry.Name())
			fmt.Printf("Archivo Markdown encontrado: %s\n", path)

			*files = append(*files, path)
		}
	}
	return nil
}

func main() {
	// Usamos flag para obtener los argumentos de la línea de comandos.
	rootDir := flag.String("b", ".", "El directorio base del proyecto")
	recursive := flag.Bool("r", false, "Si se deben buscar archivos de forma recursiva")
	outputFile := flag.String("o", "output.pdf", "El nombre del archivo de salida")
	defaultsFile := flag.String("d", ".\\defaults.yaml", "El archivo de configuración de pandoc")
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error al obtener el directorio de trabajo actual: %v", err)
	}

	fullDefaultsFile := filepath.Join(cwd, *defaultsFile)
	_, err = os.Stat(fullDefaultsFile)
	if err != nil {
		log.Fatalf("Error archivo de configuración no encontrado: %v", err)
	}

	fullDefaultsFile = filepath.ToSlash(fullDefaultsFile)

	// Aquí se construye la lista de argumentos para el comando pandoc.
	args := make([]string, 0, 64)

	fmt.Printf("Ruta completa al archivo de configuración: %s\n", fullDefaultsFile)
	args = append(args, "-d", fullDefaultsFile)

	fmt.Printf("Archivo de salida: %s\n", *outputFile)
	args = append(args, "-o", *outputFile)

	fullRootDir := filepath.Join(cwd, *rootDir)

	fmt.Printf("Directorio base: %s\n", fullRootDir)
	// Obtener los subdirectorios del directorio base.

	mdFiles := make([]string, 0, 128)
	if err := getMarkdownFiles(fullRootDir, *recursive, &mdFiles); err != nil {
		log.Fatalf("Error al obtener los archivos Markdown: %v", err)
	}

	fmt.Printf("Número total de archivos Markdown encontrados: %d\n", len(mdFiles))
	fmt.Printf("Archivos Markdown encontrados:\n")
	for _, mdFile := range mdFiles {
		fmt.Println(mdFile)
	}

	// args = append(args, buildMarkdownFiles(mdFiles)..)
	args = append(args, mdFiles...)

	// cmd := exec.Command(pandocCommand, args...)

	// Intentando usar una shell para no volverme loco.
	// singleArg := strings.Join(args, " ")
	// cmd := exec.Command("pwsh.exe", "/C", fmt.Sprintf("%s %s", pandocCommand, singleArg))
	cmd := exec.Command(pandocCommand, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Comando a ejecutar: %s\n", cmd.String())
	if err := cmd.Run(); err != nil {
		log.Printf("Error al ejecutar el comando pandoc: %v", err)
	}
}
