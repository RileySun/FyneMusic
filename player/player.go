package player

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/data/binding"
	
	//"fmt"
	
	"time"
	"strconv"
	
	"github.com/RileySun/FynePod/song"
)

var podcast *song.Song
var timeUpdate chan bool

var timeStr binding.String
var currentTime *widget.Label
var maxTimeStr binding.String
var maxTime *widget.Label

var sliderFloat binding.Float
var slider *widget.Slider

//Init
func StartPlayer(current *song.Song) {
	podcast = current
	
	timeStr = binding.NewString()
	timeStr.Set(strconv.Itoa(int(podcast.Current)))
	
	maxTimeStr = binding.NewString()
	maxTimeStr.Set(strconv.Itoa(int(podcast.Length)))
	
	sliderFloat = binding.NewFloat()
	sliderFloat.Set(float64(podcast.Current))
}

//Create
func Render() *fyne.Container {
	//Slider
	//currentTime = widget.NewLabelWithData(timeStr)
	slider = widget.NewSliderWithData(0, float64(podcast.Length), sliderFloat)
	slider.OnChanged = func(f float64) {dragSlider(f)}
	//maxTime = widget.NewLabelWithData(maxTimeStr)
	
	//Buttons
	prevButton  := widget.NewButtonWithIcon("Prev", theme.MediaSkipPreviousIcon(), func() {prev()})
	rewindButton := widget.NewButtonWithIcon("Rewind", theme.MediaFastRewindIcon(), func() {rewind()})
	playButton := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {play()})
	//pauseButton := widget.NewButtonWithIcon("Pause", theme.MediaPauseIcon(), func() {pause()})
	forwardButton := widget.NewButtonWithIcon("Forward", theme. MediaFastForwardIcon(), func() {forward()})
	nextButton := widget.NewButtonWithIcon("Next", theme. MediaSkipNextIcon(), func() {next()})
	
	
	sliderContainer := container.New(layout.NewMaxLayout(), slider)
	//sliderContainer := container.New(layout.NewHBoxLayout(), currentTime, slideContainer, maxTime)
	buttonContainer := container.New(layout.NewHBoxLayout(), prevButton, rewindButton, playButton, forwardButton, nextButton)
	playerContainer := container.New(layout.NewVBoxLayout(), sliderContainer, buttonContainer)
	
	updateTime()
	
	return playerContainer
}

//Util
func updateTime() {
	go func() {
		for {
			select {
				case <- timeUpdate:
					return
				default:
					timeStr.Set(strconv.Itoa(int(podcast.Current)))
					sliderFloat.Set(float64(podcast.Current))
					time.Sleep(time.Second)
	   		 }
		}
	}()
}

//Slider
func dragSlider(newValue float64) {
	//podcast.Seek(int64(newValue))
}

//Buttons

func prev() {
	sliderFloat.Set(0.0)
	podcast.Restart()
}

func rewind() {
	podcast.Rewind()
	sliderFloat.Set(float64(podcast.Current))
}

func play() {
	podcast.Play()
}

func pause() {

}

func forward() {
	podcast.Forward()
	sliderFloat.Set(float64(podcast.Current))
}

func next() {

}

