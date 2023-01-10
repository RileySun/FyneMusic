package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	
	"github.com/RileySun/FyneMusic/player"
	"github.com/RileySun/FyneMusic/playlist"
	"github.com/RileySun/FyneMusic/song"
	"github.com/RileySun/FyneMusic/settings"
)

var config *settings.Config
var playList *playlist.Playlist
var playerObj *player.Player
var window fyne.Window

//Init
func init() {
	fmt.Println("FyneMusic")
	
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
	window = app.NewWindow("FyneMusic")
	
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
		if (playerObj.Song.Path == playList.Songs[index].Path) {
			playerContainer := playerObj.Render()
			window.SetContent(playerContainer)
			playerObj.UpdateWidgets()
			return
		} else {
			playerObj.Close()
		}
	}
	
	//Get New Queue based off selected song (index)
	playerObj.NewQueue(playList.PlaylistPaths(), index)
	if playerObj.Queue == nil {
		panic("Queue Error")
	}
	
	playerObj.Song = song.NewSong(playerObj.Queue.Songs[playerObj.Queue.Index])

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
	//If a song has been loaded
	if playerObj.Song != nil {
		playerObj.UpdateWidgets()
	}
	window.SetContent(list)
}

func openSettings() {
	settingsPage := settings.Render()
	settings.ReturnToMenu = returnToMenu
	settings.ChangeVolume = changeVolume
	window.SetContent(settingsPage)
}

//Util
func changeVolume(v float64) {
	if playerObj.Song != nil {
		playerObj.Song.Player.SetVolume(v)
	}
}