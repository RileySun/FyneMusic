package song

import(
	"os"
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

func OpenSong(filepath string) Song {
	var newSong Song
	
	var fileErr error

	//Open File	
	newSong.file, fileErr = os.Open(filepath)
	
	if fileErr != nil {
		panic("Can't open file. Error: " + fileErr.Error())
	}
	
	//defer newSong.file.Close()
	
	//Mp3-Decoder
	decoder, decoderErr := mp3.NewDecoder(newSong.file)
	if decoderErr != nil {
		panic("Can't decode file. Error: " + decoderErr.Error())
	}
	
	newSong.decoder = decoder
	splitName := strings.Split(filepath, "/")
	newSong.Name = splitName[len(splitName) - 1]
	newSong.Length = newSong.decoder.Length()
	newSong.Paused = true
	
	
	context, ready, contextError := oto.NewContext(newSong.decoder.SampleRate(), 2, 2)
	if contextError != nil {
		panic("Can't decode file. Error: " + contextError.Error())
	}
	
	<-ready
	
	newSong.Player = context.NewPlayer(newSong.decoder)
	//defer newSong.Player.Close()
	
	return newSong
}