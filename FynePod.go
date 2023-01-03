package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	
	"github.com/RileySun/FynePod/playlist"
	"github.com/RileySun/FynePod/player"
	"github.com/RileySun/FynePod/song"
	"github.com/RileySun/FynePod/settings"
)

var playList *playlist.Playlist
var config *settings.Config
var podcast *song.Song
var window fyne.Window

//Init
func init() {
	fmt.Println("FynePod")
	
	config = settings.GetSettings()
	
	playList = playlist.NewPlaylist(config.Dir)
	playList.Select = func(id int64) {selectSong(id)}
	playList.Settings = func() {openSettings()}
}

//Main
func main() {
	app := app.New()
	window = app.NewWindow("FynePod")
	
	//Settings window
	settings.ParentWindow = window
	
	list := playList.Render()
	list.Resize(fyne.NewSize(400, 600))
	
	content := list
	
	window.SetContent(content)
	
	window.CenterOnScreen()
	
	window.Resize(fyne.NewSize(400, 600))
	
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
	window.SetContent(playerContainer)
	
	//Play
	podcast.Play()
	player.UpdateWidgets()
}

func returnToMenu() {
	list := playList.Render()
	list.Resize(fyne.NewSize(400, 600))
	window.SetContent(list)
}

func openSettings() {
	settingsPage := settings.Render()
	window.SetContent(settingsPage)
}