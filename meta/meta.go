package meta

import(
	"os"
	
	"github.com/dhowden/tag"
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

func Get(file *os.File, path string) *Meta {
	meta := new(Meta)
	
	//File Stat
	s, fileErr := os.Stat(path)
	meta.File = s.Name()
	meta.Size = s.Size()
	if fileErr != nil {
		panic("meta:" + fileErr.Error())
	}
	
	//Get file tags
	m, tagErr := tag.ReadFrom(file)
	
	//If no tags, return stat info
	if tagErr != nil {
		return meta
	}
	
	//Else fill in tag meta
	meta.Title = m.Title()
	meta.Artist = m.Artist()
	meta.Image = m.Picture().Data
	
	return meta
}