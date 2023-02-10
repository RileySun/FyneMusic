package tray_test

import(
	"testing"
	"github.com/stretchr/testify/assert"
	
	"fyne.io/fyne/v2/app"
	
	"github.com/RileySun/FyneMusic/song"
	"github.com/RileySun/FyneMusic/tray"
	"github.com/RileySun/FyneMusic/player"
)

var newPlayer *player.Player

func createTestingTray() *tray.Tray {
	newApp := app.New()
	
	newPlayer = player.NewPlayer()
	paths := []string{"../Music/Intro.mp3", "../Music/No Meta.mp3"}
	newPlayer.NewQueue(paths, 0)
	newPlayer.Song = song.NewSong("../Music/Intro.mp3")
	_ = newPlayer.Render()
	
	defer newPlayer.Close()
	
	return tray.NewTray(newApp, newPlayer)
}

func Test_NewTray(t *testing.T) {
	newTray := createTestingTray()
	
	//Check when no song present
	assert.True(t, newTray.Current.Label == "No Song")
}

func Test_Refresh(t *testing.T) {
	newTray := createTestingTray()
	newTray.Refresh()
	
	//Check with song present
	assert.True(t, newTray.Current.Label == "Intro - RileySun")
}

func Test_Play(t *testing.T) {
	newTray := createTestingTray()
	newTray.Refresh()
	newTray.Play()
	
	//Check song is playing
	assert.True(t, !newPlayer.Song.Paused)
}