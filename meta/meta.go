package meta

import(
	"os"
	"github.com/dhowden/tag"
	"io/ioutil"
	"strconv"
)

type Meta struct {
	Title string
	Artist string
	Album string
	Year string
	Image []byte
	File string
	Size int64
}

func Get(path string) *Meta {
	meta := new(Meta)
	
	//File Stat
	s, statErr := os.Stat(path)
	meta.File = s.Name()
	meta.Size = s.Size()
	if statErr != nil {
		panic("meta:" + statErr.Error())
	}
	
	newFile, fileErr := os.Open(path)
	defer newFile.Close()
	if fileErr != nil {
		panic("meta:" + fileErr.Error())
	}

	//Get file tags
	m, tagErr := tag.ReadFrom(newFile)
	
	//If no tags, return stat info
	if tagErr != nil {
		//panic("Error TagErr:" + tagErr.Error())
		meta.Title = s.Name()
		return meta
	}
	
	//Else fill in tag meta
	
	//Title
	if m.Title() != " " && m.Title() != "" {
		meta.Title = m.Title()
	} else {
		meta.Title = meta.File
	}
	
	//Artist
	if m.Artist() != " " && m.Artist() != "" {
		meta.Artist = m.Artist()
	} else {
		meta.Artist = "Unknown"
	}
	
	//Album
	if m.Album() != " " && m.Album() != "" {
		meta.Album = m.Album()
	} else {
		meta.Album = "Unknown"
	}
	
	//Year
	
	if m.Year() > 1000 && m.Year() < 9999 {
		meta.Year = strconv.Itoa(m.Year())
	}
	
	//Art
	if m.Picture() != nil {
		meta.Image = m.Picture().Data
	} else {
		data, _ := ioutil.ReadFile("Default.json")
		meta.Image = data
	}
	
	return meta
}
