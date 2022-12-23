package main

import (
	"os"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	
	"github.com/RileySun/FynePod/player"
	"github.com/RileySun/FynePod/song"
)

var podcast *song.Song

func init() {
	fmt.Println("FynePod")
	dir, _ := os.Getwd()
	podcast = song.NewSong(dir + "/Intro.mp3")	
	player.StartPlayer(podcast)
}

func main() {
	app := app.New()
	window := app.NewWindow("FynePod")
	window.Resize(fyne.NewSize(400, 600))
	
	playerContainer := player.Render()
	
	content := container.New(layout.NewCenterLayout(), playerContainer)
	
	window.SetContent(container.New(layout.NewCenterLayout(), content))
	window.CenterOnScreen()
	window.ShowAndRun()
}