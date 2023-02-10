package track_test

import(
	"strconv"
	"testing"
	"github.com/stretchr/testify/assert"
	
	"github.com/RileySun/FyneMusic/song"
	"github.com/RileySun/FyneMusic/track"
)

var newSong *song.Song

func Test_NewTrack(t *testing.T) {
	newSong = song.NewSong("../Music/Intro.mp3")
	newTrack := track.NewTrack(newSong)
	
	timeStr, _ :=  newTrack.TimeStr.Get()
	maxTimeStr, _ :=  newTrack.TimeStr.Get()
	
	assert.True(t, newTrack.Song != nil)
	assert.True(t, timeStr == "0")
	assert.True(t, maxTimeStr != strconv.Itoa(int(newSong.Length)))
	
	defer newTrack.Close()
	defer newSong.Close()
}

func Test_SetTime(t *testing.T) {
	newSong = song.NewSong("../Music/Intro.mp3")
	newTrack := track.NewTrack(newSong)
	
	newSong.Seek(3)
	//Didnt use built in function cause it requires a fyne app to be running (binding.String)
	//Instead just uses first line of that function (t.SetTime())
	newTrack.TimeStr.Set(strconv.Itoa(int(newSong.Current)))
	
	timeStr, _ :=  newTrack.TimeStr.Get()
	
	assert.True(t, timeStr == strconv.Itoa(int(newSong.Current)))
	
	defer newTrack.Close()
	defer newSong.Close()
}