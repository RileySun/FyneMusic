package playlist

import(
	"os"
	"strings"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	
	"github.com/RileySun/FynePod/meta"
)

//Struct
type Playlist struct {
	Songs []string
	Meta []*meta.Meta
	Index int64
	Length int64 //-1 for 0 index, 0 if only one song
	Select func(int64)
}

//Create
func NewPlaylist(dirPath string) *Playlist {
	playlist := new(Playlist)
	
	playlist.Songs, playlist.Meta = getSongs(dirPath)
	playlist.Index = 0
	playlist.Length = int64(len(playlist.Songs)) - 1 //-1 for 0 index
	
	return playlist
}

//Util
func getSongs(dirPath string) ([]string, []*meta.Meta) {
	var songList []string
	var metaList []*meta.Meta
	
	//Open dir
	dir, dirErr := os.Open(dirPath)
	if dirErr != nil {
		panic("Can't open music dir. Error: " + dirErr.Error())
	}
	
	//Get file names
	fileList, fileErr := dir.Readdir(0)//0 to get all files
	if fileErr != nil {
		panic("Can't open music dir. Error: " + fileErr.Error())
	}
	
	//Check extension, if is dir, and get Meta
	for _, f := range fileList {
		name := f.Name()
		split := strings.Split(name, ".")
		ext := split[len(split) - 1]//-1 for 0 index
		
		if !f.IsDir() && ext == "mp3" {
			//Append 
			path := dirPath + "/" + name
			songList = append(songList, path)
			
			//Get meta for playlist tabs
			metaData := new(meta.Meta)
			file, _ := os.Open(path)
			metaData = meta.Get(file, dirPath + "/" + name)
			metaList = append(metaList, metaData)
		}
    }
    
    return songList, metaList
}

//Actions
func (p *Playlist) Render() *fyne.Container {
	list := widget.NewList(
		func() int {
			return len(p.Songs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			name := p.Meta[i].Title
			if p.Meta[i].Artist != "" {
				name += " - " + p.Meta[i].Artist
			}
			o.(*widget.Label).SetText(name)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {p.OnSelect(id)}
	
	return container.New(layout.NewMaxLayout(), list)
}

func (p *Playlist) OnSelect(id widget.ListItemID) {
	index := int64(id)
	p.Select(index)
}