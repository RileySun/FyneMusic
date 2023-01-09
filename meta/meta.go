package meta

import(
	"os"
	"github.com/dhowden/tag"
	"io/ioutil"
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
		panic("Error MAN :" + tagErr.Error())
	}
	
	//Else fill in tag meta
	meta.Title = m.Title()
	meta.Artist = m.Artist()
	if m.Picture() != nil {
		meta.Image = m.Picture().Data
	} else {
		data, _ := ioutil.ReadFile("Default.json")
		meta.Image = data
	}
	
	return meta
}
