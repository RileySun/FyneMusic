package song

import(
	"io"
	"os"
	"time"
	"strings"
	
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

type Song struct {
	file *os.File
	decoder *mp3.Decoder
	Player oto.Player
	Name string
	Length int64
	Current int64
	Paused bool
}

//Constructor

func NewSong(filepath string) Song {
	var currentSong Song
	
	var fileErr error
	
	//Open File	
	currentSong.file, fileErr = os.Open(filepath)
	
	if fileErr != nil {
		panic("Can't open file. Error: " + fileErr.Error())
	}
	
	//
	
	//Mp3-Decoder
	decoder, decoderErr := mp3.NewDecoder(currentSong.file)
	if decoderErr != nil {
		panic("Can't decode file. Error: " + decoderErr.Error())
	}
	
	currentSong.decoder = decoder
	splitName := strings.Split(filepath, "/")
	currentSong.Name = splitName[len(splitName) - 1]
	currentSong.Length = currentSong.decoder.Length()
	currentSong.Paused = true
	
	
	context, ready, contextError := oto.NewContext(currentSong.decoder.SampleRate(), 2, 2)
	if contextError != nil {
		panic("Can't create context. Error: " + contextError.Error())
	}
	
	<-ready
	
	currentSong.Player = context.NewPlayer(currentSong.decoder)
	
	
	return currentSong
}

//Actions
func (s Song) Restart() {
	 _, _ = s.Player.(io.Seeker).Seek(0, io.SeekStart)
}

func (s Song) Rewind() {
	//I think this is 10 seconds reverse, no clue thou XD
	_, _ = s.Player.(io.Seeker).Seek(-1000000, io.SeekCurrent)
}

func (s Song) Play() {
	s.Player.Play()
	for s.Player.IsPlaying() {
        time.Sleep(time.Millisecond)
    } 
}

func (s Song) Forward() {
	_, _ = s.Player.(io.Seeker).Seek(1000000, io.SeekCurrent)
}
	
func (s Song) Close() {
	fileError := s.file.Close()
	if fileError != nil {
		panic("song.File.Close failed: " + fileError.Error())
	}
	
	playerError := s.Player.Close()
	if playerError != nil {
		panic("song.Player.Close failed: " + playerError.Error())
	}
}