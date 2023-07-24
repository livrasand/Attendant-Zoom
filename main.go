package main

import (
	"flag"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

const (
	RES240      = "240p"
	RES360      = "360p"
	RES480      = "480p"
	RES720      = "720p"
	CONFIG_FILE = ".meeting-media"
	WM          = "WM"
	MM          = "MM"
)

func main() {
	config := NewConfig()
	a := app.New()

	config.DebugMode = flag.Bool("d", false, "descarga falsa; imprimir información de depuración")
	flag.Parse()
	if *config.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("FUNCIONANDO EN MODO DE DEPURACIÓN")
	}

	progressBar := widget.NewProgressBar()
	config.Progress = &progress{0, "", progressBar}
	pbFormatter := func() string { return config.Progress.Title }
	config.Progress.ProgressBar.TextFormatter = pbFormatter

	settingsTab := container.NewTabItem("", config.settingsGUI())
	settingsTab.Icon = theme.SettingsIcon()
	tabs := container.NewAppTabs(
		container.NewTabItem("Entresemana", config.mGUI(MM)),
		container.NewTabItem("Fin de semana", config.mGUI(WM)),
		settingsTab,
	)

	w := a.NewWindow("Attendant Zoom - Media Downloader")
	w.SetContent(container.NewVBox(tabs))

	w.ShowAndRun()
}
