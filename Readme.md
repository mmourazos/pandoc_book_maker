# Notas

La idea de este proyecto es hacer un programa en Go (o cualquier otro lenguaje) al que, indicándole un directorio raíz, genere un _libro_ utilizando todos los ficheros markdown que se encuentren dentro de dicho directorio. Debería de incluirse también una opción para que también pueda hacerlo de forma recursiva.

## Requisitos

### Pandoc

En principio usaremos Pandoc para generar el pdf a partir de los ficheros markdown. Por lo tanto, el programa que se haga en Go (o cualquier otro lenguaje) simplemente se encargará de generar un comando de Pandoc con los argumentos adecuados para que se genere el pdf.

Por lo tanto será necesarios tener Pandoc instalado en el sistema. Si tenemos permisos de administrador lo instaríamos con el comando `winget install JohnMacFarlane.Pandoc`. Si no, podemos instalar [scoop](https://scoop.sh/) y luego instalar Pandoc con el comando `scoop install pandoc`.

### Distribución de LaTeX

Pandoc usará $LaTeX$ para generar el pdf, por lo que también es necesario tenerlo instalado. Para ello se puede usar [MiKTeX](https://miktex.org) o [TeX Live](https://www.tug.org/texlive/). Si estamos en windows: `winget install MiKTeX` o `winget install TeXLive`. Si no, se pueden instalar con el gestor de paquetes de la distribución que se esté usando. También se puede instalar MikTeX con scoop: `scoop install miktex`.

### Mermaid

Si usamos diagramas [mermaid](https://mermaid.js.org/) en nuestros documentos markdown, también es necesario tener instalado el filtro [`mermaid-filter`](https://github.com/raghur/mermaid-filter). Hay que tener en cuenta que el filtro no se actualiza desde hace un par de años, por lo que es posible que no funcione con los último diagramas de mermaid.

Hay que tener en cuenta que, si estamos en Windows, y el filtro se instaló utilizando `npm install --global mermaid-filter`, éste se habrá añadido en la ruta: `C:/Users/USUARIO/AppData/npm`. Esta ruta ha de añadirse al path para que Pandoc pueda encontrar el filtro cuando se invoque. Además ha de invocarse como `mermaid-filter.cmd` y no como `mermaid-filter`, que es como se haría en sistemas linux.

Alternativamente se puede copiar el filtro `mermaid-filter.cmd` en el directorio `C:/Users/USUARIO/AppData/Roaming/pandoc/filters`, que es el directorio donde Pandoc busca los filtros por defecto.

Para indicar a Pandoc que use el filtro `mermaid-filter`, simplemente hay que indicar su ruta en el comando de Pandoc con el argumento `-F mermaid-filter.cmd`. Opcionalmente se puede indicar el filtro o filtros que deseamos usar en el fichero con la configuración por defecto `defaults.yaml` para no tener que repetirlo cada vez que se quiera generar el pdf.

```yaml
filters:
  - mermaid-filter.cmd
```

### Plantilla de Pandoc

Para mejorar la calidad del pdf generado es posible utilizar palantillas de Pandoc. Estas plantillas se pueden encontrar en la [página de plantillas de Pandoc](https://pandoc-templates.org). Para usar una plantilla, simplemente hay que indicar su ruta en el comando de Pandoc con el argumento `--template`. Para facilitar que Pandoc pueda encontrar la plantilla, es recomendable copiarla en el directorio `C:/Users/USUARIO/AppData/Roaming/pandoc/templates`, que es el directorio donde Pandoc busca las plantillas por defecto.

En nuestro ejemplo, se ha usado la plantilla `eisvogel`, que se puede descargar desde la página de plantillas de Pandoc. Para usar esta plantilla, simplemente hay que indicar su ruta en el comando de Pandoc con el argumento `--template eisvogel`.

#### Variables de la plantilla

Las plantillas de Pandoc pueden tener _opciones_ para activar o desactivar ciertas características. Estas opciones se denominan _variables_ de la plantilla. Para establecer estas variables (aunque se puede hacer desde la línea de comandos con el argumento `-V`), es más cómodo incluirlas en el fichero de configuración `defaults.yaml` para no tener que repetirlas una y otra vez cada vez que se quiera generar el pdf. Para ello hemos de usar la clave `variables` en el fichero de configuración, seguida de un bloque con las variables que se desean establecer. Por ejemplo:

```yaml
variables:
  titlepage: true
  toc-own-page: true
  title: "Servicios de red"
  author: Manuel C. Piñeiro Mourazos
  date: 2026-05-04
```

## Comando Pandoc

Determinar si se puede usar simplemente un comando de Pandoc para que genere un pdf a partir de múltiples fichero markdown.

El siguiente comando:

```Powershell
pandoc --defaults .\defaults.yaml -o book.pdf  '.\md_test_files\00 - Servicios de red.md' '.\md_test_files\01 - Email - SMTP.md' -F mermaid-filter.cmd
```

**NOTA:** Hay que indicar como nombre del filtro `mermaid-filter.cmd` y no `mermaid-filter`. Probablemente el segundo funcione en sistemas linux.

### Explicación del comando

Por defecto, Pandoc ya convierte múltiples ficheros markdown en un único pdf, por lo que no es necesario usar ningún argumento especial para ello.

#### `deaults.yaml`

En primer lugar hay que aclarar que, aunque en este documento siempre hablemos de `defaults.yaml`, el nombre del fichero de configuración puede ser cualquiera, siempre y cuando se indique su ruta en el comando de Pandoc con el argumento `--defaults`. Por ejemplo, si el fichero de configuración se llama `config.yaml`, el comando de Pandoc sería:

```Powershell
pandoc --defaults .\config.yaml -o book.pdf  '.\md_test_files\00 - Servicios de red.md' '.\md_test_files\01 - Email - SMTP.md' -F mermaid-filter.cmd
```

El fichero `defaults.yaml` se usa para indicar a Pandoc que argumentos se desean aplicar ([documentación de Pandoc](https://pandoc.org/MANUAL.html#defaults-files)). Estos argumentos se podría indicar directamente en la línea de comandos pero, debido a la cantidad de comandos que vamos a usar, es más cómodo usar un fichero de configuración. Además, si se desea usar el mismo comando para generar el pdf a partir de diferentes ficheros markdown, es necesario usar un fichero de configuración para no tener que repetir los mismos argumentos una y otra vez.

En este fichero de configuración se pueden indicar en primer lugar todas las opciones que admite Pandoc, como por ejemplo la plantilla a usar, los filtros a usar, las variables de la plantilla, etc. Para más información sobre las opciones se puede consultar la [documentación de Pandoc](https://pandoc.org/MANUAL.html#options). En el fichero incluido en este proyecto se ha intentado incluir comentarios para explicar cada una de las opciones que se han usado.

## Trabajo en progreso

Por el momento hay valores que se indican en el preámbulo de los ficheros markdown pero la idea sería migrar estos valores al fichero de configuración `defaults.yaml` para no tener que repetirlos en cada uno de los ficheros markdown. Por ejemplo, el título del libro, el autor, la fecha, etc. Aún no tengo claro como pasar todos los valores (que no se apliquen específicamente a un documento) al `defaults.yaml` pero según lo vaya aclarando lo iré añadiendo a este documento.
