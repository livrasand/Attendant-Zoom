package main

import (
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

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

	playlistOption := widget.NewCheck("Crear lista de reproducción", func(p bool) {
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

		// reset in case of subsequent runs
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
	lang.SetPlaceHolder("Símbolo de lenguaje MEPS (ej. S)")
	lang.SetText(c.Language)

	pubs := widget.NewEntry()
	pubs.SetPlaceHolder("Símbolos de publicación vinculados para permitir (ej. th, lff)")
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
