package main

import (
	"os"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	
	"github.com/RileySun/FynePod/playlist"
	"github.com/RileySun/FynePod/player"
	"github.com/RileySun/FynePod/song"
)

var playList *playlist.Playlist
var podcast *song.Song
var window fyne.Window

//Init
func init() {
	fmt.Println("FynePod")
	dir, _ := os.Getwd()
	
	playList = playlist.NewPlaylist(dir + "/Music")
	playList.Select = func(id int64) {selectSong(id)}
}

//Main
func main() {
	app := app.New()
	window = app.NewWindow("FynePod")
	window.Resize(fyne.NewSize(400, 600))
	
	list := playList.Render()
	list.Resize(fyne.NewSize(400, 600))
	
	content := list
	
	window.SetContent(content)
	window.CenterOnScreen()
	window.ShowAndRun()
}

//Change Tabs
func selectSong(index int64) {
	//Close old song, if there is one
	if podcast != nil {
		podcast.Close()
	}
	
	//Get New Song from playlist
	selected := playList.Songs[index]
	podcast = song.NewSong(selected)	
	player.StartPlayer(podcast)

	playerContainer := player.Render(returnToMenu)
	content := container.New(layout.NewCenterLayout(), playerContainer)
	window.SetContent(container.New(layout.NewCenterLayout(), content))
	
	//Play
	podcast.Play()
	player.UpdateWidgets()
}

func returnToMenu() {
	list := playList.Render()
	list.Resize(fyne.NewSize(400, 600))
	window.SetContent(list)
}