package main

import (
	"os"
	"strings"
	"fmt"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	//"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/layout"
	
	
	"github.com/hajimehoshi/go-mp3"
	
	
	//"github.com/RileySun/FynePod/player"
)

type Song struct {
	decoder *mp3.Decoder
	name string
	length int64
	current int64
	paused bool
}

var podcast Song

func init() {
	fmt.Println("FynePod")
	openSong("Intro.mp3")
	
	fmt.Println(podcast)
}

func main() {
	app := app.New()
	window := app.NewWindow("FynePod")
	window.Resize(fyne.NewSize(400, 600))
	
	//window.SetContent(container.New(layout.NewCenterLayout(), content))
	window.CenterOnScreen()
	window.ShowAndRun()
}

func openSong(filename string) {
	//Open File
	dir, _ := os.Getwd()	
	file, fileErr := os.Open(dir + "/" + filename)
	
	if fileErr != nil {
		panic("Can't open file. Error: " + fileErr.Error())
	}
	
	defer file.Close()
	
	//Mp3-Decoder
	decoder, decoderErr := mp3.NewDecoder(file)
	if decoderErr != nil {
		panic("Can't decode file. Error: " + decoderErr.Error())
	}
	
	podcast.decoder = decoder
	podcast.name = strings.Split(filename, ".")[0]
	podcast.length = podcast.decoder.Length()
	podcast.paused = true
}