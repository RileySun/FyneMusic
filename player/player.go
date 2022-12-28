package player

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	
	//"fmt"
	
	"github.com/RileySun/FynePod/song"
	"github.com/RileySun/FynePod/track"
	"github.com/RileySun/FynePod/playbutton"
)

var podcast *song.Song
var playButton *playbutton.PlayButton

//Init
func StartPlayer(current *song.Song) {
	podcast = current
}

//Create
func Render() *fyne.Container {
	//Slider
	slider := track.NewTrack(podcast)
	
	//Buttons
	prevButton  := widget.NewButtonWithIcon("Prev", theme.MediaSkipPreviousIcon(), func() {prev()})
	rewindButton := widget.NewButtonWithIcon("Rewind", theme.MediaFastRewindIcon(), func() {rewind()})
	forwardButton := widget.NewButtonWithIcon("Forward", theme. MediaFastForwardIcon(), func() {forward()})
	nextButton := widget.NewButtonWithIcon("Next", theme. MediaSkipNextIcon(), func() {next()})
	
	playButton = playbutton.NewPlayButton(podcast)
	
	//Containers
	sliderContainer := container.New(layout.NewMaxLayout(), slider)
	buttonContainer := container.New(layout.NewHBoxLayout(), prevButton, rewindButton, playButton, forwardButton, nextButton)
	playerContainer := container.New(layout.NewVBoxLayout(), sliderContainer, buttonContainer)
	
	return playerContainer
}

//Util
func UpdateWidgets() {
	track.SetTime()
	playButton.UpdateState()
}

//Buttons
func prev() {
	podcast.Restart()
	UpdateWidgets()
}

func rewind() {
	podcast.Rewind()
	UpdateWidgets()
}

func forward() {
	podcast.Forward()
	UpdateWidgets()
}

func next() {
	track.SetTime()
	UpdateWidgets()
}