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
	"time"
	
	"github.com/RileySun/FynePod/song"
	"github.com/RileySun/FynePod/track"
	"github.com/RileySun/FynePod/playbutton"
)

type Player struct {
	Song *song.Song
	Queue *Queue
	stopUpdating chan bool
	PlayButton *playbutton.PlayButton
	ReturnToMenu func()
	
	artContainer *fyne.Container
	title *widget.Label
	slider *track.TrackSlider
	sliderContainer *fyne.Container
}

type Queue struct {
	Songs []string
	Index int64
	Length int64
}

	//Create
//Player
func NewPlayer() *Player {
	player := new(Player)
	
	player.Queue = new(Queue)
	player.stopUpdating  = make(chan bool, 100)
	player.startUpdate()
	
	return player
}

//Queue
func(p *Player) NewQueue(songList []string, index int64) {
	queue := new(Queue)
	queue.Index = 0//Not the same thing as index^ (selected songList index)
	queue.Length = int64(len(songList))
	
	//Re-order playlist starting with selected index, then add remaining on to the end
	newOrder := songList[index:]
	
	//If not first selected, Add rest of the songs
	if index > 0 {
		for i, _ := range songList {
			if i < int(index) {
				newOrder = append(newOrder, songList[i])
			}
		}
	}
	
	//Get Music Paths
	for _, songItem := range newOrder {
		queue.Songs = append(queue.Songs, songItem)
	}
	
	p.Queue = queue
}


	//Update
func (p *Player) startUpdate() {
	go func() {
		for {
			select {
				case <- p.stopUpdating:
					return
				default:
					if p.Song != nil {
						if p.Song.IsEnded() {
							p.GetQueueNext()
						}
					}
					time.Sleep(time.Second/2)
	   		 }
		}
	}()
}

func (p *Player) endUpdate() {
	p.stopUpdating <- true
}

	//Render
//Player
func (p *Player) Render() *fyne.Container {
	//BackButton
	backButton := widget.NewButtonWithIcon("", theme. MenuIcon(), func() {p.ReturnToMenu()})
	backSpacer := layout.NewSpacer()
	backContainer := container.New(layout.NewHBoxLayout(), backSpacer, backButton)
	
	//Spacers
	spacerTop := layout.NewSpacer()
	spacerBottom := layout.NewSpacer()
	
	//Meta
	artwork, title := p.GetMeta()
	p.artContainer = container.New(layout.NewMaxLayout(), artwork)
	
	p.title = widget.NewLabel(title)
	p.title.Alignment = 1
	
	//Slider
	p.slider = track.NewTrack(p.Song)
	
	//Buttons
	prevButton  := widget.NewButtonWithIcon("Prev", theme.MediaSkipPreviousIcon(), func() {p.Prev()})
	rewindButton := widget.NewButtonWithIcon("Rewind", theme.MediaFastRewindIcon(), func() {p.Rewind()})
	forwardButton := widget.NewButtonWithIcon("Forward", theme. MediaFastForwardIcon(), func() {p.Forward()})
	nextButton := widget.NewButtonWithIcon("Next", theme. MediaSkipNextIcon(), func() {p.Next()})
	p.PlayButton = playbutton.NewPlayButton(p.Song)
	
	//Containers
	p.sliderContainer = container.New(layout.NewMaxLayout(), p.slider)
	buttonContainer := container.New(layout.NewHBoxLayout(), prevButton, rewindButton, p.PlayButton, forwardButton, nextButton)
	playerContainer := container.New(layout.NewVBoxLayout(), backContainer, spacerTop, p.artContainer, spacerBottom, p.title, p.sliderContainer, buttonContainer)
	
	return playerContainer
}

func (p *Player) RenderUpdate() {
	//Meta
	artwork, title := p.GetMeta()
	p.title.SetText(title)
	p.artContainer.RemoveAll()
	p.artContainer.Add(artwork)
	p.artContainer.Refresh()
	
	//Slider
	p.slider = track.NewTrack(p.Song)
	p.sliderContainer.RemoveAll()
	p.sliderContainer.Add(p.slider)
	p.sliderContainer.Refresh()
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

//Queue
func (p *Player) RenderQueue() {//*fyneContainer {
	
}

	//Utils
//Player
func (p *Player) UpdateWidgets() {
	p.slider.SetTime()
	p.PlayButton.UpdateState()
}

func (p *Player) GetMeta() (*canvas.Image, string) {
	//Artwork
	var art *canvas.Image
	
	if len(p.Song.Meta.Image) != 0 {
		reader := bytes.NewReader(p.Song.Meta.Image)
		art = canvas.NewImageFromReader(reader, p.Song.Meta.Title)
	} else {
		dir, _ := os.Getwd()
		art = canvas.NewImageFromFile(dir + "/Default.jpg")
	}
	art.FillMode = canvas.ImageFillOriginal

	//Title
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
	
	return art, titleString
}

func (p *Player) Close() {
	p.Song.Close()
	p.slider.Close()
	p.stopUpdating <- true
}

//Queue
func (p *Player) newQueueSong(path string) {
	p.Song.Close()
	p.Song = song.NewSong(path)
	p.PlayButton.Song = p.Song
	p.UpdateWidgets()
	p.RenderUpdate()
	p.Song.Play()
}

func (p *Player) GetQueueNext() {
	//if queue exists
	if p.Queue.Length != 0 {
		if p.Queue.Index < p.Queue.Length - 1 {
			p.Queue.Index++
		} else {
			p.Queue.Index = 0
		}
		newSongPath := p.Queue.Songs[p.Queue.Index]
		p.newQueueSong(newSongPath)
	}
}

func (p *Player) GetQueuePrev(){
	if p.Queue.Length != 0 {
		if p.Queue.Index > 0 {
			p.Queue.Index--
		} else {
			p.Queue.Index = p.Queue.Length - 1
		}
		newSongPath := p.Queue.Songs[p.Queue.Index]
		p.newQueueSong(newSongPath)
	}
}

//Actions
func (p *Player) Prev() {
	if p.Song != nil {
		p.GetQueuePrev()
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
		p.GetQueueNext()
	}
}