package song

import(
	"io"
	"os"
	"log"
	"time"
	"strings"
	
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"github.com/RileySun/FyneMusic/meta"
	"github.com/RileySun/FyneMusic/settings"
)

type Song struct {
	file *os.File			//File refrence
	decoder *mp3.Decoder	//MP3 Decoder
	Player oto.Player		//Oto Player
	Path string				//Filepath
	Name string				//File Name
	Length int64			//Seek Length
	Current int64			//Current Seek
	Paused bool				//Is paused or not
	Meta *meta.Meta			//All Available Metadata
}

var stopUpdating chan bool

//This is super important when closing a song (see Close())
var context *oto.Context

//Constructor

func NewSong(filepath string) *Song {
	//Get Congif
	config := settings.GetSettings()
	
	//Song pointer
	currentSong := new(Song)
	
	//Current Track Seek Update Channel
	stopUpdating = make(chan bool, 100)
	
	//Set file Path
	currentSong.Path = filepath
	
	//Open File	
	currentSong.file = OpenSongFile(filepath)
	
	//Get Metadata
	currentSong.Meta = meta.Get(filepath)
	
	//Mp3-Decoder
	decoder, decoderErr := mp3.NewDecoder(currentSong.file)
	if decoderErr != nil {
		log.Fatal("song: Can't decode file. Error: " + decoderErr.Error())
	}
	
	currentSong.decoder = decoder
	splitName := strings.Split(filepath, "/")
	currentSong.Name = splitName[len(splitName) - 1]
	
	currentSong.Paused = true
	currentSong.Current = 0
	
	//Get Song Duration
	currentSong.Length = currentSong.IOtoSec(currentSong.decoder.Length())
	
	//NewContext
	var ready chan struct{}
	var contextError error
	context, ready, contextError = oto.NewContext(currentSong.decoder.SampleRate(), 2, 2)
	if contextError != nil {
		log.Fatal("song: Can't create context. Error: " + contextError.Error())
	}
	
	<-ready //context ready channel
	
	//New Player
	currentSong.Player = context.NewPlayer(currentSong.decoder)
	
	//In case config is missing
	if config != nil {
		currentSong.Player.SetVolume(config.Volume)
	} else {
		currentSong.Player.SetVolume(0.5)
	}
	
	return currentSong
}

//Utils
func OpenSongFile(filepath string) *os.File {
	file, fileErr := os.Open(filepath)
	
	if fileErr != nil {
		log.Fatal("song: Can't open file. Error: " + fileErr.Error())
	}
	
	return file
}

func (s *Song) IsEnded() bool {
	//Cant be over if still playing
	if s.Player.IsPlaying() {
		return false
	}
	
	return s.Current >= s.Length
}

//Start updating seek
func (s *Song) startUpdate() {
	go func() {
		for {
			select {
				case <- stopUpdating:
					return
				default:
					if !s.Player.IsPlaying() {
						return
					}
					s.Current += 1
					time.Sleep(time.Second)
	   		 }
		}
	}()
}

func (s *Song) endUpdate() {
	stopUpdating <- true
}

//Actions
func (s *Song) Restart() {
	_, _ = s.Player.(io.Seeker).Seek(0, io.SeekStart)
	s.Current = 0
	if !s.Player.IsPlaying() {
		s.Play()
	}
}

func (s *Song) Rewind() {
	//If was playing, then pause
	wasPaused := s.Player.IsPlaying()
	if wasPaused {
		s.Pause()
	}
	
	newSeek := 10 * int64(s.decoder.SampleRate()) * 4
	
	if s.Current > 10 {
		_, _ = s.Player.(io.Seeker).Seek(newSeek, io.SeekCurrent)
		s.Current -= 10
	} else {
		_, _ = s.Player.(io.Seeker).Seek(0, io.SeekStart)
		s.Current = 0
	}
	
	
	//If was paused, continue playing
	if wasPaused {
		s.Play()
	}
}

func (s *Song) Play() {
	s.Player.Play()
	s.startUpdate()
	s.Paused = false
}

func (s *Song) Pause() {
	s.Player.Pause()
	s.endUpdate()
	s.Paused = true
}

func (s *Song) Seek(newSeek int64) {
	s.endUpdate()
	newIO := newSeek * int64(s.decoder.SampleRate()) * 4 //Sample size is 4
	_, _ = s.Player.(io.Seeker).Seek(newIO, io.SeekStart)
	s.Current = newSeek
	s.startUpdate()
}

func (s *Song) Forward() {
	//If was paused, then pause
	wasPaused := s.Player.IsPlaying()
	if wasPaused {
		s.Pause()
	}
	
	newSeek := 10 * int64(s.decoder.SampleRate()) * 4
	
	if s.Current < s.Length - 10 {
		_, _ = s.Player.(io.Seeker).Seek(newSeek, io.SeekCurrent)
		s.Current += 10
	} else {
		_, _ = s.Player.(io.Seeker).Seek(0, io.SeekEnd)
		s.Current = s.Length
	}
	
	//If was paused, continue playing
	if wasPaused {
		s.Play()
	}
}
	
func (s *Song) Close() {
	playerError := s.Player.Close()
	if playerError != nil {
		log.Print("song: Player.Close failed: " + playerError.Error())
	}
	
	fileError := s.file.Close()
	if fileError != nil {
		log.Print("song: File.Close failed: " + fileError.Error())
	}
	
	/*
		If this is not here, then the context does not quit properly and only every other song is played see
	*/
	context.Suspend()
	
	close(stopUpdating)
}

//Utils
func (s *Song) SectoIO(seconds int64) int64 {
	const sampleSize = 4 
	return seconds * int64(s.decoder.SampleRate()) * sampleSize
}

func (s *Song) IOtoSec(ioInt int64) int64 {
	const sampleSize = 4 
	samples := ioInt / sampleSize
	return samples / int64(s.decoder.SampleRate())
}