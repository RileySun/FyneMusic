package player

import (
	//"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/data/binding"
	
	"time"
	"strconv"
	
	"github.com/RileySun/FynePod/song"
)

var podcast *song.Song
var timeStr binding.String
var currentTime widget.Label
var timeUpdate chan bool

func StartPlayer(current *song.Song) {
	podcast = current
	
	timeStr = binding.NewString()
	timeStr.Set(strconv.Itoa(int(podcast.Current)))
}

func Render() *fyne.Container {
	currentTime := widget.NewLabelWithData(timeStr)
	
	prevButton  := widget.NewButtonWithIcon("Prev", theme.MediaSkipPreviousIcon(), func() {prev()})
	rewindButton := widget.NewButtonWithIcon("Rewind", theme.MediaFastRewindIcon(), func() {rewind()})
	playButton := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {play()})
	//pauseButton := widget.NewButtonWithIcon("Pause", theme.MediaPauseIcon(), func() {pause()})
	forwardButton := widget.NewButtonWithIcon("Forward", theme. MediaFastForwardIcon(), func() {forward()})
	nextButton := widget.NewButtonWithIcon("Next", theme. MediaSkipNextIcon(), func() {next()})
	
	buttonContainer := container.New(layout.NewHBoxLayout(), currentTime, prevButton, rewindButton, playButton, forwardButton, nextButton)
	
	updateTime()
	
	return buttonContainer
}

func updateTime() {
	go func() {
		for {
			select {
				case <- timeUpdate:
					return
				default:
					timeStr.Set(strconv.Itoa(int(podcast.Current)))
					time.Sleep(time.Second)
	   		 }
		}
	}()
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

