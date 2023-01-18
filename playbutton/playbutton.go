package playbutton

import(
	"fyne.io/fyne/v2/widget"
	
	"github.com/RileySun/FyneMusic/song"
	"github.com/RileySun/FyneMusic/icons"
)

//Struct with extra functions
type PlayButton struct {
	widget.Button
	Song *song.Song
}

//Global

//Create
func NewPlayButton(current *song.Song) *PlayButton {
	button := &PlayButton{}
	button.ExtendBaseWidget(button)
	button.Song = current
	button.Icon = icons.Play
	button.Text = "Play "
	button.OnTapped = func() {button.ChangeState()}
	
	return button
}

//Actions
func (b *PlayButton) ChangeState() {
	if b.Song != nil {
		if b.Song.Paused {
			b.Song.Play()
			b.Icon = icons.Pause
			b.Text = "Pause"
			b.Refresh()
		} else {
			b.Song.Pause()
			b.Icon = icons.Play
			b.Text = "Play"
			b.Refresh()
		}
	}
}

func (b *PlayButton) UpdateState() {
	if b.Song.Paused {
		b.Icon = icons.Play
		b.Text = "Play"
		b.Refresh()
	} else {
		b.Icon = icons.Pause
		b.Text = "Pause"
		b.Refresh()
	}
}