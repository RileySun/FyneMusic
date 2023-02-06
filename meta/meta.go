package meta

import(
	"os"
	"strconv"
	
	"github.com/dhowden/tag"
	"github.com/RileySun/FyneMusic/utils"
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

//Get Meta Data using filepath
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
	
	if int(m.Year()) > 1000 && int(m.Year()) < 9999 {
		meta.Year = strconv.Itoa(m.Year())
	} else {
		meta.Year = "Unknown"
	}
	
	//Art
	if m.Picture() != nil {
		meta.Image = m.Picture().Data
	} else {
		img := utils.Logo()
		meta.Image = img.StaticContent
	}
	
	return meta
}