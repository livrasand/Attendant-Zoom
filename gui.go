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

	"github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/glib"
    "log"
    "os"
)

var currentSelectedFilePath string

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
		// Agrega un mensaje de registro para verificar si se activa el controlador
		logrus.Infof("OnSelected controlador activado para el elemento %d", id)
	
		if id >= 0 && int(id) < len(fileNames) {
			selectedFileName := fileNames[id]
			newSelectedFilePath := filepath.Join(downloadedFolderPath, selectedFileName)
	
			// Verifica si la extensión del archivo corresponde a una imagen
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
					initialLabel := widget.NewLabel("Selecciona una imagen")
					mediaviewer.SetContent(container.NewMax(initialLabel))
	
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
				logrus.Infof("Este archivo no es una imagen, puedes manejarlo de manera diferente o mostrar un mensaje de error")
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

func isVideoFile(fileName string) bool {
    ext := strings.ToLower(filepath.Ext(fileName))
    return ext == ".mp4" || ext == ".avi" || ext == ".mkv" || ext == ".mov"
}