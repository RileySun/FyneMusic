package playlist

import(
	"path/filepath"
	"io/fs"
	"strings"
	"image/color"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	
	"github.com/RileySun/FyneMusic/meta"
	"github.com/RileySun/FyneMusic/player"
	"github.com/RileySun/FyneMusic/utils"
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
	
	//Search & Sort Stuff
	Original []*SongItem
	SongSort []*SongItem
	ArtistSort []*SongItem
	AlbumSort []*SongItem
	YearSort []*SongItem
	
	//Fyne Objects (change state out of scope)
	TopFinal *fyne.Container
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
	playlist.Songs = GetSongs(dirPath)
	playlist.Original = playlist.Songs
	playlist.Index = 0
	playlist.Length = int64(len(playlist.Songs)) - 1 //-1 for 0 index
	return playlist
}

//Util
func GetSongs(dirpath string) []*SongItem {
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

//Render
func (p *Playlist) Render() *fyne.Container {
	//Top Area (Logo, and Settings Button)
	logo := widget.NewLabel("Playlist")
	space := layout.NewSpacer()
	settingsButton := widget.NewButtonWithIcon("", utils.Icons.Settings, func() {p.Settings()})
	p.SearchButton = widget.NewButtonWithIcon("", utils.Icons.Search, func() {p.openSearch()})
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
		container.NewTabItem("Song", songListContainer),
		container.NewTabItem("Artist", artistListContainer),
		container.NewTabItem("Album", albumListContainer),
		container.NewTabItem("Year", yearListContainer),
	)
	tabs.SetTabLocation(container.TabLocationTop)
	tabs.OnSelected = func(t *container.TabItem) {
		switch t.Text {
			case "Song":
				p.Songs = p.SongSort
				break
			case "Artist":
				p.Songs = p.ArtistSort
				break
			case "Album":
				p.Songs = p.AlbumSort
				break
			case "Year":
				p.Songs = p.YearSort
				break
			default:
		}
	}
	
	//Set first
	p.Songs = p.SongSort
	
	return tabs
}

//Render Lists
func (p *Playlist) RenderSong(asc bool) *widget.List {
	p.SongSort = p.SortByTitle(asc)
	
	list := widget.NewList(
		func() int {
			return len(p.SongSort)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			name := p.SongSort[i].Meta.Title + " - " + p.SongSort[i].Meta.Artist
			o.(*widget.Label).SetText(name)
		},
	)
	
	return list
}

func (p *Playlist) RenderArtist(asc bool) *widget.List {
	p.ArtistSort = p.SortByArtist(asc)
	
	list := widget.NewList(
		func() int {
			return len(p.ArtistSort)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			name := p.ArtistSort[i].Meta.Artist + " - " + p.ArtistSort[i].Meta.Title
			o.(*widget.Label).SetText(name)
		},
	)
	
	return list
}

func (p *Playlist) RenderAlbum(asc bool) *widget.List {
	p.AlbumSort = p.SortByAlbum(asc)
	
	list := widget.NewList(
		func() int {
			return len(p.AlbumSort)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			name := p.AlbumSort[i].Meta.Album + " - " + p.AlbumSort[i].Meta.Title			
			o.(*widget.Label).SetText(name)
		},
	)
	
	return list
}

func (p *Playlist) RenderYear(asc bool) *widget.List {
	p.YearSort = p.SortByYear(asc)
	
	list := widget.NewList(
		func() int {
			return len(p.Songs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			name := p.YearSort[i].Meta.Year + " - " + p.YearSort[i].Meta.Title
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
	//Button Action
	p.SearchButton.OnTapped = p.closeSearch
	//New Entry
	searchEntry := widget.NewEntry()
	searchEntry.OnChanged = func(s string) {p.search(s)}
	searchEntry.SetPlaceHolder("Search")
	//New Clear Button
	clearSearch := widget.NewButtonWithIcon("", utils.Icons.Cancel, func() {p.closeSearch()})
	searchContainer := container.NewBorder(nil, nil, nil, clearSearch, searchEntry)
	//Add to View
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