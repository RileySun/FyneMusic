package playlist

import(
	"path/filepath"
	"io/fs"
	"strings"
	"image/color"
	"sort"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	
	"github.com/RileySun/FyneMusic/meta"
	"github.com/RileySun/FyneMusic/player"
)

//Struct
type Playlist struct {
	Songs []*SongItem
	Index int64
	Length int64 //-1 for 0 index, 0 if only one song
	Select func(int64)
	Settings func()
	Player *player.Player
	TabContainer *fyne.Container
	
	//Search Stuff
	Original []*SongItem
	TopFinal *fyne.Container
	SearchEntry *widget.Entry
	SearchButton *widget.Button
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
	playlist.Original = playlist.Songs
	playlist.Index = 0
	playlist.Length = int64(len(playlist.Songs)) - 1 //-1 for 0 index
	
	return playlist
}

//Util
func getSongs(dirpath string) []*SongItem {
	var songItems []*SongItem

	walkErr := filepath.Walk(dirpath, func(path string, info fs.FileInfo, err error) error {
		name := info.Name()
		split := strings.Split(name, ".")
		ext := split[len(split) - 1]//-1 for 0 index
		
		if !info.IsDir() && ext == "mp3" {
			songItem := new(SongItem)
			
			//Append 
			songItem.Path = path
			
			//Get meta for playlist tabs
			metaData := new(meta.Meta)
			metaData = meta.Get(path)
			songItem.Meta = metaData
			
			songItems = append(songItems, songItem)
		}
		
		return nil
	})
	
	if walkErr != nil {
		panic("Playlist: Rescursive Song Get Err:" + walkErr.Error())
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
	logo := widget.NewLabel("Fyne Music")
	space := layout.NewSpacer()
	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {p.Settings()})
	p.SearchButton = widget.NewButtonWithIcon("", theme.SearchIcon(), func() {p.openSearch()})
	topHBox := container.New(layout.NewHBoxLayout(), logo, space, p.SearchButton, settingsButton)
	topBorder := canvas.NewLine(color.NRGBA{R: 155, G: 155, B: 155, A: 255})
	topBorder.StrokeWidth = 2
	topContainer := container.NewBorder(nil, topBorder, nil, nil, topHBox)
	p.TopFinal = container.New(layout.NewVBoxLayout(), topContainer)
	
	//Song Sort Tabs
	tabs := p.RenderTabs()
	p.TabContainer = container.New(layout.NewMaxLayout(), tabs)
	
	//Mini Player
	mini := p.Player.RenderMini()
	
	return container.NewBorder(p.TopFinal, mini, nil, nil, p.TabContainer)
}

func (p *Playlist) RenderTabs() *container.AppTabs {	
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
	
	return tabs
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
	p.closeSearch()
}

func (p *Playlist) openSearch() {
	p.SearchButton.OnTapped = p.closeSearch
	p.SearchEntry = widget.NewEntry()
	p.SearchEntry.OnChanged = func(s string) {p.search(s)}
	p.SearchEntry.SetPlaceHolder("Search")
	searchContainer := container.New(layout.NewMaxLayout(), p.SearchEntry)
	p.TopFinal.Add(searchContainer)
	p.TopFinal.Refresh()
}

func (p *Playlist) closeSearch() {
	//if search is open
	if (len(p.TopFinal.Objects) > 1) {
		p.SearchButton.OnTapped = p.openSearch
		last := p.TopFinal.Objects[len(p.TopFinal.Objects) - 1]//get last object in container
		p.TopFinal.Remove(last)
		p.TopFinal.Refresh()
	}
	
	//Refresh List
	p.Songs = p.Original
	tabs := p.RenderTabs()
	p.TabContainer.RemoveAll()
	p.TabContainer.Add(tabs)
	p.TabContainer.Refresh()
}

func (p *Playlist) search(searchString string) {	
	if searchString == "" {
		p.Songs = p.Original
		tabs := p.RenderTabs()
		p.TabContainer.RemoveAll()
		p.TabContainer.Add(tabs)
		p.TabContainer.Refresh()
		return
	}
	
	var searchSlice []*SongItem
	lowerSearch := strings.ToLower(searchString) //Lower case to make search case insensitive
	for _, song := range p.Original {
		file := strings.Contains(strings.ToLower(song.Meta.File), lowerSearch)
		title := strings.Contains(strings.ToLower(song.Meta.Title), lowerSearch)
		artist := strings.Contains(strings.ToLower(song.Meta.Artist), lowerSearch)
		album := strings.Contains(strings.ToLower(song.Meta.Album), lowerSearch)
		year := strings.Contains(strings.ToLower(song.Meta.Year), lowerSearch)
		
		if file || title || artist || album || year {
			searchSlice = append(searchSlice, song)
		}
	}
	
	p.Songs = searchSlice
	tabs := p.RenderTabs()
	p.TabContainer.RemoveAll()
	p.TabContainer.Add(tabs)
	p.TabContainer.Refresh()
}