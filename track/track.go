package track

import(
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/data/binding"
	
	"time"
	"strconv"
	
	"fmt"
	
	"github.com/RileySun/FynePod/song"
)

//Struct with extra functions
type TrackSlider struct {
	widget.Slider
}

//Channels
var TimeUpdateDone chan bool

//Declarations
var timeStr binding.String
var maxTimeStr binding.String
var sliderFloat binding.Float
var slider *TrackSlider

//Global
var podcast *song.Song

//Create
func NewTrack(current *song.Song) *TrackSlider {
	podcast = current
	
	timeStr = binding.NewString()
	timeStr.Set(strconv.Itoa(int(podcast.Current)))
	
	maxTimeStr = binding.NewString()
	maxTimeStr.Set(strconv.Itoa(int(podcast.Length)))
	
	sliderFloat = binding.NewFloat()
	sliderFloat.Set(float64(podcast.Current))
	
	TimeUpdateDone = make(chan bool, 100)
	
	//currentTime = widget.NewLabelWithData(timeStr)
	//maxTime = widget.NewLabelWithData(maxTimeStr)
	
	slider = NewTrackSliderWithData(0, float64(podcast.Length), sliderFloat)
	
	UpdateTime()
	
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

func NewTrackSliderWithData(min float64, max float64, data binding.Float) *TrackSlider {
	track := &TrackSlider{}
	track.ExtendBaseWidget(track)
	track.Min = min
	track.Max = max
	track.Step = 1
	track.Bind(data)

	return track
}

//Util
func UpdateTime() {
	go func() {
		for {
			fmt.Println(TimeUpdateDone)
			select {
				case <- TimeUpdateDone:
					return
				default:
					timeStr.Set(strconv.Itoa(int(podcast.Current)))
					sliderFloat.Set(float64(podcast.Current))
					time.Sleep(time.Second)
	   		 }
		}
	}()
}

//Extra functions
func (s *TrackSlider) Dragged(e *fyne.DragEvent) {
	TimeUpdateDone <- true
	s.Slider.Dragged(e)
}

func (s *TrackSlider) DragEnd() {
	podcast.Seek(int64(s.Value))
	UpdateTime()
}