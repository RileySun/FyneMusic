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

var config *settings.Config
var playList *playlist.Playlist
var playerObj *player.Player
var window fyne.Window

//Init
func init() {
	fmt.Println("FynePod")
	
	config = settings.GetSettings()
	
	//Player Module (needed for Playlist, must come first)
	playerObj = player.NewPlayer()
	playerObj.ReturnToMenu = func() {returnToMenu()}
	
	//Playlist Module
	playList = playlist.NewPlaylist(config.Dir, playerObj)
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
	//If there is a song already (close old, or if same open player)
	if playerObj.Song != nil {
		//If selecting same song, open player, if not, close song
		if (playerObj.Song.Path == playList.Songs[index]) {
			playerContainer := playerObj.Render()
			window.SetContent(playerContainer)
			playerObj.UpdateWidgets()
			return
		} else {
			playerObj.Song.Close()
		}
	}
	
	//Get New Song from playlist, assign to player
	selected := playList.Songs[index]
	playerObj.Song = song.NewSong(selected)

	//Render Player
	playerContainer := playerObj.Render()
	window.SetContent(playerContainer)
	
	//Play Selected Song
	playerObj.Song.Play()
	playerObj.UpdateWidgets()
}

func returnToMenu() {
	list := playList.Render()
	list.Resize(fyne.NewSize(400, 600))
	playerObj.UpdateWidgets()
	window.SetContent(list)
}

func openSettings() {
	settingsPage := settings.Render()
	window.SetContent(settingsPage)
}