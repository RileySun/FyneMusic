package playbutton

import(
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	
	"github.com/RileySun/FynePod/song"
)

//Struct with extra functions
type PlayButton struct {
	widget.Button
	podcast *song.Song
}

//Global

//Create
func NewPlayButton(current *song.Song) *PlayButton {
	button := &PlayButton{}
	button.ExtendBaseWidget(button)
	button.podcast = current
	button.Icon = theme.MediaPlayIcon()
	button.Text = "Play"
	button.OnTapped = func() {button.ChangeState()}
	
	return button
}

//Actions
func (b *PlayButton) ChangeState() {
	if b.podcast != nil {
		if b.podcast.Paused {
			b.podcast.Play()
			b.Icon = theme.MediaPauseIcon()
			b.Text = "Pause"
			b.Refresh()
		} else {
			b.podcast.Pause()
			b.Icon = theme.MediaPlayIcon()
			b.Text = "Play"
			b.Refresh()
		}
	}
}

func (b *PlayButton) UpdateState() {
	if b.podcast.Paused {
		b.Icon = theme.MediaPlayIcon()
		b.Text = "Play"
		b.Refresh()
	} else {
		b.Icon = theme.MediaPauseIcon()
		b.Text = "Pause"
		b.Refresh()
	}
}