package playbutton

import(
	"fyne.io/fyne/v2/widget"
	
	"github.com/RileySun/FyneMusic/song"
	"github.com/RileySun/FyneMusic/utils"
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
	button.Icon = utils.Icons.Play
	button.Text = ""
	button.OnTapped = func() {button.ChangeState()}
	
	return button
}

//Actions
func (b *PlayButton) ChangeState() {
	if b.Song != nil {
		if b.Song.Paused {
			b.Song.Play()
			b.Icon = utils.Icons.Pause
			b.Refresh()
		} else {
			b.Song.Pause()
			b.Icon = utils.Icons.Play
			b.Refresh()
		}
	}
}

func (b *PlayButton) UpdateState() {
	if b.Song.Paused {
		b.Icon = utils.Icons.Play
		b.Refresh()
	} else {
		b.Icon = utils.Icons.Pause
		b.Refresh()
	}
}