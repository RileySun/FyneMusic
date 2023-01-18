package player

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	
	"bytes"
	"os"
	"time"
	"image/color"
	"math/rand"
	
	"github.com/RileySun/FyneMusic/song"
	"github.com/RileySun/FyneMusic/icons"
	"github.com/RileySun/FyneMusic/track"
	"github.com/RileySun/FyneMusic/playbutton"
)

type Player struct {
	Song *song.Song
	Queue *Queue
	Shuffle bool
	Repeat bool
	stopUpdating chan bool
	PlayButton *playbutton.PlayButton
	ReturnToMenu func()
	ResumeSong func()
	
	artContainer *fyne.Container
	title *widget.Label
	slider *track.TrackSlider
	sliderContainer *fyne.Container
	miniContainer *fyne.Container
	shuffleButton *widget.Button
	repeatButton *widget.Button
}

type Queue struct {
	Songs []string
	Original []string
	Index int64
	Length int64
}

	//Create
//Player
func NewPlayer() *Player {
	player := new(Player)
	
	player.Queue = new(Queue)
	player.Shuffle = false
	player.Repeat = false
	player.stopUpdating  = make(chan bool, 100)
	player.startUpdate()
	
	rand.Seed(time.Now().UnixNano())
	
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
	queue.Original = queue.Songs
	
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
	//Repeat and Shuffle
	p.shuffleButton = widget.NewButtonWithIcon("", icons.Shuffle, func() {p.ShuffleState()})
	p.repeatButton = widget.NewButtonWithIcon("", icons.Repeat, func() {p.RepeatState()})
	if p.Shuffle {
		p.shuffleButton.Importance = widget.HighImportance
	} else {
		p.shuffleButton.Importance = widget.MediumImportance
	} // These HAVE to be here for some reason or they dont work.
	if p.Repeat {
		p.repeatButton.Importance = widget.HighImportance
	} else {
		p.repeatButton.Importance = widget.MediumImportance
	}

	//BackButton
	backButton := widget.NewButtonWithIcon("", icons.Menu, func() {p.ReturnToMenu()})
	backSpacer := layout.NewSpacer()
	backContainer := container.New(layout.NewHBoxLayout(), backSpacer, p.shuffleButton, p.repeatButton, backButton)
	
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
	p.sliderContainer = container.New(layout.NewMaxLayout(), p.slider)
	
	//Buttons
	prevButton := widget.NewButtonWithIcon("Prev", icons.Prev, func() {p.Prev()})
	rewindButton := widget.NewButtonWithIcon("Rewind", icons.Rewind, func() {p.Rewind()})
	forwardButton := widget.NewButtonWithIcon("Forward", icons.Forward, func() {p.Forward()})
	nextButton := widget.NewButtonWithIcon("Next", icons.Next, func() {p.Next()})
	p.PlayButton = playbutton.NewPlayButton(p.Song)
	buttonContainer := container.New(layout.NewHBoxLayout(), prevButton, rewindButton, p.PlayButton, forwardButton, nextButton)
	
	//Output
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
	//Get Meta Mini
	p.miniContainer = container.New(layout.NewMaxLayout())
	if p.Song != nil {
		artwork, title := p.GetRawMeta()
		miniButton := widget.NewButtonWithIcon(title, artwork, p.ResumeSong)
		p.miniContainer.Add(miniButton)
	}
	topBorder := canvas.NewLine(color.NRGBA{R: 155, G: 155, B: 155, A: 255})
	topBorder.StrokeWidth = 2
	topContainer := container.NewBorder(topBorder, nil, nil, nil, p.miniContainer)

	//Controls
	prevButton := widget.NewButtonWithIcon("Prev", icons.Prev, func() {p.Prev()})
	rewindButton := widget.NewButtonWithIcon("Rewind", icons.Rewind, func() {p.Rewind()})
	forwardButton := widget.NewButtonWithIcon("Forward", icons.Forward, func() {p.Forward()})
	nextButton := widget.NewButtonWithIcon("Next", icons.Next, func() {p.Next()})
	
	p.PlayButton = playbutton.NewPlayButton(p.Song)
	buttonContainer := container.New(layout.NewHBoxLayout(), prevButton, rewindButton, p.PlayButton, forwardButton, nextButton)
	
	finalContainer := container.New(layout.NewVBoxLayout(), topContainer, buttonContainer)
	
	return finalContainer
}

func (p *Player) RenderMiniUpdate() {
	//Meta
	artwork, title := p.GetRawMeta()
	miniButton := widget.NewButtonWithIcon(title, artwork, p.ResumeSong)
	p.miniContainer.RemoveAll()
	p.miniContainer.Add(miniButton)
	p.PlayButton.UpdateState()
}

//Queue
func (p *Player) RenderQueue() {//*fyneContainer {
	
}

	//Utils
//Player
func (p *Player) UpdateWidgets() {
	//Slider and PlayButton
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
	if p.Song.Meta.Title != " " {
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

func (p *Player) GetRawMeta() (fyne.Resource, string) {
	//Art
	var art *canvas.Image
	if len(p.Song.Meta.Image) != 0 {
		reader := bytes.NewReader(p.Song.Meta.Image)
		art = canvas.NewImageFromReader(reader, p.Song.Meta.Title)
	} else {
		dir, _ := os.Getwd()
		art = canvas.NewImageFromFile(dir + "/Default.jpg")
	}
	artByte := art.Resource
	
	//Title
	var titleString string
	if p.Song.Meta.Title != " " {
		//Add Tag Title, and Artist if available
		titleString = p.Song.Meta.Title
		if p.Song.Meta.Artist != "" {
			titleString += " - " + p.Song.Meta.Artist
		}
	} else {
		titleString = p.Song.Meta.File
	}
	
	return artByte, titleString
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
	p.Song.Play()
	p.UpdateWidgets()
	p.RenderUpdate()
	p.RenderMiniUpdate()
}

func (p *Player) GetQueueNext() {
	//if queue exists
	if p.Queue.Length != 0 {
		//If repeat is on
		if p.Repeat {
			newSongPath := p.Queue.Songs[p.Queue.Index]
			p.newQueueSong(newSongPath)
		} else {		
			if p.Queue.Index < p.Queue.Length - 1 {
				p.Queue.Index++
			} else {
				p.Queue.Index = 0
			}
			newSongPath := p.Queue.Songs[p.Queue.Index]
			p.newQueueSong(newSongPath)
		}
	}
}

func (p *Player) GetQueuePrev() {
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

func (p *Player) ShuffleQueue() {
	currentSong := p.Queue.Songs[p.Queue.Index]
	songs := p.Queue.Songs
	rand.Shuffle(len(songs), func(a, b int) {
		songs[a], songs[b] = songs[b], songs[a]
	})
	//get new index
	for i, v := range songs {
		if v == currentSong {
			p.Queue.Index = int64(i)
			p.Queue.Songs = songs
		}
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

func (p *Player) ShuffleState() {
	if p.Shuffle {
		p.shuffleButton.Importance = widget.MediumImportance
		p.Queue.Songs = p.Queue.Original
		p.Shuffle = false
	} else {
		p.shuffleButton.Importance = widget.HighImportance
		p.ShuffleQueue()
		p.Shuffle = true
	}
}

func (p *Player) RepeatState() {
	if p.Repeat {
		p.repeatButton.Importance = widget.MediumImportance
		p.Repeat = false
	} else {
		p.repeatButton.Importance = widget.HighImportance
		p.Repeat = true
	}
}