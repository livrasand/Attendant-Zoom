# Attendant Zoom
### 🇪🇸 Español
Facilita la compartición de contenido multimedia durante reuniones, ya sea en plataformas de videoconferencias o en el Salón del Reino, ofreciendo una adaptación perfecta para encuentros híbridos, remotos o presenciales. Nuestra aplicación es completamente gratuita y no incluye publicidad.

Desarrollada en Go con Fyne, lo que permite que la aplicacion sea multiplataforma, disponible para Windows, macOS y Linux, esta aplicación se presenta como una excelente alternativa a JW Library. Está diseñada especialmente para aquellos dispositivos que no cuentan con la compatibilidad de JW Library, no tienen acceso a una tienda de aplicaciones o experimentan lentitud o bloqueos al usar JW Library en sus computadoras. Simplifica tus reuniones y estudios con nuestra herramienta de confianza.

### 🇺🇸 English
It facilitates the sharing of multimedia content during meetings, whether on video conferencing platforms or at the Kingdom Hall, offering a perfect adaptation for hybrid, remote or in-person meetings. Our application is completely free and does not include advertising.

Developed in Go with Fyne, which allows the application to be cross-platform, available for Windows, macOS and Linux, this application is presented as an excellent alternative to JW Library. It is designed especially for devices that do not have JW Library support, do not have access to an app store, or experience slowness or freezes when using JW Library on their computers. Simplify your meetings and studies with our trusted tool.

# Cómo empezar
Para instalar todas las dependencias del proyecto, usa:
```go
go get ...
```
Si te obtienes algún error como respuesta de alguna dependencia, prueba instalándola individualmente.

Para probar la aplicación, utilice el comando `go run`:
```go
go run main.go
```

Usando `go build`, puede generar un binario ejecutable para la aplicación, lo que le permitirá implementarlo donde lo desee (Windows, macOS, Linux, Android, iOS o iPad).

Pruebe esto con `main.go`:
```go
go build
```

Si no proporciona un argumento para este comando, `go build` compilará automáticamente el programa `main.go` en su directorio actual. El comando incluirá todos sus archivos `*.go` en el directorio. También creará todo el código de soporte necesario para poder ejecutar el binario en computadoras que tengan la misma arquitectura de sistema que su computadora, independientemente de que este tenga los archivos de origen `.go` o incluso, sin una instalación de Go.
