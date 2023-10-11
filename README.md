# Attendant Zoom
![AttendantZoomLogo](https://github.com/livrasand/Attendant-Zoom/assets/104039397/e4535a92-68fb-45bb-adec-f63ff20aed72)
> ##### Si consideras útil este proyecto, apóyalo haciendo "★ Star" en el repositorio. ¡Gracias!

### 🇲🇽 Español 🌮
#### ¿Qué es Attendant Zoom?
Facilita la compartición de contenido multimedia durante reuniones, ya sea en plataformas de videoconferencias o en el Salón del Reino, ofreciendo una adaptación perfecta para encuentros híbridos, remotos o presenciales. Nuestra aplicación es completamente gratuita y no incluye publicidad.

Desarrollada en Go con Fyne, lo que permite que la aplicacion sea multiplataforma, disponible para Windows, macOS y Linux, esta aplicación se presenta como una excelente alternativa a JW Library. Está diseñada especialmente para aquellos dispositivos que no cuentan con la compatibilidad de JW Library, no tienen acceso a una tienda de aplicaciones o experimentan lentitud o bloqueos al usar JW Library en sus computadoras. Simplifica tus reuniones y estudios con nuestra herramienta de confianza.

#### ¿Cómo empiezo?
Para comenzar a utilizar esta aplicación, dirígete a la sección de Releases de este repositorio. Descarga el instalador correspondiente a tu sistema operativo (Windows, macOS o Linux), y ejecuta el instalador y sigue las instrucciones para completar la instalación de la aplicación en tu sistema.

¡Listo! Ahora deberías tener la aplicación instalada y lista para su uso.

#### ¿Qué puede hacer Attendant Zoom?
Attendant Zoom le permite descargar, sincronizar, compartir y mostrar fácil y automáticamente toda la multimedia de la reunión. Para reuniones de congregación híbridas o en persona, el modo de presentación de multimedia integrado tiene todas las características necesarias para simplificar la tarea de compartir multimedia con la congregación, que incluyen:

- Miniaturas de la multimedia.
- Proyectar imágenes y videos.
- Reproducción de audio.
- Funciones de _pausa/reproducir/detener_ fáciles de usar para gestionar la reproducción de archivos multimedia.
- Reproducción sencilla de música de fondo, con parada automática antes del inicio de las reuniones programadas periódicamente.
- Reconocimiento y gestión automática de monitores externos.
- Convertir imágenes a MP4.

En cuanto a las reuniones de Zoom de congregación totalmente remotas, la función de conversión a MP4 incorporada le permite compartir archivos multimedia de todo tipo fácilmente, utilizando la función para compartir MP4 nativa de Zoom.

En general, Attendant Zoom tiene todas las funciones respecto al departamento de audio y video de JW Library, con algunas ventajas sobresalientes en compatibilidad con Zoom, lo cual convierte a Attendant Zoom en una herramienta completa y valiosa para los Testigos de Jehová, diseñada especialmente para ayudarles en su departamento.

#### ¿Attendant Zoom funciona en mi idioma?
**¡Sí!** La multimedia para las reuniones de los testigos de Jehová se pueden descargar automáticamente en cualquiera de los miles de idiomas que están disponibles en JW.ORG. La lista de idiomas disponibles se actualiza dinámicamente. Todo lo que necesitas hacer es seleccionar cuál quieres.

¡Además, constantemente el propio Attendant Zoom se está traduciendo a varios idiomas! Por lo tanto, puede configurar el idioma que desea que se muestre en la interfaz de Attendant Zoom. ¿Quieres ayudar a traducir Attendant Zoom a tu idioma? Consulte nuestro archivo CONTRIBUTING.md para obtener instrucciones sobre cómo hacerlo.

### 🇺🇸 English 🗽
#### What is Attendant Zoom?
It facilitates the sharing of multimedia content during meetings, whether on video conferencing platforms or at the Kingdom Hall, offering a perfect adaptation for hybrid, remote or in-person meetings. Our application is completely free and does not include advertising.

Developed in Go with Fyne, which allows the application to be cross-platform, available for Windows, macOS and Linux, this application is presented as an excellent alternative to JW Library. It is designed especially for devices that do not have JW Library support, do not have access to an app store, or experience slowness or freezes when using JW Library on their computers. Simplify your meetings and studies with our trusted tool.

#### How do I start?
To start using this application, go to the Releases section of this repository. Download the installer corresponding to your operating system (Windows, macOS or Linux), and run the installer and follow the instructions to complete the installation of the application on your system.

Ready! You should now have the app installed and ready to use.

#### What can Attendant Zoom do?
Attendant Zoom allows you to easily and automatically download, sync, share and display all meeting multimedia. For in-person or hybrid congregation meetings, the integrated multimedia presentation mode has all the features needed to simplify the task of sharing multimedia with the congregation, including:

- Multimedia thumbnails.
- Project images and videos.
- Audio playback.
- Easy-to-use _pause/play/stop_ functions to manage playback of media files.
- Simple background music playback, with automatic stopping before the start of regularly scheduled meetings.
- Automatic recognition and management of external monitors.
- Convert images to MP4.

As for fully remote congregation Zoom meetings, the built-in MP4 conversion feature allows you to easily share media files of all types, using Zoom's native MP4 sharing feature.

Overall, Attendant Zoom is fully featured in JW Library's audio and video department, with some notable advantages in Zoom compatibility, making Attendant Zoom a complete and valuable tool for Jehovah's Witnesses, designed especially for help them in their department.

#### Does Attendant Zoom work in my language?
**Yes!** Multimedia for Jehovah's Witness meetings can be automatically downloaded in any of the thousands of languages that are available on JW.ORG. The list of available languages is updated dynamically. All you need to do is select which one you want.

Plus, Attendant Zoom itself is constantly being translated into multiple languages! Therefore, you can configure the language you want to be displayed on the Attendant Zoom interface. Do you want to help translate Attendant Zoom into your language? See our CONTRIBUTING.md file for instructions on how to do this.

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
