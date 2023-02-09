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
	"math"
	"strconv"
	
	"github.com/RileySun/FyneMusic/song"
	"github.com/RileySun/FyneMusic/track"
	"github.com/RileySun/FyneMusic/utils"
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
	RefreshTray func()
	
	artContainer *fyne.Container
	title *widget.Label
	slider *track.TrackSlider
	sliderCurrent *widget.Label
	sliderEnd *widget.Label
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
						p.sliderCurrent.SetText(p.timeToString(p.Song.Current))
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
	//Player Label
	playerLabel := widget.NewLabel("Player")
	
	//Repeat and Shuffle
	p.shuffleButton = widget.NewButtonWithIcon("", utils.Icons.Shuffle, func() {p.ShuffleState()})
	p.repeatButton = widget.NewButtonWithIcon("", utils.Icons.Repeat, func() {p.RepeatState()})
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
	backButton := widget.NewButtonWithIcon("", utils.Icons.Menu, func() {p.ReturnToMenu()})
	backSpacer := layout.NewSpacer()
	backContainer := container.New(layout.NewHBoxLayout(), playerLabel, backSpacer, p.shuffleButton, p.repeatButton, backButton)
	
	//Top Container
	topBorder := canvas.NewLine(color.NRGBA{R: 155, G: 155, B: 155, A: 255})
	topBorder.StrokeWidth = 2
	topBorderContainer := container.NewBorder(nil, topBorder, nil, nil, backContainer)
	//Label works better as a spacer here than spacer does. weird.
	labelSpace := widget.NewLabel("")
	topContainer := container.New(layout.NewVBoxLayout(), topBorderContainer, labelSpace)
	
	//Meta
	artwork, title := p.GetMeta()
	p.artContainer = container.New(layout.NewMaxLayout(), artwork)
	
	//Title (if too long, truncate)
	titleLen := len([]rune(title))
	if titleLen > 55 {
		remove := titleLen - 55
		title = title[:len(title)-remove] + "..."
	}
	p.title = widget.NewLabel(title)
	p.title.Alignment = 1
	p.title.Wrapping = 1
	
	//Slider
	p.sliderCurrent = widget.NewLabel("0:00")
	p.slider = track.NewTrack(p.Song)
	p.sliderEnd = widget.NewLabel(p.timeToString(p.Song.Length))
	sliderBorder := container.NewBorder(nil, nil, p.sliderCurrent, p.sliderEnd, p.slider)
	p.sliderContainer = container.New(layout.NewMaxLayout(), sliderBorder)
	
	//Buttons
	prevButton := widget.NewButtonWithIcon("", utils.Icons.Prev, func() {p.Prev()})
	rewindButton := widget.NewButtonWithIcon("", utils.Icons.Rewind, func() {p.Rewind()})
	forwardButton := widget.NewButtonWithIcon("", utils.Icons.Forward, func() {p.Forward()})
	nextButton := widget.NewButtonWithIcon("", utils.Icons.Next, func() {p.Next()})
	p.PlayButton = playbutton.NewPlayButton(p.Song)
	buttonContainer := container.New(layout.NewAdaptiveGridLayout(5), prevButton, rewindButton, p.PlayButton, forwardButton, nextButton)
	
	//Bottom Container
	bottomContainer := container.New(layout.NewVBoxLayout(), p.title, p.sliderContainer, buttonContainer)
	
	//Output
	playerContainer := container.NewBorder(topContainer, bottomContainer, nil, nil, p.artContainer)
	
	return playerContainer
}

func (p *Player) RenderUpdate() {
	//Meta
	artwork, title := p.GetMeta()
	
	//Title (if too long, truncate) ((*widget.Label).TextWrap also truncates, not sure if I wanna do it that way.)
	titleLen := len([]rune(title))
	if titleLen > 55 {
		remove := titleLen - 55
		title = title[:len(title)-remove] + "..."
	}
	p.title.SetText(title)
	p.title.Wrapping = 1
	
	//Art
	p.artContainer.RemoveAll()
	p.artContainer.Add(artwork)
	artwork.Refresh()
	p.artContainer.Refresh()
	
	//Slider
	p.slider = track.NewTrack(p.Song)
	p.sliderCurrent = widget.NewLabel(p.timeToString(p.Song.Current))
	p.sliderEnd = widget.NewLabel(p.timeToString(p.Song.Length))
	sliderBorder := container.NewBorder(nil, nil, p.sliderCurrent, p.sliderEnd, p.slider)
	
	//Refresh
	p.sliderContainer.RemoveAll()
	p.sliderContainer.Add(sliderBorder)	
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
	prevButton := widget.NewButtonWithIcon("", utils.Icons.Prev, func() {p.Prev()})
	rewindButton := widget.NewButtonWithIcon("", utils.Icons.Rewind, func() {p.Rewind()})
	forwardButton := widget.NewButtonWithIcon("", utils.Icons.Forward, func() {p.Forward()})
	nextButton := widget.NewButtonWithIcon("", utils.Icons.Next, func() {p.Next()})
	
	p.PlayButton = playbutton.NewPlayButton(p.Song)
	buttonContainer := container.New(layout.NewAdaptiveGridLayout(5), prevButton, rewindButton, p.PlayButton, forwardButton, nextButton)
	//buttonCenter := container.New(layout.NewMaxLayout(), buttonContainer)
	
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

func (p *Player) UpdateWidgets() {
	//Slider and PlayButton
	p.slider.SetTime()
	p.PlayButton.UpdateState()
	p.sliderCurrent.SetText(p.timeToString(p.Song.Current))
	p.sliderEnd.SetText(p.timeToString(p.Song.Length))
}

//Queue
func (p *Player) RenderQueue() {//*fyneContainer {
	
}

	//Utils
//Player
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
	art.FillMode = canvas.ImageFillStretch
	art.Resize(fyne.NewSize(435, 435))

	//Title
	titleString := p.Song.Meta.Title + " - " + p.Song.Meta.Artist
	
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

func (p *Player) timeToString(time int64) string {
	min := int(time/60)
	output := strconv.Itoa(min) + ":"
	
	sec := math.Mod(float64(time), 60)
	if sec < 10 {
		output += "0" + strconv.Itoa(int(sec))
	} else {
		output += strconv.FormatFloat(sec, 'f', 0, 64)
	}
	
	return output
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
	p.RefreshTray()
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