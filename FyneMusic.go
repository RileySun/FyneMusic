package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	
	"github.com/RileySun/FyneMusic/song"
	"github.com/RileySun/FyneMusic/tray"
	"github.com/RileySun/FyneMusic/utils"
	"github.com/RileySun/FyneMusic/player"
	"github.com/RileySun/FyneMusic/playlist"
	"github.com/RileySun/FyneMusic/settings"
)

var fyneApp fyne.App
var window fyne.Window
var config *settings.Config
var playList *playlist.Playlist
var playerObj *player.Player
var systemTray *tray.Tray

//Init
func init() {
	fmt.Println("FyneMusic")
}

func setup() {
	fyneApp = app.NewWithID("com.sunshine.fynemusic")
	window = fyneApp.NewWindow("FyneMusic")
	
	//Config(need fyneApp & window for Prefrences api, and copying to clipboard)
	settings.LoadSettings(fyneApp, window)
	config = settings.GetSettings()
	
	//Theme
	fyneApp.Settings().SetTheme(&utils.NewTheme{})
	
	//Player Module (needed for Playlist, must come first)
	playerObj = player.NewPlayer()
	playerObj.ReturnToMenu = func() {returnToMenu()}
	playerObj.ResumeSong = func() {resumeSong()}
	playerObj.RefreshTray = func() {refreshTray()}
	
	//Playlist Module
	playList = playlist.NewPlaylist(config.Dir, playerObj)
	playList.Select = func(id int64) {selectSong(id)}
	playList.Settings = func() {openSettings()}
	
	//System Tray Notifications (desktop only, no fyne tray for mobile)
	device := fyne.CurrentDevice()
	if !device.IsMobile() {
		systemTray = tray.NewTray(fyneApp, playerObj)	
	}
}

//Main
func main() {
	setup()
	
	//Settings window
	settings.ParentWindow = window
	
	list := playList.Render()
	list.Resize(fyne.NewSize(400, 600))
	
	content := list
	
	window.SetContent(content)
	
	window.CenterOnScreen()
	
	window.Resize(fyne.NewSize(400, 600))
	window.SetFixedSize(true)
	
	window.SetMaster()
	
	//If setup not done, show settings page
	if !config.Setup {
		openSettings()
	}
	
	window.ShowAndRun()
}

//Change Tabs
func selectSong(index int64) {	
	//If there is a song already (close old, or if same open player)
	if playerObj.Song != nil {
		//If selecting same song, open player, if not, close song
		if (playerObj.Song.Path == playList.Songs[index].Path) {
			resumeSong()
			return
		} else {
			playerObj.Close()
		}
	}
	
	//Get New Queue based off selected song (index)
	playerObj.NewQueue(playList.PlaylistPaths(), index)
	if playerObj.Queue == nil {
		log.Fatal("Failed to create player queue")
	}
	
	playerObj.Song = song.NewSong(playerObj.Queue.Songs[playerObj.Queue.Index])

	//Render Player
	playerContainer := playerObj.Render()
	window.SetContent(playerContainer)
	
	//Play Selected Song
	playerObj.Song.Play()
	playerObj.UpdateWidgets()
	refreshTray()
}

func resumeSong() {
	playerContainer := playerObj.Render()
	window.SetContent(playerContainer)
	playerObj.UpdateWidgets()
}

func returnToMenu() {
	playList.Songs = playlist.GetSongs(config.Dir)
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

func refreshTray() {
	//Only if on desktop where tray is supported
	//Used in player module to refresh system tray on player changes
	if systemTray != nil {
		systemTray.Refresh()
	}
}