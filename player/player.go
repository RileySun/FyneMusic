package player

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	
	"bytes"
	"os"
	
	"github.com/RileySun/FynePod/song"
	"github.com/RileySun/FynePod/track"
	"github.com/RileySun/FynePod/playbutton"
)

//Declarations
var podcast *song.Song
var playButton *playbutton.PlayButton

//Init
func StartPlayer(current *song.Song) {
	podcast = current
}

//Create
func Render(menu func()) *fyne.Container {
	//BackButton
	backButton := widget.NewButtonWithIcon("", theme. MenuIcon(), func() {menu()})
	backContainer := container.New(layout.NewHBoxLayout(), backButton)
	
	//Meta
	artwork, title := CreateMeta()
	
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
	playerContainer := container.New(layout.NewVBoxLayout(), backContainer, artwork, title, sliderContainer, buttonContainer)
	
	return playerContainer
}

//Util
func UpdateWidgets() {
	track.SetTime()
	playButton.UpdateState()
}

func CreateMeta() (*canvas.Image, *widget.Label) {
	var art *canvas.Image
	
	if len(podcast.Meta.Image) != 0 {
		reader := bytes.NewReader(podcast.Meta.Image)
		art = canvas.NewImageFromReader(reader, podcast.Meta.Title)
	} else {
		dir, _ := os.Getwd()
		art = canvas.NewImageFromFile(dir + "/Default.jpg")
	}
	art.FillMode = canvas.ImageFillOriginal
	
	var titleString string
	if podcast.Meta.Title != "" {
		//Add Tag Title, and Artist if available
		titleString = podcast.Meta.Title
		if podcast.Meta.Artist != "" {
			titleString += " - " + podcast.Meta.Artist
		}
	} else {
		titleString = podcast.Meta.File
	}
	title := widget.NewLabel(titleString)
	title.Alignment = 1
	
	return art, title
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