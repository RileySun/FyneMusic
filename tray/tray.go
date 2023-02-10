package tray

import(
	"log"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	
	"github.com/RileySun/FyneMusic/player"
	"github.com/RileySun/FyneMusic/utils"
)

//Struct
type Tray struct {
	Desktop desktop.App
	App fyne.App
	Player *player.Player
	Menu *fyne.Menu
	ButtonMenu *fyne.Menu
	
	Current *fyne.MenuItem
	PlayItem *fyne.MenuItem
}

//Create

func NewTray(app fyne.App, playerObj *player.Player) *Tray {
	tray := new(Tray)
	tray.App = app
	tray.Player = playerObj
	tray.NewDesktopTray()
	return tray
}

func (t *Tray) NewDesktopTray() {
	var deskOk bool
	t.Desktop, deskOk = t.App.(desktop.App)
	
	if !deskOk {
		log.Print("tray: Desktop Tray Can Not Init")
	}
	
	//Menu Items
	t.Current = fyne.NewMenuItem("No Song", func() {})
	
	t.PlayItem = fyne.NewMenuItem("", func() {t.Play()})
	t.PlayItem.Icon = utils.Icons.Play
	
	prev := fyne.NewMenuItem("", func() {t.Prev()})
	prev.Icon = utils.Icons.Prev
	
	next := fyne.NewMenuItem("", func() {t.Next()})
	next.Icon = utils.Icons.Next
	
	
	t.Menu = fyne.NewMenu("Fyne Music", t.Current, prev, t.PlayItem, next)
	t.Desktop.SetSystemTrayIcon(utils.Logo())
	//t.Desktop.SetSystemTrayMenu(t.Menu)
}

//Utils
func (t *Tray) Refresh() {
	//Change label to current Song
	t.Current.Label = t.Player.Song.Meta.Title + " - " + t.Player.Song.Meta.Artist
	
	//Set play/pause state icon
	if t.Player.Song.Paused {
		t.PlayItem.Icon = utils.Icons.Play
	} else {
		t.PlayItem.Icon = utils.Icons.Pause
	}
	
	//Change Icon to match song image
	if len(t.Player.Song.Meta.Image) != 0 {
		icon := fyne.NewStaticResource("current", t.Player.Song.Meta.Image)
		t.Desktop.SetSystemTrayIcon(icon)
	} else {
		t.Desktop.SetSystemTrayIcon(utils.Credit())
	}
	
	//Update Everything
	t.Player.UpdateWidgets()
	t.Menu.Refresh()
}

//Actions

func (t *Tray) Play() {
	if t.Player.Song == nil {
		return
	}
	
	if t.Player.Song.Paused {
		t.Player.Song.Play()
	} else {
		t.Player.Song.Pause()
	}
	t.Refresh()
}

func (t *Tray) Prev() {
	t.Player.GetQueuePrev()
	t.Refresh()
}

func (t *Tray) Next() {
	t.Player.GetQueueNext()
	t.Refresh()
}

