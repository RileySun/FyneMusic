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


type Player struct {
	Song *song.Song
	PlayButton *playbutton.PlayButton
	ReturnToMenu func()
}

//Declarations
var podcast *song.Song
var playButton *playbutton.PlayButton
var slider *track.TrackSlider

//Create
func NewPlayer() *Player {
	player := new(Player)
	
	return player
}

//Render
func (p *Player) Render() *fyne.Container {
	//BackButton
	backButton := widget.NewButtonWithIcon("", theme. MenuIcon(), func() {p.ReturnToMenu()})
	backSpacer := layout.NewSpacer()
	backContainer := container.New(layout.NewHBoxLayout(), backSpacer, backButton)
	
	//Spacers
	spacerTop := layout.NewSpacer()
	spacerBottom := layout.NewSpacer()
	
	//Meta
	artwork, title := p.CreateMeta()
	
	//Slider
	slider = track.NewTrack(p.Song)
	
	//Buttons
	prevButton  := widget.NewButtonWithIcon("Prev", theme.MediaSkipPreviousIcon(), func() {p.Prev()})
	rewindButton := widget.NewButtonWithIcon("Rewind", theme.MediaFastRewindIcon(), func() {p.Rewind()})
	forwardButton := widget.NewButtonWithIcon("Forward", theme. MediaFastForwardIcon(), func() {p.Forward()})
	nextButton := widget.NewButtonWithIcon("Next", theme. MediaSkipNextIcon(), func() {p.Next()})
	p.PlayButton = playbutton.NewPlayButton(p.Song)
	
	//Containers
	sliderContainer := container.New(layout.NewMaxLayout(), slider)
	buttonContainer := container.New(layout.NewHBoxLayout(), prevButton, rewindButton, p.PlayButton, forwardButton, nextButton)
	playerContainer := container.New(layout.NewVBoxLayout(), backContainer, spacerTop, artwork, spacerBottom, title, sliderContainer, buttonContainer)
	
	return playerContainer
}

func (p *Player) RenderMini() *fyne.Container {
	prevButton  := widget.NewButtonWithIcon("Prev", theme.MediaSkipPreviousIcon(), func() {p.Prev()})
	rewindButton := widget.NewButtonWithIcon("Rewind", theme.MediaFastRewindIcon(), func() {p.Rewind()})
	forwardButton := widget.NewButtonWithIcon("Forward", theme. MediaFastForwardIcon(), func() {p.Forward()})
	nextButton := widget.NewButtonWithIcon("Next", theme. MediaSkipNextIcon(), func() {p.Next()})
	
	p.PlayButton = playbutton.NewPlayButton(p.Song)
	buttonContainer := container.New(layout.NewHBoxLayout(), prevButton, rewindButton, p.PlayButton, forwardButton, nextButton)
	
	return buttonContainer
}

//Utils
func (p *Player) UpdateWidgets() {
	slider.SetTime()
	p.PlayButton.UpdateState()
}

func (p *Player) CreateMeta() (*fyne.Container, *widget.Label) {
	var art *canvas.Image
	
	if len(p.Song.Meta.Image) != 0 {
		reader := bytes.NewReader(p.Song.Meta.Image)
		art = canvas.NewImageFromReader(reader, p.Song.Meta.Title)
	} else {
		dir, _ := os.Getwd()
		art = canvas.NewImageFromFile(dir + "/Default.jpg")
	}
	art.FillMode = canvas.ImageFillOriginal
	artContainer := container.New(layout.NewMaxLayout(), art)
	
	var titleString string
	if p.Song.Meta.Title != "" {
		//Add Tag Title, and Artist if available
		titleString = p.Song.Meta.Title
		if p.Song.Meta.Artist != "" {
			titleString += " - " + p.Song.Meta.Artist
		}
	} else {
		titleString = p.Song.Meta.File
	}
	title := widget.NewLabel(titleString)
	title.Alignment = 1
	
	return artContainer, title
}

//Buttons
func (p *Player) Prev() {
	if p.Song != nil {
		p.Song.Restart()
		p.UpdateWidgets()
	}
}

func (p *Player) Rewind() {
	if p.Song != nil {
		p.Song.Rewind()
		p.UpdateWidgets()
	}
}

func (p *Player) Forward() {
	if p.Song != nil {
		p.Song.Forward()
		p.UpdateWidgets()
	}
}

func (p *Player) Next() {
	if p.Song != nil {
		slider.SetTime()
		p.UpdateWidgets()
	}
}