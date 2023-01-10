package track

import(
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/data/binding"
	
	"time"
	"strconv"
	
	"github.com/RileySun/FyneMusic/song"
)

//Struct with extra functions
type TrackSlider struct {
	widget.Slider
	
	Song *song.Song
	
	TimeUpdateDone chan bool
	TimeStr binding.String
	MaxTimeStr binding.String
	SliderFloat binding.Float
}

//Channels
var TimeUpdateDone chan bool

//Declarations
var slider *TrackSlider

//Global

//Create
func NewTrack(current *song.Song) *TrackSlider {
	//New TrackSlider
	slider = NewTrackSlider(0, float64(current.Length))
	
	//Set Song
	slider.Song = current
	
	//Channel
	slider.TimeUpdateDone = make(chan bool, 100)
	
	//Current Time
	slider.TimeStr = binding.NewString()
	slider.TimeStr.Set(strconv.Itoa(int(slider.Song.Current)))
	
	//Max Time
	slider.MaxTimeStr = binding.NewString()
	slider.MaxTimeStr.Set(strconv.Itoa(int(slider.Song.Length)))
	
	//Float Current Time
	slider.SliderFloat = binding.NewFloat()
	slider.SliderFloat.Set(float64(slider.Song.Current))
	
	slider.Bind(slider.SliderFloat)
	
	slider.UpdateTime()
	
	return slider
}

func NewTrackSlider(min float64, max float64) *TrackSlider {
	track := &TrackSlider{}
	track.ExtendBaseWidget(track)
	track.Min = min
	track.Max = max
	track.Step = 1

	return track
}

//Util
func (t *TrackSlider) UpdateTime() {
	go func() {
		for {
			select {
				case <- t.TimeUpdateDone:
					return
				default:
					t.SetTime()
					time.Sleep(time.Second/2)
	   		 }
		}
	}()
} //runs SetTime every 0.5seconds until TimeUpdateDone is true

func (t *TrackSlider) SetTime() {
	t.TimeStr.Set(strconv.Itoa(int(t.Song.Current)))
	t.SliderFloat.Set(float64(t.Song.Current))
}

func  (t *TrackSlider) Close() {
	t.TimeUpdateDone <- true
}

//Extra functions
func (t *TrackSlider) Dragged(e *fyne.DragEvent) {
	wasPlaying := t.Song.Player.IsPlaying()
	if wasPlaying {
		t.TimeUpdateDone <- true
		t.Song.Pause()
	}
	
	t.Slider.Dragged(e)
	
	t.Song.Seek(int64(t.Value))
	t.Song.Current = int64(t.Value)
	t.SliderFloat.Set(float64(t.Value))
	
	if wasPlaying {
		t.Song.Play()
		t.UpdateTime()
	}
}

func (t *TrackSlider) DragEnd() {
	
	//t.UpdateTime()
}