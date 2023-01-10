package playbutton

import(
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	
	"github.com/RileySun/FyneMusic/song"
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
	button.Icon = theme.MediaPlayIcon()
	button.Text = "Play"
	button.OnTapped = func() {button.ChangeState()}
	
	return button
}

//Actions
func (b *PlayButton) ChangeState() {
	if b.Song != nil {
		if b.Song.Paused {
			b.Song.Play()
			b.Icon = theme.MediaPauseIcon()
			b.Text = "Pause"
			b.Refresh()
		} else {
			b.Song.Pause()
			b.Icon = theme.MediaPlayIcon()
			b.Text = "Play"
			b.Refresh()
		}
	}
}

func (b *PlayButton) UpdateState() {
	if b.Song.Paused {
		b.Icon = theme.MediaPlayIcon()
		b.Text = "Play"
		b.Refresh()
	} else {
		b.Icon = theme.MediaPauseIcon()
		b.Text = "Pause"
		b.Refresh()
	}
}