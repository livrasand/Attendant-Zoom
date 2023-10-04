# Attendant-Zoom-Downloader-Media
Una aplicación multiplataforma para descargar multimedia para las reuniones, un Plugin para Attendant Zoom.

Especial agradecimiento a <a href="https://github.com/jomast">@jomast</a> por el código fuente.

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
