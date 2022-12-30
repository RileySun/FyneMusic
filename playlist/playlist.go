package playlist

import(
	"os"
	"strings"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

//Struct
type Playlist struct {
	Songs []string
	Index int64
	Length int64 //-1 for 0 index, 0 if only one song
	Select func(int64)
}

//Create
func NewPlaylist(dirPath string) *Playlist {
	playlist := new(Playlist)
	
	playlist.Songs = getSongs(dirPath)
	playlist.Index = 0
	playlist.Length = int64(len(playlist.Songs)) - 1 //-1 for 0 index
	
	return playlist
}

//Util
func getSongs(dirPath string) []string {
	var songList []string
	
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
	
	//Check is mp3 file and not a dir
	for _, f := range fileList {
		name := f.Name()
		split := strings.Split(name, ".")
		ext := split[len(split) - 1]//-1 for 0 index
		
		if !f.IsDir() && ext == "mp3" {
			path := dirPath + "/" + name
			songList = append(songList, path)
		}
    }
    
    return songList
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
			split := strings.Split(p.Songs[i], "/")
			name := split[len(split) - 1]//-1 for 0 index
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