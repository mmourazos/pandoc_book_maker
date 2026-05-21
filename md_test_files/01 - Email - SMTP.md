---
title: "Email usando net/smtp"
author: Manuel C. Piñeiro Mourazos
date: 2026-05-04
geometry: margin=3cm
fontsize: 11pt
mainfont: "Noto Serif"

# mainfont: "Auge-Trial"

# monofont: "Hack"

monofont: "Hack Nerd Font"

monofontoptions:
  - Scale=0.9

lang: es-ES

header-includes:

- \usepackage{tcolorbox}
- \newtcolorbox{myquote}{colback=red!5!white, colframe=red!75!black}
- \renewenvironment{quote}{\begin{myquote}}{\end{myquote}}

# Variables de la template eisvogel
toc-own-page: true

output:
  pdf_document:
    highlight: tango 
    # highlight: zenburn 
    # highlight: breezedark

    figure_caption: true
    toc: true
    number_sections: true

    latex_engine: xelatex
---

# Paquete `net.smtp`

El paquete `net/smtp` es una biblioteca estándar de Go que proporciona funciones para enviar correos electrónicos utilizando el protocolo SMTP. Permite a los desarrolladores enviar correos electrónicos desde sus aplicaciones Go de manera sencilla y eficiente. El paquete `net/smtp` ofrece una interfaz simple para configurar el servidor SMTP, autenticar usuarios y enviar mensajes de correo electrónico con diferentes formatos, como texto plano o HTML.

En los ejemplos de uso de servidores de email veréis que la mayoría dispone de sus propias bibliotecas para enviar correos electrónicos, como en el caso de [jetmail](https://www.mailjet.com/) o [mailersend](https://www.mailersend.com/).

## Ejemplo de uso de `net/smtp`

En este ejemplo veremos como enviar un email usando el servidor SMTP de Google (`host: smtp.gmail.com`) utilizando el paquete `net/smtp` de Go. Para ello, necesitaremos configurar el servidor SMTP de Gmail, autenticar nuestra cuenta de Gmail y luego enviar un correo electrónico a un destinatario específico.

Para poder utilizar Gmail para enviar un correo electrónico hemos de obtener una **app password** para luego utilizarla para autenticar nuestra cuenta de Gmail en el servidor SMTP, y así poder enviar correos electrónicos desde nuestra aplicación Go. A continuación, se detallan los pasos para generar una contraseña de aplicación y configurar el servidor SMTP de Gmail.

#### Generar la contraseña de aplicación

Una contraseña de aplicación es una contraseña específica que se utiliza para autenticar aplicaciones de terceros con tu cuenta de Google. Es un método prácticamente obsoleto, ya que Google ha deshabilitado el acceso a aplicaciones menos seguras, pero aún es posible generar una contraseña de aplicación para ciertas aplicaciones que no admiten la autenticación de dos factores. Este es nuestro caso con `net/smtp`.

Para generar una contraseña de aplicación, seguiremos los siguientes pasos:

1. Hemos de acceder a la página de seguridad de tu cuenta de Google: <https://myaccount.google.com/security>
2. En la sección "Iniciar sesión en Google", haz clic en "Contraseñas de aplicaciones".
3. Si se te solicita, ingresa tu contraseña de Google para verificar tu identidad.
4. Finalmente, selecciona la aplicación para la que deseas generar la contraseña de aplicación (en este caso, "Correo") y el dispositivo para el que deseas generar la contraseña (puedes seleccionar "Otro" y escribir un nombre personalizado). Luego, haz clic en "Generar" para obtener la contraseña de aplicación.

**Ten en cuenta que has de copiar la contraseña de aplicación generada, ya que no podrás volver a verla después de cerrar la ventana. Esta contraseña de aplicación es la que utilizarás para autenticar tu cuenta de Gmail en el servidor SMTP y enviar correos electrónicos desde tu aplicación Go.**

> Si no aparece la opción "Contraseñas de aplicaciones", es posible que tu cuenta de Google no tenga habilitada la autenticación de dos factores, lo que es necesario para generar contraseñas de aplicación. En ese caso, deberás habilitar la autenticación de dos factores en tu cuenta de Google antes de poder generar una contraseña de aplicación.
>
>Si, teniendo activada la autenticación de dos factores, no aparece la opción "Contraseñas de aplicaciones", se puede acceder directamente a la página de generación de contraseñas de aplicación a través del siguiente enlace: [https://myaccount.google.com/apppasswords](https://myaccount.google.com/apppasswords)
>En algunas cuentas, aún teniendo activada la autenticación de dos factores, no aparece la opción "Contraseñas de aplicaciones". En ese caso, se puede acceder directamente a la página de generación de contraseñas de aplicación a través del siguiente enlace: [https://myaccount.google.com/apppasswords](https://myaccount.google.com/apppasswords)

Una vez obtenida la contraseña no está de más recordar que, para mantener la seguridad de tu cuenta de Google, es importante no compartir esta contraseña de aplicación con nadie y no incluirla directamente en el código fuente de tu aplicación (_hardcodear_). En su lugar, es recomendable utilizar variables de entorno o un archivo de configuración para almacenar la contraseña de aplicación de forma segura.

#### ¿Cómo almacenar las credenciales de forma segura?

En primer lugar es importante mencionar que, si utilizamos un fichero para almacenar las credenciales, hemos de asegurarnos de que este fichero **no se incluya en el control de versiones** (incluirlo en el fichero `.gitignore`) para evitar que se suban al repositorio y se expongan _a todo el mundo_. Por ejemplo, si utilizamos un archivo `.env` para almacenar las credenciales, debemos agregar la siguiente línea al archivo `.gitignore`:

```gitignore
.env
```

Otra opción es excluir todos los archivos en el fichero `.gitignore` e incluir solo los archivos que queremos subir al repositorio:

```gitignore
# Ignore everything
*

# But not these files...
!/.gitignore

!*.go
!go.sum
!go.mod

!README.md
!LICENSE
```

La idea de utilizar este fichero `.env` se debe a que en Go disponemos de un paquete llamado [`godotenv`](https://github.com/joho/godotenv) que permite cargar variables de entorno desde dicho archivo. Esto que facilita la gestión de credenciales y otros parámetros de configuración sin exponerlos en el código fuente. Para usar este método, simplemente hemos de crear un archivo `.env` en el mismo directorio que nuestro código fuente con el siguiente contenido:

```env
VAR_1=valor_1
VAR_2=valor_2
```

Como podemos ver se trata de parejas clave-valor, donde cada línea representa una variable de entorno y su valor correspondiente. En este ejemplo, hemos definido dos variables de entorno: `VAR_1` con el valor `valor_1` y `VAR_2` con el valor `valor_2`. Podemos incluir tantas variables de entorno como necesitemos para configurar nuestra aplicación.

Para recuperar el valor de estas variables de entorno en nuestro código Go, hemos de invocar la función `godotenv.Load()` al inicio de nuestra función `main` para luego acceder a los valores de las variables de entorno, como cualquier otra variable de entorno, utilizando la función
`os.Getenv` del paquete `os`. Por ejemplo, para obtener el valor de `VAR_1` y `VAR_2`, podemos hacer lo siguiente:

```go
import (
  "log"
  "os"

// Incluir el paquete `godotenv` en nuestro código Go para cargar las variables de entorno desde el archivo `.env`:
  "github.com/joho/godotenv"
)
```

A continuación, podemos cargar las variables de entorno al inicio de nuestra función `main`:

```go
func main() {
  // Cargar las variables de entorno desde el archivo .env
  err := godotenv.Load()
  if err != nil {
    log.Fatalf("Error loading .env file: %v", err)
  }

  // Ahora podemos acceder a las variables de entorno utilizando os.Getenv
  var1 := os.Getenv("VAR_1")
  var2 := os.Getenv("VAR_2")

  // Resto del código...
}
```

## Envío de emails usando `net/smtp` con el servidor SMTP de Gmail

Ahora que hemos obtenido la contraseña de aplicación, podemos configurar el servidor SMTP de Gmail en nuestro código Go utilizando el paquete `net/smtp`. Para ello, necesitaremos especificar:

* La dirección del servidor SMTP de Gmail: `smtp.gmail.com`
* El puerto que se utilizará para la comunicación (587 para TLS o 465 para SSL).
* Las credenciales de autenticación:
   - La dirección de correo electrónico del remitente (nuestra cuenta de Gmail).
   - La contraseña de la aplicación (obtenida en el paso anterior).

```go
package main

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

const (
	server      = "smtp.gmail.com"
	port        = 587
	appPassFile = "app.pass" // Nombre del archivo que contiene la contraseña de aplicación.
	crlf        = "\r\n"
)

func main() {
	// Hemos de obtener el password de la aplicación desde el entorno o, para este ejemplo, desde un fichero.
	appPass, err := os.ReadFile(appPassFile)
	if err != nil {
		// Si no disponemos del password, no podemos continuar.
		log.Fatalf("Error reading app password: %v", err)
	}

	// Ahora hemos de establecer los valores necesarios para la conexión con el servidor SMTP de Gmail.
	// El email del remitente, el destinatario, el asunto y el cuerpo del mensaje.
	from := "asincrono@gmail.com"
	to := "mmourazos@gmail.com"
	subject := "Correo de prueba desde Go"
	// El cuerpo del mensaje.
	body := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

	// Creamos un buffer ya que al final lo que le hemos de pasar al servidor SMTP es una secuencia de bytes con el mensaje completo, incluyendo los encabezados.
	msg := bytes.Buffer{}
	fmt.Fprintf(&msg, "From: %s", from)
	msg.WriteString(crlf)
	fmt.Fprintf(&msg, "To: %s", to)
	msg.WriteString(crlf)
	fmt.Fprintf(&msg, "Subject: %s", subject)
	msg.WriteString(crlf)
	msg.WriteString(crlf) // Línea en blanco entre los encabezados y el cuerpo del mensaje.
	msg.WriteString(body)

  // Como en el caso del envío de mensajes http a través de un socket tcp, cada línea del mensaje debe terminar con un CRLF (carriage return + line feed) para que el servidor SMTP pueda interpretar correctamente el mensaje.
  // También valdría terminar en "\n" pero es recomendable usar `\r\n` para asegurar la compatibilidad con todos los servidores SMTP.

	address := fmt.Sprintf("%s:%d", server, port)
	err = smtp.SendMail(address, smtp.PlainAuth("", from, string(appPass), server), from, []string{to}, msg.Bytes())
	if err != nil {
		log.Printf("Error sending email: %v", err)
	}
}
```

En el siguiente documento veremos cómo utilizar el paquete [go-mail](https://github.com/wneessen/go-mail) para realizar el envío de correos electrónicos de una forma más sencilla y con más funcionalidades que el paquete `net/smtp`. Este paquete es una biblioteca de terceros que proporciona una interfaz más amigable y fácil de usar para enviar correos electrónicos desde aplicaciones Go, permitiendo la configuración de encabezados, adjuntos, formatos de mensaje, entre otras características avanzadas.

