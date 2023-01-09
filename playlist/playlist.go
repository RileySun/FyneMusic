package playlist

import(
	"os"
	"strings"
	"image/color"
	"sort"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	
	"github.com/RileySun/FynePod/meta"
	"github.com/RileySun/FynePod/player"
)

//Struct
type Playlist struct {
	Songs []*SongItem
	Index int64
	Length int64 //-1 for 0 index, 0 if only one song
	Select func(int64)
	Settings func()
	Player *player.Player
}

type SongItem struct {
	Path string
	Meta *meta.Meta
}

//Create
func NewPlaylist(dirPath string, playerObj *player.Player) *Playlist {
	playlist := new(Playlist)
	
	playlist.Player = playerObj
	
	playlist.Songs = getSongs(dirPath)
	playlist.Index = 0
	playlist.Length = int64(len(playlist.Songs)) - 1 //-1 for 0 index
	
	return playlist
}

//Util
func getSongs(dirPath string) ([]*SongItem) {
	//Output Slice
	var songItems []*SongItem
	 
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
			//New Song Item 
			songItem := new(SongItem)
			
			//Append 
			path := dirPath + "/" + name
			songItem.Path = path
			
			//Get meta for playlist tabs
			metaData := new(meta.Meta)
			metaData = meta.Get(dirPath + "/" + name)
			songItem.Meta = metaData
			
			songItems = append(songItems, songItem)
		}
    }
    
    return songItems
}

func (p *Playlist) PlaylistPaths() []string {
	var paths []string
	for _, song := range p.Songs {
		paths = append(paths, song.Path)
	}
	return paths
}

func (p *Playlist) sortDesc(by string) {
	switch by {
		case "Title":
			sort.Slice(p.Songs, func(a, b int) bool {
				return p.Songs[a].Meta.Title < p.Songs[b].Meta.Title
			})
			break
		case "Artist":
			sort.Slice(p.Songs, func(a, b int) bool {
				return p.Songs[a].Meta.Artist < p.Songs[b].Meta.Artist
			})
			break
		case "Album":
			sort.Slice(p.Songs, func(a, b int) bool {
				return p.Songs[a].Meta.Album < p.Songs[b].Meta.Album
			})
			break
		case "Year":
			sort.Slice(p.Songs, func(a, b int) bool {
				return p.Songs[a].Meta.Year < p.Songs[b].Meta.Year
			})
			break
		default:
	}
}//by = Artist, Title

func (p *Playlist) sortAsc(by string) {
	switch by {
		case "Title":
			sort.Slice(p.Songs, func(a, b int) bool {
				return p.Songs[a].Meta.Title > p.Songs[b].Meta.Title
			})
			break
		case "Artist":
			sort.Slice(p.Songs, func(a, b int) bool {
				return p.Songs[a].Meta.Artist > p.Songs[b].Meta.Artist
			})
			break
		case "Album":
			sort.Slice(p.Songs, func(a, b int) bool {
				return p.Songs[a].Meta.Album > p.Songs[b].Meta.Album
			})
			break
		case "Year":
			sort.Slice(p.Songs, func(a, b int) bool {
				return p.Songs[a].Meta.Year > p.Songs[b].Meta.Year
			})
			break
		default:
	}
}//by = Artist, Title

//Render
func (p *Playlist) Render() *fyne.Container {
	//Top Area (Logo, and Settings Button)
	logo := widget.NewLabel("Fyne Pod")
	space := layout.NewSpacer()
	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {p.Settings()})
	topHBox := container.New(layout.NewHBoxLayout(), logo, space, settingsButton)
	topBorder := canvas.NewLine(color.NRGBA{R: 155, G: 155, B: 155, A: 255})
	topBorder.StrokeWidth = 2
	topContainer := container.NewBorder(nil, topBorder, nil, nil, topHBox)
	
	
	//Song List
	songList := p.RenderSong(false)
	songList.OnSelected = func(id widget.ListItemID) {p.OnSelect(id)}
	songListContainer := container.New(layout.NewMaxLayout(), songList)
	
	//Artist List
	artistList := p.RenderArtist(false)
	artistList.OnSelected = func(id widget.ListItemID) {p.OnSelect(id)}
	artistListContainer := container.New(layout.NewMaxLayout(), artistList)
	
	//Album List
	albumList := p.RenderAlbum(false)
	albumList.OnSelected = func(id widget.ListItemID) {p.OnSelect(id)}
	albumListContainer := container.New(layout.NewMaxLayout(), albumList)
	
	//Year List
	yearList := p.RenderYear(false)
	yearList.OnSelected = func(id widget.ListItemID) {p.OnSelect(id)}
	yearListContainer := container.New(layout.NewMaxLayout(), yearList)
	
	//Tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Title", songListContainer),
		container.NewTabItem("Artist", artistListContainer),
		container.NewTabItem("Album", albumListContainer),
		container.NewTabItem("Year", yearListContainer),
	)
	tabs.SetTabLocation(container.TabLocationTop)
	
	
	//Mini Player
	mini := p.Player.RenderMini()
	
	return container.NewBorder(topContainer, mini, nil, nil, tabs)
}

//Render Lists
func (p *Playlist) RenderSong(asc bool) *widget.List {
	if asc {
		p.sortAsc("Title")
	} else {
		p.sortDesc("Title")
	}
	
	list := widget.NewList(
		func() int {
			return len(p.Songs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			name := ""
			
			if p.Songs[i].Meta.Title == " " {
				name += p.Songs[i].Meta.File
			} else {
				name += p.Songs[i].Meta.Title
			}
			
			
			if p.Songs[i].Meta.Artist != "" {
				name += " - " + p.Songs[i].Meta.Artist
			} else {
				name += " - Unknown"
			}
			o.(*widget.Label).SetText(name)
		},
	)
	
	return list
}

func (p *Playlist) RenderArtist(asc bool) *widget.List {
	if asc {
		p.sortAsc("Artist")
	} else {
		p.sortDesc("Artist")
	}
	
	list := widget.NewList(
		func() int {
			return len(p.Songs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			name := ""
			//If artist tag exists, then use that, if not, write unknown
			if p.Songs[i].Meta.Artist != "" {
				name += p.Songs[i].Meta.Artist + " - "
			} else {
				name += "Unkown - "
			}
			
			if p.Songs[i].Meta.Title == " " {
				name += p.Songs[i].Meta.File
			} else {
				name += p.Songs[i].Meta.Title
			}
			
			o.(*widget.Label).SetText(name)
		},
	)
	
	return list
}

func (p *Playlist) RenderAlbum(asc bool) *widget.List {
	if asc {
		p.sortAsc("Album")
	} else {
		p.sortDesc("Album")
	}
	
	list := widget.NewList(
		func() int {
			return len(p.Songs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			name := ""
			//If artist tag exists, then use that, if not, write unknown
			if p.Songs[i].Meta.Album != "" {
				name += p.Songs[i].Meta.Album + " - "
			} else {
				name += "Unkown - "
			}
			
			if p.Songs[i].Meta.Title == " " {
				name += p.Songs[i].Meta.File
			} else {
				name += p.Songs[i].Meta.Title
			}	
			
			o.(*widget.Label).SetText(name)
		},
	)
	
	return list
}

func (p *Playlist) RenderYear(asc bool) *widget.List {
	if asc {
		p.sortAsc("Year")
	} else {
		p.sortDesc("Year")
	}
	
	list := widget.NewList(
		func() int {
			return len(p.Songs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			name := ""
			//If artist tag exists, then use that, if not, write unknown
			if p.Songs[i].Meta.Year != "" {
				name += p.Songs[i].Meta.Year + " - "
			} else {
				name += "Unkown - "
			}
			
			if p.Songs[i].Meta.Title == " " {
				name += p.Songs[i].Meta.File
			} else {
				name += p.Songs[i].Meta.Title
			}
			
			o.(*widget.Label).SetText(name)
		},
	)
	
	return list
}


//Actions
func (p *Playlist) OnSelect(id widget.ListItemID) {
	index := int64(id)
	p.Select(index)
}