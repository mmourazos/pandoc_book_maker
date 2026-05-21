# Notas

La idea es hacer un programa en Go al que, indicándole un directorio raíz, genere un _libro_ utilizando todos los ficheros markdown que se encuentren dentro de dicho directorio. Debería de incluirse una opción para que también pueda hacerlo de forma recursiva.

## Paso uno: comando

Determinar si se puede usar simplemente un comando de Pandoc para que genere un pdf a partir de múltiples fichero markdown.

El siguiente comando: 

```bash
pandoc --defaults .\defaults.yaml -o book.pdf  '.\md_test_files\00 - Servicios de red.md' '.\md_test_files\01 - Email - SMTP.md' -F mermaid-filter.cmd
```
```

Hay que tener en cuenta que estamos en Windows. El filtro `mermaid-filter` se instaló utilizando `npm install --global mermaid-filter`. Se añadió la ruta "C:/Users/USUARIO/AppData/npm" al path del usuario y accesos directos a la carpeta "C:/Users/USUARIO/Roaming/pandoc/filters".

**NOTA:** Hay que indicar como nombre del filtro `mermaid-filter.cmd` y no `mermaid-filter`. Probablemente el segundo funcione en sistemas linux.
