package main

import (
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"

	"io/ioutil"
	"path/filepath"

	"fyne.io/fyne/v2/canvas"
	"os"
    "github.com/hajimehoshi/go-mp3"
    "github.com/hajimehoshi/oto"
	"log"
	"fmt"
)

var currentSelectedFilePath string

func playAudioMP3(filePath string) {
	// Abre el archivo MP3 para lectura
	file, err := os.Open(filePath)
	if err != nil {
		errorMessage := fmt.Sprintf("Error al abrir el archivo MP3 '%s': %s", filePath, err)
		logError(errorMessage)
		return
	}
	defer file.Close()

	// Decodifica el archivo MP3
	mp3Decoder, err := mp3.NewDecoder(file)
	if err != nil {
		errorMessage := fmt.Sprintf("Error al decodificar el archivo MP3 '%s': %s", filePath, err)
		logError(errorMessage)
		return
	}

	// Prepara el contexto de audio Oto
	otoContext, err := oto.NewContext(44100, 2, 2, 8192)
	if err != nil {
		errorMessage := fmt.Sprintf("Error al crear el contexto de audio Oto: %s", err)
		logError(errorMessage)
		return
	}
	defer otoContext.Close()

	// Crea un nuevo reproductor de audio
	player := otoContext.NewPlayer()
	defer player.Close()

	// Lee y reproduce el audio MP3 en un goroutine independiente
	go func() {
		buffer := make([]byte, 8192)
		for {
			_, err := mp3Decoder.Read(buffer)
			if err != nil {
				break
			}
			player.Write(buffer)
		}
	}()
}

// La función isMP3File debe determinar si un archivo es un archivo MP3
func isMP3File(fileName string) bool {
    ext := strings.ToLower(filepath.Ext(fileName))
    return ext == ".mp3"
}

func (c *Config) mGUI(m string) *fyne.Container {

	date := widget.NewEntry()
	date.SetText(time.Now().Format("2006-01-02"))

	song1box := widget.NewEntry()
	song1box.SetPlaceHolder("Canción #1")
	song2box := widget.NewEntry()
	song2box.SetPlaceHolder("Canción #2")
	song3box := widget.NewEntry()
	song3box.SetPlaceHolder("Canción #3")

	fetchOtherMedia := widget.NewCheck("Obtener otros medios (imágenes y videos)", func(f bool) {
		c.FetchOtherMedia = f
		c.writeConfigToFile()
	})
	fetchOtherMedia.SetChecked(c.FetchOtherMedia)

	if c.AutoFetchMeetingData {
		if m == MM {
			song1box.Disabled()
		}
		song2box.Disabled()
		song3box.Disabled()
		fetchOtherMedia.Enable()
	} else {
		if m == MM {
			song1box.Enable()
		}
		song2box.Enable()
		song3box.Enable()
		fetchOtherMedia.Disable()
	}

	autoFetchMeetingData := widget.NewCheck("Obtener automáticamente los datos de la reunión", func(f bool) {
		c.AutoFetchMeetingData = f
		c.writeConfigToFile()
		if f {
			if m == MM {
				song1box.Disable()
			}
			song2box.Disable()
			song3box.Disable()
			fetchOtherMedia.Enable()
		} else {
			if m == MM {
				song1box.Enable()
			}
			song2box.Enable()
			song3box.Enable()
			fetchOtherMedia.Disable()
		}
	})
	autoFetchMeetingData.SetChecked(c.AutoFetchMeetingData)

	playlistOption := widget.NewCheck("Crear lista de reproducción para usar con VLC", func(p bool) {
		c.CreatePlaylist = p
		c.writeConfigToFile()
	})
	playlistOption.SetChecked(c.CreatePlaylist)

	fetchButton := widget.NewButton("Buscar", func() {
		dateToSet, err := time.Parse("2006-01-02", date.Text)
		if err != nil {
			logrus.Fatal(err)
		}
		c.Date = WeekOf(dateToSet)
		c.SongsToGet = []string{song1box.Text, song2box.Text, song3box.Text}

		if err := c.fetchMeetingStuff(m); err == nil {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Attendant Zoom",
				Content: "Descarga completada",
			})
		} else {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Attendant Zoom",
				Content: "Error en la descarga...",
			})
		}

		c.Pictures = []file{}
		c.Videos = []video{}
		c.SongsToGet = []string{}
		c.SongsNames = []string{}
	})

	mmBox := container.NewVBox(
		date,
		autoFetchMeetingData,
		fetchOtherMedia,
		song1box,
		song2box,
		song3box,
		playlistOption,
		fetchButton,
		c.Progress.ProgressBar,
	)

	return mmBox
}

func (c *Config) settingsGUI() *fyne.Container {
	resPicker := widget.NewRadioGroup([]string{
		RES240,
		RES360,
		RES480,
		RES720,
	}, func(res string) {
		c.Resolution = res
	})
	resPicker.SetSelected(c.Resolution)

	targetDir := widget.NewEntry()
	targetDir.SetPlaceHolder("Ruta de descarga...")
	targetDir.SetText(c.SaveLocation)

	cacheDir := widget.NewEntry()
	cacheDir.SetPlaceHolder("Ruta de caché...")
	cacheDir.SetText(c.CacheLocation)

	purgeDir := widget.NewCheck("Eliminar contenido anterior antes de descargar nuevo", func(d bool) {
		c.PurgeSaveDir = d
	})
	purgeDir.SetChecked(c.PurgeSaveDir)

	lang := widget.NewEntry()
	lang.SetPlaceHolder("Idioma del contenido multimedia (ej. S para Español)")
	lang.SetText(c.Language)

	pubs := widget.NewEntry()
	pubs.SetPlaceHolder("Símbolos de publicaciones para el contenido multimedia (ej. th, lff)")
	var pubSymbolString string
	for i, s := range c.PubSymbols {
		if i != 0 {
			pubSymbolString += ", "
		}
		pubSymbolString += s
	}
	pubs.SetText(pubSymbolString)

	save := widget.NewButton("Guardar", func() {
		c.SaveLocation = targetDir.Text
		c.CacheLocation = cacheDir.Text
		c.Language = lang.Text
		var pubSymbolSlice []string
		for _, p := range strings.Split(pubs.Text, ",") {
			pubSymbolSlice = append(pubSymbolSlice, strings.TrimSpace(strings.ToLower(p)))
		}
		c.PubSymbols = pubSymbolSlice
		c.writeConfigToFile()
	})

	settingsBox := container.NewVBox(
		resPicker,
		targetDir,
		cacheDir,
		purgeDir,
		lang,
		pubs,
		save,
	)

	return settingsBox
}

func (c *Config) createDownloadedFilesView(mediaviewer fyne.Window) *fyne.Container {
    downloadedFolderPath := "C:\\GoAttendant\\Attendant Zoom\\meetings"

	currentSelectedFilePath := ""

    files, err := ioutil.ReadDir(downloadedFolderPath)
    if err != nil {
        logrus.Warn(err)
    }

    var fileNames []string
    for _, file := range files {
        fileNames = append(fileNames, file.Name())
    }

    fileList := widget.NewList(
        func() int {
            return len(fileNames)
        },
        func() fyne.CanvasObject {
            label := widget.NewLabel("")
            return label
        },
        func(i widget.ListItemID, obj fyne.CanvasObject) {
            label := obj.(*widget.Label)
            label.SetText(fileNames[i])
        },
    )

	fileList.OnSelected = func(id widget.ListItemID) {
		logrus.Infof("OnSelected controlador activado para el elemento %d", id)
	
		if id >= 0 && int(id) < len(fileNames) {
			selectedFileName := fileNames[id]
			newSelectedFilePath := filepath.Join(downloadedFolderPath, selectedFileName)
	
			if isImageFile(selectedFileName) {
				// Limpia las rutas antes de compararlas
				newSelectedFilePath = filepath.Clean(newSelectedFilePath)
				currentSelectedFilePath = filepath.Clean(currentSelectedFilePath)
	
				// Agrega mensajes de registro para verificar las rutas antes de la comparación
				logrus.Infof("newSelectedFilePath antes de la comparación: %s", newSelectedFilePath)
				logrus.Infof("currentSelectedFilePath antes de la comparación: %s")
	
				if newSelectedFilePath == currentSelectedFilePath {
					// Si la imagen seleccionada es la misma que la actual, quítala
					currentSelectedFilePath = ""
					backgroundImage := canvas.NewImageFromFile("resources/yeartext.png")
					backgroundImage.Resize(fyne.NewSize(640, 360))

					containerv := container.NewMax(backgroundImage)
					containerv.Resize(fyne.NewSize(640, 360))

					mediaviewer.SetContent(containerv)
	
					// Deselecciona el elemento
					fileList.Unselect(id)
	
					logrus.Infof("Ocultar imagen seleccionada")
				} else {
					// Si la imagen seleccionada es diferente, muéstrala y actualiza la selección actual
					currentSelectedFilePath = newSelectedFilePath
					setImageInView(mediaviewer, currentSelectedFilePath)
	
					// Deselecciona el elemento
					fileList.Unselect(id)
	
					logrus.Infof("Mostrar nueva imagen seleccionada")
				}
			} else {
				if isMP3File(selectedFileName) {
					// Reproduce el archivo MP3
					playAudioMP3(newSelectedFilePath)
				}
				
			}
		}
	}
	
	
	fileListContainer := container.NewScroll(fileList)

    viewContainer := container.NewMax(
        fileListContainer,
    )

    return viewContainer
}


func setImageInView(mediaviewer fyne.Window, imagePath string) {
	 image := canvas.NewImageFromFile(imagePath)
	 mediaviewer.SetContent(container.NewMax(image))

    logrus.Infof("currentSelectedFilePath antes de la actualización: %s", currentSelectedFilePath)
    
    currentSelectedFilePath = imagePath
    
    logrus.Infof("currentSelectedFilePath después de la actualización: %s", currentSelectedFilePath)
}

func isImageFile(fileName string) bool {
    ext := strings.ToLower(filepath.Ext(fileName))
    return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".bmp" || ext == ".webp"
}

//func isVideoFile(fileName string) bool {
//    ext := strings.ToLower(filepath.Ext(fileName))
//    return ext == ".mp4" || ext == ".avi" || ext == ".mkv" || ext == ".mov"
//}

func logError(errorMessage string) {
    // Configura la ubicación de la carpeta de logs
    logFolder := "logs"

    // Crea la carpeta de logs si no existe
    if _, err := os.Stat(logFolder); os.IsNotExist(err) {
        os.Mkdir(logFolder, os.ModeDir)
    }

    // Genera el nombre del archivo de registro con fecha y hora
    logFileName := filepath.Join(logFolder, time.Now().Format("2006-01-02_15-04-05")+".log")

    // Abre el archivo de registro
    logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        fmt.Println("Error al abrir el archivo de registro:", err)
        return
    }
    defer logFile.Close()

    // Escribe el mensaje de error en el archivo de registro
    log.SetOutput(logFile)
    log.Println(errorMessage)
}
