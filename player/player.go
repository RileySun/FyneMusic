package player

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	
	"github.com/RileySun/FynePod/song"
)

var podcast song.Song

func Render(current song.Song) *fyne.Container {
	podcast = current
	
	prevButton  := widget.NewButtonWithIcon("Prev", theme.MediaSkipPreviousIcon(), func() {prev()})
	rewindButton := widget.NewButtonWithIcon("Rewind", theme.MediaFastRewindIcon(), func() {rewind()})
	playButton := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {play()})
	//pauseButton := widget.NewButtonWithIcon("Pause", theme.MediaPauseIcon(), func() {pause()})
	forwardButton := widget.NewButtonWithIcon("Forward", theme. MediaFastForwardIcon(), func() {forward()})
	nextButton := widget.NewButtonWithIcon("Next", theme. MediaSkipNextIcon(), func() {next()})
	
	buttonContainer := container.New(layout.NewHBoxLayout(), prevButton, rewindButton, playButton, forwardButton, nextButton)
	return buttonContainer
}

func prev() {
	podcast.Restart()
}

func rewind() {
	podcast.Rewind()
}

func play() {
	podcast.Play()
}

func pause() {

}

func forward() {
	podcast.Forward()
}

func next() {

}

