package player

import (
	"fmt"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	
	"github.com/hajimehoshi/oto/v2"
)

type Song struct {
	decoder *mp3.Decoder
	name string
	length int64
	current int64
	paused bool
}

func Render(podcast Song) *Container {
	prevButton =  = widget.NewButtonWithIcon("Prev", theme.MediaSkipPreviousIcon(), func() {prev()})
	rewindButton =  = widget.NewButtonWithIcon("Rewind", theme.MediaFastRewindIcon(), func() {rewind()})
	playButton = widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {play()})
	paused = widget.NewButtonWithIcon("Play", theme.MediaPauseIcon(), func() {pause()})
	forwardButton =  = widget.NewButtonWithIcon("Forward", theme. MediaFastForwardIcon(), func() {forward()})
	nextButton =  = widget.NewButtonWithIcon("Next", theme. MediaSkipNextIcon(), func() {next()})
	
	
}

func prev() {

}

func rewind() {

}

func play() {

}

func pause() {

}

func forward() {

}

func next() {

}

