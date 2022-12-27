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
)

var podcast *song.Song


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
	playButton := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {play()})
	//pauseButton := widget.NewButtonWithIcon("Pause", theme.MediaPauseIcon(), func() {pause()})
	forwardButton := widget.NewButtonWithIcon("Forward", theme. MediaFastForwardIcon(), func() {forward()})
	nextButton := widget.NewButtonWithIcon("Next", theme. MediaSkipNextIcon(), func() {next()})
	
	sliderContainer := container.New(layout.NewMaxLayout(), slider)
	buttonContainer := container.New(layout.NewHBoxLayout(), prevButton, rewindButton, playButton, forwardButton, nextButton)
	playerContainer := container.New(layout.NewVBoxLayout(), sliderContainer, buttonContainer)
	
	return playerContainer
}

//Buttons

func prev() {
	//sliderFloat.Set(0.0)
	podcast.Restart()
}

func rewind() {
	podcast.Rewind()
	//sliderFloat.Set(float64(podcast.Current))
}

func play() {
	podcast.Play()
}

func pause() {

}

func forward() {
	podcast.Forward()
	//sliderFloat.Set(float64(podcast.Current))
}

func next() {

}