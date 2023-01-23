package playlist

import(
	"strings"
	"sort"
)

//Title
type TitleSort []*SongItem

func (s TitleSort) Len() int {
	return len(s)
}

func (s TitleSort) Less(i, j int) bool {
	var si string = s[i].Meta.Title
    var sj string = s[j].Meta.Title
    var siL = strings.ToLower(si)
    var sjL = strings.ToLower(sj)
    if siL == sjL {
        return si < sj
    }
    return siL < sjL
}

func (s TitleSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (p *Playlist) SortByTitle(asc bool) []*SongItem {
	songList := make([]*SongItem, len(p.Songs))
	copy(songList, p.Songs)
	
	if asc {
		sort.Sort(sort.Reverse(TitleSort(songList)))
	} else {
		sort.Sort(TitleSort(songList))
	}
	
	return songList
}


//Artist
type ArtistSort []*SongItem

func (s ArtistSort) Len() int {
	return len(s)
}

func (s ArtistSort) Less(i, j int) bool {
	var si string = s[i].Meta.Artist
    var sj string = s[j].Meta.Artist
    var siL = strings.ToLower(si)
    var sjL = strings.ToLower(sj)
    if siL == sjL {
        return si < sj
    }
    return siL < sjL
}

func (s ArtistSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (p *Playlist) SortByArtist(asc bool) []*SongItem {
	songList := make([]*SongItem, len(p.Songs))
	copy(songList, p.Songs)
	
	if asc {
		sort.Sort(sort.Reverse(ArtistSort(songList)))
	} else {
		sort.Sort(ArtistSort(songList))
	}
	
	return songList
}

//Album
type AlbumSort []*SongItem

func (s AlbumSort) Len() int {
	return len(s)
}

func (s AlbumSort) Less(i, j int) bool {
	var si string = s[i].Meta.Album
    var sj string = s[j].Meta.Album
    var siL = strings.ToLower(si)
    var sjL = strings.ToLower(sj)
    if siL == sjL {
        return si < sj
    }
    return siL < sjL
}

func (s AlbumSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (p *Playlist) SortByAlbum(asc bool) []*SongItem {
	songList := make([]*SongItem, len(p.Songs))
	copy(songList, p.Songs)
	
	if asc {
		sort.Sort(sort.Reverse(AlbumSort(songList)))
	} else {
		sort.Sort(AlbumSort(songList))
	}
	
	return songList
}

//Year
type YearSort []*SongItem

func (s YearSort) Len() int {
	return len(s)
}

func (s YearSort) Less(i, j int) bool {
    if s[i].Meta.Year == s[j].Meta.Year {
        return s[i].Meta.Year < s[j].Meta.Year
    }
    return s[i].Meta.Year < s[j].Meta.Year
}

func (s YearSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (p *Playlist) SortByYear(asc bool) []*SongItem {
	songList := make([]*SongItem, len(p.Songs))
	copy(songList, p.Songs)
	
	if asc {
		sort.Sort(sort.Reverse(YearSort(songList)))
	} else {
		sort.Sort(YearSort(songList))
	}
	
	return songList
}